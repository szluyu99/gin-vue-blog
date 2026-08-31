package handle

import (
	"encoding/json"
	"errors"
	g "gin-blog/internal/global"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 构造一个只跑单个 handler 的测试上下文
func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

// 解析响应体, 业务码和消息都在 body 里
func parseBody(t *testing.T, w *httptest.ResponseRecorder) Response[any] {
	var resp Response[any]
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

func TestReturnSuccess(t *testing.T) {
	c, w := newTestContext()
	ReturnSuccess(c, map[string]any{"name": "test"})

	// 只要请求到达后端, HTTP 状态码就是 200
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseBody(t, w)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, map[string]any{"name": "test"}, resp.Data)
	assert.False(t, c.IsAborted())
}

func TestReturnErrorWithError(t *testing.T) {
	c, w := newTestContext()
	ReturnError(c, g.ErrRequest, errors.New("参数解析失败"))

	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseBody(t, w)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
	assert.Equal(t, g.ErrRequest.Msg(), resp.Message)
	// data 为 error 时, 取 Error() 作为详情
	assert.Equal(t, "参数解析失败", resp.Data)
	// 出错时必须中断后续 handler
	assert.True(t, c.IsAborted())
}

func TestReturnErrorWithString(t *testing.T) {
	c, w := newTestContext()
	ReturnError(c, g.ErrRequest, "字符串详情")

	resp := parseBody(t, w)
	assert.Equal(t, "字符串详情", resp.Data)
}

func TestReturnErrorWithoutData(t *testing.T) {
	c, w := newTestContext()
	ReturnError(c, g.ErrRequest, nil)

	resp := parseBody(t, w)
	// 没有额外数据时, 详情回落为错误码自带的消息
	assert.Equal(t, g.ErrRequest.Msg(), resp.Data)
}

func TestReturnHttpResponse(t *testing.T) {
	c, w := newTestContext()
	ReturnHttpResponse(c, http.StatusInternalServerError, g.FAIL, "系统错误", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resp := parseBody(t, w)
	assert.Equal(t, g.FAIL, resp.Code)
	assert.Equal(t, "系统错误", resp.Message)
}

func TestPageQueryBind(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page_num=2&page_size=20&keyword=vue", nil)

	var query PageQuery
	assert.Nil(t, c.ShouldBindQuery(&query))
	assert.Equal(t, 2, query.Page)
	assert.Equal(t, 20, query.Size)
	assert.Equal(t, "vue", query.Keyword)
}
