package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 创建一个带 UserInfo 的用户, 返回 UserAuth
func createUser(t *testing.T, env *testEnv, username, nickname, password string) *model.UserAuth {
	t.Helper()

	info := model.UserInfo{Nickname: nickname, Email: username + "@qq.com"}
	assert.Nil(t, env.db.Create(&info).Error)

	hash, err := utils.BcryptHash(password)
	assert.Nil(t, err)

	auth := model.UserAuth{Username: username, Password: hash, UserInfoId: info.ID}
	assert.Nil(t, env.db.Create(&auth).Error)

	auth.UserInfo = &info
	return &auth
}

// 用户列表: 分页与按用户名/昵称过滤
func TestUserGetList(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/user/list", (&User{}).GetList)

	createUser(t, env, "admin", "管理员", "123456")
	createUser(t, env, "guest", "访客", "123456")

	resp := env.do(t, http.MethodGet, "/user/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var page PageResult[model.UserAuth]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(2), page.Total)

	// 按用户名过滤
	resp = env.do(t, http.MethodGet, "/user/list?username=guest", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "guest", page.List[0].Username)

	// 密码不应出现在响应中
	assert.Empty(t, page.List[0].Password)
}

// 修改当前用户密码: 需要旧密码验证
func TestUserUpdateCurrentPassword(t *testing.T) {
	env := newTestEnv(t)
	env.engine.PUT("/user/current/password", (&User{}).UpdateCurrentPassword)

	user := createUser(t, env, "admin", "管理员", "123456")
	env.user = user // 模拟该用户已登录

	// 旧密码不正确
	resp := env.do(t, http.MethodPut, "/user/current/password", map[string]any{
		"old_password": "wrong0", "new_password": "654321",
	})
	assert.Equal(t, g.ErrOldPassword.Code(), resp.Code)

	// 新密码太短(min=4)
	resp = env.do(t, http.MethodPut, "/user/current/password", map[string]any{
		"old_password": "123456", "new_password": "1",
	})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 正常修改, 数据库中存的是新密码的 hash
	resp = env.do(t, http.MethodPut, "/user/current/password", map[string]any{
		"old_password": "123456", "new_password": "654321",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var got model.UserAuth
	assert.Nil(t, env.db.First(&got, user.ID).Error)
	assert.True(t, utils.BcryptCheck("654321", got.Password))
	assert.False(t, utils.BcryptCheck("123456", got.Password))
}

// 修改用户禁用状态
func TestUserUpdateDisable(t *testing.T) {
	env := newTestEnv(t)
	env.engine.PUT("/user/disable", (&User{}).UpdateDisable)

	user := createUser(t, env, "guest", "访客", "123456")

	resp := env.do(t, http.MethodPut, "/user/disable", map[string]any{
		"id": user.ID, "is_disable": true,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var got model.UserAuth
	assert.Nil(t, env.db.First(&got, user.ID).Error)
	assert.True(t, got.IsDisable)
}

// 获取当前用户信息: 点赞集合来自 Redis
func TestUserGetInfo(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/user/info", (&User{}).GetInfo)

	user := createUser(t, env, "admin", "管理员", "123456")
	env.user = user

	env.rdb.SAdd(rctx, g.ARTICLE_USER_LIKE_SET+itoa(user.ID), 3)
	env.rdb.SAdd(rctx, g.COMMENT_USER_LIKE_SET+itoa(user.ID), 5)

	resp := env.do(t, http.MethodGet, "/user/info", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var info model.UserInfoVO
	decodeData(t, resp.Data, &info)
	assert.Equal(t, "管理员", info.Nickname)
	assert.Equal(t, []string{"3"}, info.ArticleLikeSet)
	assert.Equal(t, []string{"5"}, info.CommentLikeSet)
}

// 在线用户与强制下线
func TestUserOnlineAndForceOffline(t *testing.T) {
	env := newTestEnv(t)
	api := User{}
	env.engine.GET("/user/online", api.GetOnlineList)
	env.engine.POST("/user/offline/:id", api.ForceOffline)

	admin := createUser(t, env, "admin", "管理员", "123456")
	guest := createUser(t, env, "guest", "访客", "123456")
	env.user = admin

	// 在线列表来自 Redis 中的 online_user:* key
	env.rdb.Set(rctx, g.ONLINE_USER+itoa(guest.ID), guest, 0)
	resp := env.do(t, http.MethodGet, "/user/online", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 不能强制下线自己
	resp = env.do(t, http.MethodPost, "/user/offline/"+itoa(admin.ID), nil)
	assert.Equal(t, g.ErrForceOfflineSelf.Code(), resp.Code)

	// 下线其他用户: 移除在线标记, 写入下线标记
	resp = env.do(t, http.MethodPost, "/user/offline/"+itoa(guest.ID), nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.False(t, env.mr.Exists(g.ONLINE_USER+itoa(guest.ID)))
	assert.True(t, env.mr.Exists(g.OFFLINE_USER+itoa(guest.ID)))

	// id 不是数字
	resp = env.do(t, http.MethodPost, "/user/offline/abc", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}

// 修改当前用户信息: 用户 ID 取自登录态, 不信前端
func TestUserUpdateCurrent(t *testing.T) {
	env := newTestEnv(t)
	env.engine.PUT("/user/current", (&User{}).UpdateCurrent)

	info := model.UserInfo{Nickname: "老昵称"}
	assert.Nil(t, env.db.Create(&info).Error)
	user := model.UserAuth{Username: "tester", Password: "x", UserInfoId: info.ID}
	assert.Nil(t, env.db.Create(&user).Error)
	env.user = &model.UserAuth{Model: model.Model{ID: user.ID}, Username: "tester", UserInfoId: info.ID}

	resp := env.do(t, http.MethodPut, "/user/current", map[string]any{
		"nickname": "新昵称",
		"avatar":   "avatar.png",
		"intro":    "简介",
		"website":  "https://test.com",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	got, err := model.GetUserInfoById(env.db, info.ID)
	assert.Nil(t, err)
	assert.Equal(t, "新昵称", got.Nickname)
	assert.Equal(t, "https://test.com", got.Website)

	// 昵称必填
	resp = env.do(t, http.MethodPut, "/user/current", map[string]any{"intro": "只改简介"})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 未登录
	env.user = nil
	resp = env.do(t, http.MethodPut, "/user/current", map[string]any{"nickname": "匿名"})
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)
}

// 后台修改用户昵称与角色: 角色是整体替换
func TestUserUpdate(t *testing.T) {
	env := newTestEnv(t)
	env.engine.PUT("/user", (&User{}).Update)

	role1 := model.Role{Name: "role1", Label: "角色1"}
	role2 := model.Role{Name: "role2", Label: "角色2"}
	assert.Nil(t, env.db.Create(&role1).Error)
	assert.Nil(t, env.db.Create(&role2).Error)

	user := model.UserAuth{Username: "tester", Password: "x", UserInfo: &model.UserInfo{Nickname: "老昵称"}}
	assert.Nil(t, env.db.Create(&user).Error)
	assert.Nil(t, env.db.Create(&model.UserAuthRole{UserAuthId: user.ID, RoleId: role1.ID}).Error)

	resp := env.do(t, http.MethodPut, "/user", map[string]any{
		"id":       user.ID,
		"nickname": "新昵称",
		"role_ids": []int{role2.ID},
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	info, err := model.GetUserInfoById(env.db, user.UserInfoId)
	assert.Nil(t, err)
	assert.Equal(t, "新昵称", info.Nickname)

	roleIds, err := model.GetRoleIdsByUserId(env.db, user.ID)
	assert.Nil(t, err)
	assert.Equal(t, []int{role2.ID}, roleIds)

	// 昵称必填
	resp = env.do(t, http.MethodPut, "/user", map[string]any{"id": user.ID})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 用户不存在
	resp = env.do(t, http.MethodPut, "/user", map[string]any{"id": 999, "nickname": "x"})
	assert.Equal(t, g.ErrDbOp.Code(), resp.Code)
}
