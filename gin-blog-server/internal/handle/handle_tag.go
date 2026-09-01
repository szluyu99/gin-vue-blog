package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Tag struct{}

type AddOrEditTagReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

// @Summary 条件查询标签列表
// @Description 标签列表, 附带每个标签下的文章数
// @Tags Tag
// @Produce json
// @Param keyword query string false "关键字"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.TagVO]]
// @Security ApiKeyAuth
// @Router /tag/list [get]
func (*Tag) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := model.GetTagList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.TagVO]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 新增或编辑标签
// @Description 标签名称不能重复
// @Tags Tag
// @Accept json
// @Produce json
// @Param form body AddOrEditTagReq true "新增或编辑标签"
// @Success 0 {object} Response[model.Tag]
// @Security ApiKeyAuth
// @Router /tag [post]
func (*Tag) SaveOrUpdate(c *gin.Context) {
	var form AddOrEditTagReq
	if err := c.ShouldBindJSON(&form); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	tag, err := model.SaveOrUpdateTag(GetDB(c), form.ID, form.Name)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, tag)
}

// TODO: 删除行为, 添加强制删除: 有关联数据则将删除关联数据
// @Summary 删除标签（批量）
// @Description 标签下存在文章时不允许删除
// @Tags Tag
// @Accept json
// @Produce json
// @Param ids body []int true "标签 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /tag [delete]
func (*Tag) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	db := GetDB(c)

	// 检查标签下面有没有文章
	count, err := model.Count(db, &model.ArticleTag{}, "tag_id in ?", ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	if count > 0 {
		// ReturnError(c, g.ERROR_TAG_ART_EXIST, nil)
		ReturnError(c, g.ErrTagHasArt, nil)
		return
	}

	result := db.Delete(model.Tag{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}

// @Summary 获取标签选项
// @Description 用于文章选择标签
// @Tags Tag
// @Produce json
// @Success 0 {object} Response[[]model.OptionVO]
// @Security ApiKeyAuth
// @Router /tag/option [get]
func (*Tag) GetOption(c *gin.Context) {
	list, err := model.GetTagOption(GetDB(c))
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}
