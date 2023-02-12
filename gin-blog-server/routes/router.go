package routes

import (
	"gin-blog/config"
	"gin-blog/dao"
	"gin-blog/utils"
	"log"
	"net/http"
	"time"
)

// 初始化全局变量
func InitGlobalVariable() {
	// 初始化 Viper
	utils.InitViper()
	// 初始化 Logger
	utils.InitLogger()
	// 初始化数据库 DB
	dao.DB = utils.InitMySQLDB() // 需要先导入 gvb.sql
	// dao.DB = utils.InitSQLiteDB("gorm.db") // TODO: 默认无数据，暂时无法使用
	// 初始化 Redis
	utils.InitRedis()
	// 初始化 Casbin
	utils.InitCasbin(dao.DB)
}

// 后台服务
func BackendServer() *http.Server {
	backPort := config.Cfg.Server.BackPort
	log.Printf("后台服务启动于 %s 端口", backPort)
	return &http.Server{
		Addr:         backPort,
		Handler:      BackRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// 前台服务
func FrontendServer() *http.Server {
	frontPort := config.Cfg.Server.FrontPort
	log.Printf("前台服务启动于 %s 端口", frontPort)
	return &http.Server{
		Addr:         frontPort,
		Handler:      FrontRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
