# 代码审查记录

针对 `gin-blog-server` 的一次全量审查（handle / model / middleware / utils / global）。
每条给出文件与行号、问题原因、当前状态。行号以审查时的代码为准，后续改动可能偏移。

状态说明：`待处理` / `已修复` / `暂缓`（明确决定先不做）/ `潜在`（当前代码路径打不到）。

## 安全

按项目当前阶段（初期，安全要求不高）整体**暂缓**，但风险是真实的，上线前必须回来处理。

### S1 JWT 密钥与 Session 盐是仓库里的固定值，启动不校验 — 暂缓

`config.yml:8` `Secret: "abc123321"`、`config.yml:26` `Salt: "salt"`，`config.docker.yml` 相同。
`internal/global/config.go` 的 `ReadConfig` 没有任何校验。

签名链路：`internal/utils/jwt/jwt.go:36` `token.SignedString([]byte(secret))`，
验签链路：`internal/middleware/auth.go:73` `jwt.ParseToken(g.Conf.JWT.Secret, parts[1])`。

任何人都能用仓库里的密钥给 `user_id: 1` 签一个 token，`PermissionCheck` 读到 `IsSuper` 直接放行。
Session cookie 同理，`internal/handle/base.go` 的 `CurrentUserAuth` 信任 session uid，
伪造 cookie 即可通过所有 `/api/front/*`（走的是 `JWTAuth(false)`）。

修法：启动时若这两个值为空或等于仓库默认值，直接 fatal。

### S2 X-Real-IP 可伪造，登录锁定被绕过 — 暂缓

`cmd/main.go:37` `r.SetTrustedProxies([]string{"*"})`，
`internal/utils/ip.go:26` 直接 `c.Request.Header.Get("X-Real-IP")`，不判断来源是否可信。

`internal/handle/handle_auth.go:88` 的失败计数键是
`g.LOGIN_FAIL + utils.MD5(req.Username+"|"+clientIP(c))`，
轮换这个请求头即可对同一账号无限次尝试密码。同时污染 `user_auth.ip_address` 与访客地域统计。

### S3 CORS 放行所有来源且允许携带凭证 — 暂缓

`internal/middleware/base.go:54` `AllowOriginFunc` 恒 `return true` + `AllowCredentials: true`。
上方注释说明了为何不能用 `AllowOrigins: ["*"]`（那个判断是对的：`*` 与凭证请求冲突），
但"全放行 + 回显 Origin"在安全上等价。

`WithCookieStore`（`internal/middleware/base.go:62`）的 `sessions.Options` 只设了 `Path` 与 `MaxAge`，
没有 `SameSite` / `Secure` / `HttpOnly`。任意第三方站点可带访客凭证调用评论、留言、上传、改资料。

修法：来源白名单走配置。

### S4 `is_disable` 只写不读，封禁功能是空的 — 暂缓

写入：`internal/handle/handle_user.go:190` → `internal/model/user.go:132`。
全仓库没有任何读取处：`handle_auth.go` 的 `Login` 与 `middleware/auth.go` 都不查。
被禁用的用户照样登录，旧 token 照样有效。角色的 `is_disable`（`handle_role.go:19`）
同样从未进入 `PermissionCheck`。

### S5 操作日志把请求体原文永久落库，含明文密码 — 暂缓

`internal/middleware/operation_log.go:83` `body, _ := io.ReadAll(c.Request.Body)` →
`RequestParam: string(body)`，`:104` 还存了响应体。

`internal/manager.go:68` 把该中间件挂在整个 `auth` 组上，
而 `manager.go:88` 的 `PUT /user/current/password` 就在这个组里，
所以旧密码与新密码明文写入 `operation_log.request_param`（`longtext`）。

修法：按路由或字段名做脱敏白名单。

### S6 注册链接携带明文密码，且该串直接作为 Redis key — 暂缓

