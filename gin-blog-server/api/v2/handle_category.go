package v2

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"time"

	"github.com/gin-gonic/gin"
)

// 添加/编辑分类对象
type AddOrEditCategoryForm struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" validate:"required,min=1" label:"分类名称"`
}

type Category struct{}

// @Summary 获取分类列表
// @Description 根据条件查询获取分类列表
// @Tags Category
// @Param page_size query int false "当前页数"
// @Param page_num query int false "每页条数"
// @Param keyword query string false "搜索关键字"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[PageResult[[]model.CategoryVO]] "成功"
// @Security ApiKeyAuth
// @Router /category/list [get]
func (*Category) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM, err)
		return
	}

	data, total, err := model.GetCategoryList(dao.DB, query.PageNum, query.PageSize, query.Keyword)
	if err != nil {
		r.Error(c, r.ERROR_DB_OPE, err)
		return
	}
	resp := PageResult[[]model.CategoryVO]{
		Total:    total,
		List:     data,
		PageSize: query.PageSize,
		PageNum:  query.PageNum,
	}

	r.Success(c, resp)
}

// @Summary 添加或修改分类
// @Description 添加或修改分类
// @Tags Category
// @Param form body AddOrEditCategoryForm true "添加或修改分类"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[any]
// @Security ApiKeyAuth
// @Router /category [post]
func (*Category) SaveOrUpdate(c *gin.Context) {
	var form AddOrEditCategoryForm
	if err := c.ShouldBindJSON(&form); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM, err)
		return
	}

	exist := dao.GetOne(model.Category{}, "name = ?", form.Name)

	// (同名存在) && (存在的id不等于当前要更新的id) -> 重复
	if !exist.IsEmpty() && exist.ID != form.ID {
		r.Error(c, r.ERROR_CATE_NAME_USED, nil)
		return
	}

	category := utils.CopyProperties[model.Category](form)

	if form.ID != 0 {
		dao.Update(&category, "name")
	} else {
		dao.Create(&category)
	}

	r.Success(c, nil)
}

// @Summary 删除分类（批量）
// @Description 根据 ID 数组删除分类
// @Tags Category
// @Param ids body []int true "分类 ID 数组"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[any]
// @Security ApiKeyAuth
// @Router /category [delete]
func (*Category) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM, err)
		return
	}

	// 检查分类下是否存在文章
	count := dao.Count(model.Article{}, "category_id in ?", ids)
	if count > 0 {
		r.Error(c, r.ERROR_CATE_ART_EXIST, nil)
		return
	}

	dao.Delete(model.Category{}, "id in ?", ids)
	r.Success(c, nil)
}

// @Summary 获取分类选项列表
// @Description 获取标签选项列表
// @Tags Category
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[model.OptionVO]
// @Security ApiKeyAuth
// @Router /category/option [get]
func (*Category) GetOption(c *gin.Context) {
	dada, err := model.GetCategoryOption(dao.DB)
	if err != nil {
		r.Error(c, r.ERROR_DB_OPE, err)
		return
	}
	r.Success(c, dada)
}
