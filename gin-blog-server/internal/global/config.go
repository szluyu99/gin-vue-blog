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
	"jwt.secret":     "JWT_SECRET",
	"session.salt":   "SESSION_SALT",
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
