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

没引图表库：分类分布用横向条形（div 宽度百分比）画，够看且省一个依赖。

### 访问趋势图 — 已完成

做的时候先发现了一个更严重的问题：**`POST /report` 全站没人调用**。前台只有错误上报，admin 里那行 `api.report()` 是注释掉的，所以 `view_count`、`unique_visitor`、`visitor_area` 三个键一直是空的——仪表盘的「访问量」和访客地域统计从来就是 0，不是趋势缺数据，是整条链路没接上。

- 前台 `App.vue` 挂载时调一次 `reportVisit()`（`utils/visit-report.js`）。用 sessionStorage 去重而不是 localStorage：关掉标签页再进来算新的一次访问，正好是趋势图要的粒度；同一会话内刷新、切路由都不重复计数。隐私模式下 sessionStorage 会抛错，catch 掉之后照样上报——统计不该挡住页面
- 后端按天记一个 `view_count:2006-01-02`，`Incr` + 30 天 TTL。两条命令放进 `TxPipelined`：`Incr` 建出来的 key 不带过期时间，单独 `Expire` 中间失败就留下一个永不过期的键（和 config 缓存踩过的同一个坑）
- **按天的不做访客去重**，这是和累计值刻意不同的口径：累计值统计「来过多少人」（同一访客只算一次），趋势要看「每天来了多少次」。放进去重分支里就只会记新访客，老访客再来趋势图上完全看不到。两个数不相等是预期，不是 bug
- `GetHomeInfo` 用一次 `MGet` 取回 14 天（不是逐天 Get，14 个来回变 1 个），**没有数据的那天补 0 而不是跳过**，否则前端横轴会缺格、趋势看着是错的
- 图还是不引图表库：`div` 高度百分比画柱状，按最大值归一化。0 的那天留 2% 的底，不然横轴上少一格像是日期断了

端到端：模拟上报后 `view_count:2026-09-07` = 3、TTL 2592000；补两天历史数据后 `/api/home` 返回 14 天升序且中间补 0；无头浏览器打开前台确认只发一次 `POST /api/report`、会话内二次导航不重复；后台仪表盘 14 根柱子、汇总「近 14 天共 36 次」、0 的那天柱高 2%、无控制台报错。

### 访客地域分布 — 已完成

`visitor_area`（Report 里按省份累加的 Hash）之前**没有任何接口读它**——前台一直没调上报，这个键一直是空的，所以后台也就没做展示。上报接通之后它才有意义，顺手补上。

- `GetHomeInfo` 多返回一个 `visitor_area`，不新增接口
- **后端排序而不是前端排**：Hash 的遍历顺序不保证，不定序的话每刷新一次仪表盘这几项就跳一次。人数相同时按名字排，保证同样的数据每次返回同样的顺序
- 截断到 8 项：小卡片里再多挤不下
- 单个值不是数字时跳过并告警，不让一条坏数据把整个仪表盘打回错误
- 取不到 IP 归属的记成「未知」，照常展示——比悄悄丢掉更实在

仪表盘底部那一行从「分类分布 + 最近登录」两列改成三列。测试：`TestGetVisitorArea`（倒序 + 同数按名字 + 坏数据跳过）、`TestGetVisitorAreaLimitAndEmpty`（空返回空、截断保留最多的），admin 3 条。端到端：往 Redis 灌四个省份后 `/api/home` 返回排好序的四项，后台页面渲染出四条并显示「12 人 / 7 人 / 7 人 / 2 人」，无控制台报错。


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

### 前端错误上报 — 已完成

前端未捕获异常原来只进浏览器 console，博主不在场就等于没发生。`error_log` 表 + 匿名上报接口 + 后台「前端错误」页，不引 Sentry。前台和后台都接，用 `source` 字段区分。

**聚合而不是流水**：一行一个 `fingerprint`（来源 + 错误信息 + 栈顶一行），重复出现只累加 `count` 并把 `updated_at` 顶到最新，首次的 stack / url / ip 保持不动。前端错误分布极不均匀，一个死循环几秒能报上千次，流水表会被单个 bug 淹掉；而真正要看的是「有几类问题、哪类最近还在发生」。只取栈顶是因为同一个 bug 在不同调用路径下整段栈会不同，用整段做键会把一个问题散成很多行。

**上报接口必须匿名**（访客没有登录态），等于开了个无鉴权写库入口，所以是四层防护：

- 按 IP 配额（1 分钟 10 条）。**超配额返回 SUCCESS 而不是错误** —— 上报逻辑挂在 `window.onerror` 上，返回错误会让它的错误处理再报一次，形成放大循环
- 前端自去重：同一条错误在单次会话里只发一次，不然死循环会把 IP 配额吃光，真正需要看的错误反而被丢掉
- 字段按 rune 截断（按字节切会把中文切成半个字符）
- `TrimErrorLogs` 总量上限 1000 行：把随机数拼进 message 能同时绕过聚合和限频，这是最后一道兜底

