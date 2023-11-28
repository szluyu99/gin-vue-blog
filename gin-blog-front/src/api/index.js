import { request } from '@/utils'

export default {
  login: (data = {}) => request.post('/login', data),
  register: (data = {}) => request.post('/register', data),
  logout: () => request.get('/logout'),
  about: () => request.get('/about'),

  /**
   * 发送验证码
   * @param {{email: String}} params
   * @returns
   */
  sendCode: params => request.get('/code', { params }),

  // 首页数据
  getHomeData: () => request.get('/home'),
  // 首页文章列表
  getArticles: params => request.get('/article/list', { params }),
  /**
   * 文章详情
   * @param {Number} id
   */
  getArticleDetail: id => request.get(`/article/${id}`),
  // 文章归档
  getArchives: (params = {}) => request.get('/article/archive', { params }),
  // 文章搜索
  searchArticles: (params = {}) => request.get('/article/search', { params }),

  // 菜单列表
  getCategorys: () => request.get('/category/list'),
  // 标签列表
  getTags: () => request.get('/tag/list'),
  // 留言列表
  getMessages: () => request.get('/message/list'),
  // 友链列表
  getLinks: () => request.get('/link/list'),
  // 评论列表
  getComments: (params = {}) => request.get('/comment/list', { params }),
  /**
   * 评论回复列表
   * @param {Number} id
   * @param {any} params
   * @returns
   */
  getCommentReplies: (id, params = {}) => request.get(`/comment/replies/${id}`, { params }),

  // ! 需要 Token 的接口
  // 根据 token 获取当前用户信息
  getUser: () => request.get('/user/info', { needToken: true }),
  // 修改当前用户信息
  updateUser: data => request.put('/user/info', data, { needToken: true }),
  // 留言
  saveMessage: data => request.post('/message', data, { needToken: true }),
  // 评论
  saveComment: data => request.post('/comment', data, { needToken: true }),
  // 点赞评论
  saveLikeComment: id => request.get(`/comment/like/${id}`, { needToken: true }),
  // 点赞文章
  saveLikeArticle: id => request.get(`/article/like/${id}`, { needToken: true }),
}
