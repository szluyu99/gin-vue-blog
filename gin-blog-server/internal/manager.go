package ginblog

import (
	"gin-blog/docs"
	"gin-blog/internal/handle"
	"gin-blog/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	// 后台管理系统接口

	categoryAPI     handle.Category     // 分类
	tagAPI          handle.Tag          // 标签
	articleAPI      handle.Article      // 文章
	userAPI         handle.User         // 用户
	userAuthAPI     handle.UserAuth     // 用户账号
	commentAPI      handle.Comment      // 评论
	uploadAPI       handle.Upload       // 文件上传
	messageAPI      handle.Message      // 留言
	linkAPI         handle.Link         // 友情链接
	roleAPI         handle.Role         // 角色
	resourceAPI     handle.Resource     // 资源
	menuAPI         handle.Menu         // 菜单
	blogInfoAPI     handle.BlogInfo     // 博客设置
	operationLogAPI handle.OperationLog // 操作日志
	pageAPI         handle.Page         // 页面

	// 博客前台接口

	frontAPI handle.Front // 博客前台接口汇总
)

// TODO: 前端修改 PUT 和 PATCH 请求
func RegisterHandlers(r *gin.Engine) {
	// Swagger
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registerBaseHandler(r)
	registerAdminHandler(r)
	registerBlogHandler(r)
}

// 通用接口: 全部不需要 登录 + 鉴权
func registerBaseHandler(r *gin.Engine) {
	base := r.Group("/api")

	// TODO: 登录, 注册 记录日志
	base.POST("/login", userAuthAPI.Login)          // 登录
	base.POST("/register", userAuthAPI.Register)    // 注册
	base.GET("/logout", userAuthAPI.Logout)         // 退出登录
	base.POST("/report", blogInfoAPI.Report)        // 上报信息
	base.GET("/config", blogInfoAPI.GetConfigMap)   // 获取配置
	base.PATCH("/config", blogInfoAPI.UpdateConfig) // 更新配置
	base.GET("/email/verify",userAuthAPI.VerifyCode)
}

// 后台管理系统的接口: 全部需要 登录 + 鉴权
func registerAdminHandler(r *gin.Engine) {
	auth := r.Group("/api")

	// !注意使用中间件的顺序
	auth.Use(middleware.JWTAuth())
	auth.Use(middleware.PermissionCheck())
	auth.Use(middleware.OperationLog())
	auth.Use(middleware.ListenOnline())

	auth.GET("/home", blogInfoAPI.GetHomeInfo) // 后台首页信息
	auth.POST("/upload", uploadAPI.UploadFile) // 文件上传

	// 博客设置
	setting := auth.Group("/setting")
	{
		setting.GET("/about", blogInfoAPI.GetAbout)    // 获取关于我
		setting.PUT("/about", blogInfoAPI.UpdateAbout) // 编辑关于我
	}
	// 用户模块
	user := auth.Group("/user")
	{
		user.GET("/list", userAPI.GetList)          // 用户列表
		user.PUT("", userAPI.Update)                // 更新用户信息
		user.PUT("/disable", userAPI.UpdateDisable) // 修改用户禁用状态
		// user.PUT("/password", userAPI.UpdatePassword)                // 修改普通用户密码
		user.PUT("/current/password", userAPI.UpdateCurrentPassword) // 修改管理员密码
		user.GET("/info", userAPI.GetInfo)                           // 获取当前用户信息
		user.PUT("/current", userAPI.UpdateCurrent)                  // 修改当前用户信息
		user.GET("/online", userAPI.GetOnlineList)                   // 获取在线用户
		user.POST("/offline/:id", userAPI.ForceOffline)              // 强制用户下线
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
		articles.GET("/list", articleAPI.GetList)                 // 文章列表
		articles.POST("", articleAPI.SaveOrUpdate)                // 新增/编辑文章
		articles.PUT("/top", articleAPI.UpdateTop)                // 更新文章置顶
		articles.GET("/:id", articleAPI.GetDetail)                // 文章详情
		articles.PUT("/soft-delete", articleAPI.UpdateSoftDelete) // 软删除文章
		articles.DELETE("", articleAPI.Delete)                    // 物理删除文章
		articles.POST("/export", articleAPI.Export)               // 导出文章
		articles.POST("/import", articleAPI.Import)               // 导入文章
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

// 博客前台的接口: 大部分不需要登录, 部分需要登录
func registerBlogHandler(r *gin.Engine) {
	base := r.Group("/api/front")

	base.GET("/about", blogInfoAPI.GetAbout) // 获取关于我
	base.GET("/home", frontAPI.GetHomeInfo)  // 前台首页
	base.GET("/page", pageAPI.GetList)       // 前台页面

	article := base.Group("/article")
	{
		article.GET("/list", frontAPI.GetArticleList)    // 前台文章列表
		article.GET("/:id", frontAPI.GetArticleInfo)     // 前台文章详情
		article.GET("/archive", frontAPI.GetArchiveList) // 前台文章归档
		article.GET("/search", frontAPI.SearchArticle)   // 前台文章搜索
	}
	category := base.Group("/category")
	{
		category.GET("/list", frontAPI.GetCategoryList) // 前台分类列表
	}
	tag := base.Group("/tag")
	{
		tag.GET("/list", frontAPI.GetTagList) // 前台标签列表
	}
	link := base.Group("/link")
	{
		link.GET("/list", frontAPI.GetLinkList) // 前台友链列表
	}
	message := base.Group("/message")
	{
		message.GET("/list", frontAPI.GetMessageList) // 前台留言列表
	}
	comment := base.Group("/comment")
	{
		comment.GET("/list", frontAPI.GetCommentList)                         // 前台评论列表
		comment.GET("/replies/:comment_id", frontAPI.GetReplyListByCommentId) // 根据评论 id 查询回复
	}

	// 需要登录才能进行的操作
	base.Use(middleware.JWTAuth())
	{
		base.POST("/upload", uploadAPI.UploadFile)    // 文件上传
		base.GET("/user/info", userAPI.GetInfo)       // 根据 Token 获取用户信息
		base.PUT("/user/info", userAPI.UpdateCurrent) // 根据 Token 更新当前用户信息

		base.POST("/message", frontAPI.SaveMessage)                 // 前台新增留言
		base.POST("/comment", frontAPI.SaveComment)                 // 前台新增评论
		base.GET("/comment/like/:comment_id", frontAPI.LikeComment) // 前台点赞评论
		base.GET("/article/like/:article_id", frontAPI.LikeArticle) // 前台点赞文章
	}
}
