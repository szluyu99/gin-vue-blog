package v1

import (
	"gin-blog/model"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type BlogInfo struct{}

// 获取博客首页信息
func (*BlogInfo) GetHomeInfo(c *gin.Context) {
	r.SuccessData(c, blogInfoService.GetHomeInfo())
}

// 获取博客配置信息
func (*BlogInfo) GetBlogConfig(c *gin.Context) {
	r.SuccessData(c, blogInfoService.GetBlogConfig())
}

// 更新博客配置信息
func (*BlogInfo) UpdateBlogConfig(c *gin.Context) {
	r.SendCode(c, blogInfoService.UpdateBlogConfig(utils.BindJson[model.BlogConfigDetail](c)))
}

// 获取关于
func (*BlogInfo) GetAbout(c *gin.Context) {
	r.SuccessData(c, blogInfoService.GetAbout())
}

// 更新关于
func (*BlogInfo) UpdateAbout(c *gin.Context) {
	r.SendCode(c, blogInfoService.UpdateAbout(utils.BindJson[model.About](c)))
}

// 用户登进后台时上报信息: 如 ip 等信息
func (*BlogInfo) Report(c *gin.Context) {
	r.SendCode(c, blogInfoService.Report(c))
}
