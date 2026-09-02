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
//
// 无论接口是否允许匿名访问, 只要请求带了 Authorization 就会解析 token,
// 把用户挂到 gin context 与 session 上 —— 前台接口靠这一步识别当前用户。
// (之前资源表中没登记的接口会直接跳过解析, 前台的登录态只能依赖 session cookie,
// session 只有 10 分钟, 过期后发评论/上传头像都会莫名返回 TOKEN 不存在)
//
// requireLogin 为 true 时, 资源表中不存在的接口也必须携带有效 token,
// 仅跳过权限校验(fail closed)。后台接口必须用 true, 否则新增接口忘记
// 在资源表登记, 该接口就会完全无鉴权。
// 前台接口用 false: 大部分是匿名可读的, 只是顺便识别一下当前用户。
func JWTAuth(requireLogin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet(g.CTX_DB).(*gorm.DB)

		// 系统管理的资源需要做验证, 没有加进来的不需要
		url, method := c.FullPath()[4:], c.Request.Method
		resource, err := model.GetResource(db, url, method)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			handle.ReturnError(c, g.ErrDbOp, err)
			return
		}
		registered := err == nil

		// 资源表中没登记的接口, 以及登记为匿名的接口, 都不做权限校验
		if !registered || resource.Anonymous {
			slog.Debug(fmt.Sprintf("[middleware-JWTAuth] resource: %s %s skip permission check", url, method))
			c.Set("skip_check", true)
		}

		// 是否必须携带有效 token:
		// 资源表中登记且非匿名的接口必须; 没登记的接口由 requireLogin 决定
		mustLogin := (registered && !resource.Anonymous) || (!registered && requireLogin)

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			if mustLogin {
				handle.ReturnError(c, g.ErrTokenNotExist, nil)
				return
			}
			// 匿名访问: 后续 handler 自己决定要不要用户信息
			slog.Debug("[middleware-JWTAuth] no authorization header, continue as anonymous")
			return
		}

		// 带了凭证就必须是有效的, 坏掉的 token 直接报错而不是降级成匿名,
		// 否则前端无法区分"没登录"和"登录过期"
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

		// session: 每次请求都刷新一次, 避免 10 分钟后失效
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
		// 任一角色拥有该资源即放行 (OR 语义):
		// 多角色时不能因为其中一个角色没有权限就拒绝, 否则多加一个弱角色反而会减少权限
		for _, role := range auth.Roles {
			slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v\n", role.Name))
			pass, err := model.CheckRoleAuth(db, role.ID, url, method)
			if err != nil {
				handle.ReturnError(c, g.ErrDbOp, err)
				return
			}
			if pass {
				slog.Debug("[middleware-PermissionCheck]: pass")
				c.Next()
				return
			}
		}

		handle.ReturnError(c, g.ErrPermission, nil)
	}
}
