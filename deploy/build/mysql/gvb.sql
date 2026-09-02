/*
 Navicat Premium Data Transfer

 Source Server         : 172.18.45.12
 Source Server Type    : MySQL
 Source Server Version : 80031
 Source Host           : 172.18.45.12:3306
 Source Schema         : ginblog

 Target Server Type    : MySQL
 Target Server Version : 80031
 File Encoding         : 65001

 Date: 29/12/2023 23:17:19
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS`gvb` DEFAULT CHARACTER SET utf8mb4;
USE `gvb`;

-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS `article`;
CREATE TABLE `article`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `desc` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `img` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `type` tinyint NULL DEFAULT NULL COMMENT '类型(1-原创 2-转载 3-翻译)',
  `status` tinyint NULL DEFAULT NULL COMMENT '状态(1-公开 2-私密)',
  `is_top` tinyint(1) NULL DEFAULT NULL,
  `is_delete` tinyint(1) NULL DEFAULT NULL,
  `original_url` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `category_id` bigint NULL DEFAULT NULL,
  `user_id` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of article
-- ----------------------------
INSERT INTO `article` VALUES (1, '2023-12-27 22:46:36.066', '2023-12-27 22:49:01.365', '项目运行成功', '', '## 恭喜你，项目已经成功运行起来了!\n\n```go\nfmt.Println(\"恭喜！\")\n```\n\n```js\nconsole.log(\"恭喜！\")\n```\n\n🆗😋\n\n## 现在项目已经支持渲染公式啦!\n\n$$\nlarge X^{2m}_{3n}\n$$\n\n$$\nLarge X^{2m}_{3n}\n$$\n\n$$\nhuge X^{2m}_{3n}\n$$\n\n$$\nHuge X^{2m}_{3n}\n$$\n\n上标：$x^2$\n下标：$Y_1$\n综合：$X^{2m}_{3n}$\n\n$$ \\frac{2x+3}{3y-1} $$\n\n$$ \\sqrt[5]{6} $$\n', 'https://cdn.hahacode.cn/config/default_article_cover.png', 1, 1, 0, 0, '', 3, 1);
INSERT INTO `article` VALUES (2, '2023-12-27 22:47:47.513', '2023-12-27 22:48:58.872', '学习有捷径', '', '学习有捷径。学习的捷径之一就是多看看别人是怎么理解这些知识的。\n\n举两个例子。\n\n如果你喜欢《水浒》，千万不要只读原著当故事看，一定要读一读各代名家的批注和点评，看他们是如何理解的。之前学 C# 时，看《CLR via C#》晦涩难懂，但是我又通过看《你必须知道的.net》而更加了解了。因为后者就是中国一个 80 后写的，我通过他对 C# 的了解，作为桥梁和阶梯，再去窥探比较高达上的书籍和知识。\n\n最后，真诚的希望你能借助别人的力量来提高自己。我也一直在这样要求我自己。\n\n$$\n1/2 + 3/4 + 5/6 + 7^{99} = 999\n$$', 'https://cdn.hahacode.cn/config/default_article_cover.png', 1, 1, 0, 0, '', 4, 1);
INSERT INTO `article` VALUES (3, '2023-12-27 22:48:43.727', '2023-12-27 23:38:34.022', '项目介绍', '', '## 博客交流群\n\n这是旧版介绍，用来显示看看效果，新版的 Readme 还没来得及写！\n\n项目交流 QQ 群号：777260310\n\n## 博客介绍\n\n<p align=\"center\">\n   <a target=\"_blank\" href=\"#\">\n      <img style=\"display: inline-block;\" src=\"https://img.shields.io/badge/Go-1.19-blue\"/>\n      <img style=\"display: inline-block;\" src=\"https://img.shields.io/badge/Gin-v1.8.1-blue\"/>\n      <img style=\"display: inline-block;\" src=\"https://img.shields.io/badge/Casbin-v2.56.0-blue\"/>\n      <img style=\"display: inline-block;\" src=\"https://img.shields.io/badge/mysql-8.0-blue\"/>\n      <img style=\"display: inline-block;\" src=\"https://img.shields.io/badge/GORM-v1.24.0-blue\"/>\n      <img style=\"display: inline-block;\" src=\"https://img.shields.io/badge/redis-7.0-red\"/>\n      <img style=\"display: inline-block;\" src=\"https://img.shields.io/badge/vue-v3.X-green\"/>\n    </a>\n</p>\n\n[在线预览](#在线预览) | [项目介绍](#项目介绍) | [技术介绍](#技术介绍) | [目录结构](#目录结构) | [环境介绍](#环境介绍) | [快速开始](#快速开始) | [总结&鸣谢](#总结鸣谢)  | [后续计划](#后续计划)\n\n您的 Star 是我坚持的动力，感谢大家的支持，欢迎提交 Pr 共同改进项目。\n\nGithub 地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\n\nGitee 地址：[https://gitee.com/szluyu99/gin-vue-blog](https://gitee.com/szluyu99/gin-vue-blog)\n\n## 在线预览\n\n博客前台链接：[hahacode.cn](https://www.hahacode.cn)（已适配移动端）\n\n博客后台链接：[hahacode.cn/admin](https://www.hahacode.cn/admin)（暂未专门适配移动端）\n\n> 博客域名已通过备案，且配置 SSL，通过 https 访问\n\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\n\n> 在线接口文档地址：[doc.hahacode.cn](http://doc.hahacode.cn/)，准备换成 Swagger\n\n## 有 Docker 环境可一键启动效果\n\nLinux/Mac 可直接运行，Windows 要使用 GitBash 运行（默认终端不能执行 shell）\n\n```bash\ngit clone https://github.com/szluyu99/gin-vue-blog \ncd gin-vue-blog/deploy\n./bootstrap.sh\n```\n\n## 项目介绍\n\nGithub 上有很多优秀的前后台框架，本项目也参考了许多开源项目，但是大多项目都比较重量级（并非坏处），如果从学习的角度来看对初学者并不是很友好。本项目在以**博客**这个业务为主的前提下，提供一个完整的全栈项目代码（前台前端 + 后台前端 + 后端），技术点基本都是最新 + 最火的技术，代码轻量级，注释完善，适合学习。\n\n同时，本项目可用于一键搭建动态博客（参考 [快速开始](#快速开始)）。\n\n前台：\n\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\n- 响应式布局，适配了移动端\n- 实现点赞，统计用户等功能 (Redis)\n- 评论 + 回复评论功能\n- 留言采用弹幕墙，效果炫酷\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\n\n后台：\n\n- 鉴权使用 JWT\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\n- 支持动态权限修改，前端菜单由后端生成（动态路由）\n- 文章编辑使用 Markdown 编辑器\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\n- 实现记录操作日志功能（GET 不记录）\n- 实现监听在线用户、强制下线功能\n- 文件上传支持七牛云、本地（后续计划支持更多）\n- 对 CRUD 操作封装了通用 Hook\n\n其他：\n\n- 采用 Restful 风格的 API\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\n- 代码整洁层次清晰，利于开发者学习\n- 技术点新颖，代码轻量级，适度封装\n- Docker Compose 一键运行，轻松搭建在线博客\n\n### 技术介绍\n\n> 这里写一些主流的通用技术，详细第三方库：前端参考 `package.json` 文件，后端参考 `go.mod` 文件\n\n前端技术栈: 使用 pnpm 包管理工具\n\n- 基于 TypeScript\n- Vue3\n- VueUse: 服务于 Vue Composition API 的工具集\n- Unocss: 原子化 CSS\n- Pinia\n- Vue Router \n- Axios \n- Naive UI\n- ...\n\n后端技术栈:\n\n- Golang\n- Docker\n- Gin\n- GORM\n- Viper: 支持 TOML (默认)、YAML 等常用格式作为配置文件\n- Casbin\n- Zap\n- MySQL\n- Redis\n- Nginx: 部署静态资源 + 反向代理\n- ...\n\n其他:\n\n- 腾讯云人机验证\n- 七牛云对象存储\n- ...\n\n### 目录结构\n\n> 这里简单列出目录结构，具体可以查看源码\n\n代码仓库目录：\n\n```bash\ngin-vue-blog\n├── gin-blog-admin      -- 博客后台前端\n├── gin-blog-front      -- 博客前台前端\n├── gin-blog-server     -- 博客后端\n├── deploy              -- 部署\n```\n\n> 项目运行参考：[快速开始](#快速开始)\n\n后端目录：简略版\n\n```bash\ngin-blog-server\n├── api             -- API\n│   ├── front       -- 前台接口\n│   └── v1          -- 后台接口\n├── dao             -- 数据库操作模块\n├── service         -- 服务模块\n├── model           -- 数据模型\n│   ├── req             -- 请求 VO 模型\n│   ├── resp            -- 响应 VO 模型\n│   ├── dto             -- 内部传输 DTO 模型\n│   └── ...             -- 数据库模型对象 PO 模型\n├── routes          -- 路由模块\n│   └── middleware      -- 路由中间件\n├── utils           -- 工具模块\n│   ├── r               -- 响应封装\n│   ├── upload          -- 文件上传\n│   └── ...\n├── routes          -- 路由模块\n├── config          -- 配置文件\n├── test            -- 测试模块\n├── assets          -- 资源文件\n├── log             -- 存放日志的目录\n├── public          -- 外部访问的静态资源\n│   └── uploaded    -- 本地文件上传目录\n├── Dockerfile\n└── main.go\n```\n\n前端目录：简略版\n\n```\ngin-vue-admin / gin-vue-front 通用目录结构\n├── src              \n│   ├── api             -- 接口\n│   ├── assets          -- 静态资源\n│   ├── styles          -- 样式\n│   ├── components      -- 组件\n│   ├── composables     -- 组合式函数\n│   ├── router          -- 路由\n│   ├── store           -- 状态管理\n│   ├── utils           -- 工具方法\n│   ├── views           -- 页面\n│   ├── App.vue\n│   └── main.ts\n├── settings         -- 项目配置\n├── build            -- 构建相关的配置\n├── public           -- 公共资源, 在打包后会被加到 dist 根目录\n├── package.json \n├── package-lock.json\n├── index.html\n├── tsconfig.json\n├── unocss.config.ts -- unocss 配置\n└── vite.config.ts   -- vite 配置\n├── .env             -- 通用环境变量\n├── .env.development -- 开发环境变量\n├── .env.production  -- 线上环境变量\n├── .gitignore\n├── .editorconfig    -- 编辑器配置\n```\n\n部署目录：简略版\n\n```\ndeploy\n├── build      -- 镜像构建\n│   ├── mysql  -- mysql 镜像构建\n│   ├── server -- 后端镜像构建 (基于 gin-blog-server 目录)\n│   └── web    -- 前端镜像构建 (基于前端项目打包的静态资源)\n└── start\n    ├── docker-compose.yml    -- 多容器管理\n    └── .env                  -- 环境变量\n    └── ...\n```\n\n## 环境介绍\n\n### 线上环境\n\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\n\n对象存储：七牛云\n\n### 开发环境\n\n| 开发工具                          | 说明                  |\n| ----------------------------- | ------------------- |\n| Vscode                        | Golang 后端 +  Vue 前端 |\n| Navicat                       | MySQL 远程连接工具        |\n| Another Redis Desktop Manager | Redis 远程连接工具        |\n| MobaXterm                     | Linux 远程工具          |\n| Apifox                        | 接口调试 + 文档生成         |\n\n| 开发环境   | 版本   |\n| ------ | ---- |\n| Golang | 1.19 |\n| MySQL  | 8.x  |\n| Redis  | 7.x  |\n\n### VsCode 插件\n\n目前推荐安装插件已经写到 `.vscode/extensions.json` 中，使用 VsCode 打开项目会推荐安装。\n\n> 注意，使用 VsCode 打开 gin-blog-admin 和 gin-blog-front 这两个项目，而不是打开 gin-vue-blog 这个目录！\n\n## 快速开始\n\n建议下载本项目后，先一键运行起来，查看本项目在本地的运行效果。\n\n需要修改源码的话，参考常规运行，前后端分开运行。\n\n本项目开发环境是 Linux，如果 Windows 下运行有奇奇怪怪的问题，可以进群交流或提 Issue\n\n### 拉取项目前的准备 (Windows)\n\n如果是 Windows 系统，需要先执行以下指令，否则 Docker 构建过程可能会出 BUG。\n\n或者直接下载 ZIP 而不是通过 git clone 克隆项目。\n\nLinux 和 Mac 不需要进行该操作。\n\n> 原因是该项目开发时基于 Linux，本项目规范使用 lf 换行符。而 Windows 的 git 在自动拉取项目时会将项目文件中换行符转换为 crlf，经过测试，构建过程会产生 BUG。\n\n```bash\n# 防止 git 自动将换行符转换为 crlf\ngit config --global core.autocrlf false\n```\n\n### 方式一：Docker Compose 一键运行\n\n需要有 Docker 和 Docker Compose 的环境\n\n> 详细运行文档（包含环境搭建）参考：[deploy/README.md](https://github.com/szluyu99/gin-vue-blog/tree/main/deploy)\n\nLinux 下可以正常启动：（Windows 请使用 `GitBash` 进行操作）\n\n```bash\ngit clone https://github.com/szluyu99/gin-vue-blog \ncd gin-vue-blog/deploy\n./bootstrap.sh\n```\n\n本地前台访问 [localhost](http://localhost/)\n\n本地后台访问 [localhost/admin](http://localhost/admin)\n\n默认用户：\n\n- 管理员 admin 123456\n- 普通用户 user 123456\n- 测试用户 test 123456\n\n如果运行遇到问题，请查看详细文章 [deploy/README.md](https://github.com/szluyu99/gin-vue-blog/tree/main/deploy)\n\n### 方式二：常规运行\n\n需要安装 Golang、Node、MySQL、Redis 环境：\n \n- Golang 安装参考 [官方文档](https://go.dev/doc/install)\n\n- Node 安装建议使用 [Nvm](https://nvm.uihtm.com/)，也可以直接去 [Node 官网](https://nodejs.org/en) 下载\n\n- MySQL、Redis 建议使用 Docker 安装\n\n> 以下使用 Docker 安装环境，未做持久化处理，仅用于开发和演示\n\nDocker 安装 MySQL：\n\n```bash\n# 注意: 必须安装 MySQL 8.0 以上版本\ndocker pull mysql:8.0\n\n# 运行 MySQL\ndocker run --name mysql8 -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 -d mysql:8.0\n\n# 查看是否运行成功, STATUS 为 Up 即成功\ndocker ps\n\n# 进入容器, CTRL + D 退出\ndocker exec -it mysql8 bash\nmysql -u root -p123456\n```\n\nDocker 安装 Redis：\n\n```bash\ndocker pull redis:7.0\n\n# 运行 Redis\ndocker run --name redis7 -p 6379:6379 -d redis:7.0\n\n# 查看是否运行成功, STATUS 为 Up 即成功\ndocker ps\n\n# 进入容器, CTRL + D 退出\ndocker exec -it redis7 bash\nredis-cli\n```\n\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\n\n拉取项目到本地：\n\n```bash\ngit clone https://github.com/szluyu99/gin-vue-blog.git\n```\n\n后端项目运行：\n\n```bash\n# 1、进入后端项目根目录 \ncd gin-blog-server\n\n# 2、修改项目运行的配置文件，默认加载位于 config/config.toml \n\n# 3、MySQL 导入 gvb.sql\n\n# 4、启动 Redis \n\n# 5、运行项目\ngo mod tidy\ngo run main.go\n```\n\n数据库中的默认用户：\n\n- 管理员 admin 123456\n- 普通用户 user 123456\n- 测试用户 test 123456\n\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 `pnpm`\n\n```bash\nnpm install -g pnpm\n```\n\n前台前端：\n\n```bash\n# 1、进入前台前端项目根目录\ncd gin-blog-front\n\n# 2、安装依赖\npnpm install\n\n# 3、运行项目\npnpm dev\n```\n\n后台前端：\n\n```bash\n# 1、进入后台前端项目根目录\ncd gin-blog-admin\n\n# 2、安装依赖\npnpm install\n\n# 3、运行项目\npnpm dev\n```\n\n### 项目部署\n\nTODO\n\n## 总结鸣谢\n\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\n\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\n\n 鸣谢项目：\n\n- [https://butterfly.js.org/](https://butterfly.js.org/)\n- [https://github.com/qifengzhang007/GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)\n- [https://github.com/zclzone/vue-naive-admin](https://github.com/zclzone/vue-naive-admin)\n- [https://github.com/antfu/vitesse](https://github.com/antfu/vitesse)\n- ...\n\n⭐ 博客后台的前端基于 [vue-naive-admin](https://github.com/zclzone/vue-naive-admin) 进行二开，感谢作者的开源。但是和原项目区别较大，详见 [gin-blog-admin/README.md](https://github.com/szluyu99/gin-vue-blog/tree/main/gin-blog-admin)\n\n> 需要感谢的绝不止以上这些开源项目，但是一时难以全部列出，后面会慢慢补上。\n\n## 后续计划\n\n高优先级: \n\n- ~~完善图片上传功能, 目前文件上传还没怎么处理~~ 🆗\n- 后台首页重新设计（目前没放什么内容）\n- ~~前台首页搜索文章（目前使用数据库模糊搜索）~~ 🆗\n- ~~博客文章导入导出 (.md 文件)~~ 🆗\n- ~~权限管理中菜单编辑时选择图标（现在只能输入图标字符串）~~ 🆗\n- 后端日志切割\n- ~~后台修改背景图片，博客配置等~~ 🆗\n- ~~后端的 IP 地址检测 BUG 待修复~~ 🆗\n- ~~博客前台适配移动端~~ 🆗\n- ~~文章详情, 目录锚点跟随~~ 🆗\n- ~~邮箱注册 + 邮件发送验证码~~ 🆗\n- 修改测试环境的数据库为 SQLite3，方便运行\n\n后续有空安排上：\n\n- 黑夜模式\n- 前台收缩侧边信息功能\n- 说说\n- 相册\n- 音乐播放器\n- 鼠标左击特效\n- 看板娘\n- 第三方登录: QQ、微信、Github ...\n- 评论时支持选择表情，参考 Valine\n- 单独部署：前后端 + 环境\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\n- 前台首页搜索集成 ElasticSearch\n- 国际化?\n\n其他：\n\n- 写一份好的文档\n- 补全 README.md\n- 完善 Apifox 生成的接口文档\n- ~~一键部署：使用 docker compose 单机一键部署整个项目（前后端 + 环境）~~ 🆗\n', 'https://cdn.hahacode.cn/config/default_article_cover.png', 1, 1, 1, 0, '', 3, 1);

-- ----------------------------
-- Table structure for article_tag
-- ----------------------------
DROP TABLE IF EXISTS `article_tag`;
CREATE TABLE `article_tag`  (
  `tag_id` bigint NOT NULL,
  `article_id` bigint NOT NULL,
  PRIMARY KEY (`tag_id`, `article_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of article_tag
-- ----------------------------
INSERT INTO `article_tag` VALUES (1, 1);
INSERT INTO `article_tag` VALUES (1, 3);
INSERT INTO `article_tag` VALUES (2, 1);
INSERT INTO `article_tag` VALUES (2, 3);
INSERT INTO `article_tag` VALUES (3, 2);

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of category
-- ----------------------------
INSERT INTO `category` VALUES (1, '2023-12-27 22:45:09.369', '2023-12-27 22:45:09.369', '后端');
INSERT INTO `category` VALUES (2, '2023-12-27 22:45:15.006', '2023-12-27 22:45:15.006', '前端');
INSERT INTO `category` VALUES (3, '2023-12-27 22:46:36.057', '2023-12-27 22:46:36.057', '项目');
INSERT INTO `category` VALUES (4, '2023-12-27 22:47:47.501', '2023-12-27 22:47:47.501', '学习');

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `user_id` bigint NULL DEFAULT NULL,
  `reply_user_id` bigint NULL DEFAULT NULL,
  `topic_id` bigint NULL DEFAULT NULL,
  `parent_id` bigint NULL DEFAULT NULL,
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` tinyint(1) NOT NULL COMMENT '评论类型(1.文章 2.友链 3.说说)',
  `is_review` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comment
-- ----------------------------

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `config` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `key` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `value` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `desc` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `key`(`key` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of config
-- ----------------------------
INSERT INTO `config` VALUES (1, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.029', '', 'website_avatar', 'https://foruda.gitee.com/avatar/1677041571085433939/5221991_szluyu99_1614389421.png', '网站头像');
INSERT INTO `config` VALUES (2, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.033', '', 'website_name', '阵雨的个人博客', '网站名称');
INSERT INTO `config` VALUES (3, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.027', '', 'website_author', '阵雨', '网站作者');
INSERT INTO `config` VALUES (4, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.023', '', 'website_intro', '往事随风而去', '网站介绍');
INSERT INTO `config` VALUES (5, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.038', '', 'website_notice', '欢迎来到阵雨的个人博客，项目还在开发中...', '网站公告');
INSERT INTO `config` VALUES (6, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.031', '', 'website_createtime', '2023-12-27 22:40:22', '网站创建日期');
INSERT INTO `config` VALUES (7, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.011', '', 'website_record', '粤ICP备2021032312号', '网站备案号');
INSERT INTO `config` VALUES (8, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.008', '', 'qq', '123456789', 'QQ');
INSERT INTO `config` VALUES (9, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.015', '', 'github', 'https://github.com/szluyu99', 'github');
INSERT INTO `config` VALUES (10, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.025', '', 'gitee', 'https://gitee.com/szluyu99', 'gitee');
INSERT INTO `config` VALUES (11, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.019', '', 'tourist_avatar', 'https://cdn.hahacode.cn/config/tourist_avatar.png', '默认游客头像');
INSERT INTO `config` VALUES (12, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.021', '', 'user_avatar', 'https://cdn.hahacode.cn/config/user_avatar.png', '默认用户头像');
INSERT INTO `config` VALUES (13, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.013', '', 'article_cover', 'https://cdn.hahacode.cn/config/default_article_cover.png', '默认文章封面');
INSERT INTO `config` VALUES (14, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.039', '', 'is_comment_review', 'true', '评论默认审核');
INSERT INTO `config` VALUES (15, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.017', '', 'is_message_review', 'true', '留言默认审核');
INSERT INTO `config` VALUES (16, '2023-12-27 22:59:20.110', '2023-12-27 23:01:35.035', '', 'about', '```javascript\nconsole.log(\"Hello World\")\n```\n\n我就是我，不一样的烟火！', '');

-- ----------------------------
-- Table structure for friend_link
-- ----------------------------
DROP TABLE IF EXISTS `friend_link`;
CREATE TABLE `friend_link`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of friend_link


-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '头像地址',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '留言内容',
  `ip_address` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'IP 地址',
  `ip_source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'IP 来源',
  `speed` tinyint(1) NULL DEFAULT NULL COMMENT '弹幕速度',
  `is_review` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of message
-- ----------------------------

-- ----------------------------
-- Table structure for operation_log
-- ----------------------------
DROP TABLE IF EXISTS `operation_log`;
CREATE TABLE `operation_log`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `opt_module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作模块',
  `opt_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作类型',
  `opt_method` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作方法',
  `opt_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作URL',
  `opt_desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作描述',
  `request_param` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '请求参数',
  `request_method` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '请求方法',
  `response_data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '响应数据',
  `user_id` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户昵称',
  `ip_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作IP',
  `ip_source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作地址',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of operation_log
-- ----------------------------

-- ----------------------------
-- Table structure for page
-- ----------------------------
DROP TABLE IF EXISTS `page`;
CREATE TABLE `page`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `label` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  UNIQUE INDEX `label`(`label` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of page
-- ----------------------------
INSERT INTO `page` VALUES (1, '2022-12-08 13:09:58.500', '2023-12-28 16:31:43.682', '首页', 'home', 'https://cdn.hahacode.cn/page/home.jpg');
INSERT INTO `page` VALUES (2, '2022-12-08 13:51:49.474', '2023-12-28 14:55:58.704', '归档', 'archive', 'https://cdn.hahacode.cn/page/tag.png');
INSERT INTO `page` VALUES (3, '2022-12-08 13:52:18.084', '2023-12-28 16:31:30.137', '分类', 'category', 'https://cdn.hahacode.cn/page/category.png');
INSERT INTO `page` VALUES (4, '2022-12-08 13:52:31.364', '2023-12-28 14:55:45.058', '标签', 'tag', 'https://cdn.hahacode.cn/page/tag.png');
INSERT INTO `page` VALUES (5, '2022-12-08 13:52:52.389', '2023-12-28 15:02:21.859', '友链', 'link', 'https://cdn.hahacode.cn/page/link.jpg');
INSERT INTO `page` VALUES (6, '2022-12-08 13:53:04.159', '2023-12-28 16:30:03.928', '关于', 'about', 'https://cdn.hahacode.cn/page/about.jpg');
INSERT INTO `page` VALUES (7, '2022-12-08 13:53:17.707', '2023-12-28 16:27:13.418', '留言', 'message', 'https://cdn.hahacode.cn/page/message.jpeg');
INSERT INTO `page` VALUES (8, '2022-12-08 13:53:30.187', '2023-12-28 14:55:25.724', '个人中心', 'user', 'https://cdn.hahacode.cn/page/user.jpg');
INSERT INTO `page` VALUES (9, '2022-12-16 23:54:52.650', '2023-12-28 14:54:42.341', '相册', 'album', 'https://cdn.hahacode.cn/page/album.png');
INSERT INTO `page` VALUES (10, '2022-12-16 23:55:36.059', '2023-12-28 14:55:09.345', '错误页面', '404', 'https://cdn.hahacode.cn/page/404.jpg');
INSERT INTO `page` VALUES (11, '2022-12-16 23:56:17.917', '2023-12-28 16:33:16.644', '文章列表', 'article_list', 'https://cdn.hahacode.cn/page/article_list.jpg');


-- ----------------------------
-- Table structure for tag
-- ----------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tag
-- ----------------------------
INSERT INTO `tag` VALUES (1, '2023-12-27 22:45:40.731', '2023-12-27 22:45:40.731', 'Golang');
INSERT INTO `tag` VALUES (2, '2023-12-27 22:46:36.082', '2023-12-27 22:46:36.082', 'Vue');
INSERT INTO `tag` VALUES (3, '2023-12-27 22:47:47.530', '2023-12-27 22:47:47.530', '感悟');

-- ----------------------------
-- Table structure for user_auth
-- ----------------------------
DROP TABLE IF EXISTS `user_auth`;
CREATE TABLE `user_auth`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `login_type` tinyint(1) NULL DEFAULT NULL COMMENT '登录类型',
  `ip_address` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录IP地址',
  `ip_source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'IP来源',
  `last_login_time` datetime(3) NULL DEFAULT NULL,
  `is_disable` tinyint(1) NULL DEFAULT NULL,
  `is_super` tinyint(1) NULL DEFAULT NULL,
  `user_info_id` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_auth
-- ----------------------------
INSERT INTO `user_auth` VALUES (1, '2023-12-27 22:40:42.565', '2023-12-29 23:13:11.500', 'superadmin', '$2a$10$R1kus4SbUJ5QzAgcUuxrbOhv10j.CaVtUdmRbZ17C2552frAj7opW', 1, '172.18.45.12', '内网IP', '2023-12-29 23:13:11.500', 0, 1, 1);
INSERT INTO `user_auth` VALUES (2, '2022-10-31 21:54:11.040', '2023-12-27 23:44:06.416', 'admin', '$2a$10$urGRaFQoLoblBUUdgi1NCei3oGnCHJwtHFmVcIfC94135KdNffy4.', 1, '172.18.45.12', '内网IP', '2023-12-27 23:44:06.416', 0, 0, 2);
INSERT INTO `user_auth` VALUES (3, '2022-11-01 10:41:13.300', '2023-12-29 23:04:48.284', 'test@qq.com', '$2a$10$FmU4jxwDlibSL9pdt.AsuODkbB4gLp3IyyXeoMmW/XALtT/HdwTsi', 1, '172.18.45.12', '内网IP', '2023-12-29 23:04:48.284', 0, 0, 3);
INSERT INTO `user_auth` VALUES (4, '2022-10-19 22:31:26.805', '2023-12-26 21:10:35.334', 'user', '$2a$10$9vHpoeT7sF4j9beiZfPsOe0jJ67gOceO2WKJzJtHRZCjNJajl7Fhq', 1, '172.12.0.6:48716', '', '2022-12-24 12:13:52.494', 0, 0, 4);


-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `email` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `nickname` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `avatar` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `website` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `nickname`(`nickname` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_info
-- ----------------------------
INSERT INTO `user_info` VALUES (1, '2023-12-27 22:40:42.495', '2023-12-28 16:34:24.836', '', 'superadmin', 'public/uploaded/4c50eef3bdaf0b4164ce179e576f2b2d_20231228163423.gif', '这个人很懒，什么都没有留下', 'https://www.hahacode.cn');
INSERT INTO `user_info` VALUES (2, '2022-10-31 21:54:10.935', '2023-12-27 23:44:01.402', '', '管理员', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我就是我，不一样的烟火', 'https://www.hahacode.cn');
INSERT INTO `user_info` VALUES (3, '2022-10-19 22:31:26.734', '2023-12-27 23:31:39.169', 'user@qq.com', '普通用户', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我是个普通用户！', 'https://www.hahacode.cn');
INSERT INTO `user_info` VALUES (4, '2022-11-01 10:41:13.234', '2023-12-27 23:31:42.587', 'test@qq.com', '测试用户', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我是测试用的！', 'https://www.hahacode.cn');

SET FOREIGN_KEY_CHECKS = 1;
