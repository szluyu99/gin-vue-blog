package v1

import "gin-blog/service"

var (
	tagService          service.Tag
	categoryService     service.Category
	messageService      service.Message
	articleService      service.Article
	userService         service.User
	commentService      service.Comment
	linkService         service.FriendLink
	roleService         service.Role
	resourceService     service.Resource
	menuService         service.Menu
	blogInfoService     service.BlogInfo
	operationLogService service.OperationLog
	pageService         service.Page
	fileService         service.File // 文件相关: 上传、导入、导出...
)
