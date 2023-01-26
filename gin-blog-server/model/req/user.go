package req

// 查询用户(后台)
type GetUsers struct {
	PageQuery
	LoginType int8   `form:"login_type"`
	Nickname  string `form:"nickname"`
}

// 更新用户(后台)
type UpdateUser struct {
	UserInfoId int    `json:"id"` // 注意这里的别名
	Nickname   string `json:"nickname"`
	RoleIds    []int  `json:"role_ids"`
}

// 更新当前用户信息
type UpdateCurrentUser struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname" validate:"required"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`
}

// 修改用户禁用状态
type UpdateUserDisable struct {
	ID        int  `json:"id"`
	IsDisable *int `json:"is_disable" validate:"required,min=0,max=1"`
}

type Register struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required,min=4,max=20" label:"密码"`
	Code     string `json:"code" validate:"required" label:"邮箱验证码"`
}

type Login struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required" label:"密码"`
}

// 修改普通用户密码
type UpdatePassword struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required" label:"密码"`
}

// 修改管理员密码: 需要旧密码验证
type UpdateAdminPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// 更新邮箱信息
type UpdateEmail struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// 强制下线需要用户信息计算其 uuid
type ForceOfflineUser struct {
	UserIndoId int    `json:"user_info_id"`
	IpAddress  string `json:"ip_address"`
	Browser    string `json:"browser"`
	OS         string `json:"os"`
}
