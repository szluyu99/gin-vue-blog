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
