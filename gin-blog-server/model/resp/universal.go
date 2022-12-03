package resp

// 分页响应结果
type PageResult[T any] struct {
	PageSize int   `json:"pageSize"`
	PageNum  int   `json:"pageNum"`
	Total    int64 `json:"total"`
	List     T     `json:"pageData"` // ! 注意这里的别名
}

// 列表响应结果
type ListResult[T any] struct {
	Total int64 `json:"total"`
	List  T     `json:"list"`
}

type TreeOptionVo struct {
	ID       int            `json:"key"`
	Label    string         `json:"label"`
	Children []TreeOptionVo `json:"children"`
}

type OptionVo struct {
	ID   int    `json:"value"`
	Name string `json:"label"`
}
