package model

import (
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func seedArticle(t *testing.T, db *gorm.DB, a Article) *Article {
	t.Helper()
	assert.Nil(t, db.Create(&a).Error)
	return &a
}

// 后台文章列表的各种过滤条件
func TestGetArticleListFilters(t *testing.T) {
	db := newModelDB(t)
	category, _ := SaveOrUpdateCategory(db, 0, "后端")
	tag, _ := SaveOrUpdateTag(db, 0, "Go")

	public := seedArticle(t, db, Article{
		Title: "公开文章", Status: STATUS_PUBLIC, Type: TYPE_ORIGINAL,
		CategoryId: category.ID, Tags: []*Tag{tag},
	})
	seedArticle(t, db, Article{Title: "草稿", Status: STATUS_DRAFT, Type: TYPE_REPRINT})
	seedArticle(t, db, Article{Title: "回收站文章", Status: STATUS_PUBLIC, IsDelete: true})

	// 无条件
	_, total, err := GetArticleList(db, 1, 10, "", nil, 0, 0, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), total)

	// 标题模糊
	list, total, err := GetArticleList(db, 1, 10, "公开", nil, 0, 0, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, public.ID, list[0].ID)

	// 回收站
	isDelete := true
	list, total, err = GetArticleList(db, 1, 10, "", &isDelete, 0, 0, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "回收站文章", list[0].Title)

	// 状态 / 类型 / 分类 / 标签
	_, total, _ = GetArticleList(db, 1, 10, "", nil, STATUS_DRAFT, 0, 0, 0)
	assert.Equal(t, int64(1), total)
	_, total, _ = GetArticleList(db, 1, 10, "", nil, 0, TYPE_REPRINT, 0, 0)
	assert.Equal(t, int64(1), total)
	_, total, _ = GetArticleList(db, 1, 10, "", nil, 0, 0, category.ID, 0)
	assert.Equal(t, int64(1), total)
	list, total, _ = GetArticleList(db, 1, 10, "", nil, 0, 0, 0, tag.ID)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, public.ID, list[0].ID)
}

// 前台列表只出公开且未删除的文章, 置顶优先
func TestGetBlogArticleList(t *testing.T) {
	db := newModelDB(t)
	category, _ := SaveOrUpdateCategory(db, 0, "后端")
	tag, _ := SaveOrUpdateTag(db, 0, "Go")

	seedArticle(t, db, Article{Title: "普通", Status: STATUS_PUBLIC, CategoryId: category.ID, Tags: []*Tag{tag}})
	top := seedArticle(t, db, Article{Title: "置顶", Status: STATUS_PUBLIC, IsTop: true})
	seedArticle(t, db, Article{Title: "私密", Status: STATUS_SECRET})
	seedArticle(t, db, Article{Title: "已删除", Status: STATUS_PUBLIC, IsDelete: true})

	list, total, err := GetBlogArticleList(db, 1, 10, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, top.ID, list[0].ID, "置顶文章排最前")

	_, total, err = GetBlogArticleList(db, 1, 10, category.ID, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)

	_, total, err = GetBlogArticleList(db, 1, 10, 0, tag.ID)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
}

