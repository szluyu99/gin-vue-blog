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
