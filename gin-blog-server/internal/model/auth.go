package model

import (
	"encoding/json"
	"errors"
	"gin-blog/internal/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// 注册时邮箱已存在, 由 handle 层翻译成对应的业务错误码
var ErrUsernameTaken = errors.New("该邮箱已经注册")

// 权限控制: 7 张表（4 模型 + 3 关联）

type UserAuth struct {
	Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"`
	Password      string     `gorm:"type:varchar(100)" json:"-"`
	LoginType     int        `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`
	IpAddress     string     `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"`
	IpSource      string     `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`
	LastLoginTime *time.Time `json:"last_login_time"`
	IsDisable     bool       `json:"is_disable"`
	IsSuper       bool       `json:"is_super"` // 超级管理员只能后台设置

	UserInfoId int       `json:"user_info_id"`
	UserInfo   *UserInfo `json:"info"`
	Roles      []*Role   `json:"roles" gorm:"many2many:user_auth_role"`
}

func (u *UserAuth) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

type Role struct {
	Model
	Name      string `gorm:"unique" json:"name"`
	Label     string `gorm:"unique" json:"label"`
	IsDisable bool   `json:"is_disable"`

	Resources []Resource `json:"resources" gorm:"many2many:role_resource"`
	Menus     []Menu     `json:"menus" gorm:"many2many:role_menu"`
	Users     []UserAuth `json:"users" gorm:"many2many:user_auth_role"`
}

type Resource struct {
	Model
	Name      string `gorm:"unique;type:varchar(50)" json:"name"`
	ParentId  int    `json:"parent_id"`
	Url       string `gorm:"type:varchar(255);index:idx_resource_api,priority:1" json:"url"`
	Method    string `gorm:"type:varchar(10);index:idx_resource_api,priority:2" json:"request_method"`
	Anonymous bool   `json:"is_anonymous"`

	Roles []*Role `json:"roles" gorm:"many2many:role_resource"`
}

/*
菜单设计:

目录: catalogue === true
  - 如果是目录, 作为单独项, 不展开子菜单（例如 "首页", "个人中心"）
  - 如果不是目录, 且 parent_id 为 0, 则为一级菜单, 可展开子菜单（例如 "文章管理" 下有 "文章列表", "文章分类", "文章标签" 等子菜单）
  - 如果不是目录, 且 parent_id 不为 0, 则为二级菜单

隐藏: hidden
  - 隐藏则不显示在菜单栏中

外链: external, external_link
  - 如果是外链, 如果设置为外链, 则点击后会在新窗口打开
*/
type Menu struct {
	Model
	ParentId     int    `json:"parent_id"`
	Name         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(20)" json:"name"` // 菜单名称
	Path         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(50)" json:"path"` // 路由地址
	Component    string `gorm:"type:varchar(50)" json:"component"`                          // 组件路径
	Icon         string `gorm:"type:varchar(50)" json:"icon"`                               // 图标
	OrderNum     int8   `json:"order_num"`                                                  // 排序
	Redirect     string `gorm:"type:varchar(50)" json:"redirect"`                           // 重定向地址
	Catalogue    bool   `json:"is_catalogue"`                                               // 是否为目录
	Hidden       bool   `json:"is_hidden"`                                                  // 是否隐藏
	KeepAlive    bool   `json:"keep_alive"`                                                 // 是否缓存
	External     bool   `json:"is_external"`                                                // 是否外链
	ExternalLink string `gorm:"type:varchar(255)" json:"external_link"`                     // 外链地址

	Roles []*Role `json:"roles" gorm:"many2many:role_menu"`
}

type RoleResource struct {
	RoleId     int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
	ResourceId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
}

type UserAuthRole struct {
	UserAuthId int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
	RoleId     int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
}

type RoleMenu struct {
	RoleId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
	MenuId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
}

type RoleVO struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	IsDisable   bool      `json:"is_disable"`
	ResourceIds []int     `json:"resource_ids" gorm:"-"`
	MenuIds     []int     `json:"menu_ids" gorm:"-"`
}

// Menu

func SaveOrUpdateMenu(db *gorm.DB, menu *Menu) error {
	var result *gorm.DB

	if menu.ID > 0 {
		result = db.Model(menu).
			Select("name", "path", "component", "icon", "redirect", "parent_id", "order_num", "catalogue", "hidden", "keep_alive", "external").
			Updates(menu)
	} else {
		result = db.Create(menu)
	}

	return result.Error
}

