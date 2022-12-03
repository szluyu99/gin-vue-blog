package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"sort"
)

type Menu struct{}

// 获取某个用户的菜单列表(树形)
func (s *Menu) GetUserMenuTree(userInfoId int) []resp.UserMenuVo {
	// 从数据库查出用户菜单列表(非树形)
	menuList := menuDao.GetMenusByUserInfoId(userInfoId)
	// 过滤出一级菜单 (parent_id = 0)
	firstLevelMenuList := s.getFirstLevelMenus(menuList)
	// 获取存储子菜单的 map
	menuChildrenMap := s.getMenuChildrenMap(menuList)
	// 将一级 Menu 列表转成 UserMenuVo 列表, 主要处理其 Children
	return s.menus2UserMenuVos(firstLevelMenuList, menuChildrenMap)
}

// 获取菜单列表(树形)
func (s *Menu) GetTreeList(req req.PageQuery) []resp.MenuVo {
	// 从数据库中查询出菜单列表(非树形)
	menuList := menuDao.GetMenus(req)
	// 过滤出一级菜单 (parent_id = 0)
	firstLevelMenuList := s.getFirstLevelMenus(menuList)
	// 获取存储子菜单的 map
	menuChildrenMap := s.getMenuChildrenMap(menuList)
	// 将一级 Menu 列表转成 MenuVo 列表, 主要处理其 Children
	return s.menus2MenuVos(firstLevelMenuList, menuChildrenMap)
}

// 获取菜单选项列表(树形)
func (s *Menu) GetOptionList() []resp.TreeOptionVo {
	var resList = make([]resp.TreeOptionVo, 0)

	menus := dao.List([]model.Menu{}, "id, name, parent_id, order_num", "", "")
	firstLevelMenus := s.getFirstLevelMenus(menus)
	childrenMap := s.getMenuChildrenMap(menus)

	for _, item := range firstLevelMenus {
		var childrenOptionVos []resp.TreeOptionVo
		children := childrenMap[item.ID]
		if len(children) > 0 {
			s.sortMenu(children)
			for _, menu := range children {
				childrenOptionVos = append(childrenOptionVos, resp.TreeOptionVo{
					ID:    menu.ID,
					Label: menu.Name,
				})
			}
		}
		resList = append(resList, resp.TreeOptionVo{
			ID:       item.ID,
			Label:    item.Name,
			Children: childrenOptionVos,
		})
	}
	return resList
}

func (*Menu) SaveOrUpdate(req req.SaveOrUpdateMenu) (code int) {
	// 检查菜单名已经存在
	existByName := dao.GetOne(model.Menu{}, "name", req.Name)
	if existByName.ID != 0 && existByName.ID != req.ID {
		return r.ERROR_MENU_NAME_EXIST
	}
	if req.ID != 0 {
		dao.UpdatesMap(&model.Menu{}, utils.Struct2Map(req), "id", req.ID)
	} else {
		data := utils.CopyProperties[model.Menu](req)
		dao.Create(&data)
	}
	return r.OK
}

// 删除菜单
func (*Menu) Delete(menuId int) (code int) {
	// 检查要删除的菜单是否存在
	existMenuById := dao.GetOne(model.Menu{}, "id", menuId)
	if existMenuById.ID == 0 {
		return r.ERROR_MENU_NAME_EXIST
	}
	// * 检查 role_menu 下是否有数据
	existRoleMenu := dao.GetOne(model.RoleMenu{}, "menu_id", menuId)
	if existRoleMenu.MenuId != 0 {
		return r.ERROR_MENU_USED_BY_ROLE
	}
	// * 如果是一级菜单, 检查其是否有子菜单
	if existMenuById.ParentId == 0 {
		if dao.Count(model.Menu{}, "parent_id", menuId) != 0 {
			return r.ERROR_MENU_HAS_CHILDREN
		}
	}
	// 删除菜单
	dao.Delete(model.Menu{}, "id", menuId)
	return r.OK
}

// 筛选出一级菜单 (parentId == 0 的菜单)
func (s *Menu) getFirstLevelMenus(menuList []model.Menu) []model.Menu {
	firstLevelMenus := make([]model.Menu, 0)
	for _, menu := range menuList {
		if menu.ParentId == 0 {
			firstLevelMenus = append(firstLevelMenus, menu)
		}
	}
	s.sortMenu(firstLevelMenus) // 以 orderNum 降序排序
	return firstLevelMenus
}

