package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 前台接口: 资源表中没登记的接口允许匿名访问
func TestJWTAuthOptionalLogin(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodGet, "/whatever", JWTAuth(false))

	w := e.get("/api/whatever")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan, "前台接口未登录时也应该执行 handler")
}

// 后台接口: 资源表中没登记的接口也必须登录, 否则会变成完全无鉴权
func TestJWTAuthRequireLogin(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodGet, "/whatever", JWTAuth(true))

	w := e.get("/api/whatever")

	assert.Equal(t, http.StatusOK, w.Code) // 业务错误也是 HTTP 200
	assert.False(t, e.handlerRan, "后台接口未登录时不应该执行 handler")
	assert.Contains(t, w.Body.String(), "TOKEN 不存在")
}

// 带了格式错误的 token 同样被拒绝
func TestJWTAuthRequireLoginBadToken(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodGet, "/whatever", JWTAuth(true))

	w := e.request(http.MethodGet, "/api/whatever", "", map[string]string{
		"Authorization": "not-a-bearer-token",
	})

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "TOKEN 格式错误")
}

// 资源表中登记为匿名的接口: 不需要 token
func TestJWTAuthAnonymousResource(t *testing.T) {
	e := newMwEnv(t)
	e.seedResource("匿名接口", "/config", http.MethodGet, true)
	e.handle(http.MethodGet, "/config", JWTAuth(true))

	w := e.get("/api/config")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan, "匿名资源未登录时也应该执行 handler")
}

// 资源表中登记且非匿名的接口: 没有 token 直接拒绝
func TestJWTAuthRegisteredResourceWithoutToken(t *testing.T) {
	e := newMwEnv(t)
	e.seedResource("配置修改", "/config", http.MethodPatch, false)
	e.handle(http.MethodPatch, "/config", JWTAuth(true))

	w := e.request(http.MethodPatch, "/api/config", "{}", nil)

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "TOKEN 不存在")
}

// 正常的 token: 放行, 并把用户信息挂到 context 上
func TestJWTAuthValidToken(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	user := e.seedUser("admin", true)
	e.seedResource("配置修改", "/config", http.MethodPatch, false)
	e.handle(http.MethodPatch, "/config", JWTAuth(true))

	w := e.request(http.MethodPatch, "/api/config", "{}", bearer(t, user.ID, 24))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan)
	// 后续请求可以直接用 session 里的用户身份
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

// 用别的密钥签出来的 token 不被接受
func TestJWTAuthWrongSecret(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	user := e.seedUser("admin", true)
	e.seedResource("配置修改", "/config", http.MethodPatch, false)
	e.handle(http.MethodPatch, "/config", JWTAuth(true))

	headers := bearerWithSecret(t, "another-secret", user.ID, 24)

	w := e.request(http.MethodPatch, "/api/config", "{}", headers)

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "TOKEN 不正确")
}

// 过期的 token 被拒绝
// 注意: 解析阶段就会报过期, 所以返回的是 1203 而不是 1202
func TestJWTAuthExpiredToken(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	user := e.seedUser("admin", true)
	e.seedResource("配置修改", "/config", http.MethodPatch, false)
	e.handle(http.MethodPatch, "/config", JWTAuth(true))

	w := e.request(http.MethodPatch, "/api/config", "{}", bearer(t, user.ID, -1))

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "TOKEN 不正确")
}

// token 合法但用户已经被删除
func TestJWTAuthUserNotExist(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	e.seedResource("配置修改", "/config", http.MethodPatch, false)
	e.handle(http.MethodPatch, "/config", JWTAuth(true))

	w := e.request(http.MethodPatch, "/api/config", "{}", bearer(t, 999, 24))

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "该用户不存在")
}

// 前台接口带了有效 token: 即使接口没登记在资源表里, 也要识别出当前用户
// (前台的登录态之前只能靠 session cookie, session 只有 10 分钟)
func TestJWTAuthOptionalLoginIdentifiesUser(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	user := e.seedUser("tester", false)
	e.handle(http.MethodPost, "/front/upload", JWTAuth(false))

	w := e.request(http.MethodPost, "/api/front/upload", "{}", bearer(t, user.ID, 24))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan)
	assert.NotNil(t, e.handlerUser, "handler 应该能从 context 里拿到用户")
	assert.Equal(t, "tester", e.handlerUser.Username)
	// session 每次请求都会刷新, 不会因为 10 分钟没动就失效
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

// 前台接口不带 token: 匿名放行, handler 自己决定要不要用户
func TestJWTAuthOptionalLoginAnonymous(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	e.handle(http.MethodGet, "/front/article/list", JWTAuth(false))

	w := e.get("/api/front/article/list")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan)
	assert.Nil(t, e.handlerUser)
}

// 前台接口带了坏掉的 token: 直接报错, 不降级成匿名
// 否则前端无法区分"没登录"和"登录过期"
func TestJWTAuthOptionalLoginRejectsBadToken(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	user := e.seedUser("tester", false)
	e.handle(http.MethodPost, "/front/upload", JWTAuth(false))

	w := e.request(http.MethodPost, "/api/front/upload", "{}", bearerWithSecret(t, "another-secret", user.ID, 24))

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "TOKEN 不正确")
}

// 前台接口带了过期的 token: 同样直接报错
func TestJWTAuthOptionalLoginRejectsExpiredToken(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	user := e.seedUser("tester", false)
	e.handle(http.MethodPost, "/front/upload", JWTAuth(false))

	w := e.request(http.MethodPost, "/api/front/upload", "{}", bearer(t, user.ID, -1))

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "TOKEN 不正确")
}

// 资源表中登记为匿名的接口: 不带 token 放行, 带 token 则识别用户
func TestJWTAuthAnonymousResourceIdentifiesUser(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	user := e.seedUser("tester", false)
	e.seedResource("匿名接口", "/config", http.MethodGet, true)
	e.handle(http.MethodGet, "/config", JWTAuth(true), PermissionCheck())

	w := e.request(http.MethodGet, "/api/config", "", bearer(t, user.ID, 24))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan, "匿名资源不做权限校验")
	assert.NotNil(t, e.handlerUser)
	assert.Equal(t, "tester", e.handlerUser.Username)
}

// 后台接口带了有效 token 但资源表中没登记: 允许访问(仅跳过权限校验), 用户要识别出来
func TestJWTAuthRequireLoginIdentifiesUser(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	user := e.seedUser("admin", true)
	e.handle(http.MethodPost, "/whatever", JWTAuth(true), PermissionCheck())

	w := e.request(http.MethodPost, "/api/whatever", "{}", bearer(t, user.ID, 24))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan)
	assert.NotNil(t, e.handlerUser)
}
