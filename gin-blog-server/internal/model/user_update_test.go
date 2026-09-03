package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 按用户名查询必须是等值匹配: 原来是 LIKE, 通配符能匹配到别人的账号
func TestGetUserAuthInfoByNameIsExactMatch(t *testing.T) {
	db := newModelDB(t)
	info := UserInfo{Nickname: "管理员", Email: "admin@qq.com"}
	assert.Nil(t, db.Create(&info).Error)
	assert.Nil(t, db.Create(&UserAuth{Username: "admin", Password: "x", UserInfoId: info.ID}).Error)

	got, err := GetUserAuthInfoByName(db, "admin")
	assert.Nil(t, err)
	assert.Equal(t, "admin", got.Username)

	// `_` 与 `%` 在 LIKE 里是通配符, 等值匹配下都查不到
	for _, name := range []string{"admi_", "adm%", "%", "_____"} {
		_, err := GetUserAuthInfoByName(db, name)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "用户名 %q 不应该匹配到 admin", name)
	}
}

// 空字段不参与更新: 调用方漏传一个字段不能等于把它清空
func TestUpdateUserInfoKeepsEmptyFields(t *testing.T) {
	db := newModelDB(t)
	assert.Nil(t, db.Create(&UserInfo{
		Model:    Model{ID: 1},
		Nickname: "原昵称",
		Avatar:   "public/uploaded/a.jpg",
		Intro:    "原简介",
		Website:  "https://a.com",
	}).Error)

	assert.Nil(t, UpdateUserInfo(db, 1, "新昵称", "", "", ""))

	var info UserInfo
	assert.Nil(t, db.First(&info, 1).Error)
	assert.Equal(t, "新昵称", info.Nickname)
	assert.Equal(t, "public/uploaded/a.jpg", info.Avatar)
	assert.Equal(t, "原简介", info.Intro)
	assert.Equal(t, "https://a.com", info.Website)

	// 全空时不发 SQL, 也不报错
	assert.Nil(t, UpdateUserInfo(db, 1, "", "", "", ""))
	assert.Nil(t, db.First(&info, 1).Error)
	assert.Equal(t, "新昵称", info.Nickname)
}
