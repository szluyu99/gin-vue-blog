package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"gin-blog/internal/utils/jwt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Login 会用 g.Conf.JWT 签发 token, 测试里给一份最小配置
func withJWTConf(t *testing.T) {
	t.Helper()

	old := g.Conf
	g.Conf = &g.Config{}
	g.Conf.JWT.Secret = "test-secret"
	g.Conf.JWT.Issuer = "test-issuer"
	g.Conf.JWT.Expire = 24
	t.Cleanup(func() { g.Conf = old })
}

func TestAuthLogin(t *testing.T) {
	withJWTConf(t)

	env := newTestEnv(t)
	api := UserAuth{}
	env.engine.POST("/login", api.Login)

	user := createUser(t, env, "admin@qq.com", "管理员", "123456")
	// 登录后需要能读到用户点赞过的文章/评论
	env.rdb.SAdd(rctx, g.ARTICLE_USER_LIKE_SET+itoa(user.ID), "1", "2")
	env.rdb.SAdd(rctx, g.COMMENT_USER_LIKE_SET+itoa(user.ID), "3")
	// 之前被强制下线的标记, 登录后要清掉
	env.rdb.Set(rctx, g.OFFLINE_USER+itoa(user.ID), true, 0)

	resp := env.do(t, http.MethodPost, "/login", map[string]any{
		"username": "admin@qq.com", "password": "123456",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var vo LoginVO
	decodeData(t, resp.Data, &vo)
	assert.Equal(t, "管理员", vo.Nickname)
	assert.NotEmpty(t, vo.Token)
	assert.ElementsMatch(t, []string{"1", "2"}, vo.ArticleLikeSet)
	assert.Equal(t, []string{"3"}, vo.CommentLikeSet)

	// token 可以被解析回用户 ID
	claims, err := jwt.ParseToken("test-secret", vo.Token)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, claims.UserId)

	assert.False(t, env.mr.Exists(g.OFFLINE_USER+itoa(user.ID)))

	// 登录信息(IP + 时间)被记录
	var updated model.UserAuth
	assert.Nil(t, env.db.First(&updated, user.ID).Error)
	assert.NotNil(t, updated.LastLoginTime)
}

func TestAuthLoginFailed(t *testing.T) {
	withJWTConf(t)

	env := newTestEnv(t)
	env.engine.POST("/login", (&UserAuth{}).Login)

	createUser(t, env, "admin@qq.com", "管理员", "123456")

	// 缺少字段
	resp := env.do(t, http.MethodPost, "/login", map[string]any{"username": "admin@qq.com"})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 用户不存在与密码错误返回同一个错误码, 否则可以靠错误码枚举用户名
	resp = env.do(t, http.MethodPost, "/login", map[string]any{
		"username": "nobody@qq.com", "password": "123456",
	})
	assert.Equal(t, g.ErrLoginFail.Code(), resp.Code)

	resp = env.do(t, http.MethodPost, "/login", map[string]any{
		"username": "admin@qq.com", "password": "wrong",
	})
	assert.Equal(t, g.ErrLoginFail.Code(), resp.Code)
}

// 连续失败达到上限后直接拒绝, 即使密码正确
func TestAuthLoginLockedAfterTooManyFails(t *testing.T) {
	withJWTConf(t)

	env := newTestEnv(t)
	env.engine.POST("/login", (&UserAuth{}).Login)

	createUser(t, env, "admin@qq.com", "管理员", "123456")
	wrong := map[string]any{"username": "admin@qq.com", "password": "wrong"}

	for i := 0; i < loginMaxFail; i++ {
		resp := env.do(t, http.MethodPost, "/login", wrong)
		assert.Equal(t, g.ErrLoginFail.Code(), resp.Code, "第 %d 次失败", i+1)
	}

	// 密码正确也被拦住
	resp := env.do(t, http.MethodPost, "/login", map[string]any{
		"username": "admin@qq.com", "password": "123456",
	})
	assert.Equal(t, g.ErrLoginLocked.Code(), resp.Code)

	// 换一个 IP 不受影响: 计数是按 用户名+IP 记的
	resp = env.doWithHeader(t, http.MethodPost, "/login", map[string]any{
		"username": "admin@qq.com", "password": "123456",
	}, map[string]string{"X-Real-IP": "10.0.0.9"})
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 冷却时间过后可以再登录
	env.mr.FastForward(loginFailWindow + time.Second)
	resp = env.do(t, http.MethodPost, "/login", map[string]any{
		"username": "admin@qq.com", "password": "123456",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)
}

// 登录成功要清掉之前的失败计数, 否则攒够 5 次后正常用户也会被锁
func TestAuthLoginSuccessResetsFailCount(t *testing.T) {
	withJWTConf(t)

	env := newTestEnv(t)
	env.engine.POST("/login", (&UserAuth{}).Login)

	createUser(t, env, "admin@qq.com", "管理员", "123456")

	for i := 0; i < loginMaxFail-1; i++ {
		env.do(t, http.MethodPost, "/login", map[string]any{
			"username": "admin@qq.com", "password": "wrong",
		})
	}

	resp := env.do(t, http.MethodPost, "/login", map[string]any{
		"username": "admin@qq.com", "password": "123456",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 计数被清掉, 之后再失败一次也不会触发锁定
	resp = env.do(t, http.MethodPost, "/login", map[string]any{
		"username": "admin@qq.com", "password": "wrong",
	})
	assert.Equal(t, g.ErrLoginFail.Code(), resp.Code)
}

func TestAuthLogout(t *testing.T) {
	withJWTConf(t)

	env := newTestEnv(t)
	api := UserAuth{}
	env.engine.POST("/login", api.Login)
	env.engine.POST("/logout", api.Logout)

	user := createUser(t, env, "admin@qq.com", "管理员", "123456")
	env.rdb.Set(rctx, g.ONLINE_USER+itoa(user.ID), true, 0)

	// Logout 会先把 context 里的用户置空, 只认 session, 所以要走一次真实登录拿 cookie
	loginReq := httptest.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"username":"admin@qq.com","password":"123456"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.engine.ServeHTTP(w, loginReq)
	assert.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	w = httptest.NewRecorder()
	env.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// 在线标记被清掉
	assert.False(t, env.mr.Exists(g.ONLINE_USER+itoa(user.ID)))

	// 没有 session 时调用也返回成功
	resp := env.do(t, http.MethodPost, "/logout", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
}

func TestAuthRegisterInvalid(t *testing.T) {
	env := newTestEnv(t)
	env.engine.POST("/register", (&UserAuth{}).Register)

	createUser(t, env, "admin@qq.com", "管理员", "123456")

	// 密码长度不满足 min=4
	resp := env.do(t, http.MethodPost, "/register", map[string]any{
		"email": "new@qq.com", "password": "1",
	})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 邮箱已注册: 不会走到发邮件
	resp = env.do(t, http.MethodPost, "/register", map[string]any{
		"email": "admin@qq.com", "password": "123456",
	})
	assert.Equal(t, g.ErrUserExist.Code(), resp.Code)
}

// 邮箱验证链接: info 由注册时写入 Redis, 命中后才真正创建用户
func TestAuthVerifyCode(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/verify", (&UserAuth{}).VerifyCode)

	// 游客角色, CreateNewUser 里固定关联 role_id = 2
	env.db.Create(&model.Role{Model: model.Model{ID: 2}, Name: "游客"})

	info := utils.GenEmailVerificationInfo("new@qq.com", "123456")
	assert.Nil(t, SetMailInfo(env.rdb, info, time.Minute))

	w := httptest.NewRecorder()
	env.engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/verify?info="+info, nil))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "注册成功")

	var auth model.UserAuth
	assert.Nil(t, env.db.Where("username", "new@qq.com").First(&auth).Error)
	assert.True(t, utils.BcryptCheck("123456", auth.Password))

	// info 用过即失效, 同一个链接不能重复注册
	w = httptest.NewRecorder()
	env.engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/verify?info="+info, nil))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "注册失败")
}

func TestAuthVerifyCodeWithoutInfo(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/verify", (&UserAuth{}).VerifyCode)

	w := httptest.NewRecorder()
	env.engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/verify", nil))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "注册失败")
}
