package service

import (
	"gin-blog/utils/r"
	"gin-blog/utils/upload"
	"mime/multipart"
)

type Upload struct{}

// 上传本地文件
func (s *Upload) UploadFile(header *multipart.FileHeader) (url string, code int) {
	oss := upload.NewOSS() // 接口
	filePath, _, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		return "", r.ERROR_FILE_UPLOAD
	}
	return filePath, r.OK
	// 利用 gin.Context 中的方法
	// 文件上传最终路径
	//timeStr := time.Now().Format("2006_01_02_15_04_05")                             // 当前时间
	//returnPath := global.G_CONFIG.Upload.Path + timeStr + path.Ext(header.Filename) // 文件上传后的返回路径
	//savePath := global.BasePath + "/" + returnPath                                  // 文件上传到的路径
	//// 进行文件上传
	//err := c.SaveUploadedFile(header, savePath)
	//if err != nil {
	//	return "", r.ERROR_FILE_UPLOAD
	//}
	//return returnPath, r.SUCCESS
}
