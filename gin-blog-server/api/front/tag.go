package front

import (
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Tag struct{}

func (*Tag) GetFrontList(c *gin.Context) {
	r.Success(c, tagService.GetFrontList())
}
