package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenAndParseToken(t *testing.T) {
	secret := "secret"
	issuer := "issuer"
	expire := 10

	token, err := GenToken(secret, issuer, expire, 1, []int{1, 2})
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	mc, err := ParseToken(secret, token)
	assert.Nil(t, err)
	assert.Equal(t, 1, mc.UserId)
	assert.Len(t, mc.RoleIds, 2)
}

func TestParseTokenError(t *testing.T) {
	tokenString := "tokenString"

	_, err := ParseToken("secret", tokenString)
	assert.ErrorIs(t, err, ErrTokenMalformed)
}
