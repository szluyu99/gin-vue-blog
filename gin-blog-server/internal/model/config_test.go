package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigMap(t *testing.T) {
	db := initDB()
	db.AutoMigrate(&Config{})

	configs := []Config{
		{Key: "name", Value: "Blog", Desc: "姓名"},
		{Key: "age", Value: "12", Desc: "年龄"},
		{Key: "enabled", Value: "true", Desc: "是否可用"},
	}
	db.Create(&configs)

	data, err := GetConfigMap(db)
	assert.Nil(t, err)
	assert.Len(t, data, 3)
	assert.Equal(t, "Blog", data["name"])
	assert.Equal(t, "12", data["age"])
	assert.Equal(t, "true", data["enabled"])
}

func TestUpdateConfigMap(t *testing.T) {
	db := initDB()
	db.AutoMigrate(&Config{})

	configs := []Config{
		{Key: "name", Value: "Blog", Desc: "姓名"},
		{Key: "age", Value: "12", Desc: "年龄"},
		{Key: "enabled", Value: "true", Desc: "是否可用"},
	}
	db.Create(&configs)

	m := map[string]string{
		"name":    "Alice",
		"age":     "15",
		"enabled": "false",
		"new_key": "new_value", // 表里还没有的配置项
	}
	err := CheckConfigMap(db, m)
	assert.Nil(t, err)

	data, err := GetConfigMap(db)
	assert.Nil(t, err)
	assert.Equal(t, "Alice", data["name"])
	assert.Equal(t, "15", data["age"])
	assert.Equal(t, "false", data["enabled"])
	// 走 upsert: 表里没有的键要插进去, 而不是像以前那样静默丢掉
	assert.Len(t, data, 4)
	assert.Equal(t, "new_value", data["new_key"])

	// 重复保存不会插出第二行
	assert.Nil(t, CheckConfigMap(db, map[string]string{"new_key": "changed"}))
	data, err = GetConfigMap(db)
	assert.Nil(t, err)
	assert.Len(t, data, 4)
	assert.Equal(t, "changed", data["new_key"])
}

func TestConfigSetGet(t *testing.T) {
	db := initDB()
	db.AutoMigrate(&Config{})

	CheckConfig(db, "name", "AAA")

	val := GetConfig(db, "name")
	assert.Equal(t, "AAA", val)

	m, _ := GetConfigMap(db)
	assert.Len(t, m, 1)
}

func TestCheckConfig(t *testing.T) {
	db := initDB()
	db.AutoMigrate(&Config{})

	{
		CheckConfig(db, "name", "AAA")
		val := GetConfig(db, "name")
		assert.Equal(t, "AAA", val)
	}

	{
		CheckConfig(db, "name", "BBB")
		val := GetConfig(db, "name")
		assert.Equal(t, "BBB", val)
	}

}
