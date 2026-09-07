package middleware

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/handle"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
/*
不能用 AllowOrigins: ["*"], 它会让响应头固定为 *,
而带凭证(AllowCredentials)的跨域请求在 * 下会被浏览器拒绝,
前后端分域名部署时 session cookie 就带不过去。所以用 AllowOriginFunc 回显具体 Origin。

但原来这个函数恒 return true, 「全放行 + 回显 Origin + 允许凭证」在安全上等同于任意站点
都能带着访客的登录态调评论、留言、改资料。改成白名单:
来源写在 config.yml 的 server.allowed-origins, 留空时只放行本机与内网(开发和内网自用场景)。
*/
func CORS() gin.HandlerFunc {
	allowed := g.Conf.AllowedOrigins()
	if len(allowed) == 0 {
		slog.Warn("未配置 server.allowed-origins, 跨域只放行本机与内网来源")
	}

	return cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return originAllowed(origin, allowed)
		},
		MaxAge: 24 * time.Hour,
	})
}

func originAllowed(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}
	if len(allowed) > 0 {
		for _, item := range allowed {
			if strings.EqualFold(strings.TrimRight(item, "/"), strings.TrimRight(origin, "/")) {
				return true
			}
		}
		return false
	}
	return isLocalOrigin(origin)
}

// 本机 / 内网来源: localhost、127.x、::1 与私有网段。公网站点一律拦掉
func isLocalOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := u.Hostname()
	if host == "localhost" || host == "::1" {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}

// WithCookieStore 基于 cookie 的 session
func WithCookieStore(name, secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessionOptions())
	return sessions.Sessions(name, store)
}

// WithMemStore 基于内存的 session
func WithMemStore(name, secret string) gin.HandlerFunc {
	store := memstore.NewStore([]byte(secret))
	store.Options(sessionOptions())
	return sessions.Sessions(name, store)
}

/*
session cookie 属性

原来只设了 Path 与 MaxAge: 没有 HttpOnly 意味着任意 XSS 脚本能读走会话,
没有 SameSite 则第三方站点发出的请求会自动带上它。
Secure 走配置(session.secure): 本地 http 调试时置 true 浏览器会直接丢掉 cookie。
*/
func sessionOptions() sessions.Options {
	return sessions.Options{
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   g.Conf.SessionSecure(),
		SameSite: http.SameSiteLaxMode,
	}
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
