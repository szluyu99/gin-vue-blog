import api from '@/api'
import { convertImgUrl } from '@/utils'

// 页面
interface Page {
  name: string // 名称
  label: string // 标签
  cover: string // 封面
}

// 全局博客信息
interface BlogInfo {
  article_count?: number
  message_count?: number
  category_count?: number
  tag_count?: number
  view_count?: number

  blog_config: {
    website_name: string
    website_avatar: string
    website_author: string
    website_intro: string
    website_notice: string
    website_createtime: string
    website_record: string

    qq: string
    github: string
    gitee: string
  }

  page_list?: [Page] // 页面列表
}

export const useAppStore = defineStore('app', {
  persist: true, // pinia 持久化插件
  state() {
    return {
      // mobileFlag: false, // 是否是移动端
      reloadFlag: true,

      searchFlag: false, // 搜索
      loginFlag: false, // 登录
      registerFlag: false, // 注册
      forgetFlag: false, // 忘记密码

      collapsed: false, // 侧边栏折叠 (移动端)

      // TODO: 将变量拆分开来
      blogInfo: <BlogInfo>{
        blog_config: {
          website_name: '阵、雨的个人博客',
          website_author: '阵、雨',
          website_intro: '往事随风而去',
        },
      },
    }
  },
  getters: {
    // 是否是移动端
    isMobile: () => !!navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i),
    articleCount: state => state.blogInfo.article_count ?? 0,
    categoryCount: state => state.blogInfo.category_count ?? 0,
    tagCount: state => state.blogInfo.tag_count ?? 0,
    viewCount: state => state.blogInfo.view_count ?? 0,
    pageList: state => state.blogInfo.page_list ?? [], // 页面列表
    blogConfig: state => state.blogInfo.blog_config, // 博客信息
  },
  actions: {
    setCollapsed(flag: boolean) {
      this.collapsed = flag
    },
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
    // 判断当前是否是移动端
    // isMobile() {
    //   return navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i)
    // },
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
        if (res.code === 0) {
          this.blogInfo = res.data
          // 处理 page 的 cover 路径
          this.blogInfo.page_list?.map(e => e.cover = convertImgUrl(e.cover))
          // 处理 website_author
          this.blogInfo.blog_config.website_avatar = convertImgUrl(this.blogInfo.blog_config.website_avatar)
        }
        else { return Promise.reject(res) }
      }
      catch (err) {
        return Promise.reject(err)
      }
    },
  },
})
