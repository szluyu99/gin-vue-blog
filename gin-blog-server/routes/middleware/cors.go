package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 跨域中间件
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许跨域请求网站
		AllowOrigins: []string{"*"},
		// 允许使用的请求方式
		AllowMethods: []string{"PUT", "POST", "GET", "DELETE", "OPTIONS", "PATCH"},
		// 允许使用的请求头
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		// 暴露的请求头
		ExposeHeaders: []string{"Content-Type"},
		// 凭证共享
		AllowCredentials: true,
		// 允许跨域的源网站
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		// 超时时间设定
		MaxAge: 24 * time.Hour,
	})
}
