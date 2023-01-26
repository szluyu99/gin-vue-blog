package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type User struct{}

// 更新当前用户信息
func (*User) UpdateCurrent(c *gin.Context) {
	currentUser := utils.BindValidJson[req.UpdateCurrentUser](c)
	currentUser.ID = utils.GetFromContext[int](c, "user_info_id")
	r.SendCode(c, userService.UpdateCurrent(currentUser))
}

// 更新当前用户信息: 关联处理 [用户-角色] 关系
func (*User) Update(c *gin.Context) {
	r.SendCode(c, userService.Update(utils.BindJson[req.UpdateUser](c)))
}

func (*User) GetList(c *gin.Context) {
	r.SuccessData(c, userService.GetList(utils.BindQuery[req.GetUsers](c)))
}

// 修改用户禁用状态
func (*User) UpdateDisable(c *gin.Context) {
	req := utils.BindValidJson[req.UpdateUserDisable](c)
	userService.UpdateDisable(req.ID, *req.IsDisable)
	r.Success(c)
}

// 修改普通用户密码: 管理员可以直接修改
func (*User) UpdatePassword(c *gin.Context) {
	r.SendCode(c, userService.UpdatePassword(utils.BindJson[req.UpdatePassword](c)))
}

// 修改当前用户密码: 需要输入旧密码进行验证
func (*User) UpdateCurrentPassword(c *gin.Context) {
	r.SendCode(c, userService.UpdateCurrentPassword(
		utils.BindJson[req.UpdateAdminPassword](c),
		utils.GetFromContext[int](c, "user_info_id")))
}

// 获取在线用户列表
func (*User) GetOnlineList(c *gin.Context) {
	r.SuccessData(c, userService.GetOnlineList(utils.BindPageQuery(c)))
}

// 强制离线
func (*User) ForceOffline(c *gin.Context) {
	// TODO: 自己不能离线自己
	// userInfoId := utils.GetFromContext[int](c, "user_info_id")
	r.SendCode(c, userService.ForceOffline(utils.BindJson[req.ForceOfflineUser](c)))
}
