// Mock 数据源: 无后端时供 mock 适配器使用
// menus / resources / roles / pages / config 与后端初始化的默认数据保持一致, 保证动态路由和权限页面可用

// 图片用本仓库 images/ 目录下的图片(GitHub 直接提供), 不依赖第三方示例图站点
const IMG_BASE = 'https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images'
export const SAMPLE_IMG = `${IMG_BASE}/common/header.jpeg`

// 菜单 (动态路由数据源, 字段与 /menu/list 一致)
export const menus = [
  { id: 1, parent_id: 0, name: '首页', path: '/home', component: '/home', icon: 'ic:sharp-home', order_num: 0, redirect: '/home', is_catalogue: true, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 2, parent_id: 0, name: '文章管理', path: '/article', component: 'Layout', icon: 'ic:twotone-article', order_num: 1, redirect: '/article/list', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 3, parent_id: 0, name: '权限管理', path: '/auth', component: 'Layout', icon: 'cib:adguard', order_num: 3, redirect: '/auth/menu', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 4, parent_id: 0, name: '消息管理', path: '/message', component: 'Layout', icon: 'ic:twotone-email', order_num: 2, redirect: '/message/comment', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 5, parent_id: 0, name: '用户管理', path: '/user', component: 'Layout', icon: 'ph:user-list-bold', order_num: 4, redirect: '/user/list', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 6, parent_id: 0, name: '日志管理', path: '/log', component: 'Layout', icon: 'material-symbols:receipt-long-outline-rounded', order_num: 6, redirect: '/log/operation', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 7, parent_id: 0, name: '系统管理', path: '/setting', component: 'Layout', icon: 'ion:md-settings', order_num: 5, redirect: '/setting/website', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 8, parent_id: 0, name: '个人中心', path: '/profile', component: '/profile', icon: 'mdi:account', order_num: 7, redirect: '/profile', is_catalogue: true, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 9, parent_id: 2, name: '发布文章', path: 'write', component: '/article/write', icon: 'icon-park-outline:write', order_num: 1, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 10, parent_id: 2, name: '文章列表', path: 'list', component: '/article/list', icon: 'material-symbols:format-list-bulleted', order_num: 2, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 11, parent_id: 2, name: '分类管理', path: 'category', component: '/article/category', icon: 'tabler:category', order_num: 3, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 12, parent_id: 2, name: '标签管理', path: 'tag', component: '/article/tag', icon: 'tabler:tag', order_num: 4, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 13, parent_id: 2, name: '修改文章', path: 'write/:id', component: '/article/write', icon: 'icon-park-outline:write', order_num: 1, redirect: '', is_catalogue: false, is_hidden: true, keep_alive: false, is_external: false, external_link: '' },
  { id: 14, parent_id: 3, name: '菜单管理', path: 'menu', component: '/auth/menu', icon: 'ic:twotone-menu-book', order_num: 1, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 15, parent_id: 3, name: '接口管理', path: 'resource', component: '/auth/resource', icon: 'mdi:api', order_num: 2, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 16, parent_id: 3, name: '角色管理', path: 'role', component: '/auth/role', icon: 'carbon:user-role', order_num: 3, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 17, parent_id: 4, name: '评论管理', path: 'comment', component: '/message/comment', icon: 'ic:twotone-comment', order_num: 1, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 18, parent_id: 4, name: '留言管理', path: 'leave-msg', component: '/message/leave-msg', icon: 'ic:twotone-message', order_num: 2, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 19, parent_id: 5, name: '用户列表', path: 'list', component: '/user/list', icon: 'mdi:account', order_num: 1, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 20, parent_id: 5, name: '在线用户', path: 'online', component: '/user/online', icon: 'ic:outline-online-prediction', order_num: 2, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 21, parent_id: 6, name: '操作日志', path: 'operation', component: '/log/operation', icon: 'mdi:book-open-page-variant-outline', order_num: 1, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 22, parent_id: 6, name: '登录日志', path: 'login', component: '/log/login', icon: 'material-symbols:login', order_num: 2, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 23, parent_id: 7, name: '网站管理', path: 'website', component: '/setting/website', icon: 'el:website', order_num: 1, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 24, parent_id: 7, name: '页面管理', path: 'page', component: '/setting/page', icon: 'iconoir:journal-page', order_num: 2, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 25, parent_id: 7, name: '友链管理', path: 'link', component: '/setting/link', icon: 'mdi:telegram', order_num: 3, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
  { id: 26, parent_id: 7, name: '关于我', path: 'about', component: '/setting/about', icon: 'cib:about-me', order_num: 4, redirect: '', is_catalogue: false, is_hidden: false, keep_alive: false, is_external: false, external_link: '' },
]

// 接口资源 (parent_id 为 0 的是模块)
export const resources = [
  { id: 1, parent_id: 0, name: '文章模块', url: '', request_method: '', is_anonymous: false },
  { id: 2, parent_id: 0, name: '分类模块', url: '', request_method: '', is_anonymous: false },
  { id: 3, parent_id: 0, name: '标签模块', url: '', request_method: '', is_anonymous: false },
  { id: 4, parent_id: 0, name: '页面模块', url: '', request_method: '', is_anonymous: false },
  { id: 5, parent_id: 0, name: '友链模块', url: '', request_method: '', is_anonymous: false },
  { id: 6, parent_id: 0, name: '菜单模块', url: '', request_method: '', is_anonymous: false },
  { id: 7, parent_id: 0, name: '角色模块', url: '', request_method: '', is_anonymous: false },
  { id: 8, parent_id: 0, name: '资源模块', url: '', request_method: '', is_anonymous: false },
  { id: 9, parent_id: 0, name: '评论模块', url: '', request_method: '', is_anonymous: false },
  { id: 10, parent_id: 0, name: '留言模块', url: '', request_method: '', is_anonymous: false },
  { id: 11, parent_id: 0, name: '文件模块', url: '', request_method: '', is_anonymous: false },
  { id: 12, parent_id: 0, name: '博客信息模块', url: '', request_method: '', is_anonymous: false },
  { id: 13, parent_id: 0, name: '用户信息模块', url: '', request_method: '', is_anonymous: false },
  { id: 14, parent_id: 0, name: '操作日志模块', url: '', request_method: '', is_anonymous: false },
  { id: 15, parent_id: 1, name: '文章列表', url: '/article/list', request_method: 'GET', is_anonymous: false },
  { id: 16, parent_id: 1, name: '文章详情', url: '/article/:id', request_method: 'GET', is_anonymous: false },
  { id: 17, parent_id: 1, name: '新增/编辑文章', url: '/article', request_method: 'POST', is_anonymous: false },
  { id: 18, parent_id: 1, name: '更新文章软删除', url: '/article/soft-delete', request_method: 'PUT', is_anonymous: false },
  { id: 19, parent_id: 1, name: '删除文章', url: '/article', request_method: 'DELETE', is_anonymous: false },
  { id: 20, parent_id: 1, name: '修改文章置顶', url: '/article/top', request_method: 'PUT', is_anonymous: false },
  { id: 21, parent_id: 1, name: '导出文章', url: '/article/export', request_method: 'POST', is_anonymous: false },
  { id: 22, parent_id: 1, name: '导入文章', url: '/article/import', request_method: 'POST', is_anonymous: false },
  { id: 23, parent_id: 2, name: '分类列表', url: '/category/list', request_method: 'GET', is_anonymous: false },
  { id: 24, parent_id: 2, name: '新增/编辑分类', url: '/category', request_method: 'POST', is_anonymous: false },
  { id: 25, parent_id: 2, name: '删除分类', url: '/category', request_method: 'DELETE', is_anonymous: false },
  { id: 26, parent_id: 2, name: '分类选项列表', url: '/category/option', request_method: 'GET', is_anonymous: false },
  { id: 27, parent_id: 3, name: '标签列表', url: '/tag/list', request_method: 'GET', is_anonymous: false },
  { id: 28, parent_id: 3, name: '新增/编辑标签', url: '/tag', request_method: 'POST', is_anonymous: false },
  { id: 29, parent_id: 3, name: '删除标签', url: '/tag', request_method: 'DELETE', is_anonymous: false },
  { id: 30, parent_id: 3, name: '标签选项列表', url: '/tag/option', request_method: 'GET', is_anonymous: false },
  { id: 31, parent_id: 4, name: '页面列表', url: '/page/list', request_method: 'GET', is_anonymous: false },
  { id: 32, parent_id: 4, name: '新增/编辑页面', url: '/page', request_method: 'POST', is_anonymous: false },
  { id: 33, parent_id: 4, name: '删除页面', url: '/page', request_method: 'DELETE', is_anonymous: false },
  { id: 34, parent_id: 5, name: '友链列表', url: '/link/list', request_method: 'GET', is_anonymous: false },
  { id: 35, parent_id: 5, name: '新增/编辑友链', url: '/link', request_method: 'POST', is_anonymous: false },
  { id: 36, parent_id: 5, name: '删除友链', url: '/link', request_method: 'DELETE', is_anonymous: false },
  { id: 37, parent_id: 6, name: '菜单列表', url: '/menu/list', request_method: 'GET', is_anonymous: false },
  { id: 38, parent_id: 6, name: '新增/编辑菜单', url: '/menu', request_method: 'POST', is_anonymous: false },
  { id: 39, parent_id: 6, name: '删除菜单', url: '/menu', request_method: 'DELETE', is_anonymous: false },
  { id: 40, parent_id: 6, name: '菜单选项列表(树形)', url: '/menu/option', request_method: 'GET', is_anonymous: false },
  { id: 41, parent_id: 6, name: '获取当前用户菜单', url: '/menu/user/list', request_method: 'GET', is_anonymous: false },
  { id: 42, parent_id: 7, name: '角色列表', url: '/role/list', request_method: 'GET', is_anonymous: false },
  { id: 43, parent_id: 7, name: '新增/编辑角色', url: '/role', request_method: 'POST', is_anonymous: false },
  { id: 44, parent_id: 7, name: '删除角色', url: '/role', request_method: 'DELETE', is_anonymous: false },
  { id: 45, parent_id: 7, name: '角色选项列表', url: '/role/option', request_method: 'GET', is_anonymous: false },
  { id: 46, parent_id: 8, name: '资源列表', url: '/resource/list', request_method: 'GET', is_anonymous: false },
  { id: 47, parent_id: 8, name: '新增/编辑资源', url: '/resource', request_method: 'POST', is_anonymous: false },
  { id: 48, parent_id: 8, name: '删除资源', url: '/resource', request_method: 'DELETE', is_anonymous: false },
  { id: 49, parent_id: 8, name: '资源选项列表(树形)', url: '/resource/option', request_method: 'GET', is_anonymous: false },
  { id: 50, parent_id: 8, name: '修改资源匿名访问', url: '/resource/anonymous', request_method: 'PUT', is_anonymous: false },
  { id: 51, parent_id: 9, name: '评论列表', url: '/comment/list', request_method: 'GET', is_anonymous: false },
  { id: 52, parent_id: 9, name: '删除评论', url: '/comment', request_method: 'DELETE', is_anonymous: false },
  { id: 53, parent_id: 9, name: '修改评论审核', url: '/comment/review', request_method: 'PUT', is_anonymous: false },
  { id: 54, parent_id: 10, name: '留言列表', url: '/message/list', request_method: 'GET', is_anonymous: false },
  { id: 55, parent_id: 10, name: '删除留言', url: '/message', request_method: 'DELETE', is_anonymous: false },
  { id: 56, parent_id: 10, name: '修改留言审核', url: '/message/review', request_method: 'PUT', is_anonymous: false },
  { id: 57, parent_id: 11, name: '文件上传', url: '/upload', request_method: 'POST', is_anonymous: false },
  { id: 58, parent_id: 12, name: '获取博客设置', url: '/setting/blog-config', request_method: 'GET', is_anonymous: false },
  { id: 59, parent_id: 12, name: '获取关于我', url: '/setting/about', request_method: 'GET', is_anonymous: false },
  { id: 60, parent_id: 12, name: '修改博客设置', url: '/setting/blog-config', request_method: 'PUT', is_anonymous: false },
  { id: 61, parent_id: 12, name: '修改关于我', url: '/setting/about', request_method: 'PUT', is_anonymous: false },
  { id: 62, parent_id: 12, name: '获取后台首页信息', url: '/home', request_method: 'GET', is_anonymous: false },
  { id: 63, parent_id: 13, name: '用户列表', url: '/user/list', request_method: 'GET', is_anonymous: false },
  { id: 64, parent_id: 13, name: '获取当前用户信息', url: '/user/info', request_method: 'GET', is_anonymous: false },
  { id: 65, parent_id: 13, name: '修改用户信息', url: '/user', request_method: 'PUT', is_anonymous: false },
  { id: 66, parent_id: 13, name: '获取在线用户列表', url: '/user/online', request_method: 'GET', is_anonymous: false },
  { id: 67, parent_id: 13, name: '强制离线用户', url: '/user/offline', request_method: 'DELETE', is_anonymous: false },
  { id: 68, parent_id: 13, name: '修改当前用户密码', url: '/user/current/password', request_method: 'PUT', is_anonymous: false },
  { id: 69, parent_id: 13, name: '修改当前用户信息', url: '/user/current', request_method: 'PUT', is_anonymous: false },
  { id: 70, parent_id: 13, name: '修改用户禁用', url: '/user/disable', request_method: 'PUT', is_anonymous: false },
  { id: 71, parent_id: 14, name: '日志列表', url: '/operation/log/list', request_method: 'GET', is_anonymous: false },
  { id: 72, parent_id: 14, name: '删除操作日志', url: '/operation/log', request_method: 'DELETE', is_anonymous: false },
]

// 角色
export const roles = [
  { id: 1, name: 'admin', label: '管理员', is_disable: false, created_at: '2024-01-01T10:00:00.000Z', resource_ids: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72], menu_ids: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26] },
  { id: 2, name: 'guest', label: '游客', is_disable: false, created_at: '2024-01-01T10:00:00.000Z', resource_ids: [15, 16, 23, 26, 27, 30, 31, 34, 37, 40, 41, 42, 45, 46, 49, 51, 54, 58, 59, 62, 63, 64, 66, 71], menu_ids: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26] },
]

