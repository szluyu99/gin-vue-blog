package routes

import (
	"gin-blog/config"
	"gin-blog/routes/middleware"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 后台管理页面的接口路由
func BackRouter() http.Handler {
	gin.SetMode(config.Cfg.Server.AppMode)

	r := gin.New()
	r.SetTrustedProxies([]string{"*"})

	// 使用本地文件上传, 需要静态文件服务, 使用七牛云不需要
	if config.Cfg.Upload.OssType == "local" {
		r.Static("/public", "./public")
		r.StaticFS("/dir", http.Dir("./public")) // 将 public 目录内的文件列举展示
	}

	r.Use(middleware.Logger())             // 自定义的 zap 日志中间件
	r.Use(middleware.ErrorRecovery(false)) // 自定义错误处理中间件
	r.Use(middleware.Cors())               // 跨域中间件

	// 初始化 session store, session 中用来传递用户的详细信息
	// ! Session 如果使用 Redis 存, 可以存进去, 但是获取不到值?
	// store, _ := redis.NewStoreWithDB(10,
	// 	"tcp",
	// 	config.Cfg.Redis.Addr,
	// 	config.Cfg.Redis.Password,
	// 	strconv.Itoa(config.Cfg.Redis.DB),
	// 	[]byte(config.Cfg.Session.Salt))

	// 基于 cookie 存储 session
	store := cookie.NewStore([]byte(config.Cfg.Session.Salt))

	// session 存储时间跟 JWT 过期时间一致
	store.Options(sessions.Options{MaxAge: int(config.Cfg.JWT.Expire) * 3600})
	r.Use(sessions.Sessions(config.Cfg.Session.Name, store)) // Session 中间件

	// 无需鉴权的接口
	base := r.Group("/api")
	{
		// TODO: 用户注册 和 后台登录 应该记录到 日志
		base.POST("/login", userAuthAPI.Login)   // 后台登录
		base.POST("/report", blogInfoAPI.Report) // 上报信息
	}

	// 需要鉴权的接口
	auth := base.Group("") // "/admin"
	// !注意使用中间件的顺序
	auth.Use(middleware.JWTAuth())      // JWT 鉴权中间件
	auth.Use(middleware.RBAC())         // casbin 权限中间件
	auth.Use(middleware.ListenOnline()) // 监听在线用户
	auth.Use(middleware.OperationLog()) // 记录操作日志
	{
		auth.GET("/home", blogInfoAPI.GetHomeInfo) // 后台首页信息
		auth.GET("/logout", userAuthAPI.Logout)    // 退出登录
		auth.POST("/upload", uploadAPI.UploadFile) // 文件上传

		// 博客设置
		setting := auth.Group("/setting")
		{
			setting.GET("/blog-config", blogInfoAPI.GetBlogConfig)    // 获取博客设置
			setting.PUT("/blog-config", blogInfoAPI.UpdateBlogConfig) // 编辑博客设置
			setting.GET("/about", blogInfoAPI.GetAbout)               // 获取关于我
			setting.PUT("/about", blogInfoAPI.UpdateAbout)            // 编辑关于我
		}
		// 用户模块
		user := auth.Group("/user")
		{
			user.GET("/list", userAPI.GetList)                           // 用户列表
			user.PUT("", userAPI.Update)                                 // 更新用户信息
			user.PUT("/disable", userAPI.UpdateDisable)                  // 修改用户禁用状态
			user.PUT("/password", userAPI.UpdatePassword)                // 修改普通用户密码
			user.PUT("/current/password", userAPI.UpdateCurrentPassword) // 修改管理员密码
			user.GET("/info", userAPI.GetInfo)                           // 获取当前用户信息
			user.PUT("/current", userAPI.UpdateCurrent)                  // 修改当前用户信息
			user.GET("/online", userAPI.GetOnlineList)                   // 获取在线用户
			user.DELETE("/offline", userAPI.ForceOffline)                // 强制用户下线
		}
		// 分类模块
		category := auth.Group("/category")
		{
			category.GET("/list", categoryAPI.GetList)     // 分类列表
			category.POST("", categoryAPI.SaveOrUpdate)    // 新增/编辑分类
			category.DELETE("", categoryAPI.Delete)        // 删除分类
			category.GET("/option", categoryAPI.GetOption) // 分类选项列表
		}
		// 标签模块
		tag := auth.Group("/tag")
		{
			tag.GET("/list", tagAPI.GetList)     // 标签列表
			tag.POST("", tagAPI.SaveOrUpdate)    // 新增/编辑标签
			tag.DELETE("", tagAPI.Delete)        // 删除标签
			tag.GET("/option", tagAPI.GetOption) // 标签选项列表
		}
		// 文章模块
		articles := auth.Group("/article")
		{
			articles.GET("/list", articleAPI.GetList)           // 文章列表
			articles.POST("", articleAPI.SaveOrUpdate)          // 新增/编辑文章
			articles.PUT("/top", articleAPI.UpdateTop)          // 更新文章置顶
			articles.GET("/:id", articleAPI.GetInfo)            // 文章详情
			articles.PUT("/soft-delete", articleAPI.SoftDelete) // 软删除文章
			articles.DELETE("", articleAPI.Delete)              // 物理删除文章
			articles.POST("/export", articleAPI.Export)         // 导出文章
			articles.POST("/import", articleAPI.Import)         // 导入文章
		}
		// 评论模块
		comment := auth.Group("/comment")
		{
			comment.GET("/list", commentAPI.GetList)        // 评论列表
			comment.DELETE("", commentAPI.Delete)           // 删除评论
			comment.PUT("/review", commentAPI.UpdateReview) // 修改评论审核
		}
		// 留言模块
		message := auth.Group("/message")
		{
			message.GET("/list", messageAPI.GetList)        // 留言列表
			message.DELETE("", messageAPI.Delete)           // 删除留言
			message.PUT("/review", messageAPI.UpdateReview) // 审核留言
		}
		// 友情链接
		link := auth.Group("/link")
		{
			link.GET("/list", linkAPI.GetList)  // 友链列表
			link.POST("", linkAPI.SaveOrUpdate) // 新增/编辑友链
			link.DELETE("", linkAPI.Delete)     // 删除友链
		}
		// 资源模块
		resource := auth.Group("/resource")
		{
			resource.GET("/list", resourceAPI.GetTreeList)          // 资源列表(树形)
			resource.POST("", resourceAPI.SaveOrUpdate)             // 新增/编辑资源
			resource.DELETE("/:id", resourceAPI.Delete)             // 删除资源
			resource.PUT("/anonymous", resourceAPI.UpdateAnonymous) // 修改资源匿名访问
			resource.GET("/option", resourceAPI.GetOption)          // 资源选项列表(树形)
		}
		// 菜单模块
		menu := auth.Group("/menu")
		{
			menu.GET("/list", menuAPI.GetTreeList)      // 菜单列表
			menu.POST("", menuAPI.SaveOrUpdate)         // 新增/编辑菜单
			menu.DELETE("/:id", menuAPI.Delete)         // 删除菜单
			menu.GET("/user/list", menuAPI.GetUserMenu) // 获取当前用户的菜单
			menu.GET("/option", menuAPI.GetOption)      // 菜单选项列表(树形)
		}
		// 角色模块
		role := auth.Group("/role")
		{
			role.GET("/list", roleAPI.GetTreeList) // 角色列表(树形)
			role.POST("", roleAPI.SaveOrUpdate)    // 新增/编辑菜单
			role.DELETE("", roleAPI.Delete)        // 删除角色
			role.GET("/option", roleAPI.GetOption) // 角色选项列表(树形)
		}
		// 操作日志模块
		operationLog := auth.Group("/operation/log")
		{
			operationLog.GET("/list", operationLogAPI.GetList) // 操作日志列表
			operationLog.DELETE("", operationLogAPI.Delete)    // 删除操作日志
		}
		// 页面模块
		page := auth.Group("/page")
		{
			page.GET("/list", pageAPI.GetList)  // 页面列表
			page.POST("", pageAPI.SaveOrUpdate) // 新增/编辑页面
			page.DELETE("", pageAPI.Delete)     // 删除页面
		}
	}
	return r
}
