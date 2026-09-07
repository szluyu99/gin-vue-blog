package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
Redis 计数的落库备份

点赞数 / 浏览数 / 站点访问量 / 访客地域这些计数只在 Redis 里累加, 业务表里没有
对应字段 —— Redis 就是它们的唯一数据源。清一次 Redis、换台机器部署、Redis 被淘汰,
所有文章的点赞和阅读量就归零了, 而这两个数字还是前台列表展示和推荐排序的依据。

这张表只做备份, 不参与任何查询路径: 读写照旧走 Redis, 这里定期存一份、启动时回填。
故意不给 article 表加 like_count / view_count 列 —— 那样会出现两个数据源,
每个查询都得决定信哪一个; 而备份表只有 flush / restore 两个入口碰它。

一行 = 一个计数。Key 是 Redis 键名, Member 对 ZSet/Hash 是成员名, 对普通计数是空串。
*/
type CounterSnapshot struct {
	Key    string `gorm:"type:varchar(50);not null;uniqueIndex:idx_counter_key_member;comment:Redis 键名" json:"key"`
	Member string `gorm:"type:varchar(50);not null;uniqueIndex:idx_counter_key_member;comment:ZSet/Hash 的成员, 普通计数为空" json:"member"`
	Count  int64  `gorm:"not null;comment:计数值" json:"count"`
}

// SaveCounterSnapshots 覆盖式写入: 同一个 (key, member) 只保留最新值
func SaveCounterSnapshots(db *gorm.DB, items []CounterSnapshot) error {
	if len(items) == 0 {
		return nil
	}

	/*
		按 (key, member) upsert, 而不是先清表再插:
		清表 + 插入之间进程被杀就等于把备份弄丢了, 那正是这张表要防的事。
		分批是因为 SQLite 有 999 个绑定变量的上限, 一行三个字段。
	*/
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}, {Name: "member"}},
		DoUpdates: clause.AssignmentColumns([]string{"count"}),
	}).CreateInBatches(items, 200).Error
}

// GetCounterSnapshots 取出某个 key 下的所有计数
func GetCounterSnapshots(db *gorm.DB, key string) ([]CounterSnapshot, error) {
	var list []CounterSnapshot
	err := db.Where("`key` = ?", key).Find(&list).Error
	return list, err
}
