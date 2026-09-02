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

// 将页面列表缓存到 Redis 中
func addPageCache(rdb *redis.Client, pages []model.Page) error {
	data, err := json.Marshal(pages)
	if err != nil {
		return err
	}
	return rdb.Set(rctx, g.PAGE, string(data), 0).Err()
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

// 将博客配置缓存到 Redis 中
func addConfigCache(rdb *redis.Client, config map[string]string) error {
	// HMSET 不接受空的 field-value 列表, 配置表为空时直接跳过
	if len(config) == 0 {
		return nil
	}
	return rdb.HMSet(rctx, g.CONFIG, config).Err()
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
