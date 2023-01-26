package utils

import (
	"crypto/tls"
	"fmt"
	"gin-blog/config"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

// Email 发送方法
// to: 以英文逗号 , 分隔的字符串 如 "a@qq.com,b@qq.com"
func Email(to, subject, body string) error {
	return send(strings.Split(to, ","), subject, body)
}

// 邮件发送
// to: 目标数组
// subject: 邮件标题
// body: 邮件内容 (HTML)
func send(to []string, subject string, body string) error {
	// 从配置文件中读取
	from := config.Cfg.Email.From
	nickname := config.Cfg.Email.Nickname
	secret := config.Cfg.Email.Secret
	host := config.Cfg.Email.Host
	port := config.Cfg.Email.Port
	isSSL := config.Cfg.Email.IsSSL

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", nickname, from) // 发件人
	e.To = to                                       // 收件人
	e.Subject = subject                             // 主题
	e.HTML = []byte(body)                           // 内容

	var err error
	auth := smtp.PlainAuth("", from, secret, host)
	addr := fmt.Sprintf("%s:%d", host, port)

	if isSSL {
		err = e.SendWithTLS(addr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(addr, auth)
	}

	return err
}
