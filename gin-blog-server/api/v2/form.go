package v2

// 获取数据（需要分页）
type PageQuery struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
}

// 分页返回数据
type PageResult[T any] struct {
	PageSize int   `json:"pageSize"` // 上次页数
	PageNum  int   `json:"pageNum"`  // 每页条数
	Total    int64 `json:"total"`    // 总条数
	List     T     `json:"pageData"` // 分页数据
}
