package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Notification struct{}

type NotificationQuery struct {
	PageQuery
	// 指针类型: 不传表示全部, 传 false 只看未读
	IsRead *bool `form:"is_read"`
}

type ReadNotificationReq struct {
	Ids []int `json:"ids"`
}

// @Summary 当前用户的站内通知列表
// @Description 分页返回, 带触发者昵称头像与文章标题
// @Tags Notification
// @Produce json
// @Param is_read query bool false "是否已读, 不传为全部"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 0 {object} Response[PageResult[model.NotificationVO]]
// @Security ApiKeyAuth
// @Router /front/notification/list [get]
func (*Notification) GetList(c *gin.Context) {
	var query NotificationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	list, total, err := model.GetNotificationList(GetDB(c), auth.ID, query.Page, query.Size, query.IsRead)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.NotificationVO]{
		Total: total,
		List:  list,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary 当前用户的未读通知数
// @Description 前台头部铃铛上的红点
// @Tags Notification
// @Produce json
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /front/notification/unread [get]
func (*Notification) GetUnreadCount(c *gin.Context) {
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	count, err := model.GetUnreadNotificationCount(GetDB(c), auth.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, count)
}

// @Summary 标记通知为已读
// @Description ids 为空时标记当前用户的全部未读
// @Tags Notification
// @Accept json
// @Produce json
// @Param data body ReadNotificationReq true "通知 ID 列表"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /front/notification/read [put]
func (*Notification) Read(c *gin.Context) {
	var req ReadNotificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	db := GetDB(c)

	// 两个查询都带 user_id 条件, 不能只按 id 更新: 否则任何登录用户
	// 都能把别人的通知标成已读
	var (
		rows int64
		err  error
	)
	if len(req.Ids) == 0 {
		rows, err = model.ReadAllNotifications(db, auth.ID)
	} else {
		rows, err = model.ReadNotifications(db, auth.ID, req.Ids)
	}
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// @Summary 删除通知
// @Description 只能删自己的
// @Tags Notification
// @Accept json
// @Produce json
// @Param data body ReadNotificationReq true "通知 ID 列表"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /front/notification [delete]
func (*Notification) Delete(c *gin.Context) {
	var req ReadNotificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	rows, err := model.DeleteNotifications(GetDB(c), auth.ID, req.Ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}
