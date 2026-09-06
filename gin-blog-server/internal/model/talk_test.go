package model

import (
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTalkUser(t *testing.T, db *gorm.DB, username, nickname string) *UserAuth {
	user := UserAuth{Username: username, Password: "x", UserInfo: &UserInfo{Nickname: nickname, Avatar: "a.png"}}
	require.Nil(t, db.Create(&user).Error)
	return &user
}

// 后台列表: 置顶排前面, status 可筛选, 分页生效
func TestGetTalkList(t *testing.T) {
	db := newModelDB(t)
	user := newTalkUser(t, db, "u1", "博主")

	_, err := SaveOrUpdateTalk(db, 0, user.ID, "第一条", STATUS_PUBLIC, false)
	assert.Nil(t, err)
	_, err = SaveOrUpdateTalk(db, 0, user.ID, "私密的", STATUS_SECRET, false)
	assert.Nil(t, err)
	top, err := SaveOrUpdateTalk(db, 0, user.ID, "置顶的", STATUS_PUBLIC, true)
	assert.Nil(t, err)

	list, total, err := GetTalkList(db, 1, 10, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, top.ID, list[0].ID, "置顶的排最前")
	assert.Equal(t, "博主", list[0].User.UserInfo.Nickname, "要带上发布者")

	// 只看私密的
	list, total, err = GetTalkList(db, 1, 10, STATUS_SECRET)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)
	assert.Equal(t, "私密的", list[0].Content)

	// 分页
	list, _, err = GetTalkList(db, 2, 2, 0)
	assert.Nil(t, err)
	assert.Len(t, list, 1)
}

// 编辑要能把置顶取消掉: Updates 默认跳过零值, 这里是显式 Select 的
func TestSaveOrUpdateTalkCancelTop(t *testing.T) {
	db := newModelDB(t)
	user := newTalkUser(t, db, "u1", "博主")

	talk, err := SaveOrUpdateTalk(db, 0, user.ID, "内容", STATUS_PUBLIC, true)
	assert.Nil(t, err)

	_, err = SaveOrUpdateTalk(db, talk.ID, user.ID, "改过的内容", STATUS_SECRET, false)
	assert.Nil(t, err)

	got, err := GetTalk(db, talk.ID)
	assert.Nil(t, err)
	assert.Equal(t, "改过的内容", got.Content)
	assert.Equal(t, STATUS_SECRET, got.Status)
	assert.False(t, got.IsTop, "取消置顶要能写进去")
}

// 删说说要连它的评论一起删, 否则评论变成孤儿数据
func TestDeleteTalksRemovesComments(t *testing.T) {
	db := newModelDB(t)
	user := newTalkUser(t, db, "u1", "博主")

	talk, err := SaveOrUpdateTalk(db, 0, user.ID, "内容", STATUS_PUBLIC, false)
	assert.Nil(t, err)
	other, err := SaveOrUpdateTalk(db, 0, user.ID, "留下的", STATUS_PUBLIC, false)
	assert.Nil(t, err)

	_, err = AddComment(db, user.ID, TYPE_TALK, talk.ID, "说说评论", true)
	assert.Nil(t, err)
	keep, err := AddComment(db, user.ID, TYPE_TALK, other.ID, "别人的评论", true)
	assert.Nil(t, err)
	// 同 id 的文章评论不能被牵连
	article, err := AddComment(db, user.ID, TYPE_ARTICLE, talk.ID, "文章评论", true)
	assert.Nil(t, err)

	rows, err := DeleteTalks(db, []int{talk.ID})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)

	var count int64
	assert.Nil(t, db.Model(&Comment{}).Where("id = ?", keep.ID).Count(&count).Error)
	assert.Equal(t, int64(1), count, "其他说说的评论不动")
	assert.Nil(t, db.Model(&Comment{}).Where("id = ?", article.ID).Count(&count).Error)
	assert.Equal(t, int64(1), count, "topic_id 相同的文章评论不能被牵连")
	assert.Nil(t, db.Model(&Comment{}).Where("type = ? AND topic_id = ?", TYPE_TALK, talk.ID).Count(&count).Error)
	assert.Zero(t, count, "被删说说的评论要清掉")
}

// 前台: 只给公开的, 带评论数(只算已审核的)
func TestGetBlogTalkList(t *testing.T) {
	db := newModelDB(t)
	user := newTalkUser(t, db, "u1", "博主")

	pub, err := SaveOrUpdateTalk(db, 0, user.ID, "公开的", STATUS_PUBLIC, false)
	assert.Nil(t, err)
	secret, err := SaveOrUpdateTalk(db, 0, user.ID, "私密的", STATUS_SECRET, false)
	assert.Nil(t, err)

	_, err = AddComment(db, user.ID, TYPE_TALK, pub.ID, "已审核", true)
	assert.Nil(t, err)
	_, err = AddComment(db, user.ID, TYPE_TALK, pub.ID, "待审核", false)
	assert.Nil(t, err)

	list, total, err := GetBlogTalkList(db, 1, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total, "私密的不算进总数")
	require.Len(t, list, 1)
	assert.Equal(t, "公开的", list[0].Content)
	assert.Equal(t, int64(1), list[0].CommentCount, "只算已审核的评论")
	assert.Equal(t, "博主", list[0].Nickname, "昵称头像摊平到顶层")
	assert.Equal(t, "a.png", list[0].Avatar)

	// 详情: 私密的当作不存在
	vo, err := GetBlogTalk(db, pub.ID)
	assert.Nil(t, err)
	assert.Equal(t, "公开的", vo.Content)
	assert.Equal(t, int64(1), vo.CommentCount)

	_, err = GetBlogTalk(db, secret.ID)
	assert.NotNil(t, err, "私密说说前台查不到")
}
