package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Link struct{}

func (*Link) GetList(c *gin.Context) {
	r.Success(c, linkService.GetList(utils.BindPageQuery(c)))
}

func (*Link) SaveOrUpdate(c *gin.Context) {
	r.Code(c, linkService.SaveOrUpdate(utils.BindValidJson[req.SaveOrUpdateLink](c)))
}

func (*Link) Delete(c *gin.Context) {
	r.Success(c, linkService.Delete(utils.BindJson[[]int](c)))
}
