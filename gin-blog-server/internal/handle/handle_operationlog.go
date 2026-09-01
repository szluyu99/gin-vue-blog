package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type OperationLog struct{}

// @Summary 条件查询操作日志列表
// @Description 关键字匹配操作模块与操作描述
// @Tags OperationLog
// @Produce json
// @Param keyword query string false "关键字"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.OperationLog]]
// @Security ApiKeyAuth
// @Router /operation/log/list [get]
func (*OperationLog) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, total, err := model.GetOperationLogList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.OperationLog]{
		Total: total,
		List:  list,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 删除操作日志（批量）
// @Description 根据 ID 数组删除操作日志
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param ids body []int true "操作日志 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /operation/log [delete]
func (*OperationLog) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(&model.OperationLog{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
