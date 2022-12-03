import { request } from '@/utils'

export function useLinkApi() {
  return {
    getLinks: (params = {}) => request.get('/link/list', { params }),
    deleteLinks: (data = []) => request.delete(`/link`, { data }),
    saveOrUpdateLink: (data) => request.post('/link', data),
  }
}
