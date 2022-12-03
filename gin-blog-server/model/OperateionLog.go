package model

// 操作日志
type OperationLog struct {
	Universal
	OptModule string `gorm:"type:varchar(50);comment:操作模块" json:"opt_module"`
	OptType   string `gorm:"type:varchar(50);comment:操作类型" json:"opt_type"`
	OptMethod string `gorm:"type:varchar(100);comment:操作方法" json:"opt_method"`
	OptUrl    string `gorm:"type:varchar(255);comment:操作URL" json:"opt_url"`
	OptDesc   string `gorm:"type:varchar(255);comment:操作描述" json:"opt_desc"`

	RequestParam  string `gorm:"type:longtext;comment:请求参数" json:"request_param"`
	RequestMethod string `gorm:"type:longtext;comment:请求方法" json:"request_method"`
	ResponseData  string `gorm:"type:longtext;comment:响应数据" json:"response_data"`

	UserId    int    `gorm:"comment:用户ID" json:"user_id"`
	Nickname  string `gorm:"comment:用户昵称" json:"nickname"`
	IpAddress string `gorm:"comment:操作IP" json:"ip_address"`
	IpSource  string `gorm:"comment:操作地址" json:"ip_source"`
}
