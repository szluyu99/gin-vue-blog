package dto

type CommentCountDTO struct {
	ID           int `json:"id"`
	CommentCount int `json:"comment_count"`
}

type ReplyCountDTO struct {
	CommentId  int `json:"comment_id"`
	ReplyCount int `json:"reply_count"`
}