func GetMenuIdsByRoleId(db *gorm.DB, roleId int) (ids []int, err error) {
	result := db.Model(&RoleMenu{}).Where("role_id = ?", roleId).Pluck("menu_id", &ids)
	return ids, result.Error
}

func GetMenuById(db *gorm.DB, id int) (menu *Menu, err error) {
	result := db.First(&menu, id)
	return menu, result.Error
}

func CheckMenuInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&RoleMenu{}).Where("menu_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func CheckMenuHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&Menu{}).Where("parent_id = ?", id).Count(&count)
	return count > 0, result.Error
}

// 获取所有菜单列表（超级管理员用）
func GetAllMenuList(db *gorm.DB) (menu []Menu, err error) {
	result := db.Find(&menu)
	return menu, result.Error
}

// 根据 user_id 获取菜单列表
func GetMenuListByUserId(db *gorm.DB, id int) (menus []Menu, err error) {
	var userAuth UserAuth
	result := db.Where(&UserAuth{Model: Model{ID: id}}).
		Preload("Roles").Preload("Roles.Menus").
		First(&userAuth)

	if result.Error != nil {
		return nil, result.Error
	}

	set := make(map[int]Menu)
	for _, role := range userAuth.Roles {
		for _, menu := range role.Menus {
			set[menu.ID] = menu
		}
	}

	for _, menu := range set {
		menus = append(menus, menu)
	}

	return menus, nil
}

func GetMenuList(db *gorm.DB, keyword string) (list []Menu, total int64, err error) {
	db = db.Model(&Menu{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	result := db.Count(&total).Find(&list)
	return list, total, result.Error
}

func DeleteMenu(db *gorm.DB, id int) error {
	result := db.Delete(&Menu{}, id)
	return result.Error
}

// Resource

func SaveOrUpdateResource(db *gorm.DB, id, pid int, name, url, method string) error {
	resource := Resource{
		Model:    Model{ID: id},
		Name:     name,
		Url:      url,
		Method:   method,
		ParentId: pid,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&resource)
	} else {
		result = db.Create(&resource)
		// TODO: ????
		// * 解决前端的 BUG: 级联选中某个父节点后, 新增的子节点默认会展示被选中, 实际上未被选中值
		// * 解决方案: 新增子节点后, 删除该节点对应的父节点与角色的关联关系
		// dao.Delete(model.RoleResource{}, "resource_id", data.ParentId)
	}
	return result.Error
}

func GetResourceIdsByRoleId(db *gorm.DB, roleId int) (ids []int, err error) {
	result := db.Model(&RoleResource{}).
		Where("role_id = ?", roleId).
		Pluck("resource_id", &ids)
	return ids, result.Error
}

func GetResourceList(db *gorm.DB, keyword string) (list []Resource, err error) {
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}

	result := db.Find(&list)
	return list, result.Error
}

func GetResourceListByIds(db *gorm.DB, ids []int) (list []Resource, err error) {
	result := db.Where("id in ?", ids).Find(&list)
	return list, result.Error
}

// Role

func SaveOrUpdateRole(db *gorm.DB, id int, name, label string, isDisable bool) error {
	role := Role{
		Model:     Model{ID: id},
		Name:      name,
		Label:     label,
		IsDisable: isDisable,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&role)
	} else {
		result = db.Create(&role)
	}

	return result.Error
}

func GetRoleOption(db *gorm.DB) (list []OptionVO, err error) {
	result := db.Model(&Role{}).Select("id", "name").Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func GetRoleList(db *gorm.DB, num, size int, keyword string) (list []RoleVO, total int64, err error) {
	db = db.Model(&Role{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	db.Count(&total)
	result := db.Select("id", "name", "label", "created_at", "is_disable").
		Scopes(Paginate(num, size)).
		Find(&list)
	return list, total, result.Error
}

func GetRoleIdsByUserId(db *gorm.DB, userAuthId int) (ids []int, err error) {
	result := db.
		Model(&UserAuthRole{UserAuthId: userAuthId}).
		Pluck("role_id", &ids)
	return ids, result.Error
}

func SaveRole(db *gorm.DB, name, label string, resourceIds, menuIds []int) error {
	role := Role{
		Name:  name,
		Label: label,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		// 新建时以前会把 resourceIds / menuIds 丢掉, 新角色一个权限都没有
		return replaceRoleRelations(tx, role.ID, resourceIds, menuIds)
	})
}

func UpdateRole(db *gorm.DB, id int, name, label string, isDisable bool, resourceIds, menuIds []int) error {
	role := Role{
		Model:     Model{ID: id},
		Name:      name,
		Label:     label,
		IsDisable: isDisable,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).Select("name", "label", "is_disable").Updates(&role).Error; err != nil {
			return err
		}
		return replaceRoleRelations(tx, role.ID, resourceIds, menuIds)
	})
}

