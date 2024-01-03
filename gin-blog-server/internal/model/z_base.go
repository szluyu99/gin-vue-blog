package model

import (
	"time"

	"gorm.io/gorm"
)

// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
// 只支持创建表、增加表中没有的字段和索引
// 为了保护数据，并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate(db *gorm.DB) error {
	// 设置表关联
	db.SetupJoinTable(&Role{}, "Menus", &RoleMenu{})
	db.SetupJoinTable(&Role{}, "Resources", &RoleResource{})
	db.SetupJoinTable(&Role{}, "Users", &UserAuthRole{})
	db.SetupJoinTable(&UserAuth{}, "Roles", &UserAuthRole{})

	return db.AutoMigrate(
		&Article{},      // 文章
		&Category{},     // 分类
		&Tag{},          // 标签
		&Comment{},      // 评论
		&Message{},      // 消息
		&FriendLink{},   // 友链
		&Page{},         // 页面
		&Config{},       // 网站设置
		&OperationLog{}, // 操作日志
		&UserInfo{},     // 用户信息

		&UserAuth{},     // 用户验证
		&Role{},         // 角色
		&Menu{},         // 菜单
		&Resource{},     // 资源（接口）
		&RoleMenu{},     // 角色-菜单 关联
		&RoleResource{}, // 角色-资源 关联
		&UserAuthRole{}, // 用户-角色 关联
	)
}

// 通用模型

type Model struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OptionVO struct {
	ID   int    `json:"value"`
	Name string `json:"label"`
}

// Gorm Scopes

// 分页
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// 通用 CRUD

// 创建数据(可以创建[单条]数据, 也可[批量]创建)
func Create[T any](db *gorm.DB, data *T) (*T, error) {
	result := db.Create(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// [单条]数据查询
func Get[T any](db *gorm.DB, data *T, query string, args ...any) (*T, error) {
	result := db.Where(query, args...).First(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// [单行]更新: 传入对应结构体[传递主键用] 和 带有对应更新字段值的[结构体]，结构体不更新零值
func Update[T any](db *gorm.DB, data T, slt ...string) error {
	db = db.Model(&data)
	if len(slt) > 0 {
		db = db.Select(slt)
	}
	result := db.Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// [批量]更新: map 的字段就是要更新的字段 (map 可以更新零值), 通过条件可以实现[单行]更新
func UpdatesMap[T any](db *gorm.DB, data *T, maps map[string]any, query string, args ...any) error {
	result := db.Model(data).Where(query, args...).Updates(maps)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// [批量]更新: 结构体的属性就是要更新的字段 (结构体不更新零值), 通过条件可以实现[单行]更新
func Updates[T any](db *gorm.DB, data T, query string, args ...any) error {
	result := db.Model(&data).Where(query, args...).Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 数据列表
func List[T any](db *gorm.DB, data T, slt, order, query string, args ...any) (T, error) {
	db = db.Model(data).Select(slt).Order(order)
	if query != "" {
		db = db.Where(query, args...)
	}
	result := db.Find(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

// [批量]删除数据, 通过条件控制可以删除单条数据
func Delete[T any](db *gorm.DB, data T, query string, args ...any) error {
	result := db.Where(query, args...).Delete(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 统计数量
func Count[T any](db *gorm.DB, data *T, where ...any) (int, error) {
	var total int64
	db = db.Model(data)
	if len(where) > 0 {
		db = db.Where(where[0], where[1:]...)
	}
	result := db.Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(total), nil
}
