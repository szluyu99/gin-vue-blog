package r

// 错误码汇总
const (
	OK   = 0
	FAIL = 500

	// code= 9000... 通用错误
	ERROR_REQUEST_PARAM = 9001
	ERROR_REQUEST_PAGE  = 9002
	ERROR_INVALID_PARAM = 9003
	ERROR_DB_OPE        = 9004

	// code 91... 文件相关错误
	ERROR_FILE_UPLOAD  = 9100
	EEROR_FILE_RECEIVE = 9101

	// code= 1000... 用户模块的错误
	ERROR_USER_NAME_USED    = 1001
	ERROR_PASSWORD_WRONG    = 1002
	ERROR_USER_NOT_EXIST    = 1003
	ERROR_USER_NO_RIGHT     = 1009
	ERROR_OLD_PASSWORD      = 1010
	ERROR_EMAIL_SEND        = 1011
	ERROR_EMAIL_HAS_SEND    = 1012
	ERROR_VERIFICATION_CODE = 1013

	// code = 1200.. 鉴权相关错误
	ERROR_TOKEN_NOT_EXIST  = 1201
	ERROR_TOKEN_RUNTIME    = 1202
	ERROR_TOKEN_WRONG      = 1203
	ERROR_TOKEN_TYPE_WRONG = 1204
	ERROR_TOKEN_CREATE     = 1205
	ERROR_PERMI_DENIED     = 1206
	FORCE_OFFLINE          = 1207
	LOGOUT                 = 1208

	// code= 2000... 文章模块的错误
	ERROR_ART_NOT_EXIST = 2001

	// code= 3000... 分类模块的错误
	ERROR_CATE_NAME_USED = 3001
	ERROR_CATE_NOT_EXIST = 3002
	ERROR_CATE_ART_EXIST = 3003

	// code= 4000... 标签模块的错误
	ERROR_TAG_EXIST     = 4001
	ERROR_TAG_NOT_EXIST = 4002
	ERROR_TAG_ART_EXIST = 4003

	// code=5000... 评论模块的错误
	ERROR_COMMENT_NOT_EXIST = 5001

	// code=6000... 权限模块的错误
	ERROR_RESOURCE_NAME_EXIST   = 6001
	ERROR_RESOURCE_NOT_EXIST    = 6002
	ERROR_RESOURCE_USED_BY_ROLE = 6003
	ERROR_RESOURCE_HAS_CHILDREN = 6004
	ERROR_MENU_NAME_EXIST       = 6005
	ERROR_MENU_NOT_EXIST        = 6006
	ERROR_MENU_USED_BY_ROLE     = 6007
	ERROR_MENU_HAS_CHILDREN     = 6008
	ERROR_ROLE_NAME_EXIST       = 60010
	ERROR_ROLE_NOT_EXIST        = 60011

	// code=7000 ... 页面模块的错误
	ERROR_PAGE_NAME_EXIST = 7001
)

var codeMsg = map[int]string{
	OK:   "OK",
	FAIL: "FAIL",

	ERROR_REQUEST_PARAM: "请求参数格式错误",
	ERROR_REQUEST_PAGE:  "分页参数错误",
	ERROR_INVALID_PARAM: "不合法的请求参数",
	ERROR_DB_OPE:        "数据库操作异常",

	EEROR_FILE_RECEIVE: "文件接收失败",
	ERROR_FILE_UPLOAD:  "文件上传失败",

	ERROR_USER_NAME_USED:    "用户名已存在",
	ERROR_USER_NOT_EXIST:    "该用户不存在",
	ERROR_PASSWORD_WRONG:    "密码错误",
	ERROR_USER_NO_RIGHT:     "该用户无权限",
	ERROR_OLD_PASSWORD:      "旧密码不正确",
	ERROR_EMAIL_SEND:        "邮件发送失败",
	ERROR_EMAIL_HAS_SEND:    "已朝该邮箱发送验证码（有效期 15 分钟），请检查回收站",
	ERROR_VERIFICATION_CODE: "验证码错误",

	ERROR_TOKEN_NOT_EXIST:  "TOKEN 不存在，请重新登陆",
	ERROR_TOKEN_RUNTIME:    "TOKEN 已过期，请重新登陆",
	ERROR_TOKEN_WRONG:      "TOKEN 不正确，请重新登陆",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN 格式错误，请重新登陆",
	ERROR_TOKEN_CREATE:     "TOKEN 生成失败",
	ERROR_PERMI_DENIED:     "权限不足",
	FORCE_OFFLINE:          "您已被强制下线",
	LOGOUT:                 "您已退出登录",

	ERROR_ART_NOT_EXIST: "文章不存在",

	ERROR_CATE_NAME_USED: "操作失败，分类名已存在",
	ERROR_CATE_NOT_EXIST: "操作失败，分类不存在",
	ERROR_CATE_ART_EXIST: "删除失败，分类下存在文章",

	ERROR_TAG_EXIST:     "操作失败，标签名已存在",
	ERROR_TAG_NOT_EXIST: "操作失败，标签不存在",
	ERROR_TAG_ART_EXIST: "删除失败，标签下存在文章",

	ERROR_COMMENT_NOT_EXIST: "评论不存在",

	ERROR_RESOURCE_NAME_EXIST:   "该资源名已经存在",
	ERROR_RESOURCE_NOT_EXIST:    "该资源不存在",
	ERROR_RESOURCE_HAS_CHILDREN: "该资源下存在子资源，无法删除",
	ERROR_RESOURCE_USED_BY_ROLE: "该资源正在被角色使用，无法删除",
	ERROR_MENU_NAME_EXIST:       "该菜单名已经存在",
	ERROR_MENU_NOT_EXIST:        "该菜单不存在",
	ERROR_MENU_USED_BY_ROLE:     "该菜单正在被角色使用，无法删除",
	ERROR_MENU_HAS_CHILDREN:     "该菜单下存在子菜单，无法删除",
	ERROR_ROLE_NAME_EXIST:       "该角色名已经存在",
	ERROR_ROLE_NOT_EXIST:        "该角色不存在",

	ERROR_PAGE_NAME_EXIST: "该页面名称已经存在",
}

func GetMsg(code int) string {
	return codeMsg[code]
}
