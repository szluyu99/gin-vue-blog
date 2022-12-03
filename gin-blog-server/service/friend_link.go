package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
)

type FriendLink struct{}

func (*FriendLink) GetList(req req.PageQuery) (data resp.PageResult[[]model.FriendLink]) {
	list, total := friendLinkDao.GetList(req)
	return resp.PageResult[[]model.FriendLink]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    total,
		List:     list,
	}
}

func (*FriendLink) SaveOrUpdate(req req.SaveOrUpdateLink) (code int) {
	friendLink := utils.CopyProperties[model.FriendLink](req) // vo -> po
	if friendLink.ID != 0 {
		dao.Update(&friendLink)
	} else { // 新增
		dao.Create(&friendLink)
	}
	return r.OK
}

func (*FriendLink) Delete(ids []int) (code int) {
	dao.Delete(model.FriendLink{}, "id in ?", ids)
	return r.OK
}

func (*FriendLink) GetFrontList() (data []model.FriendLink) {
	return dao.List([]model.FriendLink{}, "*", "", "")
}
