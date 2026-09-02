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
// @Accept json
// @Produce json
// @Param ids body []int true "评论 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /comment [delete]
func (*Comment) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, deletedIds, err := model.DeleteComments(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 与文章删除保持一致: 清掉点赞计数, 否则新评论复用 id 会继承旧数据
	cleanCommentCounters(GetRDB(c), deletedIds)

	ReturnSuccess(c, rows)
}

// @Summary 修改评论审核（批量）
// @Description 根据 ID 数组修改审核状态
// @Tags Comment
// @Accept json
// @Produce json
// @Param form body UpdateReviewReq true "修改审核状态"
// @Success 0 {object} Response[int64]
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
// @Description 支持按昵称/审核状态/类型过滤
// @Tags Comment
// @Produce json
// @Param nickname query string false "昵称"
// @Param is_review query bool false "审核状态"
// @Param type query int false "评论类型(1-文章 2-友链 3-说说)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.Comment]]
// @Security ApiKeyAuth
// @Router /comment/list [get]
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
