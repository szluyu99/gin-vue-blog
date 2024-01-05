package main

import (
	"flag"
	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

func main() {
	configPath := flag.String("c", "../../config.yml", "配置文件路径")
	typeName := flag.String("t", "all", "要初始化的数据类型: config | auth | page | all")
	flag.Parse()

	// 根据命令行参数读取配置文件, 其他变量的初始化依赖于配置文件对象
	conf := g.ReadConfig(*configPath)

	//! 处理 sqlite3 数据库路径
	conf.SQLite.Dsn = "../" + conf.SQLite.Dsn
	conf.Server.DbLogMode = "silent"

	db := ginblog.InitDatabase(conf)

	switch *typeName {
	case "config":
		generateDefaultConfigs(db)
	case "auth":
		generateDefaultAuths(db)
	case "page":
		generateDefaultPages(db)
	case "all":
		fallthrough
	default:
		generateDefaultConfigs(db)
		generateDefaultPages(db)
		generateDefaultAuths(db)
	}
}

// 生成验证相关信息: 角色, 用户, 资源, 菜单
func generateDefaultAuths(db *gorm.DB) {
	generateDefaultRolesAndUsers(db)
	generateDefaultResources(db)
	generateDefaultMenus(db)
}

// 生成默认的页面
func generateDefaultPages(db *gorm.DB) {
	slog.Info("-----初始化博客页面 start-----")

	pages := []model.Page{
		{Name: "首页", Label: "home", Cover: "https://cdn.hahacode.cn/page/home.jpg"},
		{Name: "归档", Label: "archive", Cover: "https://cdn.hahacode.cn/page/archive.png"},
		{Name: "分类", Label: "category", Cover: "https://cdn.hahacode.cn/page/category.png"},
		{Name: "标签", Label: "tag", Cover: "https://cdn.hahacode.cn/page/tag.png"},
		{Name: "友链", Label: "link", Cover: "https://cdn.hahacode.cn/page/link.jpg"},
		{Name: "关于", Label: "about", Cover: "https://cdn.hahacode.cn/page/about.jpg"},
		{Name: "留言", Label: "message", Cover: "https://cdn.hahacode.cn/page/message.jpeg"},
		{Name: "个人中心", Label: "user", Cover: "https://cdn.hahacode.cn/page/user.jpg"},
		{Name: "相册", Label: "album", Cover: "https://cdn.hahacode.cn/page/album.png"},
		{Name: "错误页面", Label: "404", Cover: "https://cdn.hahacode.cn/page/404.jpg"},
		{Name: "文章列表", Label: "article_list", Cover: "https://cdn.hahacode.cn/page/article_list.jpg"},
	}

	for _, page := range pages {
		if err := db.Create(&page).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(page.Name + " 页面数据已经存在")
			} else {
				slog.Error(page.Name + " 页面初始化失败" + err.Error())
			}
		}
	}

	slog.Info("-----初始化博客页面 end-----")
}

// 生成默认配置信息
func generateDefaultConfigs(db *gorm.DB) {
	slog.Info("-----初始化博客配置 start-----")

	configs := []model.Config{
		{Key: "website_avatar", Value: "https://foruda.gitee.com/avatar/1677041571085433939/5221991_szluyu99_1614389421.png", Desc: "网站头像"},
		{Key: "website_name", Value: "阵雨的个人博客", Desc: "网站名称"},
		{Key: "website_author", Value: "阵雨", Desc: "网站作者"},
		{Key: "website_intro", Value: "往事随风而去", Desc: "网站介绍"},
		{Key: "website_notice", Value: "欢迎来到阵雨的个人博客，项目还在开发中...", Desc: "网站公告"},
		{Key: "website_createtime", Value: time.Now().Format(time.DateTime), Desc: "网站创建日期"},
		{Key: "website_record", Value: "粤ICP备2021032312号", Desc: "网站备案号"},
		{Key: "qq", Value: "123456789", Desc: "QQ"},
		{Key: "github", Value: "https://github.com/szluyu99", Desc: "github"},
		{Key: "gitee", Value: "https://gitee.com/szluyu99", Desc: "gitee"},
		{Key: "tourist_avatar", Value: "https://cdn.hahacode.cn/config/tourist_avatar.png", Desc: "默认游客头像"},
		{Key: "user_avatar", Value: "https://cdn.hahacode.cn/config/user_avatar.png", Desc: "默认用户头像"},
		{Key: "article_cover", Value: "https://cdn.hahacode.cn/config/default_article_cover.png", Desc: "默认文章封面"},
		{Key: "is_comment_review", Value: "true", Desc: "评论默认审核"},
		{Key: "is_message_review", Value: "true", Desc: "留言默认审核"},
	}

	for _, config := range configs {
		if err := db.Create(&config).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(config.Key + " 配置已经存在")
			} else {
				slog.Error(config.Key + " 配置初始化失败" + err.Error())
			}
		}
	}

	slog.Info("-----初始化博客配置 end-----")
}

