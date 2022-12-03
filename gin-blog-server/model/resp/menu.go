package resp

import "time"

// 用户菜单
type UserMenuVo struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Path      string       `json:"path"`
	Component string       `json:"component"`
	OrderNum  int8         `json:"order_num"`
	Icon      string       `json:"icon"`
	Children  []UserMenuVo `json:"children"`
	ParentId  int          `json:"parent_id"`
	Hidden    bool         `json:"is_hidden"` // *
	KeepAlive bool         `json:"keep_alive"`
	Redirect  string       `json:"redirect"`
}

type MenuVo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	OrderNum  int8      `json:"order_num"`
	Children  []MenuVo  `json:"children"`
	ParentId  int       `json:"parent_id"`
	IsHidden  int       `json:"is_hidden"`
	KeepAlive int8      `json:"keep_alive"`
	Redirect  string    `json:"redirect"`
}
