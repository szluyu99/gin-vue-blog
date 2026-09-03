# 代码审查记录

针对 `gin-blog-server` 的一次全量审查（handle / model / middleware / utils / global），
以及针对 `gin-blog-admin` 的一次全量审查（views / components / store / utils / layout）。
每条给出文件与行号、问题原因、当前状态。行号以审查时的代码为准，后续改动可能偏移。

状态说明：`待处理` / `已修复` / `暂缓`（明确决定先不做）/ `潜在`（当前代码路径打不到）。

当前进度：server 功能 BUG F1–F13 全部已修复；server 安全 S1–S8
全部暂缓（项目初期，安全要求不高，上线前必须回来处理）；潜在问题 P1 已修复。
admin 的 A1–A14 已修复，A15–A21（P2）待处理（A21 经测试确认当前无实际影响，已加回归护栏）。
front 的 FE1–FE4 已修复。

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

### F9 文章导入接口 — 已修复

`internal/handle/handle_article.go:376`

```go
fileName := fileHeader.Filename
title := fileName[:len(fileName)-3]
```

- 文件名短于 3 字节（如 `ab`）时切片下界为负，直接 panic → 500；
  它还假设扩展名恰好 3 个字符，`.markdown` 会被砍掉一截
- 该接口对文件类型与大小**完全没有校验**（`/upload` 是有的），任何文件都会被整体读入内存当作正文
- `:385` 传的是 `auth.ID`（`user_auth` 主键），而正常新建走的是 `auth.UserInfoId`（`:101`），
  两个 id 空间不一致时作者归属错误
- 分类与标签硬编码成 `"学习"` / `"Golang"`，且 `ImportArticle` 里用 `FirstOrCreate`，
  等于每次导入都可能凭空建出这两条分类标签

修复：

- 标题改用 `strings.TrimSuffix(name, path.Ext(name))`，去掉后为空则返回 `ErrRequest`
- 新增 `allowedImportExt`（`.md` / `.markdown`）与 `maxImportSize = 5 << 20`，
  并新增业务码 `g.ErrImportType = 9104`（`ErrFileType` 的文案是"只支持上传图片"，不适用于导入）
- 作者改用 `auth.UserInfoId`
- **分类标签方案：不建。** `model.ImportArticle` 去掉这两个参数，导入后是草稿、
  分类标签留空，由用户在后台编辑时补。相应地它现在返回创建出来的 `*Article`，
  handler 也把文章返回给前端

回归测试：`TestArticleImport`（故意让 `auth.ID` 与 `UserInfoId` 错开，断言作者用的是后者，
并断言没有凭空创建分类标签）、`TestArticleImportFileName`（`a.md` / `ab.md` /
`文章.markdown` / `带.点.的.名字.md`）、`TestArticleImportRejectsBadFile`（非法后缀、
只有扩展名、超过大小上限，且都不写库）、`model.TestImportArticle`。

### F10 注册流程：验证 token 先删、建用户无事务 — 已修复

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

修复：

- `VerifyCode` 里把 `DeleteMailInfo` 挪到 `CreateNewUser` 成功之后。
  删 token 本身失败最多导致链接可重复点，而重复点会被下面的查重挡住，比现在直接废掉链接好
- `CreateNewUser` 整体包进 `db.Transaction`
- 昵称改用插入后拿到的 `userinfo.ID`（先 insert 再 update 昵称和简介），并发不会重名
- 事务内加一次 `username = ?` 的存在性检查，避免同一邮箱的两封未过期邮件都点导致建出两个账号；
  命中时返回新的 sentinel error `model.ErrUsernameTaken`
- `BcryptHash` 的 error 不再忽略

回归测试：`model.TestCreateNewUser`（断言昵称按 id 生成）、
`model.TestCreateNewUserRejectsDuplicate`、`model.TestCreateNewUserRollsBackOnFailure`
（删掉 `user_auth_role` 表制造第三步失败，断言前两张表都回滚）、
`TestAuthVerifyCodeKeepsTokenWhenCreateFails`（断言建用户失败后验证链接仍然有效、库里没有半个用户）。

### F13 `UpdateUserInfo` 会把零值写库 — 已修复

