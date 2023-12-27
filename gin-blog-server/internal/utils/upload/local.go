package upload

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

// 本地文件上传
type Local struct{}

// 文件上传到本地
func (*Local) UploadFile(file *multipart.FileHeader) (filePath, fileName string, err error) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名
	name := strings.TrimSuffix(file.Filename, ext)
	// 加密文件名
	name = utils.MD5(name)
	// 拼接新文件名
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(g.GetConfig().Upload.StorePath, os.ModePerm)
	if mkdirErr != nil {
		slog.Error("function os.MkdirAll() Filed", slog.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}

	p := g.GetConfig().Upload.StorePath + "/" + filename   // 文件存储路径
	filepath := g.GetConfig().Upload.Path + "/" + filename // 文件展示路径

	f, openError := file.Open() // 读取文件
	if openError != nil {
		slog.Error("function file.Open() Filed", slog.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		slog.Error("function os.Create() Filed", slog.Any("err", createErr.Error()))
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 拷贝文件
	if copyErr != nil {
		slog.Error("function io.Copy() Filed", slog.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

// 从本地删除文件
func (*Local) DeleteFile(key string) error {
	p := g.GetConfig().Upload.StorePath + "/" + key
	if strings.Contains(p, g.GetConfig().Upload.StorePath) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
