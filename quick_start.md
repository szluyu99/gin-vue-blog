# 快速开始

三个子项目：`gin-blog-server`（后端）、`gin-blog-front`（博客前台）、`gin-blog-admin`（博客后台）。

前台和后台都内置 Mock 模式，只想看前端效果的话不需要启动后端。

| 项目 | 地址 | 说明 |
| --- | --- | --- |
| 博客前台 | http://localhost:8888 | |
| 博客后台 | http://localhost:8889 | 默认账号 `admin / 123456` |
| 后端接口 | http://localhost:8765/api | Swagger: http://localhost:8765/swagger/index.html |

## 方式一：一键启动（联调推荐）

前置依赖：Go 1.26+、pnpm、Redis（没装的话脚本会用 docker 起一个）。数据库默认 SQLite，无需额外安装。

```bash
./dev.sh            # 启动 Redis + 后端 + 前台 + 后台
./dev.sh restart    # 改完代码重启
./dev.sh fresh      # 从零开始: 旧库备份成 gvb.db.bak, 清 Redis, 重新建表灌基础数据
./dev.sh stop
./dev.sh status
./dev.sh logs server   # server | front | admin | redis
./dev.sh seed          # 单独重新灌一次基础数据
```

几个实现上的点：

- 两个前端的 `VITE_USE_MOCK` 由脚本在命令行覆盖为 `false`（vite 里 shell 环境变量优先级高于 `.env` 文件），所以不用改文件就能打到真后端。
- Redis：端口上已经有就直接复用；否则优先本机 `redis-server`，都没有则用 docker 起 `redis:7.0-alpine`，`stop` 时一并清掉。
- 首次启动（检测不到 `gin-blog-server/cmd/gvb.db`）会先执行一次 `generate_data.sh`，它自己会 AutoMigrate 建表，再正式启动。
- `fresh` 用来从零验证一遍流程：把旧库改名成 `gvb.db.bak`（不是删除，测试数据攒久了手一抖就没了）、清掉 Redis 的 DB 7、然后走首次启动。要恢复就把 `.bak` 改回来。上传的图片不动。
- 后端是先 `go build` 再跑二进制，记录的 pid 就是服务进程本身；前端用独立进程组启动，`stop` 整组杀，不会留下孤儿占端口。
- 端口被别的进程占着会直接报错退出，不会让 vite 悄悄换到 8890。
- 日志和 pid 在 `.dev/` 下（已 gitignore）。

## 方式二：只启动前端（Mock 模式）

不需要 Go、MySQL、Redis。

博客前台：

```bash
cd gin-blog-front
pnpm install
# .env.development 中设置 VITE_USE_MOCK = true
pnpm dev
```

博客后台：

```bash
cd gin-blog-admin
pnpm install
# .env.development 中设置 VITE_USE_MOCK = true
pnpm dev
```

Mock 模式下数据来自各项目的 `src/mock/`，增删改会写入内存、刷新浏览器后重置；后台登录不校验账号密码。

## 方式三：完整启动（前端 + 后端，手动）

和方式一等价，只是每步自己敲。前置依赖：Go 1.26+、Redis（默认 `127.0.0.1:6379`，用 DB 7）。数据库默认用 SQLite，无需额外安装。

1. 启动后端：

```bash
cd gin-blog-server
go mod tidy

cd cmd
go run main.go
```

2. 首次启动后初始化基础数据（菜单、资源、角色、默认用户、网站配置、页面封面）：

```bash
cd gin-blog-server/cmd
sh generate_data.sh
```

生成默认用户 `admin` 和 `guest`，密码都是 `123456`。脚本可重复执行，已存在的数据会跳过。

3. 启动前端（把 `.env.development` 中的 `VITE_USE_MOCK` 设为 `false`）：

```bash
cd gin-blog-front && pnpm install && pnpm dev
cd gin-blog-admin && pnpm install && pnpm dev
```

`/api` 请求由 vite 代理到后端 `:8765`，前端不需要配置跨域。

## 常见问题

- **修改 `.env*` 没生效**：vite 不会热更环境变量，必须重启 dev server。注意 `pnpm dev` 只读 `.env` 和 `.env.development`，改 `.env.production` 对开发模式无效。
- **`pnpm install` 报 `ERR_PNPM_IGNORED_BUILDS`**：pnpm 10+ 默认禁止依赖执行安装脚本，两个前端项目已在 `pnpm-workspace.yaml` 的 `allowBuilds` 中批准，如仍报错执行 `pnpm approve-builds` 并把选项全部设为 `true`。
- **前台页面空白、没有文章**：`generate_data.sh` 只生成系统基础数据，不含文章/分类/标签等内容数据，需要在后台自行添加。
- **注册用户要不要配邮箱**：不用。`config.yml` 的 `Captcha.SendEmail` 默认 `false`，注册请求直接把用户建出来。想改成邮箱验证注册，把它设为 `true` 并把 `Email` 段（`Host` / `Port` / `From` / `SmtpPass` / `SmtpUser`）配全，否则注册会返回 `6101 发送邮件失败`。另外注意环境里若存在 `EMAIL` 变量会让整个 `Email` 段读不到（见 `code_audit.md` F12）。
- **启动日志出现 `[警告] JWT.Secret 还是仓库里的示例值`**：本地开发可以忽略。`Server.Mode: release` 时这两项（`JWT.Secret` / `Session.Salt`）为空或仍是示例值会直接拒绝启动，用环境变量 `JWT_SECRET` / `SESSION_SALT` 注入即可；Docker 部署由 `deploy/bootstrap.sh` 自动生成。
- **改了数据库但接口仍返回旧数据**：页面封面等缓存在 Redis 且无过期时间，执行 `redis-cli -n 7 del page` 清除。

更多细节见各子项目的 README。