`internal/model/user.go` 原来是 `Select("nickname","avatar","intro","website").Updates(userInfo)`。
显式 `Select` 会让 GORM 把零值一起写进去，所以调用方漏传一个字段就等于把它清空。
这是 FE1 那次数据丢失的后端一半：前台表单没同步时只改昵称提交，头像/简介/网站全被清成空串。

修复：只把非空字段放进 map 更新，全空时直接返回不发 SQL。

测试：`internal/model/user_update_test.go` 的 `TestUpdateUserInfoKeepsEmptyFields`。
接口级验证（`PUT /front/user/info`）：先写满四个字段，再只传 nickname、其余留空，
返回 `code 0` 且 `GET` 读回来头像/简介/网站都还在。

## 组件测试

2026-09-03 两个前端都接入了 `@vue/test-utils@2.5.0`。此前只有 store 和 utils 有测试，
组件内部的状态问题（表单快照没同步、模板 ref 拿错、render 函数空指针）全靠读代码推断，
浏览器行为也无法自测（playwright 在本机起不来，见 `project_env_docker` 记忆）。

已覆盖的回归点（每条都验证过：把修复回退后测试会红）：

- `gin-blog-admin/src/components/crud/CrudTable.spec.js` — A6 翻页只发一次请求、
  `pagination` 上不再挂 `onChange`、请求失败清空数据
- `gin-blog-admin/src/views/article/list/index.spec.js` — A1 分类为 null 不崩表且展示「无」、
  A3 `updateOrDeleteArticles` 返回 Promise
- `gin-blog-admin/src/views/article/write/index.spec.js` — A5 新建时 `tag_names` 是数组、
  编辑导入草稿（无分类无标签）不崩
- `gin-blog-front/src/views/user/UploadOne.spec.js` — FE1 `props.preview` 后到也能显示、
  FE4 根相对路径
- `gin-blog-front/src/views/user/index.spec.js` — FE1 表单在 `getUserInfo` 后同步、
  FE2 提交发原始相对路径、未登录跳首页
- `gin-blog-admin/src/views/Login.spec.js` — A8 不勾选「记住我」不保存账号密码、勾选则保存、
  登录失败时 loading 复位且不产生未捕获 rejection、账号密码为空直接提示
- `gin-blog-admin/src/views/auth/role/index.spec.js` — A11 新建时预取两个选项、
  `menu_ids`/`resource_ids` 是数组、勾选的权限会一起提交、选项不重复请求
- `gin-blog-admin/src/layout/tags/index.spec.js` — A21 `tabRefs` 顺序与 `tags` 一致
  （目前是一致的，这条作为回归护栏；顺序一旦错位激活标签的滚动定位就会指错）、
  关闭当前标签跳左边、关闭第一个标签跳第二个

写这类测试的两个坑：

- `vi.mock('vue-router', ...)` 不能整个模块替换，`src/router/index.js` 里的 `createRouter`
  会变成 undefined。要 `async importOriginal => ({ ...await importOriginal(), useRoute, useRouter })`
- naive-ui 组件用 `findComponent({ name: 'NDataTable' })` 找不到，要直接传组件引用
  `findComponent(NDataTable)`
- `vi.mock` 的工厂会被提升到文件顶部，工厂里引用的 `vi.fn()` 必须用
  `const { push } = vi.hoisted(() => ({ push: vi.fn() }))`，否则报
  `Cannot access 'push' before initialization`
- `store/modules/tag.js` 用的是 `@/router` 导出的 `router` 实例，不是 `useRouter()`，
  测跳转要 mock `@/router`
- `layout/tags/index.vue` 模板里同时用了 `useRoute()` 和全局属性 `$route`，
  不装 router 插件时后者是 undefined，要 `global.mocks.$route`
- `<script setup>` 的绑定在测试里可以从 `wrapper.vm` 访问，且 ref 已自动解包
  （`wrapper.vm.isRemember = false`，不要写 `.value`）
- mock 的 api 函数要在 `beforeEach` 里 `mockReset()` 再重新设实现，
  否则调用次数会跨用例累计

## 潜在问题（当前打不到）
### F11 注册强依赖邮件，`Captcha.SendEmail` 是死开关 — 已修复

现象：博客前台注册用户提示「发送邮件失败」（`6101`）。两层原因：

