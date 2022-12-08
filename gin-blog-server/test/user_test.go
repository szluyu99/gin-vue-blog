package test

import (
	"testing"
)

func TestUserGetList(t *testing.T) {
	BaseSuccessTest(t, Get("/api/user/list", backRouter))
}

func TestUserUpdatePassword(t *testing.T) {
	param := map[string]any{"username": "username", "password": "password"}
	BaseHttpSuccessTest(t, PutJson("/api/user/password", param, backRouter))
}

// 需要 token: 无法测试
// func TestUserUpdateAdminPassword(t *testing.T) {}
