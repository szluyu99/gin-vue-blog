package model

type Comment struct {
	Universal
	UserId      int    `gorm:"comment:评论用户id" json:"user_id"`
	ReplyUserId int    `gorm:"comment:回复用户id" json:"reply_user_id"`
	TopicId     int    `gorm:"comment:评论主题id" json:"topic_id"`
	ParentId    int    `gorm:"comment:父评论id" json:"parent_id"`
	Content     string `gorm:"type:varchar(500);not null;comment:评论内容" json:"content"`
	Type        int8   `gorm:"type:tinyint(1);not null;comment:评论类型(1.文章 2.友链 3.说说)" json:"type"`
	IsDelete    *int8  `gorm:"type:tinyint(1);not null;default:0;comment:是否删除(0.否 1.是)" json:"is_delete"`
	IsReview    *int8  `gorm:"type:tinyint(1);not null;default:0;comment:是否审核(0.否 1.是)" json:"is_review"`
}
