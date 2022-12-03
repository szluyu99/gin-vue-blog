package resp

import "time"

type TagVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Articles     []*model.Article `json:"articles"`
	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}

func (TagVO) TableName() string {
	return "tag"
}
