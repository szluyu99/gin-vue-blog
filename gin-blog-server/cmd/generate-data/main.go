package main

import (
	"errors"
	"flag"
	"fmt"
	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

// 种子数据里的图片: 原来指向 cdn.hahacode.cn, 该图床已经失效(连接超时),
// 换成本仓库 images/ 目录下的图片, 由 GitHub 直接提供
const imgBase = "https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images"

func main() {
	configPath := flag.String("c", "../../config.yml", "配置文件路径")
	typeName := flag.String("t", "all", "要初始化的数据类型: config | auth | page | all")
	// 从 cmd/generate-data 目录运行时, sqlite 文件在上一级(与 server 的工作目录 cmd/ 一致);
	// 容器里二进制和配置文件同级, 需要显式关掉
	sqliteParent := flag.Bool("sqlite-parent", true, "sqlite 数据库文件是否在上级目录")
	flag.Parse()

	// 根据命令行参数读取配置文件, 其他变量的初始化依赖于配置文件对象
	conf := g.ReadConfig(*configPath)

	if *sqliteParent {
		conf.SQLite.Dsn = "../" + conf.SQLite.Dsn
	}
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
		{Name: "首页", Label: "home", Cover: imgBase + "/page/home.jpg"},
		{Name: "归档", Label: "archive", Cover: imgBase + "/page/archive.png"},
		{Name: "分类", Label: "category", Cover: imgBase + "/page/category.png"},
		{Name: "标签", Label: "tag", Cover: imgBase + "/page/tag.png"},
		{Name: "友链", Label: "link", Cover: imgBase + "/page/link.jpg"},
		{Name: "关于", Label: "about", Cover: imgBase + "/page/about.jpg"},
		{Name: "留言", Label: "message", Cover: imgBase + "/page/message.jpeg"},
		{Name: "个人中心", Label: "user", Cover: imgBase + "/page/user.jpg"},
		{Name: "相册", Label: "album", Cover: imgBase + "/page/album.png"},
		{Name: "错误页面", Label: "404", Cover: imgBase + "/page/404.jpg"},
		{Name: "文章列表", Label: "article_list", Cover: imgBase + "/page/article_list.jpg"},
	}

	for _, page := range pages {
		if err := db.Create(&page).Error; err != nil {
			if isDuplicate(err) {
				slog.Debug(page.Name + " 页面数据已经存在")
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
		{Key: "website_avatar", Value: imgBase + "/common/header.jpeg", Desc: "网站头像"},
		{Key: "website_name", Value: "阵雨的个人博客", Desc: "网站名称"},
		{Key: "website_author", Value: "阵雨", Desc: "网站作者"},
		{Key: "website_intro", Value: "往事随风而去", Desc: "网站介绍"},
		{Key: "website_notice", Value: "欢迎来到阵雨的个人博客，项目还在开发中...", Desc: "网站公告"},
		{Key: "website_createtime", Value: time.Now().Format(time.DateTime), Desc: "网站创建日期"},
		{Key: "website_record", Value: "粤ICP备2021032312号", Desc: "网站备案号"},
		{Key: "qq", Value: "123456789", Desc: "QQ"},
		{Key: "github", Value: "https://github.com/szluyu99", Desc: "github"},
		{Key: "gitee", Value: "https://gitee.com/szluyu99", Desc: "gitee"},
		{Key: "tourist_avatar", Value: imgBase + "/config/tourist_avatar.jpeg", Desc: "默认游客头像"},
		{Key: "user_avatar", Value: imgBase + "/config/user_avatar.jpeg", Desc: "默认用户头像"},
		{Key: "article_cover", Value: imgBase + "/config/default_article_cover.png", Desc: "默认文章封面"},
		{Key: "is_comment_review", Value: "true", Desc: "评论默认审核"},
		{Key: "is_message_review", Value: "true", Desc: "留言默认审核"},
	}

	for _, config := range configs {
		if err := db.Create(&config).Error; err != nil {
			if isDuplicate(err) {
				slog.Debug(config.Key + " 配置已经存在")
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
			if isDuplicate(err) {
				slog.Debug(roles[i].Name + " 角色已经存在")
				// 取回已有 ID, 否则后面的关联关系会写成 0
				db.Where("name", roles[i].Name).First(&roles[i])
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
				Avatar:   imgBase + "/config/user_avatar.jpeg",
			},
		},
		{
			Username: "guest",
			Password: pwd,
			UserInfo: &model.UserInfo{
				Nickname: "guest",
				Avatar:   imgBase + "/config/user_avatar.jpeg",
			},
		},
	}

	for i := range auths {
		if err := db.Create(&auths[i]).Error; err != nil {
			if isDuplicate(err) {
				slog.Debug(auths[i].Username + " 用户已经存在")
				// 取回已有 ID, 否则下面会插入 user_auth_id = 0 的脏数据
				db.Where("username", auths[i].Username).First(&auths[i])
			} else {
				slog.Error(auths[i].Username + " 用户初始化失败" + err.Error())
			}
		}
		// 创建用户角色关联关系
		if auths[i].ID != 0 && roles[i].ID != 0 {
			db.Create(&model.UserAuthRole{UserAuthId: auths[i].ID, RoleId: roles[i].ID})
		}
	}

	slog.Info("-----初始化默认角色用户 end-----")
}

// 生成默认的接口资源
//
// 资源定义见 internal/model/seed_resource.go, 与后台路由一一对应。
// 这里是对账而不是只增不删: 接口下线或改名后, 旧资源如果留在表里会继续
// 挂在角色上, 之后路由被复用时权限就凭空对上了。
func generateDefaultResources(db *gorm.DB) {
	slog.Info("-----初始化接口资源 start-----")

	// 期望存在的接口资源, key 为 "METHOD URL"
	wanted := make(map[string]model.Resource)
	// 期望存在的模块(父资源)名称
	moduleNames := make(map[string]bool)

	for _, module := range model.AdminResources {
		moduleNames[module.Name] = true

		parent := model.Resource{Name: module.Name}
		err := db.Where("name", module.Name).First(&parent).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&parent).Error; err != nil {
				slog.Error(module.Name + " 资源初始化失败" + err.Error())
				continue
			}
		} else if err != nil {
			slog.Error(module.Name + " 资源查询失败" + err.Error())
			continue
		}

		for _, item := range module.Items {
			wanted[item.Method+" "+item.Url] = model.Resource{
				Name:     item.Name,
				ParentId: parent.ID,
				Url:      item.Url,
				Method:   item.Method,
			}
		}
	}

	// 新增或更新: 以 url + method 定位, 名称和所属模块允许变更
	for _, want := range wanted {
		var exist model.Resource
		err := db.Where(&model.Resource{Url: want.Url, Method: want.Method}).First(&exist).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&want).Error; err != nil {
				slog.Error(want.Name + " 资源初始化失败" + err.Error())
			}
			continue
		}
		if err != nil {
			slog.Error(want.Name + " 资源查询失败" + err.Error())
			continue
		}
		if exist.Name != want.Name || exist.ParentId != want.ParentId {
			if err := db.Model(&exist).
				Updates(map[string]any{"name": want.Name, "parent_id": want.ParentId}).Error; err != nil {
				slog.Error(want.Name + " 资源更新失败" + err.Error())
			}
		}
	}

	// 清理代码中已经不存在的资源, 连同它的角色关联
	// 注意只删资源本身, 不动管理员在后台手动给角色配的其他权限
	var all []model.Resource
	if err := db.Find(&all).Error; err != nil {
		slog.Error("资源列表查询失败" + err.Error())
		return
	}

	var stale []int
	for _, r := range all {
		if r.Url == "" && r.Method == "" { // 模块(父资源)
			if !moduleNames[r.Name] {
				stale = append(stale, r.ID)
			}
			continue
		}
		if _, ok := wanted[r.Method+" "+r.Url]; !ok {
			stale = append(stale, r.ID)
		}
	}

	if len(stale) > 0 {
		if err := db.Delete(&model.RoleResource{}, "resource_id in ?", stale).Error; err != nil {
			slog.Error("清理过期资源的角色关联失败" + err.Error())
		}
		if err := db.Delete(&model.Resource{}, "id in ?", stale).Error; err != nil {
			slog.Error("清理过期资源失败" + err.Error())
		}
		slog.Info(fmt.Sprintf("清理了 %d 条代码中已不存在的资源", len(stale)))
	}

	// 重新加载, 下面按最新的资源表建立角色关联
	if err := db.Find(&all).Error; err != nil {
		slog.Error("资源列表查询失败" + err.Error())
		return
	}

	// 给 admin 角色添加所有资源访问权限
	var adminRole model.Role
	if err := db.Where("name", "admin").First(&adminRole).Error; err == nil {
		bindRoleResources(db, adminRole, all, func(model.Resource) bool { return true })
	}

	// 给 guest 添加查询资源访问权限
	var guestRole model.Role
	if err := db.Where("name", "guest").First(&guestRole).Error; err == nil {
		bindRoleResources(db, guestRole, all, func(r model.Resource) bool { return r.Method == "GET" })
	}

	slog.Info("-----初始化接口资源 end-----")
}

