package model

// 用户个人信息
type UserInfo struct {
	Universal
	Email     string `gorm:"type:varchar(30);comment:邮箱" json:"email"`
	Nickname  string `gorm:"type:varchar(30);not null;comment:昵称" json:"nickname"`
	Avatar    string `gorm:"type:varchar(1024);not null;comment:头像地址" json:"avatar"`
	Intro     string `gorm:"type:varchar(255);comment:个人简介" json:"intro"`
	Website   string `gorm:"type:varchar(255);comment:个人网站" json:"website"`
	IsDisable int8   `gorm:"type:tinyint(1);comment:是否禁用(0-否 1-是)" json:"is_disable"`
}
