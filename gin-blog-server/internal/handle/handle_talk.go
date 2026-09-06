package handle

import (
	"strconv"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Talk struct{}

type AddOrEditTalkReq struct {
	ID      int    `json:"id"`
	Content string `json:"content" binding:"required"`
	// 1-公开 2-私密, 不传按公开处理
	Status int  `json:"status" binding:"omitempty,oneof=1 2"`
	IsTop  bool `json:"is_top"`
}

type TalkQuery struct {
	PageQuery
	Status int `form:"status" binding:"omitempty,oneof=1 2"`
}

// @Summary 说说列表
// @Description 后台说说列表, 可按状态筛选
// @Tags Talk
// @Produce json
// @Param status query int false "状态(1公开 2私密)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.Talk]]
// @Security ApiKeyAuth
// @Router /talk/list [get]
func (*Talk) GetList(c *gin.Context) {
	var query TalkQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := model.GetTalkList(GetDB(c), query.Page, query.Size, query.Status)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.Talk]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 说说详情
// @Description 后台说说详情, 私密的也能取到(要能编辑)
// @Tags Talk
// @Produce json
// @Param id path int true "说说 ID"
// @Success 0 {object} Response[model.Talk]
// @Security ApiKeyAuth
// @Router /talk/{id} [get]
func (*Talk) GetDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	talk, err := model.GetTalk(GetDB(c), id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, talk)
}

// @Summary 新增或编辑说说
// @Description 新增或编辑说说
// @Tags Talk
// @Accept json
// @Produce json
// @Param form body AddOrEditTalkReq true "新增或编辑说说"
// @Success 0 {object} Response[model.Talk]
// @Security ApiKeyAuth
// @Router /talk [post]
func (*Talk) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditTalkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 发布者取当前登录用户, 不信前端传的
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	status := req.Status
	if status == 0 {
		status = model.STATUS_PUBLIC
	}

	talk, err := model.SaveOrUpdateTalk(GetDB(c), req.ID, auth.ID, req.Content, status, req.IsTop)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, talk)
}

// @Summary 删除说说（批量）
// @Description 根据 ID 数组删除说说, 同时删掉挂在它们下面的评论
// @Tags Talk
// @Accept json
// @Produce json
// @Param ids body []int true "说说 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /talk [delete]
func (*Talk) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.DeleteTalks(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}
