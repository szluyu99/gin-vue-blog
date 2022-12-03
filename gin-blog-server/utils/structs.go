package utils

import (
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// 结构体 转 map[string]any, 需要配合 `mapstructure`
func Struct2Map(s any) map[string]any {
	maps := map[string]any{}
	if err := mapstructure.Decode(s, &maps); err != nil {
		Logger.Error("Struct2Map: ", zap.Error(err))
		panic(err)
	}
	return maps
}

// 拷贝属性, 一般用于 vo -> po
func CopyProperties[T any](from any) (to T) {
	if err := copier.Copy(&to, from); err != nil {
		Logger.Error("CopyProperties: ", zap.Error(err))
		panic(err)
	}
	return to
}
