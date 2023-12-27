package handle

import (
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type Page struct{}

// 将页面列表缓存到 Redis 中
func addPageCache(rdb *redis.Client, pages []model.Page) error {
	data, err := json.Marshal(pages)
	if err != nil {
		return err
	}
	return rdb.Set(ctx(), g.PAGE, string(data), 0).Err()
}

// 删除 Redis 中页面列表缓存
func removePageCache(rdb *redis.Client) error {
	return rdb.Del(ctx(), g.PAGE).Err()
}

// 从 Redis 中获取页面列表缓存
func getPageCache(rdb *redis.Client) (cache []model.Page, err error) {
	s, err := rdb.Get(ctx(), g.PAGE).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(s), &cache); err != nil {
		return nil, err
	}

	return cache, nil
}

// @Summary 获取页面列表
// @Description 根据条件查询获取页面列表
// @Tags Page
// @Accept json
// @Produce json
// @Success 0 {object} Response[[]model.Page]
// @Security ApiKeyAuth
// @Router /page/list [get]
func (*Page) GetList(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	cache, err := getPageCache(rdb)
	if cache != nil && err == nil {
		slog.Debug("[handle-page-GetList] get page list from cache")
		ReturnSuccess(c, cache)
		return
	}

	switch err {
	case redis.Nil:
		break
	default:
		ReturnError(c, g.ERROR_REDIS_OPERATION, err)
		return
	}

	data, _, err := model.GetPageList(db)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	if err := addPageCache(GetRDB(c), data); err != nil {
		ReturnError(c, g.ERROR_REDIS_OPERATION, err)
		return
	}

	ReturnSuccess(c, data)
}

// @Summary 添加或修改页面
// @Description 添加或修改页面
// @Tags Page
// @Param form body AddOrEditPageReq true "添加或修改页面"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.Page]
// @Security ApiKeyAuth
// @Router /page [post]
func (*Page) SaveOrUpdate(c *gin.Context) {
	var req model.Page
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	page, err := model.SaveOrUpdatePage(db, req.ID, req.Name, req.Label, req.Cover)
	if err != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, err)
		return
	}

	if err := removePageCache(rdb); err != nil {
		ReturnError(c, g.ERROR_REDIS_OPERATION, err)
		return
	}

	ReturnSuccess(c, page)
}

// @Summary 删除页面
// @Description 根据 ID 数组删除页面
// @Tags Page
// @Param ids body []int true "页面 ID 数组"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /page [delete]
func (*Page) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ERROR_REQUEST_PARAM, err)
		return
	}

	result := GetDB(c).Delete(model.Page{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ERROR_DB_OPERATION, result.Error)
		return
	}

	if err := removePageCache(GetRDB(c)); err != nil {
		ReturnError(c, g.ERROR_REDIS_OPERATION, err)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
