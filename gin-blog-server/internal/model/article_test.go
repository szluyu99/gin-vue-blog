package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createRole(t *testing.T, db *gorm.DB, name, label string) *Role {
	role := Role{
		Name:  name,
		Label: label,
	}
	err := db.Create(&role).Error
	assert.Nil(t, err)
	return &role
}

func createUser(t *testing.T, db *gorm.DB, username, nickname string) *UserAuth {
	userAuth := UserAuth{
		Username: username,
		Password: "123456",
		UserInfo: &UserInfo{
			Nickname: nickname,
		},
	}
	db.Create(&userAuth)

	val, err := GetUserAuthInfoById(db, userAuth.ID)
	assert.Nil(t, err)
	assert.Equal(t, userAuth.ID, val.ID)
	assert.Equal(t, username, val.Username)
	assert.Equal(t, nickname, val.UserInfo.Nickname)
	return &userAuth
}

func TestArticlePreload(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	MakeMigrate(db)

	u1 := createUser(t, db, "u1", "n1")
	assert.NotNil(t, u1)

	article := Article{
		Title:   "title",
		Content: "content",
		UserId:  u1.ID,
		Category: &Category{
			Name: "category",
		},
		Tags: []*Tag{
			{Name: "tag1"},
			{Name: "tag2"},
		},
	}
	err := db.Create(&article).Error
	assert.Nil(t, err)

	var val Article
	err = db.Where("id = ?", article.ID).
		Preload("User").
		Preload("Category").
		Preload("Tags").
		First(&val).Error
	assert.Nil(t, err)

	assert.Equal(t, article.ID, val.ID)
	assert.Equal(t, article.Title, val.Title)
	assert.Equal(t, article.UserId, val.UserId)
	assert.Equal(t, article.Category.Name, val.Category.Name)
	assert.Equal(t, article.Tags[0].Name, val.Tags[0].Name)
	assert.Equal(t, article.Tags[1].Name, val.Tags[1].Name)
	assert.Equal(t, u1.ID, val.User.ID)
}
