package resp

import "time"

// 前后台通用
type CategoryVo struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ArticleCount int       `json:"article_count"` // 文章数量
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
