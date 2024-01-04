package handle

import (
	"context"
	"gin-blog/internal/model"
	"log"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
)

// 需要 Redis 环境
func initRdb() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       11,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis 连接失败: ", err)
	}

	return rdb
}

func TestPageCache(t *testing.T) {
	rdb := initRdb()

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
	rdb := initRdb()

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
