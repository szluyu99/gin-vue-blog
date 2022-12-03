import { request } from '@/utils'

export function useCategoryApi() {
  return {
    getCategorys: (params = {}) => request.get('/category/list', { params }),
    saveOrUpdateCategory: (data) => request.post('/category', data),
    deleteCategory: (data = []) => request.delete(`/category`, { data }),
    getCategoryOption: () => request.get('/category/option'),
    // addCategory: (data) => request.post('/category', data),
    // updateCategory: (data) => request.post(`/category`, data),
  }
}
