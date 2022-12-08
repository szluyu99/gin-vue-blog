package test

import (
	"testing"
)

func TestCommentGetList(t *testing.T) {
	BaseSuccessTest(t, Get("/api/comment/list", backRouter))
}

func TestCommentDelete(t *testing.T) {
	ids := []int{11}
	BaseSuccessTest(t, DeleteJson("/api/comment", ids, backRouter))
}

func TestUpdateCommentReview(t *testing.T) {
	param := map[string]any{"ids": []int{11, 12}, "is_review": 1}
	BaseSuccessTest(t, PutJson("/api/comment/review", param, backRouter))
}

// func TestCommentSave(t *testing.T) {
// 	param := map[string]any{"topic_id": 1, "content": "go test 新增评论", "type": 1}
// 	BaseSuccessTest(t, PostJson("/api/comment", param, router))
// }
