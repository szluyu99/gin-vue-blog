package middleware

import (
	g "gin-blog/internal/global"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 跨域中间件放行预检请求
func TestCORSPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORS())
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"code": 0}) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "http://localhost:3333")
	req.Header.Set("Access-Control-Request-Method", "GET")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	// 带凭证的跨域请求必须回显具体 Origin, 不能是 *
	assert.Equal(t, "http://localhost:3333", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

// 普通请求也会带上跨域响应头
func TestCORSSimpleRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORS())
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"code": 0}) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "http://localhost:3333")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:3333", w.Header().Get("Access-Control-Allow-Origin"))
}

// handler panic 时不能把整个进程带走, 返回 500
func TestRecoveryFromPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery(true))
	r.GET("/boom", func(c *gin.Context) { panic("boom") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// 不打印堆栈的分支同样能兜住 panic
func TestRecoveryWithoutStack(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery(false))
	r.GET("/boom", func(c *gin.Context) { panic("boom") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// db / redis 客户端注入到 gin.Context, handler 里靠 MustGet 取
func TestWithGormAndRedisDB(t *testing.T) {
	e := newMwEnv(t)

	r := gin.New()
	r.Use(WithGormDB(e.db), WithRedisDB(e.rdb))
	r.GET("/ping", func(c *gin.Context) {
		assert.IsType(t, &gorm.DB{}, c.MustGet(g.CTX_DB))
		assert.IsType(t, &redis.Client{}, c.MustGet(g.CTX_RDB))
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))

	assert.Equal(t, http.StatusOK, w.Code)
}

// 日志中间件不影响响应
func TestLoggerPassThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Logger())
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"code": 0}) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping?a=1", nil))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"code":0`)
}

// session 中间件: 两次请求之间能带住 session
func TestWithCookieStore(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(WithCookieStore("test-session", "test-salt"))
	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, 1)
		_ = session.Save()
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})
	r.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": sessions.Default(c).Get(g.CTX_USER_AUTH)})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/set", nil))
	cookie := w.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	w2 := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.Header.Set("Cookie", cookie)
	r.ServeHTTP(w2, req)

	assert.Contains(t, w2.Body.String(), `"data":1`)
}

// session cookie 必须带 HttpOnly 与 SameSite: 前者防 XSS 读走会话, 后者防第三方站点自动带上
func TestCookieStoreOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(WithCookieStore("test-session", "test-salt"))
	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, 1)
		_ = session.Save()
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/set", nil))

	cookie := w.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, "HttpOnly")
	assert.Contains(t, cookie, "SameSite=Lax")
	// 本地 http 调试默认不加 Secure, 否则浏览器会直接丢掉 cookie
	assert.NotContains(t, cookie, "Secure")
}

// 跨域白名单: 不在名单里的站点不能带着访客登录态调接口
func TestCORSRejectsForeignOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORS())
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"code": 0}) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	r.ServeHTTP(w, req)

	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestOriginAllowed(t *testing.T) {
	// 配了白名单: 只认名单里的, 末尾斜杠和大小写不影响
	allowed := []string{"https://blog.example.com/"}
	assert.True(t, originAllowed("https://blog.example.com", allowed))
	assert.True(t, originAllowed("https://BLOG.example.com", allowed))
	assert.False(t, originAllowed("https://evil.example.com", allowed))
	// 白名单里有域名时, 内网来源也不再自动放行
	assert.False(t, originAllowed("http://localhost:3333", allowed))

	// 没配白名单: 本机与内网放行, 公网拒绝
	assert.True(t, originAllowed("http://localhost:3333", nil))
	assert.True(t, originAllowed("http://127.0.0.1:8080", nil))
	assert.True(t, originAllowed("http://10.27.244.168:8889", nil))
	assert.True(t, originAllowed("http://192.168.1.5", nil))
	assert.False(t, originAllowed("https://evil.example.com", nil))
	assert.False(t, originAllowed("http://8.8.8.8", nil))
	assert.False(t, originAllowed("", nil))
}
