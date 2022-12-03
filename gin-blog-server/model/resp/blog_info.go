package resp

import "gin-blog/model"

type BlogHomeVO struct {
	ArticleCount int64 `json:"article_count"` // 文章数量
	UserCount    int64 `json:"user_count"`    // 用户数量
	MessageCount int64 `json:"message_count"` // 留言数量
	ViewCount    int   `json:"view_count"`    // 访问量
	// CategoryCount int64 `json:"category_count"` // 分类数量
	// TagCount      int64 `json:"tag_count"`      // 标签数量
	// BlogConfig    model.BlogConfigDetail `json:"blog_config"`    // 博客信息
	// PageList      []Page                 `json:"pageList"`
}

type FrontBlogHomeVO struct {
	ArticleCount  int64                  `json:"article_count"`  // 文章数量
	UserCount     int64                  `json:"user_count"`     // 用户数量
	MessageCount  int64                  `json:"message_count"`  // 留言数量
	CategoryCount int64                  `json:"category_count"` // 分类数量
	TagCount      int64                  `json:"tag_count"`      // 标签数量
	BlogConfig    model.BlogConfigDetail `json:"blog_config"`    // 博客信息
	ViewCount     int                    `json:"view_count"`     // 访问量
	// PageList      []Page                 `json:"pageList"`
}
