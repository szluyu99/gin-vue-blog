package g

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func confWith(mode, secret, salt string) *Config {
	conf := &Config{}
	conf.Server.Mode = mode
	conf.JWT.Secret = secret
	conf.Session.Salt = salt
	return conf
}

// release 模式下空密钥或仓库里的示例密钥必须拒绝启动:
// 以前完全不校验, 部署的人不改也能跑起来, 等于 token 和 session 都可以被伪造
func TestCheckSecretsRejectsWeakSecretsInRelease(t *testing.T) {
	cases := map[string]*Config{
		"密钥为空":   confWith("release", "", "a-real-salt"),
		"盐为空":    confWith("release", "a-real-secret", ""),
		"密钥是示例值": confWith("release", sampleJWTSecret, "a-real-salt"),
		"盐是示例值":  confWith("release", "a-real-secret", sampleSessionSalt),
	}

	for name, conf := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Panics(t, func() { CheckSecrets(conf) })
		})
	}
}

// release 模式下配了真实密钥就正常通过
func TestCheckSecretsPassesWithRealSecrets(t *testing.T) {
	assert.NotPanics(t, func() {
		CheckSecrets(confWith("release", "a-real-secret", "a-real-salt"))
	})
}

// debug 模式只告警不打断: 本地开发用仓库自带的配置文件也要能跑起来
func TestCheckSecretsOnlyWarnsInDebug(t *testing.T) {
	assert.NotPanics(t, func() {
		CheckSecrets(confWith("debug", sampleJWTSecret, sampleSessionSalt))
	})
	assert.NotPanics(t, func() {
		CheckSecrets(confWith("debug", "", ""))
	})
}
