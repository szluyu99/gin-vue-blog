package v1

import (
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Upload struct{}

func (*Upload) UploadFile(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	// 处理文件接收错误
	if err != nil {
		code := r.EEROR_FILE_RECEIVE
		utils.Logger.Error(r.GetMsg(code), zap.Error(err))
		r.Code(c, code)
		return
	}
	// 上传文件, 获取文件路径
	url, code := fileService.UploadFile(fileHeader)
	r.Data(c, code, url)
}