// 页面封面
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

// 网站配置
export const config = {
  website_avatar: `${IMG_BASE}/common/header.jpeg`,
  website_name: '阵雨的个人博客 (Mock)',
  website_author: '阵雨',
  website_intro: '当前为 Mock 模式, 数据均为本地假数据',
  website_notice: '欢迎来到阵雨的个人博客，项目还在开发中...',
  website_createtime: '2026-08-31 11:55:23',
  website_record: '粤ICP备2021032312号',
  qq: '123456789',
  github: 'https://github.com/szluyu99',
  gitee: 'https://gitee.com/szluyu99',
  tourist_avatar: `${IMG_BASE}/config/tourist_avatar.jpeg`,
  user_avatar: `${IMG_BASE}/config/user_avatar.jpeg`,
  article_cover: `${IMG_BASE}/config/default_article_cover.png`,
  is_comment_review: 'true',
  is_message_review: 'true',
}

// 分类
export const categories = [
  { id: 1, name: '后端', created_at: '2023-12-27T22:45:09.369Z', updated_at: '2023-12-27T22:45:09.369Z' },
  { id: 2, name: '前端', created_at: '2023-12-27T22:45:15.006Z', updated_at: '2023-12-27T22:45:15.006Z' },
  { id: 3, name: '项目', created_at: '2023-12-27T22:46:36.057Z', updated_at: '2023-12-27T22:46:36.057Z' },
  { id: 4, name: '学习', created_at: '2023-12-27T22:47:47.501Z', updated_at: '2023-12-27T22:47:47.501Z' },
]

