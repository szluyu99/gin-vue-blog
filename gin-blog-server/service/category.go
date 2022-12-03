package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
)

type Category struct{}

func (*Category) SaveOrUpdate(req req.AddOrEditCategory) int {
	// 查询分类是否存在
	existByName := dao.GetOne(model.Category{}, "name = ?", req.Name)
	// (同名存在) && (存在的id不等于当前要更新的id) -> 重复
	if !existByName.IsEmpty() && existByName.ID != req.ID {
		return r.ERROR_CATE_NAME_USED
	}

	category := utils.CopyProperties[model.Category](req) // vo -> po

	if req.ID != 0 {
		dao.Update(&category, "name")
	} else {
		dao.Create(&category)
	}
	return r.OK
}

func (*Category) Delete(ids []int) (code int) {
	// 查看分类下是否存在文章
	count := dao.Count(model.Article{}, "category_id in ?", ids)
	if count > 0 {
		return r.ERROR_CATE_ART_EXIST
	}
	dao.Delete(model.Category{}, "id in ?", ids)
	return r.OK
}

// 条件查询列表(后台)
func (*Category) GetList(req req.PageQuery) resp.PageResult[[]resp.CategoryVo] {
	data, total := categoryDao.GetList(req)
	return resp.PageResult[[]resp.CategoryVo]{
		Total:    total,
		List:     data,
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
}

func (*Category) GetOption() []resp.OptionVo {
	return categoryDao.GetOption()
}

// 查询列表(前台)
func (*Category) GetFrontList() resp.ListResult[[]resp.CategoryVo] {
	data, total := categoryDao.GetList(req.PageQuery{})
	return resp.ListResult[[]resp.CategoryVo]{
		Total: total,
		List:  data,
	}
}
