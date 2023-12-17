package dao

import (
	"gin-blog/model"
)

type Message struct{}

func (*Message) GetList(pageNum, pageSize int, nickname string, isReview *int8) ([]model.Message, int64) {
	var list = make([]model.Message, 0)
	var total int64

	db := DB.Model(&Message{})
	if nickname != "" {
		db = db.Where("nickname like ?", "%"+nickname+"%")
	}
	if isReview != nil {
		db = db.Where("is_review = ?", isReview)
	}
	db.Count(&total)
	db.Order("id DESC").
		Limit(pageSize).Offset(pageSize * (pageNum - 1)).
		Find(&list)
	return list, total
}
