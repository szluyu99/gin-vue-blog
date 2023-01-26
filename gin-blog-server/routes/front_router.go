package routes

import (
	"gin-blog/config"
	"gin-blog/routes/middleware"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 前台页面接口路由
func FrontRouter() http.Handler {
	gin.SetMode(config.Cfg.Server.AppMode)

	r := gin.New()
	r.SetTrustedProxies([]string{"*"})

	// 使用本地文件上传, 需要静态文件服务, 使用七牛云不需要
	// if config.Cfg.Upload.OssType == "local" {
	// 	r.Static("/public", "./public")
	// 	r.StaticFS("/dir", http.Dir("./public")) // 将 public 目录内的文件列举展示
	// }

	// 开发模式同时把日志写到控制台
	// if config.Cfg.Server.AppMode == "debug" {
	// 	r.Use(gin.Logger()) // gin 默认日志挺好看的
	// }
	r.Use(middleware.Logger())             // 自定义的 zap 日志中间件
	r.Use(middleware.ErrorRecovery(false)) // 自定义错误处理中间件
	r.Use(middleware.Cors())               // 跨域中间件

	// ! Session 如果使用 Redis 存, 可以存进去, 但是获取不到值?
	// store, _ := redis.NewStoreWithDB(10,
	// 	"tcp",
	// 	config.Cfg.Redis.Addr,
	// 	config.Cfg.Redis.Password,
	// 	strconv.Itoa(config.Cfg.Redis.DB),
	// 	[]byte(config.Cfg.Session.Salt))

	// 基于 cookies 存储 session
	store := cookie.NewStore([]byte(config.Cfg.Session.Salt))

	store.Options(sessions.Options{MaxAge: config.Cfg.Session.MaxAge})
	r.Use(sessions.Sessions(config.Cfg.Session.Name, store)) // Session 中间件 (使用 Redis 存储引擎)

	// 无需监权的接口
	base := r.Group("/api/front")
	{
		base.POST("/login", userAuthAPI.Login)           // 登录
		base.GET("/logout", userAuthAPI.Logout)          // 退出登录
		base.POST("/register", userAuthAPI.Register)     // 注册
		base.GET("/code", userAuthAPI.SendCode)          // 验证码
		base.GET("/home", fBlogInfoAPI.GetFrontHomeInfo) // 前台首页
		base.GET("/about", blogInfoAPI.GetAbout)         // 获取关于我

		article := base.Group("/article")
		{
			article.GET("/list", fArticleAPI.GetFrontList)      // 前台文章列表
			article.GET("/:id", fArticleAPI.GetFrontInfo)       // 前台文章详情
			article.GET("/archive", fArticleAPI.GetArchiveList) // 前台文章归档
			article.GET("/search", fArticleAPI.Search)          // 前台文章搜索
		}
		category := base.Group("/category")
		{
			category.GET("/list", fCategoryAPI.GetFrontList) // 前台分类列表
		}
		tag := base.Group("/tag")
		{
			tag.GET("/list", fTagAPI.GetFrontList) // 前台标签列表
		}
		link := base.Group("/link")
		{
			link.GET("/list", fLinkAPI.GetFrontList) // 前台友链列表
		}
		message := base.Group("/message")
		{
			message.GET("/list", fMessageAPI.GetFrontList) // 前台留言列表
		}
		comment := base.Group("/comment")
		{
			comment.GET("/list", fCommentAPI.GetFrontList)                           // 前台评论列表
			comment.GET("/replies/:comment_id", fCommentAPI.GetReplyListByCommentId) // 根据评论 id 查询回复
		}
	}

	// 需要登录才能进行的操作
	base.Use(middleware.JWTAuth())
	// base.Use(middleware.RBAC()) // TODO: 前端不做权限
	// base.Use(middleware.ListenOnline())
	{
		base.POST("/upload", uploadAPI.UploadFile)                  // 文件上传
		base.GET("/user/info", userAPI.GetInfo)                     // 根据 Token 获取用户信息
		base.PUT("/user/info", userAPI.UpdateCurrent)               // 根据 Token 更新当前用户信息
		base.POST("/message", fMessageAPI.Save)                     // 前台新增留言
		base.POST("/comment", fCommentAPI.Save)                     // 前台新增评论
		base.GET("/comment/like/:comment_id", fCommentAPI.SaveLike) // 前台点赞评论
		base.GET("/article/like/:article_id", fArticleAPI.SaveLike) // 前台点赞文章
	}

	return r
}
