package model

import "gorm.io/gorm"

// 登录状态
const (
	LOGIN_SUCCESS = 1
	LOGIN_FAIL    = 2
)

// 登录日志
//
// 失败的记录同样要留: 只有成功记录看不出撞库, 而 user_auth 上只有"最后一次登录"
// 没有历史。失败时拿不到用户 id 和昵称, 只有请求里的用户名, 所以这两个字段允许为空。
type LoginLog struct {
	Model

	UserId    int    `gorm:"comment:用户ID, 登录失败时为 0" json:"user_id"`
	Username  string `gorm:"type:varchar(50);comment:登录用户名" json:"username"`
	Nickname  string `gorm:"type:varchar(50);comment:用户昵称" json:"nickname"`
	IpAddress string `gorm:"type:varchar(255);comment:登录IP" json:"ip_address"`
	IpSource  string `gorm:"type:varchar(255);comment:登录地址" json:"ip_source"`
	Status    int    `gorm:"type:tinyint;comment:状态(1-成功 2-失败)" json:"status"`
	Message   string `gorm:"type:varchar(255);comment:失败原因" json:"message"`
}

func AddLoginLog(db *gorm.DB, log *LoginLog) error {
	return db.Create(log).Error
}

// 关键字匹配用户名、昵称与登录 IP
func GetLoginLogList(db *gorm.DB, num, size int, keyword string) (data []LoginLog, total int64, err error) {
	db = db.Model(&LoginLog{})
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("username LIKE ? OR nickname LIKE ? OR ip_address LIKE ?", like, like, like)
	}
	db.Count(&total)
	result := db.Order("created_at DESC").
		Scopes(Paginate(num, size)).
		Find(&data)
	return data, total, result.Error
}