// 生成 2 个默认角色及验证信息: admin, guest
func generateDefaultRolesAndUsers(db *gorm.DB) {
	slog.Info("-----初始化默认角色用户 start-----")

	roles := []model.Role{
		{Name: "admin", Label: "管理员"},
		{Name: "guest", Label: "游客"},
	}

	for i := range roles {
		if err := db.Create(&roles[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(roles[i].Name + " 角色已经存在")
			} else {
				slog.Error(roles[i].Name + " 角色初始化失败" + err.Error())
			}
		}
	}

	pwd, _ := utils.BcryptHash("123456")
	auths := []model.UserAuth{
		{
			Username: "admin",
			Password: pwd,
			UserInfo: &model.UserInfo{
				Nickname: "admin",
				Avatar:   "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
			},
		},
		{
			Username: "guest",
			Password: pwd,
			UserInfo: &model.UserInfo{
				Nickname: "guest",
				Avatar:   "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
			},
		},
	}

	for i := range auths {
		if err := db.Create(&auths[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(auths[i].Username + " 用户已经存在")
			} else {
				slog.Error(auths[i].Username + " 用户初始化失败" + err.Error())
			}
		}
		// 创建用户角色关联关系
		db.Create(&model.UserAuthRole{UserAuthId: auths[i].ID, RoleId: roles[i].ID})
	}

	slog.Info("-----初始化默认角色用户 end-----")
}

