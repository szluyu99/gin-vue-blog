# 功能规划

从产品角度对 gin-vue-blog 的现状做的一次梳理，2026-09-04。
和 `code_audit.md` 的分工：那份记录「已有代码哪里不对」，这份记录「接下来做什么、为什么、以及为什么不做某些事」。

状态标记：`待做` / `进行中` / `已完成` / `不做`（明确决定不做，附理由）。

## 一、断头路

访客或管理员能点进去、但走不通的路。优先级高于任何新功能。

### P1 分类 / 标签下的文章列表没有分页 — 已完成

`handle_front.go` 的 `GetArticleList` 把 model 返回的 total 丢掉了：

```go
list, _, err := model.GetBlogArticleList(...)  // total 被丢弃
ReturnSuccess(c, list)
```

前台 `views/article/list` 也没传页码。结果：某个分类下有 30 篇文章时，访客只能看到第一页，**剩下的没有任何入口能进去**。归档页有分页，这里没有，是漏做。

已改为：

- `GetArticleList` 返回 `PageResult[model.Article]`（`page_data` / `total` / `page_num` / `page_size`），
  和归档、评论列表等其它分页接口一致
- 前台 `views/article/list` 每页 9 条，加 `NPagination`，切页滚回顶部；只有一页时不渲染分页器；
  列表为空时给出「这里还没有文章」提示（以前是一片空白）
- 首页的无限加载改读 `page_data`；前台 mock 的 `/front/article/list` 同步返回分页结构

测试：后端 `TestFrontGetArticleList` 改断言分页结构，新增 `TestFrontGetArticleListPaging`
（total 是过滤后的总数、末页条数、越界页码返回空列表）；前台 `article/list` 补 4 条
（按分页参数请求、切页重新请求并滚动、单页不渲染分页器、空列表提示）。
端到端：`page_size=2` 翻到第 2 页、越界页码、按分类过滤时 total 为过滤后的数量都验过。

### P2 后台「登录日志」是占位页 — 已完成

`generate-data` 里登记了这个菜单，但没有 `login_log` 表和接口，页面只有一个「该页面还在开发中」。

选择了做掉——`user_auth` 里已经有 `ip_address` / `ip_source` / `last_login_time`，只是没有历史记录；做完能和已有的「登录失败次数限制」形成完整一条线。

已实现：

- `model.LoginLog`：用户名、昵称、IP、归属地、状态（1-成功 2-失败）、失败原因。
  **失败的记录同样落库**——只有成功记录看不出撞库
- `Login` 的五条路径各写一条：成功、密码错误、用户不存在、账号被禁用、失败次数过多。
  写日志失败只告警，不影响登录本身
- 失败时拿不到昵称（用户可能都不存在），所以 `nickname` 允许为空，前端显示 `-`；
  密码错误时用户是存在的，`user_id` 会记下来，用户不存在时是 0
- 接口 `GET /login/log/list`（关键字匹配用户名/昵称/IP）+ `DELETE /login/log`，
  在 `seed_resource.go` 登记为「登录日志模块」，权限对账会自动补上
- 后台 `views/log/login` 从占位页换成真实列表：状态标签、失败原因、精确到秒的登录时间、
  行内删除与批量删除；admin mock 同步补了数据和路由

测试：`handle_loginlog_test.go`（列表、三种关键字命中、批量删除、请求体非数组）、
`TestAuthLoginWritesLoginLog`（成功/密码错误/用户不存在三条记录的字段）、
`TestAuthLoginLogForDisabledUser`；前端 `views/log/login/index.spec.js` 6 条。
端到端：四种登录路径各造一条记录，字段与关键字过滤都验过。

### P3 前台「说说 / 相册」是占位页 — 说说已完成，相册待做

