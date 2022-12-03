package front

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Article struct{}

// TODO: 前台相关: 文章归档, 点赞, 文章列表, 根据 id 获取文章具体内容...
// 获取前台文章列表
func (*Article) GetFrontList(c *gin.Context) {
	r.SuccessData(c, articleService.GetFrontList(utils.BindQuery[req.GetFrontArts](c)))
}

// 根据 id 获取前台文章详情
func (*Article) GetFrontInfo(c *gin.Context) {
	r.SuccessData(c, articleService.GetFrontInfo(c, utils.GetIntParam(c, "id")))
}

// 获取文章归档
func (*Article) GetArchiveList(c *gin.Context) {
	r.SuccessData(c, articleService.GetArchiveList(utils.BindQuery[req.GetFrontArts](c)))
}

// 点赞文章
func (*Article) SaveLike(c *gin.Context) {
	uid := utils.GetFromContext[int](c, "user_info_id")
	articleId := utils.GetIntParam(c, "article_id")
	r.SendCode(c, articleService.SaveLike(uid, articleId))
}
