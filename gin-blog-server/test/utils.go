package test

import (
	"bytes"
	"encoding/json"
	"gin-blog/routes"
	"gin-blog/utils/r"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: 请求需要 token 怎么办? 临时解决方案: 在 backRouter 中去掉 JWT 中间件
var backRouter http.Handler
var frontRouter http.Handler

// TODO: 配置文件的位置
func init() {
	testing.Init()
	routes.InitGlobalVariable()
	// backRouter = routes.BackRouter()
	// frontRouter = routes.FrontRouter()
}

func BackRouter() http.Handler {
	if backRouter != nil {
		return backRouter
	} else {
		backRouter = routes.BackRouter()
		return backRouter
	}
}

func FrontRouter() http.Handler {
	if frontRouter != nil {
		return frontRouter
	} else {
		frontRouter = routes.FrontRouter()
		return frontRouter
	}
}

// 基础测试: http 状态码正常 + 反序列化数据正常
func BaseHttpSuccessTest(t *testing.T, w *httptest.ResponseRecorder) r.Response {
	assert.Equal(t, http.StatusOK, w.Code) // HTTP 请求状态码
	var resp r.Response
	err := json.Unmarshal([]byte(w.Body.Bytes()), &resp)
	assert.Nil(t, err) // 错误
	return resp
}

// 基础测试: http 状态码正常 + 反序列化数据正常 + dataCode 正常
func BaseSuccessTest(t *testing.T, w *httptest.ResponseRecorder) r.Response {
	assert.Equal(t, http.StatusOK, w.Code) // HTTP 请求状态码
	var resp r.Response
	err := json.Unmarshal([]byte(w.Body.Bytes()), &resp)
	assert.Nil(t, err)                            // 错误
	assert.Equal(t, r.OK, resp.Code)              // 数据状态码
	assert.Equal(t, r.GetMsg(r.OK), resp.Message) // 请求信息
	return resp
}

// 将 map 中的键值对输出成 querystring 形式
func ParseToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		values += "&" + key + "=" + val
	}
	temp := values[1:]
	values = "?" + temp
	return values
}

func Get(uri string, router http.Handler) *httptest.ResponseRecorder {
	// 构造 GET 请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的 handler 接口
	router.ServeHTTP(w, req)
	return w
}

func PostForm(uri string, param map[string]string, router http.Handler) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", uri+ParseToStr(param), nil)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应 handler 接口
	router.ServeHTTP(w, req)
	return w
}

func PostJson(uri string, param map[string]interface{}, router http.Handler) *httptest.ResponseRecorder {
	return RequestJson(uri, "POST", param, router)
}

func DeleteJson(uri string, ids []int, router http.Handler) *httptest.ResponseRecorder {
	jsonByte, _ := json.Marshal(ids)
	req := httptest.NewRequest("DELETE", uri, bytes.NewReader(jsonByte))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func PutJson(uri string, param map[string]interface{}, router http.Handler) *httptest.ResponseRecorder {
	return RequestJson(uri, "PUT", param, router)
}

func RequestJson(uri, method string, param map[string]any, router http.Handler) *httptest.ResponseRecorder {
	// 将参数转化为 JSON 比特流
	jsonByte, _ := json.Marshal(param)
	// 构造指定方法的请求，JSON 数据以请求 Body 的形式传递
	req := httptest.NewRequest(method, uri, bytes.NewReader(jsonByte))
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的 handler 接口
	router.ServeHTTP(w, req)
	return w
}
