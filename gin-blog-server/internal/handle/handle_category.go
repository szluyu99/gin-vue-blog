package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Category struct{}

// 添加/编辑分类对象
type AddOrEditCategoryReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

// @Summary 条件查询分类列表
// @Description 分类列表, 附带每个分类下的文章数
// @Tags Category
// @Produce json
// @Param keyword query string false "关键字"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.CategoryVO]]
// @Security ApiKeyAuth
// @Router /category/list [get]
func (*Category) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := model.GetCategoryList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.CategoryVO]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 新增或编辑分类
// @Description 分类名称不能重复
// @Tags Category
// @Accept json
// @Produce json
// @Param form body AddOrEditCategoryReq true "新增或编辑分类"
// @Success 0 {object} Response[model.Category]
// @Security ApiKeyAuth
// @Router /category [post]
func (*Category) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	category, err := model.SaveOrUpdateCategory(GetDB(c), req.ID, req.Name)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, category)
}

// @Summary 删除分类（批量）
// @Description 分类下存在文章时不允许删除
// @Tags Category
// @Accept json
// @Produce json
// @Param ids body []int true "分类 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /category [delete]
func (*Category) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// 检查分类下是否存在文章
	count, err := model.Count(db, &model.Article{}, "category_id in ?", ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	if count > 0 {
		ReturnError(c, g.ErrCateHasArt, nil)
		return
	}

	rows, err := model.DeleteCategory(db, ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, rows)
}

// @Summary 获取分类选项
// @Description 用于文章选择分类
// @Tags Category
// @Produce json
// @Success 0 {object} Response[[]model.OptionVO]
// @Security ApiKeyAuth
// @Router /category/option [get]
func (*Category) GetOption(c *gin.Context) {
	list, err := model.GetCategoryOption(GetDB(c))
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}
