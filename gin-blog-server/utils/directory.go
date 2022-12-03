package utils

import (
	"errors"
	"os"
)

// 判断文件目录是否存在
func PathExists(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	// 存在文件夹
	if fileInfo.IsDir() {
		return true, nil
	}
	// 不是文件夹, 则存在同名文件
	return false, errors.New("存在同名文件")
}

// 获取当前函数所在的 .go 文件目录
// func GetCurrentPath() string {
// 	_, filename, _, _ := runtime.Caller(1)
// 	return path.Dir(filename)
// }
