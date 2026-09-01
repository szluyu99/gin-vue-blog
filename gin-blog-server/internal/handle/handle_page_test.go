package handle

import (
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 页面接口: 覆盖 "先读缓存, 缓存没有再查库并回填, 写操作清缓存" 这条链路
func TestPageAPICache(t *testing.T) {
	env := newTestEnv(t)
	api := Page{}
	env.engine.GET("/page/list", api.GetList)
	env.engine.POST("/page", api.SaveOrUpdate)
	env.engine.DELETE("/page", api.Delete)

	// 初始没有数据, 缓存也应为空
	resp := env.do(t, http.MethodGet, "/page/list", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	var pages []model.Page
	decodeData(t, resp.Data, &pages)
	assert.Empty(t, pages)

	// 新增页面
	resp = env.do(t, http.MethodPost, "/page", map[string]any{
		"name": "home", "label": "首页", "cover": "cover.png",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 写操作后缓存应被清除
	assert.False(t, env.mr.Exists(g.PAGE))

	// 读取列表: 从库里查出来并回填缓存
	resp = env.do(t, http.MethodGet, "/page/list", nil)
	decodeData(t, resp.Data, &pages)
	assert.Len(t, pages, 1)
	assert.Equal(t, "home", pages[0].Name)
	assert.True(t, env.mr.Exists(g.PAGE))

	// 手动改掉缓存内容, 再次请求应命中缓存而不是查库
	cached, err := json.Marshal([]model.Page{{Name: "from-cache"}})
	assert.Nil(t, err)
	assert.Nil(t, env.mr.Set(g.PAGE, string(cached)))

	resp = env.do(t, http.MethodGet, "/page/list", nil)
	decodeData(t, resp.Data, &pages)
	assert.Len(t, pages, 1)
	assert.Equal(t, "from-cache", pages[0].Name)

	// 删除后缓存再次被清除
	var all []model.Page
	env.db.Find(&all)
	resp = env.do(t, http.MethodDelete, "/page", []int{all[0].ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.False(t, env.mr.Exists(g.PAGE))
}

// 博客配置接口: 配置以 hash 存在 Redis, 更新后需要失效
func TestBlogInfoConfigAPI(t *testing.T) {
	env := newTestEnv(t)
	api := BlogInfo{}
	env.engine.GET("/config", api.GetConfigMap)
	env.engine.PATCH("/config", api.UpdateConfig)

	// 配置表为空时不应报 Redis 错误(HMSET 不接受空参数)
	resp := env.do(t, http.MethodGet, "/config", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	// CheckConfigMap 只更新已存在的配置项, 所以要先有这一行
	assert.Nil(t, env.db.Create(&model.Config{Key: g.CONFIG_ABOUT, Value: "旧的关于我"}).Error)

	// 更新配置
	resp = env.do(t, http.MethodPatch, "/config", map[string]any{
		g.CONFIG_ABOUT: "关于我",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 数据库中已写入
	assert.Equal(t, "关于我", model.GetConfig(env.db, g.CONFIG_ABOUT))

	// 读取配置, 并回填缓存
	resp = env.do(t, http.MethodGet, "/config", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	var conf map[string]string
	decodeData(t, resp.Data, &conf)
	assert.Equal(t, "关于我", conf[g.CONFIG_ABOUT])
	assert.True(t, env.mr.Exists(g.CONFIG))

	// 更新后缓存被清除
	env.do(t, http.MethodPatch, "/config", map[string]any{g.CONFIG_ABOUT: "新的关于我"})
	assert.False(t, env.mr.Exists(g.CONFIG))
}

// 上报访客信息: 同一个访客只统计一次
func TestBlogInfoReport(t *testing.T) {
	env := newTestEnv(t)
	api := BlogInfo{}
	env.engine.POST("/report", api.Report)

	resp := env.do(t, http.MethodPost, "/report", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	count, err := env.rdb.Get(rctx, g.VIEW_COUNT).Int()
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	// 相同 IP + UA 再次上报, 访问量不变
	env.do(t, http.MethodPost, "/report", nil)
	count, _ = env.rdb.Get(rctx, g.VIEW_COUNT).Int()
	assert.Equal(t, 1, count)

	// 测试环境读不到 ip2region 数据库, 地域会落到 "未知"
	area := env.rdb.HGetAll(rctx, g.VISITOR_AREA).Val()
	assert.Equal(t, "1", area["未知"])
}
