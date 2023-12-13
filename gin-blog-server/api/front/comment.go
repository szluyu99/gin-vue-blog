package front

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Comment struct{}

// 目前只能新增评论不能修改
func (*Comment) Save(c *gin.Context) {
	r.Code(c, commentService.Save(
		utils.GetFromContext[int](c, "user_info_id"),
		utils.BindValidJson[req.SaveComment](c)))
}

// 获取前台评论列表
func (*Comment) GetFrontList(c *gin.Context) {
	r.Success(c, commentService.GetFrontList(
		utils.BindQuery[req.GetFrontComments](c)))
}

// 点赞评论
func (*Comment) SaveLike(c *gin.Context) {
	uid := utils.GetFromContext[int](c, "user_info_id")
	commentId := utils.GetIntParam(c, "comment_id")
	r.Code(c, commentService.SaveLike(uid, commentId))
}

// 根据 [评论id] 获取 [回复列表]
func (*Comment) GetReplyListByCommentId(c *gin.Context) {
	commentId := utils.GetIntParam(c, "comment_id")
	r.Success(c, commentService.GetReplyListByCommentId(commentId, utils.BindPageQuery(c)))
}
