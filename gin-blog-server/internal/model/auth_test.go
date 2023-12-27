package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	db := initAuthDB()
	MakeMigrate(db)

	// 添加用户信息
	userInfo := UserInfo{
		Nickname: "admin",
	}
	db.Create(&userInfo)

	// 添加用户认证信息
	userAuth := UserAuth{
		Username:   "admin",
		Password:   "123456",
		UserInfoId: userInfo.ID,
	}
	db.Create(&userAuth)

	val, err := GetUserAuthInfoById(db, userAuth.ID)
	assert.Nil(t, err)
	assert.Equal(t, userAuth.Username, val.Username)
	assert.Equal(t, userAuth.Password, val.Password)
	assert.Equal(t, userAuth.UserInfoId, val.UserInfoId)
	// preload 从 user_auth 表中获取到了 user_info 的信息
	assert.Equal(t, userInfo.Nickname, val.UserInfo.Nickname)

}
