package routes

import (
	"gin-blog/api/front"
	v1 "gin-blog/api/v1"
)

// 后台接口
var (
	categoryAPI     v1.Category     // 分类
	tagAPI          v1.Tag          // 标签
	articleAPI      v1.Article      // 文章
	userAPI         v1.User         // 用户
	userAuthAPI     v1.UserAuth     // 用户账号
	commentAPI      v1.Comment      // 评论
	uploadAPI       v1.Upload       // 文件上传
	messageAPI      v1.Message      // 留言
	linkAPI         v1.Link         // 友情链接
	roleAPI         v1.Role         // 角色
	resourceAPI     v1.Resource     // 资源
	menuAPI         v1.Menu         // 菜单
	blogInfoAPI     v1.BlogInfo     // 博客设置
	operationLogAPI v1.OperationLog // 操作日志
	pageAPI         v1.Page

	// store redis.Store // 使用 redis 作为 session 的存储引擎
)

// 前台接口
var (
	fBlogInfoAPI front.BlogInfo // 博客信息
	fArticleAPI  front.Article  // 文章
	fCategoryAPI front.Category // 分类
	fTagAPI      front.Tag      // 标签
	fMessageAPI  front.Message  // 留言
	fCommentAPI  front.Comment  // 评论
	fLinkAPI     front.Link     // 友链
)
