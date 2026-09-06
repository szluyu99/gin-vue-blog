package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"html/template"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Front struct{}

// 同一访客对同一篇文章的浏览量计数间隔
const articleViewInterval = time.Hour

// 搜索最多返回多少条
const searchResultLimit = 20

/*
把命中的关键字包上高亮标签

前台是用 v-html 渲染搜索结果的, 而 keyword 直接来自 query string,
所以必须先整体转义再插标签, 否则搜一段 <img onerror=...> 就会在浏览器里执行。
*/
func highlightKeyword(text, keyword string) string {
	escaped := template.HTMLEscapeString(text)
	escapedKeyword := template.HTMLEscapeString(keyword)
	if escapedKeyword == "" {
		return escaped
	}
	return strings.ReplaceAll(escaped, escapedKeyword,
		"<span style='color:#f47466'>"+escapedKeyword+"</span>")
}

type FAddMessageReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Content  string `json:"content" binding:"required"`
	Speed    int    `json:"speed"`
}

// 注意校验用 binding 而不是 validate: 本项目没有注册自定义 validator, Gin 只认 binding tag
type FAddCommentReq struct {
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content" binding:"required"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int    `json:"type" form:"type" binding:"required,min=1,max=3"`
}

type FCommentQuery struct {
	PageQuery
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int    `json:"type" form:"type"`
}

type FArticleQuery struct {
	PageQuery
	CategoryId int `form:"category_id"`
	TagId      int `form:"tag_id"`
}

type ArchiveVO struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Created_at time.Time `json:"created_at"`
}

