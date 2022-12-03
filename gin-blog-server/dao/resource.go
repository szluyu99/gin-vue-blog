package dao

import (
	"gin-blog/model"
)

type Resource struct{}

// 根据 [资源id列表] 获取 [资源列表]
func (*Resource) GetListByIds(ids []int) (list []model.Resource) {
	return List([]model.Resource{}, "url, request_method", "", "id in ?", ids)
}

// 根据关键字获取资源列表(非树形)
func (*Resource) GetListByKeyword(keyword string) (list []model.Resource) {
	if keyword != "" {
		list = List([]model.Resource{}, "*", "", "name like ?", "%"+keyword+"%")
	} else {
		list = List([]model.Resource{}, "*", "", "")
	}
	return
}
