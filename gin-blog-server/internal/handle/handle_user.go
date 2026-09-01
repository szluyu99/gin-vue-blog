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

// @Summary 获取当前用户信息
// @Description 根据 Token 获取用户信息与点赞记录
// @Tags User
// @Produce json
// @Success 0 {object} Response[model.UserInfoVO]
// @Security ApiKeyAuth
// @Router /user/info [get]
// @Router /front/user/info [get]
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
// @Summary 修改当前用户信息
// @Description 修改昵称/头像/简介/网站, 用户 ID 从 Token 中解析
// @Tags User
// @Accept json
// @Produce json
// @Param form body UpdateCurrentUserReq true "修改当前用户信息"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /user/current [put]
// @Router /front/user/info [put]
func (*User) UpdateCurrent(c *gin.Context) {
	var req UpdateCurrentUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}
	err := model.UpdateUserInfo(GetDB(c), auth.UserInfoId, req.Nickname, req.Avatar, req.Intro, req.Website)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// @Summary 修改用户信息
// @Description 修改指定用户的昵称与角色
// @Tags User
// @Accept json
// @Produce json
// @Param form body UpdateUserReq true "修改用户信息"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /user [put]
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

// @Summary 条件查询用户列表
// @Description 支持按登录类型/用户名/昵称过滤
// @Tags User
// @Produce json
// @Param login_type query int false "登录类型"
// @Param username query string false "用户名"
// @Param nickname query string false "昵称"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.UserAuth]]
// @Security ApiKeyAuth
// @Router /user/list [get]
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

// @Summary 修改用户禁用状态
// @Description 禁用或启用用户
// @Tags User
// @Accept json
// @Produce json
// @Param form body UpdateUserDisableReq true "修改禁用状态"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /user/disable [put]
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

// @Summary 修改当前用户密码
// @Description 需要提供旧密码进行验证
// @Tags User
// @Accept json
// @Produce json
// @Param form body UpdateCurrentPasswordReq true "修改密码"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /user/current/password [put]
func (*User) UpdateCurrentPassword(c *gin.Context) {
	var req UpdateCurrentPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

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

// @Summary 获取在线用户列表
// @Description 从 Redis 中读取在线用户, 按上次登录时间倒序
// @Tags User
// @Produce json
// @Param keyword query string false "用户名或昵称关键字"
// @Success 0 {object} Response[[]model.UserAuth]
// @Security ApiKeyAuth
// @Router /user/online [get]
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

// @Summary 强制用户下线
// @Description 强制指定用户下线, 不能强制自己下线
// @Tags User
// @Produce json
// @Param id path int true "用户 ID"
// @Success 0 {object} Response[string]
// @Security ApiKeyAuth
// @Router /user/offline/{id} [post]
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
