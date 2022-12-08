package model

type Page struct {
	Universal
	Name  string `gorm:"type:varchar(20);not null;comment:页面名称" json:"name"`
	Label string `gorm:"type:varchar(30);comment:页面标签" json:"label"`
	Cover string `gorm:"type:varchar(255);not null;comment:页面封面" json:"cover"`
}
