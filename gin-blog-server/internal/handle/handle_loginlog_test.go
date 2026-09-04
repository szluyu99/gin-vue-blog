package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 登录日志列表: 按关键字过滤(用户名/昵称/IP), 支持批量删除
func TestLoginLogList(t *testing.T) {
	env := newTestEnv(t)
	api := LoginLog{}
	env.engine.GET("/login/log/list", api.GetList)
	env.engine.DELETE("/login/log", api.Delete)

	logs := []model.LoginLog{
		{Username: "admin", Nickname: "管理员", IpAddress: "10.0.0.1", IpSource: "内网IP", Status: model.LOGIN_SUCCESS},
		{Username: "guest", Nickname: "游客", IpAddress: "10.0.0.2", IpSource: "内网IP", Status: model.LOGIN_FAIL, Message: "用户名或密码错误"},
		{Username: "nobody", IpAddress: "1.2.3.4", IpSource: "美国", Status: model.LOGIN_FAIL, Message: "用户名或密码错误"},
	}
	for i := range logs {
		assert.Nil(t, env.db.Create(&logs[i]).Error)
	}

	var page PageResult[model.LoginLog]
	resp := env.do(t, http.MethodGet, "/login/log/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(3), page.Total)

	// 关键字命中用户名
	decodeData(t, env.do(t, http.MethodGet, "/login/log/list?keyword=admin", nil).Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "管理员", page.List[0].Nickname)

	// 关键字命中昵称
	decodeData(t, env.do(t, http.MethodGet, "/login/log/list?keyword=游客", nil).Data, &page)
	assert.Equal(t, int64(1), page.Total)

	// 关键字命中 IP
	decodeData(t, env.do(t, http.MethodGet, "/login/log/list?keyword=1.2.3.4", nil).Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "nobody", page.List[0].Username)

	// 批量删除
	resp = env.do(t, http.MethodDelete, "/login/log", []int{logs[0].ID, logs[1].ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(2), resp.Data)

	decodeData(t, env.do(t, http.MethodGet, "/login/log/list", nil).Data, &page)
	assert.Equal(t, int64(1), page.Total)

	// 请求体不是 ID 数组
	resp = env.do(t, http.MethodDelete, "/login/log", map[string]any{"ids": 1})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}
