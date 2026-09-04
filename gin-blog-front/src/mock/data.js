// Mock 数据源: 无后端时供 mock 适配器使用
// 图片用本仓库 images/ 目录下的图片(GitHub 直接提供), 不依赖后端静态资源服务,
// 也不再依赖第三方示例图站点
const IMG_BASE = 'https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images'

const COVER = `${IMG_BASE}/page/article_list.jpg`
const AVATAR = `${IMG_BASE}/common/header.jpeg`
const TOURIST_AVATAR = `${IMG_BASE}/config/tourist_avatar.jpeg`

export const blogConfig = {
  website_name: 'Gin Vue Blog (Mock)',
  website_author: '阵雨',
  website_intro: '当前为 Mock 模式, 数据均为本地假数据',
  website_notice: '这是 Mock 模式下的公告, 不需要启动后端即可浏览全部页面。',
  website_createtime: '2023-12-27 22:00:00',
  website_record: '粤ICP备2021032312号',
  website_avatar: AVATAR,
  article_cover: COVER,
  user_avatar: AVATAR,
  tourist_avatar: TOURIST_AVATAR,
  github: 'https://github.com/szluyu99',
  gitee: 'https://gitee.com/szluyu99',
  qq: '123456789',
  is_comment_review: 'false',
  is_message_review: 'false',
}

export const pages = [
  { id: 1, name: '首页', label: 'home', cover: `${IMG_BASE}/page/home.jpg` },
  { id: 2, name: '归档', label: 'archive', cover: `${IMG_BASE}/page/archive.png` },
  { id: 3, name: '分类', label: 'category', cover: `${IMG_BASE}/page/category.png` },
  { id: 4, name: '标签', label: 'tag', cover: `${IMG_BASE}/page/tag.png` },
  { id: 5, name: '友链', label: 'link', cover: `${IMG_BASE}/page/link.jpg` },
  { id: 6, name: '关于', label: 'about', cover: `${IMG_BASE}/page/about.jpg` },
  { id: 7, name: '留言', label: 'message', cover: `${IMG_BASE}/page/message.jpeg` },
  { id: 8, name: '个人中心', label: 'user', cover: `${IMG_BASE}/page/user.jpg` },
  { id: 9, name: '相册', label: 'album', cover: `${IMG_BASE}/page/album.png` },
  { id: 10, name: '错误页面', label: '404', cover: `${IMG_BASE}/page/404.jpg` },
  { id: 11, name: '文章列表', label: 'article_list', cover: `${IMG_BASE}/page/article_list.jpg` },
]

export const categories = [
  { id: 1, name: '后端', created_at: '2023-12-27T22:45:09.369Z', updated_at: '2023-12-27T22:45:09.369Z' },
  { id: 2, name: '前端', created_at: '2023-12-27T22:45:15.006Z', updated_at: '2023-12-27T22:45:15.006Z' },
  { id: 3, name: '项目', created_at: '2023-12-27T22:46:36.057Z', updated_at: '2023-12-27T22:46:36.057Z' },
  { id: 4, name: '学习', created_at: '2023-12-27T22:47:47.501Z', updated_at: '2023-12-27T22:47:47.501Z' },
]

export const tags = [
  { id: 1, name: 'Golang', created_at: '2023-12-27T22:45:40.731Z', updated_at: '2023-12-27T22:45:40.731Z' },
  { id: 2, name: 'Vue', created_at: '2023-12-27T22:46:36.082Z', updated_at: '2023-12-27T22:46:36.082Z' },
  { id: 3, name: '感悟', created_at: '2023-12-27T22:47:47.530Z', updated_at: '2023-12-27T22:47:47.530Z' },
  { id: 4, name: 'Gin', created_at: '2023-12-28T09:10:00.000Z', updated_at: '2023-12-28T09:10:00.000Z' },
  { id: 5, name: '工程化', created_at: '2023-12-28T09:12:00.000Z', updated_at: '2023-12-28T09:12:00.000Z' },
]

