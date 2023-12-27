package middleware

import (
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/handle"
	"gin-blog/internal/model"
	"gin-blog/internal/utils/jwt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 资源访问权限验证
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("anonymous") {
			slog.Debug("[middleware-PermissionCheck] anonymous resource, skip check!")
			c.Next()
			return
		}

		db := c.MustGet(g.CTX_DB).(*gorm.DB)

		auth, err := handle.CurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c, g.ERROR_USER_NOT_EXIST, err)
			return
		}

		if auth.IsSuper {
			slog.Debug("[middleware-PermissionCheck]: super admin no need to check, pass!")
			c.Next()
			return
		}

		url := c.FullPath()[4:]
		method := c.Request.Method

		slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v, %v, %v\n", auth.Username, url, method))
		for _, role := range auth.Roles {
			pass, err := model.CheckRoleAuth(db, role.ID, url, method)
			if err != nil {
				handle.ReturnError(c, g.ERROR_DB_OPERATION, err)
				return
			}
			if !pass {
				handle.ReturnError(c, g.ERROR_PERMISSION_DENIED, nil)
				return
			}
		}

		slog.Debug("[middleware-PermissionCheck]: pass")
		c.Next()
	}
}

// 基于 JWT 的授权
// 如果存在 session, 则直接从 session 中获取用户信息
// 如果不存在 session, 则从 Authorization 中获取 token, 并解析 token 获取用户信息, 并设置到 session 中
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// FIXME: 前后台 session 混乱, 暂时无法将用户信息挂载在 gin context 缓存
		// auth, _ := handle.CurrentUserAuth(c)
		// if auth != nil {
		// 	slog.Debug("[middleware-JWTAuth] user auth exist, skip jwt auth")
		// 	c.Next()
		// 	return
		// }

		slog.Debug("[middleware-JWTAuth] user auth not exist, do jwt auth")

		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		resource, err := model.GetResource(db, c.FullPath()[4:], c.Request.Method)
		if err != nil {
			handle.ReturnError(c, g.ERROR_RESOURCE_NOT_EXIST, err)
			return
		}

		// 匿名资源, 直接跳过后续验证
		if resource.Anonymous {
			c.Set("anonymous", true)
			c.Next()
			c.Set("anonymous", false)
			return
		}

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			handle.ReturnError(c, g.ERROR_TOKEN_NOT_EXIST, nil)
			return
		}

		// token 的正确格式: `Bearer [tokenString]`
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handle.ReturnError(c, g.ERROR_TOKEN_TYPE_WRONG, nil)
			return
		}

		claims, err := jwt.ParseToken(g.Conf.JWT.Secret, parts[1])
		if err != nil {
			handle.ReturnError(c, g.ERROR_TOKEN_WRONG, err)
			return
		}

		// 判断 token 已过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			handle.ReturnError(c, g.ERROR_TOKEN_RUNTIME, nil)
			return
		}

		user, err := model.GetUserAuthInfoById(db, claims.UserId)
		if err != nil {
			handle.ReturnError(c, g.ERROR_USER_NOT_EXIST, err)
			return
		}

		// session
		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, claims.UserId)
		session.Save()

		// gin context
		c.Set(g.CTX_USER_AUTH, user)
	}
}
