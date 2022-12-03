package dao

import (
	"gin-blog/model"
	"gin-blog/model/req"
)

type OperationLog struct{}

func (*OperationLog) GetList(req req.PageQuery) ([]model.OperationLog, int64) {
	list := make([]model.OperationLog, 0)
	var total int64

	db := DB.Model(&model.OperationLog{})
	if req.Keyword != "" {
		db = db.Where("opt_module LIKE ?", "%"+req.Keyword+"%").
			Or("opt_desc LIKE ?", "%"+req.Keyword+"%")
	}
	db.Count(&total).
		Order("id DESC").
		Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).
		Find(&list)
	return list, total
}
