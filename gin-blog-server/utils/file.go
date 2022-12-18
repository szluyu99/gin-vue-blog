package utils

import (
	"bufio"
	"io"
	"mime/multipart"
	"os"

	"go.uber.org/zap"
)

var File fileUtils

type fileUtils struct{}

// 将文件写入到本地: 已有同名文件则覆盖
func (*fileUtils) WriteFile(name, path, content string) {
	// 指定模式打开文件: 读写|创建
	file, err := os.OpenFile(path+name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		Logger.Error("文件写入, 目标地址错误: ", zap.Error(err))
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)      // 文件写入对象
	_, err = writer.WriteString(content) // 将内容写到缓存
	writer.Flush()                       // 将缓存写到文件
	if err != nil {
		Logger.Error("文件写入失败: ", zap.Error(err))
	}
}

// 从 fileHeader 中读取内容
func (*fileUtils) ReadFromFileHeader(file *multipart.FileHeader) string {
	open, err := file.Open()
	if err != nil {
		Logger.Error("文件读取, 目标地址错误: ", zap.Error(err))
		return ""
	}
	defer open.Close()
	all, err := io.ReadAll(open)
	if err != nil {
		Logger.Error("文件读取失败: ", zap.Error(err))
		return ""
	}
	return string(all)
}
