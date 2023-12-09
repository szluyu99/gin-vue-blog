import { defineStore } from 'pinia'
import { convertImgUrl } from '@/utils'
import api from '@/api'

export const useAppStore = defineStore('app', {
  persist: true,
  state: () => ({
    searchFlag: false,
    loginFlag: false,
    registerFlag: false,

    collapsed: false, // 侧边栏折叠（移动端）

    blogInfo: {
      blog_config: {
        website_name: '阵、雨的个人博客',
        website_author: '阵、雨',
        website_intro: '往事随风而去',
        website_avatar: '',
      },
      pageList: [],
    },
  }),
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
    setSearchFlag(flag) { this.searchFlag = flag },

    async getBlogInfo() {
      try {
        const resp = await api.getHomeData()
        if (resp.code === 0) {
          this.blogInfo = resp.data
          this.blogInfo.page_list?.map(e => (e.cover = convertImgUrl(e.cover)))
          this.blogInfo.blog_config.website_avatar = convertImgUrl(this.blogInfo.blog_config.website_avatar)
        }
        else {
          return Promise.reject(resp)
        }
      }
      catch (err) {
        return Promise.reject(err)
      }
    },
  },
})
