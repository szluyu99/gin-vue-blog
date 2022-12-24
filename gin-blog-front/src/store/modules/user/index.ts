import { convertImgUrl, getToken, removeToken } from '@/utils'
import api from '@/api'

interface UserInfo {
  id?: string
  nickname?: string
  avatar?: string
  website?: string
  intro?: string
  email?: string
  articleLikeSet?: [number] // 点赞过的文章集合
  commentLikeSet?: [number] // 点赞过的评论集合
}

export const useUserStore = defineStore('user', {
  persist: true,
  state() {
    return {
      userInfo: <UserInfo>{},
    }
  },
  getters: {
    userId: state => state.userInfo.id ?? '',
    nickname: state => state.userInfo.nickname ?? '',
    avatar: state => state.userInfo.avatar ?? 'https://static.talkxj.com/avatar/user.png',
    website: state => state.userInfo.website ?? '',
    intro: state => state.userInfo.intro ?? '',
    email: state => state.userInfo.email ?? '',
    articleLikeSet(): [number?] {
      return this.userInfo.articleLikeSet || []
    },
    commentLikeSet(): [number?] {
      return this.userInfo.commentLikeSet || []
    },
  },
  actions: {
    async getUserInfo() {
      // token 存在才会获取用户信息
      const token = getToken()
      if (!token)
        return

      try {
        const res: any = await api.getUser()
        if (res.code === 0) {
          const {
            id, nickname, avatar, website, intro, email,
            article_like_set, comment_like_set,
          } = res.data
          this.userInfo = {
            id,
            nickname,
            avatar,
            website,
            intro,
            email,
            articleLikeSet: article_like_set.map((e: string) => +e), // 注意后端给的是字符串数组
            commentLikeSet: comment_like_set.map((e: string) => +e),
          }
          this.userInfo.avatar = convertImgUrl(this.userInfo.avatar)
          return Promise.resolve(res.data)
        }
        else {
          return Promise.reject(res)
        }
      }
      catch (error) {
        return Promise.reject(error)
      }
    },
    async logout() {
      removeToken()
      this.$reset()
    },
    // 维护评论点赞
    // TODO: 如果点赞完每次都 getUserInfo 也可以, 会不会调用后台接口太频繁?
    commentLike(commentId: number) {
      this.commentLikeSet.includes(commentId)
        ? this.commentLikeSet.splice(this.commentLikeSet.indexOf(commentId), 1)
        : this.commentLikeSet.push(commentId)
    },
    // 维护文章点赞
    articleLike(articleId: number) {
      this.articleLikeSet.includes(articleId)
        ? this.articleLikeSet.splice(this.articleLikeSet.indexOf(articleId), 1)
        : this.articleLikeSet.push(articleId)
    },
  },
})
