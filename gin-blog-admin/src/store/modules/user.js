import { unref } from 'vue'
import { defineStore } from 'pinia'
import { resetRouter, router } from '@/router'
import { convertImgUrl, removeToken } from '@/utils'
import { usePermissionStore, useTagStore } from '@/store'
import api from '@/api'

// 用户全局变量
export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {},
  }),
  getters: {
    userId: state => state.userInfo.id,
    nickname: state => state.userInfo.nickname,
    intro: state => state.userInfo.intro,
    website: state => state.userInfo.website,
    avatar: state => convertImgUrl(state.userInfo.avatar),
  },
  actions: {
    // 获取用户信息
    async getUserInfo() {
      try {
        const res = await api.getUser()
        const { id, nickname, avatar, intro, website } = res.data
        this.userInfo = { id, nickname, avatar, intro, website }
        return Promise.resolve(res.data)
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
    // * 退出登录: 主动行为, 需要调用退出登录接口
    async logout() {
      // * 必须先调用退出登录接口, 再清除本地 Token, 因为退出登录接口需要 Token
      await api.logout()

      removeToken()
      useTagStore().resetTags()
      usePermissionStore().resetPermission()
      resetRouter()
      this.$reset()

      toLoginWithQuery()
      window.$message.success('您已经退出登录!')
    },
    // * 被强制退出: 被动行为, 不需要调用退出登录接口
    async forceOffline() {
      removeToken()
      useTagStore().resetTags()
      usePermissionStore().resetPermission()
      resetRouter()
      this.$reset()

      toLoginWithQuery()
      window.$message.error('您已经被强制下线!')
    },
  },
})

/**
 * 携带参数跳转到登录页面
 */
function toLoginWithQuery() {
  const currentRoute = unref(router.currentRoute)
  // 跳转回去时记录 redirect 到 query 上
  const needRedirect = !currentRoute.meta.requireAuth && !['/404', '/login'].includes(currentRoute.path)
  router.replace({
    path: '/login',
    query: needRedirect ? { ...currentRoute.query, redirect: currentRoute.path } : {},
  })
}
