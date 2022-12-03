package service

import (
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/req"
	"gin-blog/model/resp"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"strconv"
)

type Comment struct{}

// TODO: reply
func (*Comment) GetList(req req.GetComments) resp.PageResult[[]resp.CommentVo] {
	var list = make([]resp.CommentVo, 0)

	total := commentDao.GetCount(req)
	if total != 0 {
		list = commentDao.GetList(req)
	}

	return resp.PageResult[[]resp.CommentVo]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    total,
		List:     list,
	}
}

func (*Comment) Delete(ids []int) (code int) {
	dao.Delete(model.Comment{}, "id IN ?", ids)
	return r.OK
}

// 修改评论审核
func (*Comment) UpdateReview(req req.UpdateReview) (code int) {
	maps := map[string]any{"is_review": req.IsReview}
	dao.UpdatesMap(&model.Comment{}, maps, "id IN ?", req.Ids)
	return r.OK
}

// 前台

// 前台文章评论列表
func (*Comment) GetFrontList(req req.GetFrontComments) resp.PageResult[[]resp.FrontCommentVO] {
	// 数据库查询评论列表
	commentList, total := commentDao.GetFrontList(req)

	// 统计 [评论 id 列表]
	commentIds := make([]int, 0)
	for _, comment := range commentList {
		commentIds = append(commentIds, comment.ID)
	}

	// [文章ID : 对应点赞数]
	likeCountMap := utils.Redis.HGetAll(KEY_COMMENT_LIKE_COUNT)

	// 获取回复列表
	replyList := commentDao.GetReplyList(commentIds)
	replyMap := make(map[int][]resp.ReplyVO)
	for i, reply := range replyList {
		replyList[i].LikeCount, _ = strconv.Atoi(likeCountMap[strconv.Itoa(reply.ID)])
		// 默认只添加 3 条回复, FIXME: SQL 中实现？
		if len(replyMap[reply.ParentId]) < 3 {
			replyMap[reply.ParentId] = append(replyMap[reply.ParentId], reply)
		}
	}

	// 获取评论的回复数量
	replyCountList := commentDao.GetReplyCountListByCommentId(commentIds)
	replyCountMap := make(map[int]int) // [评论id : 回复数量]
	for _, reply := range replyCountList {
		replyCountMap[reply.CommentId] = reply.ReplyCount
	}

	for i, comment := range commentList {
		commentList[i].LikeCount, _ = strconv.Atoi(likeCountMap[strconv.Itoa(comment.ID)])
		commentList[i].ReplyVOList = replyMap[comment.ID]
		commentList[i].ReplyCount = replyCountMap[comment.ID]
		// 给回复添加 点赞数量
		for j, reply := range commentList[i].ReplyVOList {
			commentList[i].ReplyVOList[j].LikeCount, _ = strconv.Atoi(likeCountMap[strconv.Itoa(reply.ID)])
		}
	}

	return resp.PageResult[[]resp.FrontCommentVO]{
		List:     commentList,
		Total:    total,
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
}

// 点赞评论
func (*Comment) SaveLike(uid, commentId int) (code int) {
	// 记录某个用户已经对某个评论点过赞
	commentLikeUserKey := KEY_COMMENT_USER_LIKE_SET + strconv.Itoa(uid)
	// 该评论已经被记录过, 再点赞就是取消点赞
	if utils.Redis.SIsMember(commentLikeUserKey, commentId) {
		utils.Redis.SRem(commentLikeUserKey, commentId)
		utils.Redis.HIncrBy(KEY_COMMENT_LIKE_COUNT, strconv.Itoa(commentId), -1)
	} else { // 未被记录过, 则是增加点赞
		utils.Redis.SAdd(commentLikeUserKey, commentId)
		utils.Redis.HIncrBy(KEY_COMMENT_LIKE_COUNT, strconv.Itoa(commentId), 1)
	}
	return r.OK
}

// 新增评论
func (*Comment) Save(uid int, req req.SaveComment) (code int) {
	// TODO: HTMLUtil.Filter 过滤 HTML 元素中的字符串...
	// re := regexp2.MustCompile(`(?!<(img|p|span|h1|h2|h3|h4|h5|h6).*?>)<.*?>`, 0)
	// src := req.Content
	// src, _ = re.Replace(src, "", 0, len(src))
	// req.Content = strings.TrimSpace(src)

	comment := utils.CopyProperties[model.Comment](req) // vo -> po
	comment.UserId = uid
	// 根据 blogConfig 中的配置设置默认是否审核
	isCommentReview := blogInfoService.GetBlogConfig().IsCommentReview
	comment.IsReview = &isCommentReview
	dao.Create(&comment)
	// TODO: 判断是否开启邮箱通知用户
	return r.OK
}

// 根据 [评论id] 获取 [回复列表]
func (*Comment) GetReplyListByCommentId(id int, req req.PageQuery) []resp.ReplyVO {
	replyList := commentDao.GetReplyListByCommentId(id, req)
	likeCountMap := utils.Redis.HGetAll(KEY_COMMENT_LIKE_COUNT)
	for i, reply := range replyList {
		replyList[i].LikeCount, _ = strconv.Atoi(likeCountMap[strconv.Itoa(reply.ID)])
	}
	return replyList
}
