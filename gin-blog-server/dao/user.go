package dao

import (
	"gin-blog/model/req"
	"gin-blog/model/resp"
)

type User struct{}

func (*User) GetCount(req req.GetUsers) int64 {
	var count int64
	tx := DB.Select("COUNT(1)").Table("user_auth ua").
		Joins("LEFT JOIN user_info ui ON ua.user_info_id = ui.id")
	if req.LoginType != 0 {
		tx = tx.Where("login_type = ?", req.LoginType)
	}
	if req.Nickname != "" {
		tx = tx.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	tx.Count(&count)
	return count
}

func (*User) GetList(req req.GetUsers) []resp.UserVO {
	var list = make([]resp.UserVO, 0)

	table := DB.
		Select("id, avatar, nickname, is_disable").Table("user_info")
	if req.LoginType != 0 {
		table = table.Where("id in (SELECT user_info_id FROM user_auth WHERE login_type = ?)", req.LoginType)
	}
	if req.Nickname != "" {
		table = table.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	table = table.Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1))
	DB.
		Select("ua.id, user_info_id, avatar, nickname, login_type, ip_address, ip_source, ua.created_at, last_login_time, ui.is_disable").
		Table("(?) ui", table).
		Preload("Roles").
		Joins("LEFT JOIN user_auth ua ON ua.user_info_id = ui.id").
		Find(&list)
	return list
}