// 获取 map: key 是菜单 ID, value 是该菜单对应的子菜单列表
func (*Menu) getMenuChildrenMap(menus []model.Menu) map[int][]model.Menu {
	mcMap := make(map[int][]model.Menu)
	for _, menu := range menus {
		if menu.ParentId != 0 {
			mcMap[menu.ParentId] = append(mcMap[menu.ParentId], menu)
		}
	}
	return mcMap
}

// 构建用户菜单的树形结构数据, []model.Menu => []resp.UserMenuVo
func (s *Menu) menus2UserMenuVos(firstLevelMenuList []model.Menu, childrenMap map[int][]model.Menu) []resp.UserMenuVo {
	resList := make([]resp.UserMenuVo, 0)

	// 遍历一级 Menu, 将其元素构造成 UserMenu
	for _, firstLevelMenu := range firstLevelMenuList {
		var userMenu resp.UserMenuVo           // 当前 [用户菜单]
		var userMenuChildren []resp.UserMenuVo // 当前 [用户菜单] 的 [子用户菜单]

		children := childrenMap[firstLevelMenu.ID] // 子菜单
		if len(children) > 0 {                     // 存在子菜单
			userMenu = s.menu2UserMenuVo(firstLevelMenu) // [菜单] -> [用户菜单]
			// userMenu.Path = ""                           // TODO! 外层一定是 "" 吗?
			s.sortMenu(children) // 对子菜单按照 OrderNum 排序
			// 遍历子菜单, 将其构造成 用户菜单
			for _, child := range children {
				userMenuChildren = append(userMenuChildren, s.menu2UserMenuVo(child))
			}
		} else { // 没有子菜单, 利用一级菜单构造一个用户菜单(Layout), 将原本的一级菜单作为子菜单变成新菜单的 children
			userMenu = resp.UserMenuVo{
				ID:        firstLevelMenu.ID,
				Path:      firstLevelMenu.Path,
				Name:      firstLevelMenu.Name,      // *
				Component: firstLevelMenu.Component, // ! "Layout" ?
				OrderNum:  firstLevelMenu.OrderNum,
				Hidden:    firstLevelMenu.IsHidden == 1,
				KeepAlive: firstLevelMenu.KeepAlive == 1,
				Redirect:  firstLevelMenu.Redirect,
			}
			tmpUserMenu := s.menu2UserMenuVo(firstLevelMenu)
			// tmpUserMenu.Path = "" // TODO! 考虑一下
			userMenuChildren = append(userMenuChildren, tmpUserMenu)
		}
		userMenu.Children = userMenuChildren
		resList = append(resList, userMenu)
	}
	return resList
}

// 构建菜单列表的树形结构数据, []model.Menu => []resp.MenuVo
func (s *Menu) menus2MenuVos(firstLevelMenuList []model.Menu, menuChildrenMap map[int][]model.Menu) []resp.MenuVo {
	resMenuVos := make([]resp.MenuVo, 0)
	// 遍历一级菜单
	for _, firstLevelMenu := range firstLevelMenuList {
		// 尝试将一级菜单转 MenuV, 还需要处理 Children
		menuVo := s.menu2MenuVo(firstLevelMenu)
		// 获取到当前菜单的子菜单
		menuChildren := menuChildrenMap[firstLevelMenu.ID]
		s.sortMenu(menuChildren)
		// 对子菜单进行 []menu => []menuVo
		childMenuVos := make([]resp.MenuVo, 0)
		for _, childMenu := range menuChildren {
			childMenuVos = append(childMenuVos, s.menu2MenuVo(childMenu))
		}
		menuVo.Children = childMenuVos
		resMenuVos = append(resMenuVos, menuVo)
		delete(menuChildrenMap, firstLevelMenu.ID) // 从 map 中删除已经构建完的菜单
	}
	// 处理剩下的子菜单 TODO: 思考
	if len(menuChildrenMap) > 0 {
		var menuChildren []model.Menu
		for _, v := range menuChildrenMap {
			menuChildren = append(menuChildren, v...)
		}
		s.sortMenu(menuChildren)
		for _, menu := range menuChildren {
			resMenuVos = append(resMenuVos, s.menu2MenuVo(menu))
		}
	}
	return resMenuVos
}

