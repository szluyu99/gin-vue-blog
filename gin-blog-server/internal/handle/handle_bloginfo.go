package handle

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type BlogInfo struct{}

// 趋势图上的一天
type ViewTrendVO struct {
	Date  string `json:"date"`  // 2006-01-02
	Count int    `json:"count"` // 当天访问量
}

// 访客地域分布的一项
type VisitorAreaVO struct {
	Area  string `json:"area"`  // 省份, 取不到时是 "未知"
	Count int    `json:"count"` // 独立访客数
}

type BlogHomeVO struct {
	ArticleCount int `json:"article_count"` // 文章数量
	UserCount    int `json:"user_count"`    // 用户数量
	MessageCount int `json:"message_count"` // 留言数量
	ViewCount    int `json:"view_count"`    // 访问量
	// 最近 viewTrendDays 天的访问量, 按日期升序, 没有数据的那天补 0
	ViewTrend []ViewTrendVO `json:"view_trend"`
	// 访客地域分布, 按人数倒序, 最多 visitorAreaTop 项
	VisitorArea []VisitorAreaVO `json:"visitor_area"`
	// CategoryCount int64 `json:"category_count"` // 分类数量
	// TagCount      int64 `json:"tag_count"`      // 标签数量
	// BlogConfig    model.BlogConfigDetail `json:"blog_config"`    // 博客信息
	// PageList      []Page                 `json:"pageList"`
}

const (
	viewTrendDays  = 14 // 趋势图展示的天数
	visitorAreaTop = 8  // 地域分布最多展示几项, 再多在小卡片里挤不下
	viewDayTTL     = 30 * 24 * 60 * 60 * time.Second
	viewDayLayout  = "2006-01-02"
)

type AboutReq struct {
	Content string `json:"content"`
}

// @Summary 获取博客配置
// @Description 获取博客配置, 优先读 Redis 缓存
// @Tags BlogInfo
// @Produce json
// @Success 0 {object} Response[map[string]string]
// @Router /config [get]
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

// @Summary 更新博客配置
// @Description 更新博客配置, 同时清除 Redis 缓存
// @Tags BlogInfo
// @Accept json
// @Produce json
// @Param form body map[string]string true "博客配置"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /config [patch]
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

// @Summary 获取后台首页信息
// @Description 文章数, 用户数, 留言数, 访问量
// @Tags BlogInfo
// @Produce json
// @Success 0 {object} Response[BlogHomeVO]
// @Security ApiKeyAuth
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

	trend, err := getViewTrend(rdb, time.Now())
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	area, err := getVisitorArea(rdb)
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, BlogHomeVO{
		ArticleCount: articleCount,
		UserCount:    userCount,
		MessageCount: messageCount,
		ViewCount:    viewCount,
		ViewTrend:    trend,
		VisitorArea:  area,
	})
}

/*
访客地域分布

数据是 Report 里按省份累加的 Hash, 之前没有任何接口读它 —— 前台一直没调上报,
这个键一直是空的, 所以后台也就没做展示。上报接通之后它才有意义。
Hash 的遍历顺序不保证, 这里按人数倒序并截断到 visitorAreaTop, 不然刷新一次顺序就变一次。
*/
func getVisitorArea(rdb *redis.Client) ([]VisitorAreaVO, error) {
	all, err := rdb.HGetAll(rctx, g.VISITOR_AREA).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	list := make([]VisitorAreaVO, 0, len(all))
	for area, count := range all {
		n, convErr := strconv.Atoi(count)
		if convErr != nil {
			// 单个值坏掉不该让整个仪表盘失败
			slog.Warn("访客地域计数不是数字", "area", area, "value", count)
			continue
		}
		list = append(list, VisitorAreaVO{Area: area, Count: n})
	}

	// 人数相同时按名字排, 保证同样的数据每次返回同样的顺序
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Area < list[j].Area
	})

	if len(list) > visitorAreaTop {
		list = list[:visitorAreaTop]
	}
	return list, nil
}

/*
最近 viewTrendDays 天的访问量

一次 MGet 取回所有天, 不逐天 Get: 前者是一个来回, 后者 14 个。
没有数据的那天必须补 0 而不是跳过, 否则前端画出来横轴会缺格、趋势看着是错的。
*/
func getViewTrend(rdb *redis.Client, now time.Time) ([]ViewTrendVO, error) {
	dates := make([]string, 0, viewTrendDays)
	keys := make([]string, 0, viewTrendDays)
	for i := viewTrendDays - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format(viewDayLayout)
		dates = append(dates, date)
		keys = append(keys, g.VIEW_COUNT_DAY+date)
	}

	values, err := rdb.MGet(rctx, keys...).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	trend := make([]ViewTrendVO, 0, viewTrendDays)
	for i, date := range dates {
		count := 0
		if i < len(values) {
			// MGet 的元素是 string 或 nil(该 key 不存在)
			if s, ok := values[i].(string); ok {
				count, _ = strconv.Atoi(s)
			}
		}
		trend = append(trend, ViewTrendVO{Date: date, Count: count})
	}
	return trend, nil
}

