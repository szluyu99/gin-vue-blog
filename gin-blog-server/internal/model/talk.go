package model

import (
	"gorm.io/gorm"
)

/*
说说: 低门槛的短内容更新

和文章共用评论表(Comment.Type = TYPE_TALK), 所以这里不需要自己的评论字段。
第一版不支持图片: 图片一多就要面对存储和缩略图, 而现在上传只有本地目录和七牛云两条路,
收益不如先把「能发、能看、能评论」这条链跑通。
*/
type Talk struct {
	Model
	UserId  int    `gorm:"index:idx_talk_user;comment:发布者" json:"user_id"`
	Content string `gorm:"type:varchar(1000);not null;comment:说说内容" json:"content"`
	Status  int    `gorm:"type:tinyint;default:1;comment:状态(1公开 2私密)" json:"status"`
	IsTop   bool   `gorm:"comment:是否置顶" json:"is_top"`

	// Belongs To
	User *UserAuth `gorm:"foreignKey:UserId" json:"user"`
}

// 前台列表项: 评论数是聚合出来的, 不落库
type TalkVO struct {
	Talk

	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	CommentCount int64  `json:"comment_count"`
}

// 后台说说列表: status 为 0 表示不过滤
func GetTalkList(db *gorm.DB, page, size, status int) (list []Talk, total int64, err error) {
	list = make([]Talk, 0)

	filter := func(d *gorm.DB) *gorm.DB {
		if status != 0 {
			d = d.Where("status = ?", status)
		}
		return d
	}

	if err := db.Model(&Talk{}).Scopes(filter).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := db.Model(&Talk{}).Scopes(filter).
		Preload("User").Preload("User.UserInfo").
		Order("is_top DESC, id DESC").
		Scopes(Paginate(page, size)).
		Find(&list)

	return list, total, result.Error
}

// 后台说说详情: 私密的也要能查到(要能编辑)
func GetTalk(db *gorm.DB, id int) (talk *Talk, err error) {
	result := db.Model(&Talk{}).
		Preload("User").Preload("User.UserInfo").
		Where("id = ?", id).First(&talk)
	return talk, result.Error
}

func SaveOrUpdateTalk(db *gorm.DB, id, userId int, content string, status int, isTop bool) (*Talk, error) {
	talk := Talk{
		Model:   Model{ID: id},
		UserId:  userId,
		Content: content,
		Status:  status,
		IsTop:   isTop,
	}

	var result *gorm.DB
	if id > 0 {
		// 用 Select 显式列出字段: IsTop 取消置顶时是零值,
		// Updates 默认跳过零值, 会出现"取消不掉置顶"
		result = db.Model(&talk).Select("content", "status", "is_top").Updates(talk)
	} else {
		result = db.Create(&talk)
	}

	return &talk, result.Error
}

func DeleteTalks(db *gorm.DB, ids []int) (rows int64, err error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// 说说下面的评论一起删, 否则会变成孤儿数据: 评论还在, 主题已经没了
	err = db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id IN ?", ids).Delete(&Talk{})
		if result.Error != nil {
			return result.Error
		}
		rows = result.RowsAffected

		return tx.Where("type = ? AND topic_id IN ?", TYPE_TALK, ids).Delete(&Comment{}).Error
	})

	return rows, err
}

// 前台说说列表: 只有公开的, 置顶排前面, 带上昵称头像和评论数
func GetBlogTalkList(db *gorm.DB, page, size int) (list []TalkVO, total int64, err error) {
	list = make([]TalkVO, 0)

	if err := db.Model(&Talk{}).Where("status = ?", STATUS_PUBLIC).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var talks []Talk
	if err := db.Model(&Talk{}).Where("status = ?", STATUS_PUBLIC).
		Preload("User").Preload("User.UserInfo").
		Order("is_top DESC, id DESC").
		Scopes(Paginate(page, size)).
		Find(&talks).Error; err != nil {
		return nil, 0, err
	}

	if len(talks) == 0 {
		return list, total, nil
	}

	counts, err := talkCommentCounts(db, talks)
	if err != nil {
		return nil, 0, err
	}

	for _, talk := range talks {
		list = append(list, newTalkVO(talk, counts[talk.ID]))
	}

	return list, total, nil
}

// 前台说说详情: 私密的当作不存在
func GetBlogTalk(db *gorm.DB, id int) (*TalkVO, error) {
	var talk Talk
	if err := db.Model(&Talk{}).
		Preload("User").Preload("User.UserInfo").
		Where("id = ? AND status = ?", id, STATUS_PUBLIC).
		First(&talk).Error; err != nil {
		return nil, err
	}

	counts, err := talkCommentCounts(db, []Talk{talk})
	if err != nil {
		return nil, err
	}

	vo := newTalkVO(talk, counts[talk.ID])
	return &vo, nil
}

// 一条 SQL 把本页所有说说的评论数查回来, 而不是每条说说查一次(N+1)
func talkCommentCounts(db *gorm.DB, talks []Talk) (map[int]int64, error) {
	ids := make([]int, 0, len(talks))
	for _, talk := range talks {
		ids = append(ids, talk.ID)
	}

	type row struct {
		TopicId int
		Count   int64
	}
	var rows []row

	err := db.Model(&Comment{}).
		Select("topic_id", "count(1) as count").
		Where("type = ? AND is_review = ? AND topic_id IN ?", TYPE_TALK, true, ids).
		Group("topic_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[int]int64, len(rows))
	for _, r := range rows {
		counts[r.TopicId] = r.Count
	}
	return counts, nil
}

// 昵称头像摊平到顶层: 前台只用这两个字段, 不必让页面自己层层取 user.info
func newTalkVO(talk Talk, commentCount int64) TalkVO {
	vo := TalkVO{Talk: talk, CommentCount: commentCount}
	if talk.User != nil && talk.User.UserInfo != nil {
		vo.Nickname = talk.User.UserInfo.Nickname
		vo.Avatar = talk.User.UserInfo.Avatar
	}
	return vo
}
