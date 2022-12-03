package test

import (
	"testing"
)

func TestUserRegister(t *testing.T) {
	param := map[string]any{"username": "username", "password": "password"}
	BaseHttpSuccessTest(t, PostJson("/api/register", param, router))
}

func TestUserLogin(t *testing.T) {
	param := map[string]any{"username": "yusael", "password": "123456"}
	BaseSuccessTest(t, PostJson("/api/login", param, router))
}
