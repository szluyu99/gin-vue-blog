package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"io"
	"log/slog"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Article struct{}

type SoftDeleteReq struct {
	Ids      []int `json:"ids"`
	IsDelete bool  `json:"is_delete"`
}

type AddOrEditArticleReq struct {
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Desc        string `json:"desc"`
	Content     string `json:"content" binding:"required"`
	Img         string `json:"img"`
	Type        int    `json:"type" binding:"required,min=1,max=3"`   // 类型: 1-原创 2-转载 3-翻译
	Status      int    `json:"status" binding:"required,min=1,max=3"` // 状态: 1-公开 2-私密 3-评论可见
	IsTop       bool   `json:"is_top"`
	OriginalUrl string `json:"original_url"`

	TagNames     []string `json:"tag_names"`
	CategoryName string   `json:"category_name"`
}

// TODO: 添加对标签数组的查询
type ArticleQuery struct {
	PageQuery
	Title      string `form:"title"`
	CategoryId int    `form:"category_id"`
	TagId      int    `form:"tag_id"`
	Type       int    `form:"type"`
	Status     int    `form:"status"`
	IsDelete   *bool  `form:"is_delete"`
}

type UpdateArticleTopReq struct {
	ID    int  `json:"id"`
	IsTop bool `json:"is_top"`
}

type ArticleVO struct {
	model.Article

	LikeCount    int `json:"like_count" gorm:"-"`
	ViewCount    int `json:"view_count" gorm:"-"`
	CommentCount int `json:"comment_count" gorm:"-"`
}

// @Summary 新增或编辑文章
// @Description 新增或编辑文章, 分类和标签不存在时按名称自动创建
// @Tags Article
// @Accept json
// @Produce json
// @Param form body AddOrEditArticleReq true "新增或编辑文章"
// @Success 0 {object} Response[model.Article]
// @Security ApiKeyAuth
// @Router /article [post]
func (*Article) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	if req.Img == "" {
		req.Img = model.GetConfig(db, g.CONFIG_ARTICLE_COVER) // 默认图片
	}

	if req.Type == 0 {
		req.Type = 1 // 默认为原创
	}

	article := model.Article{
		Model:       model.Model{ID: req.ID},
		Title:       req.Title,
		Desc:        req.Desc,
		Content:     req.Content,
		Img:         req.Img,
		Type:        req.Type,
		Status:      req.Status,
		OriginalUrl: req.OriginalUrl,
		IsTop:       req.IsTop,
		UserId:      auth.UserInfoId,
	}

	err := model.SaveOrUpdateArticle(db, &article, req.CategoryName, req.TagNames)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

