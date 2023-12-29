package utils

import (
	"crypto/tls"
	"fmt"
	g "gin-blog/internal/global"
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
	from := g.GetConfig().Email.From
	nickname := g.GetConfig().Email.Nickname
	secret := g.GetConfig().Email.Secret
	host := g.GetConfig().Email.Host
	port := g.GetConfig().Email.Port
	isSSL := g.GetConfig().Email.IsSSL

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
