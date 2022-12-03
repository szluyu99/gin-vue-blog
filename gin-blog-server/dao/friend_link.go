package dao

import (
	"gin-blog/model"
	"gin-blog/model/req"
)

type FriendLink struct{}

func (*FriendLink) GetList(req req.PageQuery) (list []model.FriendLink, total int64) {
	tx := DB.Model(&model.FriendLink{})
	if req.Keyword != "" {
		tx = tx.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	tx.Count(&total)
	err := tx.Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
