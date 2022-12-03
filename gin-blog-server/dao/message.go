package dao

import (
	"gin-blog/model"
	"gin-blog/model/req"
)

type Message struct{}

func (*Message) GetList(req req.GetMessages) ([]model.Message, int64) {
	var list = make([]model.Message, 0)
	var total int64

	db := DB.Model(&Message{})
	if req.Nickname != "" {
		db = db.Where("nickname like ?", "%"+req.Nickname+"%")
	}
	if req.IsReview != nil {
		db = db.Where("is_review = ?", req.IsReview)
	}
	db.Count(&total)
	db.Order("id DESC").
		Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).
		Find(&list)
	return list, total
}
