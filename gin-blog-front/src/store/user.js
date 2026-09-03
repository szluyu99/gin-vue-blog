import { defineStore } from 'pinia'
import api from '@/api'

export const useUserStore = defineStore('user', {
  persist: {
    key: 'gvb_blog_user',
    pick: ['token'],
  },
  state: () => ({
    userInfo: {
      id: '',
      nickname: '',
      avatar: 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png',
      website: '',
      intro: '',
      email: '',
      articleLikeSet: [],
      commentLikeSet: [],
    },
    token: null,
  }),
  getters: {
    userId: state => state.userInfo.id ?? '',
    nickname: state => state.userInfo.nickname ?? '',
    // 用 || 而不是 ??: 头像为空串时也要退回默认图,
    // 否则 convertImgUrl('') 会给出已失效的 dummyimage.com 占位图
    avatar: state => state.userInfo.avatar || 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png',
    website: state => state.userInfo.website ?? '',
    intro: state => state.userInfo.intro ?? '',
    email: state => state.userInfo.email ?? '',
    articleLikeSet: state => state.userInfo.articleLikeSet || [],
    commentLikeSet: state => state.userInfo.commentLikeSet || [],
  },
  actions: {
    setToken(token) {
      this.token = token
    },
    resetLoginState() {
      this.$reset()
    },
    async logout() {
      await api.logout()
      this.$reset()
    },
    async getUserInfo() {
      if (!this.token) {
        return
      }
      try {
        const resp = await api.getUser()
        if (resp.code === 0) {
          const data = resp.data
          this.userInfo = {
            id: data.id,
            nickname: data.nickname,
            // 存原始相对路径, 展示时才 convertImgUrl:
            // 否则个人中心会把带域名的绝对地址当表单初值写回库
            avatar: data.avatar ?? '',
            website: data.website,
            intro: data.intro,
            email: data.email,
            articleLikeSet: data.article_like_set.map(e => +e),
            commentLikeSet: data.comment_like_set.map(e => +e),
          }
          return Promise.resolve(resp.data)
        }
        else {
          return Promise.reject(resp)
        }
      }
      catch (error) {
        return Promise.reject(error)
      }
    },
    // 注意要写 userInfo 里的数组, 不能写 commentLikeSet / articleLikeSet 这两个 getter:
    // 底层数组为空时 getter 返回的是新建的 [], push 进去的东西会被丢掉
    commentLike(commentId) {
      const set = (this.userInfo.commentLikeSet ??= [])
      set.includes(commentId)
        ? set.splice(set.indexOf(commentId), 1)
        : set.push(commentId)
    },
    articleLike(articleId) {
      const set = (this.userInfo.articleLikeSet ??= [])
      set.includes(articleId)
        ? set.splice(set.indexOf(articleId), 1)
        : set.push(articleId)
    },
  },
})