// @Summary 软删除文章（批量）
// @Description 将文章移入或移出回收站
// @Tags Article
// @Accept json
// @Produce json
// @Param form body SoftDeleteReq true "软删除文章"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /article/soft-delete [put]
func (*Article) UpdateSoftDelete(c *gin.Context) {
	var req SoftDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.UpdateArticleSoftDelete(GetDB(c), req.Ids, req.IsDelete)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// @Summary 物理删除文章（批量）
// @Description 根据 ID 数组物理删除文章, 同时删除文章标签关联与 Redis 中的浏览量/点赞记录
// @Tags Article
// @Accept json
// @Produce json
// @Param ids body []int true "文章 ID 数组"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /article [delete]
func (*Article) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, commentIds, err := model.DeleteArticle(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 数据库里的文章没了, Redis 里的计数也要跟着清,
	// 否则新文章复用到同一个 id 时会直接继承旧文章的浏览量和点赞数
	cleanArticleCounters(GetRDB(c), ids)
	cleanCommentCounters(GetRDB(c), commentIds)

	ReturnSuccess(c, rows)
}

// 清理文章相关的 Redis 计数: 浏览量、点赞数、用户点赞记录、浏览去重 key
func cleanArticleCounters(rdb *redis.Client, ids []int) {
	if len(ids) == 0 {
		return
	}

	members := make([]any, 0, len(ids))
	fields := make([]string, 0, len(ids))
	for _, id := range ids {
		members = append(members, strconv.Itoa(id))
		fields = append(fields, strconv.Itoa(id))
	}

	if err := rdb.ZRem(rctx, g.ARTICLE_VIEW_COUNT, members...).Err(); err != nil {
		slog.Error("清理文章浏览量失败", "ids", ids, "err", err)
	}
	if err := rdb.HDel(rctx, g.ARTICLE_LIKE_COUNT, fields...).Err(); err != nil {
		slog.Error("清理文章点赞数失败", "ids", ids, "err", err)
	}

	// 用户点赞记录是按用户存的 Set, 只能扫出来把文章 id 从每个集合里移除
	if err := scanKeys(rdb, g.ARTICLE_USER_LIKE_SET+"*", func(key string) {
		rdb.SRem(rctx, key, members...)
	}); err != nil {
		slog.Error("清理用户点赞记录失败", "ids", ids, "err", err)
	}

	// 浏览去重的 key 一并删掉, 不然新文章的首次访问会被旧记录挡住
	for _, id := range ids {
		if err := scanKeys(rdb, g.ARTICLE_VIEW_VISITOR+strconv.Itoa(id)+":*", func(key string) {
			rdb.Del(rctx, key)
		}); err != nil {
			slog.Error("清理文章浏览去重记录失败", "id", id, "err", err)
		}
	}
}

// 清理评论相关的 Redis 计数: 点赞数、用户点赞记录
func cleanCommentCounters(rdb *redis.Client, ids []int) {
	if len(ids) == 0 {
		return
	}

	members := make([]any, 0, len(ids))
	fields := make([]string, 0, len(ids))
	for _, id := range ids {
		members = append(members, strconv.Itoa(id))
		fields = append(fields, strconv.Itoa(id))
	}

	if err := rdb.HDel(rctx, g.COMMENT_LIKE_COUNT, fields...).Err(); err != nil {
		slog.Error("清理评论点赞数失败", "ids", ids, "err", err)
	}
	if err := scanKeys(rdb, g.COMMENT_USER_LIKE_SET+"*", func(key string) {
		rdb.SRem(rctx, key, members...)
	}); err != nil {
		slog.Error("清理用户评论点赞记录失败", "ids", ids, "err", err)
	}
}

// 按 pattern 扫描 key 并逐个处理, 用 SCAN 而不是 KEYS, 避免阻塞 Redis
func scanKeys(rdb *redis.Client, pattern string, fn func(key string)) error {
	iter := rdb.Scan(rctx, 0, pattern, 100).Iterator()
	for iter.Next(rctx) {
		fn(iter.Val())
	}
	return iter.Err()
}

// @Summary 条件查询文章列表
// @Description 后台文章列表, 支持标题/分类/标签/类型/状态/回收站过滤
// @Tags Article
// @Produce json
// @Param title query string false "标题关键字"
// @Param category_id query int false "分类 ID"
// @Param tag_id query int false "标签 ID"
// @Param type query int false "类型(1-原创 2-转载 3-翻译)"
// @Param status query int false "状态(1-公开 2-私密 3-草稿)"
// @Param is_delete query bool false "是否在回收站"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[ArticleVO]]
// @Security ApiKeyAuth
// @Router /article/list [get]
func (*Article) GetList(c *gin.Context) {
	var query ArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	list, total, err := model.GetArticleList(db, query.Page, query.Size, query.Title, query.IsDelete, query.Status, query.Type, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ids := make([]int, 0, len(list))
	for _, article := range list {
		ids = append(ids, article.ID)
	}
	likeCountMap := hashCounts(rdb, g.ARTICLE_LIKE_COUNT, ids)
	viewCountMap := zsetCounts(rdb, g.ARTICLE_VIEW_COUNT, ids)

	data := make([]ArticleVO, 0)
	for _, article := range list {
		data = append(data, ArticleVO{
			Article:   article,
			LikeCount: likeCountMap[article.ID],
			ViewCount: viewCountMap[article.ID],
		})
	}

	ReturnSuccess(c, PageResult[ArticleVO]{
		Size:  query.Size,
		Page:  query.Page,
		Total: total,
		List:  data,
	})

}

// @Summary 获取文章详情
// @Description 根据 ID 获取文章详情, 包含分类和标签
// @Tags Article
// @Produce json
// @Param id path int true "文章 ID"
// @Success 0 {object} Response[model.Article]
// @Security ApiKeyAuth
// @Router /article/{id} [get]
func (*Article) GetDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	article, err := model.GetArticle(GetDB(c), id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

// @Summary 修改文章置顶
// @Description 修改文章置顶状态
// @Tags Article
// @Accept json
// @Produce json
// @Param form body UpdateArticleTopReq true "修改置顶"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /article/top [put]
func (*Article) UpdateTop(c *gin.Context) {
	var req UpdateArticleTopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.UpdateArticleTop(GetDB(c), req.ID, req.IsTop)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// TODO: 目前是前端导出
// @Summary 导出文章
// @Description 目前由前端导出, 后端只返回空数据
// @Tags Article
// @Produce json
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /article/export [post]
func (*Article) Export(c *gin.Context) {
	ReturnSuccess(c, nil)
}

// @Summary 导入文章
// @Description 上传 Markdown 文件导入文章, 文件名作为标题, 导入后为草稿
// @Tags Article
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Markdown 文件"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /article/import [post]
func (*Article) Import(c *gin.Context) {
	db := GetDB(c)
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}

	fileName := fileHeader.Filename
	title := fileName[:len(fileName)-3]
	content, err := readFromFileHeader(fileHeader)
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}

	defaultImg := model.GetConfig(db, g.CONFIG_ARTICLE_COVER)
	err = model.ImportArticle(db, auth.ID, title, content, defaultImg, "学习", "Golang")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

func readFromFileHeader(file *multipart.FileHeader) (string, error) {
	open, err := file.Open()
	if err != nil {
		slog.Error("文件读取, 目标地址错误", "err", err)
		return "", err
	}
	defer open.Close()
	all, err := io.ReadAll(open)
	if err != nil {
		slog.Error("文件读取失败", "err", err)
		return "", err
	}
	return string(all), nil
}