**只报未捕获的**：业务里已经 try/catch 过的是预期内错误，不该占配额。mock 模式（GitHub Pages 演示站）不上报，那里没有后端。上报用 `fetch` 而不是项目里的 axios，绕开响应拦截器 —— 上报失败绝不能再产生用户可见的动静。

顺带修掉一个真 bug：`generateDefaultMenus` 里父菜单 `db.Create` 撞唯一索引时 `parents[i].ID` 会留在 0，往已有库里新增子菜单就会变成孤儿（菜单树不渲染、路由 404）。首次灌库父子一起建所以一直没暴露，加「前端错误」菜单时才踩到。

### 缓存加 TTL 并主动失效 — 已完成

`page` / `config` 是「数据库为真源、Redis 只是副本」的读穿缓存，原来过期时间写的是 `0`（永不过期），改了库只能手动 `redis-cli del` 再重启后端 —— 换说说封面时实际踩到过。

- `internal/handle/cache.go` 加 `cacheTTL = 10 * time.Minute`，`addPageCache` / `addConfigCache` 都带上
- `config` 是 Hash，`HMSet` 不像 `Set` 那样能顺手带过期时间，要单独 `Expire`；两条命令放进 `TxPipelined`，否则中间失败会留下一个永不过期的 key，又退回原样
- 补上漏掉的主动失效：`UpdateAbout` 只写库没清缓存，而 about 就是 `config` 表里的一行、会被 `GetConfigMap` 一起缓存，导致前台 `/config` 一直返回旧的「关于我」

**故意没给点赞数 / 浏览数 / 访客地域加 TTL**：那几个键只在 Redis 里累加，数据库没有对应字段，Redis 就是唯一数据源，过期等于丢数据。TTL 只对「有真源可回落」的缓存成立，这条区分比 TTL 本身更重要。（后来给这几个键加了落库备份，见下一节。）

测试：`TestCacheHasTTL`（string 与 hash 两种类型都断言有 TTL，并用 miniredis 快进到过期后确认回落成未命中）、`TestUpdateAboutClearsConfigCache`。两条都验证过「去掉修复就会失败」。

### Redis 计数落库备份 — 已完成

上一条留下的那个「Redis 是唯一数据源」现在有了兜底。点赞数、浏览数、站点访问量、访客地域清一次 Redis 就归零，而前台列表展示和推荐排序都靠这两个数字。

- 新表 `counter_snapshot`（`key` + `member` + `count`，`(key, member)` 唯一索引）。**故意不给 `article` 表加 `like_count` / `view_count` 列**：那样就有了两个数据源，每个查询都得决定信哪一个；备份表只有 flush / restore 两个入口碰它，查询路径完全不变
- 落库用 upsert 而不是「清表再插」：清表和插入之间进程被杀，等于把备份弄丢了，而那正是这张表要防的事
- 每 10 分钟落一次，**收到 SIGTERM / SIGINT 时再落一次**。正常重启（`docker restart` / `systemctl restart`）走的都是 SIGTERM，不补这一次就会丢掉上个周期之后的所有点赞和浏览。为此 `main.go` 从 `r.Run()` 改成 `http.Server` + `signal.NotifyContext` + `Shutdown`；关闭流程用独立的 context，信号那个已经被取消了，拿它发 Redis 请求会立刻失败
- 回填**只补 Redis 里不存在的键**。Redis 还在时它的值一定比备份新（备份最多是上个周期的快照），覆盖回去等于把这段时间的增量抹掉。判断的是「键在不在」而不是逐个成员比大小——键存在就说明这份数据是活的，一个成员都不动
- 按天的访问量（`view_count:2006-01-02`）不备份：那是趋势图的原料，只留 30 天，丢一天不影响任何业务判断，备份反而要处理「哪些天已经过期」

**踩到的坑**：一开始把三个计数键都当 ZSet 处理，单测因为用同样的假设灌数据所以全绿，端到端时 `ZSCORE article_like_count 1` 直接报 `WRONGTYPE`。实际类型是——Hash：文章点赞数、评论点赞数、访客地域；ZSet：只有文章浏览数；String：站点访问量。测试跟着改成按真实类型灌数据。

顺带修掉：访客地域把内网访问记成了一个名字叫「0」的省份。ip2region 定位不到的字段是字面量 `"0"`，内网地址整条是 `0|0|0|内网IP|内网IP`，而原来直接取第 3 段。现在取不到省份就归「未知」，内网则用城市位上的「内网IP」。

