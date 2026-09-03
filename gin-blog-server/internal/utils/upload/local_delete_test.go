package upload

import (
	"os"
	"path/filepath"
	"testing"

	g "gin-blog/internal/global"

	"github.com/stretchr/testify/assert"
)

// DeleteFile 只能删 StorePath 之内的文件
func TestLocalDeleteFileRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(root, "outside.txt")
	assert.Nil(t, os.WriteFile(outside, []byte("secret"), 0o600))

	store := filepath.Join(root, "uploaded")
	assert.Nil(t, os.MkdirAll(store, 0o755))
	inside := filepath.Join(store, "a.jpg")
	assert.Nil(t, os.WriteFile(inside, []byte("img"), 0o600))

	old := g.Conf
	g.Conf = &g.Config{}
	g.Conf.Upload.StorePath = store
	t.Cleanup(func() { g.Conf = old })

	l := &Local{}

	// 穿越到上级目录被拒绝, 文件还在
	err := l.DeleteFile("../outside.txt")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "非法的文件路径")
	_, statErr := os.Stat(outside)
	assert.Nil(t, statErr)

	// 前缀相同但不是子目录, 同样被拒绝
	assert.NotNil(t, l.DeleteFile("../uploaded-evil/x.jpg"))

	// 正常的 key 可以删掉
	assert.Nil(t, l.DeleteFile("a.jpg"))
	_, statErr = os.Stat(inside)
	assert.True(t, os.IsNotExist(statErr))
}
