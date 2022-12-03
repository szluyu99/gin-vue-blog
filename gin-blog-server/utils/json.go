package utils

import (
	"encoding/json"

	"go.uber.org/zap"
)

const JSON_UTIL_ERR_PREFIX = "utils/json.go -> "

var Json = new(_json)

type _json struct{}

// data -> jsonStr
func (*_json) Marshal(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		Logger.Error(JSON_UTIL_ERR_PREFIX+"Marshal: ", zap.Error(err))
		panic(err)
	}
	return string(data)

}

// jsonStr -> data
func (*_json) Unmarshal(data string, v any) {
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		Logger.Error(JSON_UTIL_ERR_PREFIX+"Unmarshal: ", zap.Error(err))
		panic(err)
	}
}

// TODO:
// func (j *_json) DecodeMap(in, out any) {
// 	j.Unmarshal(j.Marshal(in), &out)
// }
