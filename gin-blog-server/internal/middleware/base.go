package middleware

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/handle"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// WithRedisDB 将 redis.Client 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_RDB).(*redis.Client) 来使用
func WithRedisDB(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(g.CTX_RDB, rdb)
		ctx.Next()
	}
}

// WithGormDB 将 gorm.DB 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_DB).(*gorm.DB) 来使用
func WithGormDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(g.CTX_DB, db)
		ctx.Next()
	}
}

// CORS 跨域请求
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	})
}

// WithCookieStore 基于 cookie 的 session
func WithCookieStore(name, secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{Path: "/", MaxAge: 600})
	return sessions.Sessions(name, store)
}

// WithMemStore 基于内存的 session
func WithMemStore(name, secret string) gin.HandlerFunc {
	store := memstore.NewStore([]byte(secret))
	store.Options(sessions.Options{Path: "/", MaxAge: 600})
	return sessions.Sessions(name, store)
}

// Logger 日志记录
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start)

		slog.Info("[GIN]",
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("ip", c.ClientIP()),
			slog.Int("size", c.Writer.Size()),
			slog.Duration("cost", cost),
			// slog.String("body", c.Request.PostForm.Encode()),
			// slog.String("user-agent", c.Request.UserAgent()),
			// slog.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

// Recovery 恢复中间件
func Recovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 发生 panic, 返回错误信息
				handle.ReturnHttpResponse(c, http.StatusInternalServerError, g.FAIL, g.GetMsg(g.FAIL), err)

				// 处理 panic(xxx) 的操作
				// if code, ok := err.(int); ok { // panic(code) 根据错误码获取 msg
				// 	v2.Return(c, code, nil)
				// } else if msg, ok := err.(string); ok { // panic(string) 返回 string
				// 	v2.ReturnJSON(c, http.StatusOK, g.FAIL, msg, nil)
				// } else if e, ok := err.(error); ok { // panic(error) 发送消息
				// 	v2.ReturnJSON(c, http.StatusOK, g.FAIL, e.Error(), nil)
				// } else { // 其他
				// 	v2.Return(c, g.FAIL, nil)
				// }

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					slog.Error(c.Request.URL.Path,
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // errcheck
					c.Abort()
					return
				}

				if stack {
					slog.Error("[Recovery from panic]",
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
						slog.String("stack", string(debug.Stack())),
					)
				} else {
					slog.Error("[Recovery from panic]",
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
