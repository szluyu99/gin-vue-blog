package model

import (
	"gorm.io/gorm"
)

const (
	TYPE_ARTICLE = iota + 1 // 文章
	TYPE_LINK               // 友链
	TYPE_TALK               // 说说
)

/*
如果评论类型是文章，那么 topic_id 就是文章的 id
如果评论类型是友链，不需要 topic_id
*/

type Comment struct {
	Model
	UserId      int    `gorm:"index:idx_comment_user" json:"user_id"`              // 评论者
	ReplyUserId int    `json:"reply_user_id"`                                      // 被回复者
	TopicId     int    `gorm:"index:idx_comment_topic,priority:2" json:"topic_id"` // 评论的文章
	ParentId    int    `gorm:"index:idx_comment_parent" json:"parent_id"`          // 父评论 被回复的评论
	Content     string `gorm:"type:varchar(500);not null" json:"content"`
	Type        int    `gorm:"type:tinyint(1);not null;index:idx_comment_topic,priority:1;comment:评论类型(1.文章 2.友链 3.说说)" json:"type"` // 评论类型 1.文章 2.友链 3.说说
	IsReview    bool   `json:"is_review"`

	// Belongs To
	User      *UserAuth `gorm:"foreignKey:UserId" json:"user"`
	ReplyUser *UserAuth `gorm:"foreignKey:ReplyUserId" json:"reply_user"`
	Article   *Article  `gorm:"foreignKey:TopicId" json:"article"`
}

type CommentVO struct {
	Comment
	LikeCount  int         `json:"like_count" gorm:"-"`
	ReplyCount int         `json:"reply_count" gorm:"-"`
	ReplyList  []CommentVO `json:"reply_list" gorm:"-"`
}

// 新增评论
func AddComment(db *gorm.DB, userId, typ, topicId int, content string, isReview bool) (*Comment, error) {
	comment := Comment{
		UserId:   userId,
		TopicId:  topicId,
		Content:  content,
		Type:     typ,
		IsReview: isReview,
	}
	result := db.Create(&comment)
	return &comment, result.Error
}

// 回复评论
func ReplyComment(db *gorm.DB, userId, replyUserId, parentId int, content string, isReview bool) (*Comment, error) {
	var parent Comment
	result := db.First(&parent, parentId)
	if result.Error != nil {
		return nil, result.Error
	}

	comment := Comment{
		UserId:      userId,
		Content:     content,
		ReplyUserId: replyUserId,
		ParentId:    parentId,
		IsReview:    isReview,
		TopicId:     parent.TopicId, // 主题和父评论一样
		Type:        parent.Type,    // 类型和父评论一样
	}
	result = db.Create(&comment)
	return &comment, result.Error
}

/*
批量修改审核状态

返回本次「由待审核变为已审核」的评论, 交给调用方补发站内通知 ——
发评论时如果配置要求审核, 那条评论是不可见的, 当时不发通知(否则等于把没过审的
内容推给对方), 所以过审这一刻才是通知该发出去的时机。

只挑真正发生 false -> true 的: 对已经是已审核的评论再点一次"通过"不该重复通知。
通知本身不在这里发, model 层不该关心日志和降级策略。
*/
func ReviewComments(db *gorm.DB, ids []int, isReview bool) (rows int64, approved []Comment, err error) {
	approved = make([]Comment, 0)
	if len(ids) == 0 {
		return 0, approved, nil
	}

	if isReview {
		if err := db.Where("id IN ? AND is_review = ?", ids, false).Find(&approved).Error; err != nil {
			return 0, nil, err
		}
	}

	result := db.Model(Comment{}).Where("id IN ?", ids).Update("is_review", isReview)
	if result.Error != nil {
		return 0, nil, result.Error
	}

	// 查出来时还是待审核状态, 这里补上新状态, NotifyOnComment 才会认
	for i := range approved {
		approved[i].IsReview = true
	}

	return result.RowsAffected, approved, nil
}

// 获取后台评论列表
func GetCommentList(db *gorm.DB, page, size, typ int, isReview *bool, nickname string) (data []Comment, total int64, err error) {

	// SELECT UID FROM user_info WHERE nikename LIKE nickname
	var uid int
	if nickname != "" {
		result := db.Model(&UserInfo{}).Where("nickname LIKE ?", nickname).Pluck("id", &uid)
		if result.Error != nil {
			return nil, 0, result.Error
		}
		db = db.Where("user_id = ?", uid)
	}

	if typ != 0 {
		db = db.Where("type = ?", typ)
	}
	if isReview != nil {
		db = db.Where("is_review = ?", *isReview)
	}

	result := db.Model(&Comment{}).
		Count(&total).
		Preload("User").Preload("User.UserInfo").
		Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Preload("Article").
		Order("id DESC").
		Scopes(Paginate(page, size)).
		Find(&data)

	return data, total, result.Error
}

