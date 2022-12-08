package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/utils"
	"gin-blog/utils/r"
)

type Page struct{}

func (*Page) GetList() []model.Page {
	list := make([]model.Page, 0)

	// 尝试从 Redis 中取, 取不到再查数据库
	listStr, err := utils.Redis.GetResult(KEY_PAGE)
	if listStr == "" || err != nil {
		list = dao.List([]model.Page{}, "id, name, label, cover, created_at, updated_at", "", "")
		utils.Redis.Set(KEY_PAGE, utils.Json.Marshal(list), 0)
	} else {
		utils.Json.Unmarshal(listStr, &list)
	}
	return list
}

func (*Page) SaveOrUpdate(req req.AddOrEditPage) (code int) {
	// 检查标签名称是否已经存在
	exist := dao.GetOne(model.Page{}, "name", req.Name)
	if exist.ID != 0 && exist.ID != req.ID {
		return r.ERROR_PAGE_NAME_EXIST
	}

	page := utils.CopyProperties[model.Page](req) // vo -> po
	if page.ID != 0 {
		dao.Update(&page)
	} else {
		dao.Create(&page)
	}
	utils.Redis.Del(KEY_PAGE)
	return r.OK
}

func (*Page) Delete(ids []int) (code int) {
	dao.Delete(model.Page{}, "id in ?", ids)
	utils.Redis.Del(KEY_PAGE)
	return r.OK
}
