package main

import (
	"context"
	"errors"
	"flag"
	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/middleware"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	superUsername := flag.String("username", "", "超级管理员账户")
	superPassword := flag.String("password", "", "超级管理员密码")
	configPath := flag.String("c", "../config.yml", "配置文件路径")
	migrate := flag.Bool("m", false, "是否自动迁移数据库")
	generate := flag.Bool("g", false, "是否生成博客默认配置")
	flag.Parse()

	// 根据命令行参数读取配置文件, 其他变量的初始化依赖于配置文件对象
	conf, err := g.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	_ = initLogger(conf)
	db := initDatabase(conf, *migrate)
	rdb := initRedis(conf)

	// 如果命令行参数中指定了超级管理员账户和密码, 则创建超级管理员账户
	if *superUsername != "" && *superPassword != "" {
		createSuperAdmin(db, *superUsername, *superPassword)
		return
	}

	// 如果命令行参数中指定了 -g, 则生成博客默认配置
	if *generate {
		generateConfigs(db)
		return
	}

	// 初始化 gin 服务
	gin.SetMode(conf.Server.Mode)
	r := gin.New()
	r.SetTrustedProxies([]string{"*"})
	r.Use(gin.Recovery())
	// r.Use(middleware.Recovery(true)) // TODO: 自定义错误处理中间件
	r.Use(middleware.Logger())
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
	if serverAddr[0] == ':' || strings.HasPrefix(serverAddr, "0.0.0.0:") {
		log.Printf("Serving HTTP on (http://localhost:%s/) ... \n", strings.Split(serverAddr, ":")[1])
	} else {
		log.Printf("Serving HTTP on (http://%s/) ... \n", serverAddr)
	}
	r.Run(serverAddr)

}

// 根据配置文件初始化 slog 日志
func initLogger(conf *g.Config) *slog.Logger {
	var level slog.Level
	switch conf.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	option := &slog.HandlerOptions{
		AddSource: false, // TODO: 外层
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	switch conf.Log.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, option)
	case "text":
		fallthrough
	default:
		handler = slog.NewTextHandler(os.Stdout, option)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// 根据配置文件初始化 Redis
func initRedis(conf *g.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis 连接失败: ", err)
	}

	log.Println("Redis 连接成功", conf.Redis.Addr, conf.Redis.DB, conf.Redis.Password)
	return rdb
}

// 根据配置文件初始化数据库
func initDatabase(conf *g.Config, migrate bool) *gorm.DB {
	dbtype := conf.DbType()
	dsn := conf.DbDSN()

	var db *gorm.DB
	var err error

	var level logger.LogLevel
	switch conf.Log.Level {
	case "silent":
		level = logger.Silent
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	case "error":
		fallthrough
	default:
		level = logger.Error
	}

	config := &gorm.Config{
		Logger:                                   logger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
		SkipDefaultTransaction:                   true, // 禁用默认事务（提高运行速度）
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表名
		},
	}

	switch dbtype {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), config)
	default:
		log.Fatal("不支持的数据库类型: ", dbtype)
	}

	if err != nil {
		log.Fatal("数据库连接失败", err)
	}
	log.Println("数据库连接成功", dbtype, dsn)

	if migrate {
		if err := model.MakeMigrate(db); err != nil {
			log.Fatal("数据库迁移失败", err)
		}
		log.Println("数据库自动迁移成功")
	}

	return db
}

// 创建超级管理员
func createSuperAdmin(db *gorm.DB, username, password string) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var userAuth model.UserAuth
		err := db.Where("username = ?", username).First(&userAuth).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("创建超级管理员失败: " + err.Error())
		}

		if userAuth.ID != 0 {
			return errors.New("创建超级管理员失败: 账户已存在")
		}

		slog.Info("开始创建超级管理员")

		// 默认生成一个 super admin 用户
		userInfo := model.UserInfo{
			Nickname: username,
			Avatar:   "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
			Intro:    "这个人很懒，什么都没有留下",
			Website:  "https://www.hahacode.cn",
		}

		if err := db.Create(&userInfo).Error; err != nil {
			return errors.New("创建超级管理员失败: " + err.Error())
		}

		hashPassword, err := utils.BcryptHash(password)
		if err != nil {
			return errors.New("创建超级管理员失败: " + err.Error())
		}
		userAuth = model.UserAuth{
			UserInfoId:    userInfo.ID,
			Username:      username,
			Password:      hashPassword,
			IsSuper:       true,
			LastLoginTime: time.Now(),
		}
		if err := db.Create(&userAuth).Error; err != nil {
			return errors.New("创建超级管理员失败: " + err.Error())
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("创建超级管理员成功: %s, 密码: %s\n", username, password)
}

// TODO: 将初始化写成另外的可执行工具
// 初始化博客配置
// TODO: page 初始化
// TODO: resource 初始化
// TODO: role 初始化
func generateConfigs(db *gorm.DB) {
	configs := []model.Config{
		{Key: "website_avatar", Value: "https://foruda.gitee.com/avatar/1677041571085433939/5221991_szluyu99_1614389421.png", Desc: "网站头像"},
		{Key: "website_name", Value: "阵雨的个人博客", Desc: "网站名称"},
		{Key: "website_author", Value: "阵雨", Desc: "网站作者"},
		{Key: "website_intro", Value: "往事随风而去", Desc: "网站介绍"},
		{Key: "website_notice", Value: "欢迎来到阵雨的个人博客，项目还在开发中...", Desc: "网站公告"},
		{Key: "website_createtime", Value: time.Now().Format(time.DateTime), Desc: "网站创建日期"},
		{Key: "website_record", Value: "粤ICP备2021032312号", Desc: "网站备案号"},
		{Key: "qq", Value: "123456789", Desc: "QQ"},
		{Key: "github", Value: "https://github.com/szluyu99", Desc: "github"},
		{Key: "gitee", Value: "https://gitee.com/szluyu99", Desc: "gitee"},
		{Key: "tourist_avatar", Value: "https://cdn.hahacode.cn/16815451239215dc82548dcadcd578a5bbc8d5deaa.jpg", Desc: "默认游客头像"},
		{Key: "user_avatar", Value: "https://cdn.hahacode.cn/2299fc4d14c94e6183b082973b35855d.png", Desc: "默认用户头像"},
		{Key: "article_cover", Value: "https://cdn.hahacode.cn/1679461519cc592408198d67faf1290ff8969dc614.png", Desc: "默认文章封面"},
		{Key: "is_comment_review", Value: "true", Desc: "评论默认审核"},
		{Key: "is_message_review", Value: "true", Desc: "留言默认审核"},
	}

	result := db.Create(&configs)
	if result.Error != nil {
		slog.Error("初始化博客配置失败: ", result.Error)
		os.Exit(0)
	}

	slog.Info("初始化博客配置成功")
}
