import { defineStore } from 'pinia'
import { resetRouter } from '@/router'
import { removeToken, toLogin } from '@/utils'
import { useTagsStore, usePermissionStore } from '@/store'

import { useUserApi, useAuthApi } from '@/api'
const { getUser } = useUserApi()
const { logout } = useAuthApi()

export const useUserStore = defineStore('user', {
  state() {
    return {
      userInfo: {},
    }
  },
  getters: {
    userId() {
      return this.userInfo?.id
    },
    name() {
      return this.userInfo?.nickname
    },
    avatar() {
      return this.userInfo?.avatar
    },
  },
  actions: {
    // 获取用户信息
    async getUserInfo() {
      try {
        const res = await getUser()
        const { id, nickname, avatar } = res.data
        this.userInfo = { id, nickname, avatar }
        return Promise.resolve(res.data)
      } catch (err) {
        return Promise.reject(err)
      }
    },
    // 退出登录: 主动行为, 需要调用退出登录接口
    async logoutAdmin() {
      // * 必须先调用退出登录接口, 再清除本地 Token, 因为退出登录接口需要 Token
      await logout()

      $message.success('您已经退出登录!')
      removeToken()
      useTagsStore().resetTags()
      usePermissionStore().resetPermission()
      resetRouter()
      this.$reset

      toLogin()
    },
    // 被强制退出: 被动行为, 不需要调用退出登录接口
    async forceOffline() {
      $message.error('您已经被强制下线!')
      removeToken()
      useTagsStore().resetTags()
      usePermissionStore().resetPermission()
      resetRouter()
      this.$reset

      toLogin()
    },
  },
})
