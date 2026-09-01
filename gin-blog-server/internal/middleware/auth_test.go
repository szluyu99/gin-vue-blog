package middleware

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 构造一个挂了 JWTAuth 的引擎, handlerRan 用于判断是否走到了 handler
func newAuthEngine(t *testing.T, requireLogin bool, handlerRan *bool) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	assert.Nil(t, err)
	assert.Nil(t, model.MakeMigrate(db))

	r := gin.New()
	r.Use(sessions.Sessions("test-session", memstore.NewStore([]byte("secret"))))
	r.Use(func(c *gin.Context) {
		c.Set(g.CTX_DB, db)
		c.Next()
	})

	api := r.Group("/api")
	api.Use(JWTAuth(requireLogin))
	api.GET("/whatever", func(c *gin.Context) {
		*handlerRan = true
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	return r
}

func doGet(r *gin.Engine, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

// 前台接口: 资源表中没登记的接口允许匿名访问
func TestJWTAuthOptionalLogin(t *testing.T) {
	var ran bool
	r := newAuthEngine(t, false, &ran)

	w := doGet(r, "/api/whatever")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, ran, "前台接口未登录时也应该执行 handler")
}

// 后台接口: 资源表中没登记的接口也必须登录, 否则会变成完全无鉴权
func TestJWTAuthRequireLogin(t *testing.T) {
	var ran bool
	r := newAuthEngine(t, true, &ran)

	w := doGet(r, "/api/whatever")
	assert.Equal(t, http.StatusOK, w.Code) // 业务错误也是 HTTP 200
	assert.False(t, ran, "后台接口未登录时不应该执行 handler")
	assert.Contains(t, w.Body.String(), "TOKEN 不存在")
}

// 带了格式错误的 token 同样被拒绝
func TestJWTAuthRequireLoginBadToken(t *testing.T) {
	var ran bool
	r := newAuthEngine(t, true, &ran)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/whatever", nil)
	req.Header.Set("Authorization", "not-a-bearer-token")
	r.ServeHTTP(w, req)

	assert.False(t, ran)
	assert.Contains(t, w.Body.String(), "TOKEN 格式错误")
}
