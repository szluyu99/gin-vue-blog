# 快速开始

三个子项目：`gin-blog-server`（后端）、`gin-blog-front`（博客前台）、`gin-blog-admin`（博客后台）。

前台和后台都内置 Mock 模式，只想看前端效果的话不需要启动后端。

| 项目 | 地址 | 说明 |
| --- | --- | --- |
| 博客前台 | http://localhost:8888 | |
| 博客后台 | http://localhost:8889 | 默认账号 `admin / 123456` |
| 后端接口 | http://localhost:8765/api | Swagger: http://localhost:8765/swagger/index.html |

## 方式一：只启动前端（Mock 模式）

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

## 方式二：完整启动（前端 + 后端）

前置依赖：Go 1.21+、Redis（默认 `127.0.0.1:6379`，用 DB 7）。数据库默认用 SQLite，无需额外安装。

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
- **改了数据库但接口仍返回旧数据**：页面封面等缓存在 Redis 且无过期时间，执行 `redis-cli -n 7 del page` 清除。

更多细节见各子项目的 README。