// TODO: 修改成使用 MapStruct?
// model.Menu => resp.MenuVo
func (*Menu) menu2MenuVo(menu model.Menu) resp.MenuVo {
	return resp.MenuVo{
		ID:        menu.ID,
		Name:      menu.Name,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		CreatedAt: menu.CreatedAt,
		OrderNum:  menu.OrderNum,
		IsHidden:  int(menu.IsHidden),
		ParentId:  menu.ParentId, // *
		Redirect:  menu.Redirect,
		KeepAlive: menu.KeepAlive,
		// Children:  make([]resp.MenuVo, 0),
	}
}

// model.Menu => resp.UserMenuVo
func (*Menu) menu2UserMenuVo(menu model.Menu) resp.UserMenuVo {
	return resp.UserMenuVo{
		ID:        menu.ID,
		Name:      menu.Name,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		OrderNum:  menu.OrderNum,
		Hidden:    menu.IsHidden == 1,
		KeepAlive: menu.KeepAlive == 1,
		ParentId:  menu.ParentId, // *
		Redirect:  menu.Redirect,
	}
}

// 对菜单排序: 以 OrderNum 字段进行排序
func (*Menu) sortMenu(menus []model.Menu) {
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].OrderNum < menus[j].OrderNum
	})
}

/*
菜单数据: []menuVo
{
	"id": 1,
	"name": "首页",
	"path": "/",
	"component": "/home/Home.vue",
	"icon": "el-icon-myshouye",
	"createTime": "2021-01-26T17:06:51",
	"orderNum": 1,
	"isDisable": null,
	"isHidden": 0,
	"children": []
},
{
	"id": 2,
	"name": "文章管理",
	"path": "/article-submenu",
	"component": "Layout",
	"icon": "el-icon-mywenzhang-copy",
	"createTime": "2021-01-25T20:43:07",
	"orderNum": 2,
	"isDisable": null,
	"isHidden": 0,
	"children": [
		{
			"id": 6,
			"name": "发布文章",
			"path": "/articles",
			"component": "/article/Article.vue",
			"icon": "el-icon-myfabiaowenzhang",
			"createTime": "2021-01-26T14:30:48",
			"orderNum": 1,
			"isDisable": null,
			"isHidden": 0,
			"children": null
		},
		{
			"id": 7,
			"name": "修改文章",
			"path": "/articles/*",
			"component": "/article/Article.vue",
			"icon": "el-icon-myfabiaowenzhang",
			"createTime": "2021-01-26T14:31:32",
			"orderNum": 2,
			"isDisable": null,
			"isHidden": 1,
			"children": null
		},
		{
			"id": 8,
			"name": "文章列表",
			"path": "/article-list",
			"component": "/article/ArticleList.vue",
			"icon": "el-icon-mywenzhangliebiao",
			"createTime": "2021-01-26T14:32:13",
			"orderNum": 3,
			"isDisable": null,
			"isHidden": 0,
			"children": null
		}
	]
},
*/

/*
用户菜单数据: []userMenuVo

{
	"name": null,
	"path": "/",
	"component": "Layout",
	"icon": null,
	"hidden": false,
	"children": [
		{
			"name": "首页",
			"path": "",
			"component": "/home/Home.vue",
			"icon": "el-icon-myshouye",
			"hidden": null,
			"children": null
		}
	]
},
{
	"name": "文章管理",
	"path": "/article-submenu",
	"component": "Layout",
	"icon": "el-icon-mywenzhang-copy",
	"hidden": false,
	"children": [
		{
			"name": "发布文章",
			"path": "/articles",
			"component": "/article/Article.vue",
			"icon": "el-icon-myfabiaowenzhang",
			"hidden": false,
			"children": null
		},
		{
			"name": "修改文章",
			"path": "/articles/*",
			"component": "/article/Article.vue",
			"icon": "el-icon-myfabiaowenzhang",
			"hidden": true,
			"children": null
		},
		{
			"name": "文章列表",
			"path": "/article-list",
			"component": "/article/ArticleList.vue",
			"icon": "el-icon-mywenzhangliebiao",
			"hidden": false,
			"children": null
		}
	]
}
*/
