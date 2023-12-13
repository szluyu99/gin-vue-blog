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

type Tag struct{}

// @Summary 添加或修改标签
// @Description 添加或修改标签
// @Tags Tag
// @Param form body req.AddOrEditTag true "添加或修改标签"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[any]
// @Router /tag [post]
func (*Tag) SaveOrUpdate(c *gin.Context) {
	var form req.AddOrEditTag
	if err := c.ShouldBindJSON(&form); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM)
		return
	}
	// 检查标签名称是否已经存在
	exist := dao.GetOne(model.Tag{}, "name", form.Name)
	if !exist.IsEmpty() && exist.ID != form.ID {
		r.Error(c, r.ERROR_TAG_EXIST)
		return
	}

	tag := utils.CopyProperties[model.Tag](form)
	if tag.ID != 0 {
		dao.Update(&tag, "name")
	} else {
		dao.Create(&tag)
	}

	r.Success(c, nil)
}

// @Summary 删除标签
// @Description 根据标签 ID 数组删除标签
// @Tags Tag
// @Param ids body []int true "标签 ID 列表"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[any]
// @Router /tag [delete]
func (*Tag) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM)
		return
	}
	// 检查标签下面有没有文章
	count := dao.Count(model.ArticleTag{}, "tag_id in ?", ids)
	if count > 0 {
		r.Error(c, r.ERROR_TAG_ART_EXIST)
		return
	}
	dao.Delete(model.Tag{}, "id in ?", ids)
	r.Success(c, nil)
}

// @Summary 获取标签列表
// @Description 根据条件查询获取标签列表
// @Tags Tag
// @Param page_size query int false "当前页数"
// @Param page_num query int false "每页条数"
// @Param keyword query string false "搜索关键字"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[resp.TagVO] "成功"
// @Router /tag/list [get]
func (*Tag) GetList(c *gin.Context) {
	var query req.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM)
		return
	}
	data, total := tagDao.GetList(query)
	r.Success(c, resp.PageResult[[]resp.TagVO]{
		Total:    total,
		List:     data,
		PageSize: query.PageSize,
		PageNum:  query.PageNum,
	})
}

// @Summary 获取标签选项列表
// @Description 获取标签选项列表
// @Tags Tag
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[resp.OptionVo]
// @Router /tag/option [get]
func (*Tag) GetOption(c *gin.Context) {
	r.Success(c, tagDao.GetOption())
}
