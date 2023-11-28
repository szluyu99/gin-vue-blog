import { defineStore } from 'pinia'
import { convertImgUrl, getToken, removeToken } from '@/utils'
import api from '@/api'

export const useUserStore = defineStore('user', {
  persist: true,
  state() {
    return {
      userInfo: {},
    }
  },
  getters: {
    userId: state => state.userInfo.id ?? '',
    nickname: state => state.userInfo.nickname ?? '',
    avatar: state => state.userInfo.avatar ?? 'https://static.talkxj.com/avatar/user.png',
    website: state => state.userInfo.website ?? '',
    intro: state => state.userInfo.intro ?? '',
    email: state => state.userInfo.email ?? '',
    articleLikeSet() { return this.userInfo.articleLikeSet || [] },
    commentLikeSet() { return this.userInfo.commentLikeSet || [] },
  },
  actions: {
    async getUserInfo() {
      const token = getToken()
      if (!token)
        return

      try {
        const res = await api.getUser()
        if (res.code === 0) {
          const {
            id,
            nickname,
            avatar,
            website,
            intro,
            email,
            article_like_set,
            comment_like_set,
          } = res.data
          this.userInfo = {
            id,
            nickname,
            avatar,
            website,
            intro,
            email,
            articleLikeSet: article_like_set.map(e => +e),
            commentLikeSet: comment_like_set.map(e => +e),
          }
          this.userInfo.avatar = convertImgUrl(this.userInfo.avatar)
          return Promise.resolve(res.data)
        } else {
          return Promise.reject(res)
        }
      } catch (error) {
        return Promise.reject(error)
      }
    },
    async logout() {
      removeToken()
      this.$reset()
    },
    commentLike(commentId) {
      this.commentLikeSet.includes(commentId)
        ? this.commentLikeSet.splice(this.commentLikeSet.indexOf(commentId), 1)
        : this.commentLikeSet.push(commentId)
    },
    articleLike(articleId) {
      this.articleLikeSet.includes(articleId)
        ? this.articleLikeSet.splice(this.articleLikeSet.indexOf(articleId), 1)
        : this.articleLikeSet.push(articleId)
    },
  },
})

