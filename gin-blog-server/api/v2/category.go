package v2

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Category struct{}

// 新增/编辑
func (*Category) SaveOrUpdate(c *gin.Context) {
	var form req.AddOrEditCategory
	if err := c.ShouldBindJSON(&form); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM)
		return
	}

	existByName := dao.GetOne(model.Category{}, "name = ?", form.Name)

	// (同名存在) && (存在的id不等于当前要更新的id) -> 重复
	if !existByName.IsEmpty() && existByName.ID != form.ID {
		r.Error(c, r.ERROR_CATE_NAME_USED)
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

// 删除(批量)
func (*Category) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM)
		return
	}

	// 检查分类下是否存在文章
	count := dao.Count(model.Article{}, "category_id in ?", ids)
	if count > 0 {
		r.Error(c, r.ERROR_CATE_ART_EXIST)
		return
	}

	dao.Delete(model.Category{}, "id in ?", ids)
	r.Success(c, nil)
}

// 获取列表(后台)
func (*Category) GetList(c *gin.Context) {
	var query req.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM)
		return
	}

	data, total := categoryDao.GetList(query.PageNum, query.PageSize, query.Keyword)
	resp := resp.PageResult[[]resp.CategoryVo]{
		Total:    total,
		List:     data,
		PageSize: query.PageSize,
		PageNum:  query.PageNum,
	}

	r.Success(c, resp)
}

// 获取下拉框选项数据
func (*Category) GetOption(c *gin.Context) {
	r.Success(c, categoryDao.GetOption())
}
