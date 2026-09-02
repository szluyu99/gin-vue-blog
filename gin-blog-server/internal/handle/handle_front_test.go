package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 点赞文章: 同一用户重复点赞是取消, 计数与用户集合要同步变化
func TestFrontLikeArticle(t *testing.T) {
	env := newTestEnv(t)
	user := env.loginAs(7, "tester")
	env.engine.GET("/front/article/like/:article_id", (&Front{}).LikeArticle)

	likeSetKey := g.ARTICLE_USER_LIKE_SET + strconv.Itoa(user.ID)

	// 第一次: 点赞
	resp := env.do(t, http.MethodGet, "/front/article/like/1", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.True(t, env.rdb.SIsMember(rctx, likeSetKey, 1).Val())
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "1").Val())

	// 第二次: 取消点赞
	resp = env.do(t, http.MethodGet, "/front/article/like/1", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.False(t, env.rdb.SIsMember(rctx, likeSetKey, 1).Val())
	assert.Equal(t, "0", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "1").Val())

	// 不同文章互不影响
	env.do(t, http.MethodGet, "/front/article/like/2", nil)
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "2").Val())
	assert.Equal(t, "0", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "1").Val())

	// 非法的文章 id
	resp = env.do(t, http.MethodGet, "/front/article/like/abc", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}

// 点赞评论: 与文章点赞是两套独立的 key
func TestFrontLikeComment(t *testing.T) {
	env := newTestEnv(t)
	user := env.loginAs(7, "tester")
	env.engine.GET("/front/comment/like/:comment_id", (&Front{}).LikeComment)

	likeSetKey := g.COMMENT_USER_LIKE_SET + strconv.Itoa(user.ID)

	resp := env.do(t, http.MethodGet, "/front/comment/like/3", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.True(t, env.rdb.SIsMember(rctx, likeSetKey, 3).Val())
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.COMMENT_LIKE_COUNT, "3").Val())

	// 文章点赞的计数不应受影响
	assert.Empty(t, env.rdb.HGetAll(rctx, g.ARTICLE_LIKE_COUNT).Val())

	// 取消点赞
	env.do(t, http.MethodGet, "/front/comment/like/3", nil)
	assert.False(t, env.rdb.SIsMember(rctx, likeSetKey, 3).Val())
	assert.Equal(t, "0", env.rdb.HGet(rctx, g.COMMENT_LIKE_COUNT, "3").Val())
}

// 不同用户的点赞记录互相隔离
func TestFrontLikeArticleMultiUser(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/article/like/:article_id", (&Front{}).LikeArticle)

	env.loginAs(1, "user1")
	env.do(t, http.MethodGet, "/front/article/like/9", nil)

	env.loginAs(2, "user2")
	env.do(t, http.MethodGet, "/front/article/like/9", nil)

	// 两个用户各自点赞, 计数为 2
	assert.Equal(t, "2", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "9").Val())
	assert.True(t, env.rdb.SIsMember(rctx, g.ARTICLE_USER_LIKE_SET+"1", 9).Val())
	assert.True(t, env.rdb.SIsMember(rctx, g.ARTICLE_USER_LIKE_SET+"2", 9).Val())

	// user2 取消, 只影响自己
	env.do(t, http.MethodGet, "/front/article/like/9", nil)
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "9").Val())
	assert.True(t, env.rdb.SIsMember(rctx, g.ARTICLE_USER_LIKE_SET+"1", 9).Val())
	assert.False(t, env.rdb.SIsMember(rctx, g.ARTICLE_USER_LIKE_SET+"2", 9).Val())
}

