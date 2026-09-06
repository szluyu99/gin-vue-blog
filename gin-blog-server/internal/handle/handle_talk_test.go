package handle

import (
	"net/http"
	"testing"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 后台说说: 新增 / 编辑 / 列表筛选 / 详情
func TestTalkSaveOrUpdateAndList(t *testing.T) {
	env := newTestEnv(t)
	author := env.loginAs(9, "blogger")
	assert.Nil(t, env.db.Create(&model.UserAuth{
		Model: model.Model{ID: author.ID}, Username: "blogger", Password: "x",
		UserInfo: &model.UserInfo{Nickname: "博主"},
	}).Error)

	api := &Talk{}
	env.engine.POST("/talk", api.SaveOrUpdate)
	env.engine.GET("/talk/list", api.GetList)
	env.engine.GET("/talk/:id", api.GetDetail)

	// 新增: status 不传按公开处理, 发布者取登录用户
	resp := env.do(t, http.MethodPost, "/talk", map[string]any{"content": "第一条说说"})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var talk model.Talk
	decodeData(t, resp.Data, &talk)
	assert.Equal(t, author.ID, talk.UserId, "发布者取当前登录用户, 不信前端传的")
	assert.Equal(t, model.STATUS_PUBLIC, talk.Status)

	// 内容必填
	resp = env.do(t, http.MethodPost, "/talk", map[string]any{"content": ""})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// status 只接受 1/2
	resp = env.do(t, http.MethodPost, "/talk", map[string]any{"content": "x", "status": 9})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 编辑: 改内容 + 置顶
	resp = env.do(t, http.MethodPost, "/talk", map[string]any{
		"id": talk.ID, "content": "改过的说说", "status": model.STATUS_SECRET, "is_top": true,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	resp = env.do(t, http.MethodGet, "/talk/"+itoa(talk.ID), nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	decodeData(t, resp.Data, &talk)
	assert.Equal(t, "改过的说说", talk.Content)
	assert.True(t, talk.IsTop)
	assert.Equal(t, model.STATUS_SECRET, talk.Status, "私密的在后台也要能查到")

	// 列表按状态筛选
	resp = env.do(t, http.MethodGet, "/talk/list?page_num=1&page_size=10&status=1", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	var page PageResult[model.Talk]
	decodeData(t, resp.Data, &page)
	assert.Zero(t, page.Total, "唯一一条已经改成私密了")

	resp = env.do(t, http.MethodGet, "/talk/list?page_num=1&page_size=10&status=2", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	require.Len(t, page.List, 1)
	assert.Equal(t, "博主", page.List[0].User.UserInfo.Nickname)
}

// 未登录不能发说说(后台路由上还有 JWTAuth, 这里是 handler 自己的兜底)
func TestTalkSaveRequiresLogin(t *testing.T) {
	env := newTestEnv(t)
	env.engine.POST("/talk", (&Talk{}).SaveOrUpdate)

	resp := env.do(t, http.MethodPost, "/talk", map[string]any{"content": "匿名说说"})
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)
}

// 删除说说要把它的评论一起带走
func TestTalkDelete(t *testing.T) {
	env := newTestEnv(t)
	env.loginAs(9, "blogger")
	env.engine.DELETE("/talk", (&Talk{}).Delete)

	talk, err := model.SaveOrUpdateTalk(env.db, 0, 9, "要删的", model.STATUS_PUBLIC, false)
	assert.Nil(t, err)
	_, err = model.AddComment(env.db, 9, model.TYPE_TALK, talk.ID, "评论", true)
	assert.Nil(t, err)

	resp := env.do(t, http.MethodDelete, "/talk", []int{talk.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var count int64
	assert.Nil(t, env.db.Model(&model.Comment{}).
		Where("type = ? AND topic_id = ?", model.TYPE_TALK, talk.ID).Count(&count).Error)
	assert.Zero(t, count, "说说下的评论要一起删掉")
}

// 前台: 只给公开的, 带评论数; 私密的详情当作不存在
func TestFrontTalk(t *testing.T) {
	env := newTestEnv(t)
	front := &Front{}
	env.engine.GET("/front/talk/list", front.GetTalkList)
	env.engine.GET("/front/talk/:id", front.GetTalk)

	user := model.UserAuth{Username: "blogger", Password: "x", UserInfo: &model.UserInfo{Nickname: "博主", Avatar: "a.png"}}
	assert.Nil(t, env.db.Create(&user).Error)

	pub, err := model.SaveOrUpdateTalk(env.db, 0, user.ID, "公开的", model.STATUS_PUBLIC, false)
	assert.Nil(t, err)
	secret, err := model.SaveOrUpdateTalk(env.db, 0, user.ID, "私密的", model.STATUS_SECRET, false)
	assert.Nil(t, err)
	top, err := model.SaveOrUpdateTalk(env.db, 0, user.ID, "置顶的", model.STATUS_PUBLIC, true)
	assert.Nil(t, err)

	_, err = model.AddComment(env.db, user.ID, model.TYPE_TALK, pub.ID, "已审核", true)
	assert.Nil(t, err)
	_, err = model.AddComment(env.db, user.ID, model.TYPE_TALK, pub.ID, "待审核", false)
	assert.Nil(t, err)

	resp := env.do(t, http.MethodGet, "/front/talk/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var page PageResult[model.TalkVO]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(2), page.Total, "私密的不出现在前台")
	require.Len(t, page.List, 2)
	assert.Equal(t, top.ID, page.List[0].ID, "置顶排最前")
	assert.Equal(t, "博主", page.List[0].Nickname)
	assert.Equal(t, "a.png", page.List[0].Avatar)
	assert.Equal(t, int64(1), page.List[1].CommentCount, "只算已审核的评论")

	// 详情
	resp = env.do(t, http.MethodGet, "/front/talk/"+itoa(pub.ID), nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	var vo model.TalkVO
	decodeData(t, resp.Data, &vo)
	assert.Equal(t, "公开的", vo.Content)
	assert.Equal(t, int64(1), vo.CommentCount)

	// 私密的前台取不到
	resp = env.do(t, http.MethodGet, "/front/talk/"+itoa(secret.ID), nil)
	assert.Equal(t, g.ErrDbOp.Code(), resp.Code)

	// 非法 id
	resp = env.do(t, http.MethodGet, "/front/talk/abc", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}
