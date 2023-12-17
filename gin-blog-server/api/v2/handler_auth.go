package v2

import (
	"gin-blog/config"
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/dto"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required" label:"密码"`
}

type RegisterForm struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required,min=4,max=20" label:"密码"`
	Code     string `json:"code" validate:"required" label:"邮箱验证码"`
}

type UserAuth struct{}

// @Summary 登录
// @Description 登录
// @Tags UserAuth
// @Param form body LoginForm true "登录"
// @Accept json
// @Produce json
// @Success 0 {object} r.Response[resp.LoginVO]
// @Router /login [post]
func (*UserAuth) Login(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		r.Error(c, r.ERROR_REQUEST_PARAM, err)
		return
	}

	// 检查用户是否存在
	userAuth := dao.GetOne(model.UserAuth{}, "username", form.Username)
	if userAuth.ID == 0 {
		r.Code(c, r.ERROR_USER_NOT_EXIST)
		return
	}

	// 检查密码是否正确
	if !utils.Encryptor.BcryptCheck(form.Password, userAuth.Password) {
		r.Code(c, r.ERROR_PASSWORD_WRONG)
		return
	}

	// 获取用户详细信息 DTO
	userDetailDTO := convertUserDetailDTO(userAuth, c)

	// 登录信息正确, 生成 Token
	// TODO: 目前只给用户设定一个角色, 获取第一个值就行, 后期优化: 给用户设置多个角色
	// UUID 生成方法: ip + 浏览器信息 + 操作系统信息
	uuid := utils.Encryptor.MD5(userDetailDTO.IpAddress + userDetailDTO.Browser + userDetailDTO.OS)
	token, err := utils.GetJWT().GenToken(userAuth.ID, userDetailDTO.RoleLabels[0], uuid)
	if err != nil {
		r.Error(c, r.ERROR_TOKEN_CREATE, err)
		return
	}
	userDetailDTO.Token = token

	// 更新用户验证信息: ip 信息 + 上次登录时间
	dao.Update(&model.UserAuth{
		Universal:     model.Universal{ID: userAuth.ID},
		IpAddress:     userDetailDTO.IpAddress,
		IpSource:      userDetailDTO.IpSource,
		LastLoginTime: userDetailDTO.LastLoginTime,
	}, "ip_address", "ip_source", "last_login_time")

	// 保存用户信息到 Session 和 Redis 中
	session := sessions.Default(c)
	sessionInfoStr := utils.Json.Marshal(dto.SessionInfo{UserDetailDTO: userDetailDTO})
	session.Set(KEY_USER+uuid, sessionInfoStr) // ! 确实设置到 reids 中, 但是获取不到
	utils.Redis.Set(KEY_USER+uuid, sessionInfoStr, time.Duration(config.Cfg.Session.MaxAge)*time.Second)
	session.Save()

	r.Success(c, userDetailDTO.LoginVO)
}

// @Summary 注册
// @Description 注册
// @Tags UserAuth
// @Param form body RegisterForm true "注册"
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /register [post]
func (*UserAuth) Register(c *gin.Context) {
	r.Success(c, "注册")
}

// @Summary 退出登录
// @Description 退出登录
// @Tags UserAuth
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /logout [post]
func (*UserAuth) Logout(c *gin.Context) {
	val, ok := c.Get("uuid")
	if !ok {
		r.Error(c, r.ERROR_REQUEST_PARAM, nil)
		return
	}
	uuid := val.(string)
	session := sessions.Default(c)
	session.Delete(KEY_USER + uuid) //? FIXME: 删除后 redis 还会有一条记录?
	session.Save()
	utils.Redis.Del(KEY_USER + uuid) // 删除 redis 中缓存

	r.Success(c, "退出登录")
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
	r.Success(c, "发送邮箱验证码")
}

func convertUserDetailDTO(userAuth model.UserAuth, c *gin.Context) dto.UserDetailDTO {
	// 获取 IP 相关信息 FIXME: 好像无法读取到 ip 信息
	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSourceSimpleIdle(ipAddress)
	browser, os := "unknown", "unknown"

	if userAgent := utils.IP.GetUserAgent(c); userAgent != nil {
		browser = userAgent.Name + " " + userAgent.Version.String()
		os = userAgent.OS + " " + userAgent.OSVersion.String()
	}

	// 获取用户详细信息
	userInfo := dao.GetOne(&model.UserInfo{}, "id", userAuth.ID)
	// FIXME: 获取该用户对应的角色, 没有角色默认是 "test"
	roleLabels := roleDao.GetLabelsByUserInfoId(userInfo.ID)
	if len(roleLabels) == 0 {
		roleLabels = append(roleLabels, "test")
	}
	// 用户点赞 Set
	articleLikeSet := utils.Redis.SMembers(KEY_ARTICLE_USER_LIKE_SET + strconv.Itoa(userInfo.ID))
	commentLikeSet := utils.Redis.SMembers(KEY_COMMENT_USER_LIKE_SET + strconv.Itoa(userInfo.ID))

	return dto.UserDetailDTO{
		LoginVO: resp.LoginVO{
			ID:             userAuth.ID,
			UserInfoId:     userInfo.ID,
			Email:          userInfo.Email,
			LoginType:      userAuth.LoginType,
			Username:       userAuth.Username,
			Nickname:       userInfo.Nickname,
			Avatar:         userInfo.Avatar,
			Intro:          userInfo.Intro,
			Website:        userInfo.Website,
			IpAddress:      ipAddress,
			IpSource:       ipSource,
			LastLoginTime:  time.Now(),
			ArticleLikeSet: articleLikeSet,
			CommentLikeSet: commentLikeSet,
		},
		Password:   userAuth.Password,
		RoleLabels: roleLabels,
		IsDisable:  userInfo.IsDisable,
		Browser:    browser,
		OS:         os,
	}
}
