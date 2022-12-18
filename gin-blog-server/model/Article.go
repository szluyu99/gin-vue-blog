package model

import (
	"reflect"
)

const (
	PUBLIC = iota + 1 // 公开
	SECRET            // 私密
	DRAFT             // 草稿
)

// belongTo: 一个文章 属于 一个分类
// belongTo: 一个文章 属于 一个用户
// many2many: 一个文章 可以拥有 多个标签, 多个文章 可以使用 一个标签
type Article struct {
	Universal

	CategoryId int `gorm:"type:bigint;not null;comment:分类 ID" json:"category_id"`
	// Category   *Category `gorm:"foreignkey:CategoryId" json:"category"`

	// Tags []*Tag `gorm:"many2many:article_tag;" json:"tags"`

	UserId int `gorm:"type:int;not null;comment:用户 ID" json:"user_id"`
	// User   *User `gorm:"foreignkey:UserId" json:"user"`

	Title       string `gorm:"type:varchar(100);not null;comment:文章标题" json:"title"`
	Desc        string `gorm:"type:varchar(200);comment:文章描述" json:"desc"`
	Content     string `gorm:"type:longtext;comment:文章内容" json:"content"`
	Img         string `gorm:"type:varchar(100);comment:封面图片地址" json:"img"`
	Type        int8   `gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"`
	Status      int8   `gorm:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`
	IsTop       *int8  `gorm:"type:tinyint;not null;default:0;comment:是否置顶(0-否 1-是)" json:"is_top"`
	IsDelete    *int8  `gorm:"type:tinyint;not null;default:0;comment:是否放到回收站(0-否 1-是)" json:"is_delete"`
	OriginalUrl string `gorm:"type:varchar(100);comment:源链接" json:"original_url"`
}

func (a *Article) IsEmpty() bool {
	return reflect.DeepEqual(a, &Article{})
}

// 文章-标签 关联
type ArticleTag struct {
	ArticleId int
	TagId     int
}
