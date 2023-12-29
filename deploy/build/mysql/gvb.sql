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

-- ----------------------------
-- Table structure for menu
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `parent_id` bigint NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `component` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `order_num` tinyint NULL DEFAULT NULL,
  `redirect` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `catalogue` tinyint(1) NULL DEFAULT NULL,
  `hidden` tinyint(1) NULL DEFAULT NULL,
  `keep_alive` tinyint(1) NULL DEFAULT NULL,
  `external` tinyint(1) NULL DEFAULT NULL,
  `external_link` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 49 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (2, '2022-10-31 09:41:03.000', '2023-12-27 23:26:43.807', 0, '文章管理', '/article', 'Layout', 'ic:twotone-article', 1, '/article/list', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (3, '2022-10-31 09:41:03.000', '2023-12-24 23:33:34.013', 0, '消息管理', '/message', 'Layout', 'ic:twotone-email', 2, '/message/comment	', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (4, '2022-10-31 09:41:03.000', '2023-12-24 23:32:35.177', 0, '用户管理', '/user', 'Layout', 'ph:user-list-bold', 4, '/user/list', 0, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (5, '2022-10-31 09:41:03.000', '2023-12-24 23:32:34.788', 0, '系统管理', '/setting', 'Layout', 'ion:md-settings', 5, '/setting/website', 0, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (6, '2022-10-31 09:41:03.000', '2023-12-24 23:22:29.519', 2, '发布文章', 'write', '/article/write', 'icon-park-outline:write', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (8, '2022-10-31 09:41:03.000', '2023-12-21 20:58:29.873', 2, '文章列表', 'list', '/article/list', 'material-symbols:format-list-bulleted', 2, '', 0, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (9, '2022-10-31 09:41:03.000', '2022-11-01 01:18:30.931', 2, '分类管理', 'category', '/article/category', 'tabler:category', 3, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (10, '2022-10-31 09:41:03.000', '2022-11-01 01:18:35.502', 2, '标签管理', 'tag', '/article/tag', 'tabler:tag', 4, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (16, '2022-10-31 09:41:03.000', '2022-11-01 10:11:23.195', 0, '权限管理', '/auth', 'Layout', 'cib:adguard', 3, '/auth/menu', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (17, '2022-10-31 09:41:03.000', NULL, 16, '菜单管理', 'menu', '/auth/menu', 'ic:twotone-menu-book', 1, NULL, 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (23, '2022-10-31 09:41:03.000', NULL, 16, '接口管理', 'resource', '/auth/resource', 'mdi:api', 2, NULL, 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (24, '2022-10-31 09:41:03.000', '2022-10-31 10:09:18.913', 16, '角色管理', 'role', '/auth/role', 'carbon:user-role', 3, NULL, 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (25, '2022-10-31 10:11:09.232', '2022-11-01 01:29:48.520', 3, '评论管理', 'comment', '/message/comment', 'ic:twotone-comment', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (26, '2022-10-31 10:12:01.546', '2022-11-01 01:29:54.130', 3, '留言管理', 'leave-msg', '/message/leave-msg', 'ic:twotone-message', 2, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (27, '2022-10-31 10:54:03.201', '2022-11-01 01:30:06.901', 4, '用户列表', 'list', '/user/list', 'mdi:account', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (28, '2022-10-31 10:54:34.167', '2022-11-01 01:30:13.400', 4, '在线用户', 'online', '/user/online', 'ic:outline-online-prediction', 2, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (29, '2022-10-31 10:59:33.255', '2022-11-01 01:30:20.688', 5, '网站管理', 'website', '/setting/website', 'el:website', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (30, '2022-10-31 11:00:09.997', '2022-11-01 01:30:24.097', 5, '页面管理', 'page', '/setting/page', 'iconoir:journal-page', 2, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (31, '2022-10-31 11:00:33.543', '2022-11-01 01:30:28.497', 5, '友链管理', 'link', '/setting/link', 'mdi:telegram', 3, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (32, '2022-10-31 11:01:00.444', '2022-11-01 01:30:33.186', 5, '关于我', 'about', '/setting/about', 'cib:about-me', 4, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (33, '2022-11-01 01:43:10.142', '2023-12-27 23:26:41.553', 0, '首页', '/home', '/home', 'ic:sharp-home', 0, '', 1, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (34, '2022-11-01 09:54:36.252', '2022-11-01 10:07:00.254', 2, '修改文章', 'write/:id', '/article/write', 'icon-park-outline:write', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (36, '2022-11-04 15:50:45.993', '2023-12-24 23:32:33.538', 0, '日志管理', '/log', 'Layout', 'material-symbols:receipt-long-outline-rounded', 6, '/log/operation', 0, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (37, '2022-11-04 15:53:00.251', '2023-12-24 23:15:22.034', 36, '操作日志', 'operation', '/log/operation', 'mdi:book-open-page-variant-outline', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (38, '2022-11-04 16:02:42.306', '2022-11-04 16:05:35.761', 36, '登录日志', 'login', '/log/login', 'material-symbols:login', 2, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (39, '2022-12-07 20:47:08.349', '2023-12-24 23:33:35.701', 0, '个人中心', '/profile', '/profile', 'mdi:account', 7, '', 1, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (47, '2023-12-24 20:26:14.173', '2023-12-24 23:33:36.247', 0, '测试一级菜单', '/testone', 'Layout', '', 88, '', 0, 0, 0, 1, NULL);
INSERT INTO `menu` VALUES (48, '2023-12-24 23:26:19.441', '2023-12-24 23:26:27.704', 0, '测试外链', 'https://www.baidu.com', 'Layout', 'mdi-fan-speed-3', 66, '', 1, 0, 0, 1, '');

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
-- Table structure for resource
-- ----------------------------
DROP TABLE IF EXISTS `resource`;
CREATE TABLE `resource`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `parent_id` bigint NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `anonymous` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 117 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of resource
-- ----------------------------
INSERT INTO `resource` VALUES (3, '2022-10-20 22:42:00.664', '2022-10-20 22:42:00.664', 0, '', '', '文章模块', 0);
INSERT INTO `resource` VALUES (6, '2022-10-20 22:42:23.349', '2022-10-20 22:42:23.349', 0, '', '', '留言模块', 0);
INSERT INTO `resource` VALUES (7, '2022-10-20 22:42:28.550', '2022-10-20 22:42:28.550', 0, '', '', '菜单模块', 0);
INSERT INTO `resource` VALUES (8, '2022-10-20 22:42:31.623', '2022-10-20 22:42:31.623', 0, '', '', '角色模块', 0);
INSERT INTO `resource` VALUES (9, '2022-10-20 22:42:36.262', '2022-10-20 22:42:36.262', 0, '', '', '评论模块', 0);
INSERT INTO `resource` VALUES (10, '2022-10-20 22:42:40.700', '2022-10-20 22:42:40.700', 0, '', '', '资源模块', 0);
INSERT INTO `resource` VALUES (11, '2022-10-20 22:42:51.023', '2022-10-20 22:42:51.023', 0, '', '', '博客信息模块', 0);
INSERT INTO `resource` VALUES (23, '2022-10-22 22:13:12.455', '2022-10-26 11:15:23.546', 10, '/resource/anonymous', 'PUT', '修改资源匿名访问', 0);
INSERT INTO `resource` VALUES (34, '2022-10-31 17:14:11.708', '2022-10-31 17:14:11.708', 10, '/resource', 'POST', '新增/编辑资源', 0);
INSERT INTO `resource` VALUES (35, '2022-10-31 17:14:42.320', '2022-10-31 17:18:52.508', 10, '/resource/list', 'GET', '资源列表', 0);
INSERT INTO `resource` VALUES (36, '2022-10-31 17:15:14.999', '2022-10-31 17:19:01.460', 10, '/resource/option', 'GET', '资源选项列表(树形)', 0);
INSERT INTO `resource` VALUES (37, '2022-10-31 17:16:56.830', '2022-10-31 17:16:56.830', 10, '/resource/:id', 'DELETE', '删除资源', 0);
INSERT INTO `resource` VALUES (38, '2022-10-31 17:19:28.905', '2022-10-31 17:19:28.905', 7, '/menu/list', 'GET', '菜单列表', 0);
INSERT INTO `resource` VALUES (39, '2022-10-31 18:46:33.051', '2022-10-31 18:46:33.051', 7, '/menu', 'POST', '新增/编辑菜单', 0);
INSERT INTO `resource` VALUES (40, '2022-10-31 18:46:53.804', '2022-10-31 18:46:53.804', 7, '/menu/:id', 'DELETE', '删除菜单', 0);
INSERT INTO `resource` VALUES (41, '2022-10-31 18:47:17.272', '2022-10-31 18:47:28.130', 7, '/menu/option', 'GET', '菜单选项列表(树形)', 0);
INSERT INTO `resource` VALUES (42, '2022-10-31 18:48:04.780', '2022-10-31 18:48:04.780', 7, '/menu/user/list', 'GET', '获取当前用户菜单', 0);
INSERT INTO `resource` VALUES (43, '2022-10-31 19:20:35.427', '2023-12-27 23:21:22.669', 3, '/article/list', 'GET', '文章列表', 0);
INSERT INTO `resource` VALUES (44, '2022-10-31 19:21:02.096', '2023-12-27 22:07:57.702', 3, '/article/:id', 'GET', '文章详情', 0);
INSERT INTO `resource` VALUES (45, '2022-10-31 19:26:04.763', '2022-10-31 19:26:09.709', 3, '/article', 'POST', '新增/编辑文章', 0);
INSERT INTO `resource` VALUES (46, '2022-10-31 19:26:36.453', '2022-10-31 19:26:36.453', 3, '/article/soft-delete', 'PUT', '软删除文章', 0);
INSERT INTO `resource` VALUES (47, '2022-10-31 19:26:52.344', '2022-10-31 19:26:52.344', 3, '/article', 'DELETE', '删除文章', 0);
INSERT INTO `resource` VALUES (48, '2022-10-31 19:27:07.731', '2022-10-31 19:27:07.731', 3, '/article/top', 'PUT', '修改文章置顶', 0);
INSERT INTO `resource` VALUES (49, '2022-10-31 20:19:55.588', '2022-10-31 20:19:55.588', 0, '', '', '分类模块', 0);
INSERT INTO `resource` VALUES (50, '2022-10-31 20:20:03.400', '2022-10-31 20:20:03.400', 0, '', '', '标签模块', 0);
INSERT INTO `resource` VALUES (51, '2022-10-31 20:22:03.799', '2022-10-31 20:22:03.799', 49, '/category/list', 'GET', '分类列表', 0);
INSERT INTO `resource` VALUES (52, '2022-10-31 20:22:28.840', '2022-10-31 20:22:28.840', 49, '/category', 'POST', '新增/编辑分类', 0);
INSERT INTO `resource` VALUES (53, '2022-10-31 20:31:04.577', '2022-10-31 20:31:04.577', 49, '/category', 'DELETE', '删除分类', 0);
INSERT INTO `resource` VALUES (54, '2022-10-31 20:31:36.612', '2022-10-31 20:31:36.612', 49, '/category/option', 'GET', '分类选项列表', 0);
INSERT INTO `resource` VALUES (55, '2022-10-31 20:32:57.112', '2022-10-31 20:33:13.261', 50, '/tag/list', 'GET', '标签列表', 0);
INSERT INTO `resource` VALUES (56, '2022-10-31 20:33:29.080', '2022-10-31 20:33:29.080', 50, '/tag', 'POST', '新增/编辑标签', 0);
INSERT INTO `resource` VALUES (57, '2022-10-31 20:33:39.992', '2022-10-31 20:33:39.992', 50, '/tag', 'DELETE', '删除标签', 0);
INSERT INTO `resource` VALUES (58, '2022-10-31 20:33:53.962', '2022-10-31 20:33:53.962', 50, '/tag/option', 'GET', '标签选项列表', 0);
INSERT INTO `resource` VALUES (59, '2022-10-31 20:35:05.647', '2022-10-31 20:35:05.647', 6, '/message/list', 'GET', '留言列表', 0);
INSERT INTO `resource` VALUES (60, '2022-10-31 20:35:25.551', '2022-10-31 20:35:25.551', 6, '/message', 'DELETE', '删除留言', 0);
INSERT INTO `resource` VALUES (61, '2022-10-31 20:36:20.587', '2022-10-31 20:36:20.587', 6, '/message/review', 'PUT', '修改留言审核', 0);
INSERT INTO `resource` VALUES (62, '2022-10-31 20:37:04.637', '2022-10-31 20:37:04.637', 9, '/comment/list', 'GET', '评论列表', 0);
INSERT INTO `resource` VALUES (63, '2022-10-31 20:37:29.779', '2022-10-31 20:37:29.779', 9, '/comment', 'DELETE', '删除评论', 0);
INSERT INTO `resource` VALUES (64, '2022-10-31 20:37:40.317', '2022-10-31 20:37:40.317', 9, '/comment/review', 'PUT', '修改评论审核', 0);
INSERT INTO `resource` VALUES (65, '2022-10-31 20:38:30.506', '2022-10-31 20:38:30.506', 8, '/role/list', 'GET', '角色列表', 0);
INSERT INTO `resource` VALUES (66, '2022-10-31 20:38:50.606', '2022-10-31 20:38:50.606', 8, '/role', 'POST', '新增/编辑角色', 0);
INSERT INTO `resource` VALUES (67, '2022-10-31 20:39:03.752', '2022-10-31 20:39:03.752', 8, '/role', 'DELETE', '删除角色', 0);
INSERT INTO `resource` VALUES (68, '2022-10-31 20:39:28.232', '2022-10-31 20:39:28.232', 8, '/role/option', 'GET', '角色选项', 0);
INSERT INTO `resource` VALUES (69, '2022-10-31 20:44:22.622', '2022-10-31 20:44:22.622', 0, '', '', '友链模块', 0);
INSERT INTO `resource` VALUES (70, '2022-10-31 20:44:41.334', '2022-10-31 20:44:41.334', 69, '/link/list', 'GET', '友链列表', 0);
INSERT INTO `resource` VALUES (71, '2022-10-31 20:45:01.150', '2022-10-31 20:45:01.150', 69, '/link', 'POST', '新增/编辑友链', 0);
INSERT INTO `resource` VALUES (72, '2022-10-31 20:45:12.406', '2022-10-31 20:45:12.406', 69, '/link', 'DELETE', '删除友链', 0);
INSERT INTO `resource` VALUES (74, '2022-10-31 20:46:48.330', '2022-10-31 20:47:01.505', 0, '', '', '用户信息模块', 0);
INSERT INTO `resource` VALUES (78, '2022-10-31 20:51:15.607', '2022-10-31 20:51:15.607', 74, '/user/list', 'GET', '用户列表', 0);
INSERT INTO `resource` VALUES (79, '2022-10-31 20:55:15.496', '2022-10-31 20:55:15.496', 11, '/setting/blog-config', 'GET', '获取博客设置', 0);
INSERT INTO `resource` VALUES (80, '2022-10-31 20:55:48.257', '2022-10-31 20:55:48.257', 11, '/setting/about', 'GET', '获取关于我', 0);
INSERT INTO `resource` VALUES (81, '2022-10-31 20:56:21.722', '2022-10-31 20:56:21.722', 11, '/setting/blog-config', 'PUT', '修改博客设置', 0);
INSERT INTO `resource` VALUES (82, '2022-10-31 21:57:30.021', '2022-10-31 21:57:30.021', 74, '/user/info', 'GET', '获取当前用户信息', 0);
INSERT INTO `resource` VALUES (84, '2022-10-31 22:06:18.348', '2022-10-31 22:07:38.004', 74, '/user', 'PUT', '修改用户信息', 0);
INSERT INTO `resource` VALUES (85, '2022-11-02 11:55:05.395', '2022-11-02 11:55:05.395', 11, '/setting/about', 'PUT', '修改关于我', 0);
INSERT INTO `resource` VALUES (86, '2022-11-02 13:20:09.485', '2022-11-02 13:20:09.485', 74, '/user/online', 'GET', '获取在线用户列表', 0);
INSERT INTO `resource` VALUES (91, '2022-11-03 16:42:31.558', '2022-11-03 16:42:31.558', 0, '', '', '操作日志模块', 0);
INSERT INTO `resource` VALUES (92, '2022-11-03 16:42:49.681', '2022-11-03 16:42:49.681', 91, '/operation/log/list', 'GET', '获取操作日志列表', 0);
INSERT INTO `resource` VALUES (93, '2022-11-03 16:43:04.906', '2022-11-03 16:43:04.906', 91, '/operation/log', 'DELETE', '删除操作日志', 0);
INSERT INTO `resource` VALUES (95, '2022-11-05 14:22:48.240', '2022-11-05 14:22:48.240', 11, '/home', 'GET', '获取后台首页信息', 0);
INSERT INTO `resource` VALUES (98, '2022-11-29 23:35:42.865', '2022-11-29 23:35:42.865', 74, '/user/offline', 'DELETE', '强制离线用户', 0);
INSERT INTO `resource` VALUES (99, '2022-12-07 20:48:05.939', '2022-12-07 20:48:05.939', 74, '/user/current/password', 'PUT', '修改当前用户密码', 0);
INSERT INTO `resource` VALUES (100, '2022-12-07 20:48:35.511', '2022-12-07 20:48:35.511', 74, '/user/current', 'PUT', '修改当前用户信息', 0);
INSERT INTO `resource` VALUES (101, '2022-12-07 20:55:08.271', '2022-12-07 20:55:08.271', 74, '/user/disable', 'PUT', '修改用户禁用', 0);
INSERT INTO `resource` VALUES (102, '2022-12-08 15:43:15.421', '2022-12-08 15:43:15.421', 0, '', '', '页面模块', 0);
INSERT INTO `resource` VALUES (103, '2022-12-08 15:43:26.009', '2022-12-08 15:43:26.009', 102, '/page/list', 'GET', '页面列表', 0);
INSERT INTO `resource` VALUES (104, '2022-12-08 15:43:38.570', '2022-12-08 15:43:38.570', 102, '/page', 'POST', '新增/编辑页面', 0);
INSERT INTO `resource` VALUES (105, '2022-12-08 15:43:50.879', '2022-12-08 15:43:50.879', 102, '/page', 'DELETE', '删除页面', 0);
INSERT INTO `resource` VALUES (106, '2022-12-16 11:53:57.989', '2022-12-16 11:53:57.989', 0, '', '', '文件模块', 0);
INSERT INTO `resource` VALUES (107, '2022-12-16 11:54:20.891', '2022-12-16 11:54:20.891', 106, '/upload', 'POST', '文件上传', 0);
INSERT INTO `resource` VALUES (108, '2022-12-18 01:34:47.800', '2022-12-18 01:34:47.800', 3, '/article/export', 'POST', '导出文章', 0);
INSERT INTO `resource` VALUES (109, '2022-12-18 01:34:59.255', '2022-12-18 01:34:59.255', 3, '/article/import', 'POST', '导入文章', 0);

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `label` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `is_disable` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  UNIQUE INDEX `label`(`label` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES (1, '2023-12-27 23:16:38.105', '2023-12-27 23:34:10.830', '管理员', 'admin', 0);
INSERT INTO `role` VALUES (2, '2023-12-27 23:16:50.687', '2023-12-29 23:13:46.460', '普通用户', 'user', 0);
INSERT INTO `role` VALUES (3, '2023-12-27 23:17:00.263', '2023-12-27 23:38:15.697', 'test', '测试用户', 0);

-- ----------------------------
-- Table structure for role_menu
-- ----------------------------
DROP TABLE IF EXISTS `role_menu`;
CREATE TABLE `role_menu`  (
  `menu_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`menu_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_menu
-- ----------------------------
INSERT INTO `role_menu` VALUES (2, 1);
INSERT INTO `role_menu` VALUES (2, 2);
INSERT INTO `role_menu` VALUES (2, 3);
INSERT INTO `role_menu` VALUES (3, 1);
INSERT INTO `role_menu` VALUES (3, 2);
INSERT INTO `role_menu` VALUES (3, 3);
INSERT INTO `role_menu` VALUES (4, 1);
INSERT INTO `role_menu` VALUES (4, 2);
INSERT INTO `role_menu` VALUES (4, 3);
INSERT INTO `role_menu` VALUES (5, 1);
INSERT INTO `role_menu` VALUES (5, 2);
INSERT INTO `role_menu` VALUES (5, 3);
INSERT INTO `role_menu` VALUES (6, 1);
INSERT INTO `role_menu` VALUES (6, 2);
INSERT INTO `role_menu` VALUES (6, 3);
INSERT INTO `role_menu` VALUES (8, 1);
INSERT INTO `role_menu` VALUES (8, 2);
INSERT INTO `role_menu` VALUES (8, 3);
INSERT INTO `role_menu` VALUES (9, 1);
INSERT INTO `role_menu` VALUES (9, 2);
INSERT INTO `role_menu` VALUES (9, 3);
INSERT INTO `role_menu` VALUES (10, 1);
INSERT INTO `role_menu` VALUES (10, 2);
INSERT INTO `role_menu` VALUES (10, 3);
INSERT INTO `role_menu` VALUES (16, 1);
INSERT INTO `role_menu` VALUES (16, 2);
INSERT INTO `role_menu` VALUES (16, 3);
INSERT INTO `role_menu` VALUES (17, 1);
INSERT INTO `role_menu` VALUES (17, 2);
INSERT INTO `role_menu` VALUES (17, 3);
INSERT INTO `role_menu` VALUES (23, 1);
INSERT INTO `role_menu` VALUES (23, 2);
INSERT INTO `role_menu` VALUES (23, 3);
INSERT INTO `role_menu` VALUES (24, 1);
INSERT INTO `role_menu` VALUES (24, 2);
INSERT INTO `role_menu` VALUES (24, 3);
INSERT INTO `role_menu` VALUES (25, 1);
INSERT INTO `role_menu` VALUES (25, 2);
INSERT INTO `role_menu` VALUES (25, 3);
INSERT INTO `role_menu` VALUES (26, 1);
INSERT INTO `role_menu` VALUES (26, 2);
INSERT INTO `role_menu` VALUES (26, 3);
INSERT INTO `role_menu` VALUES (27, 1);
INSERT INTO `role_menu` VALUES (27, 2);
INSERT INTO `role_menu` VALUES (27, 3);
INSERT INTO `role_menu` VALUES (28, 1);
INSERT INTO `role_menu` VALUES (28, 2);
INSERT INTO `role_menu` VALUES (28, 3);
INSERT INTO `role_menu` VALUES (29, 1);
INSERT INTO `role_menu` VALUES (29, 2);
INSERT INTO `role_menu` VALUES (29, 3);
INSERT INTO `role_menu` VALUES (30, 1);
INSERT INTO `role_menu` VALUES (30, 2);
INSERT INTO `role_menu` VALUES (30, 3);
INSERT INTO `role_menu` VALUES (31, 1);
INSERT INTO `role_menu` VALUES (31, 2);
INSERT INTO `role_menu` VALUES (31, 3);
INSERT INTO `role_menu` VALUES (32, 1);
INSERT INTO `role_menu` VALUES (32, 2);
INSERT INTO `role_menu` VALUES (32, 3);
INSERT INTO `role_menu` VALUES (33, 1);
INSERT INTO `role_menu` VALUES (33, 2);
INSERT INTO `role_menu` VALUES (33, 3);
INSERT INTO `role_menu` VALUES (34, 1);
INSERT INTO `role_menu` VALUES (34, 2);
INSERT INTO `role_menu` VALUES (34, 3);
INSERT INTO `role_menu` VALUES (36, 1);
INSERT INTO `role_menu` VALUES (36, 2);
INSERT INTO `role_menu` VALUES (36, 3);
INSERT INTO `role_menu` VALUES (37, 1);
INSERT INTO `role_menu` VALUES (37, 2);
INSERT INTO `role_menu` VALUES (37, 3);
INSERT INTO `role_menu` VALUES (38, 1);
INSERT INTO `role_menu` VALUES (38, 2);
INSERT INTO `role_menu` VALUES (38, 3);
INSERT INTO `role_menu` VALUES (39, 1);
INSERT INTO `role_menu` VALUES (39, 2);
INSERT INTO `role_menu` VALUES (39, 3);
INSERT INTO `role_menu` VALUES (47, 1);
INSERT INTO `role_menu` VALUES (48, 1);

-- ----------------------------
-- Table structure for role_resource
-- ----------------------------
DROP TABLE IF EXISTS `role_resource`;
CREATE TABLE `role_resource`  (
  `resource_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`resource_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_resource
-- ----------------------------
INSERT INTO `role_resource` VALUES (3, 1);
INSERT INTO `role_resource` VALUES (6, 1);
INSERT INTO `role_resource` VALUES (7, 1);
INSERT INTO `role_resource` VALUES (8, 1);
INSERT INTO `role_resource` VALUES (9, 1);
INSERT INTO `role_resource` VALUES (10, 1);
INSERT INTO `role_resource` VALUES (11, 1);
INSERT INTO `role_resource` VALUES (23, 1);
INSERT INTO `role_resource` VALUES (34, 1);
INSERT INTO `role_resource` VALUES (35, 1);
INSERT INTO `role_resource` VALUES (35, 2);
INSERT INTO `role_resource` VALUES (35, 3);
INSERT INTO `role_resource` VALUES (36, 1);
INSERT INTO `role_resource` VALUES (36, 2);
INSERT INTO `role_resource` VALUES (36, 3);
INSERT INTO `role_resource` VALUES (37, 1);
INSERT INTO `role_resource` VALUES (38, 1);
INSERT INTO `role_resource` VALUES (38, 2);
INSERT INTO `role_resource` VALUES (38, 3);
INSERT INTO `role_resource` VALUES (39, 1);
INSERT INTO `role_resource` VALUES (40, 1);
INSERT INTO `role_resource` VALUES (41, 1);
INSERT INTO `role_resource` VALUES (41, 2);
INSERT INTO `role_resource` VALUES (41, 3);
INSERT INTO `role_resource` VALUES (42, 1);
INSERT INTO `role_resource` VALUES (42, 2);
INSERT INTO `role_resource` VALUES (42, 3);
INSERT INTO `role_resource` VALUES (43, 1);
INSERT INTO `role_resource` VALUES (43, 2);
INSERT INTO `role_resource` VALUES (43, 3);
INSERT INTO `role_resource` VALUES (44, 1);
INSERT INTO `role_resource` VALUES (44, 2);
INSERT INTO `role_resource` VALUES (44, 3);
INSERT INTO `role_resource` VALUES (45, 1);
INSERT INTO `role_resource` VALUES (46, 1);
INSERT INTO `role_resource` VALUES (47, 1);
INSERT INTO `role_resource` VALUES (48, 1);
INSERT INTO `role_resource` VALUES (48, 2);
INSERT INTO `role_resource` VALUES (48, 3);
INSERT INTO `role_resource` VALUES (49, 1);
INSERT INTO `role_resource` VALUES (50, 1);
INSERT INTO `role_resource` VALUES (51, 1);
INSERT INTO `role_resource` VALUES (51, 2);
INSERT INTO `role_resource` VALUES (51, 3);
INSERT INTO `role_resource` VALUES (52, 1);
INSERT INTO `role_resource` VALUES (53, 1);
INSERT INTO `role_resource` VALUES (54, 1);
INSERT INTO `role_resource` VALUES (54, 2);
INSERT INTO `role_resource` VALUES (54, 3);
INSERT INTO `role_resource` VALUES (55, 1);
INSERT INTO `role_resource` VALUES (55, 2);
INSERT INTO `role_resource` VALUES (55, 3);
INSERT INTO `role_resource` VALUES (56, 1);
INSERT INTO `role_resource` VALUES (57, 1);
INSERT INTO `role_resource` VALUES (58, 1);
INSERT INTO `role_resource` VALUES (58, 2);
INSERT INTO `role_resource` VALUES (58, 3);
INSERT INTO `role_resource` VALUES (59, 1);
INSERT INTO `role_resource` VALUES (59, 2);
INSERT INTO `role_resource` VALUES (59, 3);
INSERT INTO `role_resource` VALUES (60, 1);
INSERT INTO `role_resource` VALUES (61, 1);
INSERT INTO `role_resource` VALUES (61, 2);
INSERT INTO `role_resource` VALUES (62, 1);
INSERT INTO `role_resource` VALUES (62, 2);
INSERT INTO `role_resource` VALUES (62, 3);
INSERT INTO `role_resource` VALUES (63, 1);
INSERT INTO `role_resource` VALUES (64, 1);
INSERT INTO `role_resource` VALUES (64, 2);
INSERT INTO `role_resource` VALUES (65, 1);
INSERT INTO `role_resource` VALUES (65, 2);
INSERT INTO `role_resource` VALUES (65, 3);
INSERT INTO `role_resource` VALUES (66, 1);
INSERT INTO `role_resource` VALUES (67, 1);
INSERT INTO `role_resource` VALUES (68, 1);
INSERT INTO `role_resource` VALUES (68, 2);
INSERT INTO `role_resource` VALUES (68, 3);
INSERT INTO `role_resource` VALUES (69, 1);
INSERT INTO `role_resource` VALUES (70, 1);
INSERT INTO `role_resource` VALUES (70, 2);
INSERT INTO `role_resource` VALUES (70, 3);
INSERT INTO `role_resource` VALUES (71, 1);
INSERT INTO `role_resource` VALUES (72, 1);
INSERT INTO `role_resource` VALUES (74, 1);
INSERT INTO `role_resource` VALUES (78, 1);
INSERT INTO `role_resource` VALUES (78, 2);
INSERT INTO `role_resource` VALUES (78, 3);
INSERT INTO `role_resource` VALUES (79, 1);
INSERT INTO `role_resource` VALUES (79, 2);
INSERT INTO `role_resource` VALUES (79, 3);
INSERT INTO `role_resource` VALUES (80, 1);
INSERT INTO `role_resource` VALUES (80, 2);
INSERT INTO `role_resource` VALUES (80, 3);
INSERT INTO `role_resource` VALUES (81, 1);
INSERT INTO `role_resource` VALUES (82, 1);
INSERT INTO `role_resource` VALUES (82, 2);
INSERT INTO `role_resource` VALUES (82, 3);
INSERT INTO `role_resource` VALUES (84, 1);
INSERT INTO `role_resource` VALUES (85, 1);
INSERT INTO `role_resource` VALUES (86, 1);
INSERT INTO `role_resource` VALUES (86, 2);
INSERT INTO `role_resource` VALUES (86, 3);
INSERT INTO `role_resource` VALUES (91, 1);
INSERT INTO `role_resource` VALUES (92, 1);
INSERT INTO `role_resource` VALUES (92, 2);
INSERT INTO `role_resource` VALUES (92, 3);
INSERT INTO `role_resource` VALUES (93, 1);
INSERT INTO `role_resource` VALUES (95, 1);
INSERT INTO `role_resource` VALUES (95, 2);
INSERT INTO `role_resource` VALUES (95, 3);
INSERT INTO `role_resource` VALUES (98, 1);
INSERT INTO `role_resource` VALUES (99, 1);
INSERT INTO `role_resource` VALUES (100, 1);
INSERT INTO `role_resource` VALUES (101, 1);
INSERT INTO `role_resource` VALUES (102, 1);
INSERT INTO `role_resource` VALUES (103, 1);
INSERT INTO `role_resource` VALUES (103, 2);
INSERT INTO `role_resource` VALUES (103, 3);
INSERT INTO `role_resource` VALUES (104, 1);
INSERT INTO `role_resource` VALUES (105, 1);
INSERT INTO `role_resource` VALUES (106, 1);
INSERT INTO `role_resource` VALUES (107, 1);
INSERT INTO `role_resource` VALUES (108, 1);
INSERT INTO `role_resource` VALUES (108, 2);
INSERT INTO `role_resource` VALUES (108, 3);
INSERT INTO `role_resource` VALUES (109, 1);

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
-- Table structure for user_auth_role
-- ----------------------------
DROP TABLE IF EXISTS `user_auth_role`;
CREATE TABLE `user_auth_role`  (
  `user_auth_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`user_auth_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_auth_role
-- ----------------------------
INSERT INTO `user_auth_role` VALUES (2, 1);
INSERT INTO `user_auth_role` VALUES (3, 2);
INSERT INTO `user_auth_role` VALUES (4, 3);

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
