import { defineStore } from 'pinia'
import { resetRouter } from '@/router'
import { removeToken, toLogin } from '@/utils'
import { usePermissionStore, useTagsStore } from '@/store'
import api from '@/api'

export const useUserStore = defineStore('user', {
  state() {
    return {
      userInfo: {},
    }
  },
  getters: {
    userId() {
      return this.userInfo.id
    },
    nickname() {
      return this.userInfo.nickname
    },
    avatar() {
      const SERVER_URL = 'http://localhost:8765/'
      if (this.userInfo.avatar.startsWith('http'))
        return this.userInfo.avatar
      return `${SERVER_URL}${this.userInfo.avatar}`
    },
    intro() {
      return this.userInfo.intro
    },
    website() {
      return this.userInfo.website
    },
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
    // 退出登录: 主动行为, 需要调用退出登录接口
    async logout() {
      // * 必须先调用退出登录接口, 再清除本地 Token, 因为退出登录接口需要 Token
      await api.logout()

      window.$message.success('您已经退出登录!')
      removeToken()
      useTagsStore().resetTags()
      usePermissionStore().resetPermission()
      resetRouter()
      this.$reset()

      toLogin()
    },
    // 被强制退出: 被动行为, 不需要调用退出登录接口
    async forceOffline() {
      window.$message.error('您已经被强制下线!')
      removeToken()
      useTagsStore().resetTags()
      usePermissionStore().resetPermission()
      resetRouter()
      this.$reset()

      toLogin()
    },
  },
})
