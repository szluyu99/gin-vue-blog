package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type CommentQuery struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *bool  `form:"is_review"`
	Type     int    `form:"type"`
}

type Comment struct{}

// @Summary 删除评论（批量）
// @Description 根据 ID 数组删除评论
// @Tags Comment
// @Param ids body []int true "评论 ID 数组"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /comment [delete]
func (*Comment) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(model.Comment{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}

// @Summary 修改评论审核（批量）
// @Description 根据 ID 数组修改审核状态
// @Tags Comment
// @Param form body UpdateReviewReq true "修改审核状态"
// @Accept json
// @Produce json
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /comment/review [put]
func (*Comment) UpdateReview(c *gin.Context) {
	var req UpdateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	maps := map[string]any{"is_review": req.IsReview}
	result := GetDB(c).Model(model.Comment{}).Where("id in ?", req.Ids).Updates(maps)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}

// @Summary 条件查询评论列表
// @Description 根据条件查询评论列表
// @Tags Comment
// @Param nickname query string false "昵称"
// @Param is_review query int false "审核状态"
// @Param type query int false "评论类型"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.CommentVO]
// @Security ApiKeyAuth
// @Router /comment [get]
func (*Comment) GetList(c *gin.Context) {
	var query CommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	list, total, err := model.GetCommentList(GetDB(c), query.Page, query.Size, query.Type, query.IsReview, query.Nickname)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.Comment]{
		Total: total,
		List:  list,
		Size:  query.Size,
		Page:  query.Page,
	})

}
