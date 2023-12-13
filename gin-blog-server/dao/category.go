package dao

import (
	"gin-blog/model"
	"gin-blog/model/resp"
)

type Category struct{}

func (*Category) GetList(pageNum, pageSize int, keyword string) ([]resp.CategoryVo, int64) {
	var datas = make([]resp.CategoryVo, 0)
	var total int64

	db := DB.
		Table("category c").
		Select("c.id", "c.name", "COUNT(a.id) AS article_count", "c.created_at", "c.updated_at").
		Joins("LEFT JOIN article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1")

	// 条件查询
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	db.Group("c.id").
		Order("c.id DESC").
		Count(&total).
		Limit(pageSize).Offset(pageSize * (pageNum - 1)).
		Find(&datas)
	return datas, total
}

func (*Category) GetOption() []resp.OptionVo {
	var list = make([]resp.OptionVo, 0)
	DB.Model(&model.Category{}).Select("id", "name").Find(&list)
	return list
}

// TODO 修改前台查询
// func (*Category) GetFrontList() {

// }
