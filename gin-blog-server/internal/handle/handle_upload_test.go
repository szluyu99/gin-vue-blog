package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 一个内容合法的最小 PNG: http.DetectContentType 只认前 8 个字节的签名
const pngContent = "\x89PNG\r\n\x1a\n" + "fake-image-body"

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

// 挂上上传接口, 并模拟已登录
func newUploadEnv(t *testing.T) *testEnv {
	t.Helper()
	env := newTestEnv(t)
	env.loginAs(7, "tester")
	env.engine.POST("/upload", (&Upload{}).UploadFile)
	return env
}

// 上传到本地: 返回展示路径, 文件真的落盘
func TestUploadFileLocal(t *testing.T) {
	dir := withUploadConf(t, "local")
	env := newUploadEnv(t)

	resp := env.upload(t, "/upload", "封面.png", pngContent)
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
	assert.Equal(t, pngContent, string(content))
}

// 未登录不能上传: /front/upload 走 JWTAuth(false), 中间件不会拦, handler 必须自己拦
func TestUploadFileWithoutLogin(t *testing.T) {
	dir := withUploadConf(t, "local")
	env := newTestEnv(t) // 不调用 loginAs
	env.engine.POST("/upload", (&Upload{}).UploadFile)

	resp := env.upload(t, "/upload", "a.png", pngContent)
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)

	// 没有产生任何文件
	entries, err := os.ReadDir(dir)
	assert.Nil(t, err)
	assert.Empty(t, entries)
}

// 非图片后缀被拒绝
func TestUploadFileRejectNonImageExt(t *testing.T) {
	dir := withUploadConf(t, "local")
	env := newUploadEnv(t)

	for _, name := range []string{"evil.html", "shell.php", "note.txt", "icon.svg"} {
		resp := env.upload(t, "/upload", name, pngContent)
		assert.Equal(t, g.ErrFileType.Code(), resp.Code, name)
	}

	entries, err := os.ReadDir(dir)
	assert.Nil(t, err)
	assert.Empty(t, entries)
}

// 后缀伪装成图片, 但内容不是图片
func TestUploadFileRejectFakeImage(t *testing.T) {
	withUploadConf(t, "local")
	env := newUploadEnv(t)

	resp := env.upload(t, "/upload", "evil.png", "<html><script>alert(1)</script></html>")
	assert.Equal(t, g.ErrFileType.Code(), resp.Code)
}

// 超过 10MB 的文件被拒绝
func TestUploadFileTooLarge(t *testing.T) {
	withUploadConf(t, "local")
	env := newUploadEnv(t)

	big := pngContent + strings.Repeat("x", maxUploadSize)
	resp := env.upload(t, "/upload", "big.png", big)
	assert.Equal(t, g.ErrFileSize.Code(), resp.Code)

	// 刚好在上限内的可以上传
	resp = env.upload(t, "/upload", "ok.png", pngContent+strings.Repeat("x", 1024))
	assert.Equal(t, g.SUCCESS, resp.Code)
}

// 没有带文件的请求
func TestUploadFileWithoutFile(t *testing.T) {
	withUploadConf(t, "local")
	env := newUploadEnv(t)

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

	env := newUploadEnv(t)

	resp := env.upload(t, "/upload", "a.png", pngContent)
	assert.Equal(t, g.ErrFileUpload.Code(), resp.Code)
}

// 上传接口不依赖数据库, 但登录用户信息要能取到
func TestUploadFileUsesLoginUser(t *testing.T) {
	withUploadConf(t, "local")
	env := newTestEnv(t)
	env.user = &model.UserAuth{Model: model.Model{ID: 9}, Username: "tester"}
	env.engine.POST("/upload", (&Upload{}).UploadFile)

	resp := env.upload(t, "/upload", "a.png", pngContent)
	assert.Equal(t, g.SUCCESS, resp.Code)
}
