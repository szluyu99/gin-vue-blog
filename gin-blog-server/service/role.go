package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
)

type Role struct{}

// 获取用户的角色列表
func (*Role) GetOption() []resp.OptionVo {
	return roleDao.GetOption()
}

// * 查询出角色列表(树形)
func (*Role) GetTreeList(req req.PageQuery) resp.PageResult[[]resp.RoleVo] {
	treeList := make([]resp.RoleVo, 0)
	// 角色列表(非树形)
	list, total := roleDao.GetList(req)
	// 构造角色列表(树形)
	for _, role := range list {
		// 根据 [角色id] 查询出 [资源id列表]
		role.ResourceIds = roleDao.GetResourcesByRoldId(role.ID)
		// 根据 [角色id] 查询出 [菜单id列表]
		role.MenuIds = roleDao.GetMenusByRoleId(role.ID)
		treeList = append(treeList, role)
	}
	return resp.PageResult[[]resp.RoleVo]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    total,
		List:     treeList,
	}
}

// * 新增/更新 角色, 同时维护 role_menu, role_resource, casbin_rule 中的关联数据
func (s *Role) SaveOrUpdate(req req.SaveOrUpdateRole) (code int) {
	// 检查角色名是否存在
	existByName := dao.GetOne(model.Role{}, "name", req.Name)
	if existByName.ID != 0 && existByName.ID != req.ID {
		return r.ERROR_ROLE_NAME_EXIST
	}
	role := utils.CopyProperties[model.Role](req) // vo -> po
	if role.ID == 0 {
		dao.Create(&role)
		// ! 默认添加 anonymous 和 logout 角色
		utils.Casbin.AddRoleForUser(role.Label, "anonymous")
		utils.Casbin.AddRoleForUser(role.Label, "logout")
	} else {
		// 检查要更新的角色ID是否存在
		existRoleById := dao.GetOne(model.Role{}, "id", req.ID)
		if existRoleById.ID == 0 {
			return r.ERROR_ROLE_NOT_EXIST
		}
		//! 关联更新 casbin_rule 表中的 v0 (role_label)
		utils.Casbin.UpdateRoleLabels(existRoleById.Label, role.Label)
		dao.Update(&role)
	}

	// * 处理 ResourceIds 资源列表: 先清空和 role_id 相关的, 再重新根据请求数据添加新的
	// 删除 role_resource 和 casbin_rule 中的旧数据
	dao.Delete(model.RoleResource{}, "role_id = ?", req.ID)
	utils.Casbin.DeletePermissionForRole(role.Label, "") // !
	// * 往 role_resource 和 casbin_rule 中添加最新数据
	if len(req.ResourceIds) > 0 {
		// 构造 RoleSource po 对象列表, 并往数据库中插入数据
		var rrList []model.RoleResource
		for _, rid := range req.ResourceIds {
			rrList = append(rrList, model.RoleResource{RoleId: role.ID, ResourceId: rid})
		}
		dao.Create(&rrList)
		//! 构造批量添加 casbin_rule 的 rules
		rules := make([][]string, 0)
		resources := resourceDao.GetListByIds(req.ResourceIds)
		for _, resource := range resources {
			if resource.Url != "" && resource.RequestMethod != "" {
				rules = append(rules, []string{role.Label, resource.Url, resource.RequestMethod})
			}
		}
		utils.Casbin.AddPolicies(rules) // !
	}
	// *处理 MenuIds 菜单列表: 先清空和 role_id 相关的, 再重新根据请求数据添加新的
	dao.Delete(model.RoleMenu{}, "role_id", req.ID) // 删除 role_menu 中旧数据
	if len(req.MenuIds) > 0 {
		// 往 role_menu 中添加新的数据
		var rmList []model.RoleMenu
		for _, menuId := range req.MenuIds {
			rmList = append(rmList, model.RoleMenu{RoleId: role.ID, MenuId: menuId})
		}
		dao.Create(&rmList)
	}
	return r.OK
}

// * 根据 [id列表] 删除 role 中数据, 同时删除 role_menu, role_resource, casbin_rule 中的关联数据
func (*Role) Delete(ids []int) (code int) {
	labels := roleDao.GetLabelsByRoleIds(ids)
	// 判断提前结束
	if len(labels) == 0 {
		return r.OK
	}
	// !从 casbin_rule 删除对应角色的记录
	utils.Casbin.BatchDeletePermissions(labels)
	// 删除角色的关联数据
	dao.Delete(model.RoleMenu{}, "role_id in ?", ids)     // 关联删除 role_menu 数据
	dao.Delete(model.RoleResource{}, "role_id in ?", ids) // 关联删除 role_resource 数据
	dao.Delete(model.Role{}, "id in ?", ids)              // 删除 role 数据
	return r.OK
}
