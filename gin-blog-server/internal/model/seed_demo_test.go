package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 样例内容依赖 admin / guest 两个用户, 系统数据由 generate-data 灌, 这里手动造
func seedDemoUsers(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, name := range []string{"admin", "guest"} {
		auth := UserAuth{Username: name, Password: "x", UserInfo: &UserInfo{Nickname: name}}
		assert.Nil(t, db.Create(&auth).Error)
	}
}

func TestSeedDemoContent(t *testing.T) {
	db := newModelDB(t)
	seedDemoUsers(t, db)

	assert.Nil(t, SeedDemoContent(db))

	count := func(m any) int64 {
		var n int64
		assert.Nil(t, db.Model(m).Count(&n).Error)
		return n
	}

	assert.Equal(t, int64(demoArticleCount), count(&Article{}))
	assert.Equal(t, int64(3), count(&Category{}))
	assert.Equal(t, int64(6), count(&Tag{}))
	assert.Equal(t, int64(6), count(&Comment{}))
	assert.Equal(t, int64(4), count(&Message{}))
	assert.Equal(t, int64(3), count(&FriendLink{}))

	// 前台列表只认公开且未删除的文章, 数量要超过每页 9 条才能翻页
	_, total, err := GetBlogArticleList(db, 1, 9, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(demoArticleCount-2), total) // 一篇草稿 + 一篇私密
	assert.Greater(t, total, int64(9))

	// 分类、标签、作者关联都要落上, 否则前台点进去是空的
	var articles []Article
	assert.Nil(t, db.Preload("Tags").Preload("Category").Find(&articles).Error)
	for _, a := range articles {
		assert.NotZero(t, a.CategoryId, a.Title)
		assert.NotNil(t, a.Category, a.Title)
		assert.NotZero(t, a.UserId, a.Title)
		assert.NotEmpty(t, a.Tags, a.Title)
		assert.NotEmpty(t, a.Img, a.Title)
		assert.Contains(t, a.Content, "```go", a.Title)
	}

	// 归档按月分组, 文章时间必须跨多个月份
	months := make(map[string]bool)
	for _, a := range articles {
		months[a.CreatedAt.Format("2006-01")] = true
	}
	assert.Greater(t, len(months), 1)

	// 回复评论的 topic_id / type 跟随父评论
	var replies []Comment
	assert.Nil(t, db.Where("parent_id > 0").Find(&replies).Error)
	assert.Len(t, replies, 2)
	for _, r := range replies {
		assert.NotZero(t, r.TopicId)
		assert.Equal(t, TYPE_ARTICLE, r.Type)
		assert.NotZero(t, r.ReplyUserId)
	}

	// 留一条未审核的评论和留言, 后台审核流程才有东西可点
	assert.Equal(t, int64(1), func() int64 {
		var n int64
		db.Model(&Comment{}).Where("is_review = ?", false).Count(&n)
		return n
	}())
}

// 重复执行不应该再灌一遍
func TestSeedDemoContentIdempotent(t *testing.T) {
	db := newModelDB(t)
	seedDemoUsers(t, db)

	assert.Nil(t, SeedDemoContent(db))
	assert.Nil(t, SeedDemoContent(db))

	var articles, categories int64
	db.Model(&Article{}).Count(&articles)
	db.Model(&Category{}).Count(&categories)
	assert.Equal(t, int64(demoArticleCount), articles)
	assert.Equal(t, int64(3), categories)
}

// 没有系统用户时应当报错而不是灌出 user_id = 0 的脏数据
func TestSeedDemoContentWithoutUsers(t *testing.T) {
	db := newModelDB(t)

	err := SeedDemoContent(db)
	assert.NotNil(t, err)

	var articles int64
	db.Model(&Article{}).Count(&articles)
	assert.Zero(t, articles)
}
