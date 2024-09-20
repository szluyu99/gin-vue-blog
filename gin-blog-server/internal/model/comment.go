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
	UserId      int    `json:"user_id"`       // 评论者
	ReplyUserId int    `json:"reply_user_id"` // 被回复者
	TopicId     int    `json:"topic_id"`      // 评论的文章
	ParentId    int    `json:"parent_id"`     // 父评论 被回复的评论
	Content     string `gorm:"type:varchar(500);not null" json:"content"`
	Type        int    `gorm:"type:tinyint(1);not null;comment:评论类型(1.文章 2.友链 3.说说)" json:"type"` // 评论类型 1.文章 2.友链 3.说说
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

// 获取后台评论列表
func GetCommentList(db *gorm.DB, page, size, typ int, isReview *bool, nickname string) (data []Comment, total int64, err error) {
	
	// SELECT UID FROM user_info WHERE nikename LIKE nickname
	var uid int
	if nickname != "" {
		result := db.Model(&UserInfo{}).Where("nickname LIKE ?",nickname).Pluck("id",&uid)
		if result.Error != nil{
			return nil,0,result.Error
		}
		db = db.Where("user_id = ?",uid)
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

// 获取博客评论列表
func GetCommentVOList(db *gorm.DB, page, size, topic, typ int) (data []CommentVO, total int64, err error) {
	var list []Comment

	tx := db.Model(&Comment{})
	if typ != 0 {
		tx = tx.Where("type = ?", typ)
	}
	if topic != 0 {
		tx = tx.Where("topic_id = ?", topic)
	}

	// 获取顶级评论列表
	tx.Where("parent_id = 0").
		Count(&total).
		Preload("User").Preload("User.UserInfo").
		// Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Order("id DESC").
		Scopes(Paginate(page, size))
	if err := tx.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	// 获取顶级评论的回复列表
	for _, v := range list {
		replyList := make([]CommentVO, 0)

		tx := db.Model(&Comment{})
		tx.Where("parent_id = ?", v.ID).
			Preload("User").Preload("User.UserInfo").
			// Preload("ReplyUser").Preload("ReplyUser.UserInfo")
			Order("id DESC")
		if err := tx.Find(&replyList).Error; err != nil {
			return nil, 0, err
		}

		data = append(data, CommentVO{
			ReplyCount: len(replyList),
			Comment:    v,
			ReplyList:  replyList,
		})
	}

	return data, total, nil
}

// 根据 [评论id] 获取 [回复列表]
func GetCommentReplyList(db *gorm.DB, id, page, size int) (data []Comment, err error) {
	result := db.Model(&Comment{}).
		Where(&Comment{ParentId: id}).
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