type ArticleSearchVO struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// @Summary 前台首页信息
// @Description 文章数, 分类数, 标签数, 公告与访问量
// @Tags Front
// @Produce json
// @Success 0 {object} Response[model.FrontHomeVO]
// @Router /front/home [get]
func (*Front) GetHomeInfo(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	data, err := model.GetFrontStatistics(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	data.ViewCount, _ = rdb.Get(rctx, g.VIEW_COUNT).Int64()

	ReturnSuccess(c, data)
}

// @Summary 前台标签列表
// @Description 获取全部标签
// @Tags Front
// @Produce json
// @Success 0 {object} Response[[]model.TagVO]
// @Router /front/tag/list [get]
func (*Front) GetTagList(c *gin.Context) {
	list, _, err := model.GetTagList(GetDB(c), 1, model.PageSizeAll, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// @Summary 前台分类列表
// @Description 获取全部分类
// @Tags Front
// @Produce json
// @Success 0 {object} Response[[]model.CategoryVO]
// @Router /front/category/list [get]
func (*Front) GetCategoryList(c *gin.Context) {
	list, _, err := model.GetCategoryList(GetDB(c), 1, model.PageSizeAll, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// @Summary 前台留言列表
// @Description 只返回审核通过的留言
// @Tags Front
// @Produce json
// @Success 0 {object} Response[[]model.Message]
// @Router /front/message/list [get]
func (*Front) GetMessageList(c *gin.Context) {
	isReview := true
	list, _, err := model.GetMessageList(GetDB(c), 1, model.PageSizeAll, "", &isReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// @Summary 前台友链列表
// @Description 获取全部友链
// @Tags Front
// @Produce json
// @Success 0 {object} Response[[]model.FriendLink]
// @Router /front/link/list [get]
func (*Front) GetLinkList(c *gin.Context) {
	list, _, err := model.GetLinkList(GetDB(c), 1, model.PageSizeAll, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, list)
}

/*
以下接口需要登录
*/

// TODO: 添加自定义头像和昵称留言功能（即可以不登录留言）
// @Summary 新增留言
// @Description 新增留言, 需要登录, 是否需要审核取决于博客配置
// @Tags Front
// @Accept json
// @Produce json
// @Param form body FAddMessageReq true "新增留言"
// @Success 0 {object} Response[model.Message]
// @Security ApiKeyAuth
// @Router /front/message [post]
func (*Front) SaveMessage(c *gin.Context) {
	var req FAddMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	req.Content = template.HTMLEscapeString(req.Content)
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}
	db := GetDB(c)

	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSource(ipAddress)
	// 留言要读留言自己的审核开关, 之前读的是评论的开关, 后台设置里的 is_message_review 根本没人用
	isReview := model.GetConfigBool(db, g.CONFIG_IS_MESSAGE_REVIEW)

	info := auth.UserInfo
	message, err := model.SaveMessage(db, info.Nickname, info.Nickname, req.Content, ipAddress, ipSource, req.Speed, isReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, message)
}

// 保存评论（只能新增，不能编辑）
// TODO: 添加自定义头像和昵称留言功能（即可以不登录评论）
// TODO: 开启邮箱通知用户功能
// TODO: HTMLUtil.Filter 过滤 HTML 元素中的字符串...
// @Summary 新增评论
// @Description 新增评论或回复评论, 需要登录
// @Tags Front
// @Accept json
// @Produce json
// @Param form body FAddCommentReq true "新增评论"
// @Success 0 {object} Response[model.Comment]
// @Security ApiKeyAuth
// @Router /front/comment [post]
func (*Front) SaveComment(c *gin.Context) {
	var req FAddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 过滤评论内容，防止XSS攻击
	req.Content = template.HTMLEscapeString(req.Content)
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}
	db := GetDB(c)
	isReview := model.GetConfigBool(db, g.CONFIG_IS_COMMENT_REVIEW)

	var comment *model.Comment
	var err error

	if req.ReplyUserId == 0 {
		comment, err = model.AddComment(db, auth.ID, req.Type, req.TopicId, req.Content, isReview)
	} else {
		comment, err = model.ReplyComment(db, auth.ID, req.ReplyUserId, req.ParentId, req.Content, isReview)
	}

	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 站内通知: 写失败不影响发评论本身, 只记日志
	// (邮件通知默认关闭, 这是「别人回复了你」唯一的出口)
	if err := model.NotifyOnComment(db, comment); err != nil {
		slog.Warn("写站内通知失败", "err", err, "comment_id", comment.ID)
	}

	ReturnSuccess(c, comment)
}

// @Summary 前台评论列表
// @Description 顶级评论列表, 每条最多带 3 条回复
// @Tags Front
// @Produce json
// @Param topic_id query int false "主题 ID(文章/说说)"
// @Param type query int false "评论类型(1-文章 2-友链 3-说说)"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.CommentVO]]
// @Router /front/comment/list [get]
func (*Front) GetCommentList(c *gin.Context) {
	var query FCommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	data, total, err := model.GetCommentVOList(db, query.Page, query.Size, query.TopicId, query.Type)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 顶层评论和它们的回复一起取点赞数: 以前只取了顶层的,
	// 回复的 like_count 恒为 0, 展开更多之后又会变成真实值
	ids := make([]int, 0, len(data))
	for i := range data {
		if len(data[i].ReplyList) > 3 {
			data[i].ReplyList = data[i].ReplyList[:3] // 只显示 3 条回复
		}
		ids = append(ids, data[i].ID)
		for _, reply := range data[i].ReplyList {
			ids = append(ids, reply.ID)
		}
	}
	likeCountMap := hashCounts(rdb, g.COMMENT_LIKE_COUNT, ids)
	for i := range data {
		data[i].LikeCount = likeCountMap[data[i].ID]
		for j := range data[i].ReplyList {
			data[i].ReplyList[j].LikeCount = likeCountMap[data[i].ReplyList[j].ID]
		}
	}

	ReturnSuccess(c, PageResult[model.CommentVO]{
		List:  data,
		Total: total,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 获取评论的回复列表
// @Description 根据评论 ID 分页查询其回复
// @Tags Front
// @Produce json
// @Param comment_id path int true "评论 ID"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[[]model.CommentVO]
// @Router /front/comment/replies/{comment_id} [get]
func (*Front) GetReplyListByCommentId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	replyList, err := model.GetCommentReplyList(db, id, query.Page, query.Size)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ids := make([]int, 0, len(replyList))
	for _, reply := range replyList {
		ids = append(ids, reply.ID)
	}
	likeCountMap := hashCounts(rdb, g.COMMENT_LIKE_COUNT, ids)

	data := make([]model.CommentVO, 0, len(replyList))
	for _, reply := range replyList {
		data = append(data, model.CommentVO{
			Comment:   reply,
			LikeCount: likeCountMap[reply.ID],
		})
	}

	ReturnSuccess(c, data)
}

// @Summary 点赞评论
// @Description 点赞/取消点赞评论, 需要登录
// @Tags Front
// @Produce json
// @Param comment_id path int true "评论 ID"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /front/comment/like/{comment_id} [get]
func (*Front) LikeComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rdb := GetRDB(c)
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	// 记录某个用户已经对某个评论点过赞
	commentLikeUserKey := g.COMMENT_USER_LIKE_SET + strconv.Itoa(auth.ID)
	// 该评论已经被记录过, 再点赞就是取消点赞
	if rdb.SIsMember(rctx, commentLikeUserKey, id).Val() {
		rdb.SRem(rctx, commentLikeUserKey, id)
		rdb.HIncrBy(rctx, g.COMMENT_LIKE_COUNT, strconv.Itoa(id), -1)
	} else { // 未被记录过, 则是增加点赞
		rdb.SAdd(rctx, commentLikeUserKey, id)
		rdb.HIncrBy(rctx, g.COMMENT_LIKE_COUNT, strconv.Itoa(id), 1)
	}

	ReturnSuccess(c, nil)
}

/*
文章相关接口
*/

// @Summary 前台文章列表
// @Description 只返回公开且不在回收站的文章
// @Tags Front
// @Produce json
// @Param category_id query int false "分类 ID"
// @Param tag_id query int false "标签 ID"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.Article]]
// @Router /front/article/list [get]
func (*Front) GetArticleList(c *gin.Context) {
	var query FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, total, err := model.GetBlogArticleList(GetDB(c), query.Page, query.Size, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 以前把 total 丢掉只返回当页数据, 前台拿不到总数也就做不了分页,
	// 分类/标签下超过一页的文章等于没有入口能看到
	ReturnSuccess(c, PageResult[model.Article]{
		Page:  query.Page,
		Size:  query.Size,
		Total: total,
		List:  list,
	})
}

// @Summary 前台文章详情
// @Description 文章详情, 附带上下篇/推荐/最新文章与点赞浏览评论数, 同一访客一小时内只计一次浏览量
// @Tags Front
// @Produce json
// @Param id path int true "文章 ID"
// @Success 0 {object} Response[model.BlogArticleVO]
// @Router /front/article/{id} [get]
func (*Front) GetArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	// 文章详情
	val, err := model.GetBlogArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	article := model.BlogArticleVO{Article: *val}

	// 推荐文章（6篇）
	article.RecommendArticles, err = model.GetRecommendList(db, id, 6)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 最新文章（5篇）
	article.NewestArticles, err = model.GetNewestList(db, 5)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 更新文章浏览量: 同一访客在 articleViewInterval 内重复访问只算一次,
	// 否则刷新页面就能刷量。SetNX 成功说明是这个窗口里的首次访问。
	viewKey := g.ARTICLE_VIEW_VISITOR + strconv.Itoa(id) + ":" + visitorFingerprint(c)
	first, err := rdb.SetNX(rctx, viewKey, 1, articleViewInterval).Result()
	if err != nil {
		// Redis 异常时宁可多计一次, 也不要让浏览量停止统计
		slog.Warn("文章浏览去重失败, 本次直接计数", "article_id", id, "err", err)
	}
	if first || err != nil {
		rdb.ZIncrBy(rctx, g.ARTICLE_VIEW_COUNT, 1, strconv.Itoa(id))
	}

	// 上一篇文章
	article.LastArticle, err = model.GetLastArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 下一篇文章
	article.NextArticle, err = model.GetNextArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 点赞量, 浏览量
	article.ViewCount = int64(rdb.ZScore(rctx, g.ARTICLE_VIEW_COUNT, strconv.Itoa(id)).Val())
	likeCount, _ := strconv.Atoi(rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(id)).Val())
	article.LikeCount = int64(likeCount)

	// 评论数量
	article.CommentCount, err = model.GetArticleCommentCount(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

// @Summary 前台文章归档
// @Description 按时间归档的文章列表
// @Tags Front
// @Produce json
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[ArchiveVO]]
// @Router /front/article/archive [get]
func (*Front) GetArchiveList(c *gin.Context) {
	var query FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, total, err := model.GetBlogArticleArchiveList(GetDB(c), query.Page, query.Size)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	archives := make([]ArchiveVO, 0)
	for _, article := range list {
		archives = append(archives, ArchiveVO{
			ID:         article.ID,
			Title:      article.Title,
			Created_at: article.CreatedAt,
		})
	}

	ReturnSuccess(c, PageResult[ArchiveVO]{
		Total: total,
		List:  archives,
		Page:  query.Page,
		Size:  query.Size,
	})
}

// @Summary 点赞文章
// @Description 点赞/取消点赞文章, 需要登录
// @Tags Front
// @Produce json
// @Param article_id path int true "文章 ID"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /front/article/like/{article_id} [get]
func (*Front) LikeArticle(c *gin.Context) {
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	articleId, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rdb := GetRDB(c)

	// 记录某个用户已经对某个文章点过赞
	articleLikeUserKey := g.ARTICLE_USER_LIKE_SET + strconv.Itoa(auth.ID)
	// 该文章已经被记录过, 再点赞就是取消点赞
	if rdb.SIsMember(rctx, articleLikeUserKey, articleId).Val() {
		rdb.SRem(rctx, articleLikeUserKey, articleId)
		rdb.HIncrBy(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), -1)
	} else { // 未被记录过, 则是增加点赞
		rdb.SAdd(rctx, articleLikeUserKey, articleId)
		rdb.HIncrBy(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), 1)
	}

	ReturnSuccess(c, nil)
}

// @Summary 搜索文章
// @Description 按关键字搜索标题和内容, 命中处高亮
// @Tags Front
// @Produce json
// @Param keyword query string false "搜索关键字"
// @Success 0 {object} Response[[]ArticleSearchVO]
// @Router /front/article/search [get]
func (*Front) SearchArticle(c *gin.Context) {
	result := make([]ArticleSearchVO, 0)

	keyword := c.Query("keyword")
	if keyword == "" {
		ReturnSuccess(c, result)
		return
	}

	db := GetDB(c)

	articleList, err := model.SearchArticle(db, keyword, searchResultLimit)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	for _, article := range articleList {
		content := article.Content
		// 关键字在内容中的起始位置
		keywordStartIndex := unicodeIndex(content, keyword)
		if keywordStartIndex != -1 { // 关键字在内容中
			preIndex, afterIndex := 0, 0
			if keywordStartIndex > 25 {
				preIndex = keywordStartIndex - 25
			}
			// 防止中文截取出乱码 (中文在 golang 是 3 个字符, 使用 rune 中文占一个数组下标)
			preText := substring(content, preIndex, keywordStartIndex)

			// 关键字在内容中的结束位置
			keywordEndIndex := keywordStartIndex + unicodeLen(keyword)
			afterLength := len(content) - keywordEndIndex
			if afterLength > 175 {
				afterIndex = keywordEndIndex + 175
			} else {
				afterIndex = keywordEndIndex + afterLength
			}
			afterText := substring(content, keywordStartIndex, afterIndex)
			content = preText + afterText
		}

		result = append(result, ArticleSearchVO{
			ID:      article.ID,
			Title:   highlightKeyword(article.Title, keyword),
			Content: highlightKeyword(content, keyword),
		})
	}

	ReturnSuccess(c, result)
}

// 获取带中文的字符串中子字符串的实际位置，非字节位置
func unicodeIndex(str, substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str, substr)
	if result > 0 {
		prefix := []byte(str)[0:result]
		rs := []rune(string(prefix))
		result = len(rs)
	}
	return result
}

// 获取带中文的字符串实际长度，非字节长度
func unicodeLen(str string) int {
	var r = []rune(str)
	return len(r)
}

// 解决中文获取位置不正确问题
func substring(source string, start int, end int) string {
	var unicodeStr = []rune(source)
	length := len(unicodeStr)
	if start >= end {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start <= 0 && end >= length {
		return source
	}
	var substring = ""
	for i := start; i < end; i++ {
		substring += string(unicodeStr[i])
	}
	return substring
}
