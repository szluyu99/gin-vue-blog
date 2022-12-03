package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Comment struct{}

func (*Comment) Delete(c *gin.Context) {
	r.SendCode(c, commentService.Delete(utils.BindJson[[]int](c)))
}

// 修改审核
func (*Comment) UpdateReview(c *gin.Context) {
	r.SendCode(c, commentService.UpdateReview(utils.BindValidJson[req.UpdateReview](c)))
}

// TODO: 完善
func (*Comment) GetList(c *gin.Context) {
	var req = utils.BindValidQuery[req.GetComments](c)
	utils.CheckQueryPage(&req.PageSize, &req.PageNum)
	r.SuccessData(c, commentService.GetList(req))
}
