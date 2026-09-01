package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 菜单接口: 树形列表 + 删除校验
func TestMenuAPI(t *testing.T) {
	env := newTestEnv(t)
	api := Menu{}
	env.engine.GET("/menu/list", api.GetTreeList)
	env.engine.GET("/menu/option", api.GetOption)
	env.engine.POST("/menu", api.SaveOrUpdate)
	env.engine.DELETE("/menu/:id", api.Delete)

	// 一级菜单
	resp := env.do(t, http.MethodPost, "/menu", map[string]any{
		"name": "文章管理", "path": "/article", "component": "Layout", "order_num": 1,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var parent model.Menu
	assert.Nil(t, env.db.Where("name", "文章管理").First(&parent).Error)

	// 二级菜单
	resp = env.do(t, http.MethodPost, "/menu", map[string]any{
		"name": "文章列表", "path": "list", "component": "/article/list/index.vue",
		"parent_id": parent.ID, "order_num": 1,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 树形列表
	resp = env.do(t, http.MethodGet, "/menu/list", nil)
	var tree []MenuTreeVO
	decodeData(t, resp.Data, &tree)
	assert.Len(t, tree, 1)
	assert.Equal(t, "文章管理", tree[0].Name)
	assert.Len(t, tree[0].Children, 1)

	child := tree[0].Children[0]

	// 选项列表
	resp = env.do(t, http.MethodGet, "/menu/option", nil)
	var options []TreeOptionVO
	decodeData(t, resp.Data, &options)
	assert.Len(t, options, 1)
	assert.Len(t, options[0].Children, 1)

	// 一级菜单下有子菜单, 不允许删除
	resp = env.do(t, http.MethodDelete, "/menu/"+itoa(parent.ID), nil)
	assert.Equal(t, g.ErrMenuHasChildren.Code(), resp.Code)

	// 被角色使用的菜单, 不允许删除
	env.db.Create(&model.RoleMenu{RoleId: 1, MenuId: child.ID})
	resp = env.do(t, http.MethodDelete, "/menu/"+itoa(child.ID), nil)
	assert.Equal(t, g.ErrMenuUsedByRole.Code(), resp.Code)

	// 不存在的菜单
	resp = env.do(t, http.MethodDelete, "/menu/9999", nil)
	assert.Equal(t, g.ErrMenuNotExist.Code(), resp.Code)

	// 解除角色关联后可以删除
	env.db.Where("menu_id = ?", child.ID).Delete(&model.RoleMenu{})
	resp = env.do(t, http.MethodDelete, "/menu/"+itoa(child.ID), nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
}

// 当前用户菜单: 超管拿全部, 普通用户只拿角色关联的菜单
func TestMenuGetUserMenu(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/menu/user/list", (&Menu{}).GetUserMenu)

	// 准备两个菜单, 其中一个关联给角色
	menu1 := model.Menu{Name: "文章管理", Path: "/article", OrderNum: 1}
	menu2 := model.Menu{Name: "系统管理", Path: "/system", OrderNum: 2}
	assert.Nil(t, env.db.Create(&menu1).Error)
	assert.Nil(t, env.db.Create(&menu2).Error)

	role, err := model.AddRoleWithResources(env.db, "guest", "访客", nil)
	assert.Nil(t, err)
	assert.Nil(t, env.db.Create(&model.RoleMenu{RoleId: role.ID, MenuId: menu1.ID}).Error)

	// 未登录
	resp := env.do(t, http.MethodGet, "/menu/user/list", nil)
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)

	// 超管: 拿到全部菜单
	super := env.loginAs(1, "admin")
	super.IsSuper = true
	resp = env.do(t, http.MethodGet, "/menu/user/list", nil)
	var menus []MenuTreeVO
	decodeData(t, resp.Data, &menus)
	assert.Len(t, menus, 2)

	// 普通用户: 只拿到角色关联的菜单
	user := model.UserAuth{Username: "guest", Roles: []*model.Role{role}}
	assert.Nil(t, env.db.Create(&user).Error)

	normal := env.loginAs(user.ID, "guest")
	normal.IsSuper = false
	resp = env.do(t, http.MethodGet, "/menu/user/list", nil)
	decodeData(t, resp.Data, &menus)
	assert.Len(t, menus, 1)
	assert.Equal(t, "文章管理", menus[0].Name)
}
