package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils/r"
)

type OperationLog struct{}

func (*OperationLog) GetList(req req.PageQuery) resp.PageResult[[]model.OperationLog] {
	list, total := operationLogDao.GetList(req)
	return resp.PageResult[[]model.OperationLog]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		List:     list,
		Total:    total,
	}
}

func (*OperationLog) Delete(ids []int) (code int) {
	dao.Delete(model.OperationLog{}, "id in ?", ids)
	return r.OK
}
