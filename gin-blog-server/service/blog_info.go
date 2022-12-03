package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"strings"

	"github.com/gin-gonic/gin"
)

type BlogInfo struct{}

// Cache Aside Pattern: 经典的缓存 + 数据库的读写模式

// 用户上报信息: 统计地域信息, 访问数量
func (*BlogInfo) Report(c *gin.Context) (code int) {
	ipAddress, browser, os := utils.IP.GetInfo(c)
	uuid := utils.Encryptor.MD5(ipAddress + browser + os)
	// 当前用户没有统计过访问人数 (不在 用户set 中)
	if !utils.Redis.SIsMember(KEY_UNIQUE_VISITOR_SET, uuid) {
		// 统计地域信息
		ipSource := utils.IP.GetIpSource(ipAddress)
		if ipSource != "未知" { // 获取到具体的位置, 提取出其中的 省份
			address := strings.Split(ipSource, "|")
			ipSource = strings.ReplaceAll(address[2], "省", "")
		}
		utils.Redis.HIncrBy(KEY_VISITOR_AREA, ipSource, 1)
		// 访问数量 + 1
		utils.Redis.Incr(KEY_VIEW_COUNT)
		// 将当前用户记录到 用户set
		utils.Redis.SAdd(KEY_UNIQUE_VISITOR_SET, uuid)
	}

	return r.OK
}

// 获取博客首页信息
func (b *BlogInfo) GetHomeInfo() resp.BlogHomeVO {
	articleCount := dao.Count(model.Article{}, "status = ? AND is_delete = ?", 1, 0)
	userCount := dao.Count(model.UserInfo{}, "")
	messageCount := dao.Count(model.Message{}, "")
	viewCount := utils.Redis.GetInt(KEY_VIEW_COUNT)
	// categoryCount := dao.Count(model.Category{}, "")
	// tagCount := dao.Count(model.Tag{}, "")
	// blogConfigVO := b.GetBlogConfig()
	return resp.BlogHomeVO{
		ArticleCount: articleCount,
		UserCount:    userCount,
		MessageCount: messageCount,
		ViewCount:    viewCount,
		// CategoryCount: categoryCount,
		// TagCount:      tagCount,
		// BlogConfig:    blogConfigVO,
	}
}

// 获取博客设置
func (*BlogInfo) GetBlogConfig() (respVO model.BlogConfigDetail) {
	// 尝试从 Redis 中取值
	blogConfig := utils.Redis.GetVal(KEY_BLOG_CONFIG)
	if blogConfig == "" {
		blogConfig = dao.GetOne(model.BlogConfig{}, "id", 1).Config
		utils.Redis.Set(KEY_BLOG_CONFIG, blogConfig, 0)
	}
	// 存储到 Redis
	utils.Json.Unmarshal(blogConfig, &respVO)
	return respVO
}

// 更新博客设置
func (*BlogInfo) UpdateBlogConfig(reqVO model.BlogConfigDetail) (code int) {
	blogConfig := model.BlogConfig{
		Universal: model.Universal{ID: 1},
		Config:    utils.Json.Marshal(reqVO),
	}
	dao.Update(&blogConfig, "config")
	utils.Redis.Del(KEY_BLOG_CONFIG) // 从 Redis 中清除旧值
	return r.OK
}

func (*BlogInfo) GetAbout() string {
	return utils.Redis.GetVal(KEY_ABOUT)
}

func (*BlogInfo) UpdateAbout(data model.About) (code int) {
	utils.Redis.Set(KEY_ABOUT, data.Content, 0)
	return r.OK
}

/* 前台借口 */

// 获取前台首页信息
func (b *BlogInfo) GetFrontHomeInfo() resp.FrontBlogHomeVO {
	articleCount := dao.Count(model.Article{}, "status = ? AND is_delete = ?", 1, 0)
	userCount := dao.Count(model.UserInfo{}, "")
	messageCount := dao.Count(model.Message{}, "")
	categoryCount := dao.Count(model.Category{}, "")
	tagCount := dao.Count(model.Tag{}, "")
	blogConfigVO := b.GetBlogConfig()
	viewCount := utils.Redis.GetInt(KEY_VIEW_COUNT)
	return resp.FrontBlogHomeVO{
		ArticleCount:  articleCount,
		UserCount:     userCount,
		MessageCount:  messageCount,
		CategoryCount: categoryCount,
		ViewCount:     viewCount,
		TagCount:      tagCount,
		BlogConfig:    blogConfigVO,
	}
}