// 把满足 match 的资源挂到角色下, 已经存在的关联跳过
func bindRoleResources(db *gorm.DB, role model.Role, resources []model.Resource, match func(model.Resource) bool) {
	for _, resource := range resources {
		if resource.ID == 0 || !match(resource) {
			continue
		}
		err := db.Create(&model.RoleResource{RoleId: role.ID, ResourceId: resource.ID}).Error
		if err != nil && !isDuplicate(err) {
			slog.Error(role.Name + " 角色资源关联关系初始化失败" + err.Error())
		}
	}
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
			if isDuplicate(err) {
				slog.Debug(parents[i].Name + " 菜单已经存在")
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
			if isDuplicate(err) {
				slog.Debug(menus[i].Name + " 菜单已经存在")
			} else {
				slog.Error(menus[i].Name + " 菜单初始化失败" + err.Error())
			}
		}
	}

	// 加载所有菜单
	db.Find(&menus)

	// 给 admin 和 guest 角色添加所有菜单访问权限
	for _, name := range []string{"admin", "guest"} {
		var role model.Role
		if err := db.Where("name", name).First(&role).Error; err != nil {
			continue
		}
		bindRoleMenus(db, role, menus)
	}

	slog.Info("-----初始化菜单 end-----")
}

// 把菜单挂到角色下, 已经存在的关联跳过
// 重复执行是正常的(每次容器启动都会跑一遍), 不要为此刷一堆日志
func bindRoleMenus(db *gorm.DB, role model.Role, menus []model.Menu) {
	for _, menu := range menus {
		if menu.ID == 0 {
			continue
		}
		err := db.Create(&model.RoleMenu{RoleId: role.ID, MenuId: menu.ID}).Error
		if err != nil && !isDuplicate(err) {
			slog.Error(role.Name + " 角色菜单关联关系初始化失败" + err.Error())
		}
	}
}

// sqlite 和 MySQL 的唯一约束冲突错误文案不同, 统一判断
// 种子数据允许重复执行, 冲突说明已经初始化过, 不算失败
func isDuplicate(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "Duplicate entry")
}
