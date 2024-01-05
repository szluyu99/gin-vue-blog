import { request } from '@/utils'

export default {
  // refreshToken: () => request.post('/auth/refreshToken', null, { noNeedTip: true }),
  report: () => request.post('/report'), // 上报用户信息
  getHomeInfo: () => request.get('/home'), // 获取首页信息
  login: ({ username, password }) => request.post('/login', { username, password }, { noNeedToken: true }),
  logout: () => request.get('/logout'),

  // 文章相关接口
  getArticles: (params = {}) => request.get('/article/list', { params }),
  getArticleById: id => request.get(`/article/${id}`),
  saveOrUpdateArticle: data => request.post('/article', data),
  deleteArticle: (data = []) => request.delete('/article', { data }), // 物理删除
  softDeleteArticle: (ids, is_delete) => request.put('/article/soft-delete', { ids, is_delete }), // 软删除
  updateArticleTop: (id, is_top) => request.put('/article/top', { id, is_top }), // 修改文章置顶
  exportArticles: (data = []) => request.post('/article/export', data), // 导出文章
  importArticles: data => request.post('/article/import', data), // 导入文章

  // 分类相关接口
  getCategorys: (params = {}) => request.get('/category/list', { params }),
  saveOrUpdateCategory: data => request.post('/category', data),
  deleteCategory: (data = []) => request.delete('/category', { data }),
  getCategoryOption: () => request.get('/category/option'),

  // 标签相关接口
  getTags: (params = {}) => request.get('/tag/list', { params }),
  saveOrUpdateTag: data => request.post('/tag', data),
  deleteTag: (data = []) => request.delete('/tag', { data }),
  getTagOption: () => request.get('/tag/option'),

  // 留言相关接口
  getMessages: (params = {}) => request.get('/message/list', { params }),
  deleteMessages: (data = []) => request.delete('/message', { data }),
  updateMessageReview: (ids, is_review) => request.put('/message/review', { ids, is_review }),

  // 评论相关接口
  getComments: (params = {}) => request.get('/comment/list', { params }),
  deleteComments: (data = []) => request.delete('/comment', { data }),
  updateCommentReview: (ids, is_review) => request.put('/comment/review', { ids, is_review }),

  // 友链相关接口
  getLinks: (params = {}) => request.get('/link/list', { params }),
  deleteLinks: (data = []) => request.delete('/link', { data }),
  saveOrUpdateLink: data => request.post('/link', data),
  // 日志相关接口
  getOperationLogs: (params = {}) => request.get('/operation/log/list', { params }),
  deleteOperationLogs: (data = []) => request.delete('/operation/log', { data }),

  // 用户相关接口
  getUserInfo: () => request.get('/user/info'),
  updateCurrent: data => request.put('/user/current', data), // 更新当前用户信息
  updateCurrentPassword: data => request.put('/user/current/password', data), // 修改当前用户密码
  getUsers: (params = {}) => request.get('/user/list', { params }),
  updateUser: data => request.put('/user', data),
  updateUserDisable: (id, is_disable) => request.put('/user/disable', {
    id,
    is_disable,
  }),
  getOnlineUsers: (params = { keyword: '' }) => request.get('/user/online', { params }), // 在线用户列表
  forceOfflineUser: id => request.post(`/user/offline/${id}`), // 强制离线

  // 博客设置相关接口
  getConfig: () => request.get('/config'),
  updateConfig: data => request.patch('/config', data),
  // getBlogConfig: () => request.get('/setting/blog-config'),
  // updateBlogConfig: data => request.put('/setting/blog-config', data),
  getAbout: () => request.get('/setting/about'),
  updateAbout: data => request.put('/setting/about', data),

  // 权限管理相关接口
  // 菜单
  getUserMenus: () => request.get('/menu/user/list'), // 获取当前用户的菜单
  getMenus: (params = {}) => request.get('/menu/list', { params }),
  saveOrUpdateMenu: data => request.post('/menu', data),
  deleteMenu: id => request.delete(`/menu/${id}`),
  getMenuOption: () => request.get('/menu/option'),
  // 资源
  getResources: (params = {}) => request.get('/resource/list', { params }),
  saveOrUpdateResource: data => request.post('/resource', data),
  deleteResource: id => request.delete(`/resource/${id}`),
  updateResourceAnonymous: data => request.put('/resource/anonymous', data),
  getResourceOption: () => request.get('/resource/option'),
  // 角色
  getRoles: (params = {}) => request.get('/role/list', { params }),
  saveOrUpdateRole: data => request.post('/role', data),
  deleteRole: (data = []) => request.delete('/role', { data }),
  getRoleOption: () => request.get('/role/option'),

  // 页面相关接口
  getPages: () => request.get('/page/list'),
  saveOrUpdatePage: data => request.post('/page', data),
  deletePage: (data = []) => request.delete('/page', { data }),
}
