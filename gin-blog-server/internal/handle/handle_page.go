package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type AddOrEditPageReq struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Cover string `json:"cover"`
}

type Page struct{}

// @Summary 获取页面列表
// @Description 根据条件查询获取页面列表
// @Tags Page
// @Accept json
// @Produce json
// @Success 0 {object} Response[[]model.Page]
// @Security ApiKeyAuth
// @Router /page/list [get]
func (*Page) GetList(c *gin.Context) {
	db := GetDB(c)

	data, _, err := model.GetPageList(db)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	ReturnSuccess(c, data)
}

// @Summary 添加或修改页面
// @Description 添加或修改页面
// @Tags Page
// @Param form body AddOrEditPageReq true "添加或修改页面"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.Page]
// @Security ApiKeyAuth
// @Router /page [post]
func (*Page) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	// TODO: 有必要检查吗? 可以依靠 MySQL unique 约束来保证
	// // 检查页面名称是否已经存在
	// page, err := model.GetPageByName(db, req.Name)
	// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	Error(c, g.ERROR_DB_OPERATION, err)
	// 	return
	// }
	// if page != nil && page.ID != 0 && page.ID != req.ID {
	// 	Error(c, g.ERROR_PAGE_NAME_EXIST, nil)
	// 	return
	// }

	// // 检查页面标签是否已经存在
	// page, err = model.GetPageByLabel(db, req.Label)
	// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	Error(c, g.ERROR_DB_OPERATION, err)
	// 	return
	// }
	// if page != nil && page.ID != 0 && page.ID != req.ID {
	// 	Error(c, g.ERROR_PAGE_LABEL_EXIST, nil)
	// 	return
	// }

	page, err := model.SaveOrUpdatePage(db, req.ID, req.Name, req.Label, req.Cover)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	err = rdb.Del(ctx(), g.PAGE).Err()
	if err != nil {
		ReturnError(c, g.ERROR_REDIS_OPERATION, err)
		return
	}

	ReturnSuccess(c, page)
}

// @Summary 删除页面
// @Description 根据 ID 数组删除页面
// @Tags Page
// @Param ids body []int true "页面 ID 数组"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /page [delete]
func (*Page) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	result := GetDB(c).Delete(model.Page{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, result.Error)
		return
	}

	err := GetRDB(c).Del(ctx(), g.PAGE).Err()
	if err != nil {
		ReturnError(c, g.ERROR_REDIS_OPERATION, err)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
