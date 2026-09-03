package handle

import (
	"bytes"
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/*
接口级测试脚手架

- 数据库: sqlite 内存库, 每个测试独立, 用真实的 model 迁移
- Redis: miniredis 内存实现, 不需要本地装 Redis
- 路由: 只挂被测 handler, 不经过 JWTAuth 和 PermissionCheck
  (鉴权中间件在 internal/middleware, 引入会造成 import 循环)
*/

type testEnv struct {
	engine *gin.Engine
	db     *gorm.DB
	rdb    *redis.Client
	mr     *miniredis.Miniredis
	user   *model.UserAuth // 非 nil 时模拟该用户已登录
}

// 模拟登录: CurrentUserAuth 会优先从 gin.Context 取用户, 不需要真的走 session
func (e *testEnv) loginAs(id int, username string) *model.UserAuth {
	e.user = &model.UserAuth{Model: model.Model{ID: id}, Username: username}
	return e.user
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// handler 里读配置的地方(如注册是否走邮箱验证)需要 g.Conf 已初始化。
	// 注意不能直接覆盖: withJWTConf / withUploadConf 可能已经在本测试里设过了
	if g.Conf == nil {
		old := g.Conf
		g.Conf = &g.Config{}
		t.Cleanup(func() { g.Conf = old })
	}
	g.Conf.Captcha.SendEmail = false
	g.Conf.Captcha.ExpireTime = 15

	// cache=shared + 每个测试独立的库名: file::memory: 会让连接池里的每条连接
	// 各自持有一个空库, 事务里换连接就会报 "no such table"
	dsn := "file:" + strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()) +
		"?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	assert.Nil(t, err)
	assert.Nil(t, model.MakeMigrate(db))

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	env := &testEnv{db: db, rdb: rdb, mr: mr}

	engine := gin.New()
	// CurrentUserAuth 会读 session, 生产环境在 cmd/main.go 中注册, 测试里也要有
	engine.Use(sessions.Sessions("test-session", memstore.NewStore([]byte("test-secret"))))
	engine.Use(func(c *gin.Context) {
		c.Set(g.CTX_DB, db)
		c.Set(g.CTX_RDB, rdb)
		if env.user != nil {
			c.Set(g.CTX_USER_AUTH, env.user)
		}
		c.Next()
	})
	env.engine = engine

	return env
}

// 发起请求, 返回解析后的响应
func (e *testEnv) do(t *testing.T, method, path string, body any) Response[any] {
	t.Helper()
	return e.doWithHeader(t, method, path, body, nil)
}

// 带自定义请求头的版本, 用于模拟不同访客(X-Real-IP)等场景
func (e *testEnv) doWithHeader(t *testing.T, method, path string, body any, header map[string]string) Response[any] {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		assert.Nil(t, err)
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	e.engine.ServeHTTP(w, req)

	// 业务错误也返回 HTTP 200, 非 200 说明中间件或 gin 层面出了问题
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response[any]
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

// 以 multipart/form-data 上传一个文件, 字段名固定为 file
func (e *testEnv) upload(t *testing.T, path, filename, content string) Response[any] {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	assert.Nil(t, err)
	_, err = part.Write([]byte(content))
	assert.Nil(t, err)
	assert.Nil(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	e.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response[any]
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

// 把 resp.Data 转换成具体类型, 便于断言
func decodeData[T any](t *testing.T, data any, out *T) {
	t.Helper()
	raw, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Nil(t, json.Unmarshal(raw, out))
}

// 路径参数拼接用, 避免每个测试都写 strconv
func itoa(i int) string {
	return strconv.Itoa(i)
}
