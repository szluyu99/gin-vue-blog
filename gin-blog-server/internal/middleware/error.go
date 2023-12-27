package middleware

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// TODO:
// recover 项目可能出现的 panic
func ErrorRecovery(stack bool) gin.HandlerFunc {
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