// 文章正文, 单独抽出来便于阅读
const CONTENT = {
  1: `## Mock 模式已启用\n\n当前页面的数据全部来自 \`src/mock\`, 不需要启动 Go 后端。\n\n\`\`\`go\nfmt.Println("hello mock")\n\`\`\`\n\n\`\`\`js\nconsole.log('hello mock')\n\`\`\`\n\n支持公式渲染:\n\n$$\n\\large X^{2m}_{3n}\n$$\n`,
  2: `## 学习有捷径\n\n学习的捷径之一就是多看看别人是怎么理解这些知识的。\n\n举两个例子。\n\n如果你喜欢《水浊》, 千万不要只把原著当故事看, 也去看看别人的解读。\n\n> 输入决定输出, 输入的质量决定输出的质量。\n`,
  3: `## 项目介绍\n\n这是一个前后端分离的博客系统, 由三部分组成:\n\n- \`gin-blog-server\`: Go + Gin + GORM 后端\n- \`gin-blog-front\`: Vue3 + Vite 前台\n- \`gin-blog-admin\`: Vue3 + Naive UI 后台\n\n### 特性\n\n1. 文章支持 Markdown 与代码高亮\n2. 评论、留言、友链、归档\n3. RBAC 权限管理\n`,
  4: `## Gin 中间件是怎么串起来的\n\nGin 的中间件本质上是一个 \`HandlerFunc\` 切片, 通过 \`c.Next()\` 驱动索引前进。\n\n\`\`\`go\nfunc (c *Context) Next() {\n\tc.index++\n\tfor c.index < int8(len(c.handlers)) {\n\t\tc.handlers[c.index](c)\n\t\tc.index++\n\t}\n}\n\`\`\`\n\n理解这一点, 就能理解为什么 \`c.Abort()\` 只是把索引推到末尾。\n`,
  5: `## Vue3 组合式 API 的取舍\n\n组合式 API 解决的核心问题是**逻辑复用**, 而不是少写几行代码。\n\n### 什么时候抽 composable\n\n- 同一段逻辑在两个以上组件里出现\n- 逻辑本身有独立的生命周期\n\n不满足这两条时, 直接写在组件里更好读。\n`,
  6: `## 前端工程化里那些容易被忽略的细节\n\n1. lockfile 一定要提交, 并且注意包管理器版本\n2. 依赖升级要看 CHANGELOG 的 breaking change\n3. 环境变量区分 \`.env\` / \`.env.development\` / \`.env.production\`\n\n> 一次不受控的依赖升级, 可以让整个应用白屏。\n`,
}

// 手写的几篇: 正文各不相同, 用来看 Markdown / 代码高亮 / 公式渲染
const HANDWRITTEN = [
  { id: 6, title: '前端工程化里那些容易被忽略的细节', created_at: '2024-01-05T10:20:00.000Z', category_id: 2, tag_ids: [5, 2], is_top: 0 },
  { id: 5, title: 'Vue3 组合式 API 的取舍', created_at: '2024-01-03T15:30:00.000Z', category_id: 2, tag_ids: [2, 5], is_top: 0 },
  { id: 4, title: 'Gin 中间件是怎么串起来的', created_at: '2024-01-01T09:00:00.000Z', category_id: 1, tag_ids: [1, 4], is_top: 0 },
  { id: 3, title: '项目介绍', created_at: '2023-12-27T22:48:43.727Z', category_id: 3, tag_ids: [1, 2], is_top: 1 },
  { id: 2, title: '学习有捷径', created_at: '2023-12-27T22:47:47.513Z', category_id: 4, tag_ids: [3], is_top: 0 },
  { id: 1, title: '项目运行成功', created_at: '2023-12-27T22:46:36.066Z', category_id: 3, tag_ids: [1, 2], is_top: 0 },
]

// 手写 6 篇翻不动页: 首页无限滚动一页 8 条, 文章列表页一页 9 条。
// 剩下的按模板批量生成, 只有标题和分类标签不同, 免得 data.js 变成一大坨正文
const FILLER_TITLES = [
  'GORM 的预加载与 N+1 查询',
  '为什么不要在 handler 里写 SQL',
  'JWT 与 Session 各自适合什么场景',
  'Redis 做计数器时的原子性问题',
  '接口幂等性的几种实现',
  'Go 的 context 到底传了什么',
  '错误处理: sentinel error 还是 error wrapping',
  'sync.Pool 什么时候真的有用',
  'Pinia 的 store 该拆到多细',
  'Vue Router 的动态路由怎么落地',
  '组件测试和端到端测试的分工',
  'UnoCSS 的 preset 与 shortcut',
  '为什么我的 CSS 类名没生成',
  '前端如何优雅地处理接口错误',
  'Docker 多阶段构建能省多少体积',
  'Nginx 反向代理里最常写错的三行',
  'CI 跑得慢, 时间都花在哪了',
  '写文档这件事的投入产出',
]

const FILLER = FILLER_TITLES.map((title, i) => ({
  id: 7 + i,
  title,
  // 从 2024-01-08 起每篇隔 3 天, 归档页能分出多个月份
  created_at: new Date(Date.UTC(2024, 0, 8 + i * 3, 10, 0, 0)).toISOString(),
  category_id: (i % 4) + 1,
  tag_ids: [(i % 5) + 1, ((i + 2) % 5) + 1],
  is_top: 0,
  content: `## ${title}\n\n这是 Mock 模式下批量生成的样例文章, 用来验证列表分页和首页滚动加载。\n\n- 结论先写在前面\n- 代码片段尽量能直接跑\n- 需要注意的坑单独列出来\n\n\`\`\`go\nfmt.Println("mock article ${7 + i}")\n\`\`\`\n\n> ${title}\n`,
}))

