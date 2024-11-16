package handle

import (
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct{}

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

type ForceOfflineReq struct {
	UserInfoId int `json:"user_info_id"`
}

// 根据 Token 获取用户信息
func (*User) GetInfo(c *gin.Context) {
	rdb := GetRDB(c)

	user, err := CurrentUserAuth(c)
	if err != nil {
		ReturnError(c, g.ErrTokenRuntime, err)
		return
	}

	userInfoVO := model.UserInfoVO{UserInfo: *user.UserInfo}
	userInfoVO.ArticleLikeSet, err = rdb.SMembers(rctx, g.ARTICLE_USER_LIKE_SET+strconv.Itoa(user.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	userInfoVO.CommentLikeSet, err = rdb.SMembers(rctx, g.COMMENT_USER_LIKE_SET+strconv.Itoa(user.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, userInfoVO)
}

// TODO: 用户区域分布 GetUserAreas, StatisticUserAreas
// 更新当前用户信息, 不需要传 id, 从 Token 中解析出来
func (*User) UpdateCurrent(c *gin.Context) {
	var req UpdateCurrentUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, _ := CurrentUserAuth(c)
	err := model.UpdateUserInfo(GetDB(c), auth.UserInfoId, req.Nickname, req.Avatar, req.Intro, req.Website)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 更新用户信息: 昵称 + 角色
func (*User) Update(c *gin.Context) {
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if err := model.UpdateUserNicknameAndRole(GetDB(c), req.UserAuthId, req.Nickname, req.RoleIds); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 获取用户列表
func (*User) GetList(c *gin.Context) {
	var query UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, count, err := model.GetUserList(GetDB(c), query.Page, query.Size, query.LoginType, query.Nickname, query.Username)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
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
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.UpdateUserDisable(GetDB(c), req.UserAuthId, req.IsDisable)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 修改当前用户密码: 需要输入旧密码进行验证
func (*User) UpdateCurrentPassword(c *gin.Context) {
	var req UpdateCurrentPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, _ := CurrentUserAuth(c)

	if !utils.BcryptCheck(req.OldPassword, auth.Password) {
		ReturnError(c, g.ErrOldPassword, nil)
		return
	}

	hashPassword, _ := utils.BcryptHash(req.NewPassword)
	err := model.UpdateUserPassword(GetDB(c), auth.ID, hashPassword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
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
// 		ReturnError(c, g.ErrRequestParam, err)
// 		return
// 	}

// 	hashPassword, err := utils.BcryptHash(form.Password)
// 	if err != nil {
// 		ReturnError(c, g.FAIL, err)
// 		return
// 	}
// 	err = model.UpdateUserPassword(GetDB(c), form.Username, hashPassword)
// 	if err != nil {
// 		ReturnError(c, g.ErrDbOperatioin, err)
// 		return
// 	}

// 	// TODO: 修改完密码后，强制用户下线

// 	ReturnSuccess(c, nil)
// }

// 查询当前在线用户
func (*User) GetOnlineList(c *gin.Context) {
	keyword := c.Query("keyword")

	rdb := GetRDB(c)

	onlineList := make([]model.UserAuth, 0)
	keys := rdb.Keys(rctx, g.ONLINE_USER+"*").Val()

	for _, key := range keys {
		var auth model.UserAuth
		val := rdb.Get(rctx, key).Val()
		json.Unmarshal([]byte(val), &auth)

		if keyword != "" &&
			!strings.Contains(auth.Username, keyword) &&
			!strings.Contains(auth.UserInfo.Nickname, keyword) {
			continue
		}

		onlineList = append(onlineList, auth)
	}

	// 根据上次登录时间进行排序
	sort.Slice(onlineList, func(i, j int) bool {
		return onlineList[i].LastLoginTime.Unix() > onlineList[j].LastLoginTime.Unix()
	})

	ReturnSuccess(c, onlineList)
}

// 强制离线
func (*User) ForceOffline(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.Atoi(id)
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, err := CurrentUserAuth(c)
	if err != nil {
		ReturnError(c, g.ErrUserAuth, err)
		return
	}

	// 不能离线自己
	if auth.ID == uid {
		ReturnError(c, g.ErrForceOfflineSelf, nil)
		return
	}

	rdb := GetRDB(c)
	onlineKey := g.ONLINE_USER + strconv.Itoa(uid)
	offlineKey := g.OFFLINE_USER + strconv.Itoa(uid)

	rdb.Del(rctx, onlineKey)
	rdb.Set(rctx, offlineKey, auth, time.Hour)

	ReturnSuccess(c, "强制离线成功")
}
