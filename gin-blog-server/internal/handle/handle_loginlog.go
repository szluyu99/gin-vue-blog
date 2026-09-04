package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type LoginLog struct{}

// @Summary 条件查询登录日志列表
// @Description 关键字匹配用户名、昵称与登录 IP
// @Tags LoginLog
// @Produce json
// @Param keyword query string false "关键字"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.LoginLog]]
// @Security ApiKeyAuth
// @Router /login/log/list [get]
func (*LoginLog) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, total, err := model.GetLoginLogList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.LoginLog]{
		Total: total,
		List:  list,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 删除登录日志（批量）
// @Description 根据 ID 数组删除登录日志
// @Tags LoginLog
// @Accept json
// @Produce json
// @Param ids body []int true "登录日志 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /login/log [delete]
func (*LoginLog) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(&model.LoginLog{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
