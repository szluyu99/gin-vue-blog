package ginblog

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 根据配置文件初始化 slog 日志
func InitLogger(conf *g.Config) *slog.Logger {
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

// 根据配置文件初始化数据库
func InitDatabase(conf *g.Config) *gorm.DB {
	dbtype := conf.DbType()
	dsn := conf.DbDSN()

	var db *gorm.DB
	var err error

	var level logger.LogLevel
	switch conf.Server.DbLogMode {
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
		fatal("不支持的数据库类型: " + dbtype)
	}

	if err != nil {
		fatal("数据库连接失败", "err", err)
	}
	slog.Info("数据库连接成功", "type", dbtype, "dsn", dsn)

	if conf.Server.DbAutoMigrate {
		if err := model.MakeMigrate(db); err != nil {
			fatal("数据库迁移失败", "err", err)
		}
		slog.Info("数据库自动迁移成功")
	}

	return db
}

// 根据配置文件初始化 Redis
func InitRedis(conf *g.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fatal("Redis 连接失败", "err", err)
	}

	slog.Info("Redis 连接成功", "addr", conf.Redis.Addr, "db", conf.Redis.DB)
	return rdb
}

/*
启动阶段的致命错误

slog.SetDefault 会把标准库 log 的输出接到 slog 上, 并且当成 Info 级别,
配置里 Log.Level 是 error 时(部署用的 config.docker.yml 就是), log.Fatal
的内容会被直接丢掉, 容器只能看到一个 exit 1。这里统一用 Error 级别输出。
*/
func fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
