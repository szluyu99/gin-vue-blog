package middleware

import (
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func countLogs(e *mwEnv) int64 {
	var count int64
	e.db.Model(&model.OperationLog{}).Count(&count)
	return count
}

// GET 请求不记录 (量太大)
func TestOperationLogSkipGet(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodGet, "/article/list", OperationLog())

	w := e.get("/api/article/list")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan)
	assert.Equal(t, int64(0), countLogs(e))
}

// 文件上传不记录 (请求体太长)
func TestOperationLogSkipUpload(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodPost, "/upload", OperationLog())

	w := e.request(http.MethodPost, "/api/upload", "binary...", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, e.handlerRan)
	assert.Equal(t, int64(0), countLogs(e))
}

// 写操作: 请求体和响应体都要落库, 并记录操作人
func TestOperationLogRecordWrite(t *testing.T) {
	e := newMwEnv(t)
	user := e.seedUser("test", false)
	e.loginId = user.ID
	e.handle(http.MethodPost, "/article", OperationLog())

	w := e.request(http.MethodPost, "/api/article", `{"title":"hello"}`, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), countLogs(e))

	var log model.OperationLog
	assert.Nil(t, e.db.First(&log).Error)
	assert.Equal(t, "/api/article", log.OptUrl)
	assert.Equal(t, http.MethodPost, log.RequestMethod)
	assert.Equal(t, `{"title":"hello"}`, log.RequestParam)
	assert.Contains(t, log.ResponseData, `"code":0`)
	assert.Equal(t, user.UserInfoId, log.UserId)
	assert.Equal(t, "test-nickname", log.Nickname)
	// handler 依然拿得到完整请求体 (中间件读完要还原)
	assert.Contains(t, w.Body.String(), `"code":0`)
}

// 未登录时也要能记录, 不能因为拿不到用户就 nil 解引用
func TestOperationLogWithoutLogin(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodDelete, "/article", OperationLog())

	w := e.request(http.MethodDelete, "/api/article", `[1]`, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), countLogs(e))

	var log model.OperationLog
	assert.Nil(t, e.db.First(&log).Error)
	assert.Equal(t, 0, log.UserId)
	assert.Empty(t, log.Nickname)
}

// handler 里读到的请求体没有被中间件消费掉
func TestOperationLogRequestBodyReusable(t *testing.T) {
	e := newMwEnv(t)
	var got string
	api := e.engine.Group("/api")
	api.Use(OperationLog())
	api.POST("/article", func(c *gin.Context) {
		var body map[string]string
		_ = c.ShouldBindJSON(&body)
		got = body["title"]
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	e.request(http.MethodPost, "/api/article", `{"title":"hello"}`, nil)

	assert.Equal(t, "hello", got)
}

func TestGetOptString(t *testing.T) {
	assert.Equal(t, "新增或修改", GetOptString(http.MethodPost))
	assert.Equal(t, "删除", GetOptString(http.MethodDelete))
	assert.Equal(t, "文章", GetOptString("Article"))
	assert.Empty(t, GetOptString("Unknown"))
}

func TestGetOptResource(t *testing.T) {
	assert.Equal(t, "Resource", getOptResource("gin-blog/api/v1.(*Resource).Delete-fm"))
	assert.Equal(t, "Article", getOptResource("gin-blog/internal/handle.(*Article).Save-fm"))
}
