package v1

import (
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type OperationLog struct{}

func (*OperationLog) GetList(c *gin.Context) {
	r.Success(c, operationLogService.GetList(utils.BindPageQuery(c)))
}

func (*OperationLog) Delete(c *gin.Context) {
	r.Code(c, operationLogService.Delete(utils.BindJson[[]int](c)))
}
