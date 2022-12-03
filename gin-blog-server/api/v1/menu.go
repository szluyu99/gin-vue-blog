package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Menu struct{}

// 获取当前用户菜单: 生成后台管理界面的菜单
func (*Menu) GetUserMenu(c *gin.Context) {
	r.SuccessData(c, menuService.GetUserMenuTree(
		utils.GetFromContext[int](c, "user_info_id")))
}

// 菜单列表(树形)
func (*Menu) GetTreeList(c *gin.Context) {
	r.SuccessData(c, menuService.GetTreeList(utils.BindQuery[req.PageQuery](c)))
}

func (*Menu) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, menuService.SaveOrUpdate(utils.BindValidJson[req.SaveOrUpdateMenu](c)))
}

func (*Menu) Delete(c *gin.Context) {
	r.SendCode(c, menuService.Delete(utils.GetIntParam(c, "id")))
}

// 菜单选项数据
func (*Menu) GetOption(c *gin.Context) {
	r.SuccessData(c, menuService.GetOptionList())
}
