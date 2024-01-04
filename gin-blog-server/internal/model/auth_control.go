package model

import "gorm.io/gorm"

// 权限控制相关操作

// resource

func DeleteResource(db *gorm.DB, id int) (int, error) {
	result := db.Delete(&Resource{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func AddResource(db *gorm.DB, name, uri, method string, anonymous bool) (*Resource, error) {
	resource := Resource{
		Name:      name,
		Method:    method,
		Url:       uri,
		Anonymous: anonymous,
	}
	result := db.Save(&resource)
	return &resource, result.Error
}

func GetResource(db *gorm.DB, uri, method string) (resource Resource, err error) {
	result := db.Where(&Resource{Url: uri, Method: method}).First(&resource)
	return resource, result.Error
}

func GetResourceById(db *gorm.DB, id int) (resource Resource, err error) {
	result := db.First(&resource, id)
	return resource, result.Error
}

func CheckResourceInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&RoleResource{}).Where("resource_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func CheckResourceHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&Resource{}).Where("parent_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func GetResourcesByRole(db *gorm.DB, rid int) (resources []Resource, err error) {
	var role Role
	result := db.Model(&Role{}).Preload("Resources").Take(&role, rid)
	return role.Resources, result.Error
}

func UpdateResourceAnonymous(db *gorm.DB, id int, anonymous bool) error {
	result := db.Model(&Resource{}).Where("id = ?", id).Update("anonymous", anonymous)
	return result.Error
}

// role

func AddRoleWithResources(db *gorm.DB, name, label string, rs []int) (*Role, error) {
	role := Role{
		Name:  name,
		Label: label,
	}

	result := db.Create(&role)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, rid := range rs {
		result := db.Create(&RoleResource{
			RoleId:     role.ID,
			ResourceId: rid,
		})
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return &role, nil
}

func UpdateRoleWithResources(db *gorm.DB, id int, name, label string, rs []int) (*Role, error) {
	role := Role{
		Model: Model{ID: id},
		Name:  name,
		Label: label,
	}

	result := db.Model(&role).Select("name", "label").Updates(&role)
	if result.Error != nil {
		return nil, result.Error
	}

	result = db.Delete(&RoleResource{}, "role_id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, rid := range rs {
		result := db.Create(&RoleResource{RoleId: role.ID, ResourceId: rid})
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return &role, nil
}

func DeleteRole(db *gorm.DB, rid int) error {
	result := db.Delete(&RoleResource{}, "role_id = ?", rid)
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(&RoleMenu{}, "role_id = ?", rid)
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(&Role{}, rid)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CheckRoleAuth(db *gorm.DB, rid int, uri, method string) (bool, error) {
	resources, err := GetResourcesByRole(db, rid)
	if err != nil {
		return false, err
	}

	for _, r := range resources {
		if r.Anonymous || (r.Url == uri && r.Method == method) {
			return true, nil
		}
	}

	return false, nil
}
