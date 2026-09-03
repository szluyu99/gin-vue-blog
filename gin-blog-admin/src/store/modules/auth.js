import { defineStore } from 'pinia'
import { unref } from 'vue'
import api from '@/api'
import { resetRouter, router } from '@/router'
import { usePermissionStore, useTagStore, useUserStore } from '@/store'

export const useAuthStore = defineStore('auth', {
  persist: {
    key: 'gvb_admin_auth',
    pick: ['token'],
  },
  state: () => ({
    token: null,
  }),
  actions: {
    setToken(token) {
      this.token = token
    },
    toLogin() {
      const currentRoute = unref(router.currentRoute)
      router.replace({
        path: '/login',
        query: currentRoute.query,
      })
    },
    resetLoginState() {
      useUserStore().$reset()
      usePermissionStore().$reset()
      useTagStore().$reset()
      resetRouter()
      this.$reset()
    },
    /**
     * 主动退出登录
     */
    async logout() {
      // 退出接口失败也要清掉本地登录态并跳登录页, 否则 token 留在 localStorage 里、
      // 页面不跳转也没有提示, 用户会以为还登录着
      try {
        await api.logout()
      }
      catch (err) {
        console.error(err)
      }
      this.resetLoginState()
      this.toLogin()
      window.$message.success('您已经退出登录！')
    },
    /**
     * TODO: 被强制退出
     */
    async forceOffline() {
      this.resetLoginState()
      this.toLogin()
      // window.$message.error('您已经被强制下线！')
    },
  },
})

// function toLoginWithQuery() {
//   const currentRoute = unref(router.currentRoute)
//  // 跳转回去时记录 redirect 到 query 上
//   const needRedirect = !currentRoute.meta.requireAuth && !['/404', '/login'].includes(currentRoute.path)
//   router.replace({
//     path: '/login',
//     query: needRedirect ? { ...currentRoute.query, redirect: currentRoute.path } : {},
//   })
// }
