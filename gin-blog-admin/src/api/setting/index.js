import { request } from '@/utils'

export function useSettingAPI() {
  return {
    getBlogConfig: () => request.get('/setting/blog-config'),
    updateBlogConfig: (data) => request.put('/setting/blog-config', data),
    getAbout: () => request.get('/setting/about'),
    updateAbout: (data) => request.put('/setting/about', data),
  }
}
