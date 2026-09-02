package middleware

import (
	g "gin-blog/internal/global"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常请求: 刷新 Redis 中的在线标记
func TestListenOnlineRefreshOnlineKey(t *testing.T) {
	e := newMwEnv(t)
	user := e.seedUser("test", false)
	e.loginId = user.ID
	e.handle(http.MethodGet, "/whatever", ListenOnline())

	w := e.get("/api/whatever")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan)

	onlineKey := g.ONLINE_USER + strconv.Itoa(user.ID)
	assert.True(t, e.mr.Exists(onlineKey))
	// 每次请求重新计算 10 分钟
	assert.Equal(t, 10*60, int(e.mr.TTL(onlineKey).Seconds()))
}

// 被强制下线的用户: 请求直接被拒绝
func TestListenOnlineForceOffline(t *testing.T) {
	e := newMwEnv(t)
	user := e.seedUser("test", false)
	e.loginId = user.ID
	assert.Nil(t, e.mr.Set(g.OFFLINE_USER+strconv.Itoa(user.ID), "1"))
	e.handle(http.MethodGet, "/whatever", ListenOnline())

	w := e.get("/api/whatever")

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "您已被强制下线")
	// 已下线的用户不应该被重新标记为在线
	assert.False(t, e.mr.Exists(g.ONLINE_USER+strconv.Itoa(user.ID)))
}

// 没有登录态: 认证异常
func TestListenOnlineWithoutLogin(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodGet, "/whatever", ListenOnline())

	w := e.get("/api/whatever")

	assert.False(t, e.handlerRan)
	assert.Contains(t, w.Body.String(), "用户认证异常")
}
