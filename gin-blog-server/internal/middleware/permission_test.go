package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 超级管理员不做权限校验
func TestPermissionCheckSuperAdmin(t *testing.T) {
	e := newMwEnv(t)
	user := e.seedUser("admin", true)
	e.loginId = user.ID
	e.handle(http.MethodPatch, "/config", PermissionCheck())

	w := e.request(http.MethodPatch, "/api/config", "{}", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan)
}

// 角色下挂了该资源: 放行
func TestPermissionCheckRoleHasResource(t *testing.T) {
	e := newMwEnv(t)
	res := e.seedResource("配置修改", "/config", http.MethodPatch, false)
	role := e.seedRole("admin", res.ID)
	user := e.seedUser("test", false, role.ID)
	e.loginId = user.ID
	e.handle(http.MethodPatch, "/config", PermissionCheck())

	w := e.request(http.MethodPatch, "/api/config", "{}", nil)

	assert.True(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), `"code":0`)
}

// 角色下没有该资源: 权限不足
func TestPermissionCheckRoleWithoutResource(t *testing.T) {
	e := newMwEnv(t)
	other := e.seedResource("文章列表", "/article/list", http.MethodGet, false)
	e.seedResource("配置修改", "/config", http.MethodPatch, false)
	role := e.seedRole("user", other.ID) // 只有文章列表的权限
	user := e.seedUser("test", false, role.ID)
	e.loginId = user.ID
	e.handle(http.MethodPatch, "/config", PermissionCheck())

	w := e.request(http.MethodPatch, "/api/config", "{}", nil)

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "权限不足")
}

// 同一个 url 不同 method 是两条资源, 不能互相顶用
func TestPermissionCheckMethodMatters(t *testing.T) {
	e := newMwEnv(t)
	get := e.seedResource("配置查询", "/config", http.MethodGet, false)
	e.seedResource("配置修改", "/config", http.MethodPatch, false)
	role := e.seedRole("user", get.ID)
	user := e.seedUser("test", false, role.ID)
	e.loginId = user.ID
	e.handle(http.MethodPatch, "/config", PermissionCheck())

	w := e.request(http.MethodPatch, "/api/config", "{}", nil)

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "权限不足")
}

// 没有登录态: 拿不到用户, 拒绝
func TestPermissionCheckWithoutLogin(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodPatch, "/config", PermissionCheck())

	w := e.request(http.MethodPatch, "/api/config", "{}", nil)

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "该用户不存在")
}

/*
多角色时是 AND 语义: 任一角色缺少该资源就会被拒绝。
按常见 RBAC 的理解应该是 OR (任一角色有权限即可), 这里保持现状是因为
它 fail closed, 改成 OR 属于放宽权限, 需要单独评估。这条用例用来锁住现状,
以后要改语义时会先在这里失败。
*/
func TestPermissionCheckMultiRoleIsAnd(t *testing.T) {
	e := newMwEnv(t)
	res := e.seedResource("配置修改", "/config", http.MethodPatch, false)
	full := e.seedRole("admin", res.ID)
	limited := e.seedRole("user") // 什么资源都没有
	user := e.seedUser("test", false, full.ID, limited.ID)
	e.loginId = user.ID
	e.handle(http.MethodPatch, "/config", PermissionCheck())

	w := e.request(http.MethodPatch, "/api/config", "{}", nil)

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "权限不足")
}

// JWTAuth 对资源表中没登记的接口会设置 skip_check, PermissionCheck 直接放行
func TestPermissionCheckSkippedForUnregisteredResource(t *testing.T) {
	withJWTConf(t)
	e := newMwEnv(t)
	role := e.seedRole("user") // 没有任何资源
	user := e.seedUser("test", false, role.ID)
	e.handle(http.MethodPatch, "/config", JWTAuth(true), PermissionCheck())

	w := e.request(http.MethodPatch, "/api/config", "{}", bearer(t, user.ID, 24))

	// 接口没有登记到资源表: 要求登录, 但不校验权限
	assert.True(t, e.handlerRan, "资源表中没有的接口只要登录就能访问")
	assert.Contains(t, w.Body.String(), `"code":0`)
}
