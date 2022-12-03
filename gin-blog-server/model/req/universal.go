package req

// 获取树形数据(不需要分页)
type TreeQuery struct {
	Keyword string `form:"keyword"`
}

// 获取数据(需要分页)
type PageQuery struct {
	PageSize int    `form:"page_size"`
	PageNum  int    `form:"page_num"`
	Keyword  string `form:"keyword"`
}

// 软删除请求(批量)
type SoftDelete struct {
	Ids      []int `json:"ids"`
	IsDelete *int8 `json:"is_delete" validate:"required,min=0,max=1"` // 软删除到回收站, 没有的字段不使用
}

// 修改审核(批量)
type UpdateReview struct {
	Ids      []int `json:"ids"`
	IsReview *int8 `json:"is_review" validate:"required,min=0,max=1"`
}
