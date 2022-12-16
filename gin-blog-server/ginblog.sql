/*
 Navicat Premium Data Transfer

 Source Server         : docker mysql
 Source Server Type    : MySQL
 Source Server Version : 80031 (8.0.31)
 Source Host           : localhost:3306
 Source Schema         : ginblog

 Target Server Type    : MySQL
 Target Server Version : 80031 (8.0.31)
 File Encoding         : 65001

 Date: 17/12/2022 00:31:03
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS `article`;
CREATE TABLE `article`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文章标题',
  `category_id` bigint NOT NULL COMMENT '分类 ID',
  `desc` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '文章描述',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '文章内容',
  `img` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '封面图片地址',
  `comment_count` bigint NOT NULL DEFAULT 0,
  `read_count` bigint NOT NULL DEFAULT 0,
  `is_top` tinyint NOT NULL DEFAULT 0 COMMENT '是否置顶(0-否 1-是)',
  `original_url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '源链接',
  `type` tinyint NULL DEFAULT NULL COMMENT '类型(1-原创 2-转载 3-翻译)',
  `status` tinyint NULL DEFAULT NULL COMMENT '状态(1-公开 2-私密)',
  `is_delete` tinyint NOT NULL DEFAULT 0 COMMENT '是否放到回收站(0-否 1-是)',
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_article_category`(`category_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of article
-- ----------------------------
INSERT INTO `article` VALUES (1, '2022-12-03 22:07:01.638', '2022-12-16 11:53:37.132', '测试文章', 1, '', '## 博客介绍\n\n<p  align=center>\n<a  href=\"http://www.hahacode.cn\">\n<img  src=\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\"  width=\"200\"  hight=\"200\"  alt=\"阵、雨的个人博客\"  style=\"border-radius: 50%\">\n</a>\n</p>\n\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\n\n您的 Star 是我坚持的动力，感谢大家的支持，欢迎提交 Pr 共同改进项目。\n\nGithub 地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\n\nGitee 地址：[https://gitee.com/szluyu99/gin-vue-blog](https://gitee.com/szluyu99/gin-vue-blog)\n\n## 在线预览\n\n博客前台链接：[hahacode.cn](https://www.hahacode.cn)\n\n博客后台链接：[hahacode.cn/blog-admin](https://www.hahacode.cn/blog-admin)\n\n> 博客已备案通过，且配置 SSL，可通过 https 访问\n\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\n\n在线接口文档地址：[doc.hahacode.cn](http://doc.hahacode.cn/)\n\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\n>\n> 并且由于一开始不会用 Apifox，接口文档生成的信息不全（如返回响应、响应示例），抽空完善一下\n\n## 目录结构\n\n代码仓库目录：\n\n```bash\ngin-vue-blog\n├── gin-blog-admin      -- 博客后台前端\n├── gin-blog-front      -- 博客前台前端\n├── gin-blog-server     -- 博客后端\n```\n\n> 项目运行参考：[快速开始](#快速开始)\n\n后端目录：\n\n```bash\ngin-blog-server\n├── api             -- API\n│   ├── front       -- 前台接口\n│   └── v1          -- 后台接口\n├── dao             -- 数据库操作模块\n├── service         -- 服务模块\n├── model           -- 数据模型\n│   ├── req             -- 请求 VO 模型\n│   ├── resp            -- 响应 VO 模型\n│   ├── dto             -- 内部传输 DTO 模型\n│   └── ...             -- 数据库模型对象 PO 模型\n├── routes          -- 路由模块\n│   └── middleware      -- 路由中间件\n├── utils           -- 工具模块\n│   ├── r               -- 响应封装\n│   ├── upload          -- 文件上传\n│   └── ...\n├── routes          -- 路由模块\n├── config          -- 配置文件\n├── test            -- 测试模块\n├── log             -- 日志文件\n├── Dockerfile\n└── main.go\n```\n\n前端目录：自行参考代码文件\n\n## 项目介绍\n\n市面上有太多优秀的前后台框架，本项目也是参考了很多开源项目，但是大多项目过于重量级（并非坏处），如果从学习的角度来看对初学者并不是很友好。本项目在以**博客**这个业务为主的前提下，提供一个完整的全栈项目代码（前台前端 + 后台前端 + 后端），技术点基本都是最新 + 最火的技术，代码轻量级，注释完善，适合学习。\n\n前台：\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\n- 响应式布局，适配了移动端\n- 实现点赞，统计用户等功能 (Redis)\n- 评论 + 回复评论功能\n- 留言采用弹幕墙，效果炫酷\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\n\n后台：\n- 鉴权使用 JWT\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\n- 支持动态权限修改，前端菜单由后端生成（动态路由）\n- 文章编辑使用 Markdown 编辑器\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\n- 实现记录操作日志功能（GET 不记录）\n- 实现监听在线用户、强制下线功能\n- 文件上传支持七牛云、本地（后续计划支持更多）\n- 对 CRUD 操作封装了通用 Hook\n\n其他：\n- 采用 Restful 风格的 API\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\n- 代码整洁层次清晰，利于开发者学习\n- 技术点新颖，代码轻量级，适度封装\n\n## 技术介绍\n\n> 这里只写一些主流的通用技术，详细第三方库：前端参考 `package.json` 文件，后端参考 `go.mod` 文件\n\n前台前端：使用 pnpm 包管理工具\n- 基于 TypeScript\n- Vue3\n- VueUse: 服务于 Vue Composition API 的工具集\n- Unocss: 原子化 CSS\n- Pinia\n- Vue Router \n- Axios \n- Naive UI\n- ...\n\n后台前端：使用 pnpm 包管理工具\n- 基于 JavaSciprt \n- Vue3\n- VueUse: 服务于 Vue Composition API 的工具集\n- Unocss: 原子化 CSS\n- Pinia \n- Vue Router \n- Axios \n- Naive UI\n- ...\n\n后端技术栈：\n- 基于 Golang\n- Docker\n- Gin\n- GORM\n- Viper: 使用 TOML 作为配置文件\n- Casbin\n- Zap\n- MySQL\n- Redis\n- Nginx: 部署静态资源 + 反向代理\n- ...\n\n其他：\n- 腾讯云人机验证\n- 七牛云对象存储\n- ...\n\n## 运行环境\n\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\n\n对象存储：七牛云\n\n## 开发环境\n\n| 开发工具                      | 说明                    |\n| ----------------------------- | ----------------------- |\n| Vscode                        | Golang 后端 +  Vue 前端 |\n| Navicat                       | MySQL 远程连接工具      |\n| Another Redis Desktop Manager | Redis 远程连接工具      |\n| MobaXterm                     | Linux 远程工具          |\n| Apifox                        | 接口调试 + 文档生成     |\n\n| 开发环境 | 版本 |\n| -------- | ---- |\n| Golang   | 1.19 |\n| MySQL    | 8.x  |\n| Redis    | 7.x  |\n\n\n### Vscode 插件\n\n如果使用 Vscode 开发，推荐安装一下以下插件。\n\n> 前后端相关插件属于**开发必须插件**，还有很多提升开发效果的插件后面推荐一下\n\n前端开发插件：\n\n| 插件 | 作用 |\n| -------- | ---- |\n| Volar   | Vue 官方插件 |\n| UnoCSS | Unocss 官方插件 |\n| Tailwind CSS IntelliSense | Tailwind 官方插件 |\n| Iconify IntelliSense | Iconify 提示插件 |\n\n后端开发插件：\n\n| 插件 | 作用 |\n| -------- | ---- |\n| Go | Golang 官方插件 |\n\n\n通用插件：（非开发必须，个人推荐）\n\n| 名称 | 作用 |\n| -------- | ---- |\n| Better Comments   | 优化代码注释显示效果 |\n| TODO Highlight | 高亮 TODO |\n| Highlight Matching Tag | 高亮匹配的标签 | \n\n\n## 快速开始\n\n### 本地运行\n\n> 目前需要自行安装 Golang、Node、MySQL、Redis 环境\n>\n> TODO: 完全基于 Docker 的运行教程（Docker 实在太好用了！）\n\n\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\n\n拉取项目到本地：\n\n```bash\ngit clone https://github.com/szluyu99/gin-vue-blog.git\n```\n\n后端项目运行：\n\n```bash\n# 1、进入后端项目根目录 \ncd gin-blog-server\n\n# 2、修改项目运行的配置文件，默认加载位于 config/config.toml \n\n# 3、MySQL 导入 ginblog.sql\n\n# 4、启动 Redis \n\n# 5、运行项目\ngo mod tidy\ngo run main.go\n```\n\n数据库中的默认用户：\n- 管理员 admin 123456\n- 普通用户 user 123456\n- 测试用户 test 123456\n\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\n\n```bash\nnpm install -g pnpm\n```\n\n前台前端：\n\n```bash\n# 1、进入前台前端项目根目录\ncd gin-blog-front\n\n# 2、安装依赖\npnpm install\n\n# 3、运行项目\npnpm run dev\n```\n\n后台前端：\n\n```bash\n# 1、进入后台前端项目根目录\ncd gin-blog-admin\n\n# 2、安装依赖\npnpm install\n\n# 3、运行项目\npnpm run dev\n```\n\n\n### 项目部署\n\n**目前暂时不推荐将本博客部署上生产环境，因为还有很多功能未完善**。\n\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\n\n等功能开发的差不多了，再专门针对部署写一篇文章。\n\n---\n\n这里简单介绍一下，有基础的同学可以自行折腾。\n\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\n\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\n\nDocker 打包成镜像指令：\n\n```bash\ndocker build -t ginblog .\n```\n\n> 以上只是简单说明，等功能全部完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\n\n## 项目总结\n\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\n\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\n\n## 鸣谢项目\n\n- [https://butterfly.js.org/](https://butterfly.js.org/)\n- [https://github.com/flipped-aurora/gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)\n- [https://github.com/qifengzhang007/GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)\n- [https://github.com/X1192176811/blog](https://github.com/X1192176811/blog)\n- [https://github.com/zclzone/vue-naive-admin](https://github.com/zclzone/vue-naive-admin)\n- [https://github.com/antfu/vitesse](https://github.com/antfu/vitesse)\n- ...\n\n> 需要感谢的绝不止以上这些开源项目，但是一时难以全部列出，后面会慢慢补上。\n\n## 后续计划\n\n高优先级: \n\n- ~~完善图片上传功能, 目前文件上传还没怎么处理~~\n- 后台首页重新设计（目前没放什么内容）\n- ~~前台首页搜索文章（目前使用数据库模糊搜索）~~\n- 博客文章导入导出 (.md 文件)\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\n- 后端日志切割\n- ~~后台修改背景图片，博客配置等~~ \n- 后端的 IP 地址检测 BUG 待修复\n- ~~博客前台适配移动端~~ \n- ~~文章详情, 目录锚点跟随~~ \n\n后续有空安排上：\n- 黑夜模式\n- 前台收缩侧边信息功能\n- 说说\n- 相册\n- 音乐播放器\n- 鼠标左击特效\n- 看板娘\n- 第三方登录\n- 评论时支持选择表情，参考 Valine\n- 一键部署：使用 docker compose 单机一键部署整个项目（前后端 + 环境）\n- 单独部署：前后端 + 环境\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\n- 前台首页搜索集成 ElasticSearch\n\n其他：\n- 写一份好的文档\n- 补全 README.md\n- 完善 Apifox 生成的接口文档\n\n## 更新日志\n\n描述规范定义，有以下几种行为 ACTION:\n- `FIX`: 修复\n- `REFINE`: 优化\n- `COMPLETE`: 完善\n- `ADD`: 新增\n\n---\n\n2022-12-15:\n\n博客后台：\n- 优化 菜单栏和标签栏，点击标签自动展开对应菜单，点击菜单自动滚动到显示对应标签 ⭐\n- 优化 使用响应式语法糖重构代码\n- 优化 代码结构 + 注释\n\n---\n\n2022-12-14:\n\n博客前台：\n- 优化 代码，去除 n-card 组件，使用自定义 css 实现卡片视图\n- 优化 对滚动的监听操作，添加节流函数，提升性能 ⭐\n- 优化 文章目录，根据当前滚动条自动高亮锚点 ⭐⭐\n- 新增 vite 打包优化相关插件\n\n---\n\n2022-12-13:\n\n项目部署：\n- 新增 Nginx 配置 https 访问域名 (http 自动跳转到 https) ⭐\n- 新增 七牛云添加加速域名访问图片资源\n\n博客后端：\n- 新增 文章搜索接口（数据库模糊查询）\n\n博客前台：\n- 新增 导航栏的文章搜索\n\n---\n\n2022-12-12:\n\n博客前台：\n- 新增 适配移动端 ⭐⭐\n- 优化 删除 Vuetify 相关组件及依赖 (css 样式存在冲突)，统一使用 naive-ui\n- 优化 使用 `$ref` 语法糖功能重构页面\n\n---\n\n2022-12-09:\n\n博客前台：\n- 完善 个人中心的头像上传（注意头像上传需要 Token）\n- 优化 重构通用页面的代码 ⭐\n\n---\n\n2022-12-08:\n\n博客后端：\n- 新增 页面模块 的接口 ⭐\n- 修复 单元测试无法初始化全局变量（flag 包冲突）问题\n\n博客后台：\n- 新增 pinia 数据持久化，防止刷新丢失数据 ⭐\n- 新增 页面管理页面，动态配置前台显示的封面 ⭐\n\n博客前台：\n- 新增 页面根据 `label` 从后端数据中加载封面 ⭐\n- 优化 尝试性使用 `$ref` 语法糖功能\n\n---\n\n2022-12-07:\n\n博客后台：\n- 修复 博客后台，动态加载路由模块导致的热加载失效问题 ⭐⭐\n- 完善 文件上传，抽取出单独的组件 ⭐⭐\n- 新增 个人信息页面 ⭐\n- 优化 发布文章/查看文章 时的文章标签选择、分类选择\n- 修复 发布文章/查看文章 页面数据加载错误', 'https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg', 0, 0, 0, '', 1, 1, 0, 1);
INSERT INTO `article` VALUES (2, '2022-12-16 11:56:40.880', '2022-12-16 12:36:42.106', '项目运行成功', 2, '', '恭喜你，项目已经成功运行起来了！\n\n```go\nfmt.Println(\"恭喜！\")\n```\n\n```js\nconsole.log(\"恭喜！\")\n```\n\n🆗😋', 'https://static.talkxj.com/articles/1395642324821741569.png', 0, 0, 0, '', 1, 1, 0, 1);

-- ----------------------------
-- Table structure for article_tag
-- ----------------------------
DROP TABLE IF EXISTS `article_tag`;
CREATE TABLE `article_tag`  (
  `article_id` bigint UNSIGNED NOT NULL,
  `tag_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`article_id`, `tag_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of article_tag
-- ----------------------------
INSERT INTO `article_tag` VALUES (1, 1);
INSERT INTO `article_tag` VALUES (1, 2);
INSERT INTO `article_tag` VALUES (2, 1);
INSERT INTO `article_tag` VALUES (2, 2);

-- ----------------------------
-- Table structure for blog_config
-- ----------------------------
DROP TABLE IF EXISTS `blog_config`;
CREATE TABLE `blog_config`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `config` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置信息',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of blog_config
-- ----------------------------
INSERT INTO `blog_config` VALUES (1, '2022-10-30 14:57:38.000', '2022-12-17 00:02:41.127', '{\"website_avatar\":\"https://img-blog.csdnimg.cn/9e1832a8817d4344901f2e476cc74c16.jpeg\",\"website_name\":\"阵、雨的个人博客\",\"website_author\":\"阵、雨\",\"website_intro\":\"往事随风而去\",\"website_notice\":\"博客后端基于 gin、gorm 开发\\n博客前端基于 Vue3、TS、NaiveUI 开发\\n努力开发中...冲冲冲！加油！\",\"website_createtime\":\"2022-11-01\",\"website_record\":\"鲁ICP备2022040119号\",\"social_login_list\":[],\"social_url_list\":[],\"qq\":\"123456\",\"github\":\"https://github.com/szluyu99\",\"gitee\":\"https://gitee.com/szluyu99\",\"tourist_avatar\":\"https://static.talkxj.com/photos/0bca52afdb2b9998132355d716390c9f.png\",\"user_avatar\":\"https://static.talkxj.com/avatar/user.png\",\"article_cover\":\"https://static.talkxj.com/config/e587c4651154e4da49b5a54f865baaed.jpg\",\"is_comment_review\":1,\"is_message_review\":1,\"is_email_notice\":0,\"is_reward\":0,\"wechat_qrcode\":\"http://dummyimage.com/100x100\",\"alipay_ode\":\"http://dummyimage.com/100x100\"}');

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_casbin_rule`(`ptype` ASC, `v0` ASC, `v1` ASC, `v2` ASC, `v3` ASC, `v4` ASC, `v5` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2439 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO `casbin_rule` VALUES (1416, 'g', 'admin', 'anonymous', NULL, NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (1418, 'g', 'admin', 'logout', NULL, NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (1420, 'g', 'test', 'anonymous', NULL, NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (1421, 'g', 'test', 'logout', NULL, NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (1417, 'g', 'user', 'anonymous', NULL, NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (1419, 'g', 'user', 'logout', NULL, NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (2397, 'p', 'admin', '/article', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2395, 'p', 'admin', '/article', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2394, 'p', 'admin', '/article/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2393, 'p', 'admin', '/article/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2396, 'p', 'admin', '/article/soft-delete', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2398, 'p', 'admin', '/article/top', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2401, 'p', 'admin', '/category', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2400, 'p', 'admin', '/category', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2399, 'p', 'admin', '/category/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2402, 'p', 'admin', '/category/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2411, 'p', 'admin', '/comment', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2410, 'p', 'admin', '/comment/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2412, 'p', 'admin', '/comment/review', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2430, 'p', 'admin', '/home', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2419, 'p', 'admin', '/link', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2418, 'p', 'admin', '/link', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2417, 'p', 'admin', '/link/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2389, 'p', 'admin', '/menu', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2390, 'p', 'admin', '/menu/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2388, 'p', 'admin', '/menu/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2391, 'p', 'admin', '/menu/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2392, 'p', 'admin', '/menu/user/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2408, 'p', 'admin', '/message', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2407, 'p', 'admin', '/message/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2409, 'p', 'admin', '/message/review', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2429, 'p', 'admin', '/operation/log', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2428, 'p', 'admin', '/operation/log/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2437, 'p', 'admin', '/page', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2436, 'p', 'admin', '/page', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2435, 'p', 'admin', '/page/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2384, 'p', 'admin', '/resource', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2387, 'p', 'admin', '/resource/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2383, 'p', 'admin', '/resource/anonymous', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2385, 'p', 'admin', '/resource/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2386, 'p', 'admin', '/resource/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2415, 'p', 'admin', '/role', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2414, 'p', 'admin', '/role', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2413, 'p', 'admin', '/role/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2416, 'p', 'admin', '/role/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2422, 'p', 'admin', '/setting/about', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2426, 'p', 'admin', '/setting/about', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2421, 'p', 'admin', '/setting/blog-config', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2423, 'p', 'admin', '/setting/blog-config', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2405, 'p', 'admin', '/tag', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2404, 'p', 'admin', '/tag', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2403, 'p', 'admin', '/tag/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2406, 'p', 'admin', '/tag/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2438, 'p', 'admin', '/upload', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (2425, 'p', 'admin', '/user', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2433, 'p', 'admin', '/user/current', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2432, 'p', 'admin', '/user/current/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2434, 'p', 'admin', '/user/disable', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (2424, 'p', 'admin', '/user/info', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2420, 'p', 'admin', '/user/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2431, 'p', 'admin', '/user/offline', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (2427, 'p', 'admin', '/user/online', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (1488, 'p', 'logout', '/logout', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (2365, 'p', 'test', '/article/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2364, 'p', 'test', '/article/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2366, 'p', 'test', '/category/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2367, 'p', 'test', '/category/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2371, 'p', 'test', '/comment/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2381, 'p', 'test', '/home', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2374, 'p', 'test', '/link/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2361, 'p', 'test', '/menu/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2362, 'p', 'test', '/menu/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2363, 'p', 'test', '/menu/user/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2370, 'p', 'test', '/message/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2380, 'p', 'test', '/operation/log/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2382, 'p', 'test', '/page/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2359, 'p', 'test', '/resource/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2360, 'p', 'test', '/resource/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2372, 'p', 'test', '/role/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2373, 'p', 'test', '/role/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2377, 'p', 'test', '/setting/about', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2376, 'p', 'test', '/setting/blog-config', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2368, 'p', 'test', '/tag/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2369, 'p', 'test', '/tag/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2378, 'p', 'test', '/user/info', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2375, 'p', 'test', '/user/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2379, 'p', 'test', '/user/online', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2341, 'p', 'user', '/article/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2340, 'p', 'user', '/article/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2342, 'p', 'user', '/category/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2343, 'p', 'user', '/category/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2347, 'p', 'user', '/comment/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2357, 'p', 'user', '/home', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2350, 'p', 'user', '/link/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2337, 'p', 'user', '/menu/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2338, 'p', 'user', '/menu/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2339, 'p', 'user', '/menu/user/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2346, 'p', 'user', '/message/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2356, 'p', 'user', '/operation/log/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2358, 'p', 'user', '/page/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2335, 'p', 'user', '/resource/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2336, 'p', 'user', '/resource/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2348, 'p', 'user', '/role/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2349, 'p', 'user', '/role/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2353, 'p', 'user', '/setting/about', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2352, 'p', 'user', '/setting/blog-config', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2344, 'p', 'user', '/tag/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2345, 'p', 'user', '/tag/option', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2354, 'p', 'user', '/user/info', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2351, 'p', 'user', '/user/list', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2355, 'p', 'user', '/user/online', 'GET', '', '', '');

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of category
-- ----------------------------
INSERT INTO `category` VALUES (1, '后端', '2022-12-03 22:01:29.106', '2022-12-03 22:01:29.106');
INSERT INTO `category` VALUES (2, '前端', '2022-12-03 22:01:33.074', '2022-12-03 22:01:33.074');

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `user_id` bigint NULL DEFAULT NULL COMMENT '评论用户id',
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论内容',
  `reply_user_id` bigint NULL DEFAULT NULL COMMENT '回复用户id',
  `topic_id` bigint NULL DEFAULT NULL COMMENT '评论主题id',
  `parent_id` bigint NULL DEFAULT NULL COMMENT '父评论id',
  `type` tinyint(1) NOT NULL COMMENT '评论类型(1.文章 2.友链 3.说说)',
  `is_delete` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除(0.否 1.是)',
  `is_review` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否审核(0.否 1.是)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of comment
-- ----------------------------
INSERT INTO `comment` VALUES (1, '2022-12-03 22:36:04.549', '2022-12-03 22:36:04.549', 3, '评论测试！', 0, 1, 0, 1, 0, 1);

-- ----------------------------
-- Table structure for friend_link
-- ----------------------------
DROP TABLE IF EXISTS `friend_link`;
CREATE TABLE `friend_link`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

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
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '菜单名',
  `path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '菜单路径',
  `component` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '组件',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '菜单图标',
  `order_num` tinyint NULL DEFAULT 0 COMMENT '显示排序',
  `parent_id` bigint NULL DEFAULT NULL COMMENT '父菜单id',
  `is_hidden` tinyint(1) NULL DEFAULT 0 COMMENT '是否隐藏(0-否 1-是)',
  `keep_alive` tinyint(1) NULL DEFAULT 1,
  `redirect` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '跳转路径',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 40 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (2, '2022-10-31 09:41:03.000', '2022-11-01 01:14:01.229', '文章管理', '/article', 'Layout', 'ic:twotone-article', 1, 0, 0, 1, '/article/list');
INSERT INTO `menu` VALUES (3, '2022-10-31 09:41:03.000', '2022-11-01 10:11:07.042', '消息管理', '/message', 'Layout', 'ic:twotone-email', 2, 0, 0, 1, '/message/comment	');
INSERT INTO `menu` VALUES (4, '2022-10-31 09:41:03.000', '2022-11-01 10:11:31.573', '用户管理', '/user', 'Layout', 'ph:user-list-bold', 4, 0, 0, 1, '/user/list');
INSERT INTO `menu` VALUES (5, '2022-10-31 09:41:03.000', '2022-11-01 10:11:45.624', '系统管理', '/setting', 'Layout', 'ion:md-settings', 5, 0, 0, 1, '/setting/website');
INSERT INTO `menu` VALUES (6, '2022-10-31 09:41:03.000', '2022-11-05 14:21:38.109', '发布文章', 'write', '/article/write', 'icon-park-outline:write', 1, 2, 0, 1, '');
INSERT INTO `menu` VALUES (8, '2022-10-31 09:41:03.000', '2022-11-01 01:18:25.685', '文章列表', 'list', '/article/list', 'material-symbols:format-list-bulleted', 2, 2, 0, 1, '');
INSERT INTO `menu` VALUES (9, '2022-10-31 09:41:03.000', '2022-11-01 01:18:30.931', '分类管理', 'category', '/article/category', 'tabler:category', 3, 2, 0, 1, '');
INSERT INTO `menu` VALUES (10, '2022-10-31 09:41:03.000', '2022-11-01 01:18:35.502', '标签管理', 'tag', '/article/tag', 'tabler:tag', 4, 2, 0, 1, '');
INSERT INTO `menu` VALUES (16, '2022-10-31 09:41:03.000', '2022-11-01 10:11:23.195', '权限管理', '/auth', 'Layout', 'cib:adguard', 3, 0, 0, 1, '/auth/menu');
INSERT INTO `menu` VALUES (17, '2022-10-31 09:41:03.000', NULL, '菜单管理', 'menu', '/auth/menu', 'ic:twotone-menu-book', 1, 16, 0, 1, NULL);
INSERT INTO `menu` VALUES (23, '2022-10-31 09:41:03.000', NULL, '接口管理', 'resource', '/auth/resource', 'mdi:api', 2, 16, 0, 1, NULL);
INSERT INTO `menu` VALUES (24, '2022-10-31 09:41:03.000', '2022-10-31 10:09:18.913', '角色管理', 'role', '/auth/role', 'carbon:user-role', 3, 16, 0, 1, NULL);
INSERT INTO `menu` VALUES (25, '2022-10-31 10:11:09.232', '2022-11-01 01:29:48.520', '评论管理', 'comment', '/message/comment', 'ic:twotone-comment', 1, 3, 0, 1, '');
INSERT INTO `menu` VALUES (26, '2022-10-31 10:12:01.546', '2022-11-01 01:29:54.130', '留言管理', 'leave-msg', '/message/leave-msg', 'ic:twotone-message', 2, 3, 0, 1, '');
INSERT INTO `menu` VALUES (27, '2022-10-31 10:54:03.201', '2022-11-01 01:30:06.901', '用户列表', 'list', '/user/list', 'mdi:account', 1, 4, 0, 1, '');
INSERT INTO `menu` VALUES (28, '2022-10-31 10:54:34.167', '2022-11-01 01:30:13.400', '在线用户', 'online', '/user/online', 'ic:outline-online-prediction', 2, 4, 0, 1, '');
INSERT INTO `menu` VALUES (29, '2022-10-31 10:59:33.255', '2022-11-01 01:30:20.688', '网站管理', 'website', '/setting/website', 'el:website', 1, 5, 0, 1, '');
INSERT INTO `menu` VALUES (30, '2022-10-31 11:00:09.997', '2022-11-01 01:30:24.097', '页面管理', 'page', '/setting/page', 'iconoir:journal-page', 2, 5, 0, 1, '');
INSERT INTO `menu` VALUES (31, '2022-10-31 11:00:33.543', '2022-11-01 01:30:28.497', '友链管理', 'link', '/setting/link', 'mdi:telegram', 3, 5, 0, 1, '');
INSERT INTO `menu` VALUES (32, '2022-10-31 11:01:00.444', '2022-11-01 01:30:33.186', '关于我', 'about', '/setting/about', 'cib:about-me', 4, 5, 0, 1, '');
INSERT INTO `menu` VALUES (33, '2022-11-01 01:43:10.142', '2022-12-07 20:53:27.473', '首页', '/home', '/home', 'ic:sharp-home', 0, 0, 0, 1, '');
INSERT INTO `menu` VALUES (34, '2022-11-01 09:54:36.252', '2022-11-01 10:07:00.254', '修改文章', 'write/:id', '/article/write', 'icon-park-outline:write', 1, 2, 1, 1, '');
INSERT INTO `menu` VALUES (36, '2022-11-04 15:50:45.993', '2022-11-04 15:58:14.443', '日志管理', '/log', 'Layout', 'material-symbols:receipt-long-outline-rounded', 6, 0, 0, 1, '/log/operation');
INSERT INTO `menu` VALUES (37, '2022-11-04 15:53:00.251', '2022-11-04 15:58:36.560', '操作日志', 'operation', '/log/operation', 'mdi:book-open-page-variant-outline', 1, 36, 0, 1, '');
INSERT INTO `menu` VALUES (38, '2022-11-04 16:02:42.306', '2022-11-04 16:05:35.761', '登录日志', 'login', '/log/login', 'material-symbols:login', 2, 36, 0, 1, '');
INSERT INTO `menu` VALUES (39, '2022-12-07 20:47:08.349', '2022-12-07 20:53:33.851', '个人中心', '/profile', '/profile', 'mdi:account', 7, 0, 0, 1, '');

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '头像地址',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '留言内容',
  `ip_address` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'IP 地址',
  `ip_source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'IP 来源',
  `is_review` tinyint(1) NOT NULL DEFAULT 0 COMMENT '审核状态(0-未审核,1-审核通过)',
  `speed` tinyint(1) NULL DEFAULT NULL COMMENT '弹幕速度',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of message
-- ----------------------------
INSERT INTO `message` VALUES (1, '2022-12-03 22:36:43.299', '2022-12-03 22:36:43.299', '测试用户', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '留言测试！', '[::1]:53442', '未知', 1, 0);

-- ----------------------------
-- Table structure for operation_log
-- ----------------------------
DROP TABLE IF EXISTS `operation_log`;
CREATE TABLE `operation_log`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `opt_module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作模块',
  `opt_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作类型',
  `opt_method` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作方法',
  `opt_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作URL',
  `opt_desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作描述',
  `request_param` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '请求参数',
  `request_method` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '请求方法',
  `response_data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '响应数据',
  `user_id` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `nickname` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '用户昵称',
  `ip_address` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '操作IP',
  `ip_source` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '操作地址',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 54 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of operation_log
-- ----------------------------
INSERT INTO `operation_log` VALUES (1, '2022-12-03 21:53:48.670', '2022-12-03 21:53:48.670', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":1,\"name\":\"管理员\",\"label\":\"admin\",\"created_at\":\"2022-10-20T21:24:28+08:00\",\"is_disable\":0,\"resource_ids\":[3,43,44,45,46,47,48,6,59,60,61,7,38,39,40,41,42,8,65,66,67,68,9,62,63,64,10,23,34,35,36,37,79,80,81,85,49,51,52,53,54,50,55,56,57,58,69,70,71,72,91,92,93,74,78,82,84,86,98,95,11],\"menu_ids\":[2,3,4,5,6,7,8,9,16,17,23,24,25,26,27,28,29,30,31,32,33,36,37,38,34,10]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (2, '2022-12-03 22:01:29.112', '2022-12-03 22:01:29.112', '分类', '新增或修改', 'gin-blog/api/v1.(*Category).SaveOrUpdate-fm', '/api/category', '新增或修改分类', '{\"name\":\"后端\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (3, '2022-12-03 22:01:33.080', '2022-12-03 22:01:33.080', '分类', '新增或修改', 'gin-blog/api/v1.(*Category).SaveOrUpdate-fm', '/api/category', '新增或修改分类', '{\"name\":\"前端\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (4, '2022-12-03 22:01:51.630', '2022-12-03 22:01:51.630', '标签', '新增或修改', 'gin-blog/api/v1.(*Tag).SaveOrUpdate-fm', '/api/tag', '新增或修改标签', '{\"name\":\"Golang\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (5, '2022-12-03 22:01:56.990', '2022-12-03 22:01:56.990', '标签', '新增或修改', 'gin-blog/api/v1.(*Tag).SaveOrUpdate-fm', '/api/tag', '新增或修改标签', '{\"name\":\"Vue\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (6, '2022-12-03 22:07:01.656', '2022-12-03 22:07:01.656', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"status\":1,\"title\":\"测试文章\",\"content\":\"这是一篇测试文章\",\"category_name\":\"后端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"img\":\"\",\"is_top\":0}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (7, '2022-12-03 22:19:51.943', '2022-12-03 22:19:51.943', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"这是一篇测试文章\",\"img\":\"public/uploaded//1cba77c39b4d0a81024a7aada3655a28_20221203221949.jpg\",\"category_name\":\"后端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"original_url\":\"\",\"is_top\":0,\"status\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (8, '2022-12-04 00:38:20.583', '2022-12-04 00:38:20.583', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"## 博客介绍\\n\\n<p  align=center>\\n<a  href=\\\"http://www.hahacode.cn\\\">\\n<img  src=\\\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\\\"  width=\\\"200\\\"  hight=\\\"200\\\"  alt=\\\"阵、雨的个人博客\\\"  style=\\\"border-radius: 50%\\\">\\n</a>\\n</p>\\n\\n<p align=\\\"center\\\">\\n   <a target=\\\"_blank\\\" href=\\\"#\\\">\\n      <img src=\\\"https://img.shields.io/badge/Go-1.19-blue\\\"/>\\n      <img src=\\\"https://img.shields.io/badge/Gin-v1.8.1-blue\\\"/>\\n      <img src=\\\"https://img.shields.io/badge/Casbin-v2.56.0-blue\\\"/>\\n      <img src=\\\"https://img.shields.io/badge/mysql-8.0-blue\\\"/>\\n      <img src=\\\"https://img.shields.io/badge/GORM-v1.24.0-blue\\\"/>\\n      <img src=\\\"https://img.shields.io/badge/redis-7.0-red\\\"/>\\n      <img src=\\\"https://img.shields.io/badge/vue-v3.X-green\\\"/>\\n    </a>\\n</p>\\n\\n[在线地址](#在线地址) | [目录结构](#目录结构) | [技术介绍](#技术介绍) | [运行环境](#运行环境) | [开发环境](#开发环境) | [快速开始](#快速开始) | [项目总结](#项目总结) \\n\\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\\n\\n您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。\\n\\n\\nGithub地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\nGitee地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\n## 在线地址\\n\\n博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)\\n\\n博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)\\n\\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\\n\\n在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)\\n\\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\\n\\n## 目录结构\\n\\n代码仓库目录：\\n\\n```bash\\ngin-vue-blog\\n├── gin-blog-admin      -- 博客后台前端\\n├── gin-blog-front      -- 博客前台前端\\n├── gin-blog-server     -- 博客后端\\n```\\n\\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\\n\\n> 项目运行参考：[快速开始](#快速开始)\\n\\n后端目录：\\n\\n```bash\\ngin-blog-server\\n├── api             -- API\\n│   ├── front       -- 前台接口\\n│   └── v1          -- 后台接口\\n├── dao             -- 数据库操作模块\\n├── service         -- 服务模块\\n├── model           -- 数据模型\\n│   ├── req             -- 请求 VO 模型\\n│   ├── resp            -- 响应 VO 模型\\n│   ├── dto             -- 内部传输 DTO 模型\\n│   └── ...             -- 数据库模型对象 PO 模型\\n├── routes          -- 路由模块\\n│   └── middleware      -- 路由中间件\\n├── utils           -- 工具模块\\n│   ├── r               -- 响应封装\\n│   ├── upload          -- 文件上传\\n│   └── ...\\n├── routes          -- 路由模块\\n├── config          -- 配置文件\\n├── test            -- 测试模块\\n├── log             -- 日志文件\\n├── Dockerfile\\n└── main.go\\n```\\n\\n## 项目介绍\\n\\n前台：\\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\\n- 实现点赞，统计用户等功能 (Redis)\\n- 评论 + 回复评论功能\\n- 留言采用弹幕墙，效果炫酷\\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\\n\\n后台：\\n- 鉴权使用 JWT\\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\\n- 支持动态权限修改，前端菜单由后端生成\\n- 文章编辑使用 Markdown 编辑器\\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\\n\\n其他：\\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\\n- 代码整洁层次清晰，利于开发者学习\\n\\n## 技术介绍\\n\\n> 这里只写一些主流的通用技术，详细第三方库: 前端参考 `package.json` 文件，后端参考 `go.mod` 文件\\n\\n前台前端：使用 pnpm 包管理工具\\n- 基于 TypeScript\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia\\n- Vue Router \\n- Axios \\n- Naive UI\\n- Vuetify\\n- ...\\n\\n后台前端：使用 pnpm 包管理工具\\n- 基于 JavaSciprt \\n- pnpm: 包管理工具\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia \\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后端技术栈：\\n- 基于 Golang\\n- Docker\\n- Gin\\n- GORM\\n- Viper: 使用 TOML 作为配置文件\\n- Zap\\n- MySQL\\n- Redis\\n- Nginx: 部署静态资源 + 反向代理\\n- ...\\n\\n其他：\\n- 腾讯云人机验证\\n- 七牛云对象存储\\n- ...\\n\\n## 运行环境\\n\\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\\n\\n对象存储：七牛云\\n\\n## 开发环境\\n\\n| 开发工具                      | 说明                    |\\n| ----------------------------- | ----------------------- |\\n| Vscode                        | Golang 后端 +  Vue 前端 |\\n| Navicat                       | MySQL 远程连接工具      |\\n| Another Redis Desktop Manager | Redis 远程连接工具      |\\n| MobaXterm                     | Linux 远程工具          |\\n| Apifox                        | 接口调试 + 文档生成     |\\n\\n\\n| 开发环境 | 版本 |\\n| -------- | ---- |\\n| Golang   | 1.19 |\\n| MySQL    | 8.x  |\\n| Redis    | 7.x  |\\n\\n## 快速开始\\n\\n### 本地运行\\n\\n> 自行安装 Golang、Node、MySQL、Redis 环境\\n\\n拉取项目到本地：\\n\\n```bash\\ngit clone https://github.com/szluyu99/gin-vue-blog.git\\n```\\n\\n后端项目运行：\\n\\n```bash\\n# 1、进入后端项目根目录 \\ncd gin-blog-server\\n\\n# 2、修改项目运行的配置文件，默认加载 config/config.toml \\n\\n# 3、MySQL 导入 ginblog.sql\\n\\n# 4、启动 Redis \\n\\n# 5、运行项目\\ngo mod tidy\\ngo run main.go\\n```\\n\\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\\n\\n```bash\\nnpm install -g pnpm\\n```\\n\\n前台前端：\\n\\n```bash\\n# 1、进入前台前端项目根目录\\ncd gin-blog-front\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n后台前端：\\n\\n```bash\\n# 1、进入后台前端项目根目录\\ncd gin-blog-admin\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n### 项目部署\\n\\n目前暂时不推荐将本博客部署上生产环境，因为还有太多功能未完善。\\n\\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\\n\\n等功能开发的差不多了，再专门针对部署写一篇文章。\\n\\n---\\n\\n这里简单介绍一下，有基础的同学可以自行折腾。\\n\\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\\n\\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\\n\\nDocker 打包成镜像指令：\\n\\n```bash\\ndocker build -t ginblog .\\n```\\n\\n> 以上只是简单说明，等功能大致完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\\n\\n## 项目总结\\n\\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\\n\\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\\n\\n## 后续计划\\n\\n高优先级: \\n\\n- 完善图片上传功能, 目前文件上传还没怎么处理\\n- 后台首页重新设计（目前没放什么内容）\\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\\n- 前台首页搜索文章（ElasticSearch 搜索）\\n- 博客文章导入导出 (.md 文件)\\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\\n- 后端日志切割\\n- 后台修改背景图片，博客配置等\\n- 相册\\n\\n后续有空安排上：\\n- 适配移动端\\n- 黑夜模式\\n- 前台收缩侧边信息功能\\n- 说说\\n- 音乐播放器\\n- 鼠标左击特效\\n- 看板娘\\n- 文章目录锚点跟随\\n- 第三方登录\\n- 评论时支持选择表情，参考 Valine\\n- 若干细节需要完善...\\n\",\"img\":\"https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg\",\"category_name\":\"后端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"original_url\":\"\",\"is_top\":0,\"status\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (9, '2022-12-04 00:45:54.335', '2022-12-04 00:45:54.335', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"## 博客介绍\\n\\n<p  align=center>\\n<a  href=\\\"http://www.hahacode.cn\\\">\\n<img  src=\\\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\\\"  width=\\\"200\\\"  hight=\\\"200\\\"  alt=\\\"阵、雨的个人博客\\\"  style=\\\"border-radius: 50%\\\">\\n</a>\\n</p>\\n\\n[在线地址](#在线地址) | [目录结构](#目录结构) | [技术介绍](#技术介绍) | [运行环境](#运行环境) | [开发环境](#开发环境) | [快速开始](#快速开始) | [项目总结](#项目总结) \\n\\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\\n\\n您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。\\n\\n\\nGithub地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\nGitee地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\n## 在线地址\\n\\n博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)\\n\\n博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)\\n\\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\\n\\n在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)\\n\\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\\n\\n## 目录结构\\n\\n代码仓库目录：\\n\\n```bash\\ngin-vue-blog\\n├── gin-blog-admin      -- 博客后台前端\\n├── gin-blog-front      -- 博客前台前端\\n├── gin-blog-server     -- 博客后端\\n```\\n\\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\\n\\n> 项目运行参考：[快速开始](#快速开始)\\n\\n后端目录：\\n\\n```bash\\ngin-blog-server\\n├── api             -- API\\n│   ├── front       -- 前台接口\\n│   └── v1          -- 后台接口\\n├── dao             -- 数据库操作模块\\n├── service         -- 服务模块\\n├── model           -- 数据模型\\n│   ├── req             -- 请求 VO 模型\\n│   ├── resp            -- 响应 VO 模型\\n│   ├── dto             -- 内部传输 DTO 模型\\n│   └── ...             -- 数据库模型对象 PO 模型\\n├── routes          -- 路由模块\\n│   └── middleware      -- 路由中间件\\n├── utils           -- 工具模块\\n│   ├── r               -- 响应封装\\n│   ├── upload          -- 文件上传\\n│   └── ...\\n├── routes          -- 路由模块\\n├── config          -- 配置文件\\n├── test            -- 测试模块\\n├── log             -- 日志文件\\n├── Dockerfile\\n└── main.go\\n```\\n\\n## 项目介绍\\n\\n前台：\\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\\n- 实现点赞，统计用户等功能 (Redis)\\n- 评论 + 回复评论功能\\n- 留言采用弹幕墙，效果炫酷\\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\\n\\n后台：\\n- 鉴权使用 JWT\\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\\n- 支持动态权限修改，前端菜单由后端生成\\n- 文章编辑使用 Markdown 编辑器\\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\\n\\n其他：\\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\\n- 代码整洁层次清晰，利于开发者学习\\n\\n## 技术介绍\\n\\n> 这里只写一些主流的通用技术，详细第三方库: 前端参考 `package.json` 文件，后端参考 `go.mod` 文件\\n\\n前台前端：使用 pnpm 包管理工具\\n- 基于 TypeScript\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia\\n- Vue Router \\n- Axios \\n- Naive UI\\n- Vuetify\\n- ...\\n\\n后台前端：使用 pnpm 包管理工具\\n- 基于 JavaSciprt \\n- pnpm: 包管理工具\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia \\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后端技术栈：\\n- 基于 Golang\\n- Docker\\n- Gin\\n- GORM\\n- Viper: 使用 TOML 作为配置文件\\n- Zap\\n- MySQL\\n- Redis\\n- Nginx: 部署静态资源 + 反向代理\\n- ...\\n\\n其他：\\n- 腾讯云人机验证\\n- 七牛云对象存储\\n- ...\\n\\n## 运行环境\\n\\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\\n\\n对象存储：七牛云\\n\\n## 开发环境\\n\\n| 开发工具                      | 说明                    |\\n| ----------------------------- | ----------------------- |\\n| Vscode                        | Golang 后端 +  Vue 前端 |\\n| Navicat                       | MySQL 远程连接工具      |\\n| Another Redis Desktop Manager | Redis 远程连接工具      |\\n| MobaXterm                     | Linux 远程工具          |\\n| Apifox                        | 接口调试 + 文档生成     |\\n\\n\\n| 开发环境 | 版本 |\\n| -------- | ---- |\\n| Golang   | 1.19 |\\n| MySQL    | 8.x  |\\n| Redis    | 7.x  |\\n\\n## 快速开始\\n\\n### 本地运行\\n\\n> 自行安装 Golang、Node、MySQL、Redis 环境\\n\\n拉取项目到本地：\\n\\n```bash\\ngit clone https://github.com/szluyu99/gin-vue-blog.git\\n```\\n\\n后端项目运行：\\n\\n```bash\\n# 1、进入后端项目根目录 \\ncd gin-blog-server\\n\\n# 2、修改项目运行的配置文件，默认加载 config/config.toml \\n\\n# 3、MySQL 导入 ginblog.sql\\n\\n# 4、启动 Redis \\n\\n# 5、运行项目\\ngo mod tidy\\ngo run main.go\\n```\\n\\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\\n\\n```bash\\nnpm install -g pnpm\\n```\\n\\n前台前端：\\n\\n```bash\\n# 1、进入前台前端项目根目录\\ncd gin-blog-front\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n后台前端：\\n\\n```bash\\n# 1、进入后台前端项目根目录\\ncd gin-blog-admin\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n### 项目部署\\n\\n目前暂时不推荐将本博客部署上生产环境，因为还有太多功能未完善。\\n\\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\\n\\n等功能开发的差不多了，再专门针对部署写一篇文章。\\n\\n---\\n\\n这里简单介绍一下，有基础的同学可以自行折腾。\\n\\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\\n\\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\\n\\nDocker 打包成镜像指令：\\n\\n```bash\\ndocker build -t ginblog .\\n```\\n\\n> 以上只是简单说明，等功能大致完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\\n\\n## 项目总结\\n\\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\\n\\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\\n\\n## 后续计划\\n\\n高优先级: \\n\\n- 完善图片上传功能, 目前文件上传还没怎么处理\\n- 后台首页重新设计（目前没放什么内容）\\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\\n- 前台首页搜索文章（ElasticSearch 搜索）\\n- 博客文章导入导出 (.md 文件)\\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\\n- 后端日志切割\\n- 后台修改背景图片，博客配置等\\n- 相册\\n\\n后续有空安排上：\\n- 适配移动端\\n- 黑夜模式\\n- 前台收缩侧边信息功能\\n- 说说\\n- 音乐播放器\\n- 鼠标左击特效\\n- 看板娘\\n- 文章目录锚点跟随\\n- 第三方登录\\n- 评论时支持选择表情，参考 Valine\\n- 若干细节需要完善...\\n\",\"img\":\"https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg\",\"category_name\":\"后端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"original_url\":\"\",\"is_top\":0,\"status\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (10, '2022-12-04 00:46:13.963', '2022-12-04 00:46:13.963', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"## 博客介绍\\n\\n<p  align=center>\\n<a  href=\\\"http://www.hahacode.cn\\\">\\n<img  src=\\\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\\\"  width=\\\"200\\\"  hight=\\\"200\\\"  alt=\\\"阵、雨的个人博客\\\"  style=\\\"border-radius: 50%\\\">\\n</a>\\n</p>\\n\\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\\n\\n您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。\\n\\n\\nGithub地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\nGitee地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\n## 在线地址\\n\\n博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)\\n\\n博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)\\n\\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\\n\\n在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)\\n\\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\\n\\n## 目录结构\\n\\n代码仓库目录：\\n\\n```bash\\ngin-vue-blog\\n├── gin-blog-admin      -- 博客后台前端\\n├── gin-blog-front      -- 博客前台前端\\n├── gin-blog-server     -- 博客后端\\n```\\n\\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\\n\\n> 项目运行参考：[快速开始](#快速开始)\\n\\n后端目录：\\n\\n```bash\\ngin-blog-server\\n├── api             -- API\\n│   ├── front       -- 前台接口\\n│   └── v1          -- 后台接口\\n├── dao             -- 数据库操作模块\\n├── service         -- 服务模块\\n├── model           -- 数据模型\\n│   ├── req             -- 请求 VO 模型\\n│   ├── resp            -- 响应 VO 模型\\n│   ├── dto             -- 内部传输 DTO 模型\\n│   └── ...             -- 数据库模型对象 PO 模型\\n├── routes          -- 路由模块\\n│   └── middleware      -- 路由中间件\\n├── utils           -- 工具模块\\n│   ├── r               -- 响应封装\\n│   ├── upload          -- 文件上传\\n│   └── ...\\n├── routes          -- 路由模块\\n├── config          -- 配置文件\\n├── test            -- 测试模块\\n├── log             -- 日志文件\\n├── Dockerfile\\n└── main.go\\n```\\n\\n## 项目介绍\\n\\n前台：\\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\\n- 实现点赞，统计用户等功能 (Redis)\\n- 评论 + 回复评论功能\\n- 留言采用弹幕墙，效果炫酷\\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\\n\\n后台：\\n- 鉴权使用 JWT\\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\\n- 支持动态权限修改，前端菜单由后端生成\\n- 文章编辑使用 Markdown 编辑器\\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\\n\\n其他：\\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\\n- 代码整洁层次清晰，利于开发者学习\\n\\n## 技术介绍\\n\\n> 这里只写一些主流的通用技术，详细第三方库: 前端参考 `package.json` 文件，后端参考 `go.mod` 文件\\n\\n前台前端：使用 pnpm 包管理工具\\n- 基于 TypeScript\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia\\n- Vue Router \\n- Axios \\n- Naive UI\\n- Vuetify\\n- ...\\n\\n后台前端：使用 pnpm 包管理工具\\n- 基于 JavaSciprt \\n- pnpm: 包管理工具\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia \\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后端技术栈：\\n- 基于 Golang\\n- Docker\\n- Gin\\n- GORM\\n- Viper: 使用 TOML 作为配置文件\\n- Zap\\n- MySQL\\n- Redis\\n- Nginx: 部署静态资源 + 反向代理\\n- ...\\n\\n其他：\\n- 腾讯云人机验证\\n- 七牛云对象存储\\n- ...\\n\\n## 运行环境\\n\\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\\n\\n对象存储：七牛云\\n\\n## 开发环境\\n\\n| 开发工具                      | 说明                    |\\n| ----------------------------- | ----------------------- |\\n| Vscode                        | Golang 后端 +  Vue 前端 |\\n| Navicat                       | MySQL 远程连接工具      |\\n| Another Redis Desktop Manager | Redis 远程连接工具      |\\n| MobaXterm                     | Linux 远程工具          |\\n| Apifox                        | 接口调试 + 文档生成     |\\n\\n\\n| 开发环境 | 版本 |\\n| -------- | ---- |\\n| Golang   | 1.19 |\\n| MySQL    | 8.x  |\\n| Redis    | 7.x  |\\n\\n## 快速开始\\n\\n### 本地运行\\n\\n> 自行安装 Golang、Node、MySQL、Redis 环境\\n\\n拉取项目到本地：\\n\\n```bash\\ngit clone https://github.com/szluyu99/gin-vue-blog.git\\n```\\n\\n后端项目运行：\\n\\n```bash\\n# 1、进入后端项目根目录 \\ncd gin-blog-server\\n\\n# 2、修改项目运行的配置文件，默认加载 config/config.toml \\n\\n# 3、MySQL 导入 ginblog.sql\\n\\n# 4、启动 Redis \\n\\n# 5、运行项目\\ngo mod tidy\\ngo run main.go\\n```\\n\\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\\n\\n```bash\\nnpm install -g pnpm\\n```\\n\\n前台前端：\\n\\n```bash\\n# 1、进入前台前端项目根目录\\ncd gin-blog-front\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n后台前端：\\n\\n```bash\\n# 1、进入后台前端项目根目录\\ncd gin-blog-admin\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n### 项目部署\\n\\n目前暂时不推荐将本博客部署上生产环境，因为还有太多功能未完善。\\n\\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\\n\\n等功能开发的差不多了，再专门针对部署写一篇文章。\\n\\n---\\n\\n这里简单介绍一下，有基础的同学可以自行折腾。\\n\\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\\n\\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\\n\\nDocker 打包成镜像指令：\\n\\n```bash\\ndocker build -t ginblog .\\n```\\n\\n> 以上只是简单说明，等功能大致完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\\n\\n## 项目总结\\n\\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\\n\\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\\n\\n## 后续计划\\n\\n高优先级: \\n\\n- 完善图片上传功能, 目前文件上传还没怎么处理\\n- 后台首页重新设计（目前没放什么内容）\\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\\n- 前台首页搜索文章（ElasticSearch 搜索）\\n- 博客文章导入导出 (.md 文件)\\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\\n- 后端日志切割\\n- 后台修改背景图片，博客配置等\\n- 相册\\n\\n后续有空安排上：\\n- 适配移动端\\n- 黑夜模式\\n- 前台收缩侧边信息功能\\n- 说说\\n- 音乐播放器\\n- 鼠标左击特效\\n- 看板娘\\n- 文章目录锚点跟随\\n- 第三方登录\\n- 评论时支持选择表情，参考 Valine\\n- 若干细节需要完善...\\n\",\"img\":\"https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg\",\"category_name\":\"后端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"original_url\":\"\",\"is_top\":0,\"status\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58096', '未知');
INSERT INTO `operation_log` VALUES (11, '2022-12-07 20:47:08.356', '2022-12-07 20:47:08.356', '菜单', '新增或修改', 'gin-blog/api/v1.(*Menu).SaveOrUpdate-fm', '/api/menu', '新增或修改菜单', '{\"order_num\":7,\"is_hidden\":0,\"is_catelogue\":false,\"component\":\"/profile\",\"parent_id\":0,\"name\":\"个人中心\",\"icon\":\"mdi:account\",\"path\":\"/profile\",\"keep_alive\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (12, '2022-12-07 20:47:21.159', '2022-12-07 20:47:21.159', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":1,\"name\":\"管理员\",\"label\":\"admin\",\"created_at\":\"2022-10-20T21:24:28+08:00\",\"is_disable\":0,\"resource_ids\":[3,43,44,45,46,47,48,6,59,60,61,7,38,39,40,41,42,8,65,66,67,68,9,62,63,64,10,23,34,35,36,37,79,80,81,85,49,51,52,53,54,50,55,56,57,58,69,70,71,72,91,92,93,74,78,82,84,86,98,95,11],\"menu_ids\":[2,3,4,5,6,7,8,9,16,17,23,24,25,26,27,28,29,30,31,32,33,36,37,38,34,10,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (13, '2022-12-07 20:47:23.403', '2022-12-07 20:47:23.403', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":2,\"name\":\"用户\",\"label\":\"user\",\"created_at\":\"2022-10-20T21:25:07+08:00\",\"is_disable\":0,\"resource_ids\":[43,44,59,38,42,65,68,62,35,36,51,54,55,58,70,78,82,79,80,86,92,95,41],\"menu_ids\":[1,2,6,8,9,10,3,25,26,16,17,23,24,4,27,28,29,5,30,32,31,33,34,36,37,38,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (14, '2022-12-07 20:47:25.953', '2022-12-07 20:47:25.953', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":3,\"name\":\"测试\",\"label\":\"test\",\"created_at\":\"2022-10-20T21:25:09+08:00\",\"is_disable\":0,\"resource_ids\":[43,44,59,38,41,42,65,68,62,35,36,79,80,51,54,55,58,70,78,82,92,95,86],\"menu_ids\":[1,2,3,4,33,6,34,8,9,10,25,26,16,17,23,24,5,29,30,32,31,27,28,36,37,38,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (15, '2022-12-07 20:48:05.951', '2022-12-07 20:48:05.951', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"parent_id\":74,\"name\":\"修改当前用户密码\",\"url\":\"/user/current/password\",\"request_method\":\"PUT\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (16, '2022-12-07 20:48:35.522', '2022-12-07 20:48:35.522', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"parent_id\":74,\"name\":\"修改当前用户信息\",\"url\":\"/user/current\",\"request_method\":\"PUT\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (17, '2022-12-07 20:48:56.232', '2022-12-07 20:48:56.232', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":1,\"name\":\"管理员\",\"label\":\"admin\",\"created_at\":\"2022-10-20T21:24:28+08:00\",\"is_disable\":0,\"resource_ids\":[3,43,44,45,46,47,48,6,59,60,61,7,38,39,40,41,42,8,65,66,67,68,9,62,63,64,10,23,34,35,36,37,79,80,81,85,49,51,52,53,54,50,55,56,57,58,69,70,71,72,91,92,93,78,82,84,86,98,95,11,99,100,74],\"menu_ids\":[2,3,4,5,6,7,8,9,16,17,23,24,25,26,27,28,29,30,31,32,33,36,37,38,34,10,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (18, '2022-12-07 20:49:09.822', '2022-12-07 20:49:09.822', '用户', '修改', 'gin-blog/api/v1.(*User).UpdateCurrent-fm', '/api/user/current', '修改用户', '{\"avatar\":\"https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png\",\"nickname\":\"管理员\",\"intro\":\"我是管理员用户！\",\"website\":\"https://www.baidu.com\"}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (19, '2022-12-07 20:49:20.794', '2022-12-07 20:49:20.794', '用户', '修改', 'gin-blog/api/v1.(*User).UpdateCurrentPassword-fm', '/api/user/current/password', '修改用户', '{\"old_password\":\"1234567\",\"new_password\":\"1234567\",\"confirm_password\":\"1234567\"}', 'PUT', '{\"code\":1010,\"message\":\"旧密码不正确\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (20, '2022-12-07 20:53:16.642', '2022-12-07 20:53:16.642', '文章', '修改', 'gin-blog/api/v1.(*Article).UpdateTop-fm', '/api/article/top', '修改文章', '{\"id\":1,\"created_at\":\"2022-12-03T22:07:01.638+08:00\",\"updated_at\":\"2022-12-04T00:46:13.94+08:00\",\"category_id\":1,\"category\":{\"id\":1,\"created_at\":\"2022-12-03T22:01:29.106+08:00\",\"updated_at\":\"2022-12-03T22:01:29.106+08:00\",\"name\":\"后端\",\"Articles\":null},\"tags\":[{\"id\":1,\"created_at\":\"2022-12-03T22:01:51.624+08:00\",\"updated_at\":\"2022-12-03T22:01:51.624+08:00\",\"articles\":null,\"name\":\"Golang\"},{\"id\":2,\"created_at\":\"2022-12-03T22:01:56.984+08:00\",\"updated_at\":\"2022-12-03T22:01:56.984+08:00\",\"articles\":null,\"name\":\"Vue\"}],\"user_id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"## 博客介绍\\n\\n<p  align=center>\\n<a  href=\\\"http://www.hahacode.cn\\\">\\n<img  src=\\\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\\\"  width=\\\"200\\\"  hight=\\\"200\\\"  alt=\\\"阵、雨的个人博客\\\"  style=\\\"border-radius: 50%\\\">\\n</a>\\n</p>\\n\\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\\n\\n您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。\\n\\n\\nGithub地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\nGitee地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\n## 在线地址\\n\\n博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)\\n\\n博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)\\n\\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\\n\\n在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)\\n\\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\\n\\n## 目录结构\\n\\n代码仓库目录：\\n\\n```bash\\ngin-vue-blog\\n├── gin-blog-admin      -- 博客后台前端\\n├── gin-blog-front      -- 博客前台前端\\n├── gin-blog-server     -- 博客后端\\n```\\n\\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\\n\\n> 项目运行参考：[快速开始](#快速开始)\\n\\n后端目录：\\n\\n```bash\\ngin-blog-server\\n├── api             -- API\\n│   ├── front       -- 前台接口\\n│   └── v1          -- 后台接口\\n├── dao             -- 数据库操作模块\\n├── service         -- 服务模块\\n├── model           -- 数据模型\\n│   ├── req             -- 请求 VO 模型\\n│   ├── resp            -- 响应 VO 模型\\n│   ├── dto             -- 内部传输 DTO 模型\\n│   └── ...             -- 数据库模型对象 PO 模型\\n├── routes          -- 路由模块\\n│   └── middleware      -- 路由中间件\\n├── utils           -- 工具模块\\n│   ├── r               -- 响应封装\\n│   ├── upload          -- 文件上传\\n│   └── ...\\n├── routes          -- 路由模块\\n├── config          -- 配置文件\\n├── test            -- 测试模块\\n├── log             -- 日志文件\\n├── Dockerfile\\n└── main.go\\n```\\n\\n## 项目介绍\\n\\n前台：\\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\\n- 实现点赞，统计用户等功能 (Redis)\\n- 评论 + 回复评论功能\\n- 留言采用弹幕墙，效果炫酷\\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\\n\\n后台：\\n- 鉴权使用 JWT\\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\\n- 支持动态权限修改，前端菜单由后端生成\\n- 文章编辑使用 Markdown 编辑器\\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\\n\\n其他：\\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\\n- 代码整洁层次清晰，利于开发者学习\\n\\n## 技术介绍\\n\\n> 这里只写一些主流的通用技术，详细第三方库: 前端参考 `package.json` 文件，后端参考 `go.mod` 文件\\n\\n前台前端：使用 pnpm 包管理工具\\n- 基于 TypeScript\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia\\n- Vue Router \\n- Axios \\n- Naive UI\\n- Vuetify\\n- ...\\n\\n后台前端：使用 pnpm 包管理工具\\n- 基于 JavaSciprt \\n- pnpm: 包管理工具\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia \\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后端技术栈：\\n- 基于 Golang\\n- Docker\\n- Gin\\n- GORM\\n- Viper: 使用 TOML 作为配置文件\\n- Zap\\n- MySQL\\n- Redis\\n- Nginx: 部署静态资源 + 反向代理\\n- ...\\n\\n其他：\\n- 腾讯云人机验证\\n- 七牛云对象存储\\n- ...\\n\\n## 运行环境\\n\\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\\n\\n对象存储：七牛云\\n\\n## 开发环境\\n\\n| 开发工具                      | 说明                    |\\n| ----------------------------- | ----------------------- |\\n| Vscode                        | Golang 后端 +  Vue 前端 |\\n| Navicat                       | MySQL 远程连接工具      |\\n| Another Redis Desktop Manager | Redis 远程连接工具      |\\n| MobaXterm                     | Linux 远程工具          |\\n| Apifox                        | 接口调试 + 文档生成     |\\n\\n\\n| 开发环境 | 版本 |\\n| -------- | ---- |\\n| Golang   | 1.19 |\\n| MySQL    | 8.x  |\\n| Redis    | 7.x  |\\n\\n## 快速开始\\n\\n### 本地运行\\n\\n> 自行安装 Golang、Node、MySQL、Redis 环境\\n\\n拉取项目到本地：\\n\\n```bash\\ngit clone https://github.com/szluyu99/gin-vue-blog.git\\n```\\n\\n后端项目运行：\\n\\n```bash\\n# 1、进入后端项目根目录 \\ncd gin-blog-server\\n\\n# 2、修改项目运行的配置文件，默认加载 config/config.toml \\n\\n# 3、MySQL 导入 ginblog.sql\\n\\n# 4、启动 Redis \\n\\n# 5、运行项目\\ngo mod tidy\\ngo run main.go\\n```\\n\\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\\n\\n```bash\\nnpm install -g pnpm\\n```\\n\\n前台前端：\\n\\n```bash\\n# 1、进入前台前端项目根目录\\ncd gin-blog-front\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n后台前端：\\n\\n```bash\\n# 1、进入后台前端项目根目录\\ncd gin-blog-admin\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n### 项目部署\\n\\n目前暂时不推荐将本博客部署上生产环境，因为还有太多功能未完善。\\n\\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\\n\\n等功能开发的差不多了，再专门针对部署写一篇文章。\\n\\n---\\n\\n这里简单介绍一下，有基础的同学可以自行折腾。\\n\\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\\n\\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\\n\\nDocker 打包成镜像指令：\\n\\n```bash\\ndocker build -t ginblog .\\n```\\n\\n> 以上只是简单说明，等功能大致完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\\n\\n## 项目总结\\n\\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\\n\\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\\n\\n## 后续计划\\n\\n高优先级: \\n\\n- 完善图片上传功能, 目前文件上传还没怎么处理\\n- 后台首页重新设计（目前没放什么内容）\\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\\n- 前台首页搜索文章（ElasticSearch 搜索）\\n- 博客文章导入导出 (.md 文件)\\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\\n- 后端日志切割\\n- 后台修改背景图片，博客配置等\\n- 相册\\n\\n后续有空安排上：\\n- 适配移动端\\n- 黑夜模式\\n- 前台收缩侧边信息功能\\n- 说说\\n- 音乐播放器\\n- 鼠标左击特效\\n- 看板娘\\n- 文章目录锚点跟随\\n- 第三方登录\\n- 评论时支持选择表情，参考 Valine\\n- 若干细节需要完善...\\n\",\"img\":\"https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg\",\"type\":1,\"status\":1,\"is_top\":1,\"is_delete\":0,\"original_url\":\"\",\"comment_count\":0,\"read_count\":0,\"publishing\":true}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (21, '2022-12-07 20:53:17.524', '2022-12-07 20:53:17.524', '文章', '修改', 'gin-blog/api/v1.(*Article).UpdateTop-fm', '/api/article/top', '修改文章', '{\"id\":1,\"created_at\":\"2022-12-03T22:07:01.638+08:00\",\"updated_at\":\"2022-12-07T20:53:16.635+08:00\",\"category_id\":1,\"category\":{\"id\":1,\"created_at\":\"2022-12-03T22:01:29.106+08:00\",\"updated_at\":\"2022-12-03T22:01:29.106+08:00\",\"name\":\"后端\",\"Articles\":null},\"tags\":[{\"id\":1,\"created_at\":\"2022-12-03T22:01:51.624+08:00\",\"updated_at\":\"2022-12-03T22:01:51.624+08:00\",\"articles\":null,\"name\":\"Golang\"},{\"id\":2,\"created_at\":\"2022-12-03T22:01:56.984+08:00\",\"updated_at\":\"2022-12-03T22:01:56.984+08:00\",\"articles\":null,\"name\":\"Vue\"}],\"user_id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"## 博客介绍\\n\\n<p  align=center>\\n<a  href=\\\"http://www.hahacode.cn\\\">\\n<img  src=\\\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\\\"  width=\\\"200\\\"  hight=\\\"200\\\"  alt=\\\"阵、雨的个人博客\\\"  style=\\\"border-radius: 50%\\\">\\n</a>\\n</p>\\n\\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\\n\\n您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。\\n\\n\\nGithub地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\nGitee地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\n## 在线地址\\n\\n博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)\\n\\n博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)\\n\\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\\n\\n在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)\\n\\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\\n\\n## 目录结构\\n\\n代码仓库目录：\\n\\n```bash\\ngin-vue-blog\\n├── gin-blog-admin      -- 博客后台前端\\n├── gin-blog-front      -- 博客前台前端\\n├── gin-blog-server     -- 博客后端\\n```\\n\\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\\n\\n> 项目运行参考：[快速开始](#快速开始)\\n\\n后端目录：\\n\\n```bash\\ngin-blog-server\\n├── api             -- API\\n│   ├── front       -- 前台接口\\n│   └── v1          -- 后台接口\\n├── dao             -- 数据库操作模块\\n├── service         -- 服务模块\\n├── model           -- 数据模型\\n│   ├── req             -- 请求 VO 模型\\n│   ├── resp            -- 响应 VO 模型\\n│   ├── dto             -- 内部传输 DTO 模型\\n│   └── ...             -- 数据库模型对象 PO 模型\\n├── routes          -- 路由模块\\n│   └── middleware      -- 路由中间件\\n├── utils           -- 工具模块\\n│   ├── r               -- 响应封装\\n│   ├── upload          -- 文件上传\\n│   └── ...\\n├── routes          -- 路由模块\\n├── config          -- 配置文件\\n├── test            -- 测试模块\\n├── log             -- 日志文件\\n├── Dockerfile\\n└── main.go\\n```\\n\\n## 项目介绍\\n\\n前台：\\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\\n- 实现点赞，统计用户等功能 (Redis)\\n- 评论 + 回复评论功能\\n- 留言采用弹幕墙，效果炫酷\\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\\n\\n后台：\\n- 鉴权使用 JWT\\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\\n- 支持动态权限修改，前端菜单由后端生成\\n- 文章编辑使用 Markdown 编辑器\\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\\n\\n其他：\\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\\n- 代码整洁层次清晰，利于开发者学习\\n\\n## 技术介绍\\n\\n> 这里只写一些主流的通用技术，详细第三方库: 前端参考 `package.json` 文件，后端参考 `go.mod` 文件\\n\\n前台前端：使用 pnpm 包管理工具\\n- 基于 TypeScript\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia\\n- Vue Router \\n- Axios \\n- Naive UI\\n- Vuetify\\n- ...\\n\\n后台前端：使用 pnpm 包管理工具\\n- 基于 JavaSciprt \\n- pnpm: 包管理工具\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia \\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后端技术栈：\\n- 基于 Golang\\n- Docker\\n- Gin\\n- GORM\\n- Viper: 使用 TOML 作为配置文件\\n- Zap\\n- MySQL\\n- Redis\\n- Nginx: 部署静态资源 + 反向代理\\n- ...\\n\\n其他：\\n- 腾讯云人机验证\\n- 七牛云对象存储\\n- ...\\n\\n## 运行环境\\n\\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\\n\\n对象存储：七牛云\\n\\n## 开发环境\\n\\n| 开发工具                      | 说明                    |\\n| ----------------------------- | ----------------------- |\\n| Vscode                        | Golang 后端 +  Vue 前端 |\\n| Navicat                       | MySQL 远程连接工具      |\\n| Another Redis Desktop Manager | Redis 远程连接工具      |\\n| MobaXterm                     | Linux 远程工具          |\\n| Apifox                        | 接口调试 + 文档生成     |\\n\\n\\n| 开发环境 | 版本 |\\n| -------- | ---- |\\n| Golang   | 1.19 |\\n| MySQL    | 8.x  |\\n| Redis    | 7.x  |\\n\\n## 快速开始\\n\\n### 本地运行\\n\\n> 自行安装 Golang、Node、MySQL、Redis 环境\\n\\n拉取项目到本地：\\n\\n```bash\\ngit clone https://github.com/szluyu99/gin-vue-blog.git\\n```\\n\\n后端项目运行：\\n\\n```bash\\n# 1、进入后端项目根目录 \\ncd gin-blog-server\\n\\n# 2、修改项目运行的配置文件，默认加载 config/config.toml \\n\\n# 3、MySQL 导入 ginblog.sql\\n\\n# 4、启动 Redis \\n\\n# 5、运行项目\\ngo mod tidy\\ngo run main.go\\n```\\n\\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\\n\\n```bash\\nnpm install -g pnpm\\n```\\n\\n前台前端：\\n\\n```bash\\n# 1、进入前台前端项目根目录\\ncd gin-blog-front\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n后台前端：\\n\\n```bash\\n# 1、进入后台前端项目根目录\\ncd gin-blog-admin\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n### 项目部署\\n\\n目前暂时不推荐将本博客部署上生产环境，因为还有太多功能未完善。\\n\\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\\n\\n等功能开发的差不多了，再专门针对部署写一篇文章。\\n\\n---\\n\\n这里简单介绍一下，有基础的同学可以自行折腾。\\n\\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\\n\\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\\n\\nDocker 打包成镜像指令：\\n\\n```bash\\ndocker build -t ginblog .\\n```\\n\\n> 以上只是简单说明，等功能大致完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\\n\\n## 项目总结\\n\\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\\n\\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\\n\\n## 后续计划\\n\\n高优先级: \\n\\n- 完善图片上传功能, 目前文件上传还没怎么处理\\n- 后台首页重新设计（目前没放什么内容）\\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\\n- 前台首页搜索文章（ElasticSearch 搜索）\\n- 博客文章导入导出 (.md 文件)\\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\\n- 后端日志切割\\n- 后台修改背景图片，博客配置等\\n- 相册\\n\\n后续有空安排上：\\n- 适配移动端\\n- 黑夜模式\\n- 前台收缩侧边信息功能\\n- 说说\\n- 音乐播放器\\n- 鼠标左击特效\\n- 看板娘\\n- 文章目录锚点跟随\\n- 第三方登录\\n- 评论时支持选择表情，参考 Valine\\n- 若干细节需要完善...\\n\",\"img\":\"https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg\",\"type\":1,\"status\":1,\"is_top\":0,\"is_delete\":0,\"original_url\":\"\",\"comment_count\":0,\"read_count\":0,\"publishing\":true}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (22, '2022-12-07 20:53:26.811', '2022-12-07 20:53:26.811', '菜单', '新增或修改', 'gin-blog/api/v1.(*Menu).SaveOrUpdate-fm', '/api/menu', '新增或修改菜单', '{\"id\":33,\"name\":\"首页\",\"path\":\"/home\",\"component\":\"/home\",\"icon\":\"ic:sharp-home\",\"created_at\":\"2022-11-01T01:43:10.142+08:00\",\"order_num\":0,\"children\":[],\"parent_id\":0,\"is_hidden\":1,\"keep_alive\":1,\"redirect\":\"\",\"publishing\":true}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (23, '2022-12-07 20:53:27.481', '2022-12-07 20:53:27.481', '菜单', '新增或修改', 'gin-blog/api/v1.(*Menu).SaveOrUpdate-fm', '/api/menu', '新增或修改菜单', '{\"id\":33,\"name\":\"首页\",\"path\":\"/home\",\"component\":\"/home\",\"icon\":\"ic:sharp-home\",\"created_at\":\"2022-11-01T01:43:10.142+08:00\",\"order_num\":0,\"children\":[],\"parent_id\":0,\"is_hidden\":0,\"keep_alive\":1,\"redirect\":\"\",\"publishing\":true}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (24, '2022-12-07 20:53:31.343', '2022-12-07 20:53:31.343', '菜单', '新增或修改', 'gin-blog/api/v1.(*Menu).SaveOrUpdate-fm', '/api/menu', '新增或修改菜单', '{\"id\":39,\"name\":\"个人中心\",\"path\":\"/profile\",\"component\":\"/profile\",\"icon\":\"mdi:account\",\"created_at\":\"2022-12-07T20:47:08.349+08:00\",\"order_num\":7,\"children\":[],\"parent_id\":0,\"is_hidden\":1,\"keep_alive\":1,\"redirect\":\"\",\"publishing\":true}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (25, '2022-12-07 20:53:33.858', '2022-12-07 20:53:33.858', '菜单', '新增或修改', 'gin-blog/api/v1.(*Menu).SaveOrUpdate-fm', '/api/menu', '新增或修改菜单', '{\"id\":39,\"name\":\"个人中心\",\"path\":\"/profile\",\"component\":\"/profile\",\"icon\":\"mdi:account\",\"created_at\":\"2022-12-07T20:47:08.349+08:00\",\"order_num\":7,\"children\":[],\"parent_id\":0,\"is_hidden\":0,\"keep_alive\":1,\"redirect\":\"\",\"publishing\":true}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (26, '2022-12-07 20:55:08.282', '2022-12-07 20:55:08.282', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"parent_id\":74,\"name\":\"修改用户禁用\",\"url\":\"/user/disable\",\"request_method\":\"PUT\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (27, '2022-12-07 20:55:21.851', '2022-12-07 20:55:21.851', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":1,\"name\":\"管理员\",\"label\":\"admin\",\"created_at\":\"2022-10-20T21:24:28+08:00\",\"is_disable\":0,\"resource_ids\":[3,43,44,45,46,47,48,6,59,60,61,7,38,39,40,41,42,8,65,66,67,68,9,62,63,64,10,23,34,35,36,37,79,80,81,85,49,51,52,53,54,50,55,56,57,58,69,70,71,72,91,92,93,78,82,84,86,98,95,11,99,100,101,74],\"menu_ids\":[2,3,4,5,6,7,8,9,16,17,23,24,25,26,27,28,29,30,31,32,33,36,37,38,34,10,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (28, '2022-12-07 20:55:28.941', '2022-12-07 20:55:28.941', '用户', '修改', 'gin-blog/api/v1.(*User).UpdateDisable-fm', '/api/user/disable', '修改用户', '{\"id\":2,\"user_info_id\":2,\"avatar\":\"https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png\",\"nickname\":\"普通用户\",\"roles\":[{\"id\":2,\"created_at\":\"2022-10-20T21:25:07+08:00\",\"updated_at\":\"2022-12-07T20:47:23.344+08:00\",\"name\":\"用户\",\"label\":\"user\",\"is_disable\":0}],\"login_type\":1,\"ip_address\":\"172.17.0.1:40280\",\"ip_source\":\"未知\",\"created_at\":\"2022-10-19T22:31:26.805+08:00\",\"last_login_time\":\"2022-12-03T12:20:36.539+08:00\",\"is_disable\":1,\"publishing\":true}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (29, '2022-12-07 20:55:30.353', '2022-12-07 20:55:30.353', '用户', '修改', 'gin-blog/api/v1.(*User).UpdateDisable-fm', '/api/user/disable', '修改用户', '{\"id\":2,\"user_info_id\":2,\"avatar\":\"https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png\",\"nickname\":\"普通用户\",\"roles\":[{\"id\":2,\"created_at\":\"2022-10-20T21:25:07+08:00\",\"updated_at\":\"2022-12-07T20:47:23.344+08:00\",\"name\":\"用户\",\"label\":\"user\",\"is_disable\":0}],\"login_type\":1,\"ip_address\":\"172.17.0.1:40280\",\"ip_source\":\"未知\",\"created_at\":\"2022-10-19T22:31:26.805+08:00\",\"last_login_time\":\"2022-12-03T12:20:36.539+08:00\",\"is_disable\":0,\"publishing\":true}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:58630', '未知');
INSERT INTO `operation_log` VALUES (30, '2022-12-08 15:43:15.429', '2022-12-08 15:43:15.429', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"name\":\"页面模块\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:38804', '未知');
INSERT INTO `operation_log` VALUES (31, '2022-12-08 15:43:26.018', '2022-12-08 15:43:26.018', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"parent_id\":102,\"name\":\"页面列表\",\"url\":\"/page/list\",\"request_method\":\"GET\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:38804', '未知');
INSERT INTO `operation_log` VALUES (32, '2022-12-08 15:43:38.579', '2022-12-08 15:43:38.579', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"parent_id\":102,\"name\":\"新增/编辑页面\",\"url\":\"/page\",\"request_method\":\"POST\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:38804', '未知');
INSERT INTO `operation_log` VALUES (33, '2022-12-08 15:43:50.888', '2022-12-08 15:43:50.888', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"parent_id\":102,\"name\":\"删除页面\",\"url\":\"/page\",\"request_method\":\"DELETE\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:38804', '未知');
INSERT INTO `operation_log` VALUES (34, '2022-12-08 15:44:04.384', '2022-12-08 15:44:04.384', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":1,\"name\":\"管理员\",\"label\":\"admin\",\"created_at\":\"2022-10-20T21:24:28+08:00\",\"is_disable\":0,\"resource_ids\":[3,43,44,45,46,47,48,6,59,60,61,7,38,39,40,41,42,8,65,66,67,68,9,62,63,64,10,23,34,35,36,37,79,80,81,85,49,51,52,53,54,50,55,56,57,58,69,70,71,72,91,92,93,78,82,84,86,98,95,11,99,100,101,74,102,103,104,105],\"menu_ids\":[2,3,4,5,6,7,8,9,16,17,23,24,25,26,27,28,29,30,31,32,33,36,37,38,34,10,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:38804', '未知');
INSERT INTO `operation_log` VALUES (35, '2022-12-08 15:44:09.129', '2022-12-08 15:44:09.129', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":2,\"name\":\"用户\",\"label\":\"user\",\"created_at\":\"2022-10-20T21:25:07+08:00\",\"is_disable\":0,\"resource_ids\":[43,44,59,38,42,65,68,62,35,36,51,54,55,58,70,78,82,79,80,86,92,95,41,103],\"menu_ids\":[1,2,6,8,9,10,3,25,26,16,17,23,24,4,27,28,29,5,30,32,31,33,34,36,37,38,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:38804', '未知');
INSERT INTO `operation_log` VALUES (36, '2022-12-08 15:44:13.962', '2022-12-08 15:44:13.962', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":3,\"name\":\"测试\",\"label\":\"test\",\"created_at\":\"2022-10-20T21:25:09+08:00\",\"is_disable\":0,\"resource_ids\":[43,44,59,38,41,42,65,68,62,35,36,79,80,51,54,55,58,70,78,82,92,95,86,103],\"menu_ids\":[1,2,3,4,33,6,34,8,9,10,25,26,16,17,23,24,5,29,30,32,31,27,28,36,37,38,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:38804', '未知');
INSERT INTO `operation_log` VALUES (37, '2022-12-16 11:51:53.942', '2022-12-16 11:51:53.942', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"## 博客介绍\\n\\n<p  align=center>\\n<a  href=\\\"http://www.hahacode.cn\\\">\\n<img  src=\\\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\\\"  width=\\\"200\\\"  hight=\\\"200\\\"  alt=\\\"阵、雨的个人博客\\\"  style=\\\"border-radius: 50%\\\">\\n</a>\\n</p>\\n\\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\\n\\n您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。\\n\\n\\nGithub地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\nGitee地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\n## 在线地址\\n\\n博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)\\n\\n博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)\\n\\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\\n\\n在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)\\n\\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\\n\\n## 目录结构\\n\\n代码仓库目录：\\n\\n```bash\\ngin-vue-blog\\n├── gin-blog-admin      -- 博客后台前端\\n├── gin-blog-front      -- 博客前台前端\\n├── gin-blog-server     -- 博客后端\\n```\\n\\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\\n\\n> 项目运行参考：[快速开始](#快速开始)\\n\\n后端目录：\\n\\n```bash\\ngin-blog-server\\n├── api             -- API\\n│   ├── front       -- 前台接口\\n│   └── v1          -- 后台接口\\n├── dao             -- 数据库操作模块\\n├── service         -- 服务模块\\n├── model           -- 数据模型\\n│   ├── req             -- 请求 VO 模型\\n│   ├── resp            -- 响应 VO 模型\\n│   ├── dto             -- 内部传输 DTO 模型\\n│   └── ...             -- 数据库模型对象 PO 模型\\n├── routes          -- 路由模块\\n│   └── middleware      -- 路由中间件\\n├── utils           -- 工具模块\\n│   ├── r               -- 响应封装\\n│   ├── upload          -- 文件上传\\n│   └── ...\\n├── routes          -- 路由模块\\n├── config          -- 配置文件\\n├── test            -- 测试模块\\n├── log             -- 日志文件\\n├── Dockerfile\\n└── main.go\\n```\\n\\n## 项目介绍\\n\\n前台：\\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\\n- 实现点赞，统计用户等功能 (Redis)\\n- 评论 + 回复评论功能\\n- 留言采用弹幕墙，效果炫酷\\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\\n\\n后台：\\n- 鉴权使用 JWT\\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\\n- 支持动态权限修改，前端菜单由后端生成\\n- 文章编辑使用 Markdown 编辑器\\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\\n\\n其他：\\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\\n- 代码整洁层次清晰，利于开发者学习\\n\\n## 技术介绍\\n\\n> 这里只写一些主流的通用技术，详细第三方库: 前端参考 `package.json` 文件，后端参考 `go.mod` 文件\\n\\n前台前端：使用 pnpm 包管理工具\\n- 基于 TypeScript\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia\\n- Vue Router \\n- Axios \\n- Naive UI\\n- Vuetify\\n- ...\\n\\n后台前端：使用 pnpm 包管理工具\\n- 基于 JavaSciprt \\n- pnpm: 包管理工具\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia \\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后端技术栈：\\n- 基于 Golang\\n- Docker\\n- Gin\\n- GORM\\n- Viper: 使用 TOML 作为配置文件\\n- Zap\\n- MySQL\\n- Redis\\n- Nginx: 部署静态资源 + 反向代理\\n- ...\\n\\n其他：\\n- 腾讯云人机验证\\n- 七牛云对象存储\\n- ...\\n\\n## 运行环境\\n\\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\\n\\n对象存储：七牛云\\n\\n## 开发环境\\n\\n| 开发工具                      | 说明                    |\\n| ----------------------------- | ----------------------- |\\n| Vscode                        | Golang 后端 +  Vue 前端 |\\n| Navicat                       | MySQL 远程连接工具      |\\n| Another Redis Desktop Manager | Redis 远程连接工具      |\\n| MobaXterm                     | Linux 远程工具          |\\n| Apifox                        | 接口调试 + 文档生成     |\\n\\n\\n| 开发环境 | 版本 |\\n| -------- | ---- |\\n| Golang   | 1.19 |\\n| MySQL    | 8.x  |\\n| Redis    | 7.x  |\\n\\n## 快速开始\\n\\n### 本地运行\\n\\n> 自行安装 Golang、Node、MySQL、Redis 环境\\n\\n拉取项目到本地：\\n\\n```bash\\ngit clone https://github.com/szluyu99/gin-vue-blog.git\\n```\\n\\n后端项目运行：\\n\\n```bash\\n# 1、进入后端项目根目录 \\ncd gin-blog-server\\n\\n# 2、修改项目运行的配置文件，默认加载 config/config.toml \\n\\n# 3、MySQL 导入 ginblog.sql\\n\\n# 4、启动 Redis \\n\\n# 5、运行项目\\ngo mod tidy\\ngo run main.go\\n```\\n\\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\\n\\n```bash\\nnpm install -g pnpm\\n```\\n\\n前台前端：\\n\\n```bash\\n# 1、进入前台前端项目根目录\\ncd gin-blog-front\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n后台前端：\\n\\n```bash\\n# 1、进入后台前端项目根目录\\ncd gin-blog-admin\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n### 项目部署\\n\\n目前暂时不推荐将本博客部署上生产环境，因为还有太多功能未完善。\\n\\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\\n\\n等功能开发的差不多了，再专门针对部署写一篇文章。\\n\\n---\\n\\n这里简单介绍一下，有基础的同学可以自行折腾。\\n\\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\\n\\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\\n\\nDocker 打包成镜像指令：\\n\\n```bash\\ndocker build -t ginblog .\\n```\\n\\n> 以上只是简单说明，等功能大致完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\\n\\n## 项目总结\\n\\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\\n\\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\\n\\n## 后续计划\\n\\n高优先级: \\n\\n- 完善图片上传功能, 目前文件上传还没怎么处理\\n- 后台首页重新设计（目前没放什么内容）\\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\\n- 前台首页搜索文章（ElasticSearch 搜索）\\n- 博客文章导入导出 (.md 文件)\\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\\n- 后端日志切割\\n- 后台修改背景图片，博客配置等\\n- 相册\\n\\n后续有空安排上：\\n- 适配移动端\\n- 黑夜模式\\n- 前台收缩侧边信息功能\\n- 说说\\n- 音乐播放器\\n- 鼠标左击特效\\n- 看板娘\\n- 文章目录锚点跟随\\n- 第三方登录\\n- 评论时支持选择表情，参考 Valine\\n- 若干细节需要完善...\\n\",\"img\":\"https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg\",\"category_name\":\"后端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"original_url\":\"\",\"is_top\":0,\"status\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:35424', '未知');
INSERT INTO `operation_log` VALUES (38, '2022-12-16 11:52:39.332', '2022-12-16 11:52:39.332', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"## 博客介绍\\n\\n<p  align=center>\\n<a  href=\\\"http://www.hahacode.cn\\\">\\n<img  src=\\\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\\\"  width=\\\"200\\\"  hight=\\\"200\\\"  alt=\\\"阵、雨的个人博客\\\"  style=\\\"border-radius: 50%\\\">\\n</a>\\n</p>\\n\\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\\n\\n您的 star 是我坚持的动力，感谢大家的支持，欢迎提交 pr 共同改进项目。\\n\\n\\nGithub地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\nGitee地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\n## 在线地址\\n\\n博客前台链接：[www.hahacode.cn](http://www.hahacode.cn)\\n\\n博客后台链接：[www.hahacode.cn/blog-admin](http://www.hahacode.cn/blog-admin)\\n\\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\\n\\n在线接口文档地址：[https://is68368smh.apifox.cn/](https://is68368smh.apifox.cn/)\\n\\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\\n\\n## 目录结构\\n\\n代码仓库目录：\\n\\n```bash\\ngin-vue-blog\\n├── gin-blog-admin      -- 博客后台前端\\n├── gin-blog-front      -- 博客前台前端\\n├── gin-blog-server     -- 博客后端\\n```\\n\\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\\n\\n> 项目运行参考：[快速开始](#快速开始)\\n\\n后端目录：\\n\\n```bash\\ngin-blog-server\\n├── api             -- API\\n│   ├── front       -- 前台接口\\n│   └── v1          -- 后台接口\\n├── dao             -- 数据库操作模块\\n├── service         -- 服务模块\\n├── model           -- 数据模型\\n│   ├── req             -- 请求 VO 模型\\n│   ├── resp            -- 响应 VO 模型\\n│   ├── dto             -- 内部传输 DTO 模型\\n│   └── ...             -- 数据库模型对象 PO 模型\\n├── routes          -- 路由模块\\n│   └── middleware      -- 路由中间件\\n├── utils           -- 工具模块\\n│   ├── r               -- 响应封装\\n│   ├── upload          -- 文件上传\\n│   └── ...\\n├── routes          -- 路由模块\\n├── config          -- 配置文件\\n├── test            -- 测试模块\\n├── log             -- 日志文件\\n├── Dockerfile\\n└── main.go\\n```\\n\\n## 项目介绍\\n\\n前台：\\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\\n- 实现点赞，统计用户等功能 (Redis)\\n- 评论 + 回复评论功能\\n- 留言采用弹幕墙，效果炫酷\\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\\n\\n后台：\\n- 鉴权使用 JWT\\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\\n- 支持动态权限修改，前端菜单由后端生成\\n- 文章编辑使用 Markdown 编辑器\\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\\n\\n其他：\\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\\n- 代码整洁层次清晰，利于开发者学习\\n\\n## 技术介绍\\n\\n> 这里只写一些主流的通用技术，详细第三方库: 前端参考 `package.json` 文件，后端参考 `go.mod` 文件\\n\\n前台前端：使用 pnpm 包管理工具\\n- 基于 TypeScript\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia\\n- Vue Router \\n- Axios \\n- Naive UI\\n- Vuetify\\n- ...\\n\\n后台前端：使用 pnpm 包管理工具\\n- 基于 JavaSciprt \\n- pnpm: 包管理工具\\n- Vue3\\n- Unocss: 原子化 CSS\\n- Pinia \\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后端技术栈：\\n- 基于 Golang\\n- Docker\\n- Gin\\n- GORM\\n- Viper: 使用 TOML 作为配置文件\\n- Zap\\n- MySQL\\n- Redis\\n- Nginx: 部署静态资源 + 反向代理\\n- ...\\n\\n其他：\\n- 腾讯云人机验证\\n- 七牛云对象存储\\n- ...\\n\\n## 运行环境\\n\\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\\n\\n对象存储：七牛云\\n\\n## 开发环境\\n\\n| 开发工具                      | 说明                    |\\n| ----------------------------- | ----------------------- |\\n| Vscode                        | Golang 后端 +  Vue 前端 |\\n| Navicat                       | MySQL 远程连接工具      |\\n| Another Redis Desktop Manager | Redis 远程连接工具      |\\n| MobaXterm                     | Linux 远程工具          |\\n| Apifox                        | 接口调试 + 文档生成     |\\n\\n\\n| 开发环境 | 版本 |\\n| -------- | ---- |\\n| Golang   | 1.19 |\\n| MySQL    | 8.x  |\\n| Redis    | 7.x  |\\n\\n## 快速开始\\n\\n### 本地运行\\n\\n> 自行安装 Golang、Node、MySQL、Redis 环境\\n\\n拉取项目到本地：\\n\\n```bash\\ngit clone https://github.com/szluyu99/gin-vue-blog.git\\n```\\n\\n后端项目运行：\\n\\n```bash\\n# 1、进入后端项目根目录 \\ncd gin-blog-server\\n\\n# 2、修改项目运行的配置文件，默认加载 config/config.toml \\n\\n# 3、MySQL 导入 ginblog.sql\\n\\n# 4、启动 Redis \\n\\n# 5、运行项目\\ngo mod tidy\\ngo run main.go\\n```\\n\\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\\n\\n```bash\\nnpm install -g pnpm\\n```\\n\\n前台前端：\\n\\n```bash\\n# 1、进入前台前端项目根目录\\ncd gin-blog-front\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n后台前端：\\n\\n```bash\\n# 1、进入后台前端项目根目录\\ncd gin-blog-admin\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n### 项目部署\\n\\n目前暂时不推荐将本博客部署上生产环境，因为还有太多功能未完善。\\n\\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\\n\\n等功能开发的差不多了，再专门针对部署写一篇文章。\\n\\n---\\n\\n这里简单介绍一下，有基础的同学可以自行折腾。\\n\\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\\n\\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\\n\\nDocker 打包成镜像指令：\\n\\n```bash\\ndocker build -t ginblog .\\n```\\n\\n> 以上只是简单说明，等功能大致完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\\n\\n## 项目总结\\n\\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\\n\\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\\n\\n## 后续计划\\n\\n高优先级: \\n\\n- 完善图片上传功能, 目前文件上传还没怎么处理\\n- 后台首页重新设计（目前没放什么内容）\\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\\n- 前台首页搜索文章（ElasticSearch 搜索）\\n- 博客文章导入导出 (.md 文件)\\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\\n- 后端日志切割\\n- 后台修改背景图片，博客配置等\\n- 相册\\n\\n后续有空安排上：\\n- 适配移动端\\n- 黑夜模式\\n- 前台收缩侧边信息功能\\n- 说说\\n- 音乐播放器\\n- 鼠标左击特效\\n- 看板娘\\n- 文章目录锚点跟随\\n- 第三方登录\\n- 评论时支持选择表情，参考 Valine\\n- 若干细节需要完善...\\n\\n## 更新日志\\n\\n描述规范定义，有以下几种行为 ACTION:\\n- `FIX`: 修复\\n- `REFINE`: 优化\\n- `COMPLETE`: 完善\\n- `ADD`: 新增\\n\\n---\\n\\n2022-12-15:\\n\\n博客后台：\\n- 优化 菜单栏和标签栏，点击标签自动展开对应菜单，点击菜单自动滚动到显示对应标签 ⭐\\n- 优化 使用响应式语法糖重构代码\\n- 优化 代码结构 + 注释\\n\\n---\\n\\n2022-12-14:\\n\\n博客前台：\\n- 优化 代码，去除 n-card 组件，使用自定义 css 实现卡片视图\\n- 优化 对滚动的监听操作，添加节流函数，提升性能 ⭐\\n- 优化 文章目录，根据当前滚动条自动高亮锚点 ⭐⭐\\n- 新增 vite 打包优化相关插件\\n\\n---\\n\\n2022-12-13:\\n\\n项目部署：\\n- 新增 Nginx 配置 https 访问域名 (http 自动跳转到 https) ⭐\\n- 新增 七牛云添加加速域名访问图片资源\\n\\n博客后端：\\n- 新增 文章搜索接口（数据库模糊查询）\\n\\n博客前台：\\n- 新增 导航栏的文章搜索\\n\\n---\\n\\n2022-12-12:\\n\\n博客前台：\\n- 新增 适配移动端 ⭐⭐\\n- 优化 删除 Vuetify 相关组件及依赖 (css 样式存在冲突)，统一使用 naive-ui\\n- 优化 使用 `$ref` 语法糖功能重构页面\\n\\n---\\n\\n2022-12-09:\\n\\n博客前台：\\n- 完善 个人中心的头像上传（注意头像上传需要 Token）\\n- 优化 重构通用页面的代码 ⭐\\n\\n---\\n\\n2022-12-08:\\n\\n博客后端：\\n- 新增 页面模块 的接口 ⭐\\n- 修复 单元测试无法初始化全局变量（flag 包冲突）问题\\n\\n博客后台：\\n- 新增 pinia 数据持久化，防止刷新丢失数据 ⭐\\n- 新增 页面管理页面，动态配置前台显示的封面 ⭐\\n\\n博客前台：\\n- 新增 页面根据 `label` 从后端数据中加载封面 ⭐\\n- 优化 尝试性使用 `$ref` 语法糖功能\\n\\n---\\n\\n2022-12-07:\\n\\n博客后台：\\n- 修复 博客后台，动态加载路由模块导致的热加载失效问题 ⭐⭐\\n- 完善 文件上传，抽取出单独的组件 ⭐⭐\\n- 新增 个人信息页面 ⭐\\n- 优化 发布文章/查看文章 时的文章标签选择、分类选择\\n- 修复 发布文章/查看文章 页面数据加载错误\\n\",\"img\":\"https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg\",\"category_name\":\"后端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"original_url\":\"\",\"is_top\":0,\"status\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:35424', '未知');
INSERT INTO `operation_log` VALUES (39, '2022-12-16 11:53:37.174', '2022-12-16 11:53:37.174', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"id\":1,\"title\":\"测试文章\",\"desc\":\"\",\"content\":\"## 博客介绍\\n\\n<p  align=center>\\n<a  href=\\\"http://www.hahacode.cn\\\">\\n<img  src=\\\"https://img-blog.csdnimg.cn/fe2f1034cf7c4bd795552d47373ee405.jpeg\\\"  width=\\\"200\\\"  hight=\\\"200\\\"  alt=\\\"阵、雨的个人博客\\\"  style=\\\"border-radius: 50%\\\">\\n</a>\\n</p>\\n\\n因最近忙于学业，本项目开发周期不是很长且断断续续，可能会存在不少 BUG，但是我会逐步修复的。\\n\\n您的 Star 是我坚持的动力，感谢大家的支持，欢迎提交 Pr 共同改进项目。\\n\\nGithub 地址：[https://github.com/szluyu99/gin-vue-blog](https://github.com/szluyu99/gin-vue-blog)\\n\\nGitee 地址：[https://gitee.com/szluyu99/gin-vue-blog](https://gitee.com/szluyu99/gin-vue-blog)\\n\\n## 在线预览\\n\\n博客前台链接：[hahacode.cn](https://www.hahacode.cn)\\n\\n博客后台链接：[hahacode.cn/blog-admin](https://www.hahacode.cn/blog-admin)\\n\\n> 博客已备案通过，且配置 SSL，可通过 https 访问\\n\\n测试账号：test@qq.com，密码：11111，前后台都可用这个账号登录\\n\\n在线接口文档地址：[doc.hahacode.cn](http://doc.hahacode.cn/)\\n\\n> 本项目在线接口文档由 Apifox 生成，由于项目架构调整，有些接口待完善和修改\\n>\\n> 并且由于一开始不会用 Apifox，接口文档生成的信息不全（如返回响应、响应示例），抽空完善一下\\n\\n## 目录结构\\n\\n代码仓库目录：\\n\\n```bash\\ngin-vue-blog\\n├── gin-blog-admin      -- 博客后台前端\\n├── gin-blog-front      -- 博客前台前端\\n├── gin-blog-server     -- 博客后端\\n```\\n\\n> 项目运行参考：[快速开始](#快速开始)\\n\\n后端目录：\\n\\n```bash\\ngin-blog-server\\n├── api             -- API\\n│   ├── front       -- 前台接口\\n│   └── v1          -- 后台接口\\n├── dao             -- 数据库操作模块\\n├── service         -- 服务模块\\n├── model           -- 数据模型\\n│   ├── req             -- 请求 VO 模型\\n│   ├── resp            -- 响应 VO 模型\\n│   ├── dto             -- 内部传输 DTO 模型\\n│   └── ...             -- 数据库模型对象 PO 模型\\n├── routes          -- 路由模块\\n│   └── middleware      -- 路由中间件\\n├── utils           -- 工具模块\\n│   ├── r               -- 响应封装\\n│   ├── upload          -- 文件上传\\n│   └── ...\\n├── routes          -- 路由模块\\n├── config          -- 配置文件\\n├── test            -- 测试模块\\n├── log             -- 日志文件\\n├── Dockerfile\\n└── main.go\\n```\\n\\n前端目录：自行参考代码文件\\n\\n## 项目介绍\\n\\n市面上有太多优秀的前后台框架，本项目也是参考了很多开源项目，但是大多项目过于重量级（并非坏处），如果从学习的角度来看对初学者并不是很友好。本项目在以**博客**这个业务为主的前提下，提供一个完整的全栈项目代码（前台前端 + 后台前端 + 后端），技术点基本都是最新 + 最火的技术，代码轻量级，注释完善，适合学习。\\n\\n前台：\\n- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁\\n- 响应式布局，适配了移动端\\n- 实现点赞，统计用户等功能 (Redis)\\n- 评论 + 回复评论功能\\n- 留言采用弹幕墙，效果炫酷\\n- 文章详情页有文章目录、推荐文章等功能，优化用户体验\\n\\n后台：\\n- 鉴权使用 JWT\\n- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理\\n- 支持动态权限修改，前端菜单由后端生成（动态路由）\\n- 文章编辑使用 Markdown 编辑器\\n- 常规后台功能齐全：侧边栏、面包屑、标签栏等\\n- 实现记录操作日志功能（GET 不记录）\\n- 实现监听在线用户、强制下线功能\\n- 文件上传支持七牛云、本地（后续计划支持更多）\\n- 对 CRUD 操作封装了通用 Hook\\n\\n其他：\\n- 采用 Restful 风格的 API\\n- 前后端分离部署，前端使用 Nginx，后端使用 Docker\\n- 代码整洁层次清晰，利于开发者学习\\n- 技术点新颖，代码轻量级，适度封装\\n\\n## 技术介绍\\n\\n> 这里只写一些主流的通用技术，详细第三方库：前端参考 `package.json` 文件，后端参考 `go.mod` 文件\\n\\n前台前端：使用 pnpm 包管理工具\\n- 基于 TypeScript\\n- Vue3\\n- VueUse: 服务于 Vue Composition API 的工具集\\n- Unocss: 原子化 CSS\\n- Pinia\\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后台前端：使用 pnpm 包管理工具\\n- 基于 JavaSciprt \\n- Vue3\\n- VueUse: 服务于 Vue Composition API 的工具集\\n- Unocss: 原子化 CSS\\n- Pinia \\n- Vue Router \\n- Axios \\n- Naive UI\\n- ...\\n\\n后端技术栈：\\n- 基于 Golang\\n- Docker\\n- Gin\\n- GORM\\n- Viper: 使用 TOML 作为配置文件\\n- Casbin\\n- Zap\\n- MySQL\\n- Redis\\n- Nginx: 部署静态资源 + 反向代理\\n- ...\\n\\n其他：\\n- 腾讯云人机验证\\n- 七牛云对象存储\\n- ...\\n\\n## 运行环境\\n\\n服务器：腾讯云 2核 4G Ubuntu 22.04 LTS\\n\\n对象存储：七牛云\\n\\n## 开发环境\\n\\n| 开发工具                      | 说明                    |\\n| ----------------------------- | ----------------------- |\\n| Vscode                        | Golang 后端 +  Vue 前端 |\\n| Navicat                       | MySQL 远程连接工具      |\\n| Another Redis Desktop Manager | Redis 远程连接工具      |\\n| MobaXterm                     | Linux 远程工具          |\\n| Apifox                        | 接口调试 + 文档生成     |\\n\\n| 开发环境 | 版本 |\\n| -------- | ---- |\\n| Golang   | 1.19 |\\n| MySQL    | 8.x  |\\n| Redis    | 7.x  |\\n\\n\\n### Vscode 插件\\n\\n如果使用 Vscode 开发，推荐安装一下以下插件。\\n\\n> 前后端相关插件属于**开发必须插件**，还有很多提升开发效果的插件后面推荐一下\\n\\n前端开发插件：\\n\\n| 插件 | 作用 |\\n| -------- | ---- |\\n| Volar   | Vue 官方插件 |\\n| UnoCSS | Unocss 官方插件 |\\n| Tailwind CSS IntelliSense | Tailwind 官方插件 |\\n| Iconify IntelliSense | Iconify 提示插件 |\\n\\n后端开发插件：\\n\\n| 插件 | 作用 |\\n| -------- | ---- |\\n| Go | Golang 官方插件 |\\n\\n\\n通用插件：（非开发必须，个人推荐）\\n\\n| 名称 | 作用 |\\n| -------- | ---- |\\n| Better Comments   | 优化代码注释显示效果 |\\n| TODO Highlight | 高亮 TODO |\\n| Highlight Matching Tag | 高亮匹配的标签 | \\n\\n\\n## 快速开始\\n\\n### 本地运行\\n\\n> 目前需要自行安装 Golang、Node、MySQL、Redis 环境\\n>\\n> TODO: 完全基于 Docker 的运行教程（Docker 实在太好用了！）\\n\\n\\n需要先运行后端服务，再运行前端项目，因为很多前端配置由后端动态加载（如菜单等）。\\n\\n拉取项目到本地：\\n\\n```bash\\ngit clone https://github.com/szluyu99/gin-vue-blog.git\\n```\\n\\n后端项目运行：\\n\\n```bash\\n# 1、进入后端项目根目录 \\ncd gin-blog-server\\n\\n# 2、修改项目运行的配置文件，默认加载位于 config/config.toml \\n\\n# 3、MySQL 导入 ginblog.sql\\n\\n# 4、启动 Redis \\n\\n# 5、运行项目\\ngo mod tidy\\ngo run main.go\\n```\\n\\n数据库中的默认用户：\\n- 管理员 admin 123456\\n- 普通用户 user 123456\\n- 测试用户 test 123456\\n\\n前端项目运行： 本项目使用 pnpm 进行包管理，建议全局安装 pnpm\\n\\n```bash\\nnpm install -g pnpm\\n```\\n\\n前台前端：\\n\\n```bash\\n# 1、进入前台前端项目根目录\\ncd gin-blog-front\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n后台前端：\\n\\n```bash\\n# 1、进入后台前端项目根目录\\ncd gin-blog-admin\\n\\n# 2、安装依赖\\npnpm install\\n\\n# 3、运行项目\\npnpm run dev\\n```\\n\\n\\n### 项目部署\\n\\n**目前暂时不推荐将本博客部署上生产环境，因为还有很多功能未完善**。\\n\\n但是相信本项目对于 Golang 学习者绝对是个合适的项目！\\n\\n等功能开发的差不多了，再专门针对部署写一篇文章。\\n\\n---\\n\\n这里简单介绍一下，有基础的同学可以自行折腾。\\n\\n本项目前端采用 Nginx 部署静态资源，后端使用 Docker 部署。\\n\\n后端 Docker 部署参考 `Dockerfile`，Docker 运行对应的配置文件是 `config/config.docker.toml`\\n\\nDocker 打包成镜像指令：\\n\\n```bash\\ndocker build -t ginblog .\\n```\\n\\n> 以上只是简单说明，等功能全部完成，会从 `安装 Docker`、`Docker 安装运行环境`、`Docker 部署项目` 等多个角度写几篇关于部署的教程。\\n\\n## 项目总结\\n\\n这个项目不管是前端，还是后端，都是花了比较大心血去架构的，并且从技术选型上，都是选择了目前最火 + 最新的技术栈。当然，这也是个人的学习之作，很多知识点都是边学边开发的（例如 Casbin），这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加美好，这也是开源本项目的目的之一。本项目中仍有很多不足，后续会继续更新。\\n\\n最后，项目整体代码风格很优秀，注释完善，适合 Golang 后端开发者、前端开发者学习。\\n\\n## 鸣谢项目\\n\\n- [https://butterfly.js.org/](https://butterfly.js.org/)\\n- [https://github.com/flipped-aurora/gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)\\n- [https://github.com/qifengzhang007/GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)\\n- [https://github.com/X1192176811/blog](https://github.com/X1192176811/blog)\\n- [https://github.com/zclzone/vue-naive-admin](https://github.com/zclzone/vue-naive-admin)\\n- [https://github.com/antfu/vitesse](https://github.com/antfu/vitesse)\\n- ...\\n\\n> 需要感谢的绝不止以上这些开源项目，但是一时难以全部列出，后面会慢慢补上。\\n\\n## 后续计划\\n\\n高优先级: \\n\\n- ~~完善图片上传功能, 目前文件上传还没怎么处理~~\\n- 后台首页重新设计（目前没放什么内容）\\n- ~~前台首页搜索文章（目前使用数据库模糊搜索）~~\\n- 博客文章导入导出 (.md 文件)\\n- 权限管理中菜单编辑时选择图标（现在只能输入图标字符串）\\n- 后端日志切割\\n- ~~后台修改背景图片，博客配置等~~ \\n- 后端的 IP 地址检测 BUG 待修复\\n- ~~博客前台适配移动端~~ \\n- ~~文章详情, 目录锚点跟随~~ \\n\\n后续有空安排上：\\n- 黑夜模式\\n- 前台收缩侧边信息功能\\n- 说说\\n- 相册\\n- 音乐播放器\\n- 鼠标左击特效\\n- 看板娘\\n- 第三方登录\\n- 评论时支持选择表情，参考 Valine\\n- 一键部署：使用 docker compose 单机一键部署整个项目（前后端 + 环境）\\n- 单独部署：前后端 + 环境\\n- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通\\n- 前台首页搜索集成 ElasticSearch\\n\\n其他：\\n- 写一份好的文档\\n- 补全 README.md\\n- 完善 Apifox 生成的接口文档\\n\\n## 更新日志\\n\\n描述规范定义，有以下几种行为 ACTION:\\n- `FIX`: 修复\\n- `REFINE`: 优化\\n- `COMPLETE`: 完善\\n- `ADD`: 新增\\n\\n---\\n\\n2022-12-15:\\n\\n博客后台：\\n- 优化 菜单栏和标签栏，点击标签自动展开对应菜单，点击菜单自动滚动到显示对应标签 ⭐\\n- 优化 使用响应式语法糖重构代码\\n- 优化 代码结构 + 注释\\n\\n---\\n\\n2022-12-14:\\n\\n博客前台：\\n- 优化 代码，去除 n-card 组件，使用自定义 css 实现卡片视图\\n- 优化 对滚动的监听操作，添加节流函数，提升性能 ⭐\\n- 优化 文章目录，根据当前滚动条自动高亮锚点 ⭐⭐\\n- 新增 vite 打包优化相关插件\\n\\n---\\n\\n2022-12-13:\\n\\n项目部署：\\n- 新增 Nginx 配置 https 访问域名 (http 自动跳转到 https) ⭐\\n- 新增 七牛云添加加速域名访问图片资源\\n\\n博客后端：\\n- 新增 文章搜索接口（数据库模糊查询）\\n\\n博客前台：\\n- 新增 导航栏的文章搜索\\n\\n---\\n\\n2022-12-12:\\n\\n博客前台：\\n- 新增 适配移动端 ⭐⭐\\n- 优化 删除 Vuetify 相关组件及依赖 (css 样式存在冲突)，统一使用 naive-ui\\n- 优化 使用 `$ref` 语法糖功能重构页面\\n\\n---\\n\\n2022-12-09:\\n\\n博客前台：\\n- 完善 个人中心的头像上传（注意头像上传需要 Token）\\n- 优化 重构通用页面的代码 ⭐\\n\\n---\\n\\n2022-12-08:\\n\\n博客后端：\\n- 新增 页面模块 的接口 ⭐\\n- 修复 单元测试无法初始化全局变量（flag 包冲突）问题\\n\\n博客后台：\\n- 新增 pinia 数据持久化，防止刷新丢失数据 ⭐\\n- 新增 页面管理页面，动态配置前台显示的封面 ⭐\\n\\n博客前台：\\n- 新增 页面根据 `label` 从后端数据中加载封面 ⭐\\n- 优化 尝试性使用 `$ref` 语法糖功能\\n\\n---\\n\\n2022-12-07:\\n\\n博客后台：\\n- 修复 博客后台，动态加载路由模块导致的热加载失效问题 ⭐⭐\\n- 完善 文件上传，抽取出单独的组件 ⭐⭐\\n- 新增 个人信息页面 ⭐\\n- 优化 发布文章/查看文章 时的文章标签选择、分类选择\\n- 修复 发布文章/查看文章 页面数据加载错误\",\"img\":\"https://static.talkxj.com/articles/771941739cbc70fbe40e10cf441e02e5.jpg\",\"category_name\":\"后端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"original_url\":\"\",\"is_top\":0,\"status\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:35424', '未知');
INSERT INTO `operation_log` VALUES (40, '2022-12-16 11:53:58.004', '2022-12-16 11:53:58.004', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"name\":\"文件模块\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:35424', '未知');
INSERT INTO `operation_log` VALUES (41, '2022-12-16 11:54:20.907', '2022-12-16 11:54:20.907', '资源权限', '新增或修改', 'gin-blog/api/v1.(*Resource).SaveOrUpdate-fm', '/api/resource', '新增或修改资源权限', '{\"parent_id\":106,\"name\":\"文件上传\",\"url\":\"/upload\",\"request_method\":\"POST\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:35424', '未知');
INSERT INTO `operation_log` VALUES (42, '2022-12-16 11:54:27.916', '2022-12-16 11:54:27.916', '角色', '新增或修改', 'gin-blog/api/v1.(*Role).SaveOrUpdate-fm', '/api/role', '新增或修改角色', '{\"id\":1,\"name\":\"管理员\",\"label\":\"admin\",\"created_at\":\"2022-10-20T21:24:28+08:00\",\"is_disable\":0,\"resource_ids\":[3,43,44,45,46,47,48,6,59,60,61,7,38,39,40,41,42,8,65,66,67,68,9,62,63,64,10,23,34,35,36,37,79,80,81,85,49,51,52,53,54,50,55,56,57,58,69,70,71,72,91,92,93,78,82,84,86,98,95,11,99,100,101,74,102,103,104,105,106,107],\"menu_ids\":[2,3,4,5,6,7,8,9,16,17,23,24,25,26,27,28,29,30,31,32,33,36,37,38,34,10,39]}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:35424', '未知');
INSERT INTO `operation_log` VALUES (43, '2022-12-16 11:56:40.911', '2022-12-16 11:56:40.911', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"status\":1,\"is_top\":0,\"title\":\"项目运行成功\",\"content\":\"恭喜你，项目已经成功运行起来了！\\n\\n```go\\nfmt.Println(\\\"恭喜！\\\")\\n```\\n\\n```js\\nconsole.log(\\\"恭喜！\\\")\\n```\",\"category_name\":\"前端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:35424', '未知');
INSERT INTO `operation_log` VALUES (44, '2022-12-16 12:36:42.125', '2022-12-16 12:36:42.125', '文章', '新增或修改', 'gin-blog/api/v1.(*Article).SaveOrUpdate-fm', '/api/article', '新增或修改文章', '{\"id\":2,\"title\":\"项目运行成功\",\"desc\":\"\",\"content\":\"恭喜你，项目已经成功运行起来了！\\n\\n```go\\nfmt.Println(\\\"恭喜！\\\")\\n```\\n\\n```js\\nconsole.log(\\\"恭喜！\\\")\\n```\\n\\n🆗😋\",\"img\":\"https://static.talkxj.com/articles/1395642324821741569.png\",\"category_name\":\"前端\",\"tag_names\":[\"Golang\",\"Vue\"],\"type\":1,\"original_url\":\"\",\"is_top\":0,\"status\":1}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:35424', '未知');
INSERT INTO `operation_log` VALUES (45, '2022-12-16 23:39:05.800', '2022-12-16 23:39:05.800', '博客信息', '修改', 'gin-blog/api/v1.(*BlogInfo).UpdateBlogConfig-fm', '/api/setting/blog-config', '修改博客信息', '{\"website_avatar\":\"public/uploaded/d0b80b7abde6b30fae6e950eb1ebf11a_20221216171019.png\",\"website_name\":\"阵、雨的个人博客\",\"website_author\":\"阵、雨\",\"website_intro\":\"往事随风而去\",\"website_notice\":\"博客后端基于 gin、gorm 开发\\n博客前端基于 Vue3、TS、NaiveUI 开发\\n努力开发中...冲冲冲！加油！\",\"website_createtime\":\"2022-11-01\",\"website_record\":\"qui aute sint ea dolor\",\"social_login_list\":[],\"social_url_list\":[],\"qq\":\"123456\",\"github\":\"https://github.com/szluyu99\",\"gitee\":\"https://gitee.com/szluyu99\",\"tourist_avatar\":\"https://static.talkxj.com/photos/0bca52afdb2b9998132355d716390c9f.png\",\"user_avatar\":\"https://static.talkxj.com/config/2cd793c8744199053323546875655f32.jpg\",\"article_cover\":\"https://static.talkxj.com/config/e587c4651154e4da49b5a54f865baaed.jpg\",\"is_comment_review\":1,\"is_message_review\":1,\"is_email_notice\":0,\"is_reward\":0,\"wechat_qrcode\":\"http://dummyimage.com/100x100\",\"alipay_ode\":\"http://dummyimage.com/100x100\"}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');
INSERT INTO `operation_log` VALUES (46, '2022-12-16 23:39:30.860', '2022-12-16 23:39:30.860', '博客信息', '修改', 'gin-blog/api/v1.(*BlogInfo).UpdateBlogConfig-fm', '/api/setting/blog-config', '修改博客信息', '{\"website_avatar\":\"public/uploaded/d0b80b7abde6b30fae6e950eb1ebf11a_20221216171019.png\",\"website_name\":\"阵、雨的个人博客\",\"website_author\":\"阵、雨\",\"website_intro\":\"往事随风而去\",\"website_notice\":\"博客后端基于 gin、gorm 开发\\n博客前端基于 Vue3、TS、NaiveUI 开发\\n努力开发中...冲冲冲！加油！\",\"website_createtime\":\"2022-11-01\",\"website_record\":\"鲁ICP备2022040119号\",\"social_login_list\":[],\"social_url_list\":[],\"qq\":\"123456\",\"github\":\"https://github.com/szluyu99\",\"gitee\":\"https://gitee.com/szluyu99\",\"tourist_avatar\":\"https://static.talkxj.com/photos/0bca52afdb2b9998132355d716390c9f.png\",\"user_avatar\":\"https://static.talkxj.com/config/2cd793c8744199053323546875655f32.jpg\",\"article_cover\":\"https://static.talkxj.com/config/e587c4651154e4da49b5a54f865baaed.jpg\",\"is_comment_review\":1,\"is_message_review\":1,\"is_email_notice\":0,\"is_reward\":0,\"wechat_qrcode\":\"http://dummyimage.com/100x100\",\"alipay_ode\":\"http://dummyimage.com/100x100\"}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');
INSERT INTO `operation_log` VALUES (47, '2022-12-16 23:42:32.129', '2022-12-16 23:42:32.129', '用户', '修改', 'gin-blog/api/v1.(*User).UpdateCurrent-fm', '/api/user/current', '修改用户', '{\"avatar\":\"https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png\",\"nickname\":\"管理员\",\"intro\":\"我是管理员用户！\",\"website\":\"https://www.hahacode.cn\"}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');
INSERT INTO `operation_log` VALUES (48, '2022-12-16 23:46:51.072', '2022-12-16 23:46:51.072', '页面', '新增或修改', 'gin-blog/api/v1.(*Page).SaveOrUpdate-fm', '/api/page', '新增或修改页面', '{\"id\":7,\"created_at\":\"2022-12-08T13:53:17.707+08:00\",\"updated_at\":\"2022-12-08T13:53:17.707+08:00\",\"name\":\"留言\",\"label\":\"message\",\"cover\":\"https://img-blog.csdnimg.cn/50e0ec7fde824633b2cd7870028670b2.jpeg\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');
INSERT INTO `operation_log` VALUES (49, '2022-12-16 23:52:38.392', '2022-12-16 23:52:38.392', '博客信息', '修改', 'gin-blog/api/v1.(*BlogInfo).UpdateBlogConfig-fm', '/api/setting/blog-config', '修改博客信息', '{\"website_avatar\":\"public/uploaded/d0b80b7abde6b30fae6e950eb1ebf11a_20221216171019.png\",\"website_name\":\"阵、雨的个人博客\",\"website_author\":\"阵、雨\",\"website_intro\":\"往事随风而去\",\"website_notice\":\"博客后端基于 gin、gorm 开发\\n博客前端基于 Vue3、TS、NaiveUI 开发\\n努力开发中...冲冲冲！加油！\",\"website_createtime\":\"2022-11-01\",\"website_record\":\"鲁ICP备2022040119号\",\"social_login_list\":[],\"social_url_list\":[],\"qq\":\"123456\",\"github\":\"https://github.com/szluyu99\",\"gitee\":\"https://gitee.com/szluyu99\",\"tourist_avatar\":\"https://static.talkxj.com/photos/0bca52afdb2b9998132355d716390c9f.png\",\"user_avatar\":\"https://static.talkxj.com/config/2cd793c8744199053323546875655f32.jpg\",\"article_cover\":\"https://static.talkxj.com/config/e587c4651154e4da49b5a54f865baaed.jpg\",\"is_comment_review\":1,\"is_message_review\":1,\"is_email_notice\":0,\"is_reward\":0,\"wechat_qrcode\":\"http://dummyimage.com/100x100\",\"alipay_ode\":\"http://dummyimage.com/100x100\"}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');
INSERT INTO `operation_log` VALUES (50, '2022-12-16 23:54:52.660', '2022-12-16 23:54:52.660', '页面', '新增或修改', 'gin-blog/api/v1.(*Page).SaveOrUpdate-fm', '/api/page', '新增或修改页面', '{\"name\":\"相册\",\"label\":\"album\",\"cover\":\"https://static.talkxj.com/config/e587c4651154e4da49b5a54f865baaed.jpg\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');
INSERT INTO `operation_log` VALUES (51, '2022-12-16 23:55:36.065', '2022-12-16 23:55:36.065', '页面', '新增或修改', 'gin-blog/api/v1.(*Page).SaveOrUpdate-fm', '/api/page', '新增或修改页面', '{\"cover\":\"https://img-blog.csdnimg.cn/2c0e923329974daabb91373d0834359f.jpeg\",\"name\":\"错误页面\",\"label\":\"404\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');
INSERT INTO `operation_log` VALUES (52, '2022-12-16 23:56:17.923', '2022-12-16 23:56:17.923', '页面', '新增或修改', 'gin-blog/api/v1.(*Page).SaveOrUpdate-fm', '/api/page', '新增或修改页面', '{\"name\":\"文章列表\",\"label\":\"article_list\",\"cover\":\"https://static.talkxj.com/config/924d65cc8312e6cdad2160eb8fce6831.jpg\"}', 'POST', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');
INSERT INTO `operation_log` VALUES (53, '2022-12-17 00:02:41.144', '2022-12-17 00:02:41.144', '博客信息', '修改', 'gin-blog/api/v1.(*BlogInfo).UpdateBlogConfig-fm', '/api/setting/blog-config', '修改博客信息', '{\"website_avatar\":\"https://img-blog.csdnimg.cn/9e1832a8817d4344901f2e476cc74c16.jpeg\",\"website_name\":\"阵、雨的个人博客\",\"website_author\":\"阵、雨\",\"website_intro\":\"往事随风而去\",\"website_notice\":\"博客后端基于 gin、gorm 开发\\n博客前端基于 Vue3、TS、NaiveUI 开发\\n努力开发中...冲冲冲！加油！\",\"website_createtime\":\"2022-11-01\",\"website_record\":\"鲁ICP备2022040119号\",\"social_login_list\":[],\"social_url_list\":[],\"qq\":\"123456\",\"github\":\"https://github.com/szluyu99\",\"gitee\":\"https://gitee.com/szluyu99\",\"tourist_avatar\":\"https://static.talkxj.com/photos/0bca52afdb2b9998132355d716390c9f.png\",\"user_avatar\":\"https://static.talkxj.com/avatar/user.png\",\"article_cover\":\"https://static.talkxj.com/config/e587c4651154e4da49b5a54f865baaed.jpg\",\"is_comment_review\":1,\"is_message_review\":1,\"is_email_notice\":0,\"is_reward\":0,\"wechat_qrcode\":\"http://dummyimage.com/100x100\",\"alipay_ode\":\"http://dummyimage.com/100x100\"}', 'PUT', '{\"code\":0,\"message\":\"OK\",\"data\":null}', 1, '管理员', '127.0.0.1:41990', '未知');

-- ----------------------------
-- Table structure for page
-- ----------------------------
DROP TABLE IF EXISTS `page`;
CREATE TABLE `page`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '页面名称',
  `label` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '页面标签',
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '页面封面',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of page
-- ----------------------------
INSERT INTO `page` VALUES (1, '2022-12-08 13:09:58.500', '2022-12-08 15:08:38.248', '首页', 'home', 'https://static.talkxj.com/config/0bee7ba5ac70155766648e14ae2a821f.jpg');
INSERT INTO `page` VALUES (2, '2022-12-08 13:51:49.474', '2022-12-08 13:51:49.474', '归档', 'archive', 'https://static.talkxj.com/config/643f28683e1c59a80ccfc9cb19735a9c.jpg');
INSERT INTO `page` VALUES (3, '2022-12-08 13:52:18.084', '2022-12-08 13:52:18.084', '分类', 'category', 'https://static.talkxj.com/config/83be0017d7f1a29441e33083e7706936.jpg');
INSERT INTO `page` VALUES (4, '2022-12-08 13:52:31.364', '2022-12-08 13:52:31.364', '标签', 'tag', 'https://static.talkxj.com/config/a6f141372509365891081d755da963a1.png');
INSERT INTO `page` VALUES (5, '2022-12-08 13:52:52.389', '2022-12-08 13:52:52.389', '友链', 'link', 'https://static.talkxj.com/config/9034edddec5b8e8542c2e61b0da1c1da.jpg');
INSERT INTO `page` VALUES (6, '2022-12-08 13:53:04.159', '2022-12-08 13:53:04.159', '关于', 'about', 'https://static.talkxj.com/config/2a56d15dd742ff8ac238a512d9a472a1.jpg');
INSERT INTO `page` VALUES (7, '2022-12-08 13:53:17.707', '2022-12-16 23:46:51.066', '留言', 'message', 'https://img-blog.csdnimg.cn/50e0ec7fde824633b2cd7870028670b2.jpeg');
INSERT INTO `page` VALUES (8, '2022-12-08 13:53:30.187', '2022-12-08 13:53:30.187', '个人中心', 'profile', 'https://static.talkxj.com/config/ebae4c93de1b286a8d50aa62612caa59.jpeg');
INSERT INTO `page` VALUES (9, '2022-12-16 23:54:52.650', '2022-12-16 23:54:52.650', '相册', 'album', 'https://static.talkxj.com/config/e587c4651154e4da49b5a54f865baaed.jpg');
INSERT INTO `page` VALUES (10, '2022-12-16 23:55:36.059', '2022-12-16 23:55:36.059', '错误页面', '404', 'https://img-blog.csdnimg.cn/2c0e923329974daabb91373d0834359f.jpeg');
INSERT INTO `page` VALUES (11, '2022-12-16 23:56:17.917', '2022-12-16 23:56:17.917', '文章列表', 'article_list', 'https://static.talkxj.com/config/924d65cc8312e6cdad2160eb8fce6831.jpg');

-- ----------------------------
-- Table structure for resource
-- ----------------------------
DROP TABLE IF EXISTS `resource`;
CREATE TABLE `resource`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '资源路径(接口URL)',
  `request_method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '请求方式',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '资源名(接口名)',
  `parent_id` bigint NULL DEFAULT NULL COMMENT '父权限id',
  `is_anonymous` tinyint(1) NULL DEFAULT NULL COMMENT '是否匿名访问(0-否 1-是)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 108 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of resource
-- ----------------------------
INSERT INTO `resource` VALUES (3, '2022-10-20 22:42:00.664', '2022-10-20 22:42:00.664', '', '', '文章模块', 0, 0);
INSERT INTO `resource` VALUES (6, '2022-10-20 22:42:23.349', '2022-10-20 22:42:23.349', '', '', '留言模块', 0, 0);
INSERT INTO `resource` VALUES (7, '2022-10-20 22:42:28.550', '2022-10-20 22:42:28.550', '', '', '菜单模块', 0, 0);
INSERT INTO `resource` VALUES (8, '2022-10-20 22:42:31.623', '2022-10-20 22:42:31.623', '', '', '角色模块', 0, 0);
INSERT INTO `resource` VALUES (9, '2022-10-20 22:42:36.262', '2022-10-20 22:42:36.262', '', '', '评论模块', 0, 0);
INSERT INTO `resource` VALUES (10, '2022-10-20 22:42:40.700', '2022-10-20 22:42:40.700', '', '', '资源模块', 0, 0);
INSERT INTO `resource` VALUES (11, '2022-10-20 22:42:51.023', '2022-10-20 22:42:51.023', '', '', '博客信息模块', 0, 0);
INSERT INTO `resource` VALUES (23, '2022-10-22 22:13:12.455', '2022-10-26 11:15:23.546', '/resource/anonymous', 'PUT', '修改资源匿名访问', 10, 0);
INSERT INTO `resource` VALUES (34, '2022-10-31 17:14:11.708', '2022-10-31 17:14:11.708', '/resource', 'POST', '新增/编辑资源', 10, 0);
INSERT INTO `resource` VALUES (35, '2022-10-31 17:14:42.320', '2022-10-31 17:18:52.508', '/resource/list', 'GET', '资源列表', 10, 0);
INSERT INTO `resource` VALUES (36, '2022-10-31 17:15:14.999', '2022-10-31 17:19:01.460', '/resource/option', 'GET', '资源选项列表(树形)', 10, 0);
INSERT INTO `resource` VALUES (37, '2022-10-31 17:16:56.830', '2022-10-31 17:16:56.830', '/resource/:id', 'DELETE', '删除资源', 10, 0);
INSERT INTO `resource` VALUES (38, '2022-10-31 17:19:28.905', '2022-10-31 17:19:28.905', '/menu/list', 'GET', '菜单列表', 7, 0);
INSERT INTO `resource` VALUES (39, '2022-10-31 18:46:33.051', '2022-10-31 18:46:33.051', '/menu', 'POST', '新增/编辑菜单', 7, 0);
INSERT INTO `resource` VALUES (40, '2022-10-31 18:46:53.804', '2022-10-31 18:46:53.804', '/menu/:id', 'DELETE', '删除菜单', 7, 0);
INSERT INTO `resource` VALUES (41, '2022-10-31 18:47:17.272', '2022-10-31 18:47:28.130', '/menu/option', 'GET', '菜单选项列表(树形)', 7, 0);
INSERT INTO `resource` VALUES (42, '2022-10-31 18:48:04.780', '2022-10-31 18:48:04.780', '/menu/user/list', 'GET', '获取当前用户菜单', 7, 0);
INSERT INTO `resource` VALUES (43, '2022-10-31 19:20:35.427', '2022-10-31 19:20:35.427', '/article/list', 'GET', '文章列表', 3, 0);
INSERT INTO `resource` VALUES (44, '2022-10-31 19:21:02.096', '2022-10-31 19:21:02.096', '/article/:id', 'GET', '文章详情', 3, 0);
INSERT INTO `resource` VALUES (45, '2022-10-31 19:26:04.763', '2022-10-31 19:26:09.709', '/article', 'POST', '新增/编辑文章', 3, 0);
INSERT INTO `resource` VALUES (46, '2022-10-31 19:26:36.453', '2022-10-31 19:26:36.453', '/article/soft-delete', 'PUT', '软删除文章', 3, 0);
INSERT INTO `resource` VALUES (47, '2022-10-31 19:26:52.344', '2022-10-31 19:26:52.344', '/article', 'DELETE', '删除文章', 3, 0);
INSERT INTO `resource` VALUES (48, '2022-10-31 19:27:07.731', '2022-10-31 19:27:07.731', '/article/top', 'PUT', '修改文章置顶', 3, 0);
INSERT INTO `resource` VALUES (49, '2022-10-31 20:19:55.588', '2022-10-31 20:19:55.588', '', '', '分类模块', 0, 0);
INSERT INTO `resource` VALUES (50, '2022-10-31 20:20:03.400', '2022-10-31 20:20:03.400', '', '', '标签模块', 0, 0);
INSERT INTO `resource` VALUES (51, '2022-10-31 20:22:03.799', '2022-10-31 20:22:03.799', '/category/list', 'GET', '分类列表', 49, 0);
INSERT INTO `resource` VALUES (52, '2022-10-31 20:22:28.840', '2022-10-31 20:22:28.840', '/category', 'POST', '新增/编辑分类', 49, 0);
INSERT INTO `resource` VALUES (53, '2022-10-31 20:31:04.577', '2022-10-31 20:31:04.577', '/category', 'DELETE', '删除分类', 49, 0);
INSERT INTO `resource` VALUES (54, '2022-10-31 20:31:36.612', '2022-10-31 20:31:36.612', '/category/option', 'GET', '分类选项列表', 49, 0);
INSERT INTO `resource` VALUES (55, '2022-10-31 20:32:57.112', '2022-10-31 20:33:13.261', '/tag/list', 'GET', '标签列表', 50, 0);
INSERT INTO `resource` VALUES (56, '2022-10-31 20:33:29.080', '2022-10-31 20:33:29.080', '/tag', 'POST', '新增/编辑标签', 50, 0);
INSERT INTO `resource` VALUES (57, '2022-10-31 20:33:39.992', '2022-10-31 20:33:39.992', '/tag', 'DELETE', '删除标签', 50, 0);
INSERT INTO `resource` VALUES (58, '2022-10-31 20:33:53.962', '2022-10-31 20:33:53.962', '/tag/option', 'GET', '标签选项列表', 50, 0);
INSERT INTO `resource` VALUES (59, '2022-10-31 20:35:05.647', '2022-10-31 20:35:05.647', '/message/list', 'GET', '留言列表', 6, 0);
INSERT INTO `resource` VALUES (60, '2022-10-31 20:35:25.551', '2022-10-31 20:35:25.551', '/message', 'DELETE', '删除留言', 6, 0);
INSERT INTO `resource` VALUES (61, '2022-10-31 20:36:20.587', '2022-10-31 20:36:20.587', '/message/review', 'PUT', '修改留言审核', 6, 0);
INSERT INTO `resource` VALUES (62, '2022-10-31 20:37:04.637', '2022-10-31 20:37:04.637', '/comment/list', 'GET', '评论列表', 9, 0);
INSERT INTO `resource` VALUES (63, '2022-10-31 20:37:29.779', '2022-10-31 20:37:29.779', '/comment', 'DELETE', '删除评论', 9, 0);
INSERT INTO `resource` VALUES (64, '2022-10-31 20:37:40.317', '2022-10-31 20:37:40.317', '/comment/review', 'PUT', '修改评论审核', 9, 0);
INSERT INTO `resource` VALUES (65, '2022-10-31 20:38:30.506', '2022-10-31 20:38:30.506', '/role/list', 'GET', '角色列表', 8, 0);
INSERT INTO `resource` VALUES (66, '2022-10-31 20:38:50.606', '2022-10-31 20:38:50.606', '/role', 'POST', '新增/编辑角色', 8, 0);
INSERT INTO `resource` VALUES (67, '2022-10-31 20:39:03.752', '2022-10-31 20:39:03.752', '/role', 'DELETE', '删除角色', 8, 0);
INSERT INTO `resource` VALUES (68, '2022-10-31 20:39:28.232', '2022-10-31 20:39:28.232', '/role/option', 'GET', '角色选项', 8, 0);
INSERT INTO `resource` VALUES (69, '2022-10-31 20:44:22.622', '2022-10-31 20:44:22.622', '', '', '友链模块', 0, 0);
INSERT INTO `resource` VALUES (70, '2022-10-31 20:44:41.334', '2022-10-31 20:44:41.334', '/link/list', 'GET', '友链列表', 69, 0);
INSERT INTO `resource` VALUES (71, '2022-10-31 20:45:01.150', '2022-10-31 20:45:01.150', '/link', 'POST', '新增/编辑友链', 69, 0);
INSERT INTO `resource` VALUES (72, '2022-10-31 20:45:12.406', '2022-10-31 20:45:12.406', '/link', 'DELETE', '删除友链', 69, 0);
INSERT INTO `resource` VALUES (74, '2022-10-31 20:46:48.330', '2022-10-31 20:47:01.505', '', '', '用户信息模块', 0, 0);
INSERT INTO `resource` VALUES (78, '2022-10-31 20:51:15.607', '2022-10-31 20:51:15.607', '/user/list', 'GET', '用户列表', 74, 0);
INSERT INTO `resource` VALUES (79, '2022-10-31 20:55:15.496', '2022-10-31 20:55:15.496', '/setting/blog-config', 'GET', '获取博客设置', 11, 0);
INSERT INTO `resource` VALUES (80, '2022-10-31 20:55:48.257', '2022-10-31 20:55:48.257', '/setting/about', 'GET', '获取关于我', 11, 0);
INSERT INTO `resource` VALUES (81, '2022-10-31 20:56:21.722', '2022-10-31 20:56:21.722', '/setting/blog-config', 'PUT', '修改博客设置', 11, 0);
INSERT INTO `resource` VALUES (82, '2022-10-31 21:57:30.021', '2022-10-31 21:57:30.021', '/user/info', 'GET', '获取当前用户信息', 74, 0);
INSERT INTO `resource` VALUES (84, '2022-10-31 22:06:18.348', '2022-10-31 22:07:38.004', '/user', 'PUT', '修改用户信息', 74, 0);
INSERT INTO `resource` VALUES (85, '2022-11-02 11:55:05.395', '2022-11-02 11:55:05.395', '/setting/about', 'PUT', '修改关于我', 11, 0);
INSERT INTO `resource` VALUES (86, '2022-11-02 13:20:09.485', '2022-11-02 13:20:09.485', '/user/online', 'GET', '获取在线用户列表', 74, 0);
INSERT INTO `resource` VALUES (91, '2022-11-03 16:42:31.558', '2022-11-03 16:42:31.558', '', '', '操作日志模块', 0, 0);
INSERT INTO `resource` VALUES (92, '2022-11-03 16:42:49.681', '2022-11-03 16:42:49.681', '/operation/log/list', 'GET', '获取操作日志列表', 91, 0);
INSERT INTO `resource` VALUES (93, '2022-11-03 16:43:04.906', '2022-11-03 16:43:04.906', '/operation/log', 'DELETE', '删除操作日志', 91, 0);
INSERT INTO `resource` VALUES (95, '2022-11-05 14:22:48.240', '2022-11-05 14:22:48.240', '/home', 'GET', '获取后台首页信息', 11, 0);
INSERT INTO `resource` VALUES (98, '2022-11-29 23:35:42.865', '2022-11-29 23:35:42.865', '/user/offline', 'DELETE', '强制离线用户', 74, 0);
INSERT INTO `resource` VALUES (99, '2022-12-07 20:48:05.939', '2022-12-07 20:48:05.939', '/user/current/password', 'PUT', '修改当前用户密码', 74, 0);
INSERT INTO `resource` VALUES (100, '2022-12-07 20:48:35.511', '2022-12-07 20:48:35.511', '/user/current', 'PUT', '修改当前用户信息', 74, 0);
INSERT INTO `resource` VALUES (101, '2022-12-07 20:55:08.271', '2022-12-07 20:55:08.271', '/user/disable', 'PUT', '修改用户禁用', 74, 0);
INSERT INTO `resource` VALUES (102, '2022-12-08 15:43:15.421', '2022-12-08 15:43:15.421', '', '', '页面模块', 0, 0);
INSERT INTO `resource` VALUES (103, '2022-12-08 15:43:26.009', '2022-12-08 15:43:26.009', '/page/list', 'GET', '页面列表', 102, 0);
INSERT INTO `resource` VALUES (104, '2022-12-08 15:43:38.570', '2022-12-08 15:43:38.570', '/page', 'POST', '新增/编辑页面', 102, 0);
INSERT INTO `resource` VALUES (105, '2022-12-08 15:43:50.879', '2022-12-08 15:43:50.879', '/page', 'DELETE', '删除页面', 102, 0);
INSERT INTO `resource` VALUES (106, '2022-12-16 11:53:57.989', '2022-12-16 11:53:57.989', '', '', '文件模块', 0, 0);
INSERT INTO `resource` VALUES (107, '2022-12-16 11:54:20.891', '2022-12-16 11:54:20.891', '/upload', 'POST', '文件上传', 106, 0);

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '角色名',
  `label` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '角色描述',
  `is_disable` tinyint(1) NULL DEFAULT NULL COMMENT '是否禁用(0-否 1-是)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES (1, '2022-10-20 21:24:28.000', '2022-12-16 11:54:27.862', '管理员', 'admin', 0);
INSERT INTO `role` VALUES (2, '2022-10-20 21:25:07.000', '2022-12-08 15:44:09.091', '用户', 'user', 0);
INSERT INTO `role` VALUES (3, '2022-10-20 21:25:09.000', '2022-12-08 15:44:13.917', '测试', 'test', 0);

-- ----------------------------
-- Table structure for role_menu
-- ----------------------------
DROP TABLE IF EXISTS `role_menu`;
CREATE TABLE `role_menu`  (
  `role_id` bigint NULL DEFAULT NULL,
  `menu_id` bigint NULL DEFAULT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_menu
-- ----------------------------
INSERT INTO `role_menu` VALUES (2, 1);
INSERT INTO `role_menu` VALUES (2, 2);
INSERT INTO `role_menu` VALUES (2, 6);
INSERT INTO `role_menu` VALUES (2, 8);
INSERT INTO `role_menu` VALUES (2, 9);
INSERT INTO `role_menu` VALUES (2, 10);
INSERT INTO `role_menu` VALUES (2, 3);
INSERT INTO `role_menu` VALUES (2, 25);
INSERT INTO `role_menu` VALUES (2, 26);
INSERT INTO `role_menu` VALUES (2, 16);
INSERT INTO `role_menu` VALUES (2, 17);
INSERT INTO `role_menu` VALUES (2, 23);
INSERT INTO `role_menu` VALUES (2, 24);
INSERT INTO `role_menu` VALUES (2, 4);
INSERT INTO `role_menu` VALUES (2, 27);
INSERT INTO `role_menu` VALUES (2, 28);
INSERT INTO `role_menu` VALUES (2, 29);
INSERT INTO `role_menu` VALUES (2, 5);
INSERT INTO `role_menu` VALUES (2, 30);
INSERT INTO `role_menu` VALUES (2, 32);
INSERT INTO `role_menu` VALUES (2, 31);
INSERT INTO `role_menu` VALUES (2, 33);
INSERT INTO `role_menu` VALUES (2, 34);
INSERT INTO `role_menu` VALUES (2, 36);
INSERT INTO `role_menu` VALUES (2, 37);
INSERT INTO `role_menu` VALUES (2, 38);
INSERT INTO `role_menu` VALUES (2, 39);
INSERT INTO `role_menu` VALUES (3, 1);
INSERT INTO `role_menu` VALUES (3, 2);
INSERT INTO `role_menu` VALUES (3, 3);
INSERT INTO `role_menu` VALUES (3, 4);
INSERT INTO `role_menu` VALUES (3, 33);
INSERT INTO `role_menu` VALUES (3, 6);
INSERT INTO `role_menu` VALUES (3, 34);
INSERT INTO `role_menu` VALUES (3, 8);
INSERT INTO `role_menu` VALUES (3, 9);
INSERT INTO `role_menu` VALUES (3, 10);
INSERT INTO `role_menu` VALUES (3, 25);
INSERT INTO `role_menu` VALUES (3, 26);
INSERT INTO `role_menu` VALUES (3, 16);
INSERT INTO `role_menu` VALUES (3, 17);
INSERT INTO `role_menu` VALUES (3, 23);
INSERT INTO `role_menu` VALUES (3, 24);
INSERT INTO `role_menu` VALUES (3, 5);
INSERT INTO `role_menu` VALUES (3, 29);
INSERT INTO `role_menu` VALUES (3, 30);
INSERT INTO `role_menu` VALUES (3, 32);
INSERT INTO `role_menu` VALUES (3, 31);
INSERT INTO `role_menu` VALUES (3, 27);
INSERT INTO `role_menu` VALUES (3, 28);
INSERT INTO `role_menu` VALUES (3, 36);
INSERT INTO `role_menu` VALUES (3, 37);
INSERT INTO `role_menu` VALUES (3, 38);
INSERT INTO `role_menu` VALUES (3, 39);
INSERT INTO `role_menu` VALUES (1, 2);
INSERT INTO `role_menu` VALUES (1, 3);
INSERT INTO `role_menu` VALUES (1, 4);
INSERT INTO `role_menu` VALUES (1, 5);
INSERT INTO `role_menu` VALUES (1, 6);
INSERT INTO `role_menu` VALUES (1, 7);
INSERT INTO `role_menu` VALUES (1, 8);
INSERT INTO `role_menu` VALUES (1, 9);
INSERT INTO `role_menu` VALUES (1, 16);
INSERT INTO `role_menu` VALUES (1, 17);
INSERT INTO `role_menu` VALUES (1, 23);
INSERT INTO `role_menu` VALUES (1, 24);
INSERT INTO `role_menu` VALUES (1, 25);
INSERT INTO `role_menu` VALUES (1, 26);
INSERT INTO `role_menu` VALUES (1, 27);
INSERT INTO `role_menu` VALUES (1, 28);
INSERT INTO `role_menu` VALUES (1, 29);
INSERT INTO `role_menu` VALUES (1, 30);
INSERT INTO `role_menu` VALUES (1, 31);
INSERT INTO `role_menu` VALUES (1, 32);
INSERT INTO `role_menu` VALUES (1, 33);
INSERT INTO `role_menu` VALUES (1, 36);
INSERT INTO `role_menu` VALUES (1, 37);
INSERT INTO `role_menu` VALUES (1, 38);
INSERT INTO `role_menu` VALUES (1, 34);
INSERT INTO `role_menu` VALUES (1, 10);
INSERT INTO `role_menu` VALUES (1, 39);

-- ----------------------------
-- Table structure for role_resource
-- ----------------------------
DROP TABLE IF EXISTS `role_resource`;
CREATE TABLE `role_resource`  (
  `role_id` bigint NULL DEFAULT NULL,
  `resource_id` bigint NULL DEFAULT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_resource
-- ----------------------------
INSERT INTO `role_resource` VALUES (2, 43);
INSERT INTO `role_resource` VALUES (2, 44);
INSERT INTO `role_resource` VALUES (2, 59);
INSERT INTO `role_resource` VALUES (2, 38);
INSERT INTO `role_resource` VALUES (2, 42);
INSERT INTO `role_resource` VALUES (2, 65);
INSERT INTO `role_resource` VALUES (2, 68);
INSERT INTO `role_resource` VALUES (2, 62);
INSERT INTO `role_resource` VALUES (2, 35);
INSERT INTO `role_resource` VALUES (2, 36);
INSERT INTO `role_resource` VALUES (2, 51);
INSERT INTO `role_resource` VALUES (2, 54);
INSERT INTO `role_resource` VALUES (2, 55);
INSERT INTO `role_resource` VALUES (2, 58);
INSERT INTO `role_resource` VALUES (2, 70);
INSERT INTO `role_resource` VALUES (2, 78);
INSERT INTO `role_resource` VALUES (2, 82);
INSERT INTO `role_resource` VALUES (2, 79);
INSERT INTO `role_resource` VALUES (2, 80);
INSERT INTO `role_resource` VALUES (2, 86);
INSERT INTO `role_resource` VALUES (2, 92);
INSERT INTO `role_resource` VALUES (2, 95);
INSERT INTO `role_resource` VALUES (2, 41);
INSERT INTO `role_resource` VALUES (2, 103);
INSERT INTO `role_resource` VALUES (3, 43);
INSERT INTO `role_resource` VALUES (3, 44);
INSERT INTO `role_resource` VALUES (3, 59);
INSERT INTO `role_resource` VALUES (3, 38);
INSERT INTO `role_resource` VALUES (3, 41);
INSERT INTO `role_resource` VALUES (3, 42);
INSERT INTO `role_resource` VALUES (3, 65);
INSERT INTO `role_resource` VALUES (3, 68);
INSERT INTO `role_resource` VALUES (3, 62);
INSERT INTO `role_resource` VALUES (3, 35);
INSERT INTO `role_resource` VALUES (3, 36);
INSERT INTO `role_resource` VALUES (3, 79);
INSERT INTO `role_resource` VALUES (3, 80);
INSERT INTO `role_resource` VALUES (3, 51);
INSERT INTO `role_resource` VALUES (3, 54);
INSERT INTO `role_resource` VALUES (3, 55);
INSERT INTO `role_resource` VALUES (3, 58);
INSERT INTO `role_resource` VALUES (3, 70);
INSERT INTO `role_resource` VALUES (3, 78);
INSERT INTO `role_resource` VALUES (3, 82);
INSERT INTO `role_resource` VALUES (3, 92);
INSERT INTO `role_resource` VALUES (3, 95);
INSERT INTO `role_resource` VALUES (3, 86);
INSERT INTO `role_resource` VALUES (3, 103);
INSERT INTO `role_resource` VALUES (1, 3);
INSERT INTO `role_resource` VALUES (1, 43);
INSERT INTO `role_resource` VALUES (1, 44);
INSERT INTO `role_resource` VALUES (1, 45);
INSERT INTO `role_resource` VALUES (1, 46);
INSERT INTO `role_resource` VALUES (1, 47);
INSERT INTO `role_resource` VALUES (1, 48);
INSERT INTO `role_resource` VALUES (1, 6);
INSERT INTO `role_resource` VALUES (1, 59);
INSERT INTO `role_resource` VALUES (1, 60);
INSERT INTO `role_resource` VALUES (1, 61);
INSERT INTO `role_resource` VALUES (1, 7);
INSERT INTO `role_resource` VALUES (1, 38);
INSERT INTO `role_resource` VALUES (1, 39);
INSERT INTO `role_resource` VALUES (1, 40);
INSERT INTO `role_resource` VALUES (1, 41);
INSERT INTO `role_resource` VALUES (1, 42);
INSERT INTO `role_resource` VALUES (1, 8);
INSERT INTO `role_resource` VALUES (1, 65);
INSERT INTO `role_resource` VALUES (1, 66);
INSERT INTO `role_resource` VALUES (1, 67);
INSERT INTO `role_resource` VALUES (1, 68);
INSERT INTO `role_resource` VALUES (1, 9);
INSERT INTO `role_resource` VALUES (1, 62);
INSERT INTO `role_resource` VALUES (1, 63);
INSERT INTO `role_resource` VALUES (1, 64);
INSERT INTO `role_resource` VALUES (1, 10);
INSERT INTO `role_resource` VALUES (1, 23);
INSERT INTO `role_resource` VALUES (1, 34);
INSERT INTO `role_resource` VALUES (1, 35);
INSERT INTO `role_resource` VALUES (1, 36);
INSERT INTO `role_resource` VALUES (1, 37);
INSERT INTO `role_resource` VALUES (1, 79);
INSERT INTO `role_resource` VALUES (1, 80);
INSERT INTO `role_resource` VALUES (1, 81);
INSERT INTO `role_resource` VALUES (1, 85);
INSERT INTO `role_resource` VALUES (1, 49);
INSERT INTO `role_resource` VALUES (1, 51);
INSERT INTO `role_resource` VALUES (1, 52);
INSERT INTO `role_resource` VALUES (1, 53);
INSERT INTO `role_resource` VALUES (1, 54);
INSERT INTO `role_resource` VALUES (1, 50);
INSERT INTO `role_resource` VALUES (1, 55);
INSERT INTO `role_resource` VALUES (1, 56);
INSERT INTO `role_resource` VALUES (1, 57);
INSERT INTO `role_resource` VALUES (1, 58);
INSERT INTO `role_resource` VALUES (1, 69);
INSERT INTO `role_resource` VALUES (1, 70);
INSERT INTO `role_resource` VALUES (1, 71);
INSERT INTO `role_resource` VALUES (1, 72);
INSERT INTO `role_resource` VALUES (1, 91);
INSERT INTO `role_resource` VALUES (1, 92);
INSERT INTO `role_resource` VALUES (1, 93);
INSERT INTO `role_resource` VALUES (1, 78);
INSERT INTO `role_resource` VALUES (1, 82);
INSERT INTO `role_resource` VALUES (1, 84);
INSERT INTO `role_resource` VALUES (1, 86);
INSERT INTO `role_resource` VALUES (1, 98);
INSERT INTO `role_resource` VALUES (1, 95);
INSERT INTO `role_resource` VALUES (1, 11);
INSERT INTO `role_resource` VALUES (1, 99);
INSERT INTO `role_resource` VALUES (1, 100);
INSERT INTO `role_resource` VALUES (1, 101);
INSERT INTO `role_resource` VALUES (1, 74);
INSERT INTO `role_resource` VALUES (1, 102);
INSERT INTO `role_resource` VALUES (1, 103);
INSERT INTO `role_resource` VALUES (1, 104);
INSERT INTO `role_resource` VALUES (1, 105);
INSERT INTO `role_resource` VALUES (1, 106);
INSERT INTO `role_resource` VALUES (1, 107);

-- ----------------------------
-- Table structure for tag
-- ----------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of tag
-- ----------------------------
INSERT INTO `tag` VALUES (1, '2022-12-03 22:01:51.624', '2022-12-03 22:01:51.624', 'Golang');
INSERT INTO `tag` VALUES (2, '2022-12-03 22:01:56.984', '2022-12-03 22:01:56.984', 'Vue');

-- ----------------------------
-- Table structure for user_auth
-- ----------------------------
DROP TABLE IF EXISTS `user_auth`;
CREATE TABLE `user_auth`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `user_info_id` bigint NULL DEFAULT NULL COMMENT '用户信息ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '用户名',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '密码',
  `login_type` tinyint(1) NULL DEFAULT NULL COMMENT '登录类型',
  `ip_address` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '登录IP地址',
  `ip_source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'IP来源',
  `last_login_time` datetime(3) NULL DEFAULT NULL COMMENT '上次登录时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_auth
-- ----------------------------
INSERT INTO `user_auth` VALUES (1, '2022-10-31 21:54:11.040', '2022-12-03 21:51:55.809', 1, 'admin', '$2a$10$np.P54Jep7GB/H5vG1PcbudYcxAAf1iiBf7NzsQJT9ZfsYz6tFPcm', 1, '127.0.0.1:58096', '未知', '2022-12-03 21:51:55.809');
INSERT INTO `user_auth` VALUES (2, '2022-10-19 22:31:26.805', '2022-12-03 12:20:36.539', 2, 'user', '$2a$10$9vHpoeT7sF4j9beiZfPsOe0jJ67gOceO2WKJzJtHRZCjNJajl7Fhq', 1, '172.17.0.1:40280', '未知', '2022-12-03 12:20:36.539');
INSERT INTO `user_auth` VALUES (3, '2022-11-01 10:41:13.300', '2022-12-03 22:35:53.842', 3, 'test@qq.com', '$2a$10$FmU4jxwDlibSL9pdt.AsuODkbB4gLp3IyyXeoMmW/XALtT/HdwTsi', 1, '[::1]:36896', '未知', '2022-12-03 22:35:53.842');

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `email` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '邮箱',
  `nickname` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '昵称',
  `avatar` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '头像地址',
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '个人简介',
  `website` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '个人网站',
  `is_disable` tinyint(1) NULL DEFAULT NULL COMMENT '是否禁用(0-否 1-是)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_info
-- ----------------------------
INSERT INTO `user_info` VALUES (1, '2022-10-31 21:54:10.935', '2022-12-16 23:42:32.109', '', '管理员', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我是管理员用户！', 'https://www.hahacode.cn', 0);
INSERT INTO `user_info` VALUES (2, '2022-10-19 22:31:26.734', '2022-12-07 20:55:30.347', 'user@qq.com', '普通用户', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我是个普通用户！', 'https://www.hahacode.cn', 0);
INSERT INTO `user_info` VALUES (3, '2022-11-01 10:41:13.234', '2022-11-30 13:51:51.799', 'test@qq.com', '测试用户', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我是测试用的！', 'https://www.hahacode.cn', 0);

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role`  (
  `user_id` bigint NULL DEFAULT NULL,
  `role_id` bigint NULL DEFAULT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_role
-- ----------------------------
INSERT INTO `user_role` VALUES (1, 1);
INSERT INTO `user_role` VALUES (2, 2);
INSERT INTO `user_role` VALUES (3, 3);

SET FOREIGN_KEY_CHECKS = 1;
