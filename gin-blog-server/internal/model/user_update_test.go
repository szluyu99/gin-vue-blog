package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
