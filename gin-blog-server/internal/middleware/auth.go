package middleware

import (
	"errors"
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

// 基于 JWT 的授权
// 如果存在 session, 则直接从 session 中获取用户信息
// 如果不存在 session, 则从 Authorization 中获取 token, 并解析 token 获取用户信息, 并设置到 session 中
// requireLogin 为 true 时, 资源表中不存在的接口也必须携带有效 token,
// 仅跳过权限校验(fail closed)。后台接口必须用 true, 否则新增接口忘记
// 在资源表登记, 该接口就会完全无鉴权。
// 前台接口用 false: 大部分是匿名可读的, 只是顺便识别一下当前用户。
func JWTAuth(requireLogin bool) gin.HandlerFunc {
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

		// 系统管理的资源需要做验证, 没有加进来的不需要
		url, method := c.FullPath()[4:], c.Request.Method
		resource, err := model.GetResource(db, url, method)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				handle.ReturnError(c, g.ErrDbOp, err)
				return
			}
			// 资源表中没有登记的接口: 不做权限校验
			slog.Debug("[middleware-JWTAuth] resource not exist, skip permission check")
			c.Set("skip_check", true)
			// 前台接口允许匿名访问, 后台接口仍然要求登录
			if !requireLogin {
				c.Next()
				c.Set("skip_check", false)
				return
			}
		}

		// 匿名资源, 直接跳过后续验证 (资源表中没登记时 resource 为零值)
		if resource.Anonymous {
			slog.Debug(fmt.Sprintf("[middleware-JWTAuth] resource: %s %s is anonymous, skip jwt auth!", url, method))
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			handle.ReturnError(c, g.ErrTokenNotExist, nil)
			return
		}

		// token 的正确格式: `Bearer [tokenString]`
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handle.ReturnError(c, g.ErrTokenType, nil)
			return
		}

		claims, err := jwt.ParseToken(g.Conf.JWT.Secret, parts[1])
		if err != nil {
			handle.ReturnError(c, g.ErrTokenWrong, err)
			return
		}

		// 判断 token 已过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			handle.ReturnError(c, g.ErrTokenRuntime, nil)
			return
		}

		user, err := model.GetUserAuthInfoById(db, claims.UserId)
		if err != nil {
			handle.ReturnError(c, g.ErrUserNotExist, err)
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

// 资源访问权限验证
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("skip_check") {
			c.Next()
			return
		}

		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		auth, err := handle.CurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c, g.ErrUserNotExist, err)
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
			slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v\n", role.Name))
			pass, err := model.CheckRoleAuth(db, role.ID, url, method)
			if err != nil {
				handle.ReturnError(c, g.ErrDbOp, err)
				return
			}
			if !pass {
				handle.ReturnError(c, g.ErrPermission, nil)
				return
			}
		}

		slog.Debug("[middleware-PermissionCheck]: pass")
		c.Next()
	}
}
