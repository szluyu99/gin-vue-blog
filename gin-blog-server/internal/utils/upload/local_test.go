package upload

import (
	"bytes"
	g "gin-blog/internal/global"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 造一个真实的 multipart.FileHeader
func newFileHeader(t *testing.T, filename, content string) *multipart.FileHeader {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	assert.Nil(t, err)
	_, err = part.Write([]byte(content))
	assert.Nil(t, err)
	assert.Nil(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.Nil(t, req.ParseMultipartForm(1<<20))

	return req.MultipartForm.File["file"][0]
}

// 指向临时目录的上传配置, 用例结束自动清理
func withUploadConf(t *testing.T, ossType string) string {
	t.Helper()
	dir := t.TempDir()
	old := g.Conf
	g.Conf = &g.Config{}
	g.Conf.Upload.OssType = ossType
	g.Conf.Upload.Path = "/public/uploaded"
	g.Conf.Upload.StorePath = dir
	t.Cleanup(func() { g.Conf = old })
	return dir
}

func TestLocalUploadAndDeleteFile(t *testing.T) {
	dir := withUploadConf(t, "local")
	local := &Local{}

	filePath, fileName, err := local.UploadFile(newFileHeader(t, "头像.png", "content"))
	assert.Nil(t, err)

	// 文件名被 MD5 加密并拼上时间戳, 后缀保留
	assert.True(t, strings.HasSuffix(fileName, ".png"))
	assert.NotContains(t, fileName, "头像")
	// 返回的是展示路径, 不是磁盘路径
	assert.Equal(t, "/public/uploaded/"+fileName, filePath)

	content, err := os.ReadFile(filepath.Join(dir, fileName))
	assert.Nil(t, err)
	assert.Equal(t, "content", string(content))

	assert.Nil(t, local.DeleteFile(fileName))
	_, err = os.Stat(filepath.Join(dir, fileName))
	assert.True(t, os.IsNotExist(err))
}

// 存储目录不存在时自动创建
func TestLocalUploadCreatesStorePath(t *testing.T) {
	dir := withUploadConf(t, "local")
	g.Conf.Upload.StorePath = filepath.Join(dir, "nested", "dir")

	_, fileName, err := (&Local{}).UploadFile(newFileHeader(t, "a.png", "content"))
	assert.Nil(t, err)
	assert.FileExists(t, filepath.Join(g.Conf.Upload.StorePath, fileName))
}

// 删除不存在的文件应该报错, 而不是静默成功
func TestLocalDeleteMissingFile(t *testing.T) {
	withUploadConf(t, "local")
	err := (&Local{}).DeleteFile("not-exist.png")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "本地文件删除失败")
}

// 工厂方法按配置返回实现
func TestNewOSS(t *testing.T) {
	withUploadConf(t, "local")
	assert.IsType(t, &Local{}, NewOSS())

	withUploadConf(t, "qiniu")
	assert.IsType(t, &Qiniu{}, NewOSS())

	// 未知类型兜底用本地存储
	withUploadConf(t, "unknown")
	assert.IsType(t, &Local{}, NewOSS())
}
