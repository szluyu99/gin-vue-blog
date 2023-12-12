package utils

import (
	"fmt"
	"gin-blog/config"
	"gin-blog/model"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitSQLiteDB(dsn string) *gorm.DB {
	if dsn == "" {
		dsn = "dev.db"
	}
	db, err := gorm.Open(sqlite.Open(dsn), gormConfig())

	if err != nil {
		log.Fatal("SQLite 连接失败, 请检查参数")
	}

	log.Println("SQLite 连接成功")

	// MakeMigrate(db)

	return db
}

func InitMySQLDB() *gorm.DB {
	mysqlCfg := config.Cfg.Mysql
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlCfg.Username,
		mysqlCfg.Password,
		mysqlCfg.Host,
		mysqlCfg.Port,
		mysqlCfg.Dbname,
	)

	db, err := gorm.Open(mysql.Open(dns), gormConfig())

	if err != nil {
		fmt.Println("dns: ", dns)
		log.Fatalln("MySQL 连接失败, 请检查参数", err)
	}

	log.Println("MySQL 连接成功")

	// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
	// MakeMigrate(db)

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)                  // 设置连接池中的最大闲置连接
	sqlDB.SetMaxOpenConns(100)                 // 设置数据库的最大连接数量
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 设置连接的最大可复用时间

	return db
}

func gormConfig() *gorm.Config {
	return &gorm.Config{
		// gorm 日志模式
		Logger: logger.Default.LogMode(getLogMode(config.Cfg.Mysql.LogMode)),
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	}
}

// 根据字符串获取对应 LogLevel
func getLogMode(str string) logger.LogLevel {
	switch str {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}

// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
// 只支持创建表、增加表中没有的字段和索引
// 为了保护数据，并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Article{},
		&model.Category{},
		&model.Comment{},
		&model.Tag{},
		&model.Message{},
		&model.UserInfo{},
		&model.FriendLink{},
		// 权限相关
		&model.UserAuth{},     // 用户
		&model.Role{},         // 角色
		&model.Menu{},         // 菜单
		&model.Resource{},     // 资源(接口)
		&model.UserRole{},     // 用户-角色 关联
		&model.RoleMenu{},     // 角色-菜单 关联
		&model.RoleResource{}, // 角色-资源 关联

		&model.Page{},         // 页面
		&model.BlogConfig{},   // 网站设置
		&model.OperationLog{}, // 操作日志
	)

	if err != nil {
		log.Println("gorm 自动迁移失败: ", err)
	} else {
		log.Println("gorm 自动迁移成功")
	}
}
