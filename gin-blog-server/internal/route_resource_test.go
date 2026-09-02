package ginblog

import (
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/*
后台路由必须与资源表(model.AdminResources)一一对应。

JWTAuth 只对资源表中登记过的接口做权限校验, 漏登记的接口虽然仍要求登录
(见 middleware.JWTAuth 的 requireLogin), 但不会经过 PermissionCheck,
任何登录用户都能调用。曾经因为路由从 /menu 改成 /menu/:id 而漏掉
3 个接口, 该测试用于防止再次漂移。
*/

// 资源表中的 url 不含 /api 前缀, 与 middleware.JWTAuth 中 FullPath()[4:] 一致
const apiPrefix = "/api"

func registeredAdminRoutes() []string {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerAdminHandler(r)

	routes := make([]string, 0, len(r.Routes()))
	for _, route := range r.Routes() {
		routes = append(routes, route.Method+" "+strings.TrimPrefix(route.Path, apiPrefix))
	}
	sort.Strings(routes)
	return routes
}

func seededResources() []string {
	seeded := make([]string, 0)
	for _, module := range model.AdminResources {
		for _, item := range module.Items {
			seeded = append(seeded, item.Method+" "+item.Url)
		}
	}
	sort.Strings(seeded)
	return seeded
}

func TestAdminRoutesMatchResources(t *testing.T) {
	routes := registeredAdminRoutes()
	seeded := seededResources()

	assert.NotEmpty(t, routes, "后台路由不应为空")

	inSeed := make(map[string]bool, len(seeded))
	for _, s := range seeded {
		inSeed[s] = true
	}
	inRoutes := make(map[string]bool, len(routes))
	for _, r := range routes {
		inRoutes[r] = true
	}

	var missing, stale []string
	for _, r := range routes {
		if !inSeed[r] {
			missing = append(missing, r)
		}
	}
	for _, s := range seeded {
		if !inRoutes[s] {
			stale = append(stale, s)
		}
	}

	assert.Empty(t, missing, "以下后台接口没有在 model.AdminResources 中登记, 不会被权限校验保护")
	assert.Empty(t, stale, "以下资源在 model.AdminResources 中登记, 但已经没有对应的后台路由")
}

// 资源定义本身不应有重复项, 否则 generate-data 会插入失败
func TestAdminResourcesNoDuplicate(t *testing.T) {
	seen := make(map[string]bool)
	for _, module := range model.AdminResources {
		for _, item := range module.Items {
			key := item.Method + " " + item.Url
			assert.False(t, seen[key], "资源重复定义: "+key)
			seen[key] = true

			assert.NotEmpty(t, item.Name, "资源缺少名称: "+key)
			assert.True(t, strings.HasPrefix(item.Url, "/"), "资源 url 需要以 / 开头: "+key)
		}
	}
}

// 资源名称是 unique 字段, 重名会导致 generate-data 插入失败
func TestAdminResourcesNameUnique(t *testing.T) {
	seen := make(map[string]bool)
	for _, module := range model.AdminResources {
		assert.False(t, seen[module.Name], "模块名重复: "+module.Name)
		seen[module.Name] = true

		for _, item := range module.Items {
			assert.False(t, seen[item.Name], "资源名重复: "+item.Name)
			seen[item.Name] = true
		}
	}
}

/*
前台只读接口也必须经过 JWTAuth。

gin 的中间件只对之后注册的路由生效, 如果 base.Use(JWTAuth) 写在这些路由后面,
它们就完全不过鉴权: 坏 token 被当成匿名放行, handler 里也拿不到当前用户。
*/
func TestBlogReadRoutesGoThroughJWTAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file:blog_route_jwt?mode=memory&cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	assert.Nil(t, err)
	assert.Nil(t, model.MakeMigrate(db))

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(g.CTX_DB, db)
		c.Next()
	})
	registerBlogHandler(r)

	// 只读接口: 不带 token 可以匿名访问, 带了坏 token 必须报错
	for _, path := range []string{
		"/api/front/home",
		"/api/front/article/list",
		"/api/front/category/list",
		"/api/front/tag/list",
		"/api/front/comment/list",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "broken-token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp struct {
			Code int `json:"code"`
		}
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, g.ErrTokenType.Code(), resp.Code, path+" 没有经过 JWTAuth")
	}
}
