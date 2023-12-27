package model

import "time"

type CommentCountDTO struct {
	ID           int `json:"id"`
	CommentCount int `json:"comment_count"`
}

type ReplyCountDTO struct {
	CommentId  int `json:"comment_id"`
	ReplyCount int `json:"reply_count"`
}

// Session 信息: 记录用户详细信息 + 是否被强退
type SessionInfo struct {
	UserInfoId int

	Browser   string
	OS        string
	IsOffline int
}

// 用户详细信息: 仅用于在后端内部进行传输
type UserDetailDTO struct {
	LoginVO
	Password   string   `json:"password"`
	IsDisable  bool     `json:"is_disable"`
	Browser    string   `json:"browser"`
	OS         string   `json:"os"`
	RoleLabels []string `json:"role_labels"`
}

type LoginVO struct {
	UserAuthId int `json:"id"`
	UserInfoId int `json:"user_info_id"`

	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`

	IpAddress     string    `json:"ip_address"`
	IpSource      string    `json:"ip_source"`
	LastLoginTime time.Time `json:"last_login_time"`
	LoginType     int       `json:"login_type"`

	// 点赞 Set: 用于记录用户点赞过的文章, 评论
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
}
