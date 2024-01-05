package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Menu struct{}

type MenuTreeVO struct {
	model.Menu
	Children []MenuTreeVO `json:"children"`
}

// 获取当前用户菜单: 生成后台管理界面的菜单
func (*Menu) GetUserMenu(c *gin.Context) {
	db := GetDB(c)
	auth, _ := CurrentUserAuth(c)

	var menus []model.Menu
	var err error

	if auth.IsSuper {
		menus, err = model.GetAllMenuList(db)
	} else {
		menus, err = model.GetMenuListByUserId(GetDB(c), auth.ID)
	}

	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, menus2MenuVos(menus))
}

func (*Menu) GetTreeList(c *gin.Context) {
	keyword := c.Query("keyword")

	menuList, _, err := model.GetMenuList(GetDB(c), keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, menus2MenuVos(menuList))
}

func (*Menu) SaveOrUpdate(c *gin.Context) {
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if err := model.SaveOrUpdateMenu(GetDB(c), &req); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

func (*Menu) Delete(c *gin.Context) {
	menuId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// 检查要删除的菜单是否被角色使用
	use, _ := model.CheckMenuInUse(db, menuId)
	if use {
		ReturnError(c, g.ErrMenuUsedByRole, nil)
		return
	}

	// 如果是一级菜单, 检查其是否有子菜单
	menu, err := model.GetMenuById(db, menuId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrMenuNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 一级菜单下有子菜单, 不允许删除
	if menu.ParentId == 0 {
		has, _ := model.CheckMenuHasChild(db, menuId)
		if has {
			ReturnError(c, g.ErrMenuHasChildren, nil)
			return
		}
	}

	if err = model.DeleteMenu(db, menuId); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

func (*Menu) GetOption(c *gin.Context) {
	menus, _, err := model.GetMenuList(GetDB(c), "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	result := make([]TreeOptionVO, 0)
	for _, menu := range menus2MenuVos(menus) {
		option := TreeOptionVO{ID: menu.ID, Label: menu.Name}
		for _, child := range menu.Children {
			option.Children = append(option.Children, TreeOptionVO{ID: child.ID, Label: child.Name})
		}
		result = append(result, option)
	}

	ReturnSuccess(c, result)
}

// 构建菜单列表的树形结构数据, []Menu => []MenuVo
func menus2MenuVos(menus []model.Menu) []MenuTreeVO {
	result := make([]MenuTreeVO, 0)

	firstLevelMenus := getFirstLevelMenus(menus)
	childrenMap := getMenuChildrenMap(menus)

	for _, first := range firstLevelMenus {
		menu := MenuTreeVO{Menu: first}
		for _, childMenu := range childrenMap[first.ID] {
			menu.Children = append(menu.Children, MenuTreeVO{Menu: childMenu})
		}
		delete(childrenMap, first.ID)
		result = append(result, menu)
	}

	sortMenu(result)
	return result
}

// 筛选出一级菜单 (parentId == 0 的菜单)
func getFirstLevelMenus(menuList []model.Menu) []model.Menu {
	firstLevelMenus := make([]model.Menu, 0)
	for _, menu := range menuList {
		if menu.ParentId == 0 {
			firstLevelMenus = append(firstLevelMenus, menu)
		}
	}
	return firstLevelMenus
}

// key 是菜单 ID, value 是该菜单对应的子菜单列表
func getMenuChildrenMap(menus []model.Menu) map[int][]model.Menu {
	childrenMap := make(map[int][]model.Menu)
	for _, menu := range menus {
		if menu.ParentId != 0 {
			childrenMap[menu.ParentId] = append(childrenMap[menu.ParentId], menu)
		}
	}
	return childrenMap
}

// 以 orderNum 降序排序，包括子菜单
func sortMenu(menus []MenuTreeVO) {
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].OrderNum < menus[j].OrderNum
	})
	for i := range menus {
		sort.Slice(menus[i].Children, func(j, k int) bool {
			return menus[i].Children[j].OrderNum < menus[i].Children[k].OrderNum
		})
	}
}

// 构建用户菜单的树形结构数据, []Menu => []MenuVO
// func menus2UserMenuVos(menus []model.Menu) []MenuTreeVO {
// 	firstLevelMenuList := getFirstLevelMenus(menus)
// 	childrenMap := getMenuChildrenMap(menus)

// 	result := make([]MenuTreeVO, 0)

// 	// 遍历一级 Menu, 将其元素构造成 UserMenu
// 	for _, firstLevelMenu := range firstLevelMenuList {
// 		var menuVO MenuTreeVO             // 当前 [用户菜单]
// 		var userMenuChildren []MenuTreeVO // 当前 [用户菜单] 的 [子用户菜单]

// 		children := childrenMap[firstLevelMenu.ID] // 子菜单
// 		if len(children) > 0 {                     // 存在子菜单
// 			menuVO = menu2UserMenuVo(firstLevelMenu) // [菜单] -> [用户菜单]
// 			// userMenu.Path = ""                           // TODO! 外层一定是 "" 吗?
// 			sortMenu(children) // 对子菜单按照 OrderNum 排序
// 			// 遍历子菜单, 将其构造成 用户菜单
// 			for _, child := range children {
// 				userMenuChildren = append(userMenuChildren, menu2UserMenuVo(child))
// 			}
// 		} else { // 没有子菜单, 利用一级菜单构造一个用户菜单(Layout), 将原本的一级菜单作为子菜单变成新菜单的 children
// 			menuVO = MenuTreeVO{
// 				Menu: firstLevelMenu,
// 				// ID:        firstLevelMenu.ID,
// 				// Path:      firstLevelMenu.Path,
// 				// Name:      firstLevelMenu.Name,      // *
// 				// Component: firstLevelMenu.Component, // ! "Layout" ?
// 				// OrderNum:  firstLevelMenu.OrderNum,
// 				// IsHidden:  firstLevelMenu.Hidden,
// 				// KeepAlive: firstLevelMenu.KeepAlive,
// 				// Redirect:  firstLevelMenu.Redirect,
// 			}
// 			tmpUserMenu := menu2UserMenuVo(firstLevelMenu)
// 			// tmpUserMenu.Path = "" // TODO! 考虑一下
// 			userMenuChildren = append(userMenuChildren, tmpUserMenu)
// 		}
// 		menuVO.Children = userMenuChildren
// 		result = append(result, menuVO)
// 	}
// 	return result
// }

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
