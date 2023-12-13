package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Message struct{}

func (*Message) Delete(c *gin.Context) {
	r.Code(c, messageService.Delete(utils.BindJson[[]int](c)))
}

// 修改审核
func (*Message) UpdateReview(c *gin.Context) {
	r.Code(c, messageService.UpdateReview(
		utils.BindValidJson[req.UpdateReview](c)))
}

// 条件查询列表(后台)
func (*Message) GetList(c *gin.Context) {
	var req = utils.BindValidQuery[req.GetMessages](c)
	utils.CheckQueryPage(&req.PageSize, &req.PageNum)
	r.Success(c, messageService.GetList(req))
}
