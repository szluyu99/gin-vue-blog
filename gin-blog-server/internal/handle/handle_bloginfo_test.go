package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 后台首页统计: 文章只算已发布的, 访问量来自 Redis
func TestBlogInfoGetHomeInfo(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/home", (&BlogInfo{}).GetHomeInfo)

	assert.Nil(t, env.db.Create(&model.Article{Title: "公开", Status: model.STATUS_PUBLIC}).Error)
	assert.Nil(t, env.db.Create(&model.Article{Title: "草稿", Status: model.STATUS_DRAFT}).Error)
	assert.Nil(t, env.db.Create(&model.Article{Title: "回收站", Status: model.STATUS_PUBLIC, IsDelete: true}).Error)
	assert.Nil(t, env.db.Create(&model.UserInfo{Nickname: "u1"}).Error)
	_, err := model.SaveMessage(env.db, "张三", "", "留言", "", "", 10, true)
	assert.Nil(t, err)

	// Redis 里没有访问量时不能报错 (redis.Nil 要被忽略)
	resp := env.do(t, http.MethodGet, "/home", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var data BlogHomeVO
	decodeData(t, resp.Data, &data)
	assert.Equal(t, 1, data.ArticleCount)
	assert.Equal(t, 1, data.UserCount)
	assert.Equal(t, 1, data.MessageCount)
	assert.Zero(t, data.ViewCount)

	env.rdb.Set(rctx, g.VIEW_COUNT, 66, 0)
	decodeData(t, env.do(t, http.MethodGet, "/home", nil).Data, &data)
	assert.Equal(t, 66, data.ViewCount)
}

// 关于我: 读写同一个配置项
func TestBlogInfoAbout(t *testing.T) {
	env := newTestEnv(t)
	env.engine.GET("/setting/about", (&BlogInfo{}).GetAbout)
	env.engine.PUT("/setting/about", (&BlogInfo{}).UpdateAbout)

	// 还没设置过时是空字符串
	resp := env.do(t, http.MethodGet, "/setting/about", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Empty(t, resp.Data)

	resp = env.do(t, http.MethodPut, "/setting/about", map[string]any{"content": "# 关于我"})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, "# 关于我", resp.Data)

	resp = env.do(t, http.MethodGet, "/setting/about", nil)
	assert.Equal(t, "# 关于我", resp.Data)

	// 再改一次是覆盖而不是新增
	env.do(t, http.MethodPut, "/setting/about", map[string]any{"content": "改过了"})
	resp = env.do(t, http.MethodGet, "/setting/about", nil)
	assert.Equal(t, "改过了", resp.Data)

	var count int64
	env.db.Model(&model.Config{}).Where("`key` = ?", g.CONFIG_ABOUT).Count(&count)
	assert.Equal(t, int64(1), count)
}
