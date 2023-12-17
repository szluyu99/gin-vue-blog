package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
)

type Tag struct{}

// 新增 / 修改
func (*Tag) SaveOrUpdate(req req.AddOrEditTag) int {
	// 检查标签名称是否已经存在
	exist := dao.GetOne(model.Tag{}, "name", req.Name)
	if !exist.IsEmpty() && exist.ID != req.ID {
		return r.ERROR_TAG_EXIST
	}

	tag := utils.CopyProperties[model.Tag](req) // vo -> po

	if req.ID != 0 {
		dao.Update(&tag, "name")
	} else {
		dao.Create(&tag)
	}
	return r.OK
}

// 删除标签
func (*Tag) Delete(ids []int) (code int) {
	// 检查标签下面有没有文章
	count := dao.Count(model.ArticleTag{}, "tag_id in ?", ids)
	if count > 0 {
		return r.ERROR_TAG_ART_EXIST
	}
	dao.Delete(model.Tag{}, "id in ?", ids)
	return r.OK
}

// 分页查询标签
func (*Tag) GetList(req req.PageQuery) resp.PageResult[[]resp.TagVO] {
	data, total := tagDao.GetList(req.PageNum, req.PageSize, req.Keyword)
	return resp.PageResult[[]resp.TagVO]{
		Total:    total,
		List:     data,
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
}

func (*Tag) GetOption() []resp.OptionVo {
	return tagDao.GetOption()
}

// 查询标签列表(前台)
func (*Tag) GetFrontList() resp.ListResult[[]resp.TagVO] {
	data, total := tagDao.GetList(0, 0, "")
	return resp.ListResult[[]resp.TagVO]{
		Total: total,
		List:  data,
	}
}
