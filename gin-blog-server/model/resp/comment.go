package resp

import "time"

// 后台评论 VO
type CommentVo struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Avatar        string    `json:"avatar"`
	Nickname      string    `json:"nickname"`
	ReplyNickname string    `json:"reply_nickname"`
	ArticleTitle  string    `json:"article_title"`
	Content       string    `json:"content"`
	Type          int       `json:"type"`
	IsReview      int       `json:"is_review"`
}

// 前台评论 VO
type FrontCommentVO struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserId      int       `json:"user_id"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	Website     string    `json:"website"`
	Content     string    `json:"content"`
	LikeCount   int       `json:"like_count"`
	ReplyCount  int       `json:"reply_count"`
	ReplyVOList []ReplyVO `json:"reply_vo_list" gorm:"-"`
}

// 评论回复 VO
type ReplyVO struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	ParentId      int       `json:"parent_id"`
	UserId        int       `json:"user_id"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	Website       string    `json:"website"`
	ReplyUserId   int       `json:"reply_user_id"`
	ReplyNickname string    `json:"reply_nickname"`
	ReplyWebsite  string    `json:"reply_website"`
	Content       string    `json:"content"`
	LikeCount     int       `json:"like_count"`
}
