package ginblog

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// sqlite 内存库 + miniredis, 和 handle 包的脚手架一套路子
func newCounterEnv(t *testing.T) (*gorm.DB, *redis.Client, *miniredis.Miniredis) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	assert.Nil(t, err)
	assert.Nil(t, model.MakeMigrate(db))

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { rdb.Close() })

	return db, rdb, mr
}

func TestFlushCounters(t *testing.T) {
	db, rdb, _ := newCounterEnv(t)
	ctx := context.Background()

	// 各键的真实类型: 点赞是 Hash, 浏览是 ZSet, 站点访问量是普通计数
	rdb.HIncrBy(ctx, g.ARTICLE_LIKE_COUNT, "1", 3)
	rdb.ZIncrBy(ctx, g.ARTICLE_VIEW_COUNT, 12, "1")
	rdb.HIncrBy(ctx, g.COMMENT_LIKE_COUNT, "7", 5)
	rdb.Set(ctx, g.VIEW_COUNT, 66, 0)
	rdb.HSet(ctx, g.VISITOR_AREA, "江苏", 4, "未知", 1)

	assert.Nil(t, FlushCounters(ctx, rdb, db))

	like, err := model.GetCounterSnapshots(db, g.ARTICLE_LIKE_COUNT)
	assert.Nil(t, err)
	assert.Equal(t, []model.CounterSnapshot{{Key: g.ARTICLE_LIKE_COUNT, Member: "1", Count: 3}}, like)

	view, _ := model.GetCounterSnapshots(db, g.ARTICLE_VIEW_COUNT)
	assert.Equal(t, int64(12), view[0].Count)

	commentLike, _ := model.GetCounterSnapshots(db, g.COMMENT_LIKE_COUNT)
	assert.Equal(t, int64(5), commentLike[0].Count)

	total, _ := model.GetCounterSnapshots(db, g.VIEW_COUNT)
	assert.Len(t, total, 1)
	assert.Equal(t, int64(66), total[0].Count)
	assert.Empty(t, total[0].Member, "普通计数没有成员")

	area, _ := model.GetCounterSnapshots(db, g.VISITOR_AREA)
	assert.Len(t, area, 2)
}

// 反复落库是覆盖同一行, 不是不断堆新行
func TestFlushCountersIsUpsert(t *testing.T) {
	db, rdb, _ := newCounterEnv(t)
	ctx := context.Background()

	rdb.ZIncrBy(ctx, g.ARTICLE_VIEW_COUNT, 1, "1")
	assert.Nil(t, FlushCounters(ctx, rdb, db))

	rdb.ZIncrBy(ctx, g.ARTICLE_VIEW_COUNT, 9, "1")
	assert.Nil(t, FlushCounters(ctx, rdb, db))

	rows, _ := model.GetCounterSnapshots(db, g.ARTICLE_VIEW_COUNT)
	assert.Len(t, rows, 1, "同一个 (key, member) 只该有一行")
	assert.Equal(t, int64(10), rows[0].Count)
}

// Redis 被清空后, 启动时要能把计数补回来
func TestRestoreCountersFillsEmptyRedis(t *testing.T) {
	db, rdb, mr := newCounterEnv(t)
	ctx := context.Background()

	rdb.HIncrBy(ctx, g.ARTICLE_LIKE_COUNT, "1", 3)
	rdb.ZIncrBy(ctx, g.ARTICLE_VIEW_COUNT, 12, "2")
	rdb.Set(ctx, g.VIEW_COUNT, 66, 0)
	rdb.HSet(ctx, g.VISITOR_AREA, "江苏", 4)
	assert.Nil(t, FlushCounters(ctx, rdb, db))

	mr.FlushAll() // 模拟清了 Redis / 换台机器部署

	assert.Nil(t, RestoreCounters(ctx, rdb, db))

	// 恢复回去的类型也要对: 点赞是 Hash, 浏览是 ZSet
	assert.Equal(t, "3", rdb.HGet(ctx, g.ARTICLE_LIKE_COUNT, "1").Val())
	assert.Equal(t, float64(12), rdb.ZScore(ctx, g.ARTICLE_VIEW_COUNT, "2").Val())
	assert.Equal(t, "66", rdb.Get(ctx, g.VIEW_COUNT).Val())
	assert.Equal(t, "4", rdb.HGet(ctx, g.VISITOR_AREA, "江苏").Val())
}

/*
Redis 里还有数据时一个字节都不许动

备份最多是上个周期的快照, 覆盖回去等于把这段时间新增的点赞和浏览抹掉。
所以判断的是"键在不在", 键存在就整个跳过。
*/
func TestRestoreCountersNeverOverwritesLiveData(t *testing.T) {
	db, rdb, _ := newCounterEnv(t)
	ctx := context.Background()

	// 备份里是旧值
	rdb.ZIncrBy(ctx, g.ARTICLE_VIEW_COUNT, 5, "1")
	rdb.HIncrBy(ctx, g.ARTICLE_LIKE_COUNT, "1", 2)
	rdb.Set(ctx, g.VIEW_COUNT, 10, 0)
	rdb.HSet(ctx, g.VISITOR_AREA, "江苏", 1)
	assert.Nil(t, FlushCounters(ctx, rdb, db))

	// Redis 又涨了, 此时备份是落后的
	rdb.ZIncrBy(ctx, g.ARTICLE_VIEW_COUNT, 20, "1")
	rdb.HIncrBy(ctx, g.ARTICLE_LIKE_COUNT, "1", 6)
	rdb.Set(ctx, g.VIEW_COUNT, 99, 0)
	rdb.HSet(ctx, g.VISITOR_AREA, "江苏", 8)

	assert.Nil(t, RestoreCounters(ctx, rdb, db))

	assert.Equal(t, float64(25), rdb.ZScore(ctx, g.ARTICLE_VIEW_COUNT, "1").Val(), "不能被旧备份覆盖")
	assert.Equal(t, "8", rdb.HGet(ctx, g.ARTICLE_LIKE_COUNT, "1").Val())
	assert.Equal(t, "99", rdb.Get(ctx, g.VIEW_COUNT).Val())
	assert.Equal(t, "8", rdb.HGet(ctx, g.VISITOR_AREA, "江苏").Val())
}

// 空库空 Redis 都不该报错(首次启动就是这个状态)
func TestCountersEmptyIsNoop(t *testing.T) {
	db, rdb, _ := newCounterEnv(t)
	ctx := context.Background()

	assert.Nil(t, FlushCounters(ctx, rdb, db))
	assert.Nil(t, RestoreCounters(ctx, rdb, db))

	var count int64
	db.Model(&model.CounterSnapshot{}).Count(&count)
	assert.Zero(t, count)
}
