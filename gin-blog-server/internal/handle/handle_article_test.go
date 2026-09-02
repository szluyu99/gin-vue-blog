package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newArticleEnv(t *testing.T) *testEnv {
	t.Helper()

	env := newTestEnv(t)
	api := Article{}
	env.engine.GET("/article/list", api.GetList)
	env.engine.GET("/article/:id", api.GetDetail)
	env.engine.POST("/article", api.SaveOrUpdate)
	env.engine.PUT("/article/soft-delete", api.UpdateSoftDelete)
	env.engine.PUT("/article/top", api.UpdateTop)
	env.engine.DELETE("/article", api.Delete)
	env.engine.POST("/article/import", api.Import)

	// SaveOrUpdate 用 auth.UserInfoId 作为作者
	auth := env.loginAs(1, "admin")
	auth.UserInfoId = 1

	return env
}

func TestArticleSaveOrUpdate(t *testing.T) {
	env := newArticleEnv(t)

	// 新增: 分类和标签按名称自动创建
	resp := env.do(t, http.MethodPost, "/article", map[string]any{
		"title":         "第一篇",
		"content":       "内容",
		"type":          model.TYPE_ORIGINAL,
		"status":        model.STATUS_PUBLIC,
		"category_name": "后端",
		"tag_names":     []string{"Go", "Gin"},
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var created model.Article
	decodeData(t, resp.Data, &created)
	assert.NotZero(t, created.ID)
	assert.NotZero(t, created.CategoryId)

	var tagCount, cateCount int64
	env.db.Model(&model.Tag{}).Count(&tagCount)
	env.db.Model(&model.Category{}).Count(&cateCount)
	assert.Equal(t, int64(2), tagCount)
	assert.Equal(t, int64(1), cateCount)

	// 详情: 带上分类和标签
	resp = env.do(t, http.MethodGet, "/article/"+itoa(created.ID), nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var detail model.Article
	decodeData(t, resp.Data, &detail)
	assert.Equal(t, "第一篇", detail.Title)
	assert.Equal(t, "后端", detail.Category.Name)
	assert.Len(t, detail.Tags, 2)

	// 编辑: 换标签后旧关联要被清掉
	resp = env.do(t, http.MethodPost, "/article", map[string]any{
		"id":            created.ID,
		"title":         "第一篇(改)",
		"content":       "内容",
		"type":          model.TYPE_REPRINT,
		"status":        model.STATUS_PUBLIC,
		"category_name": "后端",
		"tag_names":     []string{"Go"},
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	resp = env.do(t, http.MethodGet, "/article/"+itoa(created.ID), nil)
	decodeData(t, resp.Data, &detail)
	assert.Equal(t, "第一篇(改)", detail.Title)
	assert.Equal(t, model.TYPE_REPRINT, detail.Type)
	assert.Len(t, detail.Tags, 1)

	// 缺少必填字段 / type 越界
	resp = env.do(t, http.MethodPost, "/article", map[string]any{"title": "只有标题"})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	resp = env.do(t, http.MethodPost, "/article", map[string]any{
		"title": "t", "content": "c", "type": 9, "status": model.STATUS_PUBLIC,
	})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 详情 id 不合法
	resp = env.do(t, http.MethodGet, "/article/abc", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}

// 不传标签时也要能保存, 且没有默认封面时用配置表里的封面
func TestArticleSaveWithoutTag(t *testing.T) {
	env := newArticleEnv(t)
	env.db.Create(&model.Config{Key: g.CONFIG_ARTICLE_COVER, Value: "default.png"})

	resp := env.do(t, http.MethodPost, "/article", map[string]any{
		"title":         "无标签",
		"content":       "内容",
		"type":          model.TYPE_ORIGINAL,
		"status":        model.STATUS_DRAFT,
		"category_name": "随笔",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var created model.Article
	decodeData(t, resp.Data, &created)
	assert.Equal(t, "default.png", created.Img)

	var count int64
	env.db.Model(&model.ArticleTag{}).Where("article_id", created.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestArticleSaveWithoutLogin(t *testing.T) {
	env := newArticleEnv(t)
	env.user = nil // 模拟未登录

	resp := env.do(t, http.MethodPost, "/article", map[string]any{
		"title": "t", "content": "c",
		"type": model.TYPE_ORIGINAL, "status": model.STATUS_PUBLIC,
	})
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)
}

func TestArticleList(t *testing.T) {
	env := newArticleEnv(t)

	cate := model.Category{Name: "后端"}
	env.db.Create(&cate)
	tag := model.Tag{Name: "Go"}
	env.db.Create(&tag)

	public := model.Article{Title: "公开", Status: model.STATUS_PUBLIC, Type: model.TYPE_ORIGINAL, CategoryId: cate.ID}
	secret := model.Article{Title: "私密", Status: model.STATUS_SECRET, Type: model.TYPE_REPRINT}
	deleted := model.Article{Title: "回收站", Status: model.STATUS_PUBLIC, Type: model.TYPE_ORIGINAL, IsDelete: true}
	env.db.Create(&public)
	env.db.Create(&secret)
	env.db.Create(&deleted)
	env.db.Create(&model.ArticleTag{ArticleId: public.ID, TagId: tag.ID})

	// 点赞数来自 Redis Hash, 浏览量来自 Redis ZSet
	env.rdb.HSet(rctx, g.ARTICLE_LIKE_COUNT, itoa(public.ID), "7")
	env.rdb.ZIncrBy(rctx, g.ARTICLE_VIEW_COUNT, 3, itoa(public.ID))

	var page PageResult[ArticleVO]

	resp := env.do(t, http.MethodGet, "/article/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(3), page.Total)

	// 第一条是 id 最大的(未置顶时按 id 倒序), 逐个找出公开文章校验计数
	for _, vo := range page.List {
		if vo.ID == public.ID {
			assert.Equal(t, 7, vo.LikeCount)
			assert.Equal(t, 3, vo.ViewCount)
		}
	}

	// 各种过滤条件
	cases := []struct {
		query string
		total int64
	}{
		{"?title=公开", 1},
		{"?status=" + itoa(model.STATUS_SECRET), 1},
		{"?type=" + itoa(model.TYPE_REPRINT), 1},
		{"?is_delete=true", 1},
		{"?is_delete=false", 2},
		{"?category_id=" + itoa(cate.ID), 1},
		{"?tag_id=" + itoa(tag.ID), 1},
		{"?title=不存在", 0},
	}
	for _, tc := range cases {
		resp = env.do(t, http.MethodGet, "/article/list"+tc.query, nil)
		decodeData(t, resp.Data, &page)
		assert.Equal(t, tc.total, page.Total, tc.query)
	}
}

func TestArticleSoftDeleteAndTop(t *testing.T) {
	env := newArticleEnv(t)

	article := model.Article{Title: "文章", Status: model.STATUS_PUBLIC, Type: model.TYPE_ORIGINAL}
	env.db.Create(&article)
	env.db.Create(&model.ArticleTag{ArticleId: article.ID, TagId: 1})

	// 置顶
	resp := env.do(t, http.MethodPut, "/article/top", map[string]any{"id": article.ID, "is_top": true})
	assert.Equal(t, g.SUCCESS, resp.Code)
	env.db.First(&article, article.ID)
	assert.True(t, article.IsTop)

	// 移入回收站
	resp = env.do(t, http.MethodPut, "/article/soft-delete", map[string]any{
		"ids": []int{article.ID}, "is_delete": true,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(1), resp.Data)
	env.db.First(&article, article.ID)
	assert.True(t, article.IsDelete)

	// 从回收站恢复
	resp = env.do(t, http.MethodPut, "/article/soft-delete", map[string]any{
		"ids": []int{article.ID}, "is_delete": false,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)
	env.db.First(&article, article.ID)
	assert.False(t, article.IsDelete)

	// 物理删除, 同时清掉标签关联
	resp = env.do(t, http.MethodDelete, "/article", []int{article.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(1), resp.Data)

	var count int64
	env.db.Model(&model.Article{}).Count(&count)
	assert.Equal(t, int64(0), count)
	env.db.Model(&model.ArticleTag{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

// 物理删除文章时要把 Redis 里的计数一起清掉, 否则新文章复用 id 会继承旧计数
func TestArticleDeleteCleansRedis(t *testing.T) {
	env := newArticleEnv(t)

	kept := model.Article{Title: "保留", Status: model.STATUS_PUBLIC, Type: model.TYPE_ORIGINAL}
	gone := model.Article{Title: "删除", Status: model.STATUS_PUBLIC, Type: model.TYPE_ORIGINAL}
	env.db.Create(&kept)
	env.db.Create(&gone)

	likeSetKey := g.ARTICLE_USER_LIKE_SET + "7"
	viewVisitorKey := g.ARTICLE_VIEW_VISITOR + itoa(gone.ID) + ":fingerprint"

	// 被删文章下的评论, 其点赞数据也要一起清掉
	comment, err := model.AddComment(env.db, 7, model.TYPE_ARTICLE, gone.ID, "评论", true)
	assert.Nil(t, err)
	keptComment, err := model.AddComment(env.db, 7, model.TYPE_ARTICLE, kept.ID, "别动我", true)
	assert.Nil(t, err)
	commentLikeSetKey := g.COMMENT_USER_LIKE_SET + "7"
	env.rdb.HSet(rctx, g.COMMENT_LIKE_COUNT, itoa(comment.ID), 4)
	env.rdb.HSet(rctx, g.COMMENT_LIKE_COUNT, itoa(keptComment.ID), 2)
	env.rdb.SAdd(rctx, commentLikeSetKey, itoa(comment.ID), itoa(keptComment.ID))

	env.rdb.ZIncrBy(rctx, g.ARTICLE_VIEW_COUNT, 5, itoa(gone.ID))
	env.rdb.ZIncrBy(rctx, g.ARTICLE_VIEW_COUNT, 3, itoa(kept.ID))
	env.rdb.HSet(rctx, g.ARTICLE_LIKE_COUNT, itoa(gone.ID), 2)
	env.rdb.HSet(rctx, g.ARTICLE_LIKE_COUNT, itoa(kept.ID), 1)
	env.rdb.SAdd(rctx, likeSetKey, itoa(gone.ID), itoa(kept.ID))
	env.rdb.Set(rctx, viewVisitorKey, 1, articleViewInterval)

	resp := env.do(t, http.MethodDelete, "/article", []int{gone.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 被删文章的痕迹全部清掉
	assert.Equal(t, float64(0), env.rdb.ZScore(rctx, g.ARTICLE_VIEW_COUNT, itoa(gone.ID)).Val())
	assert.False(t, env.rdb.HExists(rctx, g.ARTICLE_LIKE_COUNT, itoa(gone.ID)).Val())
	assert.False(t, env.rdb.SIsMember(rctx, likeSetKey, itoa(gone.ID)).Val())
	assert.Equal(t, int64(0), env.rdb.Exists(rctx, viewVisitorKey).Val())

	// 文章下的评论连带删除, 评论点赞数据同步清理
	var commentCount int64
	env.db.Model(&model.Comment{}).Where("topic_id = ?", gone.ID).Count(&commentCount)
	assert.Equal(t, int64(0), commentCount)
	assert.False(t, env.rdb.HExists(rctx, g.COMMENT_LIKE_COUNT, itoa(comment.ID)).Val())
	assert.False(t, env.rdb.SIsMember(rctx, commentLikeSetKey, itoa(comment.ID)).Val())

	// 其他文章的计数不受影响
	assert.Equal(t, float64(3), env.rdb.ZScore(rctx, g.ARTICLE_VIEW_COUNT, itoa(kept.ID)).Val())
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, itoa(kept.ID)).Val())
	assert.True(t, env.rdb.SIsMember(rctx, likeSetKey, itoa(kept.ID)).Val())
	assert.Equal(t, "2", env.rdb.HGet(rctx, g.COMMENT_LIKE_COUNT, itoa(keptComment.ID)).Val())
	assert.True(t, env.rdb.SIsMember(rctx, commentLikeSetKey, itoa(keptComment.ID)).Val())
}

func TestArticleImport(t *testing.T) {
	env := newArticleEnv(t)

	resp := env.upload(t, "/article/import", "导入的文章.md", "# 标题\n正文")
	assert.Equal(t, g.SUCCESS, resp.Code)

	var article model.Article
	assert.Nil(t, env.db.First(&article).Error)
	// 文件名去掉 .md 后作为标题, 导入的文章是草稿
	assert.Equal(t, "导入的文章", article.Title)
	assert.Equal(t, "# 标题\n正文", article.Content)
	assert.Equal(t, model.STATUS_DRAFT, article.Status)

	// 固定归到 "学习" 分类 + "Golang" 标签
	var cate model.Category
	assert.Nil(t, env.db.First(&cate, article.CategoryId).Error)
	assert.Equal(t, "学习", cate.Name)

	var count int64
	env.db.Model(&model.ArticleTag{}).Where("article_id", article.ID).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestArticleImportWithoutFile(t *testing.T) {
	env := newArticleEnv(t)

	resp := env.do(t, http.MethodPost, "/article/import", nil)
	assert.Equal(t, g.ErrFileReceive.Code(), resp.Code)
}

// 标签下拉选项
func TestTagGetOption(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/tag/option", (&Tag{}).GetOption)

	assert.Nil(t, env.db.Create(&model.Tag{Name: "Go"}).Error)
	assert.Nil(t, env.db.Create(&model.Tag{Name: "Vue"}).Error)

	resp := env.do(t, http.MethodGet, "/tag/option", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var options []model.OptionVO
	decodeData(t, resp.Data, &options)
	assert.Len(t, options, 2)
	assert.NotZero(t, options[0].ID)
}

// 导出文章目前由前端完成, 后端只返回成功
func TestArticleExport(t *testing.T) {
	env := newTestEnv(t)
	env.engine.POST("/article/export", (&Article{}).Export)

	resp := env.do(t, http.MethodPost, "/article/export", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Nil(t, resp.Data)
}
