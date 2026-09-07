package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// @title gin-vue-blog API
// @version 1.0
// @description gin-vue-blog 后端接口文档。业务状态码 0 表示成功, 其他表示失败, HTTP 状态码统一为 200
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configPath := flag.String("c", "../config.yml", "配置文件路径")
	flag.Parse()

	// 根据命令行参数读取配置文件, 其他变量的初始化依赖于配置文件对象
	conf := g.ReadConfig(*configPath)

	_ = ginblog.InitLogger(conf)
	db := ginblog.InitDatabase(conf)
	rdb := ginblog.InitRedis(conf)

	/*
		点赞数 / 浏览数 / 站点访问量 / 访客地域只在 Redis 里累加, 数据库没有对应字段。
		启动时把备份里有、Redis 里没有的补回去(见 internal/counter.go),
		这样清过 Redis 或换机器部署之后, 文章的点赞和阅读量不会归零。
	*/
	if err := ginblog.RestoreCounters(context.Background(), rdb, db); err != nil {
		slog.Error("从备份恢复计数失败", "err", err)
	}

	// 初始化 gin 服务
	gin.SetMode(conf.Server.Mode)
	r := gin.New()
	// 只信任配置里的代理网段: 之前是 "*", 访客随手加个 X-Real-IP 就能伪造来源 IP
	if err := r.SetTrustedProxies(conf.TrustedProxies()); err != nil {
		slog.Error("可信代理网段配置有误", "trusted-proxies", conf.TrustedProxies(), "err", err)
		os.Exit(1)
	}
	// 开发模式使用 gin 自带的日志和恢复中间件, 生产模式使用自定义的中间件
	if conf.Server.Mode == "debug" {
		r.Use(gin.Logger(), gin.Recovery()) // gin 自带的日志和恢复中间件, 挺好用的
	} else {
		r.Use(middleware.Recovery(true), middleware.Logger())
	}
	r.Use(middleware.CORS())
	r.Use(middleware.WithGormDB(db))
	r.Use(middleware.WithRedisDB(rdb))
	r.Use(middleware.WithCookieStore(conf.Session.Name, conf.Session.Salt))
	ginblog.RegisterHandlers(r)

	// 使用本地文件上传, 需要静态文件服务, 使用七牛云不需要
	if conf.Upload.OssType == "local" {
		r.Static(conf.Upload.Path, conf.Upload.StorePath)
	}

	serverAddr := conf.Server.Port
	/*
		启动横幅不走日志: slog 的级别由配置控制 (部署用的 config.docker.yml 是 error),
		用 log/slog 输出会被直接丢掉, 容器起来后看不到监听地址。
	*/
	if serverAddr[0] == ':' || strings.HasPrefix(serverAddr, "0.0.0.0:") {
		fmt.Printf("Serving HTTP on (http://localhost:%s/) ... \n", strings.Split(serverAddr, ":")[1])
	} else {
		fmt.Printf("Serving HTTP on (http://%s/) ... \n", serverAddr)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 定期把 Redis 计数抄一份进数据库, 进程被 kill -9 时最多丢这一个周期的增量
	go flushCountersPeriodically(ctx, rdb, db)

	srv := &http.Server{Addr: serverAddr, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP 服务启动失败", "addr", serverAddr, "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("收到退出信号, 开始关闭")

	/*
		退出前再落一次库: 正常重启(docker restart / systemctl restart)走的都是 SIGTERM,
		不落这一次就会丢掉上个周期之后的所有点赞和浏览。
		用独立的 context: 上面那个已经被信号取消了, 拿它做 Redis 请求会立刻失败。
	*/
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ginblog.FlushCounters(shutdownCtx, rdb, db); err != nil {
		slog.Error("退出前落库计数失败", "err", err)
	}
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("关闭 HTTP 服务失败", "err", err)
		os.Exit(1)
	}
	slog.Info("已退出")
}

// 每 counterFlushInterval 落库一次
const counterFlushInterval = 10 * time.Minute

func flushCountersPeriodically(ctx context.Context, rdb *redis.Client, db *gorm.DB) {
	ticker := time.NewTicker(counterFlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 失败只告警: 备份写不进去不该影响正在跑的服务, 下个周期还会再试
			if err := ginblog.FlushCounters(ctx, rdb, db); err != nil {
				slog.Warn("定时落库计数失败", "err", err)
			}
		}
	}
}
