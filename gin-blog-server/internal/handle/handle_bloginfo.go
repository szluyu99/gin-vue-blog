package handle

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type BlogInfo struct{}

type BlogHomeVO struct {
	ArticleCount int `json:"article_count"` // 文章数量
	UserCount    int `json:"user_count"`    // 用户数量
	MessageCount int `json:"message_count"` // 留言数量
	ViewCount    int `json:"view_count"`    // 访问量
	// CategoryCount int64 `json:"category_count"` // 分类数量
	// TagCount      int64 `json:"tag_count"`      // 标签数量
	// BlogConfig    model.BlogConfigDetail `json:"blog_config"`    // 博客信息
	// PageList      []Page                 `json:"pageList"`
}

type AboutReq struct {
	Content string `json:"content"`
}

func (*BlogInfo) GetConfigMap(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	// get from redis cache
	cache, err := getConfigCache(rdb)
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	if len(cache) > 0 {
		slog.Debug("get config from redis cache")
		ReturnSuccess(c, cache)
		return
	}

	// get from db
	data, err := model.GetConfigMap(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// add to redis cache
	if err := addConfigCache(rdb, data); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, data)
}

func (*BlogInfo) UpdateConfig(c *gin.Context) {
	var m map[string]string
	if err := c.ShouldBindJSON(&m); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if err := model.CheckConfigMap(GetDB(c), m); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// delete cache
	if err := removeConfigCache(GetRDB(c)); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// @Summary 获取博客首页信息
// @Description 获取博客首页信息
// @Tags blog_info
// @Produce json
// @Success 0 {object} Response[model.BlogHomeVO]
// @Router /home [get]
func (*BlogInfo) GetHomeInfo(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	articleCount, err := model.Count(db, &model.Article{}, "status = ? AND is_delete = ?", 1, 0)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	userCount, err := model.Count(db, &model.UserInfo{})
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	messageCount, err := model.Count(db, &model.Message{})
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	viewCount, err := rdb.Get(rctx, g.VIEW_COUNT).Int()
	if err != nil && err != redis.Nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, BlogHomeVO{
		ArticleCount: articleCount,
		UserCount:    userCount,
		MessageCount: messageCount,
		ViewCount:    viewCount,
	})
}

// @Summary 获取关于
// @Description 获取关于
// @Tags blog_info
// @Produce json
// @Success 0 {object} Response[string]
// @Router /about [get]
func (*BlogInfo) GetAbout(c *gin.Context) {
	ReturnSuccess(c, model.GetConfig(GetDB(c), g.CONFIG_ABOUT))
}

// @Summary 更新关于
// @Description 更新关于
// @Tags blog_info
// @Accept json
// @Produce json
// @Param data body object true "关于"
// @Success 0 {object} Response[string]
// @Router /about [put]
func (*BlogInfo) UpdateAbout(c *gin.Context) {
	var req AboutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.CheckConfig(GetDB(c), g.CONFIG_ABOUT, req.Content)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, req.Content)
}

// @Summary 上报用户信息
// @Description 用户登进后台时上报信息
// @Tags blog_info
// @Accept json
// @Produce json
// @Param data body object true "用户信息"
// @Success 0 {object} Response[any]
// @Router /report [post]
func (*BlogInfo) Report(c *gin.Context) {
	rdb := GetRDB(c)

	ipAddress := utils.IP.GetIpAddress(c)
	userAgent := utils.IP.GetUserAgent(c)
	browser := userAgent.Name + " " + userAgent.Version.String()
	os := userAgent.OS + " " + userAgent.OSVersion.String()
	uuid := utils.MD5(ipAddress + browser + os)

	ctx := context.Background()

	// 当前用户没有统计过访问人数 (不在 用户set 中)
	if !rdb.SIsMember(ctx, g.KEY_UNIQUE_VISITOR_SET, uuid).Val() {
		// 统计地域信息
		ipSource := utils.IP.GetIpSource(ipAddress)
		if ipSource != "" { // 获取到具体的位置, 提取出其中的 省份
			address := strings.Split(ipSource, "|")
			province := strings.ReplaceAll(address[2], "省", "")
			rdb.HIncrBy(ctx, g.VISITOR_AREA, province, 1)
		} else {
			rdb.HIncrBy(ctx, g.VISITOR_AREA, "未知", 1)
		}
		// 访问数量 + 1
		rdb.Incr(ctx, g.VIEW_COUNT)
		// 将当前用户记录到 用户set
		rdb.SAdd(ctx, g.KEY_UNIQUE_VISITOR_SET, uuid)
	}

	ReturnSuccess(c, nil)
}

// 获取博客设置
// func GetBlogConfig() model.BlogConfigDetail {
// 	// 尝试从 Redis 中取值
// 	blogConfig := utils.Redis.GetVal(KEY_BLOG_CONFIG)
// 	// Redis 中没有值, 再查数据库, 查到后设置到 Redis 中
// 	if blogConfig == "" {
// 		blogConfig = dao.GetOne(model.BlogConfig{}, "id", 1).Config
// 		utils.Redis.Set(KEY_BLOG_CONFIG, blogConfig, 0)
// 	}
// 	// 反序列化字符串为 golang 对象
// 	var result model.BlogConfigDetail
// 	utils.Json.Unmarshal(blogConfig, &result)
// 	return result
// }
