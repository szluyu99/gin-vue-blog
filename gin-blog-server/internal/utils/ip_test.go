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

// 构造一个只信任内网代理的 context, 贴近真实部署
func newProxyCtx(t *testing.T, remoteAddr string, headers map[string]string) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)
	e := gin.New()
	if err := e.SetTrustedProxies([]string{"10.0.0.0/8"}); err != nil {
		t.Fatal(err)
	}
	c := gin.CreateTestContextOnly(httptest.NewRecorder(), e)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.RemoteAddr = remoteAddr
	for k, v := range headers {
		c.Request.Header.Set(k, v)
	}
	return c
}

// 来自可信代理的转发头才采信
func TestGetIpAddressTrustsProxyHeader(t *testing.T) {
	c := newProxyCtx(t, "10.0.0.1:1234", map[string]string{"X-Real-IP": "1.2.3.4"})
	assert.Equal(t, "1.2.3.4", IP.GetIpAddress(c))
}

// 直连的访客自己带 X-Real-IP / X-Forwarded-For 不算数, 只能拿到对端地址。
// 否则轮换这个头就能绕过按 IP 的登录失败次数限制
func TestGetIpAddressIgnoresSpoofedHeader(t *testing.T) {
	c := newProxyCtx(t, "203.0.113.7:1234", map[string]string{
		"X-Real-IP":       "1.2.3.4",
		"X-Forwarded-For": "5.6.7.8",
	})
	assert.Equal(t, "203.0.113.7", IP.GetIpAddress(c))
}

// 取到的地址不带端口: 端口每次请求都变, 带上去 按 IP 的限流和去重就失效了
func TestGetIpAddressHasNoPort(t *testing.T) {
	c := newProxyCtx(t, "203.0.113.7:1234", nil)
	assert.Equal(t, "203.0.113.7", IP.GetIpAddress(c))
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
