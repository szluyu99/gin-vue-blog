package ginblog

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

/*
Redis 计数的落库与回填

点赞数 / 浏览数 / 站点访问量 / 访客地域只在 Redis 里累加, 数据库没有对应字段,
Redis 就是唯一数据源。清一次 Redis 或者换台机器部署, 这些数字就归零了。
这里定期把它们抄一份进 counter_snapshot 表, 启动时再把缺的补回 Redis。

按天的访问量(view_count:2006-01-02)不在备份范围内: 那是趋势图的原料, 只留 30 天,
丢一天不影响任何业务判断, 备份反而要额外处理"哪些天已经过期了"。

**三种键的类型不一样, 不能一把抓**(这点踩过: 一开始按 ZSet 统一处理, 端到端时
`ZSCORE article_like_count` 直接报 WRONGTYPE):
  - Hash:   文章点赞数、评论点赞数、访客地域
  - ZSet:   文章浏览数
  - String: 站点访问量
*/
var (
	counterHashes  = []string{g.ARTICLE_LIKE_COUNT, g.COMMENT_LIKE_COUNT, g.VISITOR_AREA}
	counterZSets   = []string{g.ARTICLE_VIEW_COUNT}
	counterStrings = []string{g.VIEW_COUNT}
)

// FlushCounters 把 Redis 里的计数抄一份进数据库
func FlushCounters(ctx context.Context, rdb *redis.Client, db *gorm.DB) error {
	items := make([]model.CounterSnapshot, 0, 64)

	for _, key := range counterHashes {
		values, err := rdb.HGetAll(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		for member, raw := range values {
			n, convErr := strconv.ParseInt(raw, 10, 64)
			if convErr != nil {
				// 单个成员坏掉不该让整次备份失败
				slog.Warn("计数不是数字, 跳过备份", "key", key, "member", member, "value", raw)
				continue
			}
			items = append(items, model.CounterSnapshot{Key: key, Member: member, Count: n})
		}
	}

	for _, key := range counterZSets {
		members, err := rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		for _, m := range members {
			member, ok := m.Member.(string)
			if !ok {
				continue
			}
			items = append(items, model.CounterSnapshot{Key: key, Member: member, Count: int64(m.Score)})
		}
	}

	for _, key := range counterStrings {
		n, err := rdb.Get(ctx, key).Int64()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return err
		}
		// 没有成员, Member 留空串
		items = append(items, model.CounterSnapshot{Key: key, Count: n})
	}

	if err := model.SaveCounterSnapshots(db, items); err != nil {
		return err
	}
	slog.Debug("Redis 计数已落库", "rows", len(items))
	return nil
}

/*
RestoreCounters 用数据库里的备份补回 Redis

**只补 Redis 里不存在的键**。Redis 还在的时候它的值一定比备份新(备份最多是上个周期的
快照), 覆盖回去等于把这段时间的点赞和浏览抹掉。所以判断的是"键在不在",
而不是逐个成员比大小 —— 键存在就说明这份数据是活的, 一个成员都不动。
*/
func RestoreCounters(ctx context.Context, rdb *redis.Client, db *gorm.DB) error {
	for _, key := range counterHashes {
		rows, err := missingKeyRows(ctx, rdb, db, key)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			continue
		}

		values := make(map[string]any, len(rows))
		for _, row := range rows {
			values[row.Member] = row.Count
		}
		if err := rdb.HSet(ctx, key, values).Err(); err != nil {
			return err
		}
		slog.Info("从备份恢复计数", "key", key, "members", len(values))
	}

	for _, key := range counterZSets {
		rows, err := missingKeyRows(ctx, rdb, db, key)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			continue
		}

		members := make([]redis.Z, 0, len(rows))
		for _, row := range rows {
			members = append(members, redis.Z{Member: row.Member, Score: float64(row.Count)})
		}
		if err := rdb.ZAdd(ctx, key, members...).Err(); err != nil {
			return err
		}
		slog.Info("从备份恢复计数", "key", key, "members", len(members))
	}

	for _, key := range counterStrings {
		rows, err := missingKeyRows(ctx, rdb, db, key)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			continue
		}
		if err := rdb.Set(ctx, key, rows[0].Count, 0).Err(); err != nil {
			return err
		}
		slog.Info("从备份恢复计数", "key", key, "count", rows[0].Count)
	}

	return nil
}

// 键已经在 Redis 里就返回空: 活数据优先, 一律不覆盖
func missingKeyRows(ctx context.Context, rdb *redis.Client, db *gorm.DB, key string) ([]model.CounterSnapshot, error) {
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, nil
	}
	return model.GetCounterSnapshots(db, key)
}
