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

// 改密码接口就挂在这个中间件下, 明文密码不能落库
func TestOperationLogMasksPassword(t *testing.T) {
	e := newMwEnv(t)
	e.handle(http.MethodPut, "/user/current/password", OperationLog())

	body := `{"old_password":"111111","new_password":"222222"}`
	w := e.request(http.MethodPut, "/api/user/current/password", body, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var log model.OperationLog
	assert.Nil(t, e.db.First(&log).Error)
	assert.NotContains(t, log.RequestParam, "111111")
	assert.NotContains(t, log.RequestParam, "222222")
	assert.Contains(t, log.RequestParam, sensitiveMask)
	// handler 拿到的仍是原文
	assert.Contains(t, w.Body.String(), `"code":0`)
}

func TestMaskSensitive(t *testing.T) {
	// 不含敏感字段: 原样返回, 不重新序列化
	plain := `{"title":"hello","desc":"world"}`
	assert.Equal(t, plain, maskSensitive(plain))
	assert.Equal(t, "", maskSensitive(""))

	// 嵌套与数组里的敏感字段同样要脱敏, 非敏感字段保留
	got := maskSensitive(`{"user":{"name":"n","password":"p"},"list":[{"token":"t"}]}`)
	assert.NotContains(t, got, `"p"`)
	assert.NotContains(t, got, `"t"`)
	assert.Contains(t, got, `"name":"n"`)

	// 大小写与变体字段名都算敏感
	assert.NotContains(t, maskSensitive(`{"Password":"p","secretKey":"s"}`), `"p"`)

	// 不是合法 JSON 又带敏感字样(表单编码): 整体丢弃
	assert.Equal(t, sensitiveMask, maskSensitive("password=111111&name=x"))

	// 不是 JSON 也不含敏感字样: 保持原样
	assert.Equal(t, "[1,2", maskSensitive("[1,2"))
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
