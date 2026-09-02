package model

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type user struct {
	UUID      uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Age       int
	Enabled   bool
}

type product struct {
	UUID      string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	CanBuy    bool
}

func initDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	db.AutoMigrate(user{}, product{})
	return db
}

func TestCreate(t *testing.T) {
	db := initDB()

	val, err := Create(db, &user{
		Name:    "mockname",
		Age:     11,
		Enabled: true,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, val.UUID)

	p, err := Create(db, &product{
		UUID:   "aaaa",
		Name:   "demoproduct",
		CanBuy: true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, p)
}

func TestCount(t *testing.T) {
	db := initDB()

	db.Create(&user{Name: "user1", Email: "user1@example.com", Age: 10})
	count, err := Count(db, &user{})
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	db.Create(&user{Name: "user2", Email: "user2@example.com", Age: 20})
	count, err = Count(db, &user{})
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	db.Create(&user{Name: "user3", Email: "user3@example.com", Age: 30})
	count, err = Count(db, &user{})
	assert.Nil(t, err)
	assert.Equal(t, 3, count)

	count, err = Count(db, &user{}, "age >= ?", 20)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}

// Paginate 会把 size 截到 100, 前台"取全部"的列表必须走 PageSizeAll 才不会丢数据
func TestPaginate(t *testing.T) {
	db := initDB()

	for i := 0; i < 120; i++ {
		assert.Nil(t, db.Create(&user{Name: "u", Age: i}).Error)
	}

	count := func(page, size int) int {
		var list []user
		assert.Nil(t, db.Model(&user{}).Scopes(Paginate(page, size)).Find(&list).Error)
		return len(list)
	}

	assert.Equal(t, 10, count(1, 0), "size <= 0 时用默认的 10")
	assert.Equal(t, 20, count(1, 20))
	assert.Equal(t, 100, count(1, 1000), "普通分页最多 100 条")
	assert.Equal(t, 120, count(1, PageSizeAll), "PageSizeAll 不分页")
	// page 非法时按第一页处理
	assert.Equal(t, 20, count(0, 20))
	assert.Equal(t, 20, count(-5, 20))
}
