package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Link struct{}

// 添加或修改友链
type AddOrEditLinkReq struct {
	ID      int    `json:"id"`
	Name    string `json:"name" binding:"required"`
	Avatar  string `json:"avatar"`
	Address string `json:"address" binding:"required"`
	Intro   string `json:"intro"`
}

// @Summary 条件查询友链列表
// @Description 关键字匹配名称/地址/简介
// @Tags Link
// @Produce json
// @Param keyword query string false "关键字"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.FriendLink]]
// @Security ApiKeyAuth
// @Router /link/list [get]
func (*Link) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := model.GetLinkList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.FriendLink]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 新增或编辑友链
// @Description 新增或编辑友链
// @Tags Link
// @Accept json
// @Produce json
// @Param form body AddOrEditLinkReq true "新增或编辑友链"
// @Success 0 {object} Response[model.FriendLink]
// @Security ApiKeyAuth
// @Router /link [post]
func (*Link) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	link, err := model.SaveOrUpdateLink(GetDB(c), req.ID, req.Name, req.Avatar, req.Address, req.Intro)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, link)
}

// @Summary 删除友链（批量）
// @Description 根据 ID 数组删除友链
// @Tags Link
// @Accept json
// @Produce json
// @Param ids body []int true "友链 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /link [delete]
func (*Link) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(&model.FriendLink{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
