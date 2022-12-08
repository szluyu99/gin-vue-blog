package test

import (
	"testing"
)

func TestTagGetList(t *testing.T) {
	BaseSuccessTest(t, Get("/api/tag/list", backRouter))
}

// func TestTagGetFrontList(t *testing.T) {
// 	resp := BaseSuccessTest(t, Get("/api/tag/list/front", router))

// 	m := resp.Data.(map[string]any)
// 	total := int(m["total"].(float64))

// 	assert.GreaterOrEqual(t, total, 0)
// }

func TestTagDelete(t *testing.T) {
	ids := []int{11, 12}
	BaseSuccessTest(t, DeleteJson("/api/tag", ids, backRouter))
}

func TestTagSaveOrUpdate(t *testing.T) {
	param := map[string]any{"id": 1, "name": "前端"}
	BaseSuccessTest(t, PostJson("/api/tag", param, backRouter))
}
