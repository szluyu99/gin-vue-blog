package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 资源接口: 树形结构 + 删除前的各种校验
func TestResourceAPI(t *testing.T) {
	env := newTestEnv(t)
	api := Resource{}
	env.engine.GET("/resource/list", api.GetTreeList)
	env.engine.GET("/resource/option", api.GetOption)
	env.engine.POST("/resource", api.SaveOrUpdate)
	env.engine.PUT("/resource/anonymous", api.UpdateAnonymous)
	env.engine.DELETE("/resource/:id", api.Delete)

	// 新增模块(一级资源)与其下的接口
	resp := env.do(t, http.MethodPost, "/resource", map[string]any{"name": "文章模块"})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var module model.Resource
	assert.Nil(t, env.db.Where("name", "文章模块").First(&module).Error)

	resp = env.do(t, http.MethodPost, "/resource", map[string]any{
		"name": "文章列表", "url": "/article/list", "request_method": "GET", "parent_id": module.ID,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)

	// 树形列表: 模块下带一个子资源
	resp = env.do(t, http.MethodGet, "/resource/list", nil)
	var tree []ResourceTreeVO
	decodeData(t, resp.Data, &tree)
	assert.Len(t, tree, 1)
	assert.Equal(t, "文章模块", tree[0].Name)
	assert.Len(t, tree[0].Children, 1)
	assert.Equal(t, "/article/list", tree[0].Children[0].Url)

	// 选项列表结构一致
	resp = env.do(t, http.MethodGet, "/resource/option", nil)
	var options []TreeOptionVO
	decodeData(t, resp.Data, &options)
	assert.Len(t, options, 1)
	assert.Len(t, options[0].Children, 1)

	child := tree[0].Children[0]

	// 修改匿名访问
	resp = env.do(t, http.MethodPut, "/resource/anonymous", map[string]any{
		"id": child.ID, "is_anonymous": true,
	})
	assert.Equal(t, g.SUCCESS, resp.Code)
	got, _ := model.GetResourceById(env.db, child.ID)
	assert.True(t, got.Anonymous)

	// 模块下有子资源, 不允许删除
	resp = env.do(t, http.MethodDelete, "/resource/"+itoa(module.ID), nil)
	assert.Equal(t, g.ErrResourceHasChildren.Code(), resp.Code)

	// 被角色使用的资源, 不允许删除
	model.AddRoleWithResources(env.db, "role", "角色", []int{child.ID})
	resp = env.do(t, http.MethodDelete, "/resource/"+itoa(child.ID), nil)
	assert.Equal(t, g.ErrResourceUsedByRole.Code(), resp.Code)

	// 不存在的资源
	resp = env.do(t, http.MethodDelete, "/resource/9999", nil)
	assert.Equal(t, g.ErrResourceNotExist.Code(), resp.Code)

	// id 不是数字
	resp = env.do(t, http.MethodDelete, "/resource/abc", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 解除角色关联后可以删除, 模块随之变成空模块也能删除
	env.db.Where("resource_id = ?", child.ID).Delete(&model.RoleResource{})
	resp = env.do(t, http.MethodDelete, "/resource/"+itoa(child.ID), nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	resp = env.do(t, http.MethodDelete, "/resource/"+itoa(module.ID), nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
}