// 生成默认的接口资源
func generateDefaultResources(db *gorm.DB) {
	slog.Info("-----初始化接口资源 start-----")

	parents := []model.Resource{
		{Name: "文章模块"},
		{Name: "分类模块"},
		{Name: "标签模块"},
		{Name: "页面模块"},
		{Name: "友链模块"},
		{Name: "菜单模块"},
		{Name: "角色模块"},
		{Name: "资源模块"},
		{Name: "评论模块"},
		{Name: "留言模块"},
		{Name: "文件模块"},
		{Name: "博客信息模块"},
		{Name: "用户信息模块"},
		{Name: "操作日志模块"},
	}
	for i := range parents {
		if err := db.Create(&parents[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(parents[i].Name + " 资源已经存在")
			} else {
				slog.Error(parents[i].Name + " 资源初始化失败" + err.Error())
			}
		}
	}

	resources := []model.Resource{
		// 文章模块
		{Name: "文章列表", ParentId: parents[0].ID, Url: "/article/list", Method: "GET"},
		{Name: "文章详情", ParentId: parents[0].ID, Url: "/article/:id", Method: "GET"},
		{Name: "新增/编辑文章", ParentId: parents[0].ID, Url: "/article", Method: "POST"},
		{Name: "更新文章软删除", ParentId: parents[0].ID, Url: "/article/soft-delete", Method: "PUT"},
		{Name: "删除文章", ParentId: parents[0].ID, Url: "/article", Method: "DELETE"},
		{Name: "修改文章置顶", ParentId: parents[0].ID, Url: "/article/top", Method: "PUT"},
		{Name: "导出文章", ParentId: parents[0].ID, Url: "/article/export", Method: "POST"},
		{Name: "导入文章", ParentId: parents[0].ID, Url: "/article/import", Method: "POST"},
		// 分类模块
		{Name: "分类列表", ParentId: parents[1].ID, Url: "/category/list", Method: "GET"},
		{Name: "新增/编辑分类", ParentId: parents[1].ID, Url: "/category", Method: "POST"},
		{Name: "删除分类", ParentId: parents[1].ID, Url: "/category", Method: "DELETE"},
		{Name: "分类选项列表", ParentId: parents[1].ID, Url: "/category/option", Method: "GET"},
		// 标签模块
		{Name: "标签列表", ParentId: parents[2].ID, Url: "/tag/list", Method: "GET"},
		{Name: "新增/编辑标签", ParentId: parents[2].ID, Url: "/tag", Method: "POST"},
		{Name: "删除标签", ParentId: parents[2].ID, Url: "/tag", Method: "DELETE"},
		{Name: "标签选项列表", ParentId: parents[2].ID, Url: "/tag/option", Method: "GET"},
		// 页面模块
		{Name: "页面列表", ParentId: parents[3].ID, Url: "/page/list", Method: "GET"},
		{Name: "新增/编辑页面", ParentId: parents[3].ID, Url: "/page", Method: "POST"},
		{Name: "删除页面", ParentId: parents[3].ID, Url: "/page", Method: "DELETE"},
		// 友链模块
		{Name: "友链列表", ParentId: parents[4].ID, Url: "/link/list", Method: "GET"},
		{Name: "新增/编辑友链", ParentId: parents[4].ID, Url: "/link", Method: "POST"},
		{Name: "删除友链", ParentId: parents[4].ID, Url: "/link", Method: "DELETE"},
		// 菜单模块
		{Name: "菜单列表", ParentId: parents[5].ID, Url: "/menu/list", Method: "GET"},
		{Name: "新增/编辑菜单", ParentId: parents[5].ID, Url: "/menu", Method: "POST"},
		{Name: "删除菜单", ParentId: parents[5].ID, Url: "/menu", Method: "DELETE"},
		{Name: "菜单选项列表(树形)", ParentId: parents[5].ID, Url: "/menu/option", Method: "GET"},
		{Name: "获取当前用户菜单", ParentId: parents[5].ID, Url: "/menu/user/list", Method: "GET"},
		// 角色模块
		{Name: "角色列表", ParentId: parents[6].ID, Url: "/role/list", Method: "GET"},
		{Name: "新增/编辑角色", ParentId: parents[6].ID, Url: "/role", Method: "POST"},
		{Name: "删除角色", ParentId: parents[6].ID, Url: "/role", Method: "DELETE"},
		{Name: "角色选项列表", ParentId: parents[6].ID, Url: "/role/option", Method: "GET"},
		// 资源模块
		{Name: "资源列表", ParentId: parents[7].ID, Url: "/resource/list", Method: "GET"},
		{Name: "新增/编辑资源", ParentId: parents[7].ID, Url: "/resource", Method: "POST"},
		{Name: "删除资源", ParentId: parents[7].ID, Url: "/resource", Method: "DELETE"},
		{Name: "资源选项列表(树形)", ParentId: parents[7].ID, Url: "/resource/option", Method: "GET"},
		{Name: "修改资源匿名访问", ParentId: parents[7].ID, Url: "/resource/anonymous", Method: "PUT"},
		// 评论模块
		{Name: "评论列表", ParentId: parents[8].ID, Url: "/comment/list", Method: "GET"},
		{Name: "删除评论", ParentId: parents[8].ID, Url: "/comment", Method: "DELETE"},
		{Name: "修改评论审核", ParentId: parents[8].ID, Url: "/comment/review", Method: "PUT"},
		// 留言模块
		{Name: "留言列表", ParentId: parents[9].ID, Url: "/message/list", Method: "GET"},
		{Name: "删除留言", ParentId: parents[9].ID, Url: "/message", Method: "DELETE"},
		{Name: "修改留言审核", ParentId: parents[9].ID, Url: "/message/review", Method: "PUT"},
		// 文件模块
		{Name: "文件上传", ParentId: parents[10].ID, Url: "/upload", Method: "POST"},
		// 博客信息模块
		{Name: "获取博客设置", ParentId: parents[11].ID, Url: "/setting/blog-config", Method: "GET"},
		{Name: "获取关于我", ParentId: parents[11].ID, Url: "/setting/about", Method: "GET"},
		{Name: "修改博客设置", ParentId: parents[11].ID, Url: "/setting/blog-config", Method: "PUT"},
		{Name: "修改关于我", ParentId: parents[11].ID, Url: "/setting/about", Method: "PUT"},
		{Name: "获取后台首页信息", ParentId: parents[11].ID, Url: "/home", Method: "GET"},
		// 用户信息模块
		{Name: "用户列表", ParentId: parents[12].ID, Url: "/user/list", Method: "GET"},
		{Name: "获取当前用户信息", ParentId: parents[12].ID, Url: "/user/info", Method: "GET"},
		{Name: "修改用户信息", ParentId: parents[12].ID, Url: "/user", Method: "PUT"},
		{Name: "获取在线用户列表", ParentId: parents[12].ID, Url: "/user/online", Method: "GET"},
		{Name: "强制离线用户", ParentId: parents[12].ID, Url: "/user/offline", Method: "DELETE"},
		{Name: "修改当前用户密码", ParentId: parents[12].ID, Url: "/user/current/password", Method: "PUT"},
		{Name: "修改当前用户信息", ParentId: parents[12].ID, Url: "/user/current", Method: "PUT"},
		{Name: "修改用户禁用", ParentId: parents[12].ID, Url: "/user/disable", Method: "PUT"},
		// 操作日志模块
		{Name: "日志列表", ParentId: parents[13].ID, Url: "/operation/log/list", Method: "GET"},
		{Name: "删除操作日志", ParentId: parents[13].ID, Url: "/operation/log", Method: "DELETE"},
	}

	for i := range resources {
		if err := db.Create(&resources[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(resources[i].Name + " 资源已经存在")
			} else {
				slog.Error(resources[i].Name + " 资源初始化失败" + err.Error())
			}
		}
	}

	// 加载所有资源
	db.Find(&resources)

	// 给 admin 角色添加所有资源访问权限
	var adminRole model.Role
	if err := db.Where("name", "admin").First(&adminRole).Error; err == nil {
		for _, resource := range resources {
			if resource.ID != 0 {
				if err := db.Create(&model.RoleResource{RoleId: adminRole.ID, ResourceId: resource.ID}).Error; err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
						slog.Info("admin 角色菜单关联关系初始化失败" + err.Error())
					} else {
						slog.Error("admin 角色菜单关联关系初始化失败" + err.Error())
					}
				}
			}
		}
	}

	// 给 guest 添加查询资源访问权限
	var guestRole model.Role
	if err := db.Where("name", "guest").First(&guestRole).Error; err == nil {
		for _, resource := range resources {
			if resource.ID != 0 && resource.Method == "GET" {
				if err := db.Create(&model.RoleResource{RoleId: guestRole.ID, ResourceId: resource.ID}).Error; err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
						slog.Info("guest 角色菜单关联关系初始化失败" + err.Error())
					} else {
						slog.Error("guest 角色菜单关联关系初始化失败" + err.Error())
					}
				}
			}
		}
	}

	slog.Info("-----初始化接口资源 end-----")
}