- **说说** ✅：`talk` 表 + 前台 `/talks` 列表 + `/talk/:id` 详情 + 后台「说说管理」（公开/私密、置顶、评论），首页 `TalkingCarousel` 从一言改成了轮播真实说说。第一版只做纯文本 + 评论，图片支持留到后面（要一起解决存储与缩略图）。
- **相册** 优先级低：本质是图片管理，而上传目前只支持本地目录和七牛云，图片一多就要面对存储与缩略图，收益/成本不如说说。

## 二、值得做的新功能

### 后台首页做成仪表盘 — 已完成

`gin-blog-admin/src/views/home/index.vue` 重做：待处理事项（待审核评论 / 留言）、最新文章、分类分布、最近登录四块，数据用 `Promise.allSettled` 并行拉，单个接口失败不拖垮整页。

没引图表库：分类分布用横向条形（div 宽度百分比）画，够看且省一个依赖。访问趋势没做——`ARTICLE_VIEW_COUNT` 是 Redis ZSet 只有累计值，要趋势得先加按天的 key（如 `view:2026-09-04`），留着。

### 站内通知 — 已完成

`notification` 表 + 前台头部铃铛 + `/notification` 独立页 + 移动端侧栏入口。触发点：评论被回复、自己的文章被评论。开了评论审核时不在发评论那刻通知（等于把没过审的内容推给对方），而是在后台点「通过」时补发（`ReviewComments` 只返回真正 false→true 的那些）。

关联信息只存 id，昵称和文章标题查询时 join 出来，这样对方改昵称、文章改标题之后通知里显示的会跟着变；只有评论内容截断存一份，因为评论可能被删但摘要还要看得到。

点通知会跳到 `/article/:id?comment=:commentId`，前台自动翻页定位到那条评论并高亮 4 秒。

**后续可关注**：定位是前端一页页加载找过去的（每页还有 0.8s 的展示延时），目标评论在第 8 页时要连着请求 8 次，体感会慢。评论量真上来了再考虑后端加个「按评论 id 算它在第几页」的接口一次跳到位；当前评论量下不值得加这个复杂度。

### 阅读体验小补丁 — 已完成

- 阅读进度条：顶部细线，随滚动位置变化
- 代码块一键复制：渲染完成后往每个 `pre` 里注入按钮（`src/utils/code-block.js`）
- 老文章提示：超过 90 天的文章顶部提示内容可能已过时

### 本地样例内容 — 已完成

空库能验证「服务起得来」，但验证不了分页、归档分组、评论回复这些必须有内容才看得出问题的地方。`internal/model/seed_demo.go` 灌一批可控的样例内容（3 分类 / 6 标签 / 15 篇文章 / 评论与回复 / 留言 / 友链），入口是 `generate-data -t demo`，或 `./dev.sh fresh --demo`。只在库里一篇文章都没有时才灌。

### 前台界面统一 — 已完成

分三轮把前台各页面过了一遍, 沉淀出几条约定, 后续加页面照着来:

