import { request } from '@/utils'

export function useHomeAPI() {
  return {
    getHomeInfo: () => request.get('/home'),
  }
}
