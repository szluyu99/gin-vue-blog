## 博客介绍

<p  align=center>
<a  href="http://www.hahacode.cn">
<img  src="https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg"  width="200"  hight="200"  alt="阵、雨的个人博客"  style="border-radius: 50%">
</a>
</p>

<p align="center">
   <a target="_blank" href="#">
      <img src="https://img.shields.io/badge/Go-1.19-blue"/>
      <img src="https://img.shields.io/badge/Gin-v1.8.1-blue"/>
      <img src="https://img.shields.io/badge/Casbin-v2.56.0-blue"/>
      <img src="https://img.shields.io/badge/mysql-8.0-blue"/>
      <img src="https://img.shields.io/badge/GORM-v1.24.0-blue"/>
      <img src="https://img.shields.io/badge/redis-7.0-red"/>
      <img src="https://img.shields.io/badge/vue-v3.X-green"/>
    </a>
</p>

[在线预览](#在线预览) | [目录结构](#目录结构) | [技术介绍](#技术介绍) | [运行环境](#运行环境) | [开发环境](#开发环境) | [快速开始](#快速开始) | [项目总结](#项目总结)  | [鸣谢项目](#鸣谢项目) | [后续计划](#后续计划) | [更新日志](#更新日志)

因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。

您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。


Github地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)

Gitee地址：[https://gitee.com/szluyu99/gin-vue-blog](https://gitee.com/szluyu99/gin-vue-blog)

## 在线预览

博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)

博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)

