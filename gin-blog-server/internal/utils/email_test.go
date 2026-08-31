package utils

import (
	g "gin-blog/internal/global"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	assert.Equal(t, "abc@qq.com", Format("ABC@qq.com"))
	assert.Equal(t, "abc@qq.com", Format("  abc@qq.com  "))
	assert.Equal(t, "abc@qq.com", Format(" AbC@Qq.Com\n"))
	assert.Equal(t, "", Format("   "))
}

func TestGetCode(t *testing.T) {
	code := GetCode()
	assert.Len(t, code, 24)
	// 随机生成, 两次不应相同
	assert.NotEqual(t, code, GetCode())
}

func TestEncodeDecode(t *testing.T) {
	cases := []string{"", "abc@qq.com", "带中文的内容|123456", "a|b|c"}
	for _, c := range cases {
		got, err := Decode(Encode(c))
		assert.Nil(t, err)
		assert.Equal(t, c, got)
	}

	// 非法 base64
	_, err := Decode("这不是 base64")
	assert.NotNil(t, err)
}

func TestEmailVerificationInfo(t *testing.T) {
	// 正常流程: 生成的信息能解析出邮箱和密码
	info := GenEmailVerificationInfo("abc@qq.com", "123456")
	email, password, err := ParseEmailVerificationInfo(info)
	assert.Nil(t, err)
	assert.Equal(t, "abc@qq.com", email)
	assert.Equal(t, "123456", password)

	// 密码中含分隔符时, 会被切成多段, 属于已知限制
	info = GenEmailVerificationInfo("abc@qq.com", "12|34")
	_, _, err = ParseEmailVerificationInfo(info)
	assert.NotNil(t, err)

	// 被篡改: 不是 base64
	_, _, err = ParseEmailVerificationInfo("not-base64-@@@")
	assert.NotNil(t, err)

	// 被篡改: 是 base64 但段数不对
	_, _, err = ParseEmailVerificationInfo(Encode("abc@qq.com|123456"))
	assert.NotNil(t, err)
}

func TestGetEmailData(t *testing.T) {
	// GetEmailVerifyURL 依赖全局配置中的端口
	g.Conf = &g.Config{}
	g.Conf.Server.Port = ":8765"
	defer func() { g.Conf = nil }()

	info := GenEmailVerificationInfo("abc@qq.com", "123456")
	data := GetEmailData("abc@qq.com", info)

	assert.Equal(t, "abc@qq.com", data.UserName)
	assert.NotEmpty(t, data.Subject)
	assert.True(t, strings.Contains(string(data.URL), "/api/email/verify?info="+info))
	// 端口以 : 开头时会补上 localhost
	assert.True(t, strings.HasPrefix(string(data.URL), "localhost:8765/"))
}
