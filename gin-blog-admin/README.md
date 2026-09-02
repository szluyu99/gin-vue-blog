# gin-blog-admin

博客后台。Vue 3 + Vite + Naive UI + UnoCSS，开发端口 `8889`。项目骨架来自 [vue-naive-admin](https://github.com/zclzone/vue-naive-admin)。

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

打开 http://localhost:8889 即可，不需要启动 Go 后端和 Redis。

实现方式是替换 axios 的 `adapter`（`src/utils/http.js`），请求不会发出浏览器，但仍会走原有的响应拦截器，所以业务状态码、token 注入、鉴权跳转等逻辑与真实环境一致。

- 假数据在 `src/mock/data.js`，接口路由表在 `src/mock/index.js`，覆盖 `src/api.js` 中的全部接口
- 菜单、资源、角色三份数据与后端 `generate_data.sh` 生成的默认数据一致，保证动态路由和权限页面可用
- 列表页支持分页与查询条件；新增、编辑、删除、置顶、审核、禁用、强制下线都会写入内存，**刷新浏览器后重置**
- 登录不校验账号密码，任意输入都会以 admin 身份登录

## 接后端模式

`.env.development` 中设置：

```ini
VITE_USE_MOCK = false
VITE_SERVER_URL = 'http://localhost:8765'
```

然后按 `gin-blog-server/README.md` 启动后端并初始化数据，再 `pnpm dev`，用 `admin / 123456` 登录。`/api` 请求由 vite 代理到 `VITE_SERVER_URL`。

> 修改 `.env*` 后必须重启 dev server，vite 不会热更环境变量。

## 路由与权限

`VITE_BACK_ROUTER = true` 时走后端路由：后端 `/menu/user/list` 返回当前用户的菜单树，前端在 `src/store/modules/permission.js` 的 `buildRoutes` 中组装成 vue-router 路由并动态注册。

`VITE_BACK_ROUTER = false` 时走前端路由：加载 `src/views/**/route.js`，按 `meta.requireAuth` 和角色过滤。

> 注意 vue-router 4.3+ 不允许子路由与父路由同名，会直接抛异常导致动态路由全部注册失败（表现为登录后所有页面跳 404）。新增菜单时不要让父子同名。

## 其他命令

```bash
pnpm build    # 生产构建
pnpm lint     # 代码检查
pnpm test     # 单元测试 (vitest)
```

代码风格使用 [antfu/eslint-config](https://github.com/antfu/eslint-config)，不使用 Prettier。

## 打包体积

体积大的依赖都做了按需加载，改动时注意别把它们拉回首屏 chunk：

- `exceljs`（gzip 约 250KB）在 `CrudTable.vue` 里用动态 `import()`，点击导出才加载
- `md-editor-v3` + CodeMirror 只被文章编辑页和「关于」页引入
- `highlight.js` 只有操作日志页的 `NCode` 需要，在该页引入后通过 `:hljs` 传入，不要挂到 `App.vue` 的 `NConfigProvider` 上

`pnpm build` 会生成 `stats.html`（rollup-plugin-visualizer），可以直接打开看各 chunk 的构成。
