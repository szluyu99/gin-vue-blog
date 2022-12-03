import { defineStore } from 'pinia'
import api from '@/api'

// 全局博客信息
interface BlogInfo {
  article_count?: number
  message_count?: number
  category_count?: number
  tag_count?: number
  view_count?: number

  blog_config?: {
    website_name?: string
    website_avatar?: string
    website_author?: string
    website_intro?: string
    website_notice?: string
    website_createtime?: string
    website_record?: string

    qq?: string
    github?: string
    gitee?: string
  }
}

export const useAppStore = defineStore('app', {
  state() {
    return {
      reloadFlag: true,
      searchFlag: false, // 搜索
      loginFlag: false, // 登录
      registerFlag: false, // 注册
      forgetFlag: false, // 忘记密码
      blogInfo: <BlogInfo>{},
    }
  },
  getters: {
    articleCount(): number {
      return this.blogInfo.article_count ?? 0
    },
  },
  actions: {
    setLoginFlag(flag: boolean) {
      this.loginFlag = flag
    },
    setRegisterFlag(flag: boolean) {
      this.registerFlag = flag
    },
    setForgetFlag(flag: boolean) {
      this.forgetFlag = flag
    },
    setSearchFlag(flag: boolean) {
      this.searchFlag = flag
    },
    async reloadPage() {
      window.$loadingBar?.start()
      this.reloadFlag = false
      await nextTick()
      this.reloadFlag = true

      setTimeout(() => {
        document.documentElement.scrollTo({ left: 0, top: 0 })
        window.$loadingBar?.finish()
      }, 100)
    },
    async getBlogInfo() {
      try {
        const res: any = await api.getHomeData()
        if (res.code === 0)
          this.blogInfo = res.data
        else return Promise.reject(res)
      }
      catch (err) {
        return Promise.reject(err)
      }
    },
  },
})
