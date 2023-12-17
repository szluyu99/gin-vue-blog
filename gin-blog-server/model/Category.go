package model

import (
	"reflect"
	"time"

	"gorm.io/gorm"
)

// hasMany: 一个分类下可以有多篇文章
type Category struct {
	Universal
	Name     string    `gorm:"type:varchar(20);not null;comment:分类名称" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryId"` // 重写外键
}

func (c *Category) IsEmpty() bool {
	return reflect.DeepEqual(c, &Category{})
}

type CategoryVO struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ArticleCount int       `json:"article_count"` // 文章数量
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetCategoryList(db *gorm.DB, pageNum, pageSize int, keyword string) ([]CategoryVO, int64, error) {
	var list = make([]CategoryVO, 0)
	var total int64

	db = db.Table("category c").
		Select("c.id", "c.name", "COUNT(a.id) AS article_count", "c.created_at", "c.updated_at").
		Joins("LEFT JOIN article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1")

	// 条件查询
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	err := db.Group("c.id").
		Order("c.id DESC").
		Count(&total).
		Limit(pageSize).Offset(pageSize * (pageNum - 1)).
		Find(&list).Error

	return list, total, err
}

func GetCategoryOption(db *gorm.DB) ([]OptionVO, error) {
	var list = make([]OptionVO, 0)
	err := db.Model(&Category{}).Select("id", "name").Find(&list).Error
	return list, err
}
