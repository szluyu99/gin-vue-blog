package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentAPI(t *testing.T) {
	env := newTestEnv(t)
	api := Comment{}
	env.engine.GET("/comment/list", api.GetList)
	env.engine.PUT("/comment/review", api.UpdateReview)
	env.engine.DELETE("/comment", api.Delete)

	// 评论列表会 Preload 用户和文章
	user := model.UserInfo{Nickname: "张三", Avatar: "a.png"}
	env.db.Create(&user)
	article := model.Article{Title: "文章", Status: model.STATUS_PUBLIC}
	env.db.Create(&article)

	articleComment := model.Comment{
		UserId: user.ID, TopicId: article.ID,
		Content: "文章评论", Type: model.TYPE_ARTICLE,
	}
	linkComment := model.Comment{
		UserId:  user.ID,
		Content: "友链评论", Type: model.TYPE_LINK, IsReview: true,
	}
	env.db.Create(&articleComment)
	env.db.Create(&linkComment)

	var page PageResult[model.Comment]

	resp := env.do(t, http.MethodGet, "/comment/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(2), page.Total)

	// 按类型 / 审核状态过滤
	resp = env.do(t, http.MethodGet, "/comment/list?type="+itoa(model.TYPE_LINK), nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "友链评论", page.List[0].Content)

	resp = env.do(t, http.MethodGet, "/comment/list?is_review=false", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "文章评论", page.List[0].Content)

	// 批量通过审核
	resp = env.do(t, http.MethodPut, "/comment/review", map[string]any{
		"ids": []int{articleComment.ID}, "is_review": true,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(1), resp.Data)

	resp = env.do(t, http.MethodGet, "/comment/list?is_review=true", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(2), page.Total)

	// 批量删除
	resp = env.do(t, http.MethodDelete, "/comment", []int{articleComment.ID, linkComment.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(2), resp.Data)

	resp = env.do(t, http.MethodGet, "/comment/list", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(0), page.Total)

	// 请求体不是 ID 数组
	resp = env.do(t, http.MethodDelete, "/comment", map[string]any{"ids": 1})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}

// 删除评论要连带删掉它的回复, 并清掉 Redis 里的点赞数据,
// 否则回复变成查得到的孤儿数据, 新评论复用 id 还会继承旧计数
func TestCommentDeleteCleansRepliesAndRedis(t *testing.T) {
	env := newTestEnv(t)
	env.engine.DELETE("/comment", (&Comment{}).Delete)

	article := model.Article{Title: "文章", Status: model.STATUS_PUBLIC}
	assert.Nil(t, env.db.Create(&article).Error)
	user := model.UserAuth{Username: "u", Password: "x", UserInfo: &model.UserInfo{Nickname: "n"}}
	assert.Nil(t, env.db.Create(&user).Error)

	top, err := model.AddComment(env.db, user.ID, model.TYPE_ARTICLE, article.ID, "顶级评论", true)
	assert.Nil(t, err)
	reply, err := model.ReplyComment(env.db, user.ID, user.ID, top.ID, "回复", true)
	assert.Nil(t, err)
	// 另一条评论及其回复, 用来确认没被误删
	other, err := model.AddComment(env.db, user.ID, model.TYPE_ARTICLE, article.ID, "另一条", true)
	assert.Nil(t, err)
	otherReply, err := model.ReplyComment(env.db, user.ID, user.ID, other.ID, "另一条的回复", true)
	assert.Nil(t, err)

	likeKey := g.COMMENT_USER_LIKE_SET + itoa(user.ID)
	for _, id := range []int{top.ID, reply.ID, other.ID, otherReply.ID} {
		env.rdb.HSet(rctx, g.COMMENT_LIKE_COUNT, itoa(id), 2)
		env.rdb.SAdd(rctx, likeKey, itoa(id))
	}

	// 只删顶级评论, 它的回复应该一起消失
	resp := env.do(t, http.MethodDelete, "/comment", []int{top.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(2), resp.Data, "顶级评论 + 1 条回复")

	var left []model.Comment
	assert.Nil(t, env.db.Find(&left).Error)
	assert.Len(t, left, 2)

	var orphan int64
	assert.Nil(t, env.db.Model(&model.Comment{}).Where("parent_id = ?", top.ID).Count(&orphan).Error)
	assert.Zero(t, orphan, "不能留下 parent_id 指向已删评论的回复")

	// 被删的两条的 Redis 数据清掉了, 没被删的还在
	for _, id := range []int{top.ID, reply.ID} {
		assert.False(t, env.rdb.HExists(rctx, g.COMMENT_LIKE_COUNT, itoa(id)).Val(), itoa(id))
		assert.False(t, env.rdb.SIsMember(rctx, likeKey, itoa(id)).Val(), itoa(id))
	}
	for _, id := range []int{other.ID, otherReply.ID} {
		assert.True(t, env.rdb.HExists(rctx, g.COMMENT_LIKE_COUNT, itoa(id)).Val(), itoa(id))
		assert.True(t, env.rdb.SIsMember(rctx, likeKey, itoa(id)).Val(), itoa(id))
	}
}

func TestMessageAPI(t *testing.T) {
	env := newTestEnv(t)
	api := Message{}
	env.engine.GET("/message/list", api.GetList)
	env.engine.PUT("/message/review", api.UpdateReview)
	env.engine.DELETE("/message", api.Delete)

	first := model.Message{Nickname: "张三", Content: "你好"}
	second := model.Message{Nickname: "李四", Content: "在吗", IsReview: true}
	env.db.Create(&first)
	env.db.Create(&second)

	var page PageResult[model.Message]

	resp := env.do(t, http.MethodGet, "/message/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(2), page.Total)

	// 昵称模糊匹配
	resp = env.do(t, http.MethodGet, "/message/list?nickname=张", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "你好", page.List[0].Content)

	// 未审核的留言
	resp = env.do(t, http.MethodGet, "/message/list?is_review=false", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "张三", page.List[0].Nickname)

	resp = env.do(t, http.MethodPut, "/message/review", map[string]any{
		"ids": []int{first.ID}, "is_review": true,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(1), resp.Data)

	resp = env.do(t, http.MethodDelete, "/message", []int{first.ID, second.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(2), resp.Data)

	resp = env.do(t, http.MethodGet, "/message/list", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(0), page.Total)
}

func TestLinkAPI(t *testing.T) {
	env := newTestEnv(t)
	api := Link{}
	env.engine.GET("/link/list", api.GetList)
	env.engine.POST("/link", api.SaveOrUpdate)
	env.engine.DELETE("/link", api.Delete)

	resp := env.do(t, http.MethodPost, "/link", map[string]any{
		"name": "友链", "address": "https://example.com", "intro": "示例站点",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var created model.FriendLink
	decodeData(t, resp.Data, &created)
	assert.NotZero(t, created.ID)

	// 修改
	resp = env.do(t, http.MethodPost, "/link", map[string]any{
		"id": created.ID, "name": "友链(改)", "address": "https://example.org",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var page PageResult[model.FriendLink]
	resp = env.do(t, http.MethodGet, "/link/list?page_num=1&page_size=10", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "友链(改)", page.List[0].Name)

	// 关键字命中名称/地址/简介中的任意一个
	resp = env.do(t, http.MethodGet, "/link/list?keyword=example.org", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)

	resp = env.do(t, http.MethodGet, "/link/list?keyword=不存在", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(0), page.Total)

	// 缺少必填字段
	resp = env.do(t, http.MethodPost, "/link", map[string]any{"name": "只有名称"})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	resp = env.do(t, http.MethodDelete, "/link", []int{created.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(1), resp.Data)
}

func TestOperationLogAPI(t *testing.T) {
	env := newTestEnv(t)
	api := OperationLog{}
	env.engine.GET("/operation/log/list", api.GetList)
	env.engine.DELETE("/operation/log", api.Delete)

	first := model.OperationLog{OptModule: "文章模块", OptDesc: "新增文章", Nickname: "admin"}
	second := model.OperationLog{OptModule: "用户模块", OptDesc: "修改密码", Nickname: "admin"}
	env.db.Create(&first)
	env.db.Create(&second)

	var page PageResult[model.OperationLog]

	resp := env.do(t, http.MethodGet, "/operation/log/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(2), page.Total)

	// 关键字匹配模块名或操作描述
	resp = env.do(t, http.MethodGet, "/operation/log/list?keyword=文章", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "文章模块", page.List[0].OptModule)

	resp = env.do(t, http.MethodGet, "/operation/log/list?keyword=修改密码", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)

	resp = env.do(t, http.MethodDelete, "/operation/log", []int{first.ID, second.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(2), resp.Data)

	resp = env.do(t, http.MethodGet, "/operation/log/list", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(0), page.Total)
}
