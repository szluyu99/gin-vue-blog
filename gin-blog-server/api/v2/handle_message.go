package v2

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

// 修改审核（批量）
type UpdateReviewForm struct {
	Ids      []int `json:"ids"`
	IsReview *int8 `json:"is_review" validate:"required,min=0,max=1"`
}

// 条件查询列表
type MessageQueryForm struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *int8  `form:"is_review"`
}

type Message struct{}

// @Summary 删除留言（批量）
// @Description 根据 ID 数组删除留言
// @Tags Category
// @Param ids body []int true "留言 ID 数组"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[any]
// @Security ApiKeyAuth
// @Router /category [delete]
func (*Message) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM, err)
		return
	}

	dao.Delete(model.Category{}, "id in ?", ids)
	r.Success(c, nil)
}

// @Summary 修改留言审核（批量）
// @Description 根据 ID 数组修改审核状态
// @Tags Message
// @Param form body UpdateReviewForm true "修改审核状态"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[any]
// @Security ApiKeyAuth
// @Router /message/review [put]
func (*Message) UpdateReview(c *gin.Context) {
	var form UpdateReviewForm
	if err := c.ShouldBindJSON(&form); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM, err)
		return
	}

	maps := map[string]any{"is_review": form.IsReview}
	dao.UpdatesMap(&model.Message{}, maps, "id IN ?", form.Ids)

	r.Success(c, nil)
}

// @Summary 条件查询留言列表
// @Description 根据条件查询留言列表
// @Tags Message
// @Param nickname query string false "昵称"
// @Param is_review query int false "审核状态"
// @Param page_size query int false "当前页数"
// @Param page_num query int false "每页条数"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[PageResult[[]model.Message]]
// @Security ApiKeyAuth
// @Router /message/list [get]
func (*Message) GetList(c *gin.Context) {
	var form MessageQueryForm
	if err := c.ShouldBindQuery(&form); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM, err)
		return
	}

	data, total := messageDao.GetList(form.PageNum, form.PageSize, form.Nickname, form.IsReview)
	resp := PageResult[[]model.Message]{
		Total: total,
		List:  data,
	}
	r.Success(c, resp)
}
