package dao

import (
	"gin-blog/model"
	"gin-blog/model/dto"
	"gin-blog/model/req"
	"gin-blog/model/resp"
)

type Comment struct{}

// TODO:
func (*Comment) Save(c model.Comment) {
}

func (*Comment) GetCount(req req.GetComments) int64 {
	var count int64
	tx := DB.Select("COUNT(1)").Table("comment c").
		Joins("LEFT JOIN user_info ui ON c.user_id = ui.id")
	if req.Type != 0 {
		tx = tx.Where("c.type = ?", req.Type)
	}
	if req.IsReview != nil {
		tx = tx.Where("c.is_review = ?", req.IsReview)
	}
	if req.Nickname != "" {
		tx = tx.Where("ui.nickname LIKE ?", "%"+req.Nickname+"%")
	}
	tx.Count(&count)
	return count
}

func (*Comment) GetList(req req.GetComments) []resp.CommentVo {
	list := make([]resp.CommentVo, 0)
	tx := DB.
		Select("c.id, u1.avatar, u1.nickname, u2.nickname reply_nickname, a.title article_title, c.content, c.type, c.created_at, c.is_review").
		Table("comment c").
		Joins("LEFT JOIN article a ON c.topic_id = a.id").
		Joins("LEFT JOIN user_info u1 ON c.user_id = u1.id").
		Joins("LEFT JOIN user_info u2 ON c.reply_user_id = u2.id")
	if req.Type != 0 {
		tx = tx.Where("c.type = ?", req.Type)
	}
	if req.IsReview != nil {
		tx = tx.Where("c.is_review = ?", req.IsReview)
	}
	if req.Nickname != "" {
		tx = tx.Where("u1.nickname LIKE ?", "%"+req.Nickname+"%")
	}
	tx.Order("id DESC").Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).Find(&list)
	return list
}

// 获取前台评论列表
func (*Comment) GetFrontList(req req.GetFrontComments) ([]resp.FrontCommentVO, int64) {
	list := make([]resp.FrontCommentVO, 0)
	var total int64

	tx := DB.
		Select("u.nickname, u.avatar, u.website, c.id, c.user_id, c.content, c.created_at").
		Table("comment c").
		Joins("LEFT JOIN user_info u ON c.user_id = u.id").
		Where("type = ? AND c.is_review = 1 AND parent_id = 0", req.Type)
	if req.TopicId != 0 {
		tx = tx.Where("topic_id", req.TopicId)
	}
	tx.Count(&total).
		Order("c.id DESC").
		Limit(req.PageSize).Offset((req.PageNum - 1) * req.PageSize).
		Find(&list)
	return list, total
}

// 根据 [评论id列表] 获取 回复列表
func (*Comment) GetReplyList(ids []int) []resp.ReplyVO {
	list := make([]resp.ReplyVO, 0)
	DB.
		Select("c.user_id, u1.nickname, u1.avatar, u1.website, c.reply_user_id, u2.nickname AS reply_nickname, "+
			"u2.website AS reply_website, c.id, c.parent_id, c.content, c.created_at").
		Table("comment c").
		Joins("JOIN user_info u1 ON c.user_id = u1.id").
		Joins("JOIN user_info u2 ON c.reply_user_id = u2.id").
		Where("c.is_review = 1 AND parent_id IN ?", ids).
		Order("c.parent_id, c.created_at ASC").
		Find(&list)
	// TODO! 理解一下
	// t := DB.
	// 	Select("if(@mno=b.parent_id,@rank:=@rank+1,@rank:=1) as row_number,@mno:=b.parent_id,b.*").
	// 	Table("(?) b", b)
	// DB.
	// 	Table("(?) t", t).
	// 	Where("row_numbser < 4").
	// 	Find(&list)
	return list
}

// 根据 parent_ids 获取已经审核的回复数量列表
func (*Comment) GetReplyCountListByCommentId(ids []int) []dto.ReplyCountDTO {
	list := make([]dto.ReplyCountDTO, 0)
	DB.
		Select("parent_id AS comment_id, COUNT(1) AS reply_count").
		Table("comment").
		Where("is_review = 1 AND parent_id in ?", ids).
		Group("parent_id").
		Find(&list)
	return list
}

// 根据 [评论id] 获取 回复列表
func (*Comment) GetReplyListByCommentId(id int, req req.PageQuery) []resp.ReplyVO {
	list := make([]resp.ReplyVO, 0)
	DB.
		Select("c.user_id, u1.nickname, u1.avatar, u1.website, c.reply_user_id, u2.nickname AS reply_nickname, "+
			"u2.website AS reply_website, c.id, c.parent_id, c.content, c.created_at").
		Table("comment c").
		Joins("JOIN user_info u1 ON c.user_id = u1.id").
		Joins("JOIN user_info u2 ON c.reply_user_id = u2.id").
		Where("c.is_review = 1 AND parent_id = ?", id).
		Order("c.id ASC").
		Limit(req.PageSize).Offset(req.PageSize * (req.PageNum - 1)).
		Find(&list)
	return list
}

// 查询文章的评论数量
func (*Comment) GetArticleCommentCount(articleId int) (count int64) {
	DB.Model(&model.Comment{}).
		Where("topic_id = ? AND type = 1 AND is_review = 1", articleId).
		Count(&count)
	return
}
