package model

import "reflect"

// 留言
type Message struct {
	Universal
	Nickname  string `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Avatar    string `gorm:"type:varchar(255);comment:头像地址" json:"avatar"`
	Content   string `gorm:"type:varchar(255);comment:留言内容" json:"content"`
	IpAddress string `gorm:"type:varchar(50);comment:IP 地址" json:"ip_address"`
	IpSource  string `gorm:"type:varchar(255);comment:IP 来源" json:"ip_source"`
	Speed     int    `gorm:"type:tinyint(1);comment:弹幕速度" json:"speed"`
	IsReview  int8   `gorm:"type:tinyint(1);not null;default:0;comment:审核状态(0-未审核,1-审核通过)" json:"is_review"`
}

func (c *Message) IsEmpty() bool {
	return reflect.DeepEqual(c, &Message{})
}
