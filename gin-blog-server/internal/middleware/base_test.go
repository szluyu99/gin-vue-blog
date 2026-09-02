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
	// 注意: 配置里 AllowOrigins 是 "*", 所以不会回显 Origin
	// 与 AllowCredentials: true 同时使用时, 浏览器会拒绝带 cookie 的跨域请求
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
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
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
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
