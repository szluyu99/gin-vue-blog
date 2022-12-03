package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"

	"github.com/gin-gonic/gin"
)

type Message struct{}

// 保存时, 获取发送留言的 IP 地址
func (*Message) Save(c *gin.Context, req req.AddMessage) {
	// 通过 ip 工具获取 ip地址 和 ip来源
	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSource(ipAddress)

	// vo -> po
	message := utils.CopyProperties[model.Message](req)
	message.IpAddress = ipAddress
	message.IpSource = ipSource
	// 根据 blogConfig 中的配置设置默认是否审核
	message.IsReview = blogInfoService.GetBlogConfig().IsMessageReview

	dao.Create(&message)
}

func (*Message) Delete(ids []int) (code int) {
	dao.Delete(model.Message{}, "id in ?", ids)
	return r.OK
}

// 修改审核状态
func (*Message) UpdateReview(req req.UpdateReview) (code int) {
	maps := map[string]any{"is_review": req.IsReview}
	dao.UpdatesMap(&model.Message{}, maps, "id IN ?", req.Ids)
	return r.OK
}

func (*Message) GetList(req req.GetMessages) resp.PageResult[[]model.Message] {
	list, total := messageDao.GetList(req)
	return resp.PageResult[[]model.Message]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    total,
		List:     list,
	}
}

// 前台
func (*Message) GetFrontList() []model.Message {
	return dao.List([]model.Message{}, "*", "", "is_review = 1")
}
