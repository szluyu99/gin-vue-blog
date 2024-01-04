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
