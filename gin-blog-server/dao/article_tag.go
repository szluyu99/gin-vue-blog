package dao

import (
	"gin-blog/model"
)

type ArticleTag struct{}

// 新增 文章-分类 的关联关系
func (*ArticleTag) CreateArticleTag(artId, tagId int) error {
	data := &model.ArticleTag{
		ArticleId: artId,
		TagId:     tagId,
	}
	err := DB.Create(data).Error
	return err
}

// 批量新增/编辑 文章-分类的关联关系
func (*ArticleTag) SaveOrUpdateArticlesTags(artId int, tagIds []int) {
	// 删除该文章原本对应的标签
	if artId != 0 {
		Delete(model.ArticleTag{}, "article_id = ?", artId)
	}
}
