package model

import (
	"gorm.io/gorm"
)

type FriendLink struct {
	Model
	Name    string `gorm:"type:varchar(50)" json:"name"`
	Avatar  string `gorm:"type:varchar(255)" json:"avatar"`
	Address string `gorm:"type:varchar(255)" json:"address"`
	Intro   string `gorm:"type:varchar(255)" json:"intro"`
}

func GetLinkList(db *gorm.DB, num, size int, keyword string) (list []FriendLink, total int64, err error) {
	db = db.Model(&FriendLink{})
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
		db = db.Or("address LIKE ?", "%"+keyword+"%")
		db = db.Or("intro LIKE ?", "%"+keyword+"%")
	}
	db.Count(&total)
	result := db.Order("created_at DESC").
		Scopes(Paginate(num, size)).
		Find(&list)
	return list, total, result.Error
}

func SaveOrUpdateLink(db *gorm.DB, id int, name, avatar, address, intro string) (*FriendLink, error) {
	link := FriendLink{
		Model:   Model{ID: id},
		Name:    name,
		Avatar:  avatar,
		Address: address,
		Intro:   intro,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&link)
	} else {
		result = db.Create(&link)
	}

	return &link, result.Error
}
