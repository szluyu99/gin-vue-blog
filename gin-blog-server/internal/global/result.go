package g

import (
	"fmt"
)

const (
	SUCCESS = 0   // 成功业务码
	FAIL    = 500 // 失败业务码
)

// 自定义业务 error 类型
type Result struct {
	code int
	msg  string
}

func (e Result) Code() int {
	return e.code
}

func (e Result) Msg() string {
	return e.msg
}

var (
	_codes    = map[int]struct{}{}   // 注册过的错误码集合, 防止重复
	_messages = make(map[int]string) // 根据错误码获取错误信息
)

// 注册一个响应码, 不允许重复注册
func RegisterResult(code int, msg string) Result {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	if msg == "" {
		panic("错误码消息不能为空")
	}

	_codes[code] = struct{}{}
	_messages[code] = msg

	return Result{
		code: code,
		msg:  msg,
	}
}

// 根据响应码获取响应信息
func GetMsg(code int) string {
	return _messages[code]
}

var (
	OkResult   = RegisterResult(SUCCESS, "OK")
	FailResult = RegisterResult(FAIL, "FAIL")
)

var (
	ErrRequest  = RegisterResult(9001, "请求参数格式错误")
	ErrDbOp     = RegisterResult(9004, "数据库操作异常")
	ErrRedisOp  = RegisterResult(9005, "Redis 操作异常")
	ErrUserAuth = RegisterResult(9006, "用户认证异常")

	ErrPassword     = RegisterResult(1002, "密码错误")
	ErrUserNotExist = RegisterResult(1003, "该用户不存在")
	ErrOldPassword  = RegisterResult(1010, "旧密码不正确")

	ErrTokenNotExist    = RegisterResult(1201, "TOKEN 不存在，请重新登陆")
	ErrTokenRuntime     = RegisterResult(1202, "TOKEN 已过期，请重新登陆")
	ErrTokenWrong       = RegisterResult(1203, "TOKEN 不正确，请重新登陆")
	ErrTokenType        = RegisterResult(1204, "TOKEN 格式错误，请重新登陆")
	ErrTokenCreate      = RegisterResult(1205, "TOKEN 生成失败")
	ErrPermission       = RegisterResult(1206, "权限不足")
	ErrForceOffline     = RegisterResult(1207, "您已被强制下线")
	ErrForceOfflineSelf = RegisterResult(1208, "不能强制下线自己")

	ErrFileUpload  = RegisterResult(9100, "文件上传失败")
	ErrFileReceive = RegisterResult(9101, "文件接收失败")

	ErrTagHasArt  = RegisterResult(4003, "删除失败，标签下存在文章")
	ErrCateHasArt = RegisterResult(3003, "删除失败，分类下存在文章")

	ErrResourceNotExist    = RegisterResult(6002, "该资源不存在")
	ErrResourceUsedByRole  = RegisterResult(6003, "该资源正在被角色使用，无法删除")
	ErrResourceHasChildren = RegisterResult(6004, "该资源下存在子资源，无法删除")
	ErrMenuNotExist        = RegisterResult(6006, "该菜单不存在")
	ErrMenuUsedByRole      = RegisterResult(6007, "该菜单正在被角色使用，无法删除")
	ErrMenuHasChildren     = RegisterResult(6008, "该菜单下存在子菜单，无法删除")

	ErrSendEmail = RegisterResult(6101, "发送邮件失败")
	ErrCodeNoexit = RegisterResult(6102, "Code不存在 请重新注册")
	ErrParseEmailCode = RegisterResult(6103, "解析邮件Code失败 请重试")
	ErrUserExist = RegisterResult(6104, "该邮箱已经注册 请重新注册")
)
