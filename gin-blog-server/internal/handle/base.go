package handle

import (
	"context"
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

/*
响应设计方案：不使用 HTTP 码来表示业务状态, 采用业务状态码的方式
- 只要能到达后端的请求, HTTP 状态码都为 200
- 业务状态码为 0 表示成功, 其他都表示失败
- 当后端发生 panic 并且被 gin 中间件捕获时, 才会返回 HTTP 500 状态码
*/

// 响应结构体
type Response[T any] struct {
	Code    int    `json:"code"`    // 业务状态码
	Message string `json:"message"` // 响应消息
	Data    T      `json:"data"`    // 响应数据
}

// HTTP 码 + 业务码 + 消息 + 数据
func ReturnHttpResponse(c *gin.Context, httpCode, code int, msg string, data any) {
	c.JSON(httpCode, Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// 业务码 + 数据
func ReturnResponse(c *gin.Context, code int, data any) {
	ReturnHttpResponse(c, http.StatusOK, code, g.GetMsg(code), data)
}

// 成功业务码 + 数据
func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, g.SUCCESS, data)
}

// 处理错误
func ReturnError(c *gin.Context, code int, err error) {
	var msg string

	if err == nil {
		msg = g.GetMsg(code)
	} else {
		msg = err.Error()
	}
	slog.Error(msg) // 记录错误日志

	c.AbortWithStatusJSON(
		http.StatusOK,
		Response[string]{
			Code:    code,
			Message: msg,
			Data:    msg,
		},
	)
}

// 分页获取数据
type PageQuery struct {
	Page    int    `form:"page_num"`  // 当前页数（从1开始）
	Size    int    `form:"page_size"` // 每页条数
	Keyword string `form:"keyword"`   // 搜索关键字
}

// 分页响应数据
type PageResult[T any] struct {
	Page  int   `json:"page_num"`  // 每页条数
	Size  int   `json:"page_size"` // 上次页数
	Total int64 `json:"total"`     // 总条数
	List  []T   `json:"page_data"` // 分页数据
}

// 获取 *gorm.DB
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(g.CTX_DB).(*gorm.DB)
}

// 获取 *redis.Client
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(g.CTX_RDB).(*redis.Client)
}

// 获取 context.Context
func ctx() context.Context {
	return context.Background()
}

/*
获取当前登录用户信息
1. 能从 gin Context 上获取到 user 对象, 说明本次请求链路中获取过了
2. 从 session 中获取到 uid
3. 根据 uid 获取用户信息, 并设置到 gin Context 上
*/
func CurrentUserAuth(c *gin.Context) (*model.UserAuth, error) {
	key := g.CTX_USER_AUTH

	// 1
	if cache, exist := c.Get(key); exist {
		slog.Debug("[Func-CurrentUserAuth] get from cache: " + cache.(*model.UserAuth).Username)
		return cache.(*model.UserAuth), nil
	}

	// 2
	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("session 中没有 user_auth_id")
	}

	// 3
	db := GetDB(c)
	user, err := model.GetUserAuthInfoById(db, id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}
