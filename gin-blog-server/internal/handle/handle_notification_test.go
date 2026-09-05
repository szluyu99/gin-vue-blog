package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 通知接口只认 token 里的用户: 未登录一律拒绝, 也不接受前端传 user_id
func TestNotificationRequireLogin(t *testing.T) {
	env := newTestEnv(t)
	api := &Notification{}
	env.engine.GET("/notification/list", api.GetList)
	env.engine.GET("/notification/unread", api.GetUnreadCount)
	env.engine.PUT("/notification/read", api.Read)
	env.engine.DELETE("/notification", api.Delete)

	cases := []struct {
		method string
		path   string
		body   any
	}{
		{http.MethodGet, "/notification/list", nil},
		{http.MethodGet, "/notification/unread", nil},
		{http.MethodPut, "/notification/read", ReadNotificationReq{}},
		{http.MethodDelete, "/notification", ReadNotificationReq{Ids: []int{1}}},
	}
	for _, c := range cases {
		resp := env.do(t, c.method, c.path, c.body)
		assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code, c.path)
	}
}

func TestNotificationListAndRead(t *testing.T) {
	env := newTestEnv(t)
	api := &Notification{}
	env.engine.GET("/notification/list", api.GetList)
	env.engine.GET("/notification/unread", api.GetUnreadCount)
	env.engine.PUT("/notification/read", api.Read)

	me := env.loginAs(1, "me")
	other := model.UserAuth{Model: model.Model{ID: 2}, Username: "other", Password: "x"}
	assert.Nil(t, env.db.Create(&other).Error)

	for i := 0; i < 2; i++ {
		assert.Nil(t, model.AddNotification(env.db, &model.Notification{
			UserId: me.ID, FromUserId: other.ID, Type: model.NOTIFY_COMMENT_REPLY, Content: "回复你了",
		}))
	}
	// 别人的通知不该出现在我的列表和未读数里
	assert.Nil(t, model.AddNotification(env.db, &model.Notification{
		UserId: other.ID, FromUserId: me.ID, Type: model.NOTIFY_COMMENT_REPLY, Content: "别人的",
	}))

	resp := env.do(t, http.MethodGet, "/notification/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	var page PageResult[model.NotificationVO]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(2), page.Total)
	assert.Len(t, page.List, 2)

	resp = env.do(t, http.MethodGet, "/notification/unread", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.EqualValues(t, 2, resp.Data)

	// ids 为空 = 全部已读
	resp = env.do(t, http.MethodPut, "/notification/read", ReadNotificationReq{})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.EqualValues(t, 2, resp.Data)

	resp = env.do(t, http.MethodGet, "/notification/unread", nil)
	assert.EqualValues(t, 0, resp.Data)

	// 别人的未读数不受影响
	count, err := model.GetUnreadNotificationCount(env.db, other.ID)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

// 越权: 传别人的通知 id, 一条也不该被改动
func TestNotificationReadOnlyOwn(t *testing.T) {
	env := newTestEnv(t)
	api := &Notification{}
	env.engine.PUT("/notification/read", api.Read)

	env.loginAs(1, "me")
	his := model.Notification{UserId: 2, FromUserId: 1, Content: "别人的"}
	assert.Nil(t, model.AddNotification(env.db, &his))

	resp := env.do(t, http.MethodPut, "/notification/read", ReadNotificationReq{Ids: []int{his.ID}})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.EqualValues(t, 0, resp.Data)

	count, _ := model.GetUnreadNotificationCount(env.db, 2)
	assert.Equal(t, int64(1), count)
}

// 发评论会顺带产生通知, 但通知写失败不能影响发评论本身
func TestSaveCommentCreatesNotification(t *testing.T) {
	env := newTestEnv(t)
	env.engine.POST("/comment", (&Front{}).SaveComment)

	author := model.UserAuth{Model: model.Model{ID: 2}, Username: "author", Password: "x"}
	assert.Nil(t, env.db.Create(&author).Error)
	article := model.Article{Title: "文章", Status: model.STATUS_PUBLIC, UserId: author.ID}
	assert.Nil(t, env.db.Create(&article).Error)
	// 评论免审核, 否则通知会被跳过
	assert.Nil(t, env.db.Create(&model.Config{Key: g.CONFIG_IS_COMMENT_REVIEW, Value: "true"}).Error)

	env.loginAs(1, "guest")
	resp := env.do(t, http.MethodPost, "/comment", FAddCommentReq{
		TopicId: article.ID, Type: model.TYPE_ARTICLE, Content: "写得不错",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	count, err := model.GetUnreadNotificationCount(env.db, author.ID)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count, "文章作者应该收到通知")

	// 自己评论自己的文章不产生通知
	env.loginAs(author.ID, "author")
	resp = env.do(t, http.MethodPost, "/comment", FAddCommentReq{
		TopicId: article.ID, Type: model.TYPE_ARTICLE, Content: "作者补充",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)
	count, _ = model.GetUnreadNotificationCount(env.db, author.ID)
	assert.Equal(t, int64(1), count)
}