// 用给定的 id 列表整体替换角色的资源与菜单关联, 必须在事务里调用
func replaceRoleRelations(tx *gorm.DB, roleId int, resourceIds, menuIds []int) error {
	if err := tx.Delete(&RoleResource{}, "role_id = ?", roleId).Error; err != nil {
		return err
	}
	if len(resourceIds) > 0 {
		rows := make([]RoleResource, 0, len(resourceIds))
		for _, rid := range resourceIds {
			rows = append(rows, RoleResource{RoleId: roleId, ResourceId: rid})
		}
		if err := tx.Create(&rows).Error; err != nil {
			return err
		}
	}

	if err := tx.Delete(&RoleMenu{}, "role_id = ?", roleId).Error; err != nil {
		return err
	}
	if len(menuIds) > 0 {
		rows := make([]RoleMenu, 0, len(menuIds))
		for _, mid := range menuIds {
			rows = append(rows, RoleMenu{RoleId: roleId, MenuId: mid})
		}
		if err := tx.Create(&rows).Error; err != nil {
			return err
		}
	}
	return nil
}

// 删除角色: 事务删除 role, role_resource, role_menu
func DeleteRoles(db *gorm.DB, ids []int) error {
	return db.Transaction(func(tx *gorm.DB) error {

		result := tx.Delete(&Role{}, "id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Delete(&RoleResource{}, "role_id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Delete(&RoleMenu{}, "role_id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

// UserAuth

func GetUserAuthInfoById(db *gorm.DB, id int) (*UserAuth, error) {
	var userAuth = UserAuth{Model: Model{ID: id}}
	result := db.Model(&userAuth).
		Preload("Roles").Preload("UserInfo").
		First(&userAuth)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userAuth, nil
}

// 新注册用户的默认头像: 原来是 bing 的一张外链图, 换成本仓库 images/ 下的图片
const DefaultAvatar = "https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images/config/user_avatar.jpeg"

// 注册新用户
//
// 三张表必须一起成功: 以前没有事务, 中途失败会留下孤儿 user_info,
// 或者一个没有任何角色的用户(能登录, 但 PermissionCheck 查不到角色)。
// 昵称也不再按 Count 生成, 并发注册会重名, 改用插入后拿到的自增 id。
func CreateNewUser(db *gorm.DB, username, password string) (*UserAuth, *UserInfo, *UserAuthRole, error) {
	pass, err := utils.BcryptHash(password)
	if err != nil {
		return nil, nil, nil, err
	}

	userinfo := &UserInfo{
		Email:  username,
		Avatar: DefaultAvatar,
	}
	userauth := &UserAuth{}
	userRole := &UserAuthRole{}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 邮箱验证链接是一次性的, 但同一个邮箱可能有多封未过期的邮件,
		// 这里再查一次, 避免两个链接都点导致建出两个账号
		var exist int64
		if err := tx.Model(&UserAuth{}).Where("username = ?", username).Count(&exist).Error; err != nil {
			return err
		}
		if exist > 0 {
			return ErrUsernameTaken
		}

		if err := tx.Create(userinfo).Error; err != nil {
			return err
		}
		// 拿到 id 之后再补昵称和简介, 保证唯一
		number := strconv.Itoa(userinfo.ID)
		userinfo.Nickname = "游客" + number
		userinfo.Intro = "我是这个程序的第" + number + "个用户"
		if err := tx.Model(userinfo).
			Select("nickname", "intro").Updates(userinfo).Error; err != nil {
			return err
		}

		userauth.Username = username
		userauth.Password = pass
		userauth.UserInfoId = userinfo.ID
		if err := tx.Create(userauth).Error; err != nil {
			return err
		}

		userRole.UserAuthId = userauth.ID
		userRole.RoleId = 2 // 默认身份为游客
		return tx.Create(userRole).Error
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return userauth, userinfo, userRole, nil
}