1. `config.yml` 的 `Email.From` / `SmtpPass` / `SmtpUser` 都是空串，没有可用的 SMTP 凭据。
2. `Captcha.SendEmail` 这个开关只在 `internal/global/config.go:59` 定义过，**全仓库没有任何地方读它**，
   `handle_auth.go` 的 `Register` 无条件调用 `utils.SendEmail`，所以改成 false 也不会跳过发邮件。

修复：让开关真正生效。`Register` 在查重之后加一个分支，`Captcha.SendEmail` 为 false 时直接
`model.CreateNewUser` 建号并返回成功，不写 Redis 验证 token、不发邮件；为 true 时保持原来的
邮箱验证流程。`config.yml` 与 `config.docker.yml` 都改为 `SendEmail: false`。
前台 `RegisterModal.vue` 的提示从「邮件已发送」改成「注册成功, 请登录」并自动切到登录弹窗，
同时给 `api.register` 补上 `catch`（原来是裸 await，失败会产生 unhandled rejection）。

测试：`TestAuthRegisterWithoutEmail`（注册即建号、Redis 里没有验证 token、重复注册被挡）。
`testenv_test.go` 的 `newTestEnv` 现在会初始化 `g.Conf`，注意它不能直接覆盖 —— `withJWTConf` /
`withUploadConf` 可能已经在同一个测试里设过，覆盖会让 `TestAuthLogin` 拿不到 JWT 密钥。

顺带修正 `config.docker.yml` 的 Email 段键名：原来写的 `IsSSL` / `Secret` / `Nickname` 和
`g.Config.Email` 的字段对不上（结构体要的是 `SmtpPass` / `SmtpUser`），即使填了也绑定不上。

### F12 `AutomaticEnv` 让同名环境变量吃掉整段配置 — 已修复

`internal/global/config.go:93` 的 `v.AutomaticEnv()` 会让 viper 在解析每个 key 时先查环境变量。
如果环境里存在和某个顶层配置段同名的变量，`Get("Email")` 返回的是那个**字符串**而不是 yaml 里的
map，于是整段反序列化成零值。

实测（本机环境里有 `EMAIL=xxx@xxx.com`，Linux 上 git 相关工具常设）：

```
有 EMAIL 环境变量:  Email = {From: Host: Port:0 SmtpPass: SmtpUser:}
env -u EMAIL:      Email = {From: Host:smtp.qq.com Port:465 SmtpPass: SmtpUser:}
```

后果是 SMTP 拨号变成 `dial tcp :0: connect: connection refused`，报错完全指不到原因。
`Server` / `Captcha` 这些没有同名环境变量的段不受影响，所以只有邮件这一处发作。

这不只影响邮件：任何一段配置只要撞上同名环境变量都会被静默吃掉。

修复：去掉 `AutomaticEnv()` 与 `SetEnvKeyReplacer`，改成 `envBindings` 白名单显式 `BindEnv`，
只绑定 `deploy/start/docker-compose.yml` 里实际用到的 8 个 key（`SERVER_PORT`、`MYSQL_HOST`、
`MYSQL_PORT`、`MYSQL_DBNAME`、`MYSQL_USERNAME`、`MYSQL_PASSWORD`、`REDIS_ADDR`、
`REDIS_PASSWORD`）。这样父路径遮蔽消失，compose 的变量名也不用改。

验证（临时探针，跑完已删）：

```
EMAIL=someone@example.com 时  Email = {From: Host:smtp.qq.com Port:465 ...}   不再被吃掉
REDIS_ADDR=1.2.3.4:6379 时    Redis.Addr = "1.2.3.4:6379"                     覆盖仍生效
```

## 潜在问题（当前打不到）
### P1 本地上传的路径穿越防护恒真 — 已修复

`internal/utils/upload/local.go:60`

```go
p := g.GetConfig().Upload.StorePath + "/" + key
if strings.Contains(p, g.GetConfig().Upload.StorePath) {
```

`p` 由拼接而来，必然包含 `StorePath`，`key` 传 `../../config.yml` 也能通过。
目前 `DeleteFile` 全仓库没有调用方，只在 `upload/oss.go:11` 的接口里声明，所以一直是潜在问题。

