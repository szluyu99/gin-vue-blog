import { defineStore } from 'pinia'
import { useDark } from '@vueuse/core'
import api from '@/api'

// 应用全局变量
export const useAppStore = defineStore('app', {
  // persist: true,
  state: () => ({
    isDark: useDark(),
    blogConfig: {
      website_record: '暂无备案',
    },
  }),
  actions: {
    // 加载博客设置信息
    async getBlogInfo() {
      try {
        const resp = await api.getBlogConfig()
        if (resp.code === 0) {
          this.blogConfig = resp.data
        }
        else return Promise.reject(resp)
      }
      catch (err) {
        return Promise.reject(err)
      }
    },
  },
})