// 未登录访问需要登录的接口: 应返回业务错误而不是 panic 成 500
// (JWTAuth 对资源表中不存在的接口会跳过鉴权, 因此 handler 必须自己兜住)
func TestFrontLikeWithoutLogin(t *testing.T) {
	env := newTestEnv(t) // 不调用 loginAs
	env.engine.GET("/front/article/like/:article_id", (&Front{}).LikeArticle)
	env.engine.GET("/front/comment/like/:comment_id", (&Front{}).LikeComment)

	resp := env.do(t, http.MethodGet, "/front/article/like/1", nil)
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)

	resp = env.do(t, http.MethodGet, "/front/comment/like/1", nil)
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)

	// 未登录不应产生任何点赞数据
	assert.Empty(t, env.rdb.HGetAll(rctx, g.ARTICLE_LIKE_COUNT).Val())
	assert.Empty(t, env.rdb.HGetAll(rctx, g.COMMENT_LIKE_COUNT).Val())
}

// 前台首页统计 + Redis 中的访问量
func TestFrontGetHomeInfo(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/home", (&Front{}).GetHomeInfo)

	assert.Nil(t, env.db.Create(&model.Article{Title: "公开", Status: model.STATUS_PUBLIC}).Error)
	assert.Nil(t, env.db.Create(&model.Article{Title: "草稿", Status: model.STATUS_DRAFT}).Error)
	assert.Nil(t, env.db.Create(&model.Category{Name: "后端"}).Error)
	assert.Nil(t, env.db.Create(&model.Tag{Name: "Go"}).Error)
	env.rdb.Set(rctx, g.VIEW_COUNT, 233, 0)

	resp := env.do(t, http.MethodGet, "/front/home", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var data model.FrontHomeVO
	decodeData(t, resp.Data, &data)
	assert.Equal(t, int64(1), data.ArticleCount, "只统计已发布的文章")
	assert.Equal(t, int64(1), data.CategoryCount)
	assert.Equal(t, int64(1), data.TagCount)
	assert.Equal(t, int64(233), data.ViewCount)
}

// 前台的几个只读列表接口
func TestFrontSimpleLists(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/tag/list", (&Front{}).GetTagList)
	env.engine.GET("/front/category/list", (&Front{}).GetCategoryList)
	env.engine.GET("/front/message/list", (&Front{}).GetMessageList)
	env.engine.GET("/front/link/list", (&Front{}).GetLinkList)

	assert.Nil(t, env.db.Create(&model.Tag{Name: "Go"}).Error)
	assert.Nil(t, env.db.Create(&model.Category{Name: "后端"}).Error)
	_, err := model.SaveMessage(env.db, "已审核", "", "内容", "", "", 10, true)
	assert.Nil(t, err)
	_, err = model.SaveMessage(env.db, "待审核", "", "内容", "", "", 10, false)
	assert.Nil(t, err)
	_, err = model.SaveOrUpdateLink(env.db, 0, "友链", "", "https://test.com", "简介")
	assert.Nil(t, err)

	var tags []model.TagVO
	decodeData(t, env.do(t, http.MethodGet, "/front/tag/list", nil).Data, &tags)
	assert.Len(t, tags, 1)

	var categories []model.CategoryVO
	decodeData(t, env.do(t, http.MethodGet, "/front/category/list", nil).Data, &categories)
	assert.Len(t, categories, 1)

	// 前台只展示审核通过的留言
	var messages []model.Message
	decodeData(t, env.do(t, http.MethodGet, "/front/message/list", nil).Data, &messages)
	assert.Len(t, messages, 1)
	assert.Equal(t, "已审核", messages[0].Nickname)

	var links []model.FriendLink
	decodeData(t, env.do(t, http.MethodGet, "/front/link/list", nil).Data, &links)
	assert.Len(t, links, 1)
}

// 前台文章列表: 只出公开且未删除的文章, 支持分类和标签过滤
func TestFrontGetArticleList(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/article/list", (&Front{}).GetArticleList)

	category := model.Category{Name: "后端"}
	assert.Nil(t, env.db.Create(&category).Error)
	tag := model.Tag{Name: "Go"}
	assert.Nil(t, env.db.Create(&tag).Error)

	assert.Nil(t, env.db.Create(&model.Article{
		Title: "公开", Status: model.STATUS_PUBLIC, CategoryId: category.ID, Tags: []*model.Tag{&tag},
	}).Error)
	assert.Nil(t, env.db.Create(&model.Article{Title: "私密", Status: model.STATUS_SECRET}).Error)
	assert.Nil(t, env.db.Create(&model.Article{Title: "回收站", Status: model.STATUS_PUBLIC, IsDelete: true}).Error)

	var list []model.Article
	decodeData(t, env.do(t, http.MethodGet, "/front/article/list?page_num=1&page_size=10", nil).Data, &list)
	assert.Len(t, list, 1)
	assert.Equal(t, "公开", list[0].Title)

	// 分类 / 标签过滤
	decodeData(t, env.do(t, http.MethodGet,
		"/front/article/list?page_num=1&page_size=10&category_id="+itoa(category.ID), nil).Data, &list)
	assert.Len(t, list, 1)

	decodeData(t, env.do(t, http.MethodGet,
		"/front/article/list?page_num=1&page_size=10&tag_id="+itoa(tag.ID), nil).Data, &list)
	assert.Len(t, list, 1)
}

// 文章详情: 附带上下篇/推荐/最新, 且浏览量 +1
func TestFrontGetArticleInfo(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/article/:id", (&Front{}).GetArticleInfo)

	tag := model.Tag{Name: "Go"}
	assert.Nil(t, env.db.Create(&tag).Error)
	a1 := model.Article{Title: "第一篇", Status: model.STATUS_PUBLIC, Tags: []*model.Tag{&tag}}
	a2 := model.Article{Title: "第二篇", Status: model.STATUS_PUBLIC, Tags: []*model.Tag{&tag}}
	a3 := model.Article{Title: "第三篇", Status: model.STATUS_PUBLIC, Tags: []*model.Tag{&tag}}
	assert.Nil(t, env.db.Create(&a1).Error)
	assert.Nil(t, env.db.Create(&a2).Error)
	assert.Nil(t, env.db.Create(&a3).Error)

	resp := env.do(t, http.MethodGet, "/front/article/"+itoa(a2.ID), nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var data model.BlogArticleVO
	decodeData(t, resp.Data, &data)
	assert.Equal(t, "第二篇", data.Title)
	assert.Equal(t, a1.ID, data.LastArticle.ID)
	assert.Equal(t, a3.ID, data.NextArticle.ID)
	assert.Len(t, data.RecommendArticles, 2, "同标签的其他文章")
	assert.Len(t, data.NewestArticles, 3)
	assert.Equal(t, int64(1), data.ViewCount, "访问一次浏览量 +1")

	// 同一访客再访问一次不再计数, 否则刷新就能刷量
	decodeData(t, env.do(t, http.MethodGet, "/front/article/"+itoa(a2.ID), nil).Data, &data)
	assert.Equal(t, int64(1), data.ViewCount)

	// 换一个访客(不同 IP)才继续累加
	resp = env.doWithHeader(t, http.MethodGet, "/front/article/"+itoa(a2.ID), nil,
		map[string]string{"X-Real-IP": "10.0.0.9"})
	decodeData(t, resp.Data, &data)
	assert.Equal(t, int64(2), data.ViewCount)

	// 去重窗口过期后同一访客可以再计一次
	env.mr.FastForward(articleViewInterval + time.Second)
	decodeData(t, env.do(t, http.MethodGet, "/front/article/"+itoa(a2.ID), nil).Data, &data)
	assert.Equal(t, int64(3), data.ViewCount)

	// 私密文章拿不到
	secret := model.Article{Title: "私密", Status: model.STATUS_SECRET}
	assert.Nil(t, env.db.Create(&secret).Error)
	resp = env.do(t, http.MethodGet, "/front/article/"+itoa(secret.ID), nil)
	assert.Equal(t, g.ErrDbOp.Code(), resp.Code)

	// 非法 id
	resp = env.do(t, http.MethodGet, "/front/article/abc", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}

// 归档列表
func TestFrontGetArchiveList(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/article/archive", (&Front{}).GetArchiveList)

	assert.Nil(t, env.db.Create(&model.Article{Title: "公开", Status: model.STATUS_PUBLIC}).Error)
	assert.Nil(t, env.db.Create(&model.Article{Title: "私密", Status: model.STATUS_SECRET}).Error)

	resp := env.do(t, http.MethodGet, "/front/article/archive?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var page PageResult[ArchiveVO]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "公开", page.List[0].Title)
	assert.NotZero(t, page.List[0].ID)
}

// 搜索文章: 命中标题或内容, 结果里关键字被高亮包裹
func TestFrontSearchArticle(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/article/search", (&Front{}).SearchArticle)

	assert.Nil(t, env.db.Create(&model.Article{
		Title: "Go 并发编程", Content: "goroutine 与 channel", Status: model.STATUS_PUBLIC,
	}).Error)
	assert.Nil(t, env.db.Create(&model.Article{
		Title: "中文标题", Content: "这里讲的是数据库索引的原理, 索引很重要", Status: model.STATUS_PUBLIC,
	}).Error)
	assert.Nil(t, env.db.Create(&model.Article{
		Title: "私密文章", Content: "索引", Status: model.STATUS_SECRET,
	}).Error)

	var result []ArticleSearchVO

	// 关键字为空时直接返回空列表
	decodeData(t, env.do(t, http.MethodGet, "/front/article/search", nil).Data, &result)
	assert.Empty(t, result)

	// 命中标题
	decodeData(t, env.do(t, http.MethodGet, "/front/article/search?keyword=并发", nil).Data, &result)
	assert.Len(t, result, 1)
	assert.Contains(t, result[0].Title, "<span style='color:#f47466'>并发</span>")

	// 命中内容(含中文截取), 私密文章不参与搜索
	decodeData(t, env.do(t, http.MethodGet, "/front/article/search?keyword=索引", nil).Data, &result)
	assert.Len(t, result, 1)
	assert.Contains(t, result[0].Content, "<span style='color:#f47466'>索引</span>")

	// 没有命中
	decodeData(t, env.do(t, http.MethodGet, "/front/article/search?keyword=不存在的词", nil).Data, &result)
	assert.Empty(t, result)
}

// 登录用户的昵称等信息, 生产环境由 JWTAuth 预加载
func (e *testEnv) loginAsFullUser(id int, username, nickname string) *model.UserAuth {
	e.user = &model.UserAuth{
		Model:    model.Model{ID: id},
		Username: username,
		UserInfo: &model.UserInfo{Nickname: nickname},
	}
	return e.user
}

// 新增留言: 内容做 HTML 转义, 审核开关由博客配置决定
func TestFrontSaveMessage(t *testing.T) {
	env := newTestEnv(t)
	env.loginAsFullUser(7, "tester", "测试昵称")
	env.engine.POST("/front/message", (&Front{}).SaveMessage)

	resp := env.do(t, http.MethodPost, "/front/message", map[string]any{
		"nickname": "测试昵称",
		"content":  "<script>alert(1)</script>留言",
		"speed":    10,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var message model.Message
	decodeData(t, resp.Data, &message)
	assert.Equal(t, "测试昵称", message.Nickname, "昵称取自登录用户, 不信前端传的")
	assert.NotContains(t, message.Content, "<script>", "内容要转义, 防 XSS")
	assert.False(t, message.IsReview, "没有配置时按需要审核处理")

	// 评论的审核开关不能影响留言
	assert.Nil(t, env.db.Create(&model.Config{Key: g.CONFIG_IS_COMMENT_REVIEW, Value: "true"}).Error)
	resp = env.do(t, http.MethodPost, "/front/message", map[string]any{
		"nickname": "测试昵称",
		"content":  "第二条留言",
	})
	decodeData(t, resp.Data, &message)
	assert.False(t, message.IsReview, "留言只看 is_message_review")

	// is_message_review 为 true 表示留言免审核, 直接展示 (与后台设置页的选项对应)
	assert.Nil(t, env.db.Create(&model.Config{Key: g.CONFIG_IS_MESSAGE_REVIEW, Value: "true"}).Error)
	resp = env.do(t, http.MethodPost, "/front/message", map[string]any{
		"nickname": "测试昵称",
		"content":  "第三条留言",
	})
	decodeData(t, resp.Data, &message)
	assert.True(t, message.IsReview)

	// 缺少必填字段
	resp = env.do(t, http.MethodPost, "/front/message", map[string]any{"speed": 10})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}

// 新增评论与回复评论
func TestFrontSaveComment(t *testing.T) {
	env := newTestEnv(t)
	user := env.loginAsFullUser(7, "tester", "测试昵称")
	env.engine.POST("/front/comment", (&Front{}).SaveComment)

	article := model.Article{Title: "文章", Status: model.STATUS_PUBLIC}
	assert.Nil(t, env.db.Create(&article).Error)

	// 顶级评论
	resp := env.do(t, http.MethodPost, "/front/comment", map[string]any{
		"topic_id": article.ID,
		"type":     1,
		"content":  "<b>写得不错</b>",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var comment model.Comment
	decodeData(t, resp.Data, &comment)
	assert.Equal(t, user.ID, comment.UserId)
	assert.Zero(t, comment.ParentId)
	assert.NotContains(t, comment.Content, "<b>")

	// 回复该评论
	resp = env.do(t, http.MethodPost, "/front/comment", map[string]any{
		"topic_id":      article.ID,
		"type":          1,
		"content":       "谢谢",
		"parent_id":     comment.ID,
		"reply_user_id": user.ID,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var reply model.Comment
	decodeData(t, resp.Data, &reply)
	assert.Equal(t, comment.ID, reply.ParentId)
	assert.Equal(t, user.ID, reply.ReplyUserId)

	// 未登录不能评论
	env.user = nil
	resp = env.do(t, http.MethodPost, "/front/comment", map[string]any{
		"topic_id": article.ID, "type": 1, "content": "匿名",
	})
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)
}

// 评论列表: 顶级评论最多带 3 条回复, 点赞数来自 Redis
func TestFrontGetCommentList(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/comment/list", (&Front{}).GetCommentList)
	env.engine.GET("/front/comment/replies/:comment_id", (&Front{}).GetReplyListByCommentId)

	article := model.Article{Title: "文章", Status: model.STATUS_PUBLIC}
	assert.Nil(t, env.db.Create(&article).Error)
	user := model.UserAuth{Username: "u", Password: "x", UserInfo: &model.UserInfo{Nickname: "n"}}
	assert.Nil(t, env.db.Create(&user).Error)

	// 前台只展示审核通过的评论, 所以这里造的数据都是已审核的
	top, err := model.AddComment(env.db, user.ID, 1, article.ID, "顶级评论", true)
	assert.Nil(t, err)
	for i := 0; i < 4; i++ {
		_, err = model.ReplyComment(env.db, user.ID, user.ID, top.ID, "回复"+itoa(i), true)
		assert.Nil(t, err)
	}
	env.rdb.HSet(rctx, g.COMMENT_LIKE_COUNT, itoa(top.ID), 5)

	resp := env.do(t, http.MethodGet,
		"/front/comment/list?page_num=1&page_size=10&topic_id="+itoa(article.ID)+"&type=1", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var page PageResult[model.CommentVO]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total, "只统计顶级评论")
	assert.Len(t, page.List[0].ReplyList, 3, "列表里最多带 3 条回复")
	assert.Equal(t, 5, page.List[0].LikeCount)

	// 回复列表接口能拿到全部回复
	var replies []model.CommentVO
	decodeData(t, env.do(t, http.MethodGet,
		"/front/comment/replies/"+itoa(top.ID)+"?page_num=1&page_size=10", nil).Data, &replies)
	assert.Len(t, replies, 4)

	// 非法的评论 id
	resp = env.do(t, http.MethodGet, "/front/comment/replies/abc", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 未审核的评论不出现在前台
	_, err = model.AddComment(env.db, user.ID, 1, article.ID, "待审核", false)
	assert.Nil(t, err)
	decodeData(t, env.do(t, http.MethodGet,
		"/front/comment/list?page_num=1&page_size=10&topic_id="+itoa(article.ID)+"&type=1", nil).Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Len(t, page.List, 1)
}
