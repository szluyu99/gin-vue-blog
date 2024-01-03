package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateCurrentUserReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`
}

type UpdateCurrentPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=4,max=20"`
	OldPassword string `json:"old_password" binding:"required,min=4,max=20"`
}

type UpdateUserReq struct {
	UserAuthId int    `json:"id"`
	Nickname   string `json:"nickname" binding:"required"`
	RoleIds    []int  `json:"role_ids"`
}

type UpdateUserDisableReq struct {
	UserAuthId int  `json:"id"`
	IsDisable  bool `json:"is_disable"`
}

type UserQuery struct {
	PageQuery
	LoginType int8   `form:"login_type"`
	Username  string `form:"username"`
	Nickname  string `form:"nickname"`
}

type OnlineUserVO struct {
	UserInfoId    int       `json:"g.CTX_USER_INFO_ID"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	IpAddress     string    `json:"ip_address"`
	IpSource      string    `json:"ip_source"`
	Browser       string    `json:"browser"`
	OS            string    `json:"os"`
	LastLoginTime time.Time `json:"last_login_time"`
}

type User struct{}

// 根据 Token 获取用户信息
func (*User) GetInfo(c *gin.Context) {
	rdb := GetRDB(c)

	user, err := CurrentUserAuth(c)
	if err != nil {
		ReturnError(c, g.ERROR_TOKEN_RUNTIME, err)
		return
	}

	userInfoVO := model.UserInfoVO{UserInfo: *user.UserInfo}
	userInfoVO.ArticleLikeSet, err = rdb.SMembers(ctx(), g.ARTICLE_USER_LIKE_SET+strconv.Itoa(user.UserInfoId)).Result()
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}
	userInfoVO.CommentLikeSet, err = rdb.SMembers(ctx(), g.COMMENT_USER_LIKE_SET+strconv.Itoa(user.UserInfoId)).Result()
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	ReturnSuccess(c, userInfoVO)
}

// TODO: 用户区域分布 GetUserAreas, StatisticUserAreas
// 更新当前用户信息, 不需要传 id, 从 Token 中解析出来
func (*User) UpdateCurrent(c *gin.Context) {
	var req UpdateCurrentUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	auth, _ := CurrentUserAuth(c)
	err := model.UpdateUserInfo(GetDB(c), auth.UserInfoId, req.Nickname, req.Avatar, req.Intro, req.Website)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 更新用户信息: 昵称 + 角色
func (*User) Update(c *gin.Context) {
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	if err := model.UpdateUserNicknameAndRole(GetDB(c), req.UserAuthId, req.Nickname, req.RoleIds); err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 获取用户列表
func (*User) GetList(c *gin.Context) {
	var query UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	list, count, err := model.GetUserList(GetDB(c), query.Page, query.Size, query.LoginType, query.Nickname, query.Username)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	ReturnSuccess(c, PageResult[model.UserAuth]{
		Size:  query.Size,
		Page:  query.Page,
		Total: count,
		List:  list,
	})
}

// 修改用户禁用状态
func (*User) UpdateDisable(c *gin.Context) {
	var req UpdateUserDisableReq

	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	err := model.UpdateUserDisable(GetDB(c), req.UserAuthId, req.IsDisable)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 修改当前用户密码: 需要输入旧密码进行验证
func (*User) UpdateCurrentPassword(c *gin.Context) {
	var req UpdateCurrentPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	auth, _ := CurrentUserAuth(c)

	if !utils.BcryptCheck(req.OldPassword, auth.Password) {
		ReturnError(c, g.ERROR_OLD_PASSWORD, nil)
		return
	}

	hashPassword, _ := utils.BcryptHash(req.NewPassword)
	err := model.UpdateUserPassword(GetDB(c), auth.ID, hashPassword)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	// TODO: 修改完密码后，强制用户下线

	ReturnSuccess(c, nil)
}

// TODO: 修改普通用户密码（管理员可以直接修改）
// func (*User) UpdatePassword(c *gin.Context) {
// 	type UpdatePasswordForm struct {
// 		Username string `json:"username" validate:"required" label:"用户名"`
// 		Password string `json:"password" validate:"required" label:"密码"`
// 	}

// 	var form UpdatePasswordForm
// 	if err := c.ShouldBindJSON(&form); err != nil {
// 		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
// 		return
// 	}

// 	hashPassword, err := utils.BcryptHash(form.Password)
// 	if err != nil {
// 		ReturnError(c, g.FAIL, err)
// 		return
// 	}
// 	err = model.UpdateUserPassword(GetDB(c), form.Username, hashPassword)
// 	if err != nil {
// 		ReturnError(c, g.ERROR_DB_OPERATION, err)
// 		return
// 	}

// 	// TODO: 修改完密码后，强制用户下线

// 	ReturnSuccess(c, nil)
// }

// 查询当前在线用户
// FIXME: 需要修复
// TODO: 分页 + 条件搜索
func (*User) GetOnlineList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	// onlineList := make([]OnlineUserVO, 0)
	// keys := utils.Redis.Keys(KEY_USER + "*")
	// for _, key := range keys {
	// 	var sessionInfo model.SessionInfo
	// 	utils.Json.Unmarshal(utils.Redis.GetVal(key), &sessionInfo)

	// 	// 查询关键字不为空, 且不满足查询条件
	// 	if query.Keyword != "" && !strings.Contains(sessionInfo.Nickname, query.Keyword) {
	// 		continue
	// 	}

	// 	onlineUser := utils.CopyProperties[OnlineUserVO](sessionInfo)
	// 	onlineUser.UserInfoId = sessionInfo.UserInfoId // *
	// 	onlineList = append(onlineList, onlineUser)
	// }

	// // 根据上次登录时间进行排序
	// sort.Slice(onlineList, func(i, j int) bool {
	// 	return onlineList[i].LastLoginTime.Unix() > onlineList[j].LastLoginTime.Unix()
	// })

	// Success(c, PageResult[[]OnlineUserVO]{
	// 	Total: int64(len(keys)),
	// 	List:  onlineList,
	// })

	ReturnSuccess(c, PageResult[OnlineUserVO]{
		Total: 0,
		List:  make([]OnlineUserVO, 0),
	})
}

// FIXME: 强制离线
func (*User) ForceOffline(c *gin.Context) {
	type ForceOfflineReq struct {
		UserInfoId int    `json:"g.CTX_USER_INFO_ID"`
		IpAddress  string `json:"ip_address"`
		Browser    string `json:"browser"`
		OS         string `json:"os"`
	}

	var req ForceOfflineReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	// TODO: 自己不能离线自己

	// uuid := utils.Encrypt.MD5(form.IpAddress + form.Browser + form.OS)
	// var sessionInfo model.SessionInfo
	// utils.Json.Unmarshal(utils.Redis.GetVal(KEY_USER+uuid), &sessionInfo)
	// sessionInfo.IsOffline = 1 // *
	// utils.Redis.Del(KEY_USER + uuid)
	// // ? 这里设置强制离线后 redis 中存储的 delete:xxx 时间和 Token 过期时间一致
	// utils.Redis.Set(KEY_DELETE+uuid, utils.Json.Marshal(sessionInfo), time.Duration(g.GetConfig().JWT.Expire)*time.Hour)

	// Success(c, nil)
	c.JSON(http.StatusNotImplemented, nil)
}
