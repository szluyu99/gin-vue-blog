import { request } from '@/utils'

export function useUserApi() {
  return {
    getUser: () => request.get('/user/info'),
    getUsers: (params = {}) => request.get('/user/list', { params }),
    updateUser: (data) => request.put('/user', data),
    updateUserDisable: (data) => request.put('/user/disable', data),

    getOnlineUsers: () => request.get('/user/online'),
    forceOfflineUser: (data) => request.delete('/user/offline', { data }),
  }
}
