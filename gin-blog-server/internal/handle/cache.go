package handle

import (
	"context"
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"time"
	"github.com/go-redis/redis/v9"
)

// redis context
var rctx = context.Background()

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
func SetMailInfo (rdb *redis.Client,info string,expire time.Duration) error{
	return rdb.Set(rctx,info,true,expire).Err()
}
func GetMailInfo (rdb *redis.Client,info string) (bool,error){
	return rdb.Get(rctx,info).Bool()
}
func DeleteMailInfo(rdb *redis.Client,info string) error{
	return rdb.Del(rctx,info).Err()
}