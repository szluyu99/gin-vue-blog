package dto

import (
	"gin-blog/model/resp"
)

// Session 信息: 记录用户详细信息 + 是否被强退
type SessionInfo struct {
	UserDetailDTO
	IsOffline int `json:"is_offline"`
}

// 用户详细信息: 仅用于在后端内部进行传输
type UserDetailDTO struct {
	resp.LoginVO
	Password   string   `json:"password"`
	IsDisable  int8     `json:"is_disable"`
	Browser    string   `json:"browser"`
	OS         string   `json:"os"`
	RoleLabels []string `json:"role_labels"`
}
