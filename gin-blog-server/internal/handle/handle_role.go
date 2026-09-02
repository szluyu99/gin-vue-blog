package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Role struct{}

// 新增/编辑 角色, 关联维护 role_resource, role_menu
type AddOrEditRoleReq struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Label       string `json:"label" binding:"required"`
	IsDisable   bool   `json:"is_disable"`
	ResourceIds []int  `json:"resource_ids"` // 资源 id 列表
	MenuIds     []int  `json:"menu_ids"`     // 菜单 id 列表
}

// @Summary 获取角色选项
// @Description 用于给用户分配角色
// @Tags Role
// @Produce json
// @Success 0 {object} Response[[]model.OptionVO]
// @Security ApiKeyAuth
// @Router /role/option [get]
func (*Role) GetOption(c *gin.Context) {
	list, err := model.GetRoleOption(GetDB(c))
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, list)
}

// @Summary 条件查询角色列表
// @Description 角色列表, 附带其资源和菜单 ID
// @Tags Role
// @Produce json
// @Param keyword query string false "关键字"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.RoleVO]]
// @Security ApiKeyAuth
// @Router /role/list [get]
func (*Role) GetTreeList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	result := make([]model.RoleVO, 0)

	list, total, err := model.GetRoleList(db, query.Page, query.Size, query.Keyword)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	for _, role := range list {
		role.ResourceIds, _ = model.GetResourceIdsByRoleId(db, role.ID)
		role.MenuIds, _ = model.GetMenuIdsByRoleId(db, role.ID)
		result = append(result, role)
	}

	ReturnSuccess(c, PageResult[model.RoleVO]{
		Size:  query.Size,
		Page:  query.Page,
		Total: total,
		List:  result,
	})
}

// @Summary 新增或编辑角色
// @Description 新增或编辑角色, 同时维护角色的资源与菜单关联
// @Tags Role
// @Accept json
// @Produce json
// @Param form body AddOrEditRoleReq true "新增或编辑角色"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /role [post]
func (*Role) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	if req.ID == 0 {
		err := model.SaveRole(db, req.Name, req.Label, req.ResourceIds, req.MenuIds)
		if err != nil {
			ReturnError(c, g.ErrDbOp, err)
			return
		}
	} else {
		err := model.UpdateRole(db, req.ID, req.Name, req.Label, req.IsDisable, req.ResourceIds, req.MenuIds)
		if err != nil {
			ReturnError(c, g.ErrDbOp, err)
			return
		}
	}

	ReturnSuccess(c, nil)
}

// @Summary 删除角色（批量）
// @Description 根据 ID 数组删除角色, 同时清理资源与菜单关联
// @Tags Role
// @Accept json
// @Produce json
// @Param ids body []int true "角色 ID 数组"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /role [delete]
func (*Role) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.DeleteRoles(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}
