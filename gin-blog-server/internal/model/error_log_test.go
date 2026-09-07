package model

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
同类错误聚合成一行并累加 count

这是整个设计的核心: 前端错误分布极不均匀, 一个死循环能在几秒内报上千次,
不聚合的话单个 bug 就把表淹了。上一版用 FirstOrCreate 写的, 命中已有记录时
只查不写, count 永远是 1 —— 聚合静默失效, 所以这条断言必须有。
*/
func TestAddErrorLogAggregates(t *testing.T) {
	db := newModelDB(t)

	report := func() error {
		return AddErrorLog(db, &ErrorLog{
			Source:  ERR_SOURCE_FRONT,
			Message: "Cannot read properties of undefined",
			Stack:   "at foo (a.js:1:1)\nat bar (b.js:2:2)",
			Url:     "/article/1",
		})
	}

	assert.Nil(t, report())
	assert.Nil(t, report())
	assert.Nil(t, report())

	var logs []ErrorLog
	assert.Nil(t, db.Find(&logs).Error)
	require.Len(t, logs, 1, "同一个错误只该有一行")
	assert.Equal(t, 3, logs[0].Count, "重复上报要累加 count")
	assert.Equal(t, "/article/1", logs[0].Url, "首次的上下文保持不变")
}

// 栈顶不同就是不同的问题, 要分成两行
func TestAddErrorLogFingerprint(t *testing.T) {
	db := newModelDB(t)

	base := ErrorLog{Source: ERR_SOURCE_FRONT, Message: "boom"}

	first := base
	first.Stack = "at foo (a.js:1:1)\nat x (x.js:9:9)"
	assert.Nil(t, AddErrorLog(db, &first))

	// 栈顶一致、后续帧不同: 视为同一个问题
	same := base
	same.Stack = "at foo (a.js:1:1)\nat y (y.js:8:8)"
	assert.Nil(t, AddErrorLog(db, &same))

	// 栈顶不同: 另一个问题
	other := base
	other.Stack = "at baz (c.js:3:3)"
	assert.Nil(t, AddErrorLog(db, &other))

	// 来源不同: 也算另一个问题
	admin := base
	admin.Source = ERR_SOURCE_ADMIN
	admin.Stack = "at foo (a.js:1:1)"
	assert.Nil(t, AddErrorLog(db, &admin))

	var count int64
	assert.Nil(t, db.Model(&ErrorLog{}).Count(&count).Error)
	assert.Equal(t, int64(3), count, "栈顶相同的合并, 栈顶或来源不同的分开")
}

// 上报接口匿名可访问, 字段长度不能信前端
func TestAddErrorLogTruncates(t *testing.T) {
	db := newModelDB(t)

	assert.Nil(t, AddErrorLog(db, &ErrorLog{
		Source:    "伪造的来源",
		Message:   strings.Repeat("中", errMessageLimit+50),
		Stack:     strings.Repeat("s", errStackLimit+50),
		Url:       strings.Repeat("u", errUrlLimit+50),
		UserAgent: strings.Repeat("a", errUserAgentLimit+50),
	}))

	var log ErrorLog
	assert.Nil(t, db.First(&log).Error)
	// 按 rune 截断: 按字节切会把中文切成半个字符
	assert.Equal(t, errMessageLimit, len([]rune(log.Message)))
	assert.Equal(t, errStackLimit, len(log.Stack))
	assert.Equal(t, errUrlLimit, len(log.Url))
	assert.Equal(t, errUserAgentLimit, len(log.UserAgent))
	assert.Equal(t, ERR_SOURCE_FRONT, log.Source, "来源只认 front/admin, 其余归到前台")
}

// 列表: 按最近出现倒序, 支持关键字与来源筛选
func TestGetErrorLogList(t *testing.T) {
	db := newModelDB(t)

	assert.Nil(t, AddErrorLog(db, &ErrorLog{Source: ERR_SOURCE_FRONT, Message: "前台炸了", Url: "/article/1", Stack: "a"}))
	assert.Nil(t, AddErrorLog(db, &ErrorLog{Source: ERR_SOURCE_ADMIN, Message: "后台炸了", Url: "/setting/page", Stack: "b"}))

	list, total, err := GetErrorLogList(db, 1, 10, "", "")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)

	// 按来源
	list, total, err = GetErrorLogList(db, 1, 10, "", ERR_SOURCE_ADMIN)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, "后台炸了", list[0].Message)

	// 关键字同时匹配 message 与 url
	_, total, err = GetErrorLogList(db, 1, 10, "前台", "")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	_, total, err = GetErrorLogList(db, 1, 10, "setting", "")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
}

func TestDeleteErrorLogs(t *testing.T) {
	db := newModelDB(t)

	assert.Nil(t, AddErrorLog(db, &ErrorLog{Message: "boom", Stack: "a"}))
	var log ErrorLog
	assert.Nil(t, db.First(&log).Error)

	// 空数组不发 SQL
	rows, err := DeleteErrorLogs(db, nil)
	assert.Nil(t, err)
	assert.Zero(t, rows)

	rows, err = DeleteErrorLogs(db, []int{log.ID})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rows)
}

/*
总量兜底

聚合压住的是「同一个错误」, 限频压住的是「同一个 IP」。
把随机数拼进 message 的话两道都绕得过去, 行数仍会无限涨, 所以要有总量上限。
*/
func TestTrimErrorLogs(t *testing.T) {
	db := newModelDB(t)

	// 没超上限时什么都不做
	assert.Nil(t, AddErrorLog(db, &ErrorLog{Message: "boom", Stack: "a"}))
	assert.Nil(t, TrimErrorLogs(db))
	var count int64
	assert.Nil(t, db.Model(&ErrorLog{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)

	// 直接造行: 走 AddErrorLog 要凑一千多个不同指纹, 没必要
	logs := make([]ErrorLog, 0, errorLogKeepRows+20)
	for i := 0; i < errorLogKeepRows+20; i++ {
		logs = append(logs, ErrorLog{
			Message:     "msg" + strconv.Itoa(i),
			Fingerprint: "fp" + strconv.Itoa(i),
			Count:       1,
		})
	}
	assert.Nil(t, db.CreateInBatches(logs, 200).Error)

	assert.Nil(t, TrimErrorLogs(db))
	assert.Nil(t, db.Model(&ErrorLog{}).Count(&count).Error)
	assert.Equal(t, int64(errorLogKeepRows), count, "超出的最旧记录要被删掉")
}
