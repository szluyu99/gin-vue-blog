package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryAPI(t *testing.T) {
	env := newTestEnv(t)
	api := Category{}
	env.engine.GET("/category/list", api.GetList)
	env.engine.GET("/category/option", api.GetOption)
	env.engine.POST("/category", api.SaveOrUpdate)
	env.engine.DELETE("/category", api.Delete)

	// 新增
	resp := env.do(t, http.MethodPost, "/category", map[string]any{"name": "后端"})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var created model.Category
	decodeData(t, resp.Data, &created)
	assert.Equal(t, "后端", created.Name)
	assert.NotZero(t, created.ID)

	// 列表
	resp = env.do(t, http.MethodGet, "/category/list?page_num=1&page_size=10", nil)
	assert.Equal(t, g.SUCCESS, resp.Code)

	var page PageResult[model.CategoryVO]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	assert.Len(t, page.List, 1)

	// 关键字过滤: 不匹配时为空
	resp = env.do(t, http.MethodGet, "/category/list?keyword=不存在", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(0), page.Total)

	// 选项列表
	resp = env.do(t, http.MethodGet, "/category/option", nil)
	var options []model.OptionVO
	decodeData(t, resp.Data, &options)
	assert.Len(t, options, 1)

	// 缺少必填字段
	resp = env.do(t, http.MethodPost, "/category", map[string]any{})
	assert.Equal(t, g.ErrRequest.Code(), resp.Code)

	// 分类下有文章时不允许删除
	env.db.Create(&model.Article{Title: "文章", CategoryId: created.ID})
	resp = env.do(t, http.MethodDelete, "/category", []int{created.ID})
	assert.Equal(t, g.ErrCateHasArt.Code(), resp.Code)

	// 清掉文章后可以删除
	env.db.Where("category_id = ?", created.ID).Delete(&model.Article{})
	resp = env.do(t, http.MethodDelete, "/category", []int{created.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)
	assert.Equal(t, float64(1), resp.Data)
}

func TestTagAPI(t *testing.T) {
	env := newTestEnv(t)
	api := Tag{}
	env.engine.GET("/tag/list", api.GetList)
	env.engine.POST("/tag", api.SaveOrUpdate)
	env.engine.DELETE("/tag", api.Delete)

	resp := env.do(t, http.MethodPost, "/tag", map[string]any{"name": "Go"})
	assert.Equal(t, g.SUCCESS, resp.Code)

	var created model.Tag
	decodeData(t, resp.Data, &created)
	assert.Equal(t, "Go", created.Name)

	resp = env.do(t, http.MethodGet, "/tag/list?page_num=1&page_size=10", nil)
	var page PageResult[model.TagVO]
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)

	resp = env.do(t, http.MethodDelete, "/tag", []int{created.ID})
	assert.Equal(t, g.SUCCESS, resp.Code)

	resp = env.do(t, http.MethodGet, "/tag/list", nil)
	decodeData(t, resp.Data, &page)
	assert.Equal(t, int64(0), page.Total)
}
