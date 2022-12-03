package front

import (
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type BlogInfo struct{}

// 获取前台首页信息
func (*BlogInfo) GetFrontHomeInfo(c *gin.Context) {
	r.SuccessData(c, blogInfoService.GetFrontHomeInfo())
}
