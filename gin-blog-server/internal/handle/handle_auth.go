package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"gin-blog/internal/utils/jwt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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
	Username string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
	
}

type LoginVO struct {
	model.UserInfo

	// 点赞 Set: 用于记录用户点赞过的文章, 评论
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
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
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	userAuth, err := model.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 检查密码是否正确
	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		ReturnError(c, g.ErrPassword, nil)
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
			ReturnError(c, g.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	roleIds, err := model.GetRoleIdsByUserId(db, userAuth.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	articleLikeSet, err := rdb.SMembers(rctx, g.ARTICLE_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}
	commentLikeSet, err := rdb.SMembers(rctx, g.COMMENT_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	// 登录信息正确, 生成 Token

	// UUID 生成方法: ip + 浏览器信息 + 操作系统信息
	// uuid := utils.MD5(ipAddress + browser + os)
	conf := g.Conf.JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), userAuth.ID, roleIds)
	if err != nil {
		ReturnError(c, g.ErrTokenCreate, err)
		return
	}

	// 更新用户验证信息: ip 信息 + 上次登录时间
	err = model.UpdateUserLoginInfo(db, userAuth.ID, ipAddress, ipSource)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	slog.Info("用户登录成功: " + userAuth.Username)

	session := sessions.Default(c)
	session.Set(g.CTX_USER_AUTH, userAuth.ID)
	session.Save()

	// 删除 Redis 中的离线状态
	offlineKey := g.OFFLINE_USER + strconv.Itoa(userAuth.ID)
	rdb.Del(rctx, offlineKey).Result()

	ReturnSuccess(c, LoginVO{
		UserInfo: *userInfo,

		ArticleLikeSet: articleLikeSet,
		CommentLikeSet: commentLikeSet,
		Token:          token,
	})
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
	
	// 已经退出登录
	auth, _ := CurrentUserAuth(c)
	if auth == nil {
		ReturnSuccess(c, nil)
		return
	}

	session := sessions.Default(c)
	session.Delete(g.CTX_USER_AUTH)
	session.Save()
	
	// 删除 Redis 中的在线状态
	rdb := GetRDB(c)
	onlineKey := g.ONLINE_USER + strconv.Itoa(auth.ID)
	rdb.Del(rctx, onlineKey)

	ReturnSuccess(c, nil)
}

// 完成注册功能
// 首先检查用户名是否存在，避免重复注册；其次把用户输入的信息加密保存在redis中，等待验证
// 在以下情况下会出错：1. 用户邮箱已经注册过 2.用户邮箱无效等原因导致的发送邮件失败
func (*UserAuth) Register(c *gin.Context) {
	var regreq RegisterReq
	if err := c.ShouldBindJSON(&regreq); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	// 格式化用户名
	regreq.Username = utils.Format(regreq.Username)

	// 检查用户名是否存在，避免重复注册
	auth,err := model.GetUserAuthInfoByName(GetDB(c),regreq.Username)
	if err != nil {
		var flag bool = false
		if errors.Is(err,gorm.ErrRecordNotFound) {
			flag = true
		}
		if !flag{
			ReturnError(c,g.ErrDbOp,err)
			return
		}
	}

	if auth != nil {
		ReturnError(c,g.ErrUserExist,err)
		return
	}
	

	// 通过邮箱验证
	info := utils.GenEmailVerificationInfo(regreq.Username,regreq.Password)
	SetMailInfo(GetRDB(c),info,15*time.Minute) // 15分钟过期
	EmailData := utils.GetEmailData(regreq.Username,info)
	err = utils.SendEmail(regreq.Username,EmailData)
	if err != nil {
		ReturnError(c,g.ErrSendEmail,err)
		return
	}

	ReturnSuccess(c,nil)
}

// 邮箱验证
// 当用户点击邮箱中的链接时，会携带info（加密后的帐号密码）向这个接口发送请求。
// Verify会检查info是否存在redis中，若存在则认证成功，完成注册
// 会在以下方面出错： 1. 发送信息中没有info 2. info不存在redis中(已过期) 3. 创造新用户失败
func (*UserAuth) VerifyCode(c *gin.Context) {
    var code string
    if code = c.Query("info"); code == "" {
        returnErrorPage(c)
        return
    }

    // 验证是否有code在数据库中
    ifExist, err := GetMailInfo(GetRDB(c), code)
    if err != nil {
        returnErrorPage(c)
        return
    }
    if !ifExist {
        returnErrorPage(c)
        return
    }

    DeleteMailInfo(GetRDB(c), code)

    username, password, err := utils.ParseEmailVerificationInfo(code)
    if err != nil {
        returnErrorPage(c)
        return
    }

    // 注册用户
      _,_,_,err = model.CreateNewUser(GetDB(c), username, password)
    if err != nil {
        returnErrorPage(c)
        return
    }

    // 注册成功，返回成功页面
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册成功</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #5cb85c;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册成功</h1>
                <p>恭喜您，注册成功！</p>
            </div>
        </body>
        </html>
    `))
}

func returnErrorPage(c *gin.Context) {
    c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册失败</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #d9534f;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册失败</h1>
                <p>请重试。</p>
            </div>
        </body>
        </html>
    `))
}