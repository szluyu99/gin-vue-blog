package test

import (
	"testing"
)

func TestCategoryGetList(t *testing.T) {
	BaseSuccessTest(t, Get("/api/category/list", backRouter))
}

func TestCategoryGetFrontList(t *testing.T) {
	BaseSuccessTest(t, Get("/api/front/category/list", frontRouter))
}

// func TestCategoryDelete(t *testing.T) {
// 	ids := []int{11, 23}
// 	BaseHttpSuccessTest(t, DeleteJson("/api/category", ids, backRouter))
// }
// func TestCategorySaveOrUpdate(t *testing.T) {
// 	param := map[string]any{"id": 1, "name": "前端"}
// 	BaseSuccessTest(t, PostJson("/api/category", param, backRouter))
// }
