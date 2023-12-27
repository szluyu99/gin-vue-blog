package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	Universal
	Name string `gorm:"unique;type:varchar(20);not null" json:"name"`

	Articles []*Article `gorm:"many2many:article_tag;" json:"articles,omitempty"`
}

type TagVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}

func GetTagList(db *gorm.DB, page, size int, keyword string) (list []TagVO, total int64, err error) {
	db = db.Table("tag t").
		Joins("LEFT JOIN article_tag at ON t.id = at.tag_id").
		Select("t.id", "t.name", "COUNT(at.article_id) AS article_count", "t.created_at", "t.updated_at")

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	result := db.
		Group("t.id").Order("t.updated_at DESC").
		Count(&total).
		Limit(size).Offset(size * (page - 1)).
		Find(&list)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return list, total, nil
}

func GetTagOption(db *gorm.DB) ([]OptionVO, error) {
	list := make([]OptionVO, 0)
	result := db.Model(&Tag{}).Select("id", "name").Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

// 根据 [文章id] 获取 [标签名称列表]
func GetTagNamesByArticleId(db *gorm.DB, id int) ([]string, error) {
	list := make([]string, 0)
	result := db.Table("tag").
		Joins("LEFT JOIN article_tag ON tag.id = article_tag.tag_id").
		Where("article_id", id).
		Pluck("name", &list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func SaveOrUpdateTag(db *gorm.DB, id int, name string) (*Tag, error) {
	tag := Tag{Name: name}

	var result *gorm.DB
	if id > 0 {
		result = db.Model(&tag).Where("id", id).Updates(tag)
	} else {
		result = db.Create(&tag)
	}

	if result.Error != nil {
		return nil, result.Error
	}
	return &tag, nil
}
