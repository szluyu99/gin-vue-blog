import { defineStore } from 'pinia'
import { resetRouter } from '@/router'
import { convertImgUrl, removeToken, toLogin } from '@/utils'
import { usePermissionStore, useTagsStore } from '@/store'
import api from '@/api'

// 用户全局变量
export const useUserStore = defineStore('user', {
  state() {
    return {
      userInfo: {},
    }
  },
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
      useTagsStore().resetTags()
      usePermissionStore().resetPermission()
      resetRouter()
      this.$reset()

      toLogin()
      window.$message.success('您已经退出登录!')
    },
    // * 被强制退出: 被动行为, 不需要调用退出登录接口
    async forceOffline() {
      removeToken()
      useTagsStore().resetTags()
      usePermissionStore().resetPermission()
      resetRouter()
      this.$reset()

      toLogin()
      window.$message.error('您已经被强制下线!')
    },
  },
})
