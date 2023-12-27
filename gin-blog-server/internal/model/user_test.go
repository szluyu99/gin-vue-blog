package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateUserPassword(t *testing.T) {
	db := initDB()
	db.AutoMigrate(&UserAuth{}, &UserInfo{})

	var auth = UserAuth{
		Username: "test",
		Password: "123456",
	}
	db.Create(&auth)

	// 测试正常修改
	err := UpdateUserPassword(db, auth.ID, "654321")
	assert.Nil(t, err)

	// 测试不存在的用户
	err = UpdateUserPassword(db, auth.ID, "654321")
	assert.Nil(t, err)
}

func TestGetUserAuthInfoById(t *testing.T) {
	db := initDB()
	db.AutoMigrate(UserAuth{}, UserInfo{})

	var userAuth = UserAuth{
		Username: "test",
		Password: "123456",
	}
	db.Create(&userAuth)

	{
		val, err := GetUserAuthInfoById(db, userAuth.ID)
		assert.Nil(t, err)
		assert.Equal(t, "test", val.Username)
	}

	// 测试不存在的用户
	{
		val, err := GetUserAuthInfoById(db, -99)
		assert.Nil(t, val)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	}
}

func TestUpdateUserInfo(t *testing.T) {
	db := initDB()
	db.AutoMigrate(UserAuth{}, UserInfo{})

	userInfo := UserInfo{
		Nickname: "nickname",
		Avatar:   "avatar",
		Intro:    "intro",
	}
	db.Create(&userInfo)

	// 测试正常修改
	err := UpdateUserInfo(db, userInfo.ID, "update_nickname", "update_avatar", "intro", "website")
	assert.Nil(t, err)

	db.First(&userInfo, userInfo.ID)
	assert.Equal(t, "update_nickname", userInfo.Nickname)
	assert.Equal(t, "update_avatar", userInfo.Avatar)
	assert.Equal(t, "intro", userInfo.Intro)
}
