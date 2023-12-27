package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"gin-blog/internal/utils/jwt"
	"log/slog"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserAuth struct{}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
	Code     string `json:"code" binding:"required"`
}

// @Summary 登录
// @Description 登录
// @Tags UserAuth
// @Param form body LoginReq true "登录"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.LoginVO]
// @Router /login [post]
func (*UserAuth) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	userAuth, err := model.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ERROR_USER_NOT_EXIST, nil)
			return
		}
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	// 检查密码是否正确
	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		ReturnError(c, g.ERROR_PASSWORD_WRONG, nil)
		return
	}

	// 获取 IP 相关信息 FIXME: 好像无法读取到 ip 信息
	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSourceSimpleIdle(ipAddress)

	// browser, os := "unknown", "unknown"
	// if userAgent := utils.IP.GetUserAgent(c); userAgent != nil {
	// 	browser = userAgent.Name + " " + userAgent.Version.String()
	// 	os = userAgent.OS + " " + userAgent.OSVersion.String()
	// }

	userInfo, err := model.GetUserInfoById(db, userAuth.UserInfoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ERROR_USER_NOT_EXIST, nil)
			return
		}
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	roleIds, err := model.GetRoleIdsByUserId(db, userAuth.ID)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	articleLikeSet, err := rdb.SMembers(ctx(), g.ARTICLE_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}
	commentLikeSet, err := rdb.SMembers(ctx(), g.COMMENT_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	// 登录信息正确, 生成 Token

	// UUID 生成方法: ip + 浏览器信息 + 操作系统信息
	// uuid := utils.MD5(ipAddress + browser + os)
	conf := g.Conf.JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), userAuth.ID, roleIds)
	if err != nil {
		ReturnError(c, g.ERROR_TOKEN_CREATE, err)
		return
	}

	// 更新用户验证信息: ip 信息 + 上次登录时间
	err = model.UpdateUserLoginInfo(db, userAuth.ID, ipAddress, ipSource)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	slog.Info("用户登录成功: " + userAuth.Username)

	session := sessions.Default(c)
	session.Set(g.CTX_USER_AUTH, userAuth.ID)
	session.Save()

	ReturnSuccess(c, model.LoginVO{
		// user info
		UserInfoId: userInfo.ID,
		Nickname:   userInfo.Nickname,
		Email:      userInfo.Email,
		Avatar:     userInfo.Avatar,
		Intro:      userInfo.Intro,
		Website:    userInfo.Website,

		// user auth
		UserAuthId:    userAuth.ID,
		Username:      userAuth.Username,
		LastLoginTime: userAuth.LastLoginTime,
		LoginType:     userAuth.LoginType,

		IpAddress:      ipAddress,
		IpSource:       ipSource,
		ArticleLikeSet: articleLikeSet,
		CommentLikeSet: commentLikeSet,
		Token:          token,
	})
}

// @Summary 注册
// @Description 注册
// @Tags UserAuth
// @Param form body RegisterReq true "注册"
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /register [post]
func (*UserAuth) Register(c *gin.Context) {
	ReturnSuccess(c, "注册")
}

// @Summary 退出登录
// @Description 退出登录
// @Tags UserAuth
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /logout [post]
func (*UserAuth) Logout(c *gin.Context) {
	c.Set(g.CTX_USER_AUTH, nil)

	session := sessions.Default(c)
	session.Delete(g.CTX_USER_AUTH)
	session.Save()

	ReturnSuccess(c, nil)
}

// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码
// @Tags UserAuth
// @Param email query string true "邮箱"
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /code [get]
func (*UserAuth) SendCode(c *gin.Context) {
	ReturnSuccess(c, "发送邮箱验证码")
}

// TODO: refresh token