> 因为博客备案还在申请中，可能无法通过域名访问，可尝试 IP 访问：
>
> 博客前台链接：[http://106.52.70.138/](http://106.52.70.138/)
>
> 博客后台链接：[http://106.52.70.138/blog-admin](http://106.52.70.138/)

测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录

在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)

> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改
>
> 并且由于一开始不会用 Apifox，接口文档生成的信息不全（如返回响应、响应示例），抽空完善一下

## 目录结构

代码仓库目录：

```bash
gin-vue-blog
├── gin-blog-admin      -- 博客后台前端
├── gin-blog-front      -- 博客前台前端
├── gin-blog-server     -- 博客后端
```

> 项目运行参考：[快速开始](#快速开始)

后端目录：

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
├── log             -- 日志文件
├── Dockerfile
└── main.go
```

前端目录：自行参考代码文件

## 项目介绍

市面上有太多优秀的前后台框架，本项目也是参考了很多开源项目，但是大多项目过于重量级（并非坏处），如果从学习的角度来看对初学者并不是很友好。本项目在以**博客**这个业务为主的前提下，提供一个完整的全栈项目代码（前台前端 + 后台前端 + 后端），技术点基本都是最新 + 最火的技术，代码轻量级，注释完善，适合学习。

前台：
- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁
- 实现点赞，统计用户等功能 (Redis)
- 评论 + 回复评论功能
- 留言采用弹幕墙，效果炫酷
- 文章详情页有文章目录、推荐文章等功能，优化用户体验

后台：
- 鉴权使用 JWT
- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理
- 支持动态权限修改，前端菜单由后端生成
- 文章编辑使用 Markdown 编辑器
- 常规后台功能齐全：侧边栏、面包屑、标签栏等
- 实现记录操作日志功能（GET 不记录）
- 实现监听在线用户、强制下线功能

其他：
- 采用 Restful 风格的 API
- 前后端分离部署，前端使用 Nginx，后端使用 Docker
- 代码整洁层次清晰，利于开发者学习
- 技术点新颖，代码轻量级

## 技术介绍

> 这里只写一些主流的通用技术，详细第三方库：前端参考 `package.json` 文件，后端参考 `go.mod` 文件

前台前端：使用 pnpm 包管理工具
- 基于 TypeScript
- Vue3
- Unocss: 原子化 CSS
- Pinia
- Vue Router 
- Axios 
- Naive UI
- Vuetify
- ...

后台前端：使用 pnpm 包管理工具
- 基于 JavaSciprt 
- Vue3
- Unocss: 原子化 CSS
- Pinia 
- Vue Router 
- Axios 
- Naive UI
- ...

后端技术栈：
- 基于 Golang
- Docker
- Gin
- GORM
- Viper: 使用 TOML 作为配置文件
- Casbin
- Zap
- MySQL
- Redis
- Nginx: 部署静态资源 + 反向代理
- ...

其他：
- 腾讯云人机验证
- 七牛云对象存储
- ...

## 运行环境

服务器：腾讯云 2核 4G Ubuntu 22.04 LTS

对象存储：七牛云

## 开发环境

| 开发工具                      | 说明                    |
| ----------------------------- | ----------------------- |
| Vscode                        | Golang 后端 +  Vue 前端 |
| Navicat                       | MySQL 远程连接工具      |
| Another Redis Desktop Manager | Redis 远程连接工具      |
| MobaXterm                     | Linux 远程工具          |
| Apifox                        | 接口调试 + 文档生成     |

| 开发环境 | 版本 |
| -------- | ---- |
| Golang   | 1.19 |
| MySQL    | 8.x  |
| Redis    | 7.x  |


### Vscode 插件

如果使用 Vscode 开发，推荐安装一下以下插件。

> 前后端相关插件属于**开发必须插件**，还有很多提升开发效果的插件后面推荐一下

前端开发插件：

| 插件 | 作用 |
| -------- | ---- |
| Volar   | Vue 官方插件 |
| UnoCSS | Unocss 官方插件 |
| Tailwind CSS IntelliSense | Tailwind 官方插件 |
| Iconify IntelliSense | Iconify 提示插件 |

后端开发插件：

| 插件 | 作用 |
| -------- | ---- |
| Go | Golang 官方插件 |


通用插件：（非开发必须，个人推荐）

| 名称 | 作用 |
| -------- | ---- |
| Better Comments   | 优化代码注释显示效果 |
| TODO Highlight | 高亮 TODO |
| Highlight Matching Tag | 高亮匹配的标签 | 


## 快速开始

### 本地运行

> 目前需要自行安装 Golang、Node、MySQL、Redis 环境
>
> TODO: 完全基于 Docker 的运行教程（Docker 实在太好用了！）


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

# 3、MySQL 导入 ginblog.sql

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
pnpm run dev
```

后台前端：

```bash
# 1、进入后台前端项目根目录
cd gin-blog-admin

# 2、安装依赖
pnpm install

# 3、运行项目
pnpm run dev
```


### 项目部署

**目前暂时不推荐将本博客部署上生产环境，因为还有很多功能未完善**。

但是相信本项目对于 Golang 学习者绝对是个合适的项目！

等功能开发的差不多了，再专门针对部署写一篇文章。

---

这里简单介绍一下，有基础的同学可以自行折腾。

本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。

后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`

Docker 打包成镜像指令：

```bash
docker build -t ginblog .
```

> 以上只是简单说明，等功能全部完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。

## 项目总结

这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。

最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。

## 鸣谢项目

- [https://butterfly.js.org/](https://butterfly.js.org/)
- [https://github.com/flipped-aurora/gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)
- [https://github.com/qifengzhang007/GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)
- [https://github.com/X1192176811/blog](https://github.com/X1192176811/blog)
- [https://github.com/zclzone/vue-naive-admin](https://github.com/zclzone/vue-naive-admin)
- ...

> 需要感谢的绝不止以上这些开源项目，但是一时难以全部列出，后面会慢慢补上。

## 后续计划

高优先级: 

- 完善图片上传功能, 目前文件上传还没怎么处理
- 后台首页重新设计（目前没放什么内容）
- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通
- 前台首页搜索文章（ElasticSearch 搜索）
- 博客文章导入导出 (.md 文件)
- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）
- 后端日志切割
- 后台修改背景图片，博客配置等
- 后端的 IP 地址检测 BUG 待修复
- 博客前台适配移动端

后续有空安排上：
- 黑夜模式
- 前台收缩侧边信息功能
- 说说
- 相册
- 音乐播放器
- 鼠标左击特效
- 看板娘
- 文章目录锚点跟随
- 第三方登录
- 评论时支持选择表情，参考 Valine
- 若干细节需要完善...

其他：
- 写一份好的文档
- 补全 README.md
- 完善 Apifox 生成的接口文档

## 更新日志

描述规范定义，有以下几种行为 ACTION:
- `FIX`: 修复
- `REFINE`: 优化
- `COMPLETE`: 完善
- `ADD`: 新增

标准描述形式：`[项目]: ` `[ACTION]`  `[模块]`，`[DETAIL 详情]`

---

2022-12-07:

博客后台：
- 修复 博客后台，动态加载路由模块导致的热加载失效问题 ⭐⭐⭐
- 完善 文件上传，抽取出单独的组件 ⭐⭐
- 修复 发布文章/查看文章 页面数据加载错误 ⭐ 
- 优化 发布文章/查看文章 时的文章标签选择、分类选择 ⭐ 
- 新增 个人信息页面 ⭐