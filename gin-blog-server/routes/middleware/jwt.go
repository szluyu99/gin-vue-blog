package middleware

import (
	"gin-blog/utils"
	"gin-blog/utils/r"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JWT 中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 约定 Token 放在 Header 的 Authorization 中, 并使用 Bearer 开头
		token := c.Request.Header.Get("Authorization")
		// token 为空
		if token == "" {
			r.SendCode(c, r.ERROR_TOKEN_NOT_EXIST)
			c.Abort()
			return
		}

		// token 的正确格式: `Bearer [tokenString]`
		parts := strings.Split(token, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			r.SendCode(c, r.ERROR_TOKEN_TYPE_WRONG)
			c.Abort()
			return
		}

		// parts[1] 是获取到的 tokenString, 使用 JWT 解析函数解析它
		claims, err := utils.GetJWT().ParseToken(parts[1])
		// token 解析失败
		if err != nil {
			r.SendData(c, r.ERROR_TOKEN_WRONG, err.Error())
			c.Abort()
			return
		}

		// 判断 token 已过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			r.SendCode(c, r.ERROR_TOKEN_RUNTIME)
			c.Abort()
			return
		}

		// 将当前请求的相关信息保存到请求的上下文 c 上
		// 后续的处理函数可以用过 c.Get("xxx") 来获取当前请求的用户信息
		c.Set("user_info_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Set("uuid", claims.UUID)
		c.Next()
	}
}
