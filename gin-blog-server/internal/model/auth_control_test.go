package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func initModelDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //
		},
	})
	if err != nil {
		return nil, err
	}

	MakeMigrate(db)
	return db, nil
}

func TestAuth(t *testing.T) {
	db, _ := initModelDB()

	p1, _ := AddResource(db, "api_1", "/v1/api_1", "GET", false)
	p2, _ := AddResource(db, "api_2", "/v1/api_2", "POST", false)

	// 添加拥有两个资源的角色
	role1, _ := AddRoleWithResources(db, "admin", "管理员", []int{p1.ID, p2.ID})
	resources, _ := GetResourcesByRole(db, role1.ID)
	assert.Len(t, resources, 2)

	// 修改角色资源为一个
	role2, _ := UpdateRoleWithResources(db, role1.ID, "super", "超管", []int{p1.ID})
	resources, _ = GetResourcesByRole(db, role2.ID)
	assert.Len(t, resources, 1)

	// 测试角色资源鉴权
	{
		flag, _ := CheckRoleAuth(db, role2.ID, "/v1/api_1", "GET")
		assert.True(t, flag)
		flag, _ = CheckRoleAuth(db, role2.ID, "/v1/api_99", "POST")
		assert.False(t, flag)
	}
}

func TestCheckRoleAuthMethod(t *testing.T) {
	db, _ := initModelDB()

	r, _ := AddResource(db, "api_get", "/v1/api", "GET", false)
	role, _ := AddRoleWithResources(db, "role", "角色", []int{r.ID})

	// url 相同但请求方法不同, 不应通过
	flag, err := CheckRoleAuth(db, role.ID, "/v1/api", "POST")
	assert.Nil(t, err)
	assert.False(t, flag)

	// 角色不存在时不应报错, 直接判定为无权限
	flag, _ = CheckRoleAuth(db, 999, "/v1/api", "GET")
	assert.False(t, flag)
}

// 匿名资源只代表该资源本身不需要登录, 不能让拥有它的角色获得所有权限
func TestCheckRoleAuthAnonymous(t *testing.T) {
	db, _ := initModelDB()

	anon, _ := AddResource(db, "public", "/v1/public", "GET", true)
	secret, _ := AddResource(db, "secret", "/v1/secret", "DELETE", false)
	role, _ := AddRoleWithResources(db, "role", "角色", []int{anon.ID})

	// 拥有匿名资源, 不能越权访问其他资源
	flag, err := CheckRoleAuth(db, role.ID, "/v1/secret", "DELETE")
	assert.Nil(t, err)
	assert.False(t, flag)

	// 匿名资源自身仍然能通过
	flag, _ = CheckRoleAuth(db, role.ID, "/v1/public", "GET")
	assert.True(t, flag)

	// 显式分配后才有权限
	UpdateRoleWithResources(db, role.ID, "role", "角色", []int{anon.ID, secret.ID})
	flag, _ = CheckRoleAuth(db, role.ID, "/v1/secret", "DELETE")
	assert.True(t, flag)
}

func TestResourceQuery(t *testing.T) {
	db, _ := initModelDB()

	r, err := AddResource(db, "api", "/v1/api", "GET", false)
	assert.Nil(t, err)

	// 按 url + method 查询
	got, err := GetResource(db, "/v1/api", "GET")
	assert.Nil(t, err)
	assert.Equal(t, r.ID, got.ID)

	// 按 id 查询
	got, err = GetResourceById(db, r.ID)
	assert.Nil(t, err)
	assert.Equal(t, "api", got.Name)

	// 不存在的资源返回 ErrRecordNotFound, 中间件依赖该错误跳过鉴权
	_, err = GetResource(db, "/v1/not_exist", "GET")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestResourceAnonymousUpdate(t *testing.T) {
	db, _ := initModelDB()

	r, _ := AddResource(db, "api", "/v1/api", "GET", false)
	assert.Nil(t, UpdateResourceAnonymous(db, r.ID, true))

	got, _ := GetResourceById(db, r.ID)
	assert.True(t, got.Anonymous)
}

func TestResourceInUseAndChild(t *testing.T) {
	db, _ := initModelDB()

	parent, _ := AddResource(db, "parent", "", "", false)
	child := Resource{Name: "child", ParentId: parent.ID, Url: "/v1/child", Method: "GET"}
	db.Create(&child)

	// 父资源有子资源, 子资源没有
	has, err := CheckResourceHasChild(db, parent.ID)
	assert.Nil(t, err)
	assert.True(t, has)
	has, _ = CheckResourceHasChild(db, child.ID)
	assert.False(t, has)

	// 未被角色引用
	inUse, err := CheckResourceInUse(db, child.ID)
	assert.Nil(t, err)
	assert.False(t, inUse)

	// 被角色引用后
	AddRoleWithResources(db, "role", "角色", []int{child.ID})
	inUse, _ = CheckResourceInUse(db, child.ID)
	assert.True(t, inUse)
}

func TestDeleteResourceAndRole(t *testing.T) {
	db, _ := initModelDB()

	r, _ := AddResource(db, "api", "/v1/api", "GET", false)
	role, _ := AddRoleWithResources(db, "role", "角色", []int{r.ID})
	db.Create(&RoleMenu{RoleId: role.ID, MenuId: 1})

	// 删除资源
	count, err := DeleteResource(db, r.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
	_, err = GetResourceById(db, r.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	// 删除角色时, 关联的 role_resource 和 role_menu 一并清理
	assert.Nil(t, DeleteRole(db, role.ID))

	var resourceCount, menuCount int64
	db.Model(&RoleResource{}).Where("role_id = ?", role.ID).Count(&resourceCount)
	db.Model(&RoleMenu{}).Where("role_id = ?", role.ID).Count(&menuCount)
	assert.Equal(t, int64(0), resourceCount)
	assert.Equal(t, int64(0), menuCount)
}
