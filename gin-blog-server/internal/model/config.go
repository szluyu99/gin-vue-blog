package model

import (
	"strconv"

	"gorm.io/gorm"
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

func CheckConfigMap(db *gorm.DB, m map[string]string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for k, v := range m {
			result := tx.Model(Config{}).Where("key", k).Update("value", v)
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