// 标签
export const tags = [
  { id: 1, name: 'Golang', created_at: '2023-12-27T22:45:40.731Z', updated_at: '2023-12-27T22:45:40.731Z' },
  { id: 2, name: 'Vue', created_at: '2023-12-27T22:46:36.082Z', updated_at: '2023-12-27T22:46:36.082Z' },
  { id: 3, name: '感悟', created_at: '2023-12-27T22:47:47.530Z', updated_at: '2023-12-27T22:47:47.530Z' },
  { id: 4, name: 'Gin', created_at: '2023-12-28T09:10:00.000Z', updated_at: '2023-12-28T09:10:00.000Z' },
  { id: 5, name: '工程化', created_at: '2023-12-28T09:12:00.000Z', updated_at: '2023-12-28T09:12:00.000Z' },
]

// 文章 (type: 1-原创 2-转载 3-翻译, status: 1-公开 2-私密 3-草稿)
export const articles = [
  { id: 6, title: '前端工程化里那些容易被忽略的细节', created_at: '2024-01-05T10:20:00.000Z', category_id: 2, tag_ids: [5, 2], is_top: 0, status: 1, is_delete: false },
  { id: 5, title: 'Vue3 组合式 API 的取舍', created_at: '2024-01-03T15:30:00.000Z', category_id: 2, tag_ids: [2, 5], is_top: 0, status: 1, is_delete: false },
  { id: 4, title: 'Gin 中间件是怎么串起来的', created_at: '2024-01-01T09:00:00.000Z', category_id: 1, tag_ids: [1, 4], is_top: 0, status: 1, is_delete: false },
  { id: 3, title: '项目介绍', created_at: '2023-12-27T22:48:43.727Z', category_id: 3, tag_ids: [1, 2], is_top: 1, status: 1, is_delete: false },
  { id: 2, title: '学习有捷径', created_at: '2023-12-27T22:47:47.513Z', category_id: 4, tag_ids: [3], is_top: 0, status: 2, is_delete: false },
  { id: 1, title: '项目运行成功', created_at: '2023-12-27T22:46:36.066Z', category_id: 3, tag_ids: [1, 2], is_top: 0, status: 1, is_delete: true },
].map(e => ({
  ...e,
  updated_at: e.created_at,
  desc: '这是 Mock 模式下的文章摘要',
  content: `## ${e.title}\n\n这是 Mock 模式下的文章正文, 可以在 \`src/mock/data.js\` 中修改。\n`,
  img: SAMPLE_IMG,
  type: 1,
  original_url: '',
  like_count: (e.id * 7) % 23,
  view_count: e.id * 31,
  comment_count: e.id % 3,
}))

