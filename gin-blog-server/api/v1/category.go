package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Category struct{}

// 新增/编辑
func (*Category) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, categoryService.SaveOrUpdate(utils.BindValidJson[req.AddOrEditCategory](c)))
}

// 删除(批量)
func (*Category) Delete(c *gin.Context) {
	r.SendCode(c, categoryService.Delete(utils.BindJson[[]int](c)))
}

// 获取下拉框选项数据
func (*Category) GetOption(c *gin.Context) {
	r.SuccessData(c, categoryService.GetOption())
}

// 条件查询列表(后台)
func (*Category) GetList(c *gin.Context) {
	r.SuccessData(c, categoryService.GetList(utils.BindPageQuery(c)))
}
