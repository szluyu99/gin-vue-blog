package model

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 错误来源
const (
	ERR_SOURCE_FRONT = "front" // 博客前台
	ERR_SOURCE_ADMIN = "admin" // 博客后台
)

/*
前端错误日志

前端异常原来只进浏览器 console, 博主不在场就等于没发生。这张表让 window.onerror
和 unhandledrejection 有个落点, 不引入 Sentry 这类外部服务。

上报接口必须允许匿名访问(访客没有登录态), 所以这张表的写入方是不可信的:
- 所有字符串字段都截断, 否则一次请求就能塞进几 MB
- 相同错误按 fingerprint 聚合成一行只累加 count, 一个死循环里的报错不会刷出十万行
- 频率限制在 handler 层做(按 IP), 总量上限由 TrimErrorLogs 兜底
*/
type ErrorLog struct {
	Model

	Source      string `gorm:"type:varchar(20);index:idx_error_source;comment:来源(front/admin)" json:"source"`
	Message     string `gorm:"type:varchar(500);not null;comment:错误信息" json:"message"`
	Stack       string `gorm:"type:varchar(2000);comment:调用栈" json:"stack"`
	Url         string `gorm:"type:varchar(500);comment:出错页面" json:"url"`
	UserAgent   string `gorm:"type:varchar(300);comment:浏览器 UA" json:"user_agent"`
	IpAddress   string `gorm:"type:varchar(50);comment:上报 IP" json:"ip_address"`
	UserId      int    `gorm:"comment:登录用户 id, 匿名为 0" json:"user_id"`
	Fingerprint string `gorm:"type:varchar(32);uniqueIndex;comment:同类错误的聚合键" json:"fingerprint"`
	Count       int    `gorm:"default:1;comment:出现次数" json:"count"`
}

// 各字段的长度上限, 与 gorm tag 里的 varchar 对齐
const (
	errMessageLimit   = 500
	errStackLimit     = 2000
	errUrlLimit       = 500
	errUserAgentLimit = 300
)

// 按 rune 截断, 按字节切会把中文切成半个字符
func truncate(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit])
}

/*
聚合键: 来源 + 错误信息 + 调用栈第一行

只取栈顶一行: 同一个 bug 在不同调用路径下整段栈会不同, 但栈顶通常一致,
用整段栈做键会导致同一个问题散成很多行。
*/
func errorFingerprint(source, message, stack string) string {
	topFrame := stack
	if idx := strings.IndexByte(stack, '\n'); idx >= 0 {
		topFrame = stack[:idx]
	}

	sum := md5.Sum([]byte(source + "|" + message + "|" + strings.TrimSpace(topFrame)))
	return hex.EncodeToString(sum[:])
}

/*
记录一条前端错误, 同类错误只累加次数

用 OnConflict 而不是 FirstOrCreate: 后者命中已有记录时只查不写, count 会永远是 1,
聚合等于静默失效。也不能先查再更新 —— 上报是并发的, 两个请求中间插进来会撞唯一索引。
*/
func AddErrorLog(db *gorm.DB, log *ErrorLog) error {
	if log.Source != ERR_SOURCE_ADMIN {
		log.Source = ERR_SOURCE_FRONT // 只认这两种, 其余一律归到前台
	}

	log.Message = truncate(log.Message, errMessageLimit)
	log.Stack = truncate(log.Stack, errStackLimit)
	log.Url = truncate(log.Url, errUrlLimit)
	log.UserAgent = truncate(log.UserAgent, errUserAgentLimit)
	log.Fingerprint = errorFingerprint(log.Source, log.Message, log.Stack)
	log.Count = 1

	// 命中已有指纹时只累加 count 并把 updated_at 顶到最新:
	// 首次的 stack / url / ip 保持不动, 第一次出现的上下文通常比后续更有价值
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "fingerprint"}},
		DoUpdates: clause.Assignments(map[string]any{
			"count":      gorm.Expr("count + 1"),
			"updated_at": time.Now(),
		}),
	}).Create(log).Error
}

// 后台列表: 按最近出现时间倒序, keyword 匹配错误信息与页面地址
func GetErrorLogList(db *gorm.DB, page, size int, keyword, source string) (list []ErrorLog, total int64, err error) {
	list = make([]ErrorLog, 0)

	filter := func(d *gorm.DB) *gorm.DB {
		if keyword != "" {
			like := "%" + keyword + "%"
			d = d.Where("message LIKE ? OR url LIKE ?", like, like)
		}
		if source != "" {
			d = d.Where("source = ?", source)
		}
		return d
	}

	if err := db.Model(&ErrorLog{}).Scopes(filter).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := db.Model(&ErrorLog{}).Scopes(filter).
		Order("updated_at DESC").
		Scopes(Paginate(page, size)).
		Find(&list)

	return list, total, result.Error
}

func DeleteErrorLogs(db *gorm.DB, ids []int) (rows int64, err error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := db.Where("id IN ?", ids).Delete(&ErrorLog{})
	return result.RowsAffected, result.Error
}

// 保留的最大行数: 聚合之后不同错误种类不会太多, 这个上限只防异常情况
const errorLogKeepRows = 1000

/*
超出上限时删掉最旧的几行

上报接口是匿名的, 聚合和限频都只能压住"同一个错误"和"同一个 IP",
不断变化的错误信息(比如把随机数拼进 message)仍然能让行数无限涨,
这是最后一道兜底。按 updated_at 排, 留下最近还在发生的。
*/
func TrimErrorLogs(db *gorm.DB) error {
	var total int64
	if err := db.Model(&ErrorLog{}).Count(&total).Error; err != nil {
		return err
	}
	if total <= errorLogKeepRows {
		return nil
	}

	var ids []int
	if err := db.Model(&ErrorLog{}).
		Order("updated_at ASC").
		Limit(int(total-errorLogKeepRows)).
		Pluck("id", &ids).Error; err != nil {
		return err
	}

	return db.Where("id IN ?", ids).Delete(&ErrorLog{}).Error
}
