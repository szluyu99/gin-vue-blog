package main

import (
	"gin-blog/routes"
	"log"

	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

// @title Gin-Vue-Blog API 文档
// @version 1.0 版本
// @description Gin-Vue-Blog 的 Swagger API 文档

// @host		localhost:8765
// @BasePath /api/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	// 初始化全局变量
	routes.InitGlobalVariable()

	// 前台接口服务
	g.Go(func() error {
		return routes.FrontendServer().ListenAndServe()
	})

	// 后台接口服务
	g.Go(func() error {
		return routes.BackendServer().ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
