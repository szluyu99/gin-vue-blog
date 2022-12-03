import { request } from '@/utils'

export function useArticleApi() {
  return {
    getArticles: (params = {}) => request.get('/article/list', { params }),
    getArticleById: (id) => request.get(`/article/${id}`),

    saveOrUpdateArticle: (data) => request.post('/article', data),
    // addArticle: (data) => request.post('/article/add', data),
    // updateArticle: (data) => request.put(`/article/${data.id}`, data),

    deleteArticle: (data = []) => request.delete(`/article`, { data }), // 物理删除
    softDeleteArticle: (data) => request.put(`/article/soft-delete`, data), // 软删除

    updateArticleTop: (data) => request.put(`/article/top`, data),
  }
}
