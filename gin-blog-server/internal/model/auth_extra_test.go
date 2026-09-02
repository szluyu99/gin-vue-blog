package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func seedMenu(t *testing.T, db *gorm.DB, name, path string, parentId int) *Menu {
	t.Helper()
	menu := Menu{Name: name, Path: path, ParentId: parentId, Component: "Layout"}
	assert.Nil(t, SaveOrUpdateMenu(db, &menu))
	assert.NotZero(t, menu.ID)
	return &menu
}

func TestMenuCRUD(t *testing.T) {
	db := newModelDB(t)

	parent := seedMenu(t, db, "文章管理", "/article", 0)
	child := seedMenu(t, db, "文章列表", "/article/list", parent.ID)

	// id > 0 走更新
	parent.Name = "内容管理"
	assert.Nil(t, SaveOrUpdateMenu(db, parent))
	got, err := GetMenuById(db, parent.ID)
	assert.Nil(t, err)
	assert.Equal(t, "内容管理", got.Name)

	all, err := GetAllMenuList(db)
	assert.Nil(t, err)
	assert.Len(t, all, 2)

	list, total, err := GetMenuList(db, "文章")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, child.ID, list[0].ID)

	hasChild, err := CheckMenuHasChild(db, parent.ID)
	assert.Nil(t, err)
	assert.True(t, hasChild)

	inUse, err := CheckMenuInUse(db, child.ID)
	assert.Nil(t, err)
	assert.False(t, inUse)

	// 挂到角色下之后就算被使用
	role := Role{Name: "admin", Label: "管理员"}
	assert.Nil(t, db.Create(&role).Error)
	assert.Nil(t, db.Create(&RoleMenu{RoleId: role.ID, MenuId: child.ID}).Error)
	inUse, err = CheckMenuInUse(db, child.ID)
	assert.Nil(t, err)
	assert.True(t, inUse)

	ids, err := GetMenuIdsByRoleId(db, role.ID)
	assert.Nil(t, err)
	assert.Equal(t, []int{child.ID}, ids)

	assert.Nil(t, DeleteMenu(db, child.ID))
	_, total, _ = GetMenuList(db, "")
	assert.Equal(t, int64(1), total)
}

// 用户菜单来自其所有角色的并集, 多个角色共有的菜单只出现一次
func TestGetMenuListByUserIdDedup(t *testing.T) {
	db := newModelDB(t)
	home := seedMenu(t, db, "首页", "/home", 0)
	article := seedMenu(t, db, "文章", "/article", 0)

	role1 := Role{Name: "role1", Label: "角色1"}
	role2 := Role{Name: "role2", Label: "角色2"}
	assert.Nil(t, db.Create(&role1).Error)
	assert.Nil(t, db.Create(&role2).Error)
	// 两个角色都有首页, role2 额外有文章
	assert.Nil(t, db.Create(&RoleMenu{RoleId: role1.ID, MenuId: home.ID}).Error)
	assert.Nil(t, db.Create(&RoleMenu{RoleId: role2.ID, MenuId: home.ID}).Error)
	assert.Nil(t, db.Create(&RoleMenu{RoleId: role2.ID, MenuId: article.ID}).Error)

	user := UserAuth{Username: "test", Password: "x", UserInfo: &UserInfo{Nickname: "测试"}}
	assert.Nil(t, db.Create(&user).Error)
	assert.Nil(t, db.Create(&UserAuthRole{UserAuthId: user.ID, RoleId: role1.ID}).Error)
	assert.Nil(t, db.Create(&UserAuthRole{UserAuthId: user.ID, RoleId: role2.ID}).Error)

	menus, err := GetMenuListByUserId(db, user.ID)
	assert.Nil(t, err)
	assert.Len(t, menus, 2)
}

