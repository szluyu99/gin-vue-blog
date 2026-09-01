package model

// 后台接口的资源定义(权限控制的最小单位)
//
// cmd/generate-data 用它初始化资源表; internal 包的测试用它校验
// 与实际注册的后台路由完全一致 —— 资源表缺失的接口不会被 PermissionCheck
// 保护, 曾经因为路由改名而漏掉过 3 个接口。
type ResourceModule struct {
	Name  string         // 父资源(模块)名称
	Items []ResourceItem // 模块下的接口
}

type ResourceItem struct {
	Name   string
	Url    string
	Method string
}

var AdminResources = []ResourceModule{
	{
		Name: "文章模块",
		Items: []ResourceItem{
			{Name: "文章列表", Url: "/article/list", Method: "GET"},
			{Name: "文章详情", Url: "/article/:id", Method: "GET"},
			{Name: "新增/编辑文章", Url: "/article", Method: "POST"},
			{Name: "更新文章软删除", Url: "/article/soft-delete", Method: "PUT"},
			{Name: "删除文章", Url: "/article", Method: "DELETE"},
			{Name: "修改文章置顶", Url: "/article/top", Method: "PUT"},
			{Name: "导出文章", Url: "/article/export", Method: "POST"},
			{Name: "导入文章", Url: "/article/import", Method: "POST"},
		},
	},
	{
		Name: "分类模块",
		Items: []ResourceItem{
			{Name: "分类列表", Url: "/category/list", Method: "GET"},
			{Name: "新增/编辑分类", Url: "/category", Method: "POST"},
			{Name: "删除分类", Url: "/category", Method: "DELETE"},
			{Name: "分类选项列表", Url: "/category/option", Method: "GET"},
		},
	},
	{
		Name: "标签模块",
		Items: []ResourceItem{
			{Name: "标签列表", Url: "/tag/list", Method: "GET"},
			{Name: "新增/编辑标签", Url: "/tag", Method: "POST"},
			{Name: "删除标签", Url: "/tag", Method: "DELETE"},
			{Name: "标签选项列表", Url: "/tag/option", Method: "GET"},
		},
	},
	{
		Name: "页面模块",
		Items: []ResourceItem{
			{Name: "页面列表", Url: "/page/list", Method: "GET"},
			{Name: "新增/编辑页面", Url: "/page", Method: "POST"},
			{Name: "删除页面", Url: "/page", Method: "DELETE"},
		},
	},
	{
		Name: "友链模块",
		Items: []ResourceItem{
			{Name: "友链列表", Url: "/link/list", Method: "GET"},
			{Name: "新增/编辑友链", Url: "/link", Method: "POST"},
			{Name: "删除友链", Url: "/link", Method: "DELETE"},
		},
	},
	{
		Name: "菜单模块",
		Items: []ResourceItem{
			{Name: "菜单列表", Url: "/menu/list", Method: "GET"},
			{Name: "新增/编辑菜单", Url: "/menu", Method: "POST"},
			{Name: "删除菜单", Url: "/menu/:id", Method: "DELETE"},
			{Name: "菜单选项列表(树形)", Url: "/menu/option", Method: "GET"},
			{Name: "获取当前用户菜单", Url: "/menu/user/list", Method: "GET"},
		},
	},
	{
		Name: "角色模块",
		Items: []ResourceItem{
			{Name: "角色列表", Url: "/role/list", Method: "GET"},
			{Name: "新增/编辑角色", Url: "/role", Method: "POST"},
			{Name: "删除角色", Url: "/role", Method: "DELETE"},
			{Name: "角色选项列表", Url: "/role/option", Method: "GET"},
		},
	},
	{
		Name: "资源模块",
		Items: []ResourceItem{
			{Name: "资源列表", Url: "/resource/list", Method: "GET"},
			{Name: "新增/编辑资源", Url: "/resource", Method: "POST"},
			{Name: "删除资源", Url: "/resource/:id", Method: "DELETE"},
			{Name: "资源选项列表(树形)", Url: "/resource/option", Method: "GET"},
			{Name: "修改资源匿名访问", Url: "/resource/anonymous", Method: "PUT"},
		},
	},
	{
		Name: "评论模块",
		Items: []ResourceItem{
			{Name: "评论列表", Url: "/comment/list", Method: "GET"},
			{Name: "删除评论", Url: "/comment", Method: "DELETE"},
			{Name: "修改评论审核", Url: "/comment/review", Method: "PUT"},
		},
	},
	{
		Name: "留言模块",
		Items: []ResourceItem{
			{Name: "留言列表", Url: "/message/list", Method: "GET"},
			{Name: "删除留言", Url: "/message", Method: "DELETE"},
			{Name: "修改留言审核", Url: "/message/review", Method: "PUT"},
		},
	},
	{
		Name: "文件模块",
		Items: []ResourceItem{
			{Name: "文件上传", Url: "/upload", Method: "POST"},
		},
	},
	{
		Name: "博客信息模块",
		Items: []ResourceItem{
			{Name: "获取关于我", Url: "/setting/about", Method: "GET"},
			{Name: "修改关于我", Url: "/setting/about", Method: "PUT"},
			{Name: "获取后台首页信息", Url: "/home", Method: "GET"},
			{Name: "修改博客配置", Url: "/config", Method: "PATCH"},
		},
	},
	{
		Name: "用户信息模块",
		Items: []ResourceItem{
			{Name: "用户列表", Url: "/user/list", Method: "GET"},
			{Name: "获取当前用户信息", Url: "/user/info", Method: "GET"},
			{Name: "修改用户信息", Url: "/user", Method: "PUT"},
			{Name: "获取在线用户列表", Url: "/user/online", Method: "GET"},
			{Name: "强制离线用户", Url: "/user/offline/:id", Method: "POST"},
			{Name: "修改当前用户密码", Url: "/user/current/password", Method: "PUT"},
			{Name: "修改当前用户信息", Url: "/user/current", Method: "PUT"},
			{Name: "修改用户禁用", Url: "/user/disable", Method: "PUT"},
		},
	},
	{
		Name: "操作日志模块",
		Items: []ResourceItem{
			{Name: "日志列表", Url: "/operation/log/list", Method: "GET"},
			{Name: "删除操作日志", Url: "/operation/log", Method: "DELETE"},
		},
	},
}
