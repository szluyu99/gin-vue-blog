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

[在线预览](#在线预览) | [目录结构](#目录结构) | [技术介绍](#技术介绍) | [运行环境](#运行环境) | [开发环境](#开发环境) | [快速开始](#快速开始) | [项目总结](#项目总结) 

因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。

您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。


Github地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)

Gitee地址：[https://gitee.com/szluyu99/gin-vue-blog](https://gitee.com/szluyu99/gin-vue-blog)

## 在线预览

博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)

博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)

测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录

在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)

> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改

## 目录结构

代码仓库目录：

```bash
gin-vue-blog
├── gin-blog-admin      -- 博客后台前端
├── gin-blog-front      -- 博客前台前端
├── gin-blog-server     -- 博客后端
```

需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。

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

## 项目介绍

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

其他：
- 前后端分离部署，前端使用 Nginx，后端使用 Docker
- 代码整洁层次清晰，利于开发者学习

## 技术介绍

> 这里只写一些主流的通用技术，详细第三方库: 前端参考 `package.json` 文件，后端参考 `go.mod` 文件

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
- pnpm: 包管理工具
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

## 快速开始

### 本地运行

> 自行安装 Golang、Node、MySQL、Redis 环境

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

目前暂时不推荐将本博客部署上生产环境，因为还有太多功能未完善。

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

> 以上只是简单说明，等功能大致完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。

## 项目总结

这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。

最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。

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
- 相册

后续有空安排上：
- 适配移动端
- 黑夜模式
- 前台收缩侧边信息功能
- 说说
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