修复：改成 `filepath.Abs` 规整后判断是否仍在 `StorePath` 之内，前缀比较带上路径分隔符
（避免 `/data/uploaded-evil` 混过去），越界直接返回错误而不是静默跳过。
测试 `internal/utils/upload/local_delete_test.go`：`../outside.txt` 被拒且文件仍在、
`../uploaded-evil/x.jpg` 被拒、正常 key 能删。

注意 `local.go` 的 `UploadFile` 里有个局部变量叫 `filepath`，会在该函数内遮蔽新导入的
`path/filepath` 包，以后在那个函数里用 `filepath.Xxx` 要留意。

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

## gin-blog-front

### FE1 个人中心头像与资料显示为默认值, 只改昵称会清空其他字段 — 已修复

现象：前台登录后在个人中心上传头像，接口返回 `code 0` 且文件正常落盘、URL 直接访问是
`200 image/jpeg`，但页面刷新后头像仍是默认图，四个输入框也是空的。

原因是两层「拷贝一次」：

1. `views/user/index.vue:14-20` 在 setup 阶段用 `reactive({ avatar: userStore.avatar, ... })`
   把 store 的值拷成快照，而 `getUserInfo()` 是在 `onMounted` 里异步取的，此刻 store 还是默认值。
2. `views/user/UploadOne.vue:14` 又 `ref(props.preview)` 拷一次，且没有 `watch` props
   （admin 的同名组件 `:23` 有），所以父组件后来拿到真实头像也传不进来。

上传后能看到正确的 src，是因为 `previewImg.value = responseJSON.data` 是组件内部直接赋值，
那条路径本来是通的；一刷新就回到默认值。

连带的数据丢失路径（已实测）：表单没同步时四个框都是空的，全空提交会被
`UpdateCurrentUserReq.Nickname` 的 `binding:"required"` 挡住（`9001`），但只要只填昵称就提交，
`model.UpdateUserInfo` 用的是 `Select("nickname","avatar","intro","website").Updates(...)`，
显式 `Select` 会把零值一起写进去：

```
PUT /front/user/info {"nickname":"只改昵称","avatar":"","intro":"","website":""} → code 0
DB: (3, '只改昵称', '', '', '')   ← 头像/简介/网站被清空
```

修复：`onMounted` 里 `getUserInfo()` 返回后 `Object.assign(form, ...)` 重新同步；
`UploadOne.vue` 补 `watch(() => props.preview, ...)`。

### FE2 store 存的是绝对头像 URL, 会被写回库 — 已修复

`store/user.js:55` 原来存 `convertImgUrl(data.avatar)`，即 `http://localhost:8765/...`。
个人中心把它当表单初值，点「修改」就把带域名的绝对地址写回数据库——换域名或换环境就全失效。
这是 admin A12 的前台孪生。

修复：store 只存原始相对路径，转成可访问地址留给展示层的 `convertImgUrl`。已确认所有 img src
都已经套了 `convertImgUrl`（`AppHeader.vue:144`、`CommentField.vue:83`、`Comment.vue:194,229`、
`LinkList.vue:29`、`message/index.vue:126`），所以改成相对路径不影响显示。
`avatar` getter 的默认值判断从 `??` 改成 `||`：头像为空串时也要退回默认图，
否则 `convertImgUrl('')` 会给出已失效的 `dummyimage.com` 占位图。

测试：`store/user.spec.js` 原来断言的是「相对路径会被拼上后端地址」，即旧的错误行为，
已改为断言存原始路径，并新增「后端头像为空时退回默认图」。

### FE4 图片地址拼 localhost, 远程访问必然裂图 — 已修复

这是「上传头像不显示」的真正原因，前面 FE1/FE2 只是同一个页面上的另外两个问题。

`convertImgUrl` 原来把 `VITE_BACKEND_URL`（前台）/ `VITE_SERVER_URL`（后台）拼进图片地址，
而这两个变量在 `.env.development` 里写的是 `http://localhost:8765`。从服务器本机以外的浏览器
访问页面时，`localhost` 指的是浏览器所在的机器，图片必然加载失败。