func TestResourceCRUD(t *testing.T) {
	db := newModelDB(t)

	assert.Nil(t, SaveOrUpdateResource(db, 0, 0, "文章模块", "", ""))
	list, err := GetResourceList(db, "")
	assert.Nil(t, err)
	assert.Len(t, list, 1)

	parent := list[0]
	assert.Nil(t, SaveOrUpdateResource(db, 0, parent.ID, "文章列表", "/article/list", "GET"))
	assert.Nil(t, SaveOrUpdateResource(db, 0, parent.ID, "文章保存", "/article", "POST"))

	list, err = GetResourceList(db, "文章列")
	assert.Nil(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "/article/list", list[0].Url)

	// id > 0 走更新
	assert.Nil(t, SaveOrUpdateResource(db, list[0].ID, parent.ID, "文章分页列表", "/article/list", "GET"))
	got, err := GetResourceById(db, list[0].ID)
	assert.Nil(t, err)
	assert.Equal(t, "文章分页列表", got.Name)

	hasChild, err := CheckResourceHasChild(db, parent.ID)
	assert.Nil(t, err)
	assert.True(t, hasChild)

	inUse, err := CheckResourceInUse(db, got.ID)
	assert.Nil(t, err)
	assert.False(t, inUse)
}

func TestRoleListAndOption(t *testing.T) {
	db := newModelDB(t)
	assert.Nil(t, SaveRole(db, "admin", "管理员"))
	assert.Nil(t, SaveRole(db, "user", "普通用户"))

	options, err := GetRoleOption(db)
	assert.Nil(t, err)
	assert.Len(t, options, 2)

	list, total, err := GetRoleList(db, 1, 10, "adm")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "admin", list[0].Name)

	// SaveOrUpdateRole: id > 0 时更新
	assert.Nil(t, SaveOrUpdateRole(db, list[0].ID, "admin", "超级管理员", true))
	list, _, _ = GetRoleList(db, 1, 10, "adm")
	assert.Equal(t, "超级管理员", list[0].Label)
	assert.True(t, list[0].IsDisable)
}

// 更新角色: 资源与菜单关联整体替换
func TestUpdateRoleWithResourceAndMenu(t *testing.T) {
	db := newModelDB(t)
	res1, err := AddResource(db, "api_1", "/api_1", "GET", false)
	assert.Nil(t, err)
	res2, err := AddResource(db, "api_2", "/api_2", "POST", false)
	assert.Nil(t, err)
	menu := seedMenu(t, db, "首页", "/home", 0)

	role, err := AddRoleWithResources(db, "admin", "管理员", []int{res1.ID})
	assert.Nil(t, err)

	assert.Nil(t, UpdateRole(db, role.ID, "admin", "超管", false, []int{res2.ID}, []int{menu.ID}))

	resourceIds, err := GetResourceIdsByRoleId(db, role.ID)
	assert.Nil(t, err)
	assert.Equal(t, []int{res2.ID}, resourceIds)

	menuIds, err := GetMenuIdsByRoleId(db, role.ID)
	assert.Nil(t, err)
	assert.Equal(t, []int{menu.ID}, menuIds)

	// 删除角色会一并清掉关联
	assert.Nil(t, DeleteRoles(db, []int{role.ID}))
	resourceIds, _ = GetResourceIdsByRoleId(db, role.ID)
	assert.Empty(t, resourceIds)
	menuIds, _ = GetMenuIdsByRoleId(db, role.ID)
	assert.Empty(t, menuIds)
}

func TestCreateNewUser(t *testing.T) {
	db := newModelDB(t)
	// 注册的用户默认挂到 id 为 2 的角色上, 先把角色造出来
	assert.Nil(t, SaveRole(db, "admin", "管理员"))
	assert.Nil(t, SaveRole(db, "user", "普通用户"))

	auth, info, userRole, err := CreateNewUser(db, "newbie", "123456")
	assert.Nil(t, err)
	assert.NotZero(t, auth.ID)
	assert.NotZero(t, info.ID)
	assert.Equal(t, auth.ID, userRole.UserAuthId)

	got, err := GetUserAuthInfoByName(db, "newbie")
	assert.Nil(t, err)
	assert.Equal(t, auth.ID, got.ID)

	gotInfo, err := GetUserInfoById(db, info.ID)
	assert.Nil(t, err)
	assert.Equal(t, info.Nickname, gotInfo.Nickname)
}

