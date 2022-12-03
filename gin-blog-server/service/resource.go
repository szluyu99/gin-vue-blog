package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
)

type Resource struct{}

// 获取资源列表(树形)
func (s *Resource) GetTreeList(req req.PageQuery) []resp.ResourceVo {
	// 根据关键字查询出[资源列表](非树形)
	resources := resourceDao.GetListByKeyword(req.Keyword)
	// []resource -> []resourceVo
	return s.resources2ResourceVos(resources)
}

// 获取数据选项(树形)
func (s *Resource) GetOptionList() []resp.TreeOptionVo {
	resList := make([]resp.TreeOptionVo, 0)

	resources := dao.List([]model.Resource{}, "id, name, parent_id", "", "is_anonymous = ?", 0)
	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)
	for _, item := range parentList {
		// 构造 Children
		var childrenOptionVos []resp.TreeOptionVo
		for _, re := range childrenMap[item.ID] {
			childrenOptionVos = append(childrenOptionVos, resp.TreeOptionVo{
				ID:    re.ID,
				Label: re.Name,
			})
		}
		resList = append(resList, resp.TreeOptionVo{
			ID:       item.ID,
			Label:    item.Name,
			Children: childrenOptionVos,
		})
	}
	return resList
}

// 新增或编辑资源, 关联更新 casbin_rule 中数据
func (*Resource) SaveOrUpdate(req req.SaveOrUpdateResource) (code int) {
	// 检查该资源名已经存在
	existByName := dao.GetOne(model.Resource{}, "name", req.Name)
	if existByName.ID != 0 && existByName.ID != req.ID {
		return r.ERROR_RESOURCE_NAME_EXIST
	}

	if req.ID != 0 { // 更新
		oldResource := dao.GetOne(model.Resource{}, "id", req.ID)
		dao.UpdatesMap(&model.Resource{}, utils.Struct2Map(req), "id", req.ID) // map 可以更新零值
		// ! 关联更新 casbin_rule 中的信息
		utils.Casbin.UpdatePolicy(
			[]string{"", oldResource.Url, oldResource.RequestMethod},
			[]string{"", req.Url, req.RequestMethod},
		)
	} else { // 新增
		data := utils.CopyProperties[model.Resource](req)
		dao.Create(&data)
		// * 解决前端的 BUG: 级联选中某个父节点后, 新增的子节点默认会展示被选中, 实际上未被选中值
		// * 解决方案: 新增子节点后, 删除该节点对应的父节点与角色的关联关系
		dao.Delete(model.RoleResource{}, "resource_id", data.ParentId)
	}
	return r.OK
}

// 编辑资源的匿名访问, 关联更新 casbin_rule 中数据
func (*Resource) UpdateAnonymous(req req.UpdateAnonymous) (code int) {
	// 检查要更新的资源是否存在
	existById := dao.GetOne(model.Resource{}, "id", req.ID)
	if existById.ID == 0 {
		return r.ERROR_RESOURCE_NOT_EXIST
	}
	// 只更新 is_anonymous 字段
	dao.UpdatesMap(&model.Resource{}, map[string]any{"is_anonymous": *req.IsAnonymous}, "id", req.ID)
	// ! 关联处理 casbin_rule 中的 isAnonymous
	if *req.IsAnonymous == 0 {
		// fmt.Println("删除 casbin_rule anonymous")
		utils.Casbin.DeletePermissionForRole("anonymous", req.Url, req.RequestMethod)
	} else {
		utils.Casbin.AddPermissionForRole("anonymous", req.Url, req.RequestMethod)
	}
	return r.OK
}

// 删除资源, 关联删除 casbin_rule 中数据
// TODO: 考虑删除模块后, 其子资源怎么办? 目前做法是有子资源无法删除
func (*Resource) Delete(resourceId int) (code int) {
	// 检查要删除的资源是否存在
	existResourceById := dao.GetOne(model.Resource{}, "id", resourceId)
	if existResourceById.ID == 0 {
		return r.ERROR_RESOURCE_NOT_EXIST
	}
	// * 检查 role_resource 下是否有数据
	existRoleResource := dao.GetOne(model.RoleResource{}, "resource_id", resourceId)
	if existRoleResource.ResourceId != 0 {
		return r.ERROR_RESOURCE_USED_BY_ROLE
	}
	// * 如果该 resource 是模块, 检查其是否有子资源
	if existResourceById.ParentId == 0 {
		if dao.Count(model.Resource{}, "parent_id = ?", resourceId) != 0 {
			return r.ERROR_RESOURCE_HAS_CHILDREN
		}
	}
	// 删除资源
	dao.Delete(model.Resource{}, "id", resourceId)
	// ! 关联删除 casbin_rule 中的数据
	// 因为前面检查过 role_resource 下是否有数据, 理论上来说不会真正到关联删除这步
	if existResourceById.Url != "" && existResourceById.RequestMethod != "" {
		utils.Casbin.DeletePermission(existResourceById.Url, existResourceById.RequestMethod)
	}
	return r.OK
}

// []model.Resourec -> []resp.ResourceVO
func (s *Resource) resources2ResourceVos(resources []model.Resource) []resp.ResourceVo {
	resList := make([]resp.ResourceVo, 0)
	// 找出父节点列表 (parentId == 0)
	parentList := getModuleList(resources)
	// 存储每个节点对应[子资源列表]的 map
	childrenMap := getChildrenMap(resources)
	// []model.Resource -> []resp.ResourceVo
	for _, item := range parentList {
		resourceVo := s.resource2ResourceVo(item) // vo -> po
		resourceVo.Children = make([]resp.ResourceVo, 0)
		// 构建 Children 数据
		for _, child := range childrenMap[item.ID] {
			resourceVo.Children = append(resourceVo.Children, s.resource2ResourceVo(child))
		}
		resList = append(resList, resourceVo)
	}
	return resList
}

// model.Resource -> resp.ResourceVo
func (*Resource) resource2ResourceVo(r model.Resource) resp.ResourceVo {
	return resp.ResourceVo{
		ID:            r.ID,
		Name:          r.Name,
		Url:           r.Url,
		RequestMethod: r.RequestMethod,
		IsAnonymous:   r.IsAnonymous,
		CreatedAt:     r.CreatedAt,
		// Children:      make([]resp.ResourceVo, 0),
	}
}

// 存储每个节点对应[子资源列表]的 map: key 为 resourceId, value 为其对应的 childrenList
func getChildrenMap(resources []model.Resource) map[int][]model.Resource {
	childrenMap := make(map[int][]model.Resource)
	for _, r := range resources {
		if r.ParentId != 0 { // 有父节点
			childrenMap[r.ParentId] = append(childrenMap[r.ParentId], r)
		}
	}
	return childrenMap
}

// 获取一级资源 (没有 parent_id 的 resource)
func getModuleList(resources []model.Resource) (list []model.Resource) {
	var parentList []model.Resource
	for _, r := range resources {
		if r.ParentId == 0 {
			parentList = append(parentList, r)
		}
	}
	return parentList
}
