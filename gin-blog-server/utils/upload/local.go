package upload

import (
	"errors"
	"gin-blog/config"
	"gin-blog/utils"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
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
	name = utils.Encryptor.MD5(name)
	// 拼接新文件名
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(config.Cfg.Upload.StorePath, os.ModePerm)
	if mkdirErr != nil {
		utils.Logger.Error("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}

	p := config.Cfg.Upload.StorePath + "/" + filename   // 文件存储路径
	filepath := config.Cfg.Upload.Path + "/" + filename // 文件展示路径

	f, openError := file.Open() // 读取文件
	if openError != nil {
		utils.Logger.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		utils.Logger.Error("function os.Create() Filed", zap.Any("err", createErr.Error()))
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 拷贝文件
	if copyErr != nil {
		utils.Logger.Error("function io.Copy() Filed", zap.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

// 从本地删除文件
func (*Local) DeleteFile(key string) error {
	p := config.Cfg.Upload.StorePath + "/" + key
	if strings.Contains(p, config.Cfg.Upload.StorePath) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
