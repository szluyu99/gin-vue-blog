package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func seedNotifyUser(t *testing.T, db *gorm.DB, nickname string) *UserAuth {
	t.Helper()
	auth := UserAuth{Username: nickname, Password: "x", UserInfo: &UserInfo{Nickname: nickname, Avatar: "a.png"}}
	assert.Nil(t, db.Create(&auth).Error)
	return &auth
}

// 列表要把触发者昵称头像和文章标题 join 出来, 而不是在通知里存快照
func TestNotificationList(t *testing.T) {
	db := newModelDB(t)
	me := seedNotifyUser(t, db, "me")
	other := seedNotifyUser(t, db, "other")
	article := seedArticle(t, db, Article{Title: "被评论的文章", Status: STATUS_PUBLIC, UserId: me.ID})

	assert.Nil(t, AddNotification(db, &Notification{
		UserId: me.ID, FromUserId: other.ID, Type: NOTIFY_ARTICLE_COMMENT,
		ArticleId: article.ID, CommentId: 1, Content: "写得不错",
	}))
	// 别人的通知不能出现在我的列表里
	assert.Nil(t, AddNotification(db, &Notification{
		UserId: other.ID, FromUserId: me.ID, Type: NOTIFY_COMMENT_REPLY,
		ArticleId: article.ID, CommentId: 2, Content: "多谢",
	}))

	list, total, err := GetNotificationList(db, me.ID, 1, 10, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)
	assert.Equal(t, "other", list[0].FromNickname)
	assert.Equal(t, "a.png", list[0].FromAvatar)
	assert.Equal(t, "被评论的文章", list[0].ArticleTitle)
	assert.False(t, list[0].IsRead)
}

func TestNotificationUnreadAndRead(t *testing.T) {
	db := newModelDB(t)
	me := seedNotifyUser(t, db, "me")
	other := seedNotifyUser(t, db, "other")

	var ids []int
	for i := 0; i < 3; i++ {
		n := Notification{UserId: me.ID, FromUserId: other.ID, Type: NOTIFY_COMMENT_REPLY, Content: "x"}
		assert.Nil(t, AddNotification(db, &n))
		ids = append(ids, n.ID)
	}
	// 别人的一条, 用来验证越权
	his := Notification{UserId: other.ID, FromUserId: me.ID, Type: NOTIFY_COMMENT_REPLY, Content: "y"}
	assert.Nil(t, AddNotification(db, &his))

	count, err := GetUnreadNotificationCount(db, me.ID)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), count)

	rows, err := ReadNotifications(db, me.ID, ids[:1])
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)

	count, _ = GetUnreadNotificationCount(db, me.ID)
	assert.Equal(t, int64(2), count)

	// 只看未读
	unread := false
	list, total, err := GetNotificationList(db, me.ID, 1, 10, &unread)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)

	// 越权: 带上别人的 id 也改不动
	rows, err = ReadNotifications(db, me.ID, []int{his.ID})
	assert.Nil(t, err)
	assert.Zero(t, rows)
	count, _ = GetUnreadNotificationCount(db, other.ID)
	assert.Equal(t, int64(1), count)

	rows, err = ReadAllNotifications(db, me.ID)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), rows)
	count, _ = GetUnreadNotificationCount(db, me.ID)
	assert.Zero(t, count)
}

func TestNotificationDelete(t *testing.T) {
	db := newModelDB(t)
	me := seedNotifyUser(t, db, "me")
	other := seedNotifyUser(t, db, "other")

	mine := Notification{UserId: me.ID, FromUserId: other.ID, Content: "x"}
	assert.Nil(t, AddNotification(db, &mine))
	his := Notification{UserId: other.ID, FromUserId: me.ID, Content: "y"}
	assert.Nil(t, AddNotification(db, &his))

	// 空 id 列表不该删掉任何东西
	rows, err := DeleteNotifications(db, me.ID, nil)
	assert.Nil(t, err)
	assert.Zero(t, rows)

	rows, err = DeleteNotifications(db, me.ID, []int{mine.ID, his.ID})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows, "只能删自己的")

	var left int64
	db.Model(&Notification{}).Count(&left)
	assert.Equal(t, int64(1), left)
}

// 长内容按 rune 截断: 按字节切会把中文切成半个字符
func TestNotificationContentTruncate(t *testing.T) {
	db := newModelDB(t)
	me := seedNotifyUser(t, db, "me")

	long := ""
	for i := 0; i < 150; i++ {
		long += "中"
	}
	n := Notification{UserId: me.ID, Content: long}
	assert.Nil(t, AddNotification(db, &n))

	assert.Equal(t, notifyContentLimit+3, len([]rune(n.Content)))
	assert.True(t, []rune(n.Content)[notifyContentLimit] == '.')
}

// 通知的产生规则
func TestNotifyOnComment(t *testing.T) {
	db := newModelDB(t)
	author := seedNotifyUser(t, db, "author")
	guest := seedNotifyUser(t, db, "guest")
	article := seedArticle(t, db, Article{Title: "文章", Status: STATUS_PUBLIC, UserId: author.ID})

	countFor := func(userId int) int64 {
		c, err := GetUnreadNotificationCount(db, userId)
		assert.Nil(t, err)
		return c
	}

	// 1. 游客在文章下发顶级评论 -> 通知作者
	top, err := AddComment(db, guest.ID, TYPE_ARTICLE, article.ID, "顶级评论", true)
	assert.Nil(t, err)
	assert.Nil(t, NotifyOnComment(db, top))
	assert.Equal(t, int64(1), countFor(author.ID))

	// 2. 作者回复游客 -> 通知游客
	reply, err := ReplyComment(db, author.ID, guest.ID, top.ID, "多谢", true)
	assert.Nil(t, err)
	assert.Nil(t, NotifyOnComment(db, reply))
	assert.Equal(t, int64(1), countFor(guest.ID))

	// 3. 作者评论自己的文章 -> 不通知
	own, err := AddComment(db, author.ID, TYPE_ARTICLE, article.ID, "自己补充", true)
	assert.Nil(t, err)
	assert.Nil(t, NotifyOnComment(db, own))
	assert.Equal(t, int64(1), countFor(author.ID), "自己评论自己的文章不该产生通知")

	// 4. 自己回复自己 -> 不通知
	selfReply, err := ReplyComment(db, guest.ID, guest.ID, top.ID, "补充一句", true)
	assert.Nil(t, err)
	assert.Nil(t, NotifyOnComment(db, selfReply))
	assert.Equal(t, int64(1), countFor(guest.ID), "自己回复自己不该产生通知")

	// 5. 待审核的评论先不通知
	pending, err := AddComment(db, guest.ID, TYPE_ARTICLE, article.ID, "待审核", false)
	assert.Nil(t, err)
	assert.Nil(t, NotifyOnComment(db, pending))
	assert.Equal(t, int64(1), countFor(author.ID))

	// 6. 友链页的评论没有跳转目标, 不通知
	linkComment, err := AddComment(db, guest.ID, TYPE_LINK, 0, "申请友链", true)
	assert.Nil(t, err)
	assert.Nil(t, NotifyOnComment(db, linkComment))
	assert.Equal(t, int64(1), countFor(author.ID))

	// nil 不能 panic
	assert.Nil(t, NotifyOnComment(db, nil))
}
