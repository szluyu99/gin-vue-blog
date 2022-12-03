package test

import (
	"testing"
)

func TestCasbin(t *testing.T) {
	// * HasRoleForUser
	// res, err := global.G_E.HasRoleForUser("user", "role")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// * DeletePermission =>
	// b, err := global.G_E.DeletePermission("/hello", "*")
	// fmt.Println(b)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// * DeletePermissionsForUser
	// b, err := global.G_E.DeletePermissionForUser("tttt")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(b)

	// * AddPermissionForUser(a, b, c) => p, a, b, c
	// global.G_E.AddPermissionForUser("abcde", "/hello", "GET")

	// * AddRoleForUser(a, b) => g, a, b
	// _, err := global.G_E.AddRoleForUser("hhhh", "test")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
