package model

import (
	"gorm.io/gorm"
)

// 通知类型
const (
	NOTIFY_COMMENT_REPLY   = iota + 1 // 评论被回复
	NOTIFY_ARTICLE_COMMENT            // 自己的文章收到评论
)

/*
站内通知

邮件功能默认关闭, 所以「别人回复了你」这类事件原本完全没有出口 —— 只有主动回访
才看得到。这张表让评论区能真正对话, 且不依赖 SMTP 或任何外部服务。

关联信息存 id 而不是整份拷贝: 文章改了标题、用户换了昵称之后, 通知里显示的
应该跟着变, 而不是留下一条对不上的旧快照。所以列表查询用 join 取昵称和标题。
只有评论内容截断存了一份, 因为评论可能被删, 但通知里那句摘要仍然要看得到。
*/
type Notification struct {
	Model

	UserId     int    `gorm:"index:idx_notify_user,priority:1;comment:接收者(user_auth_id)" json:"user_id"`
	FromUserId int    `gorm:"comment:触发者(user_auth_id)" json:"from_user_id"`
	Type       int    `gorm:"type:tinyint;comment:类型(1-评论被回复 2-文章被评论)" json:"type"`
	ArticleId  int    `gorm:"comment:关联文章, 用于前台跳转" json:"article_id"`
	CommentId  int    `gorm:"comment:触发的评论" json:"comment_id"`
	Content    string `gorm:"type:varchar(500);comment:触发时的评论摘要" json:"content"`
	IsRead     bool   `gorm:"index:idx_notify_user,priority:2;comment:是否已读" json:"is_read"`
}

// 列表项: 昵称、头像、文章标题都是查询时 join 出来的, 不落库
type NotificationVO struct {
	Notification

	FromNickname string `json:"from_nickname"`
	FromAvatar   string `json:"from_avatar"`
	ArticleTitle string `json:"article_title"`
}

// 通知列表只是摘要, 详情点进文章看, 所以长评论截断
const notifyContentLimit = 100

func AddNotification(db *gorm.DB, n *Notification) error {
	// 按 rune 截断, 按字节切会把中文切成半个字符
	if runes := []rune(n.Content); len(runes) > notifyContentLimit {
		n.Content = string(runes[:notifyContentLimit]) + "..."
	}
	return db.Create(n).Error
}

// 某个用户的通知列表, isRead 为 nil 表示不过滤
func GetNotificationList(db *gorm.DB, userId, page, size int, isRead *bool) (list []NotificationVO, total int64, err error) {
	list = make([]NotificationVO, 0)

	filter := func(d *gorm.DB) *gorm.DB {
		d = d.Where("n.user_id = ?", userId)
		if isRead != nil {
			d = d.Where("n.is_read = ?", *isRead)
		}
		return d
	}

	if err := db.Table("notification n").Scopes(filter).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 触发者昵称/头像在 user_info 上, 而通知里存的是 user_auth_id, 中间要过一层
	result := db.Table("notification n").
		Scopes(filter).
		Joins("LEFT JOIN user_auth ua ON ua.id = n.from_user_id").
		Joins("LEFT JOIN user_info ui ON ui.id = ua.user_info_id").
		Joins("LEFT JOIN article a ON a.id = n.article_id").
		Select("n.*, ui.nickname AS from_nickname, ui.avatar AS from_avatar, a.title AS article_title").
		Order("n.id DESC").
		Scopes(Paginate(page, size)).
		Find(&list)

	return list, total, result.Error
}

func GetUnreadNotificationCount(db *gorm.DB, userId int) (count int64, err error) {
	err = db.Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userId, false).
		Count(&count).Error
	return count, err
}

// 标记已读
//
// 必须带 user_id 条件: 只按 id 更新的话, 任何登录用户都能把别人的通知标成已读
func ReadNotifications(db *gorm.DB, userId int, ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := db.Model(&Notification{}).
		Where("user_id = ? AND id IN ?", userId, ids).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

func ReadAllNotifications(db *gorm.DB, userId int) (int64, error) {
	result := db.Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userId, false).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

func DeleteNotifications(db *gorm.DB, userId int, ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := db.Where("user_id = ? AND id IN ?", userId, ids).Delete(&Notification{})
	return result.RowsAffected, result.Error
}

/*
评论产生的通知

规则:
  - 回复别人 -> 通知被回复者(NOTIFY_COMMENT_REPLY)
  - 在文章下发顶级评论 -> 通知文章作者(NOTIFY_ARTICLE_COMMENT)
  - 自己回复自己、自己评论自己的文章都不通知
  - 只有文章下的评论才通知: 友链/说说页的评论没有可跳转的目标

待审核的评论先不通知, 否则等于把还没过审的内容推给了对方;
过审后补发通知需要在后台审核动作里再触发, 目前没做。
*/
func NotifyOnComment(db *gorm.DB, comment *Comment) error {
	if comment == nil || !comment.IsReview || comment.Type != TYPE_ARTICLE {
		return nil
	}

	receiver, typ := comment.ReplyUserId, NOTIFY_COMMENT_REPLY
	if receiver == 0 {
		// 顶级评论: 通知文章作者
		var authorId int
		if err := db.Model(&Article{}).
			Select("user_id").
			Where("id = ?", comment.TopicId).
			Scan(&authorId).Error; err != nil {
			return err
		}
		receiver, typ = authorId, NOTIFY_ARTICLE_COMMENT
	}

	if receiver == 0 || receiver == comment.UserId {
		return nil
	}

	return AddNotification(db, &Notification{
		UserId:     receiver,
		FromUserId: comment.UserId,
		Type:       typ,
		ArticleId:  comment.TopicId,
		CommentId:  comment.ID,
		Content:    comment.Content,
	})
}
