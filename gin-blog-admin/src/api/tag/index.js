import { request } from '@/utils'

export function useTagApi() {
  return {
    getTags: (params = {}) => request.get('/tag/list', { params }),
    saveOrUpdateTag: (data) => request.post('/tag', data),
    deleteTag: (data = []) => request.delete(`/tag`, { data }),
    getTagOption: () => request.get('/tag/option'),
  }
}
