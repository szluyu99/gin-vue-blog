package v1

import (
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Article struct{}

func (*Article) SaveOrUpdate(c *gin.Context) {
	// 通过解析 token 获取当前文章的作者
	r.SendCode(c, articleService.SaveOrUpdate(
		utils.BindValidJson[req.SaveOrUpdateArt](c),
		utils.GetFromContext[int](c, "user_info_id")))
}

// 软删除或恢复
func (*Article) SoftDelete(c *gin.Context) {
	req := utils.BindValidJson[req.SoftDelete](c)
	r.SendCode(c, articleService.SoftDelete(req.Ids, req.IsDelete))
}

// 物理删除
func (*Article) Delete(c *gin.Context) {
	r.SendCode(c, articleService.Delete(utils.BindJson[[]int](c)))
}

// 获取文章列表
func (*Article) GetList(c *gin.Context) {
	r.SuccessData(c, articleService.GetList(utils.BindQuery[req.GetArts](c)))
}

// 获取文章详细信息
func (*Article) GetInfo(c *gin.Context) {
	r.SuccessData(c, articleService.GetInfo(utils.GetIntParam(c, "id")))
}

// 修改置顶信息
func (*Article) UpdateTop(c *gin.Context) {
	r.SendCode(c, articleService.UpdateTop(utils.BindValidJson[req.UpdateArtTop](c)))
}

// 导出文章: 获取导出后的资源链接列表
func (*Article) Export(c *gin.Context) {
	r.SuccessData(c, articleService.Export(utils.BindJson[[]int](c)))
}

// 导入文章: 题目 + 内容
func (*Article) Import(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		r.SendCode(c, r.EEROR_FILE_RECEIVE)
		return
	}
	fileName := fileHeader.Filename
	articleService.Import(
		fileName[:len(fileName)-3],
		utils.File.ReadFromFileHeader(fileHeader),
		utils.GetFromContext[int](c, "user_info_id"),
	)
	r.Success(c)
}
