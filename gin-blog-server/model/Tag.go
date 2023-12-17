package model

import (
	"reflect"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	Universal
	Articles []*Article `gorm:"many2many:article_tag;" json:"articles"`
	Name     string     `gorm:"type:varchar(20);not null" json:"name"`
}

type TagVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Articles     []*model.Article `json:"articles"`
	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}

func (t *Tag) IsEmpty() bool {
	return reflect.DeepEqual(t, &Tag{})
}

func GetTagList(db *gorm.DB, pageNum, pageSize int, keyword string) ([]TagVO, int64, error) {
	var list = make([]TagVO, 0)
	var total int64

	db = db.Table("tag t").
		Select("t.id", "name", "COUNT(at.article_id) AS article_count", "t.created_at", "t.updated_at").
		Joins("LEFT JOIN article_tag at ON t.id = at.tag_id")

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	err := db.Group("t.id").Order("t.id DESC").
		Count(&total).
		Limit(pageSize).Offset(pageSize * (pageNum - 1)).
		Find(&list).Error

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func GetTagOption(db *gorm.DB) ([]OptionVO, error) {
	var list = make([]OptionVO, 0)
	err := db.Model(&Tag{}).Select("id", "name").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, err
}

// 根据 [文章id] 获取 [标签名称列表]
func GetTagNameByArticleId(db *gorm.DB, id int) ([]string, error) {
	list := make([]string, 0)
	err := db.Table("tag").
		Joins("LEFT JOIN article_tag ON tag.id = article_tag.tag_id").
		Where("article_id", id).
		Pluck("name", &list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