`internal/utils/email.go:56` `Encode(email + "|" + password + "|" + code)` 仅 base64（可逆），
`:86` 拼进验证 URL 发邮件，`internal/handle/cache.go:130` 又把这个 blob 当 Redis 键名。
密码可从收件箱、邮件服务商日志、浏览器历史 / Referer、`KEYS *` 还原。

另外 `ParseEmailVerificationInfo`（`email.go:74`）只取 `str[0]`、`str[1]`，
拼进去的随机 `code` 取出来没被使用，只起了让 blob 唯一的作用。

### S7 邮件模块两处 — 暂缓（邮件功能当前不启用）

- `internal/utils/email.go:129` `slog.Info("User:" + User + "  Pass:" + Pass + ...)` 把 SMTP 密码打进日志
- `internal/utils/email.go:160` `d.TLSConfig = &tls.Config{InsecureSkipVerify: true}`，SMTP 凭证可被中间人截获

### S8 登录查用户用 `LIKE` 而非 `=` — 暂缓

`internal/model/user.go:34` `db.Model(&userauth).Where("username LIKE ?", name).First(&userauth)`。
`Login`（`handle_auth.go:97`）与注册查重（`:244`）都走它。

`_` 与 `%` 是 SQL 通配符：提交 `%` 会命中主键最小的用户（`First` 按主键排序，通常是 admin）；
含下划线的合法邮箱可能命中另一个账号，然后用那个账号的 hash 去校验密码。
密码仍需匹配，所以不是直接绕过，但这是认证路径上的错误运算符，也让"邮箱已注册"判断不准。

## 功能 BUG

### F1 后台文章列表零结果时返回假的数据库错误 — 已修复

`internal/handle/handle_article.go:263`

```go
list, total, err := model.GetArticleList(...)
if err != nil || list == nil {
    ReturnError(c, g.ErrDbOp, err) // err 此时是 nil
    return
}
```

`model.GetArticleList`（`internal/model/article.go:115`）中 `var list []Article` 配 `Find(&list)`，
GORM 零行时 slice 保持 `nil`。因此任何筛不到数据的条件（如标题搜索未命中）都会返回
`ErrDbOp` 且 error 为空串，而不是一个空列表。

修复：`GetArticleList` 里把 `list` 初始化成空切片，handler 只判断 `err != nil`。
`TestArticleList` 原来的 `{"?title=不存在", 0}` 用例只断言了 `page.Total`，
而错误响应解码后 `total` 也是 0，所以没抓住这个 BUG，一并补上了业务码断言。
回归测试：`TestArticleListEmptyResultIsNotError`。

### F2 `validate:` tag 全部失效，评论入参没有校验 — 已修复

全仓库搜不到 `binding.Validator` / `validator.New` / `SetTagName` / `RegisterValidation`，
Gin 的 `ShouldBind` 系列只认 `binding` tag。

`internal/handle/handle_front.go:47-53`

```go
Content  string `json:"content" form:"content"`
Type     int    `json:"type" form:"type" validate:"required,min=1,max=3" label:"评论类型"`
```

`type` 可以是 0、99、负数，`content` 可以是空串。
对照 `handle_article.go:28` 的 `AddOrEditArticleReq` 用的是正确的 `binding:"required,min=1,max=3"`。
全仓库 21 处用 `binding`、3 处用 `validate`。

修复：`FAddCommentReq` 的 `Type` 与 `Content` 改用 `binding`，
`handle_user.go` 里注释掉的那段示例代码也一起改了，避免以后被复制。
回归测试：`TestFrontSaveCommentValidatesInput`。

### F3 前台列表被 `Paginate` 静默截到 100 条 — 已修复

`internal/handle/handle_front.go:109,124,140,155` 都传 `size = 1000`：

```go
list, _, err := model.GetTagList(GetDB(c), 1, 1000, "")
```

