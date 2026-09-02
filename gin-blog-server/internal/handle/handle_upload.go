package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils/upload"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type Upload struct{}

var (
	errFileExt  = errors.New("不支持的文件后缀")
	errFileMime = errors.New("文件内容不是图片")
)

// 上传大小上限: 10MB
const maxUploadSize = 10 << 20

// 允许上传的图片后缀
// 不放 svg: svg 里可以写 <script>, 上传后从同源访问就是存储型 XSS
var allowedImageExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
}

// 允许上传的图片类型, 由文件内容嗅探得出, 不看客户端传的 Content-Type
var allowedImageMime = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	"image/bmp":  true,
}

// @Summary 上传文件
// @Description 上传图片到本地或七牛云, 返回访问地址。需要登录, 只允许图片, 单个文件不超过 10MB
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Success 0 {object} Response[string]
// @Security ApiKeyAuth
// @Router /upload [post]
// @Router /front/upload [post]
func (*Upload) UploadFile(c *gin.Context) {
	// 前台的 /front/upload 走 JWTAuth(false), 匿名请求也能到这里,
	// 不自己拦一次的话任何人都能往服务器写文件
	if _, ok := MustCurrentUserAuth(c); !ok {
		return
	}

	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}

	if fileHeader.Size > maxUploadSize {
		ReturnError(c, g.ErrFileSize, nil)
		return
	}

	if err := checkImage(fileHeader); err != nil {
		ReturnError(c, g.ErrFileType, err)
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

// 校验是不是图片: 后缀与文件内容都要对得上
func checkImage(fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(path.Ext(fileHeader.Filename))
	if !allowedImageExt[ext] {
		return errFileExt
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// http.DetectContentType 只看前 512 字节
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}

	mime := strings.Split(http.DetectContentType(buf[:n]), ";")[0]
	if !allowedImageMime[mime] {
		return errFileMime
	}
	return nil
}
