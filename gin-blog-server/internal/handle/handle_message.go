package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

// 修改审核（批量）
type UpdateReviewReq struct {
	Ids      []int `json:"ids"`
	IsReview bool  `json:"is_review"`
}

// 条件查询列表
type MessageQuery struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *bool  `form:"is_review"`
}

type Message struct{}

// @Summary 删除留言（批量）
// @Description 根据 ID 数组删除留言
// @Tags Message
// @Accept json
// @Produce json
// @Param ids body []int true "留言 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /message [delete]
func (*Message) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.DeleteMessages(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// @Summary 修改留言审核（批量）
// @Description 根据 ID 数组修改审核状态
// @Tags Message
// @Accept json
// @Produce json
// @Param form body UpdateReviewReq true "修改审核状态"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /message/review [put]
func (*Message) UpdateReview(c *gin.Context) {
	var req UpdateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.UpdateMessagesReview(GetDB(c), req.Ids, req.IsReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// @Summary 条件查询留言列表
// @Description 支持按昵称与审核状态过滤
// @Tags Message
// @Produce json
// @Param nickname query string false "昵称"
// @Param is_review query bool false "审核状态"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.Message]]
// @Security ApiKeyAuth
// @Router /message/list [get]
func (*Message) GetList(c *gin.Context) {
	var query MessageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := model.GetMessageList(GetDB(c), query.Page, query.Size, query.Nickname, query.IsReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, PageResult[model.Message]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}