// 获取博客评论列表: 顶级评论分页, 每条评论带上它的回复
// 只返回审核通过的, 未审核的不能出现在前台(GetArticleCommentCount 也是按 is_review 统计的)
func GetCommentVOList(db *gorm.DB, page, size, topic, typ int) (data []CommentVO, total int64, err error) {
	// 过滤条件写成 scope, Count 和 Find 各用一次全新的查询,
	// 之前把 Preload / Scopes(Paginate) 的返回值丢掉了(链式调用返回的是新对象),
	// 导致分页参数完全没生效, 每次都把所有顶级评论返回给前台
	filter := func(d *gorm.DB) *gorm.DB {
		d = d.Where("parent_id = 0 AND is_review = ?", true)
		if typ != 0 {
			d = d.Where("type = ?", typ)
		}
		if topic != 0 {
			d = d.Where("topic_id = ?", topic)
		}
		return d
	}

	if err := db.Model(&Comment{}).Scopes(filter).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []Comment
	if err := db.Model(&Comment{}).Scopes(filter).
		Preload("User").Preload("User.UserInfo").
		Order("id DESC").
		Scopes(Paginate(page, size)).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	data = make([]CommentVO, 0, len(list))
	if len(list) == 0 {
		return data, total, nil
	}

	// 一次把本页所有顶级评论的回复查回来, 而不是每条评论查一次(N+1)
	ids := make([]int, 0, len(list))
	for _, v := range list {
		ids = append(ids, v.ID)
	}

	var replies []Comment
	if err := db.Model(&Comment{}).
		Where("parent_id IN ? AND is_review = ?", ids, true).
		Preload("User").Preload("User.UserInfo").
		Order("id DESC").
		Find(&replies).Error; err != nil {
		return nil, 0, err
	}

	grouped := make(map[int][]CommentVO, len(ids))
	for _, reply := range replies {
		grouped[reply.ParentId] = append(grouped[reply.ParentId], CommentVO{Comment: reply})
	}

	for _, v := range list {
		replyList := grouped[v.ID]
		if replyList == nil {
			replyList = make([]CommentVO, 0)
		}
		data = append(data, CommentVO{
			Comment:    v,
			ReplyCount: len(replyList),
			ReplyList:  replyList,
		})
	}

	return data, total, nil
}

// 根据 [评论id] 获取 [回复列表], 同样只返回审核通过的
func GetCommentReplyList(db *gorm.DB, id, page, size int) (data []Comment, err error) {
	result := db.Model(&Comment{}).
		Where("parent_id = ? AND is_review = ?", id, true).
		Preload("User").Preload("User.UserInfo").
		Order("id DESC").
		Scopes(Paginate(page, size)).
		Find(&data)
	return data, result.Error
}

// 获取某篇文章的评论数
func GetArticleCommentCount(db *gorm.DB, articleId int) (count int64, err error) {
	result := db.Model(&Comment{}).
		Where("topic_id = ? AND type = 1 AND is_review = 1", articleId).
		Count(&count)
	return count, result.Error
}

// 删除评论, 同时删掉挂在它们下面的回复, 返回实际删掉的所有评论 id
//
// 回复不一起删的话会变成孤儿数据: parent_id 指向已经不存在的评论,
// 但 GetCommentReplyList 仍然能把它们查出来。
// 返回的 id 交给调用方清理 Redis 里的点赞计数, 否则新评论复用 id 会继承旧的计数。
func DeleteComments(db *gorm.DB, ids []int) (rows int64, deletedIds []int, err error) {
	if len(ids) == 0 {
		return 0, []int{}, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		var replyIds []int
		if err := tx.Model(&Comment{}).
			Where("parent_id IN ?", ids).
			Pluck("id", &replyIds).Error; err != nil {
			return err
		}

		all := make([]int, 0, len(ids)+len(replyIds))
		all = append(all, ids...)
		all = append(all, replyIds...)

		result := tx.Where("id IN ?", all).Delete(&Comment{})
		if result.Error != nil {
			return result.Error
		}
		rows = result.RowsAffected
		deletedIds = all
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	return rows, deletedIds, nil
}
