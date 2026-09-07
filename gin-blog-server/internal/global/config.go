package g

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Mode          string // debug | release
		Port          string
		DbType        string // mysql | sqlite
		DbAutoMigrate bool   // 是否自动迁移数据库表结构
		DbLogMode     string // silent | error | warn | info
		// 可信代理网段, 只有来自这些地址的请求才认 X-Forwarded-For / X-Real-IP。
		// 留空时用内网网段兜底, 见 TrustedProxies()
		TrustedProxies []string `mapstructure:"trusted-proxies"`
		// 允许跨域的前端来源, 如 https://blog.example.com。留空时只放行本机与内网, 见 AllowedOrigins()
		AllowedOrigins []string `mapstructure:"allowed-origins"`
	}
	Log struct {
		Level     string // debug | info | warn | error
		Prefix    string
		Format    string // text | json
		Directory string
	}
	JWT struct {
		Secret string
		Expire int64 // hour
		Issuer string
	}
	Mysql struct {
		Host     string // 服务器地址
		Port     string // 端口
		Config   string // 高级配置
		Dbname   string // 数据库名
		Username string // 数据库用户名
		Password string // 数据库密码
	}
	SQLite struct {
		Dsn string // Data Source Name
	}
	Redis struct {
		DB       int    // 指定 Redis 数据库
		Addr     string // 服务器地址:端口
		Password string // 密码
	}
	Session struct {
		Name   string
		Salt   string
		MaxAge int
		Secure bool // 部署在 https 之后时置 true, session cookie 只走加密连接
	}
	Email struct {
		From     string // 发件人 要发邮件的邮箱
		Host     string // 服务器地址, 例如 smtp.qq.com 前往要发邮件的邮箱查看其 smtp 协议
		Port     int    // 前往要发邮件的邮箱查看其 smtp 协议端口, 大多为 465
		SmtpPass string // 邮箱密钥 不是密码是开启smtp后给你的密钥
		SmtpUser string // 邮箱账号
	}
	Captcha struct {
		SendEmail  bool // 是否通过邮箱发送验证码
		ExpireTime int  // 过期时间
	}
	Upload struct {
		// Size      int    // 文件上传的最大值
		OssType   string // local | qiniu
		Path      string // 本地文件访问路径
		StorePath string // 本地文件存储路径
	}
	Qiniu struct {
		ImgPath       string // 外链链接
		Zone          string // 存储区域
		Bucket        string // 空间名称
		AccessKey     string // 秘钥AK
		SecretKey     string // 秘钥SK
		UseHTTPS      bool   // 是否使用https
		UseCdnDomains bool   // 上传是否使用 CDN 上传加速
	}
}

var Conf *Config

func GetConfig() *Config {
	if Conf == nil {
		log.Panic("配置文件未初始化")
		return nil
	}
	return Conf
}

// 允许用环境变量覆盖的配置项 (deploy/start/docker-compose.yml 里用到的那些)
//
// 这里必须显式绑定, 不能用 v.AutomaticEnv():
// AutomaticEnv 下 viper 解析 email.host 时会先检查父路径是否被环境变量遮蔽
// (isPathShadowedInAutoEnv), 只要环境里存在同名的 EMAIL 变量 (Linux 上 git 相关
// 工具常设), 整个 Email 段就会反序列化成零值, 表现为 SMTP 拨号 dial tcp :0。
// 任何一段配置撞上同名环境变量都会被静默吃掉, 所以改成白名单。
var envBindings = map[string]string{
	"server.port":    "SERVER_PORT",
	"server.dbtype":  "SERVER_DBTYPE",
	"jwt.secret":     "JWT_SECRET",
	"session.salt":   "SESSION_SALT",
	"sqlite.dsn":     "SQLITE_DSN",
	"mysql.host":     "MYSQL_HOST",
	"mysql.port":     "MYSQL_PORT",
	"mysql.dbname":   "MYSQL_DBNAME",
	"mysql.username": "MYSQL_USERNAME",
	"mysql.password": "MYSQL_PASSWORD",
	"redis.addr":     "REDIS_ADDR",
	"redis.password": "REDIS_PASSWORD",
}

// 仓库里配置文件自带的示例密钥, 谁都能从 GitHub 上读到,
// 用它签发的 token / session 等于没有签名
const (
	sampleJWTSecret   = "abc123321"
	sampleSessionSalt = "salt"
)

// 从指定路径读取配置文件
func ReadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)

	for key, env := range envBindings {
		if err := v.BindEnv(key, env); err != nil {
			panic("绑定环境变量失败: " + err.Error())
		}
	}

	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败: " + err.Error())
	}

	if err := v.Unmarshal(&Conf); err != nil {
		panic("配置文件反序列化失败: " + err.Error())
	}

	CheckSecrets(Conf)

	log.Println("配置文件内容加载成功: ", path)
	return Conf
}

// 校验签名密钥: 以前完全不校验, 空值直接跑起来(任何人都能自签 token),
// 示例值又原样躺在仓库里, 部署的人往往不知道要改
//
// release 模式下空值或示例值直接 panic; debug 模式只告警, 不打断本地开发
func CheckSecrets(conf *Config) {
	release := conf.Server.Mode == "release"

	for _, item := range []struct {
		name, value, sample, env string
	}{
		{"JWT.Secret", conf.JWT.Secret, sampleJWTSecret, "JWT_SECRET"},
		{"Session.Salt", conf.Session.Salt, sampleSessionSalt, "SESSION_SALT"},
	} {
		var reason string
		switch {
		case item.value == "":
			reason = "为空"
		case item.value == item.sample:
			reason = "还是仓库里的示例值"
		default:
			continue
		}

		msg := fmt.Sprintf("%s %s, 请在配置文件中修改, 或用环境变量 %s 注入", item.name, reason, item.env)
		if release {
			panic(msg)
		}
		log.Println("[警告] " + msg)
	}
}

/*
可信代理网段

留空时默认只信任内网(反向代理通常和后端在同一台机器或同一 docker 网络),
这样公网访客伪造 X-Real-IP / X-Forwarded-For 不会被采信。
原来是 SetTrustedProxies("*") —— 任何人都能自称任意 IP。
如果代理在别的网段, 在 config.yml 里写 server.trusted-proxies。
*/
func (*Config) TrustedProxies() []string {
	if Conf != nil && len(Conf.Server.TrustedProxies) > 0 {
		return Conf.Server.TrustedProxies
	}
	return []string{"127.0.0.1/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "::1/128", "fc00::/7"}
}

// 允许跨域的前端来源白名单, 留空时由中间件退回「只放行本机与内网」, 见 middleware.CORS
func (*Config) AllowedOrigins() []string {
	if Conf == nil {
		return nil
	}
	return Conf.Server.AllowedOrigins
}

// session cookie 是否只走 https
func (*Config) SessionSecure() bool {
	return Conf != nil && Conf.Session.Secure
}

// 数据库类型
func (*Config) DbType() string {
	if Conf.Server.DbType == "" {
		Conf.Server.DbType = "sqlite"
	}
	return Conf.Server.DbType
}

// 数据库连接字符串
func (*Config) DbDSN() string {
	switch Conf.Server.DbType {
	case "mysql":
		conf := Conf.Mysql
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?%s",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Config,
		)
	case "sqlite":
		return Conf.SQLite.Dsn
	// 默认使用 sqlite, 并且使用内存数据库
	default:
		Conf.Server.DbType = "sqlite"
		if Conf.SQLite.Dsn == "" {
			Conf.SQLite.Dsn = "file::memory:"
		}
		return Conf.SQLite.Dsn
	}
}
