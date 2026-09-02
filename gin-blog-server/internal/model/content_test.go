package model

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/*
内容相关模型的测试

initModelDB 用的是 file::memory:, 事务里换连接会拿到空库,
所以这里单独用 cache=shared + 每个测试独立的库名。
*/
func newModelDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()) +
		"?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	assert.Nil(t, err)
	assert.Nil(t, MakeMigrate(db))
	return db
}

func TestCategoryCRUD(t *testing.T) {
	db := newModelDB(t)

	c, err := SaveOrUpdateCategory(db, 0, "后端")
	assert.Nil(t, err)
	assert.NotZero(t, c.ID)

	// id > 0 走更新
	_, err = SaveOrUpdateCategory(db, c.ID, "服务端")
	assert.Nil(t, err)
	got, err := GetCategoryById(db, c.ID)
	assert.Nil(t, err)
	assert.Equal(t, "服务端", got.Name)

	byName, err := GetCategoryByName(db, "服务端")
	assert.Nil(t, err)
	assert.Equal(t, c.ID, byName.ID)

	options, err := GetCategoryOption(db)
	assert.Nil(t, err)
	assert.Len(t, options, 1)

	rows, err := DeleteCategory(db, []int{c.ID})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)
}

// 分类列表只统计已发布且未删除的文章
func TestGetCategoryListArticleCount(t *testing.T) {
	db := newModelDB(t)
	c, _ := SaveOrUpdateCategory(db, 0, "后端")

	assert.Nil(t, db.Create(&Article{Title: "已发布", CategoryId: c.ID, Status: 1}).Error)
	assert.Nil(t, db.Create(&Article{Title: "草稿", CategoryId: c.ID, Status: 3}).Error)
	assert.Nil(t, db.Create(&Article{Title: "回收站", CategoryId: c.ID, Status: 1, IsDelete: true}).Error)

	list, total, err := GetCategoryList(db, 1, 10, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, list[0].ArticleCount)

	// 关键字过滤
	_, total, err = GetCategoryList(db, 1, 10, "前端")
	assert.Nil(t, err)
	assert.Zero(t, total)
}

func TestTagCRUD(t *testing.T) {
	db := newModelDB(t)

	tag, err := SaveOrUpdateTag(db, 0, "Go")
	assert.Nil(t, err)
	assert.NotZero(t, tag.ID)

	_, err = SaveOrUpdateTag(db, tag.ID, "Golang")
	assert.Nil(t, err)

	options, err := GetTagOption(db)
	assert.Nil(t, err)
	assert.Len(t, options, 1)
	assert.Equal(t, "Golang", options[0].Name)
}

// 标签列表带上每个标签关联的文章数
func TestGetTagListArticleCount(t *testing.T) {
	db := newModelDB(t)
	tag, _ := SaveOrUpdateTag(db, 0, "Go")
	article := Article{Title: "文章", Status: 1, Tags: []*Tag{tag}}
	assert.Nil(t, db.Create(&article).Error)

	list, total, err := GetTagList(db, 1, 10, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, list[0].ArticleCount)

	names, err := GetTagNamesByArticleId(db, article.ID)
	assert.Nil(t, err)
	assert.Equal(t, []string{"Go"}, names)

	_, total, err = GetTagList(db, 1, 10, "Rust")
	assert.Nil(t, err)
	assert.Zero(t, total)
}

func TestMessageCRUD(t *testing.T) {
	db := newModelDB(t)

	m1, err := SaveMessage(db, "张三", "avatar", "留言1", "127.0.0.1", "内网IP", 10, false)
	assert.Nil(t, err)
	m2, err := SaveMessage(db, "李四", "avatar", "留言2", "127.0.0.1", "内网IP", 10, true)
	assert.Nil(t, err)

	// 昵称模糊查询
	list, total, err := GetMessageList(db, 1, 10, "张", nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, m1.ID, list[0].ID)

	// 审核状态过滤
	reviewed := true
	list, total, err = GetMessageList(db, 1, 10, "", &reviewed)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, m2.ID, list[0].ID)

	rows, err := UpdateMessagesReview(db, []int{m1.ID}, true)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)
	_, total, _ = GetMessageList(db, 1, 10, "", &reviewed)
	assert.Equal(t, int64(2), total)

	rows, err = DeleteMessages(db, []int{m1.ID, m2.ID})
	assert.Nil(t, err)
	assert.Equal(t, int64(2), rows)
}

func TestFriendLinkCRUD(t *testing.T) {
	db := newModelDB(t)

	link, err := SaveOrUpdateLink(db, 0, "博客", "avatar", "https://test.com", "简介")
	assert.Nil(t, err)
	assert.NotZero(t, link.ID)

	_, err = SaveOrUpdateLink(db, link.ID, "新博客", "avatar", "https://test.com", "简介")
	assert.Nil(t, err)

	list, total, err := GetLinkList(db, 1, 10, "新")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "新博客", list[0].Name)

	_, total, err = GetLinkList(db, 1, 10, "不存在")
	assert.Nil(t, err)
	assert.Zero(t, total)
}

func TestGetOperationLogList(t *testing.T) {
	db := newModelDB(t)
	assert.Nil(t, db.Create(&OperationLog{OptModule: "文章", OptDesc: "新增文章"}).Error)
	assert.Nil(t, db.Create(&OperationLog{OptModule: "标签", OptDesc: "删除标签"}).Error)

	list, total, err := GetOperationLogList(db, 1, 10, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)

	_, total, err = GetOperationLogList(db, 1, 10, "文章")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
}

// 前台首页统计: 只算已发布未删除的文章
func TestGetFrontStatistics(t *testing.T) {
	db := newModelDB(t)
	assert.Nil(t, db.Create(&Article{Title: "已发布", Status: 1}).Error)
	assert.Nil(t, db.Create(&Article{Title: "回收站", Status: 1, IsDelete: true}).Error)
	_, _ = SaveOrUpdateCategory(db, 0, "后端")
	_, _ = SaveOrUpdateTag(db, 0, "Go")
	_, _ = SaveMessage(db, "张三", "", "留言", "", "", 10, true)
	assert.Nil(t, db.Create(&UserAuth{Username: "test", Password: "x"}).Error)

	data, err := GetFrontStatistics(db)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), data.ArticleCount)
	assert.Equal(t, int64(1), data.UserCount)
	assert.Equal(t, int64(1), data.MessageCount)
	assert.Equal(t, int64(1), data.CategoryCount)
	assert.Equal(t, int64(1), data.TagCount)
	assert.NotNil(t, data.Config)
}