// 文章: 列表 / 详情 / 归档 / 搜索 共用
// 按时间倒序: 详情页的上一篇/下一篇直接按数组下标取, 顺序不能乱
export const articles = [...HANDWRITTEN, ...FILLER]
  .map(e => ({
    ...e,
    updated_at: e.created_at,
    desc: '',
    content: e.content ?? CONTENT[e.id],
    img: COVER,
    type: 1,
    status: 1,
    is_delete: false,
    original_url: '',
    like_count: (e.id * 7) % 23,
    view_count: e.id * 31,
  }))
  .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))

// 留言 (弹幕)
export const messages = [
  { id: 1, nickname: '路人甲', avatar: TOURIST_AVATAR, content: '博客做得不错!', created_at: '2024-01-02T10:00:00.000Z' },
  { id: 2, nickname: '路人乙', avatar: TOURIST_AVATAR, content: 'Mock 模式下也能看到留言弹幕', created_at: '2024-01-02T11:00:00.000Z' },
  { id: 3, nickname: '小明', avatar: TOURIST_AVATAR, content: '前端不依赖后端真方便', created_at: '2024-01-03T09:00:00.000Z' },
  { id: 4, nickname: '小红', avatar: TOURIST_AVATAR, content: '来学 Vue3 了', created_at: '2024-01-03T09:30:00.000Z' },
  { id: 5, nickname: '匿名', avatar: TOURIST_AVATAR, content: '给作者点个 Star', created_at: '2024-01-04T20:00:00.000Z' },
]

// 友链
export const links = [
  { id: 1, name: 'Gin', avatar: AVATAR, address: 'https://gin-gonic.com', intro: 'Go 的 HTTP Web 框架' },
  { id: 2, name: 'Vue', avatar: AVATAR, address: 'https://cn.vuejs.org', intro: '渐进式 JavaScript 框架' },
  { id: 3, name: 'Vite', avatar: AVATAR, address: 'https://cn.vitejs.dev', intro: '下一代前端工具链' },
  { id: 4, name: 'UnoCSS', avatar: AVATAR, address: 'https://unocss.dev', intro: '即时按需原子 CSS 引擎' },
]

export const about = '## 关于本站\n\n这是 Mock 模式下的关于页面内容, 可以在 `src/mock/data.js` 中修改。\n\n- 技术栈: Vue3 + Vite + UnoCSS\n- 后端: Gin + GORM\n'

// 当前登录用户 (mock 登录后返回)
export const currentUser = {
  id: 1,
  nickname: 'Mock 用户',
  avatar: AVATAR,
  website: 'https://github.com/szluyu99',
  intro: '这是 Mock 模式下的登录用户',
  email: 'mock@example.com',
  article_like_set: [],
  comment_like_set: [],
}

// 评论: type 1-文章 2-友链 3-留言, topic_id 为 0 表示非文章页面
export const comments = [
  { id: 1, type: 1, topic_id: 1, parent_id: 0, user_id: 2, reply_user_id: 0, content: '跑起来了, 感谢分享!', like_count: 3, created_at: '2023-12-28T10:00:00.000Z' },
  { id: 2, type: 1, topic_id: 1, parent_id: 1, user_id: 3, reply_user_id: 2, content: '我也跑通了', like_count: 1, created_at: '2023-12-28T10:30:00.000Z' },
  { id: 3, type: 1, topic_id: 3, parent_id: 0, user_id: 3, reply_user_id: 0, content: '项目结构很清晰', like_count: 0, created_at: '2023-12-29T14:00:00.000Z' },
  { id: 4, type: 2, topic_id: 0, parent_id: 0, user_id: 2, reply_user_id: 0, content: '申请友链: https://example.com', like_count: 2, created_at: '2024-01-02T08:00:00.000Z' },
  { id: 5, type: 3, topic_id: 0, parent_id: 0, user_id: 3, reply_user_id: 0, content: '留言板测试', like_count: 0, created_at: '2024-01-04T08:00:00.000Z' },
]

// 评论用户信息, key 为 user_id
export const commentUsers = {
  1: { nickname: 'Mock 用户', avatar: AVATAR, website: 'https://github.com/szluyu99' },
  2: { nickname: '路人甲', avatar: TOURIST_AVATAR, website: '' },
  3: { nickname: '小明', avatar: TOURIST_AVATAR, website: 'https://example.com' },
}
