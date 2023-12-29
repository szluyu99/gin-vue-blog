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
	typeName := flag.String("type", "all", "要初始化的数据类型: config | auth | page | resource | all")
	flag.Parse()

	// 根据命令行参数读取配置文件, 其他变量的初始化依赖于配置文件对象
	conf, err := g.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}

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
	case "resource":
		generateDefaultResource(db)
	case "all":
		fallthrough
	default:
		generateDefaultConfigs(db)
		generateDefaultAuths(db)
		generateDefaultPages(db)
		generateDefaultResource(db)
	}
}

// 生成 3 个默认角色及验证信息: admin, user, guest
func generateDefaultAuths(db *gorm.DB) {
	roles := []model.Role{
		{Name: "admin", Label: "管理员"},
		{Name: "user", Label: "用户"},
		{Name: "guest", Label: "游客"},
	}

	for i := range roles {
		if err := db.Create(&roles[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				slog.Info(roles[i].Name + "角色已经存在")
			} else {
				slog.Error(roles[i].Name + "角色初始化失败" + err.Error())
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
			Username: "user",
			Password: pwd,
			UserInfo: &model.UserInfo{
				Nickname: "user",
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
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				slog.Info(auths[i].Username + " 用户已经存在")
			} else {
				slog.Error(auths[i].Username + " 用户初始化失败" + err.Error())
			}
		}
		db.Create(&model.UserAuthRole{UserAuthId: auths[i].ID, RoleId: roles[i].ID})
	}

	slog.Info("-----初始化默认角色用户结束-----")
}

// 生成默认的页面
func generateDefaultPages(db *gorm.DB) {
	pages := []model.Page{
		{Label: "首页", Name: "home", Cover: "https://cdn.hahacode.cn/page/home.jpg"},
		{Label: "归档", Name: "archive", Cover: "https://cdn.hahacode.cn/page/archive.png"},
		{Label: "分类", Name: "category", Cover: "https://cdn.hahacode.cn/page/category.png"},
		{Label: "标签", Name: "tag", Cover: "https://cdn.hahacode.cn/page/tag.png"},
		{Label: "友链", Name: "link", Cover: "https://cdn.hahacode.cn/page/link.png"},
		{Label: "关于", Name: "about", Cover: "https://cdn.hahacode.cn/page/about.png"},
		{Label: "留言", Name: "message", Cover: "https://cdn.hahacode.cn/page/message.png"},
		{Label: "个人中心", Name: "user", Cover: "https://cdn.hahacode.cn/page/user.png"},
		{Label: "相册", Name: "album", Cover: "https://cdn.hahacode.cn/page/album.png"},
		{Label: "错误页面", Name: "404", Cover: "https://cdn.hahacode.cn/page/404.png"},
		{Label: "文章列表", Name: "article_list", Cover: "https://cdn.hahacode.cn/page/article_list.png"},
	}

	for _, page := range pages {
		if err := db.Create(&page).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				slog.Info(page.Name + " 页面数据已经存在")
			} else {
				slog.Error(page.Name + " 页面初始化失败" + err.Error())
			}
		}
	}

	slog.Info("-----初始化博客页面结束-----")
}

// TODO: 生成默认的接口资源
func generateDefaultResource(db *gorm.DB) {
}

// 生成默认配置信息
func generateDefaultConfigs(db *gorm.DB) {
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
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				slog.Info(config.Key + "配置已经存在")
			} else {
				slog.Error(config.Key + "配置初始化失败" + err.Error())
			}
		}
	}

	slog.Info("-----初始化博客配置结束-----")
}
