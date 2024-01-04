package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	db, _ := initModelDB()

	// 添加用户认证信息
	userAuth := UserAuth{
		Username: "admin",
		Password: "123456",
		UserInfo: &UserInfo{
			Nickname: "admin",
		},
	}
	db.Create(&userAuth)

	val, err := GetUserAuthInfoById(db, userAuth.ID)
	assert.Nil(t, err)
	assert.Equal(t, userAuth.Username, val.Username)
	assert.Equal(t, userAuth.Password, val.Password)
	assert.Equal(t, userAuth.UserInfoId, val.UserInfoId)
	// preload 从 user_auth 表中获取到了 user_info 的信息
	assert.Equal(t, userAuth.UserInfo.Nickname, val.UserInfo.Nickname)
}

func TestGetMenuListByUserId(t *testing.T) {
	db, _ := initModelDB()

	user := UserAuth{
		Username: "user",
		Password: "password",
		UserInfo: &UserInfo{
			Nickname: "nickname",
		},
		Roles: []*Role{
			{
				Name:  "role_1",
				Label: "label_1",
				Menus: []Menu{
					{Name: "menu1", Path: "/menu1"},
					{Name: "menu2", Path: "/menu2"},
				},
			},
			{
				Name:  "role_2",
				Label: "label_2",
				Menus: []Menu{
					{Name: "menu3", Path: "/menu3"},
					{Name: "menu4", Path: "/menu4"},
				},
			},
		},
	}

	db.Create(&user)

	var val UserAuth
	db.Preload("UserInfo").Preload("Roles").Preload("Roles.Menus").First(&val, user.ID)

	assert.Equal(t, user.Username, val.Username)
	assert.Equal(t, user.UserInfoId, val.UserInfoId)               // associate create
	assert.Equal(t, user.UserInfo.Nickname, val.UserInfo.Nickname) // preload UserInfo
	assert.Len(t, val.Roles, 2)                                    // preload Roles
	assert.Len(t, val.Roles[0].Menus, 2)                           // preload Roles.Menus

	menus, err := GetMenuListByUserId(db, user.ID)
	assert.Nil(t, err)
	assert.Len(t, menus, 4)

	var count int64
	db.Model(&Menu{}).Count(&count)
	assert.Equal(t, int64(4), count)

	db.Model(&RoleMenu{}).Count(&count)
	assert.Equal(t, int64(4), count)
}

func TestAssociateDelete(t *testing.T) {
	db, _ := initModelDB()

	user := UserAuth{
		Username: "user",
		Password: "password",
		UserInfo: &UserInfo{
			Nickname: "nickname",
		},
		Roles: []*Role{
			{Name: "role", Menus: []Menu{
				{Name: "menu1", Path: "/menu1"},
				{Name: "menu2", Path: "/menu2"},
			}},
		},
	}

	db.Create(&user)

	var val UserAuth
	db.Preload("UserInfo").Preload("Roles").Preload("Roles.Menus").First(&val, user.ID)

	{
		var roleMenus []RoleMenu
		db.Table("role_menu").Where("role_id = ?", user.Roles[0].ID).Find(&roleMenus)
		assert.Len(t, roleMenus, 2)

		var userAuthRole []UserAuthRole
		db.Table("user_auth_role").Where("user_auth_id = ?", user.ID).Find(&userAuthRole)
		assert.Len(t, userAuthRole, 1)

		db.Model(&val).Association("Roles").Unscoped().Clear()
		db.Table("user_auth_role").Where("user_auth_id = ?", user.ID).Find(&userAuthRole)
		assert.Len(t, userAuthRole, 0)
	}
}
