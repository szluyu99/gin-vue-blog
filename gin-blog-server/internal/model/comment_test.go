package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// test associate create
func TestAssociateCreate(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	MakeMigrate(db)

	userAuth := UserAuth{
		Username: "test",
		Password: "test",
		UserInfo: &UserInfo{
			Nickname: "test",
		},
	}

	db.Create(&userAuth)
	assert.Equal(t, 1, userAuth.ID)
	assert.Equal(t, userAuth.UserInfo.ID, userAuth.UserInfoId)
	assert.Equal(t, "test", userAuth.Username)
}

func TestGetCommentList(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	MakeMigrate(db)

	user := UserAuth{
		Username: "username",
		Password: "123456",
		UserInfo: &UserInfo{
			Nickname: "nickname",
		},
	}
	db.Create(&user)

	article := Article{Title: "title", Content: "content"}
	db.Create(&article)

	comment, _ := AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "content", true)
	_, _ = ReplyComment(db, user.ID, user.ID, comment.ID, "reply_content", true)

	data, total, err := GetCommentList(db, 1, 10, TYPE_ARTICLE, nil, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, int(total))
	assert.Equal(t, "reply_content", data[0].Content)
	assert.Equal(t, "content", data[1].Content)

	v1 := data[0]
	assert.Equal(t, "reply_content", v1.Content)
	assert.Equal(t, "username", v1.User.Username)               // preload userAuth
	assert.Equal(t, "nickname", v1.User.UserInfo.Nickname)      // preload userAuth.userInfo
	assert.Equal(t, "username", v1.ReplyUser.Username)          // preload replyUser
	assert.Equal(t, "nickname", v1.ReplyUser.UserInfo.Nickname) // preload replyUser.userInfo
	assert.Equal(t, "title", v1.Article.Title)                  // preload article
}

// 前台评论列表: 顶级评论分页生效, 回复一次查回并按父评论分组
func TestGetCommentVOList(t *testing.T) {
	db := newModelDB(t)

	user := UserAuth{Username: "username", Password: "123456", UserInfo: &UserInfo{Nickname: "nickname"}}
	assert.Nil(t, db.Create(&user).Error)
	article := Article{Title: "title"}
	assert.Nil(t, db.Create(&article).Error)
	other := Article{Title: "other"}
	assert.Nil(t, db.Create(&other).Error)

	// 3 条顶级评论, 其中第一条有 2 条回复
	first, err := AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "第一条", true)
	assert.Nil(t, err)
	_, err = AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "第二条", true)
	assert.Nil(t, err)
	_, err = AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "第三条", true)
	assert.Nil(t, err)
	_, err = ReplyComment(db, user.ID, user.ID, first.ID, "回复一", true)
	assert.Nil(t, err)
	_, err = ReplyComment(db, user.ID, user.ID, first.ID, "回复二", true)
	assert.Nil(t, err)
	// 其他文章的评论不能混进来
	_, err = AddComment(db, user.ID, TYPE_ARTICLE, other.ID, "别人的", true)
	assert.Nil(t, err)

	// 分页只返回 2 条顶级评论, total 是顶级评论总数(不含回复)
	data, total, err := GetCommentVOList(db, 1, 2, article.ID, TYPE_ARTICLE)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, data, 2, "分页参数要生效")
	assert.Equal(t, "第三条", data[0].Content, "按 id 倒序")

	// 第二页
	data, _, err = GetCommentVOList(db, 2, 2, article.ID, TYPE_ARTICLE)
	assert.Nil(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, "第一条", data[0].Content)
	assert.Equal(t, 2, data[0].ReplyCount)
	assert.Len(t, data[0].ReplyList, 2)
	assert.Equal(t, "回复二", data[0].ReplyList[0].Content, "回复也按 id 倒序")
	assert.Equal(t, "nickname", data[0].ReplyList[0].User.UserInfo.Nickname, "回复要预加载用户")

	// 没有回复的评论给空数组而不是 nil, 前端可以直接遍历
	data, _, err = GetCommentVOList(db, 1, 1, article.ID, TYPE_ARTICLE)
	assert.Nil(t, err)
	assert.NotNil(t, data[0].ReplyList)
	assert.Empty(t, data[0].ReplyList)
}

// 未审核的评论和回复都不能出现在前台
func TestGetCommentVOListOnlyReviewed(t *testing.T) {
	db := newModelDB(t)

	user := UserAuth{Username: "username", Password: "123456", UserInfo: &UserInfo{Nickname: "nickname"}}
	assert.Nil(t, db.Create(&user).Error)
	article := Article{Title: "title"}
	assert.Nil(t, db.Create(&article).Error)

	reviewed, err := AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "已审核", true)
	assert.Nil(t, err)
	_, err = AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "待审核", false)
	assert.Nil(t, err)
	_, err = ReplyComment(db, user.ID, user.ID, reviewed.ID, "已审核的回复", true)
	assert.Nil(t, err)
	_, err = ReplyComment(db, user.ID, user.ID, reviewed.ID, "待审核的回复", false)
	assert.Nil(t, err)

	data, total, err := GetCommentVOList(db, 1, 10, article.ID, TYPE_ARTICLE)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total, "total 不能把未审核的算进去")
	assert.Len(t, data, 1)
	assert.Equal(t, "已审核", data[0].Content)
	assert.Equal(t, 1, data[0].ReplyCount)
	assert.Len(t, data[0].ReplyList, 1)
	assert.Equal(t, "已审核的回复", data[0].ReplyList[0].Content)

	// 展示条数要和文章评论数(同样按 is_review 统计)对得上
	count, err := GetArticleCommentCount(db, article.ID)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), count, "已审核的顶级评论 + 回复")

	// 单独分页拉回复时也要过滤
	replies, err := GetCommentReplyList(db, reviewed.ID, 1, 10)
	assert.Nil(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, "已审核的回复", replies[0].Content)
}
