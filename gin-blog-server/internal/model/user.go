package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserInfo struct {
	Model
	Email    string `json:"email" gorm:"type:varchar(30)"`
	Nickname string `json:"nickname" gorm:"unique;type:varchar(30);not null"`
	Avatar   string `json:"avatar" gorm:"type:varchar(1024);not null"`
	Intro    string `json:"intro" gorm:"type:varchar(255)"`
	Website  string `json:"website" gorm:"type:varchar(255)"`
}

type UserInfoVO struct {
	UserInfo
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
}

func GetUserInfoById(db *gorm.DB, id int) (*UserInfo, error) {
	var userInfo UserInfo
	result := db.Model(&userInfo).Where("id", id).First(&userInfo)
	return &userInfo, result.Error
}

// 按用户名精确查询用户认证信息
//
// 这里必须是等值匹配: 原来写的 `username LIKE ?`, 虽然没拼 %,
// 但 MySQL 下 LIKE 的模式串里 `%` / `_` 仍然是通配符, 注册一个含 `_` 的
// 用户名(如 `admi_`)就能匹配到别人(`admin`), 登录时拿到的是别人的记录
func GetUserAuthInfoByName(db *gorm.DB, name string) (*UserAuth, error) {
	var userauth UserAuth

	result := db.Model(&userauth).Where("username = ?", name).First(&userauth)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &userauth, result.Error
}

func GetUserList(db *gorm.DB, page, size int, loginType int8, nickname, username string) (list []UserAuth, total int64, err error) {
	if loginType != 0 {
		db = db.Where("login_type = ?", loginType)
	}

	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}

	result := db.Model(&UserAuth{}).
		Joins("LEFT JOIN user_info ON user_info.id = user_auth.user_info_id").
		Where("user_info.nickname LIKE ?", "%"+nickname+"%").
		Preload("UserInfo").
		Preload("Roles").
		Count(&total).
		Scopes(Paginate(page, size)).
		Find(&list)

	return list, total, result.Error
}

// 更新用户昵称及角色信息
func UpdateUserNicknameAndRole(db *gorm.DB, authId int, nickname string, roleIds []int) error {
	userAuth, err := GetUserAuthInfoById(db, authId)
	if err != nil {
		return err
	}

	userInfo := UserInfo{
		Model:    Model{ID: userAuth.UserInfoId},
		Nickname: nickname,
	}
	result := db.Model(&userInfo).Updates(userInfo)
	if result.Error != nil {
		return result.Error
	}

	// 至少有一个角色
	if len(roleIds) == 0 {
		return nil
	}

	// 更新用户角色, 清空原本的 user_role 关系, 添加新的关系
	// 注意这里是 user_auth_id, 用 UserInfoId 会删掉别人的角色并留下自己的旧角色
	result = db.Where(UserAuthRole{UserAuthId: userAuth.ID}).Delete(UserAuthRole{})
	if result.Error != nil {
		return result.Error
	}

	var userRoles []UserAuthRole
	for _, id := range roleIds {
		userRoles = append(userRoles, UserAuthRole{
			RoleId:     id,
			UserAuthId: userAuth.ID,
		})
	}
	result = db.Create(&userRoles)

	return result.Error
}

func UpdateUserPassword(db *gorm.DB, id int, password string) error {
	userAuth := UserAuth{
		Model:    Model{ID: id},
		Password: password,
	}
	result := db.Model(&userAuth).Updates(userAuth)
	return result.Error
}

// 更新用户资料
//
// 只更新传进来的非空字段: 原来是 Select("nickname","avatar","intro","website").Updates(...),
// 显式 Select 会把零值一起写库, 调用方漏传一个字段就等于把它清空
// (前台个人中心表单没同步时, 只改昵称提交会把头像/简介/网站全清掉)
func UpdateUserInfo(db *gorm.DB, id int, nickname, avatar, intro, website string) error {
	updates := map[string]any{}
	for column, value := range map[string]string{
		"nickname": nickname,
		"avatar":   avatar,
		"intro":    intro,
		"website":  website,
	} {
		if value != "" {
			updates[column] = value
		}
	}
	if len(updates) == 0 {
		return nil
	}

	result := db.Model(&UserInfo{Model: Model{ID: id}}).Updates(updates)
	return result.Error
}

func UpdateUserDisable(db *gorm.DB, id int, isDisable bool) error {
	userAuth := UserAuth{
		Model:     Model{ID: id},
		IsDisable: isDisable,
	}
	result := db.Model(&userAuth).Select("is_disable").Updates(&userAuth)
	return result.Error
}

// 更新用户登录信息
func UpdateUserLoginInfo(db *gorm.DB, id int, ipAddress, ipSource string) error {
	now := time.Now()
	userAuth := UserAuth{
		IpAddress:     ipAddress,
		IpSource:      ipSource,
		LastLoginTime: &now,
	}

	result := db.Where("id", id).Updates(userAuth)
	return result.Error
}
