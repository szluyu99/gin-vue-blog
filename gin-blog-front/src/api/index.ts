import { request } from '@/utils'
import type { RequestConfig } from '~/types/axios'

export default {
  login: (data: any = {}) => request.post('/login', data),
  register: (data: any = {}) => request.post('/register', data),
  logout: () => request.get('/logout'),
  about: () => request.get('/about'),

  // 发送验证码
  sendCode: (params?: { email: string }) => request.get('/code', { params }),

  // 首页数据
  getHomeData: () => request.get('/home'),
  // 首页文章列表
  getArticles: (params?: any) => request.get('/article/list', { params }),
  // 文章详情
  getArticleDetail: (id: number) => request.get(`/article/${id}`),
  // 文章归档
  getArchives: (params?: any) => request.get('/article/archive', { params }),
  // 文章搜索
  searchArticles: (params?: any) => request.get('/article/search', { params }),

  // 菜单列表
  getCategorys: () => request.get('/category/list'),
  // 标签列表
  getTags: () => request.get('/tag/list'),
  // 留言列表
  getMessages: () => request.get('/message/list'),
  // 友链列表
  getLinks: () => request.get('/link/list'),
  // 评论列表
  getComments: (params?: any) => request.get('/comment/list', { params }),
  // 评论回复列表
  getCommentReplies: (id: number, params?: any) => request.get(`/comment/replies/${id}`, { params }),

  // ! 需要 Token 的接口
  // 根据 token 获取当前用户信息
  getUser: () => request.get('/user/info', { needToken: true } as RequestConfig),
  // 修改当前用户信息
  updateUser: (data: any) => request.put('/user/info', data, { needToken: true } as RequestConfig),
  // 留言
  saveMessage: (data: any) => request.post('/message', data, { needToken: true } as RequestConfig),
  // 评论
  saveComment: (data: any) => request.post('/comment', data, { needToken: true } as RequestConfig),
  // 点赞评论
  saveLikeComment: (id: number) => request.get(`/comment/like/${id}`, { needToken: true } as RequestConfig),
  // 点赞文章
  saveLikeArticle: (id: number) => request.get(`/article/like/${id}`, { needToken: true } as RequestConfig),
}
