## 博客交流群

项目交流 QQ 群号：777260310

## 博客介绍

<p align=center>
<a href="http://www.hahacode.cn">
<img src="https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg"  width="200" hight="200" alt="阵、雨的个人博客" style="border-radius: 50%">
</a>
</p>

<p align="center">
   <a target="_blank" href="#">
      <img style="display: inline-block;" src="https://img.shields.io/badge/Go-1.19-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/Gin-v1.8.1-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/Casbin-v2.56.0-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/mysql-8.0-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/GORM-v1.24.0-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/redis-7.0-red"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/vue-v3.X-green"/>
    </a>
</p>

[在线预览](#在线预览) | [项目介绍](#项目介绍) | [技术介绍](#技术介绍) | [目录结构](#目录结构) | [环境介绍](#环境介绍) | [快速开始](#快速开始) | [总结&鸣谢](#总结鸣谢)  | [后续计划](#后续计划) | [更新日志](#更新日志)

您的 Star 是我坚持的动力，感谢大家的支持，欢迎提交 Pr 共同改进项目。

Github 地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)

Gitee 地址：[https://gitee.com/szluyu99/gin-vue-blog](https://gitee.com/szluyu99/gin-vue-blog)

## 在线预览

博客前台链接：[hahacode.cn](https://www.hahacode.cn)（已适配移动端）

博客后台链接：[hahacode.cn/admin](https://www.hahacode.cn/admin)（暂未专门适配移动端）

> 博客域名已通过备案，且配置 SSL，通过 https 访问

测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录

在线接口文档地址：[doc.hahacode.cn](http://doc.hahacode.cn/)

> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改
> 
> 由于一开始不会用，接口文档生成的信息不全（如返回响应、响应示例），后续完善

以下放几张简单的预览图，强烈建议点击上面的[预览链接](#在线预览)进去体验下：

![前台首页图片](https://img-blog.csdnimg.cn/fd25f0e52cdc467893925a48d0d66662.png)

![前台首页文章列表](https://img-blog.csdnimg.cn/005cee463d3c4e729a1dd7bbee41e963.png)

![后台文章列表](https://img-blog.csdnimg.cn/18c39f63b7b64f7bbd2a4ab764552d13.png)

## 项目介绍

Github 上有很多优秀的前后台框架，本项目也参考了许多开源项目，但是大多项目都比较重量级（并非坏处），如果从学习的角度来看对初学者并不是很友好。本项目在以**博客**这个业务为主的前提下，提供一个完整的全栈项目代码（前台前端 + 后台前端 + 后端），技术点基本都是最新 + 最火的技术，代码轻量级，注释完善，适合学习。

同时，本项目可用于一键搭建动态博客（参考 [快速开始](#快速开始)）。

前台：

- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁
- 响应式布局，适配了移动端
- 实现点赞，统计用户等功能 (Redis)
- 评论 + 回复评论功能
- 留言采用弹幕墙，效果炫酷
- 文章详情页有文章目录、推荐文章等功能，优化用户体验

后台：

- 鉴权使用 JWT
- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理
- 支持动态权限修改，前端菜单由后端生成（动态路由）
- 文章编辑使用 Markdown 编辑器
- 常规后台功能齐全：侧边栏、面包屑、标签栏等
- 实现记录操作日志功能（GET 不记录）
- 实现监听在线用户、强制下线功能
- 文件上传支持七牛云、本地（后续计划支持更多）
- 对 CRUD 操作封装了通用 Hook

其他：

- 采用 Restful 风格的 API
- 前后端分离部署，前端使用 Nginx，后端使用 Docker
- 代码整洁层次清晰，利于开发者学习
- 技术点新颖，代码轻量级，适度封装
- Docker Compose 一键运行，轻松搭建在线博客

### 技术介绍

> 这里只写一些主流的通用技术，详细第三方库：前端参考 `package.json` 文件，后端参考 `go.mod` 文件

前端技术栈: 使用 pnpm 包管理工具

- 基于 TypeScript
- Vue3
- VueUse: 服务于 Vue Composition API 的工具集
- Unocss: 原子化 CSS
- Pinia
- Vue Router 
- Axios 
- Naive UI
- ...

后端技术栈:

- Golang
- Docker
- Gin
- GORM
- Viper: 支持 TOML (默认)、YAML 等常用格式作为配置文件
- Casbin
- Zap
- MySQL
- Redis
- Nginx: 部署静态资源 + 反向代理
- ...

其他:

- 腾讯云人机验证
- 七牛云对象存储
- ...

### 目录结构

> 这里简单列出目录结构，具体可以查看源码

代码仓库目录：

```bash
gin-vue-blog
├── gin-blog-admin      -- 博客后台前端
├── gin-blog-front      -- 博客前台前端
├── gin-blog-server     -- 博客后端
├── deploy              -- 部署
```

> 项目运行参考：[快速开始](#快速开始)

后端目录：简略版

```bash
gin-blog-server
├── api             -- API
│   ├── front       -- 前台接口
│   └── v1          -- 后台接口
├── dao             -- 数据库操作模块
├── service         -- 服务模块
├── model           -- 数据模型
│   ├── req             -- 请求 VO 模型
│   ├── resp            -- 响应 VO 模型
│   ├── dto             -- 内部传输 DTO 模型
│   └── ...             -- 数据库模型对象 PO 模型
├── routes          -- 路由模块
│   └── middleware      -- 路由中间件
├── utils           -- 工具模块
│   ├── r               -- 响应封装
│   ├── upload          -- 文件上传
│   └── ...
├── routes          -- 路由模块
├── config          -- 配置文件
├── test            -- 测试模块
├── assets          -- 资源文件
├── log             -- 存放日志的目录
├── public          -- 外部访问的静态资源
│   └── uploaded    -- 本地文件上传目录
├── Dockerfile
└── main.go
```

前端目录：简略版

```
gin-vue-admin / gin-vue-front 通用目录结构
├── src              
│   ├── api             -- 接口
│   ├── assets          -- 静态资源
│   ├── styles          -- 样式
│   ├── components      -- 组件
│   ├── composables     -- 组合式函数
│   ├── router          -- 路由
│   ├── store           -- 状态管理
│   ├── utils           -- 工具方法
│   ├── views           -- 页面
│   ├── App.vue
│   └── main.ts
├── settings         -- 项目配置
├── build            -- 构建相关的配置
├── public           -- 公共资源, 在打包后会被加到 dist 根目录
├── package.json 
├── package-lock.json
├── index.html
├── tsconfig.json
├── unocss.config.ts -- unocss 配置
└── vite.config.ts   -- vite 配置
├── .env             -- 通用环境变量
├── .env.development -- 开发环境变量
├── .env.production  -- 线上环境变量
├── .gitignore
├── .editorconfig    -- 编辑器配置
```

部署目录：简略版

```
deploy
├── build      -- 镜像构建
│   ├── mysql  -- mysql 镜像构建
│   ├── server -- 后端镜像构建 (基于 gin-blog-server 目录)
│   └── web    -- 前端镜像构建 (基于前端项目打包的静态资源)
└── start
    ├── docker-compose.yml    -- 多容器管理
    └── .env                  -- 环境变量
    └── ...
```

## 环境介绍

### 线上环境

服务器：腾讯云 2核 4G Ubuntu 22.04 LTS

对象存储：七牛云

### 开发环境

| 开发工具                          | 说明                  |
| ----------------------------- | ------------------- |
| Vscode                        | Golang 后端 +  Vue 前端 |
| Navicat                       | MySQL 远程连接工具        |
| Another Redis Desktop Manager | Redis 远程连接工具        |
| MobaXterm                     | Linux 远程工具          |
| Apifox                        | 接口调试 + 文档生成         |

| 开发环境   | 版本   |
| ------ | ---- |
| Golang | 1.19 |
| MySQL  | 8.x  |
| Redis  | 7.x  |

### VsCode 插件

目前推荐安装插件已经写到 `.vscode/extensions.json` 中，使用 VsCode 打开项目会推荐安装。

> 注意，一定要用 VsCode 打开 gin-blog-admin 和 gin-blog-front 这两个前端项目，而不是打开 gin-vue-blog 这个目录！

前端必备插件：

| 插件                   | 作用           |
| -------------------- | ------------ |
| Volar                | Vue 官方插件     |
| UnoCSS               | Unocss 官方插件  |
| Iconify IntelliSense | Iconify 提示插件 |

后端必备插件：

| 插件  | 作用          |
| --- | ----------- |
| Go  | Golang 官方插件 |

其他插件：（个人推荐，用于提升开发体验）

| 名称                     | 作用         |
| ---------------------- | ---------- |
| Better Comments        | 优化代码注释显示效果 |
| TODO Highlight         | 高亮 TODO    |
| Highlight Matching Tag | 高亮匹配的标签    |

## 快速开始

建议下载本项目后，先一键运行起来，查看本项目在本地的运行效果。

需要修改源码的话，参考常规运行，前后端分开运行。

### 拉取项目前的准备 (Windows)

如果是 Windows 系统，需要先执行以下指令，否则 Docker 构建过程会出 BUG。

或者直接下载 ZIP 而不是通过 git clone 克隆项目。

Linux 和 Mac 不需要进行该操作。

> 原因是该项目开发时基于 Linux，本项目规范使用 lf 换行符。而 Windows 的 git 在自动拉取项目时会将项目文件中换行符转换为 crlf，经过测试，构建过程会产生 BUG。

```bash
# 防止 git 自动将换行符转换为 crlf
git config --global core.autocrlf false
```

### 方式一：Docker Compose 一键运行

需要有 Docker 和 Docker Compose 的环境

> 详细运行文档（包含环境搭建）参考：[deploy/README.md](https://github.com/szluyu99/gin-vue-blog/tree/main/deploy)

```bash
# 拉取项目, 进入根目录
git clone https://github.com/szluyu99/gin-vue-blog 
cd gin-vue-blog

# 进入 docker-compose 目录
cd deploy/start

# 一键运行
docker-compose up -d
```

本地前台访问 [localhost](http://localhost/)

本地后台访问 [localhost/admin](http://localhost/admin)

默认用户：

- 管理员 admin 123456
- 普通用户 user 123456
- 测试用户 test 123456

如果运行遇到问题，请查看详细文章 [deploy/README.md](https://github.com/szluyu99/gin-vue-blog/tree/main/deploy)

### 方式二：常规运行

> 自行安装 Golang、Node、MySQL、Redis 环境：
> 
> Golang 安装参考 [官方文档](https://go.dev/doc/install)
> Node 安装建议使用 [Nvm](https://nvm.uihtm.com/) 
> MySQL、Redis 建议使用 Docker 运行

需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。

拉取项目到本地：

```bash
git clone https://github.com/szluyu99/gin-vue-blog.git
```

后端项目运行：

```bash
# 1、进入后端项目根目录 
cd gin-blog-server

# 2、修改项目运行的配置文件，默认加载位于 config/config.toml 

# 3、MySQL 导入 gvb.sql

# 4、启动 Redis 

# 5、运行项目
go mod tidy
go run main.go
```

数据库中的默认用户：

- 管理员 admin 123456
- 普通用户 user 123456
- 测试用户 test 123456

前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm

```bash
npm install -g pnpm
```

前台前端：

```bash
# 1、进入前台前端项目根目录
cd gin-blog-front

# 2、安装依赖
pnpm install

# 3、运行项目
pnpm dev
```

后台前端：

```bash
# 1、进入后台前端项目根目录
cd gin-blog-admin

# 2、安装依赖
pnpm install

# 3、运行项目
pnpm dev
```

### 项目部署

> TODO

## 总结鸣谢

这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。

最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。

 鸣谢项目：

- [https://butterfly.js.org/](https://butterfly.js.org/)
- [https://github.com/qifengzhang007/GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)
- [https://github.com/X1192176811/blog](https://github.com/X1192176811/blog)
- [https://github.com/zclzone/vue-naive-admin](https://github.com/zclzone/vue-naive-admin)
- [https://github.com/antfu/vitesse](https://github.com/antfu/vitesse)
- ...

⭐ 博客后台的前端基于 [vue-naive-admin](https://github.com/zclzone/vue-naive-admin) 进行二开，感谢作者的开源。

> 需要感谢的绝不止以上这些开源项目，但是一时难以全部列出，后面会慢慢补上。

## 后续计划

高优先级: 

- ~~完善图片上传功能, 目前文件上传还没怎么处理~~ 🆗
- 后台首页重新设计（目前没放什么内容）
- ~~前台首页搜索文章（目前使用数据库模糊搜索）~~ 🆗
- ~~博客文章导入导出 (.md 文件)~~ 🆗
- ~~权限管理中菜单编辑时选择图标（现在只能输入图标字符串）~~ 🆗
- 后端日志切割
- ~~后台修改背景图片，博客配置等~~ 🆗
- ~~后端的 IP 地址检测 BUG 待修复~~ 🆗
- ~~博客前台适配移动端~~ 🆗
- ~~文章详情, 目录锚点跟随~~ 🆗
- ~~邮箱注册 + 邮件发送验证码~~ 🆗

后续有空安排上：

- 黑夜模式
- 前台收缩侧边信息功能
- 说说
- 相册
- 音乐播放器
- 鼠标左击特效
- 看板娘
- 第三方登录: QQ、微信、Github ...
- 评论时支持选择表情，参考 Valine
- 单独部署：前后端 + 环境
- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通
- 前台首页搜索集成 ElasticSearch
- 国际化?

其他：

- 写一份好的文档
- 补全 README.md
- 完善 Apifox 生成的接口文档
- ~~一键部署：使用 docker compose 单机一键部署整个项目（前后端 + 环境）~~ 🆗

## 更新日志

描述规范定义，有以下几种行为 ACTION:

- `FIX`: 修复
- `REFINE`: 优化
- `COMPLETE`: 完善
- `ADD`: 新增

---

2023-06-18:

博客后台:
- 大规模重构，精简代码，更新依赖，更新结构


---

2023-01-01:

博客前台:

- 新增邮箱注册

---

2022-12-23:

博客部署:

- 新增 deploy 子模块, 使用 docker-compose 一键运行前后端 ⭐⭐

---

2022-12-22:

博客部署:

- 修复 线上图片上传问题, 本地上传则采用 Nginx 反向代理 Gin 的静态服务资源, 也可采用七牛云存储

---

2022-12-21:

博客后台:

- 新增 图标选择组件, 选择范围为自定义图标数组 ⭐

---

2022-12-20:

博客后台: 

- 完善 IP 检测工具 ⭐
- 修复 Nginx 获取的客户端 IP 是服务器 IP (docker 部署 Nginx 的问题, 修改 Nginx 配置文件即可)

---

2022-12-19:

博客后台:

- 添加简易版黑夜模式 (一行 CSS 实现)

---

2022-12-17:

博客后台:

- 新增 文章导入、导出 ⭐

---

2022-12-16:

博客后台:

- 修复 MySQL 无法存储 emoji 的问题，设置 utf8mb4 字符集编码，更新 SQL 
- 完善 添加文章时选择标签的逻辑
- 新增 全屏水印 功能

---

2022-12-15:

博客后台：

- 优化 菜单栏和标签栏，点击标签自动展开对应菜单，点击菜单自动滚动到显示对应标签 ⭐
- 优化 使用响应式语法糖重构代码
- 优化 代码结构 + 注释

---

2022-12-14:

博客前台：

- 优化 代码，去除 n-card 组件，使用自定义 css 实现卡片视图
- 优化 对滚动的监听操作，添加节流函数，提升性能 ⭐
- 优化 文章目录，根据当前滚动条自动高亮锚点 ⭐⭐
- 新增 vite 打包优化相关插件

---

2022-12-13:

项目部署：

- 新增 Nginx 配置 https 访问域名 (http 自动跳转到 https) ⭐
- 新增 七牛云添加加速域名访问图片资源

博客后端：

- 新增 文章搜索接口（数据库模糊查询）

博客前台：

- 新增 导航栏的文章搜索

---

2022-12-12:

博客前台：

- 新增 适配移动端 ⭐⭐
- 优化 删除 Vuetify 相关组件及依赖 (css 样式存在冲突)，统一使用 naive-ui
- 优化 使用 `$ref` 语法糖功能重构页面

---

2022-12-09:

博客前台：

- 完善 个人中心的头像上传（注意头像上传需要 Token）
- 优化 重构通用页面的代码 ⭐

---

2022-12-08:

博客后端：

- 新增 页面模块 的接口 ⭐
- 修复 单元测试无法初始化全局变量（flag 包冲突）问题

博客后台：

- 新增 pinia 数据持久化，防止刷新丢失数据 ⭐
- 新增 页面管理页面，动态配置前台显示的封面 ⭐

博客前台：

- 新增 页面根据 `label` 从后端数据中加载封面 ⭐
- 优化 尝试性使用 `$ref` 语法糖功能

---

2022-12-07:

博客后台：

- 修复 博客后台，动态加载路由模块导致的热加载失效问题 ⭐⭐
- 完善 文件上传，抽取出单独的组件 ⭐⭐
- 新增 个人信息页面 ⭐
- 优化 发布文章/查看文章 时的文章标签选择、分类选择
- 修复 发布文章/查看文章 页面数据加载错误
