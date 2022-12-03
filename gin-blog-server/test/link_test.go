package test

import (
	"testing"
)

func TestFriendLinkGetList(t *testing.T) {
	BaseSuccessTest(t, Get("/api/link/list", router))
}

func TestFriendLinkDelete(t *testing.T) {
	ids := []int{11}
	BaseSuccessTest(t, DeleteJson("/api/link", ids, router))
}

func TestSaveOrUpdateFriendLink(t *testing.T) {
	param := map[string]any{
		"name":    "go-test",
		"avatar":  "avatar_url",
		"address": "blog_address",
		"intro":   "blog_intro",
	}
	BaseSuccessTest(t, PostJson("/api/link", param, router))
}
