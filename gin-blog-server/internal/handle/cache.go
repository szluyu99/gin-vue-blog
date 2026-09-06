package handle

import (
	"context"
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis context
var rctx = context.Background()

/*
读穿缓存的过期时间

只给「数据库是真源、Redis 只是副本」的键用 —— 目前是 page 和 config。
写操作会主动删缓存, TTL 是给绕过接口的改动兜底: 直接 UPDATE 数据库、跑
generate-data 灌种子、或者换一个后端进程写库, 这些路径都不会触发主动失效,
以前只能手动 redis-cli del 再重启。

注意不要给点赞数/浏览数/访客地域那几个键加 TTL: 它们只在 Redis 里累加,
数据库没有对应字段, Redis 就是唯一数据源, 过期等于丢数据。
*/
const cacheTTL = 10 * time.Minute

/*
按 id 批量取计数

列表接口以前用 HGetAll / ZRange(0, -1) 把整个计数集合拉回来, 只为给当页十几条
数据填数, 开销随文章/评论总量线性增长。这里只取当页需要的字段。
计数属于展示信息, 取不到就当 0, 不让 Redis 抖动影响列表本身。
*/
func hashCounts(rdb *redis.Client, key string, ids []int) map[int]int {
	counts := make(map[int]int, len(ids))
	if len(ids) == 0 {
		return counts
	}

	fields := make([]string, 0, len(ids))
	for _, id := range ids {
		fields = append(fields, strconv.Itoa(id))
	}

	vals, err := rdb.HMGet(rctx, key, fields...).Result()
	if err != nil {
		slog.Warn("批量读取计数失败", "key", key, "err", err)
		return counts
	}
	for i, val := range vals {
		s, ok := val.(string)
		if !ok { // 字段不存在时是 nil
			continue
		}
		if n, err := strconv.Atoi(s); err == nil {
			counts[ids[i]] = n
		}
	}
	return counts
}

// hashCounts 的 ZSet 版本
func zsetCounts(rdb *redis.Client, key string, ids []int) map[int]int {
	counts := make(map[int]int, len(ids))
	if len(ids) == 0 {
		return counts
	}

	members := make([]string, 0, len(ids))
	for _, id := range ids {
		members = append(members, strconv.Itoa(id))
	}

	scores, err := rdb.ZMScore(rctx, key, members...).Result()
	if err != nil {
		slog.Warn("批量读取计数失败", "key", key, "err", err)
		return counts
	}
	for i, score := range scores {
		counts[ids[i]] = int(score) // 成员不存在时 score 为 0
	}
	return counts
}

// Page

// 将页面列表缓存到 Redis 中, 带 TTL
func addPageCache(rdb *redis.Client, pages []model.Page) error {
	data, err := json.Marshal(pages)
	if err != nil {
		return err
	}
	return rdb.Set(rctx, g.PAGE, string(data), cacheTTL).Err()
}

// 删除 Redis 中页面列表缓存
func removePageCache(rdb *redis.Client) error {
	return rdb.Del(rctx, g.PAGE).Err()
}

// 从 Redis 中获取页面列表缓存
// rdb.Get 如果不存在 key, 会返回 redis.Nil 错误
func getPageCache(rdb *redis.Client) (cache []model.Page, err error) {
	s, err := rdb.Get(rctx, g.PAGE).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(s), &cache); err != nil {
		return nil, err
	}

	return cache, nil
}

// Config

// 将博客配置缓存到 Redis 中, 带 TTL
//
// config 是 Hash, HMSet 不像 Set 那样能顺手带过期时间, 要单独 Expire 一次。
// 两条命令放进 pipeline: 分开发的话中间失败会留下一个永不过期的 key,
// 又退回到「只能手动 del」的状态
func addConfigCache(rdb *redis.Client, config map[string]string) error {
	// HMSET 不接受空的 field-value 列表, 配置表为空时直接跳过
	if len(config) == 0 {
		return nil
	}

	_, err := rdb.TxPipelined(rctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(rctx, g.CONFIG, config)
		pipe.Expire(rctx, g.CONFIG, cacheTTL)
		return nil
	})
	return err
}

// 删除 Redis 中博客配置缓存
func removeConfigCache(rdb *redis.Client) error {
	return rdb.Del(rctx, g.CONFIG).Err()
}

// 从 Redis 中获取博客配置缓存
// rdb.HGetAll 如果不存在 key, 不会返回 redis.Nil 错误, 而是返回空 map
func getConfigCache(rdb *redis.Client) (cache map[string]string, err error) {
	return rdb.HGetAll(rctx, g.CONFIG).Result()
}

// email
func SetMailInfo(rdb *redis.Client, info string, expire time.Duration) error {
	return rdb.Set(rctx, info, true, expire).Err()
}
func GetMailInfo(rdb *redis.Client, info string) (bool, error) {
	return rdb.Get(rctx, info).Bool()
}
func DeleteMailInfo(rdb *redis.Client, info string) error {
	return rdb.Del(rctx, info).Err()
}
