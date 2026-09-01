package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/utils/upload"

	"github.com/gin-gonic/gin"
)

type Upload struct{}

// @Summary 上传文件
// @Description 上传文件到本地或七牛云, 返回访问地址
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 0 {object} Response[string]
// @Security ApiKeyAuth
// @Router /upload [post]
// @Router /front/upload [post]
func (*Upload) UploadFile(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}

	oss := upload.NewOSS()
	filePath, _, err := oss.UploadFile(fileHeader)
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}

	ReturnSuccess(c, filePath)
}
