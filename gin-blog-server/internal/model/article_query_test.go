package model

import (
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
	rows, err = DeleteArticle(db, []int{article.ID})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)

	var count int64
	db.Model(&ArticleTag{}).Where("article_id = ?", article.ID).Count(&count)
	assert.Zero(t, count)
}

// 导入 markdown: 分类和标签不存在时自动创建
func TestImportArticle(t *testing.T) {
	db := newModelDB(t)
	user := UserAuth{Username: "admin", Password: "x", UserInfo: &UserInfo{Nickname: "管理员"}}
	assert.Nil(t, db.Create(&user).Error)

	assert.Nil(t, ImportArticle(db, user.ID, "标题", "内容", "", "学习", "Golang"))

	list, total, err := GetArticleList(db, 1, 10, "", nil, 0, 0, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "标题", list[0].Title)
	assert.Equal(t, STATUS_DRAFT, list[0].Status, "导入的文章是草稿")

	category, err := GetCategoryByName(db, "学习")
	assert.Nil(t, err)
	assert.Equal(t, category.ID, list[0].CategoryId)

	names, err := GetTagNamesByArticleId(db, list[0].ID)
	assert.Nil(t, err)
	assert.Equal(t, []string{"Golang"}, names)
}
