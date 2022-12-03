package req

type GetComments struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *int8  `form:"is_review"`
	Type     int    `form:"type"`
}

// 修改审核
// type UpdateCommentReview struct {
// 	Ids      []int `json:"ids"`
// 	IsReview int   `json:"is_review"`
// }

type SaveComment struct {
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int8   `json:"type" form:"type" validate:"required,min=1,max=3" label:"评论类型"`
}

type GetFrontComments struct {
	PageQuery
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int    `json:"type" form:"type"`
}
