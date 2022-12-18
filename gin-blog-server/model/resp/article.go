package resp

import (
	"gin-blog/model"
	"time"
)

// 文章列表 VO
type ArticleVO struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CategoryId int            `json:"category_id"`
	Category   model.Category `gorm:"foreignkey:CategoryId" json:"category"`
	Tags       []model.Tag    `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`

	UserId      int    `json:"user_id"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Content     string `json:"content"`
	Img         string `json:"img"`
	Type        int8   `json:"type"`
	Status      int8   `json:"status"`
	IsTop       *int8  `json:"is_top"`
	IsDelete    *int8  `json:"is_delete"`
	OriginalUrl string `json:"original_url"`

	LikeCount int `json:"like_count"`
	ViewCount int `json:"view_count"`
	// CommentCount int `json:"comment_count"`
}

// 文章详情 VO
type ArticleDetailVO struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Desc         string   `json:"desc"`
	Content      string   `json:"content"`
	Img          string   `json:"img"`
	CategoryName *string  `json:"category_name"`
	TagNames     []string `json:"tag_names"`
	Type         int      `json:"type"`
	OriginalUrl  string   `json:"original_url"`
	IsTop        int      `json:"is_top"`
	Status       int      `json:"status"`
}

// 前台首页文章列表 VO
// manyToMany 参考: https://gorm.io/zh_CN/docs/many_to_many.html
type FrontArticleVO struct {
	ID         int             `json:"id"`
	CreatedAt  time.Time       `json:"created_at"`
	Title      string          `json:"title"`
	Desc       string          `json:"desc"`
	Content    string          `json:"content"`
	Img        string          `json:"img"`
	IsTop      int             `json:"is_top"`
	Type       int             `json:"type"`
	CategoryId int             `json:"category_id"`
	Category   *model.Category `gorm:"foreignkey:CategoryId;" json:"category"`
	Tags       []*model.Tag    `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`
}

// 前端文章详情 VO
type FrontArticleDetailVO struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Img         string `json:"img"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Type        int    `json:"type"`
	OriginalUrl string `json:"original_url"`

	CommentCount int `json:"comment_count"` // 评论数量
	LikeCount    int `json:"like_count"`    // 点赞数量
	ViewCount    int `json:"view_count"`    // 访问数量

	CategoryId int            `json:"category_id"`
	Category   model.Category `gorm:"foreignkey:CategoryId;" json:"category"`
	Tags       []model.Tag    `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`

	LastArticle ArticlePaginationVO `gorm:"-" json:"last_article"` // 上一篇
	NextArticle ArticlePaginationVO `gorm:"-" json:"next_article"` // 下一篇

	RecommendArticles []RecommendArticleVO `gorm:"-" json:"recommend_articles"` // 推荐文章
	NewestArticles    []RecommendArticleVO `gorm:"-" json:"newest_articles"`    // 最新文章
	// Desc         string   `json:"desc"`
}

// 文章详情界面: 上一篇文章, 下一篇文章显示, 只需要 标题, 封面
type ArticlePaginationVO struct {
	ID    int    `json:"id"`
	Img   string `json:"img"`
	Title string `json:"title"`
}

// 推荐文章
type RecommendArticleVO struct {
	ID        int       `json:"id"`
	Img       string    `json:"img"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// 文章归档
type ArchiveVO struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Created_at time.Time `json:"created_at"`
}

// 文章搜索结果
type ArticleSearchVO struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
