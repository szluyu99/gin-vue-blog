package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Article struct{}

/* 后台接口 */

// 软删除(回收站)
func (*Article) SoftDelete(ids []int, isDelete *int8) (code int) {
	var isTop int8 = 0
	dao.Updates(&model.Article{IsTop: &isTop, IsDelete: isDelete}, "id IN ?", ids)
	return r.OK
}

func (*Article) Delete(ids []int) (code int) {
	// 删除 [文章-标签] 关联
	dao.Delete(model.ArticleTag{}, "article_id IN ?", ids)
	// 删除 [文章]
	dao.Delete(model.Article{}, "id IN ?", ids)
	return r.OK
}

func (*Article) UpdateTop(req req.UpdateArticleTop) (code int) {
	article := model.Article{
		Universal: model.Universal{ID: req.ID},
		IsTop:     req.IsTop,
	}
	dao.Update(&article, "is_top")
	return r.OK
}

func (*Article) GetList(req req.GetArts) resp.PageResult[[]resp.ArticleVO] {
	data, total := articleDao.GetList(req)
	// TODO: 浏览量, 点赞量
	return resp.PageResult[[]resp.ArticleVO]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    total,
		List:     data,
	}
}

// 根据 [文章id] 获取 [文章详情]
func (*Article) GetInfo(id int) resp.ArticleDetailVO {
	article := dao.GetOne(model.Article{}, "id", id)
	category := dao.GetOne(model.Category{}, "id", article.CategoryId)
	tagNames := tagDao.GetTagNamesByArtId(id)

	articleVo := utils.CopyProperties[resp.ArticleDetailVO](article)
	articleVo.CategoryName = category.Name
	articleVo.TagNames = tagNames

	return articleVo
}

// TODO: 添加事务
func (*Article) SaveOrUpdate(req req.SaveOrUpdateArt, userId int) (code int) {
	article := utils.CopyProperties[model.Article](req) // po -> vo
	article.UserId = userId

	// 维护 [文章-分类] 关联
	category := saveArticleCategory(req)
	if !category.IsEmpty() {
		article.CategoryId = category.ID
	}

	// 设置默认图片 (blogConfig 中配置)
	if req.Img == "" {
		req.Img = blogInfoService.GetBlogConfig().ArticleCover // 设置默认图片
	}

	// 先 添加/更新 文章, 可以获取到 ID
	if article.ID == 0 {
		dao.Create(&article)
	} else {
		dao.Update(&article)
	}

	// 维护 [文章-标签] 关联
	saveArticleTag(req, article.ID)
	return r.OK
}

// 维护 [文章-分类] 关联
func saveArticleCategory(req req.SaveOrUpdateArt) model.Category {
	category := dao.GetOne(model.Category{}, "name = ?", req.CategoryName)
	if category.IsEmpty() && req.Status != model.DRAFT {
		category.Name = req.CategoryName
		dao.Create(&category)
	}
	return category
}

// 维护 [文章-标签] 关联
func saveArticleTag(req req.SaveOrUpdateArt, articleId int) {
	// 清除文章对应的标签关联
	if req.ID != 0 {
		dao.Delete(model.ArticleTag{}, "article_id = ?", req.ID)
	}
	// 遍历 req.TagNams 中传来的标签, 不存在则新建
	var articleTags []model.ArticleTag // 文章-标签 关系
	for _, tagName := range req.TagNames {
		tag := dao.GetOne(model.Tag{}, "name = ?", tagName)
		// 标签不存在则创建
		if tag.IsEmpty() {
			tag.Name = tagName
			dao.Create(&tag)
		}
		articleTags = append(articleTags, model.ArticleTag{
			ArticleId: articleId,
			TagId:     tag.ID,
		})
	}
	dao.Create(&articleTags)
}

/* 前台接口 */
// 获取前台文章列表
func (*Article) GetFrontList(req req.GetFrontArts) []resp.FrontArticleVO {
	list, _ := articleDao.GetFrontList(req)
	return list
}

// 获取前台文章详情
func (*Article) GetFrontInfo(c *gin.Context, id int) resp.FrontArticleDetailVO {
	// 查询具体文章
	article := articleDao.GetInfoById(id)
	// 查询推荐文章 (6篇)
	article.RecommendArticles = articleDao.GetRecommendList(id, 6)
	// 查询最新文章 (5篇)
	article.NewestArticles = articleDao.GetNewestList(5)
	// 更新文章浏览量 TODO: 删除文章时删除其浏览量
	// updateArticleViewCount(c, id)
	// * 请求一次就会增加访问量, 即刷新可以刷访问量
	utils.Redis.ZincrBy(KEY_ARTICLE_VIEW_COUNT, strconv.Itoa(id), 1)
	// 获取上一篇文章, 下一篇文章
	article.LastArticle = articleDao.GetLast(id)
	article.NextArticle = articleDao.GetNext(id)
	// 点赞量, 浏览量
	article.ViewCount = utils.Redis.ZScore(KEY_ARTICLE_VIEW_COUNT, strconv.Itoa(id))
	article.LikeCount = utils.Redis.HGet(KEY_ARTICLE_LIKE_COUNT, strconv.Itoa(id))
	return article
}

// 获取前台文章归档
func (*Article) GetArchiveList(req req.GetFrontArts) resp.PageResult[[]resp.ArchiveVO] {
	articles, total := articleDao.GetFrontList(req)
	archives := make([]resp.ArchiveVO, 0)
	for _, article := range articles {
		archives = append(archives, resp.ArchiveVO{
			ID:         article.ID,
			Title:      article.Title,
			Created_at: article.CreatedAt,
		})
	}
	return resp.PageResult[[]resp.ArchiveVO]{
		Total:    total,
		List:     archives,
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
}

// 前台文章点赞
func (*Article) SaveLike(uid, articleId int) (code int) {
	// 记录某个用户已经对某个文章点过赞
	articleLikeUserKey := KEY_ARTICLE_USER_LIKE_SET + strconv.Itoa(uid)
	// 该文章已经被记录过, 再点赞就是取消点赞
	if utils.Redis.SIsMember(articleLikeUserKey, articleId) {
		utils.Redis.SRem(articleLikeUserKey, articleId)
		utils.Redis.HIncrBy(KEY_ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), -1)
	} else { // 未被记录过, 则是增加点赞
		utils.Redis.SAdd(articleLikeUserKey, articleId)
		utils.Redis.HIncrBy(KEY_ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), 1)
	}
	return r.OK

}

// ? 更新文章访问数量?? 防止刷新刷访问量?
// func updateArticleViewCount(c *gin.Context, articleId int) {
// 	session := sessions.Default(c)
// 	if session.Get(SESSION_ARTICLE_SET) == nil {
// 		var articleSet set.Set
// 		articleSet.Init()
// 		articleSet.Add(articleId)
// 		session.Set(SESSION_ARTICLE_SET, articleSet)
// 		utils.Redis.ZincrBy(KEY_ARTICLE_VIEW_COUNT, strconv.Itoa(articleId), 1)
// 	} else {
// 		articleSet := session.Get(SESSION_ARTICLE_SET).(set.Set)
// 		if !articleSet.Exist(articleId) {
// 			articleSet.Add(articleSet)
// 			session.Set(SESSION_ARTICLE_SET, articleSet)
// 			utils.Redis.ZincrBy(KEY_ARTICLE_VIEW_COUNT, strconv.Itoa(articleId), 1)
// 		}
// 	}
// 	session.Save()
// }
