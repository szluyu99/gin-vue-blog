package handle

import (
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"
	"time"

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

/*
按天的访问量不做访客去重

累计值统计"来过多少人"(同一访客只算一次), 趋势要看的是"每天来了多少次"。
所以上面同一访客第二次上报累计值不变, 而按天的会 +1。
*/
func TestBlogInfoReportCountsEveryVisitPerDay(t *testing.T) {
	env := newTestEnv(t)
	env.engine.POST("/report", (&BlogInfo{}).Report)

	dayKey := g.VIEW_COUNT_DAY + time.Now().Format(viewDayLayout)

	env.do(t, http.MethodPost, "/report", nil)
	env.do(t, http.MethodPost, "/report", nil)
	env.do(t, http.MethodPost, "/report", nil)

	// 累计值只认第一次(同一个访客指纹), 按天的三次都算
	count, _ := env.rdb.Get(rctx, g.VIEW_COUNT).Int()
	assert.Equal(t, 1, count)
	dayCount, err := env.rdb.Get(rctx, dayKey).Int()
	assert.Nil(t, err)
	assert.Equal(t, 3, dayCount)

	// 必须带上 TTL: Incr 建出来的 key 默认永不过期, 30 天后就是永久垃圾
	ttl := env.rdb.TTL(rctx, dayKey).Val()
	assert.Greater(t, ttl, time.Duration(0), "按天的 key 要有过期时间")
	assert.LessOrEqual(t, ttl, viewDayTTL)

	// 过期之后这一天就没了, 趋势图上补 0
	env.mr.FastForward(viewDayTTL + time.Second)
	assert.False(t, env.mr.Exists(dayKey))
}

// 趋势: 固定天数, 按日期升序, 没数据的那天补 0 而不是跳过
func TestGetViewTrend(t *testing.T) {
	env := newTestEnv(t)
	now := time.Now()

	today := now.Format(viewDayLayout)
	twoDaysAgo := now.AddDate(0, 0, -2).Format(viewDayLayout)
	env.rdb.Set(rctx, g.VIEW_COUNT_DAY+today, 7, 0)
	env.rdb.Set(rctx, g.VIEW_COUNT_DAY+twoDaysAgo, 3, 0)
	// 超出窗口的那天不该出现在结果里
	env.rdb.Set(rctx, g.VIEW_COUNT_DAY+now.AddDate(0, 0, -viewTrendDays).Format(viewDayLayout), 99, 0)

	trend, err := getViewTrend(env.rdb, now)
	assert.Nil(t, err)
	assert.Len(t, trend, viewTrendDays)

	// 升序: 最后一项是今天
	assert.Equal(t, today, trend[viewTrendDays-1].Date)
	assert.Equal(t, 7, trend[viewTrendDays-1].Count)
	assert.Equal(t, 3, trend[viewTrendDays-3].Count)
	assert.Equal(t, 0, trend[viewTrendDays-2].Count, "没有数据的那天要补 0")

	for _, day := range trend {
		assert.NotEqual(t, 99, day.Count, "窗口外的日期不该被带进来")
	}
}

// 访客地域: 按人数倒序, 坏数据跳过而不是让整个仪表盘失败
func TestGetVisitorArea(t *testing.T) {
	env := newTestEnv(t)

	env.rdb.HSet(rctx, g.VISITOR_AREA, map[string]any{
		"江苏": 5, "广东": 12, "北京": 12, "未知": 1, "坏数据": "abc",
	})

	area, err := getVisitorArea(env.rdb)
	assert.Nil(t, err)

	// 人数相同(广东/北京 都是 12)时按名字排: Hash 遍历顺序不保证,
	// 不定序的话每次刷新仪表盘这几项都会跳
	assert.Equal(t, []VisitorAreaVO{
		{Area: "北京", Count: 12},
		{Area: "广东", Count: 12},
		{Area: "江苏", Count: 5},
		{Area: "未知", Count: 1},
	}, area)
}

func TestGetVisitorAreaLimitAndEmpty(t *testing.T) {
	env := newTestEnv(t)

	// 一条数据都没有时返回空而不是报错
	area, err := getVisitorArea(env.rdb)
	assert.Nil(t, err)
	assert.Empty(t, area)

	for i := 0; i < visitorAreaTop+3; i++ {
		env.rdb.HSet(rctx, g.VISITOR_AREA, "省"+itoa(i), i+1)
	}

	area, err = getVisitorArea(env.rdb)
	assert.Nil(t, err)
	assert.Len(t, area, visitorAreaTop)
	// 截断保留的是人数最多的那些
	assert.Equal(t, visitorAreaTop+3, area[0].Count)
}
