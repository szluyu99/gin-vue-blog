package test

import (
	"gin-blog/utils"
	"testing"
)

func TestSendMial(t *testing.T) {
	err := utils.Email("370166800@qq.com", "邮件标题", "邮件内容: <h1> 123456 </h1>")

	if err != nil {
		t.Log("邮件发送失败: ", err)
	} else {
		t.Log("邮件发送成功")
	}
}
