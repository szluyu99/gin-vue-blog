package model

import (
	"time"
)

// 包含逻辑删除的模型 (感觉不是很好用)
// type DelUniversal struct {
// 	ID        int            `gorm:"primary_key;auto_increment" json:"id"`
// 	CreatedAt time.Time      `json:"created_at"`
// 	UpdatedAt time.Time      `json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 该字段不展示给前端
// }

// 不包含逻辑删除的模型
type Universal struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" mapstructure:"-"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"-"`
}
