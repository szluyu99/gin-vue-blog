package model

import (
	"gorm.io/gorm"
)

// hasMany: 一个分类下可以有多篇文章
type Category struct {
	Model
	Name     string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryId"`
}

type CategoryVO struct {
	Category
	ArticleCount int `json:"article_count"`
}

func GetCategoryList(db *gorm.DB, num, size int, keyword string) ([]CategoryVO, int64, error) {
	var list = make([]CategoryVO, 0)
	var total int64

	db = db.Table("category c").
		Select("c.id", "c.name", "COUNT(a.id) AS article_count", "c.created_at", "c.updated_at").
		Joins("LEFT JOIN article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1")

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	result := db.Group("c.id").
		Order("c.updated_at DESC").
		Count(&total).
		Scopes(Paginate(num, size)).
		Find(&list)

	return list, total, result.Error
}

func GetCategoryOption(db *gorm.DB) ([]OptionVO, error) {
	var list = make([]OptionVO, 0)
	result := db.Model(&Category{}).Select("id", "name").Find(&list)
	return list, result.Error
}

func GetCategoryById(db *gorm.DB, id int) (*Category, error) {
	var category Category
	result := db.Where("id", id).First(&category)
	return &category, result.Error
}

func GetCategoryByName(db *gorm.DB, name string) (*Category, error) {
	var category Category
	result := db.Where("name", name).First(&category)
	return &category, result.Error
}

func DeleteCategory(db *gorm.DB, ids []int) (int64, error) {
	result := db.Where("id IN ?", ids).Delete(Category{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func SaveOrUpdateCategory(db *gorm.DB, id int, name string) (*Category, error) {
	category := Category{
		Model: Model{ID: id},
		Name:  name,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(category)
	} else {
		result = db.Create(&category)
	}

	return &category, result.Error
}
