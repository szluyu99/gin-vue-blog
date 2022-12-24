package config

var Cfg Config

// 配置文件的结构体
type Config struct {
	Server  Server
	JWT     JWT
	Mysql   Mysql
	Redis   Redis
	Session Session
	Upload  Upload
	Zap     Zap
	Qiniu   Qiniu
}

type Server struct {
	AppMode   string
	BackPort  string
	FrontPort string
}

type JWT struct {
	Secret string // JWT 签名
	Expire int64  // 过期时间
	Issuer string // 签发者
}

type Mysql struct {
	Host     string // 服务器地址
	Port     string // 端口
	Config   string // 高级配置
	Dbname   string // 数据库名
	Username string // 数据库用户名
	Password string // 数据库密码
	LogMode  string // 日志级别
}

type Redis struct {
	DB       int    // 指定 Redis 数据库
	Addr     string // 服务器地址:端口
	Password string // 密码
}

type Session struct {
	Name   string
	Salt   string
	MaxAge int
}

type Upload struct {
	// Size int    // 文件上传的最大值
	OssType     string // 存储类型
	Path        string // 本地文件访问路径
	StorePath   string // 本地文件存储路径
	MdPath      string // Markdown 访问路径
	MdStorePath string // Markdown 存储路径
}

type Zap struct {
	Level        string // 级别
	Prefix       string // 日志前缀
	Format       string // 输出
	Directory    string // 日志文件夹
	MaxAge       int    // 日志留存时间
	ShowLine     bool   // 显示行
	LogInConsole bool   // 输出控制台
}

type Qiniu struct {
	ImgPath       string // 外链链接
	Zone          string // 存储区域
	Bucket        string // 空间名称
	AccessKey     string // 秘钥AK
	SecretKey     string // 秘钥SK
	UseHTTPS      bool   // 是否使用https
	UseCdnDomains bool   // 上传是否使用 CDN 上传加速
}
