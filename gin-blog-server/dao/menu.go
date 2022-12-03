package dao

import (
	"gin-blog/model"
	"gin-blog/model/req"
)

type Menu struct{}

// 根据 userInfoId 获取菜单列表(非树形结构): 关联 user_role, role_menu, menu 表
func (*Menu) GetMenusByUserInfoId(userInfoId int) []model.Menu {
	list := make([]model.Menu, 0)
	DB.Table("user_role ur").
		Distinct("m.id", "name", "path", "component", "icon", "is_hidden", "keep_alive", "redirect", "parent_id", "order_num").
		Where("user_id = ?", userInfoId).
		Joins("JOIN role_menu rm ON ur.role_id = rm.role_id").
		Joins("JOIN menu m ON rm.menu_id = m.id").
		Find(&list)
	return list
}

// 获取菜单列表(非树形结构)
func (*Menu) GetMenus(req req.PageQuery) []model.Menu {
	var list []model.Menu
	if req.Keyword != "" {
		list = List([]model.Menu{}, "*", "", "name like?", "%"+req.Keyword+"%")
	} else {
		list = List([]model.Menu{}, "*", "", "")
	}
	return list
}
