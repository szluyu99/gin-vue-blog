package service

import (
	"fmt"
	"gin-blog/config"
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/dto"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type User struct{}

// 登录
func (*User) Login(c *gin.Context, username, password string) (loginVo resp.LoginVO, code int) {
	// 检查用户是否存在
	userAuth := dao.GetOne(model.UserAuth{}, "username", username)
	if userAuth.ID == 0 {
		return loginVo, r.ERROR_USER_NOT_EXIST
	}
	// 检查密码是否正确
	if !utils.Encryptor.BcryptCheck(password, userAuth.Password) {
		return loginVo, r.ERROR_PASSWORD_WRONG
	}
	// 获取用户详细信息 DTO
	userDetailDTO := convertUserDetailDTO(userAuth, c)

	// 登录信息正确, 生成 Token
	// TODO: 目前只给用户设定一个角色, 获取第一个值就行, 后期优化: 给用户设置多个角色
	// UUID 生成方法: ip + 浏览器信息 + 操作系统信息
	uuid := utils.Encryptor.MD5(userDetailDTO.IpAddress + userDetailDTO.Browser + userDetailDTO.OS)
	token, err := utils.GetJWT().GenToken(userAuth.ID, userDetailDTO.RoleLabels[0], uuid)
	if err != nil {
		utils.Logger.Info("登录时生成 Token 错误: ", zap.Error(err))
		return loginVo, r.ERROR_TOKEN_CREATE
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
	// ! session 中只能存储字符串
	sessionInfoStr := utils.Json.Marshal(dto.SessionInfo{UserDetailDTO: userDetailDTO})
	session.Set(KEY_USER+uuid, sessionInfoStr) // ! 确实设置到 reids 中, 但是获取不到
	utils.Redis.Set(KEY_USER+uuid, sessionInfoStr, time.Duration(config.Cfg.Session.MaxAge)*time.Second)
	// fmt.Println("login: ", KEY_USER+uuid)
	session.Save()

	return userDetailDTO.LoginVO, r.OK
}

// 退出登录
func (*User) Logout(c *gin.Context) {
	uuid := utils.GetFromContext[string](c, "uuid")
	session := sessions.Default(c)
	session.Delete(KEY_USER + uuid) //? FIXME: 删除后 redis 还会有一条记录?
	session.Save()
	utils.Redis.Del(KEY_USER + uuid) // 删除 redis 中缓存
}

// 用户注册
func (*User) Register(req req.Register) (code int) {
	// 检查验证码是否正确
	if req.Code != utils.Redis.GetVal(KEY_CODE+req.Username) {
		return r.ERROR_VERIFICATION_CODE
	}

	// 检查用户名已存在, 则该账号已经注册过
	if exist := checkUserExistByName(req.Username); exist {
		return r.ERROR_USER_NAME_USED
	}

	userInfo := &model.UserInfo{
		Email:    req.Username,
		Nickname: "用户" + req.Username,
		// blogConfig 中配置的默认用户头像
		Avatar: blogInfoService.GetBlogConfig().UserAvatar,
	}
	dao.Create(&userInfo)
	// 设置用户默认角色
	dao.Create(&model.UserRole{
		UserId: userInfo.ID,
		RoleId: 3, // 默认角色是 "测试"
	})
	dao.Create(&model.UserAuth{
		UserInfoId:    userInfo.ID,
		Username:      req.Username,
		Password:      utils.Encryptor.BcryptHash(req.Password),
		LoginType:     1,
		LastLoginTime: time.Now(), // 注册时会更新 "上次登录时间"
	})
	return r.OK
}

// 更新用户邮箱信息
func (*User) UpdateEmail(userInfoId int, req req.UpdateEmail) (code int) {
	// 检查验证码是否正确
	if req.Code != utils.Redis.GetVal(KEY_CODE+req.Email) {
		return r.ERROR_VERIFICATION_CODE
	}

	// TODO: 更新 user表 和 user_auth表 的邮箱信息

	return r.OK
}

// 查询当前在线用户
// TODO: 分页 + 条件搜索
func (*User) GetOnlineList(req req.PageQuery) resp.PageResult[[]resp.UserOnline] {
	onlineList := make([]resp.UserOnline, 0)

	keys := utils.Redis.Keys(KEY_USER + "*")
	for _, key := range keys {
		var sessionInfo dto.SessionInfo
		utils.Json.Unmarshal(utils.Redis.GetVal(key), &sessionInfo)

		// 查询关键字不为空, 且不满足查询条件
		if req.Keyword != "" && !strings.Contains(sessionInfo.Nickname, req.Keyword) {
			continue
		}

		onlineUser := utils.CopyProperties[resp.UserOnline](sessionInfo)
		onlineUser.UserIndoId = sessionInfo.UserInfoId // *
		onlineList = append(onlineList, onlineUser)
	}

	// 根据上次登录时间进行排序
	sort.Slice(onlineList, func(i, j int) bool {
		return onlineList[i].LastLoginTime.Unix() > onlineList[j].LastLoginTime.Unix()
	})
	return resp.PageResult[[]resp.UserOnline]{
		Total: int64(len(keys)),
		List:  onlineList,
	}
}

// 强制用户离线
func (*User) ForceOffline(req req.ForceOfflineUser) (code int) {
	// TODO: 不能让自己离线? 目前设定是可以的...
	// if userInfoId == req.UserIndoId {
	// }

	uuid := utils.Encryptor.MD5(req.IpAddress + req.Browser + req.OS)
	var sessionInfo dto.SessionInfo
	utils.Json.Unmarshal(utils.Redis.GetVal(KEY_USER+uuid), &sessionInfo)
	sessionInfo.IsOffline = 1 // *
	utils.Redis.Del(KEY_USER + uuid)
	// ? 这里设置强制离线后 redis 中存储的 delete:xxx 时间和 Token 过期时间一致
	utils.Redis.Set(KEY_DELETE+uuid, utils.Json.Marshal(sessionInfo), time.Duration(config.Cfg.JWT.Expire)*time.Hour)
	return r.OK
}

// TODO: 用户区域分布 GetUserAreas, statiscalUserArea

func (*User) UpdateCurrent(current req.UpdateCurrentUser) (code int) {
	user := utils.CopyProperties[model.UserInfo](current)
	dao.Update(&user, "nickname", "intro", "website", "avatar", "email")
	return r.OK
}

// TODO: 优化
func (*User) GetInfo(id int) resp.UserInfoVO {
	var userInfo model.UserInfo
	dao.GetOne(&userInfo, "id", id)

	data := utils.CopyProperties[resp.UserInfoVO](userInfo)
	data.ArticleLikeSet = utils.Redis.SMembers(KEY_ARTICLE_USER_LIKE_SET + strconv.Itoa(id))
	data.CommentLikeSet = utils.Redis.SMembers(KEY_COMMENT_USER_LIKE_SET + strconv.Itoa(id))
	return data
}

func (*User) GetList(req req.GetUsers) resp.PageResult[[]resp.UserVO] {
	count := userDao.GetCount(req)
	list := userDao.GetList(req)
	return resp.PageResult[[]resp.UserVO]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    count,
		List:     list,
	}
}

func (*User) Update(req req.UpdateUser) int {
	userInfo := model.UserInfo{
		Universal: model.Universal{ID: req.UserInfoId},
		Nickname:  req.Nickname,
	}
	dao.Update(&userInfo)
	// 清空 user_role 关系
	dao.Delete(model.UserRole{}, "user_id = ?", req.UserInfoId)
	// 要更新的 user_role 列表
	var userRoles []model.UserRole
	for _, id := range req.RoleIds {
		userRoles = append(userRoles, model.UserRole{
			RoleId: id,
			UserId: req.UserInfoId,
		})
	}
	dao.Create(&userRoles)
	return r.OK
}

// 修改普通用户密码, 不需要旧密码
func (*User) UpdatePassword(req req.UpdatePassword) int {
	// 用户名不存在
	if exist := checkUserExistByName(req.Username); !exist {
		return r.ERROR_USER_NOT_EXIST
	}

	m := map[string]any{"password": utils.Encryptor.BcryptHash(req.Password)}
	dao.UpdatesMap(&model.UserAuth{}, m, "username = ?", req.Username)
	return r.OK
}

// 修改管理员密码, 需要旧密码验证
func (*User) UpdateCurrentPassword(req req.UpdateAdminPassword, id int) int {
	user := dao.GetOne(model.UserAuth{}, "id", id)
	if !user.IsEmpty() && utils.Encryptor.BcryptCheck(req.OldPassword, user.Password) {
		user.Password = utils.Encryptor.BcryptHash(req.NewPassword)
		dao.Update(&user, "password")
		return r.OK
	} else {
		return r.ERROR_OLD_PASSWORD
	}
}

// 更新用户禁用
func (*User) UpdateDisable(id, isDisable int) {
	dao.UpdatesMap(&model.UserInfo{}, map[string]any{"is_disable": isDisable}, "id", id)
}

// 发送验证码
func (*User) SendCode(email string) (code int) {
	// 已经发送验证码且未过期
	if utils.Redis.GetVal(KEY_CODE+email) != "" {
		return r.ERROR_EMAIL_HAS_SEND
	}

	expireTime := config.Cfg.Captcha.ExpireTime
	validateCode := utils.Encryptor.ValidateCode()
	content := fmt.Sprintf(`
		<div style="text-align:center">
			<div>你好！欢迎访问阵、雨的个人博客！</div>
			<div style="padding: 8px 40px 8px 50px;">
            	<p>
					您本次的验证码为
					<p style="font-size:75px;font-weight:blod;"> %s </p>
					为了保证账号安全，验证码有效期为 %d 分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用~
				</p>
       	 	</div>
			<div>
            	<p>发送邮件专用邮箱，请勿回复。</p>
        	</div>
		</div>
	`, validateCode, expireTime)

	if err := utils.Email(email, "博客注册验证码", content); err != nil {
		return r.ERROR_EMAIL_SEND
	}

	// 将验证码存储到 Redis 中
	utils.Redis.Set(KEY_CODE+email, validateCode, time.Duration(expireTime)*time.Minute)
	return r.OK
}

// 检查用户名是否存在
func checkUserExistByName(username string) bool {
	existUser := dao.GetOne(model.UserAuth{}, "username = ?", username)
	return existUser.ID != 0
}

// 转化 UserDetailDTO
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
