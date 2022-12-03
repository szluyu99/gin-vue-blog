package middleware

import (
	"gin-blog/utils"
	"gin-blog/utils/r"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// recover 项目可能出现的 panic, 并使用 zap 记录相关日志
func ErrorRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 处理 panic(xxx) 的操作
				if code, ok := err.(int); ok { // panic(code) 根据错误码获取 msg
					r.SendCode(c, code)
				} else if msg, ok := err.(string); ok { // panic(string) 返回 string
					r.ReturnJson(c, http.StatusOK, r.FAIL, msg, nil)
				} else if e, ok := err.(error); ok { // panic(error) 发送消息
					r.ReturnJson(c, http.StatusOK, r.FAIL, e.Error(), nil)
				} else { // 其他
					r.Send(c, http.StatusOK, r.FAIL, nil)
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					utils.Logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					utils.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					utils.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
