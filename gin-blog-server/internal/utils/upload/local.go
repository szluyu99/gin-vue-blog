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
	"path/filepath"
	"strings"
	"time"
)

// 本地文件上传
type Local struct{}

// 文件上传到本地
func (*Local) UploadFile(file *multipart.FileHeader) (filePath, fileName string, err error) {
	ext := path.Ext(file.Filename)                                     // 读取文件后缀
	name := strings.TrimSuffix(file.Filename, ext)                     // 读取文件名
	name = utils.MD5(name)                                             // 加密文件名
	filename := name + "_" + time.Now().Format("20060102150405") + ext // 拼接新文件名

	conf := g.Conf.Upload
	mkdirErr := os.MkdirAll(conf.StorePath, os.ModePerm) // 尝试创建存储路径
	if mkdirErr != nil {
		slog.Error("function os.MkdirAll() Filed", slog.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}

	storePath := conf.StorePath + "/" + filename // 文件存储路径
	filepath := conf.Path + "/" + filename       // 文件展示路径

	f, openError := file.Open() // 读取文件
	if openError != nil {
		slog.Error("function file.Open() Filed", slog.String("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(storePath)
	if createErr != nil {
		slog.Error("function os.Create() Filed", slog.String("err", createErr.Error()))
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 拷贝文件
	if copyErr != nil {
		slog.Error("function io.Copy() Filed", slog.String("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

// 从本地删除文件
//
// key 来自调用方, 可能包含 ../, 必须先规整再判断是否仍在 StorePath 之内。
// 原来的写法是 strings.Contains(p, StorePath): p 由 StorePath 拼接而来, 这个条件恒真,
// 等于没有防护, key 传 "../../config.yml" 也能删。
func (*Local) DeleteFile(key string) error {
	storePath := g.GetConfig().Upload.StorePath
	root, err := filepath.Abs(storePath)
	if err != nil {
		return errors.New("上传目录解析失败, err:" + err.Error())
	}

	p, err := filepath.Abs(filepath.Join(root, key))
	if err != nil {
		return errors.New("文件路径解析失败, err:" + err.Error())
	}

	// 必须是 root 的子路径, 前缀比较时带上分隔符, 避免 /data/uploaded-evil 混过去
	if p != root && !strings.HasPrefix(p, root+string(os.PathSeparator)) {
		return errors.New("非法的文件路径: " + key)
	}

	if err := os.Remove(p); err != nil {
		return errors.New("本地文件删除失败, err:" + err.Error())
	}
	return nil
}
