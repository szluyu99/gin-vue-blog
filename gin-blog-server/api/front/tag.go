package front

import (
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Tag struct{}

func (*Tag) GetFrontList(c *gin.Context) {
	r.SuccessData(c, tagService.GetFrontList())
}