func TestUserListAndUpdate(t *testing.T) {
	db := newModelDB(t)
	assert.Nil(t, SaveRole(db, "admin", "管理员"))
	assert.Nil(t, SaveRole(db, "user", "普通用户"))

	u1 := UserAuth{Username: "zhangsan", Password: "x", LoginType: 1, UserInfo: &UserInfo{Nickname: "张三"}}
	u2 := UserAuth{Username: "lisi", Password: "x", LoginType: 2, UserInfo: &UserInfo{Nickname: "李四"}}
	assert.Nil(t, db.Create(&u1).Error)
	assert.Nil(t, db.Create(&u2).Error)

	// 按用户名 / 昵称 / 登录方式过滤
	_, total, err := GetUserList(db, 1, 10, 0, "", "zhang")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)

	_, total, err = GetUserList(db, 1, 10, 0, "李", "")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)

	_, total, err = GetUserList(db, 1, 10, 2, "", "")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)

	assert.Nil(t, UpdateUserInfo(db, u1.UserInfoId, "张三丰", "avatar", "简介", "https://test.com"))
	info, err := GetUserInfoById(db, u1.UserInfoId)
	assert.Nil(t, err)
	assert.Equal(t, "张三丰", info.Nickname)
	assert.Equal(t, "https://test.com", info.Website)

	assert.Nil(t, UpdateUserPassword(db, u1.ID, "new-password"))
	assert.Nil(t, UpdateUserDisable(db, u1.ID, true))
	assert.Nil(t, UpdateUserLoginInfo(db, u1.ID, "127.0.0.1", "内网IP"))

	got, err := GetUserAuthInfoById(db, u1.ID)
	assert.Nil(t, err)
	assert.Equal(t, "new-password", got.Password)
	assert.True(t, got.IsDisable)
	assert.Equal(t, "127.0.0.1", got.IpAddress)
	assert.NotNil(t, got.LastLoginTime)
}

// 更新昵称与角色: 角色关联要整体替换, 不能残留旧角色
func TestUpdateUserNicknameAndRole(t *testing.T) {
	db := newModelDB(t)
	role1 := Role{Name: "role1", Label: "角色1"}
	role2 := Role{Name: "role2", Label: "角色2"}
	assert.Nil(t, db.Create(&role1).Error)
	assert.Nil(t, db.Create(&role2).Error)

	// 先占掉一个 user_info id, 让 user_auth.id 和 user_auth.user_info_id 不相等,
	// 否则用错主键的 bug 会被巧合掩盖
	assert.Nil(t, db.Create(&UserInfo{Nickname: "占位"}).Error)

	user := UserAuth{Username: "test", Password: "x", UserInfo: &UserInfo{Nickname: "测试"}}
	assert.Nil(t, db.Create(&user).Error)
	assert.NotEqual(t, user.ID, user.UserInfoId)
	assert.Nil(t, db.Create(&UserAuthRole{UserAuthId: user.ID, RoleId: role1.ID}).Error)

	assert.Nil(t, UpdateUserNicknameAndRole(db, user.ID, "新昵称", []int{role2.ID}))

	info, err := GetUserInfoById(db, user.UserInfoId)
	assert.Nil(t, err)
	assert.Equal(t, "新昵称", info.Nickname)

	roleIds, err := GetRoleIdsByUserId(db, user.ID)
	assert.Nil(t, err)
	assert.Equal(t, []int{role2.ID}, roleIds, "旧角色应该被清掉")
}

func TestUserAuthMarshalBinary(t *testing.T) {
	user := UserAuth{Username: "test", Password: "secret"}

	data, err := user.MarshalBinary()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "test")
	// 密码是 json:"-", 不能被序列化进 Redis
	assert.NotContains(t, string(data), "secret")
}