排查时的教训：在服务器上 `curl http://localhost:8765/public/uploaded/xxx.jpg` 得到
`200 image/jpeg`，据此判断链路正常是错的 —— curl 跑在服务器本机，`localhost` 恰好指对了。
故障环境是远程浏览器，验证环境必须一致。

修法不是把 localhost 换成具体 IP（换 IP 或上域名又会失效），而是改成根相对路径：
`deploy/build/web/default.conf.template:20` 的 nginx 本来就有 `location /public/uploaded`
转发到后端，生产环境的设计就是根相对，只有 dev 模式漏了这条代理才被迫拼绝对地址。

- 两个 `utils/index.js` 的 `convertImgUrl` 返回 `/public/uploaded/xxx.jpg`，并归一化前导 `/`；
  `http` 开头的外链仍原样返回
- 两个 `vite.config.js` 各加一条 `/public` 代理，target 沿用原来的后端地址
- `VITE_BACKEND_URL` / `VITE_SERVER_URL` 现在只作为代理目标，不再进入页面里的 URL

验证：经两个 dev server 代理取图都是 `200 image/jpeg 35217`。
测试里四处断言旧行为的地方已改（`front/utils/index.spec.js`、`front/store/app.spec.js` 两条、
`admin/utils/index.spec.js`）。

### FE3 头像尺寸只在 lg 断点生效 — 已修复
`views/user/UploadOne.vue:58` 的 `class="lg:h-[160px] lg:w-[160px]"` 在 1024px 以下不生效，
图片按原始尺寸撑开（父容器只有 `max-w-[300px]`），布局会炸。改成
`h-[160px] w-[160px] object-cover`。注意这一条不是上面那次故障的原因——
报告时是 1920×1080 最大化窗口，`lg:` 是生效的。

## gin-blog-admin
后台管理端的全量审查。A1–A6 是「当前就会坏」的，已修复；A7 之后待处理。

### A1 文章列表分类为 null 导致整表 render 抛异常 — 已修复

`src/views/article/list/index.vue:86` 原为 `h('div', row.category.name || '无')`。
后端 `model.Article.Category` 是指针，`category_id = 0` 时序列化为 `null`。
这条以前打不到，但 F9（导入文章为草稿、不建分类和标签）之后必然产生这种数据，
只要导入过一篇，整张表格的 render 就会抛异常。改为 `row.category?.name`。

同一个根因还有 `src/views/article/write/index.vue:84`：编辑页 `category.name` 和
`tags.map(...)` 都没有守卫，打开一篇导入生成的草稿会直接白屏。已改为
`tags?.map(e => e.name) ?? []` 和 `category?.name ?? ''`。

### A2 后台批量导入文章不带 Authorization，必然 401 — 已修复

`src/views/article/list/index.vue` 的 `<NUpload action="/api/article/import">` 没有带
认证头，而 `/article/import` 挂在 `JWTAuth(true)` 下，非 mock 构建 100% 失败。
参照 `components/UploadOne.vue` 补 `:headers="{ Authorization: \`Bearer ${authStore.token}\` }"`。
注意这里用 `useAuthStore()` 实例而不是解构出 `token`，否则重新登录后拿到的是旧值
（`UploadOne.vue:20` 就是解构的，属于 A13，未处理）。

顺带修了 `afterUpload` 里裸的 `JSON.parse(event.target.response)`：鉴权失败或被网关
拦截时响应不是 JSON，会抛在 naive-ui 的 `@finish` 回调里。现在走 `utils/parseJson`。

### A3 删除文章不提示成功、失败无法捕获 — 已修复

`src/views/article/list/index.vue:214` 的 `updateOrDeleteArticles` 少一个 `return`。
`composables/useCRUD.js:110` 里 `data = await doDelete(...)` 拿到 `undefined`，于是
`data?.code === 0` 不成立，删成功也不弹提示；失败时因为 promise 没被 await，
外层 `try/catch` 兜不住，变成 unhandled rejection，同时 `refresh()` 照跑当成功。

### A4 操作日志详情弹窗 `JSON.parse('')` 抛在 render 里 — 已修复

`src/views/log/operation/index.vue:228,237,147` 三处 `JSON.stringify(JSON.parse(x), null, 2)`。
`request_param` 对 `POST /user/offline/:id`、`DELETE /menu/:id`、`DELETE /resource/:id`
是空串，点开详情即崩。新增 `utils/formatJson`（解析不了就原样返回），三处都换过去。