// 用户 (user_auth + user_info + roles)
export const users = [
  {
    id: 1,
    username: 'admin',
    login_type: 0,
    ip_address: '127.0.0.1',
    ip_source: '内网IP',
    last_login_time: '2024-01-05T10:00:00.000Z',
    created_at: '2023-12-27T22:00:00.000Z',
    updated_at: '2024-01-05T10:00:00.000Z',
    is_disable: false,
    is_super: true,
    user_info_id: 1,
    info: { id: 1, email: 'admin@example.com', nickname: 'admin', avatar: SAMPLE_IMG, intro: 'Mock 模式下的管理员', website: 'https://github.com/szluyu99' },
    role_ids: [1],
  },
  {
    id: 2,
    username: 'guest',
    login_type: 0,
    ip_address: '10.0.0.2',
    ip_source: '内网IP',
    last_login_time: '2024-01-04T18:00:00.000Z',
    created_at: '2023-12-27T22:00:00.000Z',
    updated_at: '2024-01-04T18:00:00.000Z',
    is_disable: false,
    is_super: false,
    user_info_id: 2,
    info: { id: 2, email: 'guest@example.com', nickname: 'guest', avatar: SAMPLE_IMG, intro: 'Mock 模式下的游客', website: '' },
    role_ids: [2],
  },
]

