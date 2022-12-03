package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Role struct{}

// 获取角色选项
func (*Role) GetOption(c *gin.Context) {
	r.SuccessData(c, roleService.GetOption())
}

func (*Role) GetTreeList(c *gin.Context) {
	r.SuccessData(c, roleService.GetTreeList(utils.BindPageQuery(c)))
}

func (*Role) SaveOrUpdate(c *gin.Context) {
	r.SendCode(c, roleService.SaveOrUpdate(utils.BindJson[req.SaveOrUpdateRole](c)))
}

func (*Role) Delete(c *gin.Context) {
	r.SendCode(c, roleService.Delete(utils.BindJson[[]int](c)))
}

// TODO: 修改角色禁用
// func (*Role) UpdateDisable(c *gin.Context) {
// 	r.SendCode(c, roleService.UpdateDisable())
// }
