package req

type SaveOrUpdateMenu struct {
	ID        int    `json:"id"`
	Name      string `json:"name" mapstructure:"name"`
	Path      string `json:"path" mapstructure:"path"`
	Component string `json:"component" mapstructure:"component"`
	Icon      string `json:"icon" mapstructure:"icon"`
	OrderNum  *int8  `json:"order_num" validate:"required,min=0" mapstructure:"order_num"`
	ParentId  int    `json:"parent_id" mapstructure:"parent_id"`
	IsHidden  *int8  `json:"is_hidden" mapstructure:"is_hidden"`
	KeepAlive *int8  `json:"keep_alive" mapstructure:"keep_alive"`
	Redirect  string `json:"redirect" mapstructure:"redirect"`
}