### A5 写文章页新建时 `tag_names` 为 undefined — 已修复

`src/views/article/write/index.vue:73` 新建分支重置成
`{ status: 1, is_top: false, title: '', type: 1 }`，漏了 `tag_names`；
`:63-65` 的 watcher 紧接着 `newVal.includes(...)` → TypeError。
现在重置对象补齐 `tag_names: []` 和 `category_name: ''`。

### A6 CrudTable 每次翻页发两次请求 — 已修复

`src/components/crud/CrudTable.vue` 里 `pagination.onChange` 调 `handleQuery()`，
模板上的 `@update:page="onPageChange"` 又调一次。对照装好的 naive-ui 源码
`es/data-table/src/use-table-data.mjs:166-181`：`mergedOnUpdatePage` 先
`call(onChange, page)`，再 `doUpdatePage(page)` 触发组件的 `onUpdate:page`，所以两个
回调每次翻页都会各发一次请求；两个请求没有排序或取消，慢的那个能覆盖新的。
`onUpdatePageSize` 没有双绑，不受影响。
修法是删掉 `pagination.onChange`，只留 `@update:page` 一个入口
（`onPageChange` 里有 `props.remote &&` 判断，前端分页模式下不会多发请求）。

### 测试

`src/utils/index.spec.js` 新增 `parseJson` / `formatJson` 的 9 个断言，覆盖 A4 与 A2 的
空串、非法 JSON、HTML 响应场景。A1/A3/A5/A6 都在组件内部，要验证需要
`@vue/test-utils`（当前未引入），暂时只靠 `pnpm build` + 人工点检。

### 待处理

P1（特定路径下会坏）：

- **A7（已修复）** `src/utils/http.js:59-67`：`code === 1201` 与 `1202/1203/1207` 两个分支直接
  `return`，等于 `Promise.resolve(undefined)`，导致 `CrudTable.vue:78` 的
  `const { data } = await ...`、`Login.vue:55` 的 `resp.data.token`、
  `permission.js:24` 的 `buildRoutes(resp.data)` 二次抛错。应 `return Promise.reject(responseData)`。
  1201 分支还没清 token。
- **A8（已修复）** `src/views/Login.vue:39,60`：`const isRemember = useStorage('isRemember', false)` 是
  Ref，恒为真，取消勾选也会把用户名密码写进 localStorage，`removeLocal` 分支永远走不到。
  少一个 `.value`。（`local.js:47` 只是 base64，不是加密。）
- **A9（已修复）** `src/store/modules/auth.js:36-41` + `layout/header/components/UserAvatar.vue:32`：
  `logout` 既不 await 也不 catch，`/logout` 失败时 token 留在 localStorage、不跳转、不提示。
- **A10（已修复）** `src/layout/index.vue:18-22`：`computed(() => router.getRoutes()...)` 没有响应式
  依赖，永久缓存首次结果；登录后动态添加的路由不在 `<KeepAlive :include>` 里，刷新才生效。
- **A11（已修复）** `src/views/auth/role/index.vue`：菜单 / 资源权限树原来只在
  `modalAction === 'edit'` 下渲染，option 预取还被注释着，F6 之后后端 `model.SaveRole`
  已支持新建时一并写入，只剩前端在逼用户走「先建角色再编辑权限」两趟。
  现在新建入口改为 `handleAddRole()`：先并发拉取两个选项（已拉过则跳过），再打开弹窗，
  弹窗在 `modalAction === 'add'` 时同时给出两棵树；编辑流程保持原样（按入口按钮只展示对应那棵）。
  `initForm` 也补齐成 `{ name: '', label: '', menu_ids: [], resource_ids: [] }`，
  避免 `NTree` 的 `checked-keys` 拿到 undefined。
  测试 `src/views/auth/role/index.spec.js` 4 条；接口级验证：新建角色带
  `menu_ids:[1,2] / resource_ids:[1]`，`role_menu` 与 `role_resource` 均正确写入，
  删除角色后关联无残留。
