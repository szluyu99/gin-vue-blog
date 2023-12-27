import { defineStore } from 'pinia'
import { convertImgUrl } from '@/utils'
import api from '@/api'

// 用户全局变量
export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {
      id: null,
      nickname: '',
      intro: '',
      website: '',
      avatar: '',
      roles: [],
    },
  }),
  getters: {
    userId: state => state.userInfo.id,
    nickname: state => state.userInfo.nickname,
    intro: state => state.userInfo.intro,
    website: state => state.userInfo.website,
    avatar: state => convertImgUrl(state.userInfo.avatar),
    roles: state => state.userInfo.roles,
  },
  actions: {
    async getUserInfo() {
      try {
        const resp = await api.getUserInfo()
        const { id, nickname, avatar, intro, website } = resp.data
        this.userInfo = { id, nickname, avatar, intro, website }
        return Promise.resolve(resp.data)
      }
      catch (err) {
        return Promise.reject(err)
      }
    },
    setUserInfo(user) {
      Object.keys(user).forEach((key) => {
        this.userInfo[key] = user[key]
      })
    },
    resetUser() {
      this.$reset()
    },
  },
})

/**
 * 携带参数跳转到登录页面
 */
// function toLoginWithQuery() {
//   const currentRoute = unref(router.currentRoute)
// 跳转回去时记录 redirect 到 query 上
//   const needRedirect = !currentRoute.meta.requireAuth && !['/404', '/login'].includes(currentRoute.path)
//   router.replace({
//     path: '/login',
//     query: needRedirect ? { ...currentRoute.query, redirect: currentRoute.path } : {},
//   })
// }
