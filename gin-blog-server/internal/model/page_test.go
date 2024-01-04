package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPageList(t *testing.T) {
	db, _ := initModelDB()

	db.Create(&Page{Name: "name1", Label: "label1", Cover: "cover1"})
	db.Create(&Page{Name: "name2", Label: "label2", Cover: "cover2"})
	db.Create(&Page{Name: "name3", Label: "label3", Cover: "cover3"})
	db.Create(&Page{Name: "name4", Label: "label4", Cover: "cover4"})
	db.Create(&Page{Name: "name5", Label: "label5", Cover: "cover5"})

	list, total, err := GetPageList(db)
	assert.Nil(t, err)
	assert.Equal(t, int64(5), total)
	assert.Equal(t, 5, len(list))
}

func TestSaveOrUpdatePage(t *testing.T) {
	db, _ := initModelDB()

	page, err := SaveOrUpdatePage(db, 0, "name", "label", "cover")
	assert.Nil(t, err)
	assert.Equal(t, "name", page.Name)
	assert.Equal(t, "label", page.Label)
	assert.Equal(t, "cover", page.Cover)

	page, err = SaveOrUpdatePage(db, page.ID, "name2", "label2", "cover2")
	assert.Nil(t, err)
	assert.Equal(t, "name2", page.Name)
	assert.Equal(t, "label2", page.Label)
	assert.Equal(t, "cover2", page.Cover)

	var val Page
	db.First(&val, page.ID)
	assert.Equal(t, "name2", val.Name)
	assert.Equal(t, "label2", val.Label)
	assert.Equal(t, "cover2", val.Cover)
}
