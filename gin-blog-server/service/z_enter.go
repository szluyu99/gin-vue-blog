package service

import (
	"gin-blog/dao"
)

// redis key

const (
	KEY_CODE         = "code:"        // 验证码
	KEY_USER         = "user:"        // 记录用户
	KEY_DELETE       = "delete:"      //? 记录强制下线用户?
	KEY_ABOUT        = "about"        // 关于我信息
	KEY_BLOG_CONFIG  = "blog_config"  // 博客配置信息
	KEY_VISITOR_AREA = "visitor_area" // 地域统计
	KEY_VIEW_COUNT   = "view_count"   // 访问数量

	KEY_UNIQUE_VISITOR_SET = "unique_visitor" // 唯一用户记录 set

	KEY_ARTICLE_USER_LIKE_SET = "article_user_like:" // 文章点赞 Set
	KEY_ARTICLE_LIKE_COUNT    = "article_like_count" // 文章点赞数
	KEY_ARTICLE_VIEW_COUNT    = "article_view_count" // 文章查看数

	KEY_COMMENT_USER_LIKE_SET = "comment_user_like:" // 评论点赞 Set
	KEY_COMMENT_LIKE_COUNT    = "comment_like_count" // 评论点赞数

	KEY_PAGE = "page" // 页面封面
)

var (
	tagDao          dao.Tag
	categoryDao     dao.Category
	messageDao      dao.Message
	articleDao      dao.Article
	userDao         dao.User
	commentDao      dao.Comment
	friendLinkDao   dao.FriendLink
	roleDao         dao.Role
	resourceDao     dao.Resource
	menuDao         dao.Menu
	operationLogDao dao.OperationLog
)

var (
	blogInfoService BlogInfo
)
