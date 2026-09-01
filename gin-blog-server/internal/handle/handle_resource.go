package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Resource struct{}

type TreeOptionVO struct {
	ID       int            `json:"key"`
	Label    string         `json:"label"`
	Children []TreeOptionVO `json:"children"`
}

type ResourceTreeVO struct {
	ID        int              `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	Name      string           `json:"name"`
	Url       string           `json:"url"`
	Method    string           `json:"request_method"`
	Anonymous bool             `json:"is_anonymous"`
	Children  []ResourceTreeVO `json:"children"`
}

// TODO: 使用 oneof 标签校验数据
type AddOrEditResourceReq struct {
	ID       int    `json:"id"`
	Url      string `json:"url"`
	Method   string `json:"request_method"`
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
}

type EditAnonymousReq struct {
	ID        int  `json:"id" binding:"required"`
	Anonymous bool `json:"is_anonymous"`
}

// @Summary 获取资源列表（树形）
// @Description 获取接口资源列表, 支持关键字过滤
// @Tags Resource
// @Produce json
// @Param keyword query string false "关键字"
// @Success 0 {object} Response[[]ResourceTreeVO]
// @Security ApiKeyAuth
// @Router /resource/list [get]
func (*Resource) GetTreeList(c *gin.Context) {
	keyword := c.Query("keyword")

	resourceList, err := model.GetResourceList(GetDB(c), keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, resources2ResourceVos(resourceList))
}

// @Summary 获取资源选项（树形）
// @Description 用于角色分配资源
// @Tags Resource
// @Produce json
// @Success 0 {object} Response[[]TreeOptionVO]
// @Security ApiKeyAuth
// @Router /resource/option [get]
func (*Resource) GetOption(c *gin.Context) {
	result := make([]TreeOptionVO, 0)

	db := GetDB(c)
	resources, err := model.GetResourceList(db, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)

	for _, item := range parentList {
		var children []TreeOptionVO
		for _, re := range childrenMap[item.ID] {
			children = append(children, TreeOptionVO{
				ID:    re.ID,
				Label: re.Name,
			})
		}
		result = append(result, TreeOptionVO{
			ID:       item.ID,
			Label:    item.Name,
			Children: children,
		})
	}
	ReturnSuccess(c, result)
}

// @Summary 新增或编辑资源
// @Description 新增或编辑接口资源
// @Tags Resource
// @Accept json
// @Produce json
// @Param form body AddOrEditResourceReq true "新增或编辑资源"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /resource [post]
func (*Resource) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditResourceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	err := model.SaveOrUpdateResource(db, req.ID, req.ParentId, req.Name, req.Url, req.Method)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// @Summary 修改资源匿名访问
// @Description 修改资源是否允许匿名访问
// @Tags Resource
// @Accept json
// @Produce json
// @Param form body EditAnonymousReq true "修改匿名访问"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /resource/anonymous [put]
func (*Resource) UpdateAnonymous(c *gin.Context) {
	var req EditAnonymousReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.UpdateResourceAnonymous(GetDB(c), req.ID, req.Anonymous)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// TODO: 考虑删除模块后, 其子资源怎么办? 目前做法是有子资源无法删除
// TODO: 强制删除?
// @Summary 删除资源
// @Description 被角色使用或存在子资源时不允许删除
// @Tags Resource
// @Produce json
// @Param id path int true "资源 ID"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /resource/{id} [delete]
func (*Resource) Delete(c *gin.Context) {
	resourceId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// 检查该资源是否被角色使用
	use, _ := model.CheckResourceInUse(db, resourceId)
	if use {
		ReturnError(c, g.ErrResourceUsedByRole, nil)
		return
	}

	// 获取该资源
	resource, err := model.GetResourceById(db, resourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrResourceNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 如果作为模块, 检查模块下是否有子资源
	if resource.ParentId == 0 {
		hasChild, _ := model.CheckResourceHasChild(db, resourceId)
		if hasChild {
			ReturnError(c, g.ErrResourceHasChildren, nil)
			return
		}
	}

	rows, err := model.DeleteResource(db, resourceId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// []Resource => []ResourceVO
func resources2ResourceVos(resources []model.Resource) []ResourceTreeVO {
	list := make([]ResourceTreeVO, 0)
	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)
	for _, item := range parentList {
		resourceVO := resource2ResourceVo(item)
		resourceVO.Children = make([]ResourceTreeVO, 0)
		for _, child := range childrenMap[item.ID] {
			resourceVO.Children = append(resourceVO.Children, resource2ResourceVo(child))
		}
		list = append(list, resourceVO)
	}
	return list
}

// Resource => ResourceVO
func resource2ResourceVo(r model.Resource) ResourceTreeVO {
	return ResourceTreeVO{
		ID:        r.ID,
		Name:      r.Name,
		Url:       r.Url,
		Method:    r.Method,
		Anonymous: r.Anonymous,
		CreatedAt: r.CreatedAt,
	}
}

// 存储每个节点对应 [子资源列表] 的 map
// key: resourceId
// value: childrenList
func getChildrenMap(resources []model.Resource) map[int][]model.Resource {
	m := make(map[int][]model.Resource)
	for _, r := range resources {
		if r.ParentId != 0 {
			m[r.ParentId] = append(m[r.ParentId], r)
		}
	}
	return m
}

// 获取一级资源 (parent_id == 0)
func getModuleList(resources []model.Resource) []model.Resource {
	list := make([]model.Resource, 0)
	for _, r := range resources {
		if r.ParentId == 0 {
			list = append(list, r)
		}
	}
	return list
}
