package front

import (
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Link struct{}

func (*Link) GetFrontList(c *gin.Context) {
	r.SuccessData(c, linkService.GetFrontList())
}
