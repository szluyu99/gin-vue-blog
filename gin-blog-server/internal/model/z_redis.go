package model

import (
	"context"
	g "gin-blog/internal/global"
	"strconv"

	"github.com/go-redis/redis/v9"
)

// 获取文章点赞集合
func GetArticleLikeSet(rdb *redis.Client, id int) ([]string, error) {
	key := g.ARTICLE_LIKE_COUNT + strconv.Itoa(id)
	return rdb.SMembers(context.Background(), key).Result()
}
