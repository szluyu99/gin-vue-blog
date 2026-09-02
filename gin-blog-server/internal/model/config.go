package model

import (
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config struct {
	Model
	Key   string `gorm:"unique;type:varchar(256)" json:"key"`
	Value string `gorm:"type:varchar(256)" json:"value"`
	Desc  string `gorm:"type:varchar(256)" json:"desc"`
}

func GetConfigMap(db *gorm.DB) (map[string]string, error) {
	var configs []Config
	result := db.Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}

	m := make(map[string]string)
	for _, config := range configs {
		m[config.Key] = config.Value
	}

	return m, nil
}

/*
批量保存配置

用 upsert 而不是 Update: 配置表里没有这一行时 Update 影响 0 行且不报错,
后台设置页新增的配置项会静默保存失败(is_message_review 就吃过这个亏)。
*/
func CheckConfigMap(db *gorm.DB, m map[string]string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for k, v := range m {
			config := Config{Key: k, Value: v}
			result := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "key"}},
				DoUpdates: clause.AssignmentColumns([]string{"value"}),
			}).Create(&config)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func CheckConfig(db *gorm.DB, key, value string) error {
	var config Config

	result := db.Where("key", key).FirstOrCreate(&config)
	if result.Error != nil {
		return result.Error
	}

	config.Value = value
	result = db.Save(&config)

	return result.Error
}

func GetConfig(db *gorm.DB, key string) string {
	var config Config
	result := db.Where("key", key).First(&config)

	if result.Error != nil {
		return ""
	}

	return config.Value
}

func GetConfigBool(db *gorm.DB, key string) bool {
	val := GetConfig(db, key)
	if val == "" {
		return false
	}
	return val == "true"
}

func GetConfigInt(db *gorm.DB, key string) int {
	val := GetConfig(db, key)
	if val == "" {
		return 0
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return result
}
