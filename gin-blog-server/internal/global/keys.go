package g

// Redis Key

const (
	// MAIL_CODE = "mail_code:" // 验证码
	// DELETE       = "delete:"      //? 记录强制下线用户?
	ONLINE_USER  = "online_user:"  // 在线用户
	OFFLINE_USER = "offline_user:" // 强制下线用户
	LOGIN_FAIL   = "login_fail:"   // 登录连续失败次数 (login_fail:<用户名+IP 的哈希>), 带 TTL
	VISITOR_AREA = "visitor_area"  // 地域统计
	VIEW_COUNT   = "view_count"    // 访问数量

	/*
		按天的访问量 (view_count:2006-01-02), 带 30 天 TTL

		和 VIEW_COUNT 不是一个口径: 那个是历史独立访客累计(同一访客只算一次),
		这个是当天的访问次数(每次打开前台算一次), 所以按天的加起来不等于累计值。
		只用来画趋势图, 过期即丢 —— 数据库里没有对应字段, 不必长期保留。
	*/
	VIEW_COUNT_DAY = "view_count:"

	// 前端错误上报的按 IP 配额 (error_report:<IP>), 带 TTL。
	// 上报接口匿名可访问, 没有配额等于开了个无鉴权的写库入口
	ERROR_REPORT = "error_report:"

	KEY_UNIQUE_VISITOR_SET = "unique_visitor" // 唯一用户记录 set

	ARTICLE_USER_LIKE_SET = "article_user_like:" // 文章点赞 Set
	ARTICLE_LIKE_COUNT    = "article_like_count" // 文章点赞数
	ARTICLE_VIEW_COUNT    = "article_view_count" // 文章查看数
	// 文章浏览去重, 格式 article_view_visitor:<文章id>:<访客指纹>, 带 TTL
	ARTICLE_VIEW_VISITOR = "article_view_visitor:"

	COMMENT_USER_LIKE_SET = "comment_user_like:" // 评论点赞 Set
	COMMENT_LIKE_COUNT    = "comment_like_count" // 评论点赞数

	PAGE   = "page"   // 页面封面
	CONFIG = "config" // 博客配置
)

// Gin Context Key | Session Key

const (
	CTX_DB        = "_db_field"
	CTX_RDB       = "_rdb_field"
	CTX_USER_AUTH = "_user_auth_field"
)

// Config Key

const (
	CONFIG_ARTICLE_COVER     = "article_cover"
	CONFIG_IS_COMMENT_REVIEW = "is_comment_review"
	CONFIG_IS_MESSAGE_REVIEW = "is_message_review"
	CONFIG_ABOUT             = "about"
)
