import { nextTick } from 'vue'
import { defineStore } from 'pinia'
import api from '@/api'
import { convertImgUrl } from '@/utils'

export const useAppStore = defineStore('app', {
  persist: true, // pinia 持久化插件
  state() {
    return {
      reloadFlag: true,
      // 搜索
      searchFlag: false,
      // 登录
      loginFlag: false,
      // 注册
      registerFlag: false,
      // 忘记密码
      forgetFlag: false,

      // 侧边栏折叠（移动端）
      collapsed: false,

      blogInfo: {
        blog_config: {
          website_name: '阵、雨的个人博客',
          website_author: '阵、雨',
          website_intro: '往事随风而去',
        },
      },
    }
  },
  getters: {
    isMobile: () => !!navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i),
    articleCount: state => state.blogInfo.article_count ?? 0,
    categoryCount: state => state.blogInfo.category_count ?? 0,
    tagCount: state => state.blogInfo.tag_count ?? 0,
    viewCount: state => state.blogInfo.view_count ?? 0,
    pageList: state => state.blogInfo.page_list ?? [],
    blogConfig: state => state.blogInfo.blog_config,
  },
  actions: {
    setCollapsed(flag) { this.collapsed = flag },
    setLoginFlag(flag) { this.loginFlag = flag },
    setRegisterFlag(flag) { this.registerFlag = flag },
    setForgetFlag(flag) { this.forgetFlag = flag },
    setSearchFlag(flag) { this.searchFlag = flag },

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
        const res = await api.getHomeData()
        if (res.code === 0) {
          this.blogInfo = res.data
          this.blogInfo.page_list?.map(e => (e.cover = convertImgUrl(e.cover)))
          this.blogInfo.blog_config.website_avatar = convertImgUrl(this.blogInfo.blog_config.website_avatar)
        } else {
          return Promise.reject(res)
        }
      } catch (err) {
        return Promise.reject(err)
      }
    },
  },
})

// export default useAppStore
