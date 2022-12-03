import { request } from '@/utils'

export function useProfileApi() {
  return {
    getProfile: (id) => request.get(`/profile/${id}`),
    updateProfile: (data) => request.put(`/profile/${data.id}`, data),
  }
}
