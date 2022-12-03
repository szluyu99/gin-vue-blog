package model

import (
	"reflect"
)

type Tag struct {
	Universal
	Articles []*Article `gorm:"many2many:article_tag;" json:"articles"`
	Name     string     `gorm:"type:varchar(20);not null" json:"name"`
}

func (t *Tag) IsEmpty() bool {
	return reflect.DeepEqual(t, &Tag{})
}
