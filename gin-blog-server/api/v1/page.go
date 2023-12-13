package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Page struct{}

func (*Page) GetList(c *gin.Context) {
	r.Success(c, pageService.GetList())
}

func (*Page) SaveOrUpdate(c *gin.Context) {
	r.Code(c, pageService.SaveOrUpdate(utils.BindJson[req.AddOrEditPage](c)))
}

func (*Page) Delete(c *gin.Context) {
	r.Code(c, pageService.Delete(utils.BindJson[[]int](c)))
}
