package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"github.com/k3a/html2text"
	"github.com/thanhpk/randstr"
	"github.com/vanng822/go-premailer/premailer"
	"gopkg.in/gomail.v2"
)

// 注册的核心思想即发送邮件的同时创建code存储在本地redis中，当用户点击验证链接时即向sever发出code数据，如果这个数据在redis中存在，则验证成功，否则失败
type EmailData struct {
	URL  		template.URL // 验证链接
	UserName    string // 用户名即邮箱地址
	Subject     string // 邮件主题
}

// 将邮箱地址转换成小写，并去除空格
// 格式化邮件地址可以防止写错大小写重复注册，同时给用户预留犯错空间，输入空格和大小写错误也能正常处理
func Format(email string) string{
	return strings.ToLower(strings.TrimSpace(email))
}

// 生成随机字符串
func GetCode() string{
	code := randstr.String(24)
	return code
}

// 返回 生成base64编码
func Encode (s string) string{
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return data
}
// 返回 解码base64
func Decode (s string) (string,error){
	data ,err := base64.StdEncoding.DecodeString(s)
	if err != nil{
		return "",errors.New("emailVertify failed, decode error!!!")
	}
	return string(data),nil
}

// 返回生成加密后的base64字符串
func GenEmailVerificationInfo(email string,password string) string{
	code := GetCode()
	info := Encode(email+"|"+password+"|"+code)   // 加密
	return info
}

// 返回解析base64字符串后的 邮箱地址和code
func ParseEmailVerificationInfo(info string) (string,string,error){
	data,err := Decode(info)
	if err!= nil{
		return "","",err
	}

	str := strings.Split(data,"|")
	if len(str) != 3{   // 被篡改了格式
		return "","",errors.New("Wrong Vertifacion Info fomat")
	}

	return str[0],str[1],nil
}

//生成验证链接
func GetEmailVerifyURL(info string) string{
	baseurl := g.GetConfig().Server.Port
	if baseurl[0] ==':'{
		baseurl = fmt.Sprintf("localhost%s",baseurl)
	}
	// 如果是用docker部署,则 注释上面的代码，使用下面的代码
	// baseurl := "你的域名"   切记不需要加端口

	return fmt.Sprintf("%s/api/email/verify?info=%s",baseurl,info)  // form数据
}

//生成邮件数据
func GetEmailData(email string,info string) *EmailData{
	return &EmailData{
		URL: template.URL(GetEmailVerifyURL(info)),// 验证链接
		UserName: email,  // 用户邮箱地址
		Subject: "请完成帐号注册", // 邮件主题
	}
}

//解析模板目录
func ParseTemplateDir(dir string) (*template.Template,error){
	var paths []string
	// 遍历模板目录，将所有模板文件路径添加到paths中
	err := filepath.Walk(dir,func(path string,info os.FileInfo,err error) error{
				if err != nil {
					return err
				}
				if !info.IsDir(){
					paths = append(paths,path)
				}
				return nil
	})

	if err!= nil{
		return nil,err
	}
	return template.ParseFiles(paths...)
}

// 发送邮件
// 发送邮件需要配置邮箱服务器信息， 可以在config.yaml中配置
// 以下情况会发生错误: 1. 邮箱配置错误,smtp信息错误 2. 修改模板后,解析模板失败!
func SendEmail(email string,data *EmailData) error{
	config := g.GetConfig().Email
	from := config.From
	Pass := config.SmtpPass
	User := config.SmtpUser
	to := email
	Host := config.Host
	Port := config.Port
	slog.Info("User:"+User+"  Pass:"+Pass+"  Host:"+Host+"  Port:")

	var body bytes.Buffer
	// 解析模板
	template,err := ParseTemplateDir("../assets/templates")
	if err != nil {
		return errors.New("解析模板失败")
	}
	slog.Info("解析模板成功\n")

	fmt.Println("URL:",data.URL)
	// 执行模板
	template.ExecuteTemplate(&body,"email-verify.tpl",&data) // 把html数据存储在body中， 第二个参数是模板名称， 第三个参数是模板数据（把模板中的占位符换成data数据）
	//为了确保html文件在各个邮件客户端都能正常显示，把html转换成内联模式
    htmlString := body.String()
	prem,_ := premailer.NewPremailerFromString(htmlString,nil)
	htmlline,_ := prem.Transform()
	m :=gomail.NewMessage()   // 使用gomail库发送邮件

	slog.Info("准备发送邮件\n")
	// 设定m头
	m.SetHeader("From",from)
	m.SetHeader("To",to)
	m.SetHeader("Subject",data.Subject)

	// 设定html体 内容
	m.SetBody("text/html",htmlline)	
	m.AddAlternative("text/plain",html2text.HTML2Text(body.String()))

	//配置SMTP连接
	d := gomail.NewDialer(Host,Port,User,Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify:true}
	slog.Info("smtp 连接已建立")
	if err := d.DialAndSend(m); err != nil{
		return err
	}
	return nil
}
