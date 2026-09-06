package handle

import (
	"gin-blog/internal/model"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// 使用 miniredis 提供内存中的 Redis, 不依赖本地环境
func initRdb(t *testing.T) *redis.Client {
	t.Helper()
	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

// 同时拿到 miniredis 句柄, 用来快进时间验证 TTL
func initRdbWithMini(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()}), mr
}

func TestPageCache(t *testing.T) {
	rdb := initRdb(t)

	pages := []model.Page{
		{Name: "page1"},
		{Name: "page2"},
	}

	// 直接获取缓存
	// 不存在, 返回 redis.Nil 错误
	{
		cache, err := getPageCache(rdb)
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, cache)
	}

	// 新增, 获取 缓存
	{
		err := addPageCache(rdb, pages)
		assert.Nil(t, err)

		cache, err := getPageCache(rdb)
		assert.Nil(t, err)
		assert.Equal(t, pages, cache)
	}

	// 删除, 获取 缓存
	// 不存在, 返回 redis.Nil 错误
	{
		err := removePageCache(rdb)
		assert.Nil(t, err)

		cache, err := getPageCache(rdb)
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, cache)
	}

}

/*
读穿缓存必须带 TTL

主动失效只覆盖走接口的写操作; 直接 UPDATE 数据库、跑 generate-data 灌种子
都不会触发, 以前这类改动只能手动 redis-cli del 再重启后端才能看到。
TTL 就是这条兜底路径, 所以两个缓存都要断言「会过期」。
*/
func TestCacheHasTTL(t *testing.T) {
	rdb, mr := initRdbWithMini(t)

	assert.Nil(t, addPageCache(rdb, []model.Page{{Name: "page1"}}))
	assert.Nil(t, addConfigCache(rdb, map[string]string{"name": "name"}))

	// string 和 hash 两种类型都要真的带上过期时间
	pageTTL, err := rdb.TTL(rctx, "page").Result()
	assert.Nil(t, err)
	assert.Greater(t, pageTTL, time.Duration(0), "page 缓存要有 TTL")
	configTTL, err := rdb.TTL(rctx, "config").Result()
	assert.Nil(t, err)
	assert.Greater(t, configTTL, time.Duration(0), "config 是 Hash, HMSet 不会自带 TTL, 要额外 Expire")

	// 快进到过期之后: 两个缓存都要变成「未命中」, 让请求回落到数据库
	mr.FastForward(cacheTTL + time.Second)

	cache, err := getPageCache(rdb)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, cache)

	configCache, err := getConfigCache(rdb)
	assert.Nil(t, err)
	assert.Empty(t, configCache, "config 过期后 HGetAll 返回空 map, 上层会重新查库")
}

func TestConfigCache(t *testing.T) {
	rdb := initRdb(t)

	config := map[string]string{
		"name": "name",
		"url":  "url",
	}

	// 直接获取缓存
	// 不存在, 返回空 map
	{
		cache, err := getConfigCache(rdb)
		assert.Nil(t, err)
		assert.Empty(t, cache)
	}

	// 新增, 获取 缓存
	{
		err := addConfigCache(rdb, config)
		assert.Nil(t, err)

		cache, err := getConfigCache(rdb)
		assert.Nil(t, err)
		assert.Equal(t, config, cache)
	}

	// 删除, 获取 缓存
	// 不存在, 返回空 map
	{
		err := removeConfigCache(rdb)
		assert.Nil(t, err)

		cache, err := getConfigCache(rdb)
		assert.Nil(t, err)
		assert.Empty(t, cache)
	}
}
