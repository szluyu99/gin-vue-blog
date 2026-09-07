package handle

import (
	"log/slog"
	"time"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ErrorLog struct{}

// 上报请求: 除 message 外都可缺省, 前端能拿到多少给多少
type ReportErrorReq struct {
	Source  string `json:"source"`
	Message string `json:"message" binding:"required"`
	Stack   string `json:"stack"`
	Url     string `json:"url"`
}

const (
	errorReportQuota  = 10              // 同一 IP 在一个窗口内允许上报的条数
	errorReportWindow = 1 * time.Minute // 配额窗口
)

/*
按 IP 领配额

上报接口必须匿名可访问(访客没登录态), 所以要在入口挡住刷量。
写法照 handle_auth.go 的登录失败计数: Incr + 首次设置 TTL。
Redis 出问题时放行 —— 丢一条错误日志不算事, 因为它拦住真实上报更糟。
*/
func errorReportAllowed(rdb *redis.Client, key string) bool {
	count, err := rdb.Incr(rctx, key).Result()
	if err != nil {
		slog.Warn("前端错误上报计数出错, 本次放行", "err", err)
		return true
	}
	if count == 1 {
		if err := rdb.Expire(rctx, key, errorReportWindow).Err(); err != nil {
			slog.Warn("设置前端错误上报配额有效期出错", "err", err)
		}
	}
	return count <= errorReportQuota
}

// @Summary 上报前端错误
// @Description 前端未捕获异常上报, 允许匿名访问; 同类错误按指纹聚合只累加次数
// @Tags Front
// @Accept json
// @Produce json
// @Param form body ReportErrorReq true "错误信息"
// @Success 0 {object} Response[any]
// @Router /front/error/report [post]
func (*ErrorLog) Report(c *gin.Context) {
	var req ReportErrorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	ip := clientIP(c)

	/*
		超配额时返回成功而不是错误

		前端的上报逻辑挂在 window.onerror 上, 这里返回错误会让它的错误处理
		再报一次错, 形成放大循环。静默丢弃是这里唯一安全的选择。
	*/
	if !errorReportAllowed(GetRDB(c), g.ERROR_REPORT+ip) {
		slog.Warn("前端错误上报超出配额, 已丢弃", "ip", ip)
		ReturnSuccess(c, nil)
		return
	}

	// 登录用户记下 id, 便于复现时找人问; 匿名访客留 0
	var userId int
	if auth, err := CurrentUserAuth(c); err == nil && auth != nil {
		userId = auth.ID
	}

	log := &model.ErrorLog{
		Source:    req.Source,
		Message:   req.Message,
		Stack:     req.Stack,
		Url:       req.Url,
		UserAgent: c.Request.UserAgent(),
		IpAddress: ip,
		UserId:    userId,
	}

	db := GetDB(c)
	if err := model.AddErrorLog(db, log); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 总量兜底: 变化的 message 能绕过聚合和限频, 行数仍会涨。
	// 裁剪失败只告警, 不影响上报本身
	if err := model.TrimErrorLogs(db); err != nil {
		slog.Warn("裁剪前端错误日志失败", "err", err)
	}

	ReturnSuccess(c, nil)
}

type ErrorLogQuery struct {
	PageQuery
	Source string `form:"source" binding:"omitempty,oneof=front admin"`
}

// @Summary 前端错误日志列表
// @Description 按最近出现时间倒序, 关键字匹配错误信息与页面地址
// @Tags ErrorLog
// @Produce json
// @Param keyword query string false "关键字"
// @Param source query string false "来源(front/admin)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.ErrorLog]]
// @Security ApiKeyAuth
// @Router /error/log/list [get]
func (*ErrorLog) GetList(c *gin.Context) {
	var query ErrorLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, total, err := model.GetErrorLogList(GetDB(c), query.Page, query.Size, query.Keyword, query.Source)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.ErrorLog]{
		Total: total,
		List:  list,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 删除前端错误日志（批量）
// @Description 根据 ID 数组删除
// @Tags ErrorLog
// @Accept json
// @Produce json
// @Param ids body []int true "错误日志 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /error/log [delete]
func (*ErrorLog) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.DeleteErrorLogs(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}
