package middleware

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils/jwt"
	"net/http"
	"net/http/httptest"
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
中间件测试脚手架

- 数据库: sqlite 内存库, 每个测试独立, 用真实的 model 迁移
- Redis: miniredis 内存实现, 不需要本地装 Redis
- 路由: 由每个用例自己决定挂哪些中间件, handler 只记录"有没有被执行到"
*/

type mwEnv struct {
	t      *testing.T
	engine *gin.Engine
	db     *gorm.DB
	rdb    *redis.Client
	mr     *miniredis.Miniredis

	// 请求是否走到了最终的 handler, 用来判断中间件有没有拦住
	handlerRan bool
	// handler 里从 gin context 拿到的登录用户, nil 说明中间件没识别出用户
	handlerUser *model.UserAuth
	// 大于 0 时在中间件链最前面把用户写进 session, 模拟已登录
	loginId int
}

func newMwEnv(t *testing.T) *mwEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// cache=shared + 每个测试独立的库名: file::memory: 会让连接池里的每条连接
	// 各自持有一个空库, 换连接就会报 "no such table"
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

	env := &mwEnv{t: t, db: db, rdb: rdb, mr: mr}

	engine := gin.New()
	engine.Use(sessions.Sessions("test-session", memstore.NewStore([]byte("test-secret"))))
	engine.Use(func(c *gin.Context) {
		c.Set(g.CTX_DB, db)
		c.Set(g.CTX_RDB, rdb)
		if env.loginId > 0 {
			session := sessions.Default(c)
			session.Set(g.CTX_USER_AUTH, env.loginId)
			_ = session.Save()
		}
		c.Next()
	})
	env.engine = engine

	return env
}

// 在 /api 下注册一个被 mws 保护的路由, handler 只标记自己被执行过
func (e *mwEnv) handle(method, path string, mws ...gin.HandlerFunc) {
	e.t.Helper()
	api := e.engine.Group("/api")
	api.Use(mws...)
	api.Handle(method, path, func(c *gin.Context) {
		e.handlerRan = true
		if user, exist := c.Get(g.CTX_USER_AUTH); exist {
			e.handlerUser, _ = user.(*model.UserAuth)
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
	})
}

func (e *mwEnv) request(method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	e.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	e.engine.ServeHTTP(w, req)
	return w
}

func (e *mwEnv) get(path string) *httptest.ResponseRecorder {
	return e.request(http.MethodGet, path, "", nil)
}

// 造一条资源记录 (JWTAuth / PermissionCheck 都以此表为准)
func (e *mwEnv) seedResource(name, url, method string, anonymous bool) *model.Resource {
	e.t.Helper()
	r := model.Resource{Name: name, Url: url, Method: method, Anonymous: anonymous}
	assert.Nil(e.t, e.db.Create(&r).Error)
	return &r
}

// 造一个角色, 并把指定资源挂到该角色下
func (e *mwEnv) seedRole(name string, resourceIds ...int) *model.Role {
	e.t.Helper()
	role := model.Role{Name: name, Label: name}
	assert.Nil(e.t, e.db.Create(&role).Error)
	for _, rid := range resourceIds {
		assert.Nil(e.t, e.db.Create(&model.RoleResource{RoleId: role.ID, ResourceId: rid}).Error)
	}
	return &role
}

// 造一个用户 (带 UserInfo, 操作日志中间件会读昵称), 并绑定角色
func (e *mwEnv) seedUser(username string, isSuper bool, roleIds ...int) *model.UserAuth {
	e.t.Helper()
	info := model.UserInfo{Nickname: username + "-nickname", Email: username + "@test.com"}
	assert.Nil(e.t, e.db.Create(&info).Error)

	user := model.UserAuth{Username: username, Password: "x", IsSuper: isSuper, UserInfoId: info.ID}
	assert.Nil(e.t, e.db.Create(&user).Error)
	for _, rid := range roleIds {
		assert.Nil(e.t, e.db.Create(&model.UserAuthRole{UserAuthId: user.ID, RoleId: rid}).Error)
	}
	return &user
}

// JWTAuth 依赖全局配置里的密钥, 用例结束后还原
func withJWTConf(t *testing.T) {
	t.Helper()
	old := g.Conf
	g.Conf = &g.Config{}
	g.Conf.JWT.Secret = "test-secret"
	g.Conf.JWT.Issuer = "test-issuer"
	g.Conf.JWT.Expire = 24
	t.Cleanup(func() { g.Conf = old })
}

func bearer(t *testing.T, userId, expireHour int) map[string]string {
	t.Helper()
	return bearerWithSecret(t, g.Conf.JWT.Secret, userId, expireHour)
}

func bearerWithSecret(t *testing.T, secret string, userId, expireHour int) map[string]string {
	t.Helper()
	token, err := jwt.GenToken(secret, g.Conf.JWT.Issuer, expireHour, userId, nil)
	assert.Nil(t, err)
	return map[string]string{"Authorization": "Bearer " + token}
}
