package ginblog

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v9"
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
		log.Fatal("不支持的数据库类型: ", dbtype)
	}

	if err != nil {
		log.Fatal("数据库连接失败", err)
	}
	log.Println("数据库连接成功", dbtype, dsn)

	if conf.Server.DbAutoMigrate {
		if err := model.MakeMigrate(db); err != nil {
			log.Fatal("数据库迁移失败", err)
		}
		log.Println("数据库自动迁移成功")
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
		log.Fatal("Redis 连接失败: ", err)
	}

	log.Println("Redis 连接成功", conf.Redis.Addr, conf.Redis.DB, conf.Redis.Password)
	return rdb
}
