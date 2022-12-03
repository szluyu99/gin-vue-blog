package req

// type GetFriendLinks struct {
// 	PageQuery
// }

type SaveOrUpdateLink struct {
	ID      int    `json:"id"`
	Name    string `json:"name" validate:"required" label:"友链名称"`
	Avatar  string `json:"avatar"`
	Address string `json:"address" validate:"required" label:"友链地址"`
	Intro   string `json:"intro"`
}
