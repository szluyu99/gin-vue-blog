import { defineStore } from 'pinia'
import api from '@/api'

/*
站内通知

单独一个 store 而不是塞进 user: 未读数需要在头部铃铛、移动端侧边栏等多处共享,
而且退出登录时要连同列表一起清空(user.$reset() 只管用户自己的字段)。

不做轮询: 博客的通知量很低, 打开下拉时拉一次、发完评论刷新一次就够了,
挂一个 setInterval 反而是每个访客每分钟给后端加一次请求。
*/
export const useNotificationStore = defineStore('notification', {
  state: () => ({
    unreadCount: 0,
    list: [],
    total: 0,
    loading: false,
  }),
  actions: {
    reset() {
      this.$reset()
    },
    async fetchUnreadCount() {
      try {
        const resp = await api.getUnreadNotificationCount()
        // 后端直接返回数字, 失败时不要把 unreadCount 变成 undefined
        this.unreadCount = Number(resp.data) || 0
      }
      catch (err) {
        console.error(err)
      }
    },
    async fetchList(pageSize = 10) {
      this.loading = true
      try {
        const resp = await api.getNotifications({ page_num: 1, page_size: pageSize })
        this.list = resp.data?.page_data ?? []
        this.total = resp.data?.total ?? 0
      }
      catch (err) {
        console.error(err)
      }
      finally {
        this.loading = false
      }
    },
    // ids 为空表示全部已读, 与后端一致
    async read(ids = []) {
      try {
        await api.readNotifications(ids)
        for (const item of this.list) {
          if (!ids.length || ids.includes(item.id)) {
            item.is_read = true
          }
        }
        // 未读数从后端重新取, 而不是本地减: 别的标签页也可能已经读过了
        await this.fetchUnreadCount()
      }
      catch (err) {
        console.error(err)
      }
    },
  },
})
