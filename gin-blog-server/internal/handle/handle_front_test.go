package handle

import (
	g "gin-blog/internal/global"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 点赞文章: 同一用户重复点赞是取消, 计数与用户集合要同步变化
func TestFrontLikeArticle(t *testing.T) {
	env := newTestEnv(t)
	user := env.loginAs(7, "tester")
	env.engine.GET("/front/article/like/:article_id", (&Front{}).LikeArticle)

	likeSetKey := g.ARTICLE_USER_LIKE_SET + strconv.Itoa(user.ID)

	// 第一次: 点赞
	resp := env.do(t, http.MethodGet, "/front/article/like/1", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.True(t, env.rdb.SIsMember(rctx, likeSetKey, 1).Val())
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "1").Val())

	// 第二次: 取消点赞
	resp = env.do(t, http.MethodGet, "/front/article/like/1", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.False(t, env.rdb.SIsMember(rctx, likeSetKey, 1).Val())
	assert.Equal(t, "0", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "1").Val())

	// 不同文章互不影响
	env.do(t, http.MethodGet, "/front/article/like/2", nil)
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "2").Val())
	assert.Equal(t, "0", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "1").Val())

	// 非法的文章 id
	resp = env.do(t, http.MethodGet, "/front/article/like/abc", nil)
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)
}

// 点赞评论: 与文章点赞是两套独立的 key
func TestFrontLikeComment(t *testing.T) {
	env := newTestEnv(t)
	user := env.loginAs(7, "tester")
	env.engine.GET("/front/comment/like/:comment_id", (&Front{}).LikeComment)

	likeSetKey := g.COMMENT_USER_LIKE_SET + strconv.Itoa(user.ID)

	resp := env.do(t, http.MethodGet, "/front/comment/like/3", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.True(t, env.rdb.SIsMember(rctx, likeSetKey, 3).Val())
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.COMMENT_LIKE_COUNT, "3").Val())

	// 文章点赞的计数不应受影响
	assert.Empty(t, env.rdb.HGetAll(rctx, g.ARTICLE_LIKE_COUNT).Val())

	// 取消点赞
	env.do(t, http.MethodGet, "/front/comment/like/3", nil)
	assert.False(t, env.rdb.SIsMember(rctx, likeSetKey, 3).Val())
	assert.Equal(t, "0", env.rdb.HGet(rctx, g.COMMENT_LIKE_COUNT, "3").Val())
}

// 不同用户的点赞记录互相隔离
func TestFrontLikeArticleMultiUser(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/front/article/like/:article_id", (&Front{}).LikeArticle)

	env.loginAs(1, "user1")
	env.do(t, http.MethodGet, "/front/article/like/9", nil)

	env.loginAs(2, "user2")
	env.do(t, http.MethodGet, "/front/article/like/9", nil)

	// 两个用户各自点赞, 计数为 2
	assert.Equal(t, "2", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "9").Val())
	assert.True(t, env.rdb.SIsMember(rctx, g.ARTICLE_USER_LIKE_SET+"1", 9).Val())
	assert.True(t, env.rdb.SIsMember(rctx, g.ARTICLE_USER_LIKE_SET+"2", 9).Val())

	// user2 取消, 只影响自己
	env.do(t, http.MethodGet, "/front/article/like/9", nil)
	assert.Equal(t, "1", env.rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, "9").Val())
	assert.True(t, env.rdb.SIsMember(rctx, g.ARTICLE_USER_LIKE_SET+"1", 9).Val())
	assert.False(t, env.rdb.SIsMember(rctx, g.ARTICLE_USER_LIKE_SET+"2", 9).Val())
}

// 未登录访问需要登录的接口: 应返回业务错误而不是 panic 成 500
// (JWTAuth 对资源表中不存在的接口会跳过鉴权, 因此 handler 必须自己兜住)
func TestFrontLikeWithoutLogin(t *testing.T) {
	env := newTestEnv(t) // 不调用 loginAs
	env.engine.GET("/front/article/like/:article_id", (&Front{}).LikeArticle)
	env.engine.GET("/front/comment/like/:comment_id", (&Front{}).LikeComment)

	resp := env.do(t, http.MethodGet, "/front/article/like/1", nil)
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)

	resp = env.do(t, http.MethodGet, "/front/comment/like/1", nil)
	assert.Equal(t, g.ErrTokenNotExist.Code(), resp.Code)

	// 未登录不应产生任何点赞数据
	assert.Empty(t, env.rdb.HGetAll(rctx, g.ARTICLE_LIKE_COUNT).Val())
	assert.Empty(t, env.rdb.HGetAll(rctx, g.COMMENT_LIKE_COUNT).Val())
}
