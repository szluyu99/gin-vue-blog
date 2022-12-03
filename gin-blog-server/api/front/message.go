package front

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Message struct{}

// 查询列表(前台)
func (*Message) GetFrontList(c *gin.Context) {
	r.SuccessData(c, messageService.GetFrontList())
}

// 留言不能编辑, 只能新增
func (*Message) Save(c *gin.Context) {
	messageService.Save(c, utils.BindValidJson[req.AddMessage](c))
	r.Success(c)
}
