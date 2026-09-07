package handle

import (
	"net/http"
	"testing"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 匿名上报: 不登录也能报, 同类错误聚合成一行
func TestErrorLogReport(t *testing.T) {
	env := newTestEnv(t)
	api := &ErrorLog{}
	env.engine.POST("/front/error/report", api.Report)

	body := map[string]any{
		"source":  "front",
		"message": "Cannot read properties of undefined",
		"stack":   "at foo (a.js:1:1)",
		"url":     "/article/1",
	}

	resp := env.do(t, http.MethodPost, "/front/error/report", body)
	assert.Equal(t, g.SUCCESS, resp.Code, "未登录也要能上报")

	var logs []model.ErrorLog
	assert.Nil(t, env.db.Find(&logs).Error)
	require.Len(t, logs, 1)
	assert.Equal(t, 1, logs[0].Count)
	assert.Zero(t, logs[0].UserId, "匿名上报 user_id 为 0")
	assert.NotEmpty(t, logs[0].IpAddress, "IP 由后端取, 不信前端")

	// 再报两次: 还是一行, count 累加
	env.do(t, http.MethodPost, "/front/error/report", body)
	env.do(t, http.MethodPost, "/front/error/report", body)

	assert.Nil(t, env.db.Find(&logs).Error)
	require.Len(t, logs, 1)
	assert.Equal(t, 3, logs[0].Count)

	// message 必填
	resp = env.do(t, http.MethodPost, "/front/error/report", map[string]any{"url": "/x"})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}

// 登录用户上报时记下 user_id, 便于复现时找人问
func TestErrorLogReportWithUser(t *testing.T) {
	env := newTestEnv(t)
	user := env.loginAs(7, "tester")
	env.engine.POST("/front/error/report", (&ErrorLog{}).Report)

	resp := env.do(t, http.MethodPost, "/front/error/report", map[string]any{
		"message": "boom", "stack": "at a", "source": "admin",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var log model.ErrorLog
	assert.Nil(t, env.db.First(&log).Error)
	assert.Equal(t, user.ID, log.UserId)
	assert.Equal(t, model.ERR_SOURCE_ADMIN, log.Source)
}

/*
超配额时静默丢弃, 而且必须返回成功

前端的上报挂在 window.onerror 上, 这里返回错误会让它的错误处理再报一次,
形成放大循环。所以断言两件事: 不写库, 且响应码是 SUCCESS。
*/
func TestErrorLogReportQuota(t *testing.T) {
	env := newTestEnv(t)
	env.engine.POST("/front/error/report", (&ErrorLog{}).Report)

	// 每条 message 不同, 绕过聚合, 逼近配额
	for i := 0; i < errorReportQuota; i++ {
		resp := env.do(t, http.MethodPost, "/front/error/report", map[string]any{
			"message": "boom" + itoa(i), "stack": "at a" + itoa(i),
		})
		assert.Equal(t, g.SUCCESS, resp.Code)
	}

	var before int64
	assert.Nil(t, env.db.Model(&model.ErrorLog{}).Count(&before).Error)
	assert.Equal(t, int64(errorReportQuota), before)

	// 第 11 条超配额
	resp := env.do(t, http.MethodPost, "/front/error/report", map[string]any{
		"message": "超配额的", "stack": "at z",
	})
	assert.Equal(t, g.SUCCESS, resp.Code, "超配额也要返回成功, 否则前端错误处理会循环上报")

	var after int64
	assert.Nil(t, env.db.Model(&model.ErrorLog{}).Count(&after).Error)
	assert.Equal(t, before, after, "超配额的不写库")
}

// 后台列表与删除
func TestErrorLogListAndDelete(t *testing.T) {
	env := newTestEnv(t)
	api := &ErrorLog{}
	env.engine.GET("/error/log/list", api.GetList)
	env.engine.DELETE("/error/log", api.Delete)

	assert.Nil(t, model.AddErrorLog(env.db, &model.ErrorLog{
		Source: model.ERR_SOURCE_FRONT, Message: "前台炸了", Stack: "a", Url: "/article/1",
	}))
	assert.Nil(t, model.AddErrorLog(env.db, &model.ErrorLog{
		Source: model.ERR_SOURCE_ADMIN, Message: "后台炸了", Stack: "b", Url: "/setting/page",
	}))

	resp := env.do(t, http.MethodGet, "/error/log/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	var page PageResult[model.ErrorLog]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(2), page.Total)

	// 按来源筛选
	resp = env.do(t, http.MethodGet, "/error/log/list?page_num=1&page_size=10&source=admin", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	require.Len(t, page.List, 1)
	assert.Equal(t, "后台炸了", page.List[0].Message)

	// source 只接受 front/admin
	resp = env.do(t, http.MethodGet, "/error/log/list?page_num=1&page_size=10&source=xxx", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 关键字匹配页面地址
	resp = env.do(t, http.MethodGet, "/error/log/list?page_num=1&page_size=10&keyword=setting", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)

	// 删除
	resp = env.do(t, http.MethodDelete, "/error/log", []int{page.List[0].ID})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var count int64
	assert.Nil(t, env.db.Model(&model.ErrorLog{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}