- 颜色只用 `uno.config.js` 里的语义 token: `primary`(#49b1f5) 主色、`accent`(#ff7242) 强调/激活态、
  `surface`/`surface-soft` 底色、`main`/`muted` 文字、`divider` 分隔。不要再写死 hex ——
  硬编码在暗色模式下不跟随 CSS 变量
- `<img>` 必须同时写宽高并加 `object-cover`: 只写宽度时图片加载失败高度会塌成 0,
  整行布局跟着散(友链头像、最新文章缩略图、首页封面都踩过); 再给容器一个 `bg-surface-soft` 占位
- 卡片手感统一: 常态 `bg-surface-soft` + `shadow-sm`, 悬停上浮 + `shadow-md`,
  不要用「常态透明、悬停整块变色」那种跳变
- 列表入场动画用 `card-enter`(见 `styles/animate.css`), 三个参数在那里调; 翻页时不重复错峰
- 改完 `uno.config.js` 必须重启 dev server, UnoCSS 的 vite 插件不热更新 config;
  新增语义色只用在状态分支里时要留静态兜底, 否则规则没生成会静默变成不可见元素
- 外链一律 `target="_blank"` + `rel="noopener noreferrer"`
- 横幅已经有页面标题时, 卡片里不要再重复一遍, 改成说明数量

## 三、技术架构

结论：**当前架构对这个规模是合适的，不要提前优化。**

### 缓存加 TTL 并主动失效 — 已完成

`page` / `config` 是「数据库为真源、Redis 只是副本」的读穿缓存，原来过期时间写的是 `0`（永不过期），改了库只能手动 `redis-cli del` 再重启后端 —— 换说说封面时实际踩到过。

- `internal/handle/cache.go` 加 `cacheTTL = 10 * time.Minute`，`addPageCache` / `addConfigCache` 都带上
- `config` 是 Hash，`HMSet` 不像 `Set` 那样能顺手带过期时间，要单独 `Expire`；两条命令放进 `TxPipelined`，否则中间失败会留下一个永不过期的 key，又退回原样
- 补上漏掉的主动失效：`UpdateAbout` 只写库没清缓存，而 about 就是 `config` 表里的一行、会被 `GetConfigMap` 一起缓存，导致前台 `/config` 一直返回旧的「关于我」

**故意没给点赞数 / 浏览数 / 访客地域加 TTL**：那几个键只在 Redis 里累加，数据库没有对应字段，Redis 就是唯一数据源，过期等于丢数据。TTL 只对「有真源可回落」的缓存成立，这条区分比 TTL 本身更重要。

测试：`TestCacheHasTTL`（string 与 hash 两种类型都断言有 TTL，并用 miniredis 快进到过期后确认回落成未命中）、`TestUpdateAboutClearsConfigCache`。两条都验证过「去掉修复就会失败」。

### 可以考虑

- **CI 自动发布镜像到 GHCR**：现在只构建验证、不发布，加上之后部署可以直接 pull 而不用在服务器上 build
- **极简前端错误上报**：现在前端异常只进浏览器 console。`window.onerror` + `unhandledrejection` 收集后 POST 到后端一张表，后台加个列表页，不引 Sentry 也能解决问题

### 不做

- **ElasticSearch 搜索**：数据库 LIKE 在几百篇文章的量级完全够用，引 ES 等于多一个要运维的中间件
- **SSR / SSG**：博客确实适合，但这是把整个前台重写，且前后台以后要拆库，现在动是双倍工作
- **微服务 / 分库分表**：规模完全不适用
- **RSS / sitemap / robots.txt**：学习项目，暫时不需要
- **第三方登录（QQ / 微信 / GitHub）**：暫时不需要
- **国际化**：体力活为主，学习收益低
- **后端日志切割**：目前日志量不需要
- **后台的 `@iconify/vue` 换 UnoCSS 图标类**：前台已经换完（图标名都是静态字面量，改完运行时请求从 20 个降到 0）。后台不换：菜单图标名存在数据库的 `menu.icon` 字段里，编译期看不到，必须用 UnoCSS `safelist` 预生成 —— 要覆盖 `assets/icons.js` 里 IconPicker 提供的 228 个 `mdi-*` 加数据库现有的 24 个，其中 24 个还用到 7 个当前没加载的图标集（carbon / cib / el / icon-park-outline / iconoir / ph / tabler），dev server 启动要多解析约 15MB JSON，CSS 增加 30~40KB gzip。后台是自用工具，接受首屏图标晚 1~2 秒出现（不阻塞内容）

## 推荐顺序

1. 文章列表分页（半天）— 修断头路 ✅
2. 说说功能（1–2 天）— 完整练一遍全栈流程，填掉一个占位页 ✅
3. 登录日志（半天）— 填掉第二个占位页 ✅
4. 后台仪表盘（1–2 天）— 让后台有中枢 ✅
5. 站内通知（1–2 天）— 让评论区能对话 ✅
6. 之后再看缓存 TTL、镜像发布、错误上报