// 归档按发布时间倒序, 不受置顶影响, 也不跟着 id 走
//
// 曾经的 bug: 归档直接复用 GetBlogArticleList("is_top DESC, id DESC"),
// 于是置顶文章跑到时间轴最前, 而 id 顺序和发布时间不一致时整条时间轴都是乱的
func TestGetBlogArticleArchiveList(t *testing.T) {
	db := newModelDB(t)

	// 刻意让 id 递增方向和发布时间相反
	newest := seedArticle(t, db, Article{
		Title: "最新", Status: STATUS_PUBLIC,
		Model: Model{CreatedAt: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
	})
	middle := seedArticle(t, db, Article{
		Title: "中间", Status: STATUS_PUBLIC, IsTop: true,
		Model: Model{CreatedAt: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
	})
	oldest := seedArticle(t, db, Article{
		Title: "最早", Status: STATUS_PUBLIC,
		Model: Model{CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	})
	seedArticle(t, db, Article{Title: "私密", Status: STATUS_SECRET})
	seedArticle(t, db, Article{Title: "已删除", Status: STATUS_PUBLIC, IsDelete: true})

	list, total, err := GetBlogArticleArchiveList(db, 1, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), total, "私密和已删除的不进归档")
	assert.Equal(t,
		[]int{newest.ID, middle.ID, oldest.ID},
		[]int{list[0].ID, list[1].ID, list[2].ID},
		"按时间倒序, 置顶的那篇(中间)不该被顶到最前",
	)

	// 分页也按同一个顺序切
	page2, _, err := GetBlogArticleArchiveList(db, 2, 2)
	assert.Nil(t, err)
	assert.Len(t, page2, 1)
	assert.Equal(t, oldest.ID, page2[0].ID)
}

func TestGetArticleDetail(t *testing.T) {
	db := newModelDB(t)
	category, _ := SaveOrUpdateCategory(db, 0, "后端")
	tag, _ := SaveOrUpdateTag(db, 0, "Go")
	article := seedArticle(t, db, Article{
		Title: "文章", Status: STATUS_PUBLIC, CategoryId: category.ID, Tags: []*Tag{tag},
	})
	secret := seedArticle(t, db, Article{Title: "私密", Status: STATUS_SECRET})

	got, err := GetArticle(db, article.ID)
	assert.Nil(t, err)
	assert.Equal(t, "后端", got.Category.Name)
	assert.Len(t, got.Tags, 1)

	// 前台只能拿到公开文章
	_, err = GetBlogArticle(db, article.ID)
	assert.Nil(t, err)
	_, err = GetBlogArticle(db, secret.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

// 上一篇 / 下一篇 / 最新 / 推荐
func TestArticleNavigationAndRecommend(t *testing.T) {
	db := newModelDB(t)
	tag, _ := SaveOrUpdateTag(db, 0, "Go")

	a1 := seedArticle(t, db, Article{Title: "第一篇", Status: STATUS_PUBLIC, Tags: []*Tag{tag}})
	a2 := seedArticle(t, db, Article{Title: "第二篇", Status: STATUS_PUBLIC, Tags: []*Tag{tag}})
	a3 := seedArticle(t, db, Article{Title: "第三篇", Status: STATUS_PUBLIC, Tags: []*Tag{tag}})

	last, err := GetLastArticle(db, a2.ID)
	assert.Nil(t, err)
	assert.Equal(t, a1.ID, last.ID)

	next, err := GetNextArticle(db, a2.ID)
	assert.Nil(t, err)
	assert.Equal(t, a3.ID, next.ID)

	newest, err := GetNewestList(db, 2)
	assert.Nil(t, err)
	assert.Len(t, newest, 2)

	// 同标签的其他文章
	recommend, err := GetRecommendList(db, a1.ID, 5)
	assert.Nil(t, err)
	assert.Len(t, recommend, 2)
	for _, r := range recommend {
		assert.NotEqual(t, a1.ID, r.ID, "推荐里不能有自己")
	}
}

func TestArticleTopSoftDeleteAndDelete(t *testing.T) {
	db := newModelDB(t)
	tag, _ := SaveOrUpdateTag(db, 0, "Go")
	article := seedArticle(t, db, Article{Title: "文章", Status: STATUS_PUBLIC, Tags: []*Tag{tag}})

	assert.Nil(t, UpdateArticleTop(db, article.ID, true))
	got, _ := GetArticle(db, article.ID)
	assert.True(t, got.IsTop)

	rows, err := UpdateArticleSoftDelete(db, []int{article.ID}, true)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)
	got, _ = GetArticle(db, article.ID)
	assert.True(t, got.IsDelete)
	// 只改 is_delete, 置顶状态保留 (前台列表已经过滤了回收站文章)
	assert.True(t, got.IsTop)

	// 物理删除会一并清掉标签关联
	rows, _, err = DeleteArticle(db, []int{article.ID})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)

	var count int64
	db.Model(&ArticleTag{}).Where("article_id = ?", article.ID).Count(&count)
	assert.Zero(t, count)
}

// 物理删除文章要连带删掉文章下的评论(含回复), 并把评论 id 交给调用方清理 Redis
func TestDeleteArticleRemovesComments(t *testing.T) {
	db := newModelDB(t)
	gone := seedArticle(t, db, Article{Title: "待删除", Status: STATUS_PUBLIC})
	kept := seedArticle(t, db, Article{Title: "保留", Status: STATUS_PUBLIC})

	comment, err := AddComment(db, 1, TYPE_ARTICLE, gone.ID, "评论", true)
	assert.Nil(t, err)
	reply, err := ReplyComment(db, 2, 1, comment.ID, "回复", true)
	assert.Nil(t, err)
	keptComment, err := AddComment(db, 1, TYPE_ARTICLE, kept.ID, "别动我", true)
	assert.Nil(t, err)
	linkComment, err := AddComment(db, 1, TYPE_LINK, 0, "友链留言", true)
	assert.Nil(t, err)

	rows, commentIds, err := DeleteArticle(db, []int{gone.ID})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)
	assert.ElementsMatch(t, []int{comment.ID, reply.ID}, commentIds)

	var ids []int
	assert.Nil(t, db.Model(&Comment{}).Pluck("id", &ids).Error)
	assert.ElementsMatch(t, []int{keptComment.ID, linkComment.ID}, ids,
		"其他文章的评论和友链留言不受影响")
}

// 迁移要真的把热点查询的索引建出来, 漏掉就是全表扫描
func TestHotQueryIndexesCreated(t *testing.T) {
	db := newModelDB(t)
	m := db.Migrator()

	cases := []struct {
		model any
		index string
	}{
		{&Article{}, "idx_article_list"},     // is_delete + status
		{&Article{}, "idx_article_category"}, // 按分类筛选
		{&Comment{}, "idx_comment_topic"},    // type + topic_id
		{&Comment{}, "idx_comment_parent"},   // 查回复
		{&Resource{}, "idx_resource_api"},    // url + method, 每个请求都要查
	}
	for _, c := range cases {
		assert.True(t, m.HasIndex(c.model, c.index), "缺少索引: "+c.index)
	}
}

// 导入 markdown: 草稿状态, 不带分类和标签
func TestImportArticle(t *testing.T) {
	db := newModelDB(t)
	user := UserAuth{Username: "admin", Password: "x", UserInfo: &UserInfo{Nickname: "管理员"}}
	assert.Nil(t, db.Create(&user).Error)

	article, err := ImportArticle(db, user.UserInfo.ID, "标题", "内容", "cover.png")
	assert.Nil(t, err)
	assert.NotZero(t, article.ID)

	list, total, err := GetArticleList(db, 1, 10, "", nil, 0, 0, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "标题", list[0].Title)
	assert.Equal(t, "cover.png", list[0].Img)
	assert.Equal(t, user.UserInfo.ID, list[0].UserId)
	assert.Equal(t, STATUS_DRAFT, list[0].Status, "导入的文章是草稿")

	// 不再硬编码分类标签, 也不会凭空建出来
	assert.Zero(t, list[0].CategoryId)

	names, err := GetTagNamesByArticleId(db, list[0].ID)
	assert.Nil(t, err)
	assert.Empty(t, names)

	var count int64
	assert.Nil(t, db.Model(&Category{}).Count(&count).Error)
	assert.Zero(t, count)
	assert.Nil(t, db.Model(&Tag{}).Count(&count).Error)
	assert.Zero(t, count)
}