但 `internal/model/z_base.go` 的 `Paginate` 里有 `case size > 100: size = 100`。
标签、分类、留言、友链超过 100 条的部分在前台直接消失，
且这几个接口不返回 total、前端也没有分页。

修复：新增 `model.PageSizeAll = -1`，`Paginate` 遇到它就不加 `Offset/Limit`，
四个前台 handler 改传这个常量。同时给 `PageQuery.Size` 加了 `binding:"omitempty,min=1"`，
防止请求参数直接传 -1 变成全表查询。
回归测试：`model.TestPaginate`、`TestFrontSimpleListsReturnMoreThan100`。

### F4 改密码忽略 bcrypt 错误 — 已修复

`internal/handle/handle_user.go:225` `hashPassword, _ := utils.BcryptHash(req.NewPassword)`。
哈希失败会把空串写进库，之后 `BcryptCheck` 永不匹配，账号被永久锁死。

修复：改为接收并返回 error。没有加测试——`bcrypt.GenerateFromPassword` 只在密码超过 72 字节
或随机源出错时才失败，在 handler 层不好稳定构造。

### F5 删除评论不清 Redis 计数、不删子回复 — 已修复

`internal/handle/handle_comment.go:35` 只做了 `GetDB(c).Delete(model.Comment{}, "id in ?", ids)`。

文章删除路径是特意清理过 Redis 的（`handle_article.go:160` 调 `cleanCommentCounters`，
注释写明"新行复用 id 会继承旧计数"），这里没跟上：
`COMMENT_LIKE_COUNT` 与 `comment_user_like:*` 会残留；
`parent_id` 指向已删评论的回复仍留在库里，且仍会被 `GetCommentReplyList` 查出来。

修复：新增 `model.DeleteComments`，在事务里先把回复的 id 查出来一起删，
返回全部被删 id 交给 handler 调 `cleanCommentCounters`。
回归测试：`TestCommentDeleteCleansRepliesAndRedis`（同时验证别的评论及其回复没被误删）。

### F6 新建角色时丢弃资源与菜单关联 — 已修复

`internal/handle/handle_role.go:100`

```go
if req.ID == 0 {
    err := model.SaveRole(db, req.Name, req.Label) // ResourceIds / MenuIds 被丢掉
} else {
    err := model.UpdateRole(db, req.ID, req.Name, req.Label, req.IsDisable, req.ResourceIds, req.MenuIds)
}
```

请求体里带了 `resource_ids` / `menu_ids`，Swagger 注释也写着"同时维护角色的资源与菜单关联"。

修复：`SaveRole` 增加这两个参数，并把关联写入抽成 `replaceRoleRelations`
供新增和更新共用（顺便把原来循环里逐条 `Create` 改成批量插入）。
回归测试：`model.TestSaveRoleWithResourceAndMenu`、`TestRoleCreateKeepsRelations`。

### F7 菜单 / 资源删除的占用检查 fail-open — 已修复

`internal/handle/handle_menu.go:114,133`、`internal/handle/handle_resource.go:174,193`

```go
use, _ := model.CheckMenuInUse(db, menuId)
if use { ... }
```

查询报错时 `use` 为 false，于是删掉了仍被角色引用的菜单 / 资源。

修复：四处都改成接收 error 并返回 `ErrDbOp`。没有加测试——要让这个 `Count` 查询失败
得注入一个坏的 `*gorm.DB`，测试脚手架目前没有这个能力，成本大于收益。

### F8 前台评论列表里回复的点赞数恒为 0 — 已修复

`internal/handle/handle_front.go:279-289` 只用顶层评论的 id 去 `hashCounts`，
并且只给 `data[i].LikeCount` 赋值，嵌套的 `ReplyList` 从未填充。
而 `GetReplyListByCommentId`（`:334`）是填的，所以同一条回复在"展开更多"之后
点赞数会从 0 突然变成真实值。

修复：把截断后的回复 id 一起收进 `hashCounts` 的入参，回填时给 `ReplyList` 也赋值。
回归测试：`TestFrontCommentListFillsReplyLikeCount`。

