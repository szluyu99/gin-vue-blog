package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 角色接口: 新增/编辑时维护 role_resource 与 role_menu 关联
func TestRoleAPI(t *testing.T) {
	env := newTestEnv(t)
	api := Role{}
	env.engine.GET("/role/list", api.GetTreeList)
	env.engine.GET("/role/option", api.GetOption)
	env.engine.POST("/role", api.SaveOrUpdate)
	env.engine.DELETE("/role", api.Delete)

	// 准备一个资源和一个菜单
	resource, err := model.AddResource(env.db, "文章列表", "/article/list", "GET", false)
	assert.Nil(t, err)
	menu := model.Menu{Name: "文章管理", Path: "/article"}
	assert.Nil(t, env.db.Create(&menu).Error)

	// 新增角色: 只传名称和标签
	resp := env.do(t, http.MethodPost, "/role", map[string]any{
		"name": "editor", "label": "编辑",
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var role model.Role
	assert.Nil(t, env.db.Where("name", "editor").First(&role).Error)

	// 缺少必填的 label
	resp = env.do(t, http.MethodPost, "/role", map[string]any{"name": "no-label"})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 编辑角色: 关联资源和菜单
	resp = env.do(t, http.MethodPost, "/role", map[string]any{
		"id": role.ID, "name": "editor", "label": "编辑者",
		"resource_ids": []int{resource.ID}, "menu_ids": []int{menu.ID},
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 列表中带出关联的 id
	resp = env.do(t, http.MethodGet, "/role/list?page_num=1&page_size=10", nil)
	var page PageResult[model.RoleVO]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "编辑者", page.List[0].Label)
	assert.Equal(t, []int{resource.ID}, page.List[0].ResourceIds)
	assert.Equal(t, []int{menu.ID}, page.List[0].MenuIds)

	// 关联生效后, 权限校验应通过
	pass, err := model.CheckRoleAuth(env.db, role.ID, "/article/list", "GET")
	assert.Nil(t, err)
	assert.True(t, pass)

	// 选项列表
	resp = env.do(t, http.MethodGet, "/role/option", nil)
	var options []model.OptionVO
	decodeData(t, resp.Data, &options)
	assert.Len(t, options, 1)

	// 删除角色, 关联关系一并清理
	resp = env.do(t, http.MethodDelete, "/role", []int{role.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var count int64
	env.db.Model(&model.RoleResource{}).Where("role_id = ?", role.ID).Count(&count)
	assert.Equal(t, int64(0), count)
	env.db.Model(&model.RoleMenu{}).Where("role_id = ?", role.ID).Count(&count)
	assert.Equal(t, int64(0), count)

	resp = env.do(t, http.MethodGet, "/role/list", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(0), page.Total)
}