/*
从 IP 解析出用于地域统计的省份

ip2region 的返回是「国家|区域|省份|城市|ISP」, 定位不到的字段是字面量 "0",
内网地址整条是 "0|0|0|内网IP|内网IP"。原来直接取 address[2], 内网访问就在
仪表盘上留下一个名字叫「0」的条目。取不到省份时统一归到「未知」。
*/
func visitorProvince(ipAddress string) string {
	region := utils.IP.GetIpSource(ipAddress)
	if region == "" {
		return "未知"
	}

	fields := strings.Split(region, "|")
	if len(fields) < 3 {
		return "未知"
	}

	province := strings.ReplaceAll(fields[2], "省", "")
	if province == "" || province == "0" {
		// 内网地址在城市位上写的是「内网IP」, 拿它比「未知」更有信息量
		if len(fields) > 3 && fields[3] != "" && fields[3] != "0" {
			return fields[3]
		}
		return "未知"
	}
	return province
}

// @Summary 获取关于我
// @Description 获取关于我的内容
// @Tags BlogInfo
// @Produce json
// @Success 0 {object} Response[string]
// @Router /setting/about [get]
// @Router /front/about [get]
func (*BlogInfo) GetAbout(c *gin.Context) {
	ReturnSuccess(c, model.GetConfig(GetDB(c), g.CONFIG_ABOUT))
}

// @Summary 更新关于我
// @Description 更新关于我的内容
// @Tags BlogInfo
// @Accept json
// @Produce json
// @Param form body AboutReq true "关于我"
// @Success 0 {object} Response[string]
// @Security ApiKeyAuth
// @Router /setting/about [put]
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

	// about 也是 config 表里的一行, 会被 GetConfigMap 一起缓存起来,
	// 这里不清缓存的话 /config 会一直返回改之前的「关于我」
	if err := removeConfigCache(GetRDB(c)); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, req.Content)
}

// @Summary 上报访客信息
// @Description 统计访问量, 独立访客数与访客地域分布
// @Tags BlogInfo
// @Produce json
// @Success 0 {object} Response[any]
// @Router /report [post]
func (*BlogInfo) Report(c *gin.Context) {
	rdb := GetRDB(c)

	ipAddress := utils.IP.GetIpAddress(c)
	uuid := visitorFingerprint(c)

	ctx := context.Background()

	// 当前用户没有统计过访问人数 (不在 用户set 中)
	if !rdb.SIsMember(ctx, g.KEY_UNIQUE_VISITOR_SET, uuid).Val() {
		// 统计地域信息
		rdb.HIncrBy(ctx, g.VISITOR_AREA, visitorProvince(ipAddress), 1)
		// 访问数量 + 1
		rdb.Incr(ctx, g.VIEW_COUNT)
		// 将当前用户记录到 用户set
		rdb.SAdd(ctx, g.KEY_UNIQUE_VISITOR_SET, uuid)
	}

	// 按天计数不做去重: 累计值统计的是"来过多少人", 趋势要看的是"每天来了多少次",
	// 放在上面的 if 里就只会记新访客, 老访客再来趋势图上看不到
	if err := incrViewCountOfDay(ctx, rdb, time.Now()); err != nil {
		// 趋势数据丢一天不影响访问本身, 只告警
		slog.Warn("按天访问量记录失败", "err", err)
	}

	ReturnSuccess(c, nil)
}

/*
当天访问量 +1, 并保证这个 key 带上 TTL

Incr 建出来的 key 是没有过期时间的, 必须单独 Expire。两条命令放进 TxPipelined:
中间失败会留下一个永不过期的 key, 30 天后就成了永久垃圾(和 config 缓存踩过的坑一样)。
每次都 Expire 而不是只在首次: 少一次 Exists 往返, 代价只是把过期时间往后顶,
反正同一天内顶到的还是同一个时间点。
*/
func incrViewCountOfDay(ctx context.Context, rdb *redis.Client, now time.Time) error {
	key := g.VIEW_COUNT_DAY + now.Format(viewDayLayout)
	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, viewDayTTL)
		return nil
	})
	return err
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
