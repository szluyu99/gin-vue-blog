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
