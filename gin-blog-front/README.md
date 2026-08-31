# gin-blog-front

博客前台。Vue 3 + Vite + UnoCSS + Pinia，开发端口 `8888`。

支持两种运行模式：**Mock 模式**（不需要后端，全部数据来自本地假数据）和**接后端模式**。

## 安装依赖

```bash
pnpm install
```

> 使用 pnpm 10+ 时，依赖的 install/postinstall 脚本默认被禁止执行，未批准会直接报 `ERR_PNPM_IGNORED_BUILDS`。
> 本项目已在 `pnpm-workspace.yaml` 的 `allowBuilds` 中批准了必需的几个（`esbuild` 等），正常安装即可。

## Mock 模式（推荐用于只改前端）

`.env.development` 中设置：

```ini
VITE_USE_MOCK = true
```

```bash
pnpm dev
```

打开 http://localhost:8888 即可，不需要启动 Go 后端和 Redis。

实现方式是替换 axios 的 `adapter`（`src/utils/http.js`），请求不会发出浏览器，但仍会走原有的响应拦截器，所以业务状态码、token 注入、登录弹框等逻辑与真实环境一致。

- 假数据在 `src/mock/data.js`，接口路由表在 `src/mock/index.js`，覆盖 `src/api.js` 中的全部接口
- 文章列表真分页、支持分类/标签筛选与搜索高亮
- 点赞、发评论、发留言会写入内存，页面内立即生效，**刷新浏览器后重置**
- 登录不校验账号密码，任意输入都会登录成功

## 接后端模式

`.env.development` 中设置：

```ini
VITE_USE_MOCK = false
VITE_BACKEND_URL = 'http://localhost:8765'
```

然后按 `gin-blog-server/README.md` 启动后端，再 `pnpm dev`。`/api` 请求由 vite 代理到 `VITE_BACKEND_URL`。

> 修改 `.env*` 后必须重启 dev server，vite 不会热更环境变量。

## 其他命令

```bash
pnpm build    # 生产构建
pnpm lint     # 代码检查
```
