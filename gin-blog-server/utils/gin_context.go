package utils

import (
	"gin-blog/model/req"
	"gin-blog/utils/r"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 在内部 panic 会被 middleware 捕获到并返回错误信息

// 参数合法性校验
func Validate(c *gin.Context, data any) {
	validMsg := Validator.Validate(data)
	if validMsg != "" {
		r.ReturnJson(c, http.StatusOK, r.ERROR_INVALID_PARAM, validMsg, nil)
		panic(nil)
	}
}

// Json 绑定
func BindJson[T any](c *gin.Context) (data T) {
	if err := c.ShouldBindJSON(&data); err != nil {
		Logger.Error("BindJson", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}
	return
}

// Json 绑定验证 + 合法性校验
func BindValidJson[T any](c *gin.Context) (data T) {
	// Json 绑定
	if err := c.ShouldBindJSON(&data); err != nil {
		Logger.Error("BindValidJson", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}
	// 参数合法性校验
	Validate(c, &data)
	return data
}

// Param 绑定
func BindQuery[T any](c *gin.Context) (data T) {
	if err := c.ShouldBindQuery(&data); err != nil {
		Logger.Error("BindQuery", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}

	// TODO: 检查是否有 PageSize 或 PageQuery 字段，并处理其值
	// val := reflect.ValueOf(data)
	// pageSize := val.FieldByName("PageSize").Int()
	// fmt.Println("pageSize: ", pageSize)
	// val.FieldByName("PageSize").Elem().SetInt(12)
	return
}

// Param 分页绑定(处理了 PageSize 和 PageQuery)
func BindPageQuery(c *gin.Context) (data req.PageQuery) {
	if err := c.ShouldBindQuery(&data); err != nil {
		Logger.Error("BindQuery", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}
	// 检查分页参数
	CheckQueryPage(&data.PageSize, &data.PageNum)
	return
}

// Param 绑定验证 + 合法性校验
func BindValidQuery[T any](c *gin.Context) (data T) {
	// Query 绑定
	if err := c.ShouldBindQuery(&data); err != nil {
		Logger.Error("BindValidQuery", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}
	// 参数合法性校验
	Validate(c, &data)
	return data
}

// 检查分页参数
func CheckQueryPage(pageSize, pageNum *int) {
	switch {
	case *pageSize >= 100:
		*pageSize = 100
	case *pageSize <= 0:
		*pageSize = 10
	}
	if *pageNum <= 0 {
		*pageNum = 1
	}
}

// 从 Gin Context 上获取值, 该值是 JWT middleware 解析 Token 后设置的
// 如果该值不存在, 说明 Token 有问题
func GetFromContext[T any](c *gin.Context, key string) T {
	val, exist := c.Get(key)
	if !exist {
		panic(r.ERROR_TOKEN_RUNTIME)
	}
	return val.(T)
}

// 从 Context 获取 Int 类型 Param 参数
func GetIntParam(c *gin.Context, key string) int {
	val, err := strconv.Atoi(c.Param(key))
	if err != nil {
		Logger.Error("GetIntParam", zap.Error(err))
		panic(r.ERROR_REQUEST_PARAM)
	}
	return val
}