- **A12（已修复）** `src/views/profile/index.vue:14,22-27,33`：`infoForm.avatar = userStore.avatar` 是
  跑过 `convertImgUrl` 的展示地址，再 `api.updateCurrent` 存回库。FE4 之后它不再是带域名的
  绝对地址（改成了根相对路径），但仍会多出前导 `/`，头像为空时还会把占位图
  `http://dummyimage.com/400x400` 写进库。应发原始 `userInfo.avatar`。
  另外 F13 之后后端已不再把空值写库，所以这条的破坏面比原来小了。
- **A13（已修复）** `src/components/UploadOne.vue:20,26-28`：`const { token } = useAuthStore()` 解构后
  失去响应性；`JSON.parse(respStr)` 无保护（可用新加的 `parseJson`）。
- **A14（已修复）** `src/components/common/ScrollX.vue:33,51`：用了废弃的 `e.wheelDelta`，Firefox 下
  为 undefined，`translateX` 变 `NaN`，横向滚动失效。应换 `e.deltaY`。

P2（展示问题 / 潜在 / 清理）：

- **A15** `src/views/auth/role/index.vue:93-97`：`value: row.is_disable` 配
  `checkedValue: 1, uncheckedValue: 0`，后端字段是 `bool`，开关永远显示关闭。
  目前只是展示，`onUpdateValue` 只弹「这个功能暂时还不支持~」。与 S4 一起处理。
- **A16** `src/views/message/comment/index.vue`：`:81` 列 key 还是老的 `reply_nick_name`
  （render 函数正常，只有接上 `handleExport` 才会导出空列，`CrudTable.vue:136` 按
  `item[key]` 取值）；`:132` `commentTypeMap[row.type].tag` 没守卫；`:63-78`「评论类型」
  和 `:124-136`「来源」是同一字段渲染两遍且前者 `key: ''`；`:275` typo `filterablec`；
  `handleUpdateReview` 无 try/catch。
- **A17** `.then()` 无 `.catch`：`views/article/list/index.vue:45-46`、
  `views/article/write/index.vue:39,42`、`views/user/list/index.vue:40`。
  拦截器已弹过错误提示，只剩控制台的 unhandled rejection 噪音。
- **A18** `src/assets/config.js` 里导出的 `config` 对象（原作者的 QQ / 微博 APP_ID、
  腾讯验证码 ID）全仓库无引用，可删。**注意同文件的 `loginTypeMap` / `articleTypeMap` /
  `commentTypeMap` 及对应 Options 被 4 个页面引用，不是死代码**（和 front 那份不同）。
- **A19** `views/article/list/index.vue:302` 自己实现了一份 `downloadFile`，
  `utils/index.js:40` 已有同名导出，重复。
- **A20** `views/article/list/index.vue:281` 的 `beforeUpload` 只放行 `.md`，
  后端 F9 之后同时接受 `.md` 和 `.markdown`，两边不一致（偏严，不影响正确性）。
- **A21** `src/layout/tags/index.vue:16,42` 的 `v-for` 模板 ref 数组顺序没有保证，
  `tabRefs.value[activeIndex]` 可能取错元素。后果只是激活标签滚动位置偶尔不对。

### admin 排查过但不是问题的

- `src/store/modules/tag.js:66-77` `removeTag` 的负数索引：进入该分支的前提是
  `path === this.activeTag`，而 `path` 必然来自 `tags` 中的某一项，所以 `activeIndex >= 0`。
  sessionStorage 只持久化 `tags`（`pick: ['tags']`），刷新后 `activeTag === ''`，此时
  `if` 不成立也不会崩。`ContextMenu.vue:30` 还额外用 `tags.length <= 1` 禁用了「关闭」。
- 后台评论列表不显示点赞 / 回复数，不依赖 `CommentVO.LikeCount`，F8 与 admin 无关。
- `Paginate` 的 100 条上限在 admin 里打不到，`CrudTable` 的 `pageSizes` 是 `[5, 10, 20]`。
- `api.deleteArticle(ids)` 收到的是 `JSON.stringify(array)` 字符串而 `softDeleteArticle`
  收到的是数组，两边不对称但都能用：`handle_article.go:159` 的 `ShouldBindJSON` 不看
  Content-Type，`"[1,2]"` 作为请求体也能正确解析成 `[]int`。