// 评论 (type: 1-文章 2-友链 3-留言)
export const comments = [
  { id: 1, type: 1, topic_id: 3, parent_id: 0, user_id: 2, reply_user_id: 0, content: '跑起来了, 感谢分享!', like_count: 3, is_review: true, created_at: '2023-12-28T10:00:00.000Z' },
  { id: 2, type: 1, topic_id: 3, parent_id: 1, user_id: 1, reply_user_id: 2, content: '不客气', like_count: 1, is_review: true, created_at: '2023-12-28T10:30:00.000Z' },
  { id: 3, type: 1, topic_id: 4, parent_id: 0, user_id: 2, reply_user_id: 0, content: '中间件这段讲得清楚', like_count: 0, is_review: false, created_at: '2024-01-02T14:00:00.000Z' },
  { id: 4, type: 2, topic_id: 0, parent_id: 0, user_id: 2, reply_user_id: 0, content: '申请友链: https://example.com', like_count: 2, is_review: true, created_at: '2024-01-02T08:00:00.000Z' },
  { id: 5, type: 3, topic_id: 0, parent_id: 0, user_id: 2, reply_user_id: 0, content: '留言板测试', like_count: 0, is_review: false, created_at: '2024-01-04T08:00:00.000Z' },
]

// 留言 (弹幕)
export const messages = [
  { id: 1, nickname: '路人甲', avatar: SAMPLE_IMG, content: '博客做得不错!', ip_address: '10.0.0.3', ip_source: '内网IP', speed: 1, is_review: true, created_at: '2024-01-02T10:00:00.000Z' },
  { id: 2, nickname: '路人乙', avatar: SAMPLE_IMG, content: 'Mock 模式下也能看到留言', ip_address: '10.0.0.4', ip_source: '内网IP', speed: 1, is_review: true, created_at: '2024-01-02T11:00:00.000Z' },
  { id: 3, nickname: '小明', avatar: SAMPLE_IMG, content: '前端不依赖后端真方便', ip_address: '10.0.0.5', ip_source: '内网IP', speed: 2, is_review: false, created_at: '2024-01-03T09:00:00.000Z' },
  { id: 4, nickname: '小红', avatar: SAMPLE_IMG, content: '来学 Vue3 了', ip_address: '10.0.0.6', ip_source: '内网IP', speed: 1, is_review: false, created_at: '2024-01-03T09:30:00.000Z' },
  { id: 5, nickname: '匿名', avatar: SAMPLE_IMG, content: '给作者点个 Star', ip_address: '10.0.0.7', ip_source: '内网IP', speed: 3, is_review: true, created_at: '2024-01-04T20:00:00.000Z' },
]