### F9 文章导入接口 — 待处理

`internal/handle/handle_article.go:376`

```go
fileName := fileHeader.Filename
title := fileName[:len(fileName)-3]
```

- 文件名短于 3 字节（如 `ab`）时切片下界为负，直接 panic → 500
- 该接口对文件类型与大小**完全没有校验**（`/upload` 是有的），任何文件都会被整体读入内存当作正文
- `:385` 传的是 `auth.ID`（`user_auth` 主键），而正常新建走的是 `auth.UserInfoId`（`:101`），
  两个 id 空间不一致时作者归属错误
- 分类与标签硬编码成 `"学习"` / `"Golang"`

### F10 注册流程：验证 token 先删、建用户无事务 — 待处理

`internal/handle/handle_auth.go:303`

```go
DeleteMailInfo(GetRDB(c), code)
...
_, _, _, err = model.CreateNewUser(GetDB(c), username, password)
```

一次性 token 先被删掉，`CreateNewUser` 失败则该链接作废，用户只能重新注册。

`internal/model/auth.go:371` 的 `CreateNewUser` 连做三次 insert
（`user_info`、`user_auth`、`user_auth_role`）**没有事务**，
第 2 或第 3 步失败会留下孤儿 `user_info`，或一个没有角色的用户。
昵称按 `Count` 生成"游客N"，并发注册会重名，且那次 `Count` 的 error 只打了日志。

## 潜在问题（当前打不到）

### P1 本地上传的路径穿越防护恒真 — 潜在

`internal/utils/upload/local.go:60`

```go
p := g.GetConfig().Upload.StorePath + "/" + key
if strings.Contains(p, g.GetConfig().Upload.StorePath) {
```

`p` 由拼接而来，必然包含 `StorePath`，`key` 传 `../../config.yml` 也能通过。
目前 `DeleteFile` 全仓库没有调用方，只在 `upload/oss.go:11` 的接口里声明，所以是潜在问题。

## 排查过但不是问题的

记录下来避免重复排查。

- `internal/utils/ip.go:135` 的 `ipSource[2]`：曾怀疑 `GetIpSource` 返回 `""` 时会越界 panic。
  实测（xdb 缺失、非法 IP 串）都在 `ipSource[0] != "中国" && ipSource[0] != "0"` 这一步提前返回，
  不会走到下标 2，**不会 panic**。
- 带 `Group(...)` 的链上调用 `Count(&total)`（`model/tag.go:36`、`model/category.go:33`、
  `model/article.go:141`）是正确的：GORM 的 `Count` 在检测到 `GROUP BY` 时会回退成行数。
- `...Count(&total).…Find(&list)` 这种链式复用是安全的：GORM 每次 `Execute` 结束会
  `stmt.SQL.Reset()`，第二个 finisher 会重建 SQL。
- `db.Updates(category)` 传值而非指针仍然会带上主键 `WHERE`（GORM 在 `!CanAddr()` 时补主键条件），
  `SaveOrUpdateCategory` / `SaveOrUpdateTag` / `UpdateUserInfo` 不会变成全表更新。
- 评论与留言内容在入库前已做 HTML 转义（`handle_front.go:185`、`:228` 的
  `template.HTMLEscapeString`），前台用 `v-html` 渲染不构成存储型 XSS。

## 顺手记下的清理项（不是 BUG）

- `internal/model/auth.go` 的 `SaveOrUpdateRole` 只被 `auth_extra_test.go:137` 调用，
  生产代码里没有调用方。要删得连测试一起改，价值不高，先留着。
- `internal/handle/handle_front.go` 的 `GetTagList` / `GetCategoryList` / `GetMessageList` /
  `GetLinkList` 都不返回 total，前端也没有分页。数据量大了之后一次全量返回会变慢，
  届时要么加分页要么加缓存。
