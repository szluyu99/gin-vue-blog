package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newCtx(t *testing.T, remoteAddr string, headers map[string]string) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.RemoteAddr = remoteAddr
	for k, v := range headers {
		c.Request.Header.Set(k, v)
	}
	return c
}

// Nginx 转发时优先用 X-Real-IP
func TestGetIpAddressPreferXRealIP(t *testing.T) {
	c := newCtx(t, "10.0.0.1:1234", map[string]string{
		"X-Real-IP":       "1.2.3.4",
		"X-Forwarded-For": "5.6.7.8, 9.9.9.9",
	})
	assert.Equal(t, "1.2.3.4", IP.GetIpAddress(c))
}

// 没有 X-Real-IP 时取 X-Forwarded-For 的第一个
func TestGetIpAddressFallbackToForwardedFor(t *testing.T) {
	c := newCtx(t, "10.0.0.1:1234", map[string]string{
		"X-Forwarded-For": "5.6.7.8,9.9.9.9",
	})
	assert.Equal(t, "5.6.7.8", IP.GetIpAddress(c))
}

// X-Real-IP 为 unknown 时视为没有
func TestGetIpAddressUnknownIsIgnored(t *testing.T) {
	c := newCtx(t, "10.0.0.1:1234", map[string]string{
		"X-Real-IP":          "unknown",
		"Proxy-Client-IP":    "unknown",
		"WL-Proxy-Client-IP": "2.2.2.2",
	})
	assert.Equal(t, "2.2.2.2", IP.GetIpAddress(c))
}

// 什么头都没有时用 RemoteAddr
func TestGetIpAddressFallbackToRemoteAddr(t *testing.T) {
	c := newCtx(t, "10.0.0.1:1234", nil)
	assert.Equal(t, "10.0.0.1:1234", IP.GetIpAddress(c))
}

// 本机地址会被替换成局域网地址, 不会原样返回 127.0.0.1
func TestGetIpAddressLocalhost(t *testing.T) {
	c := newCtx(t, "127.0.0.1:1234", nil)
	assert.NotContains(t, IP.GetIpAddress(c), "127.0.0.1")
}

// ip2region 数据库在测试环境下不存在, 只保证不 panic 且返回空
func TestGetIpSourceWithoutDB(t *testing.T) {
	assert.Empty(t, IP.GetIpSource("1.2.3.4"))
}

func TestGetUserAgent(t *testing.T) {
	c := newCtx(t, "10.0.0.1:1234", map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 " +
			"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})

	ua := IP.GetUserAgent(c)

	assert.NotNil(t, ua)
	assert.Equal(t, "Chrome", ua.Name)
	assert.Equal(t, "Mac OS X", ua.OS)
}