// 生成默认的菜单
func generateDefaultMenus(db *gorm.DB) {
	slog.Info("-----初始化菜单 start-----")

	parents := []model.Menu{
		{Name: "首页", Path: "/home", Icon: "ic:sharp-home", OrderNum: 0, Component: "/home", Redirect: "/home", Catalogue: true}, // catalogue
		{Name: "文章管理", Path: "/article", Icon: "ic:twotone-article", OrderNum: 1, Component: "Layout", Redirect: "/article/list"},
		{Name: "权限管理", Path: "/auth", Icon: "cib:adguard", OrderNum: 3, Component: "Layout", Redirect: "/auth/menu"},
		{Name: "消息管理", Path: "/message", Icon: "ic:twotone-email", OrderNum: 2, Component: "Layout", Redirect: "/message/comment"},
		{Name: "用户管理", Path: "/user", Icon: "ph:user-list-bold", OrderNum: 4, Component: "Layout", Redirect: "/user/list"},
		{Name: "日志管理", Path: "/log", Icon: "material-symbols:receipt-long-outline-rounded", OrderNum: 6, Component: "Layout", Redirect: "/log/operation"},
		{Name: "系统管理", Path: "/setting", Icon: "ion:md-settings", OrderNum: 5, Component: "Layout", Redirect: "/setting/website"},
		{Name: "个人中心", Path: "/profile", Icon: "mdi:account", OrderNum: 7, Component: "/profile", Redirect: "/profile", Catalogue: true}, // catalogue
	}

	for i := range parents {
		if err := db.Create(&parents[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(parents[i].Name + " 菜单已经存在")
			} else {
				slog.Error(parents[i].Name + " 菜单初始化失败" + err.Error())
			}
		}
	}

	menus := []model.Menu{
		// 文章管理
		{Name: "发布文章", Path: "write", Component: "/article/write", Icon: "icon-park-outline:write", OrderNum: 1, ParentId: parents[1].ID},
		{Name: "文章列表", Path: "list", Component: "/article/list", Icon: "material-symbols:format-list-bulleted", OrderNum: 2, ParentId: parents[1].ID},
		{Name: "分类管理", Path: "category", Component: "/article/category", Icon: "tabler:category", OrderNum: 3, ParentId: parents[1].ID},
		{Name: "标签管理", Path: "tag", Component: "/article/tag", Icon: "tabler:tag", OrderNum: 4, ParentId: parents[1].ID},
		{Name: "修改文章", Path: "write/:id", Component: "/article/write", Icon: "icon-park-outline:write", OrderNum: 1, ParentId: parents[1].ID, Hidden: true},
		// 权限管理
		{Name: "菜单管理", Path: "menu", Component: "/auth/menu", Icon: "ic:twotone-menu-book", OrderNum: 1, ParentId: parents[2].ID},
		{Name: "接口管理", Path: "resource", Component: "/auth/resource", Icon: "mdi:api", OrderNum: 2, ParentId: parents[2].ID},
		{Name: "角色管理", Path: "role", Component: "/auth/role", Icon: "carbon:user-role", OrderNum: 3, ParentId: parents[2].ID},
		// 消息管理
		{Name: "评论管理", Path: "comment", Component: "/message/comment", Icon: "ic:twotone-comment", OrderNum: 1, ParentId: parents[3].ID},
		{Name: "留言管理", Path: "leave-msg", Component: "/message/leave-msg", Icon: "ic:twotone-message", OrderNum: 2, ParentId: parents[3].ID},
		// 用户管理
		{Name: "用户列表", Path: "list", Component: "/user/list", Icon: "mdi:account", OrderNum: 1, ParentId: parents[4].ID},
		{Name: "在线用户", Path: "online", Component: "/user/online", Icon: "ic:outline-online-prediction", OrderNum: 2, ParentId: parents[4].ID},
		// 日志管理
		{Name: "操作日志", Path: "operation", Component: "/log/operation", Icon: "mdi:book-open-page-variant-outline", OrderNum: 1, ParentId: parents[5].ID},
		{Name: "登录日志", Path: "login", Component: "/log/login", Icon: "material-symbols:login", OrderNum: 2, ParentId: parents[5].ID},
		// 系统管理
		{Name: "网站管理", Path: "website", Component: "/setting/website", Icon: "el:website", OrderNum: 1, ParentId: parents[6].ID},
		{Name: "页面管理", Path: "page", Component: "/setting/page", Icon: "iconoir:journal-page", OrderNum: 2, ParentId: parents[6].ID},
		{Name: "友链管理", Path: "link", Component: "/setting/link", Icon: "mdi:telegram", OrderNum: 3, ParentId: parents[6].ID},
		{Name: "关于我", Path: "about", Component: "/setting/about", Icon: "cib:about-me", OrderNum: 4, ParentId: parents[6].ID},
	}

	for i := range menus {
		if err := db.Create(&menus[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(menus[i].Name + " 菜单已经存在")
			} else {
				slog.Error(menus[i].Name + " 菜单初始化失败" + err.Error())
			}
		}
	}

	// 加载所有菜单
	db.Find(&menus)

	// 给 admin 角色添加所有菜单访问权限
	var adminRole model.Role
	if err := db.Where("name", "admin").First(&adminRole).Error; err == nil {
		for _, menu := range menus {
			if menu.ID != 0 {
				if err := db.Create(&model.RoleMenu{RoleId: adminRole.ID, MenuId: menu.ID}).Error; err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
						slog.Info("admin 角色菜单关联关系初始化失败" + err.Error())
					} else {
						slog.Error("admin 角色菜单关联关系初始化失败" + err.Error())
					}
				}
			}
		}
	}

	// 给 guest 角色添加所有菜单访问权限
	var guestRole model.Role
	if err := db.Where("name", "guest").First(&guestRole).Error; err == nil {
		for _, menu := range menus {
			if menu.ID != 0 {
				if err := db.Create(&model.RoleMenu{RoleId: guestRole.ID, MenuId: menu.ID}).Error; err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
						slog.Info("guest 角色菜单关联关系初始化失败" + err.Error())
					} else {
						slog.Error("guest 角色菜单关联关系初始化失败" + err.Error())
					}
				}
			}
		}
	}

	slog.Info("-----初始化菜单 end-----")
}
