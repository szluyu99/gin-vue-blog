package model

import (
	"gorm.io/gorm"
)

type Message struct {
	Model
	Nickname  string `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Avatar    string `gorm:"type:varchar(255);comment:头像地址" json:"avatar"`
	Content   string `gorm:"type:varchar(255);comment:留言内容" json:"content"`
	IpAddress string `gorm:"type:varchar(50);comment:IP 地址" json:"ip_address"`
	IpSource  string `gorm:"type:varchar(255);comment:IP 来源" json:"ip_source"`
	Speed     int    `gorm:"type:tinyint(1);comment:弹幕速度" json:"speed"`
	IsReview  bool   `json:"is_review"`
}

func GetMessageList(db *gorm.DB, num, size int, nickname string, isReview *bool) (list []Message, total int64, err error) {
	db = db.Model(&Message{})

	if nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+nickname+"%")
	}

	if isReview != nil {
		db = db.Where("is_review = ?", isReview)
	}

	db.Count(&total)
	result := db.Order("created_at DESC").Scopes(Paginate(num, size)).Find(&list)
	return list, total, result.Error
}

func DeleteMessages(db *gorm.DB, ids []int) (int64, error) {
	result := db.Where("id in ?", ids).Delete(&Message{})
	return result.RowsAffected, result.Error
}

func UpdateMessagesReview(db *gorm.DB, ids []int, isReview bool) (int64, error) {
	result := db.Model(&Message{}).Where("id in ?", ids).Update("is_review", isReview)
	return result.RowsAffected, result.Error
}

func SaveMessage(db *gorm.DB, nickname, avatar, content, address, source string, speed int, isReview bool) (*Message, error) {
	message := Message{
		Nickname:  nickname,
		Avatar:    avatar,
		Content:   content,
		IpAddress: address,
		IpSource:  source,
		Speed:     speed,
		IsReview:  isReview,
	}

	result := db.Create(&message)
	return &message, result.Error
}
