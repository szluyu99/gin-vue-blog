package handle

import (
	g "gin-blog/internal/global"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 指向临时目录的上传配置, 用例结束自动还原
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

// 上传到本地: 返回展示路径, 文件真的落盘
func TestUploadFileLocal(t *testing.T) {
	dir := withUploadConf(t, "local")
	env := newTestEnv(t)
	env.engine.POST("/upload", (&Upload{}).UploadFile)

	resp := env.upload(t, "/upload", "封面.png", "image-content")
	assert.Equal(t, g.SUCCESS, resp.Code)

	filePath, ok := resp.Data.(string)
	assert.True(t, ok)
	assert.Contains(t, filePath, "/public/uploaded/")
	assert.Contains(t, filePath, ".png")

	// 文件名被 MD5 处理过, 不会泄露原始中文名
	assert.NotContains(t, filePath, "封面")

	entries, err := os.ReadDir(dir)
	assert.Nil(t, err)
	assert.Len(t, entries, 1)
	content, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	assert.Nil(t, err)
	assert.Equal(t, "image-content", string(content))
}

// 没有带文件的请求
func TestUploadFileWithoutFile(t *testing.T) {
	withUploadConf(t, "local")
	env := newTestEnv(t)
	env.engine.POST("/upload", (&Upload{}).UploadFile)

	resp := env.do(t, http.MethodPost, "/upload", map[string]any{"file": "not-a-file"})
	assert.Equal(t, g.ErrFileReceive.Code(), resp.Code)
}

// 存储目录不可写时返回上传失败, 而不是 500
func TestUploadFileStoreFailed(t *testing.T) {
	dir := withUploadConf(t, "local")
	// 用一个已存在的文件当目录, MkdirAll 会失败
	blocker := filepath.Join(dir, "blocker")
	assert.Nil(t, os.WriteFile(blocker, []byte("x"), 0o600))
	g.Conf.Upload.StorePath = filepath.Join(blocker, "sub")

	env := newTestEnv(t)
	env.engine.POST("/upload", (&Upload{}).UploadFile)

	resp := env.upload(t, "/upload", "a.png", "content")
	assert.Equal(t, g.ErrFileUpload.Code(), resp.Code)
}