测试：`internal/counter_test.go` 5 条（落库各类型、反复落库是 upsert 不堆行、清空 Redis 后恢复且类型正确、**Redis 有数据时一个字节都不动**、空库空 Redis 不报错）、`TestVisitorProvince`。端到端：造点赞/浏览/访客后 `./dev.sh stop`，日志出现「收到退出信号 → Redis 计数已落库 rows=5 → 已退出」，备份表 5 行；然后 Redis 容器整个重建（真的空了），重启后日志四条「从备份恢复计数」，Redis 里的值和关停前一致。

### CI 发布镜像到 GHCR — 已完成（只做发布）

原来 docker job 用 `load: true` 把镜像装进 runner 本地、跑完冒烟测试就随 runner 销毁。现在 main 有新提交时，会把 `gvb-web` / `gvb-server` 推到 `ghcr.io/szluyu99/`，两个 tag：`latest` 给部署用，`main-<短 sha>` 用于回滚到具体某次提交。

几个刻意的选择：

- **冒烟测试通过之后才推**，不是把 `load: true` 换成 `push: true`。反过来的话，测失败时坏镜像已经在仓库里了
- 用 `if: github.event_name == 'push' && github.ref == 'refs/heads/main'` 圈住。PR（尤其 fork 来的）拿不到 `packages: write`，不加这个条件 PR 一律红
- 镜像名过一遍 `tr 'A-Z' 'a-z'`：owner 带大写字母时 `docker push` 直接拒绝
- 凭证用默认的 `GITHUB_TOKEN` + job 级 `permissions: packages: write`，不额外配 secret

**没有一起改 compose**：`deploy/start/docker-compose.yml` 里 server 和 web 仍是 `build:`，web 的 context 还是 `build_web.sh` 现场拼出来的 `${WEB_BUILD_CONTEXT}`，换成 `image: ghcr.io/...` 等于把本地构建链路一起换掉。等确实会去 pull 了再动。另外首次推完要在 GitHub 上把两个 package 设为 public，否则 pull 需要先 login——这一步只能在网页上点。

### 可以考虑

- **部署侧改成拉 GHCR 镜像**：见上，收益取决于是否真的在别的机器上部署（2026-09-07 评估：暂不做）

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
6. 缓存 TTL、镜像发布、错误上报 ✅
7. 访客上报 + 访问趋势 + 访客地域、Redis 计数落库 ✅

## 当前状态（2026-09-07）

这一轮告一段落，项目进入休息期。下次接手从这一节看起。

**代码状态**：`main` 分支，已推到 `github/main`（远端只剩 `main`、`api-v1` 和几个 dependabot 分支，我开的功能分支都删了）。测试：后端 8 个包全绿，前台 32 文件 197 条，后台 27 文件 165 条，两个前端 lint + build 通过，`docs/` 与 Swagger 注解同步。

**这一轮做完的**（细节各见上文对应小节和 `code_audit.md`）：
说说功能、站内通知（含评论定位高亮）、后台仪表盘、登录日志、前端错误上报、缓存 TTL + 主动失效、S2/S3/S5 三项安全修复、CI 发布镜像到 GHCR、访客上报接线 + 访问趋势 + 访客地域、Redis 计数落库备份，以及移动端顶栏按钮塌成 0×0（手机上等于没有导航）这个真 bug。

**接着做的候选**（按我的推荐顺序）：

1. **前台侧边信息模块整体优化**。窄屏下 `AuthorInfo` / `Announcement` / `WebsiteInfo` 整列消失（`views/home/index.vue:107` 的 `col-span-0 lg:col-span-3`），手机上看不到公告和网站资讯。2026-09-07 决定：不只做一个折叠入口，这个模块本身要重新设计，所以整块留到下次
2. **文章分享元信息**：`router.js:148` 只设了 `document.title`，没有 `meta description` 和 `og:*`，分享到微信 / Twitter 没有摘要和封面。纯前端，1~2 小时
3. **浏览器冒烟测试进 CI**：移动端那个 0×0 按钮单测和 build 都发现不了，只有真实布局能暴露。本机 Playwright 可用（`/root/.cache/ms-playwright`）。2026-09-07 评估为暂不做
4. **相册页**：目前是「禁止访问」占位页。本质是图片管理，上传只支持本地目录和七牛云，图片一多要处理存储与缩略图 —— 一直建议不做

**明确暂缓、下次不用重新评估的**：

- 部署侧改成拉 GHCR 镜像（CI 已在发布，compose 仍走本地 `build:`）
- `code_audit.md` 的 S6 / S7：都挂在邮件注册链路上，当前 `Captcha.SendEmail: false` 没启用
- 评论定位改成后端算页码：当前评论量下前端逐页找足够用

**本地环境要点**：`./dev.sh` 一键起 Redis + 后端 + 前台 + 后台（日志在 `.dev/log/`）；连 github 要先 `export https_proxy=http://agent.baidu.com:8891`；docker daemon 在容器外，改不了它的代理配置。
