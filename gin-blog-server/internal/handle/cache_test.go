package handle

import (
	"gin-blog/internal/model"
	"testing"

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