// 友链
export const links = [
  { id: 1, name: 'Gin', avatar: SAMPLE_IMG, address: 'https://gin-gonic.com', intro: 'Go 的 HTTP Web 框架', created_at: '2024-01-01T10:00:00.000Z', updated_at: '2024-01-01T10:00:00.000Z' },
  { id: 2, name: 'Vue', avatar: SAMPLE_IMG, address: 'https://cn.vuejs.org', intro: '渐进式 JavaScript 框架', created_at: '2024-01-01T10:05:00.000Z', updated_at: '2024-01-01T10:05:00.000Z' },
  { id: 3, name: 'Vite', avatar: SAMPLE_IMG, address: 'https://cn.vitejs.dev', intro: '下一代前端工具链', created_at: '2024-01-01T10:10:00.000Z', updated_at: '2024-01-01T10:10:00.000Z' },
  { id: 4, name: 'UnoCSS', avatar: SAMPLE_IMG, address: 'https://unocss.dev', intro: '即时按需原子 CSS 引擎', created_at: '2024-01-01T10:15:00.000Z', updated_at: '2024-01-01T10:15:00.000Z' },
]

// 操作日志
export const operationLogs = [
  { id: 1, opt_module: '文章模块', opt_type: '新增或修改', opt_url: '/api/article', opt_method: 'handle.(*Article).SaveOrUpdate', opt_desc: '新增或修改文章', request_method: 'POST', request_param: '{"title":"Gin 中间件是怎么串起来的"}', response_data: '{"code":0}', user_id: 1, nickname: 'admin', ip_address: '127.0.0.1', ip_source: '内网IP', created_at: '2024-01-05T10:20:00.000Z' },
  { id: 2, opt_module: '分类模块', opt_type: '新增或修改', opt_url: '/api/category', opt_method: 'handle.(*Category).SaveOrUpdate', opt_desc: '新增或修改分类', request_method: 'POST', request_param: '{"name":"前端"}', response_data: '{"code":0}', user_id: 1, nickname: 'admin', ip_address: '127.0.0.1', ip_source: '内网IP', created_at: '2024-01-04T15:00:00.000Z' },
  { id: 3, opt_module: '标签模块', opt_type: '删除', opt_url: '/api/tag', opt_method: 'handle.(*Tag).Delete', opt_desc: '删除标签', request_method: 'DELETE', request_param: '[9]', response_data: '{"code":0}', user_id: 1, nickname: 'admin', ip_address: '127.0.0.1', ip_source: '内网IP', created_at: '2024-01-03T11:00:00.000Z' },
  { id: 4, opt_module: '友链模块', opt_type: '新增或修改', opt_url: '/api/link', opt_method: 'handle.(*Link).SaveOrUpdate', opt_desc: '新增或修改友链', request_method: 'POST', request_param: '{"name":"Vite"}', response_data: '{"code":0}', user_id: 1, nickname: 'admin', ip_address: '127.0.0.1', ip_source: '内网IP', created_at: '2024-01-02T09:30:00.000Z' },
  { id: 5, opt_module: '用户模块', opt_type: '修改', opt_url: '/api/user/disable', opt_method: 'handle.(*User).UpdateDisable', opt_desc: '修改用户禁用状态', request_method: 'PUT', request_param: '{"id":2,"is_disable":false}', response_data: '{"code":0}', user_id: 1, nickname: 'admin', ip_address: '127.0.0.1', ip_source: '内网IP', created_at: '2024-01-01T20:00:00.000Z' },
]

// 关于我 (Markdown)
export const about = '## 关于本站\n\n这是 Mock 模式下的关于页面内容, 可以在后台编辑后立即看到效果 (刷新页面后重置)。\n'
