import { defineStore } from 'pinia'
import { nextTick } from 'vue'
import api from '@/api'

// 应用全局变量
export const useAppStore = defineStore('app', {
  persist: true,
  state() {
    return {
      reloadFlag: true,
      aliveKeys: {}, // keepAlive 路由的 key, 重新赋值可重置 keepAlive
      blogConfig: {}, // 博客设置信息
    }
  },
  actions: {
    setAliveKeys(key, val) {
      this.aliveKeys[key] = val
    },
    // 刷新页面: 效果并非按 F5 刷新整个网页, 而是模拟刷新 (nextTick + 滚动到顶部)
    async reloadPage() {
      window.$loadingBar.start()
      // 配合 v-if="reloadFlag" 有白屏效果
      this.reloadFlag = false
      await nextTick() // 将回调延迟到下次 DOM 更新循环之后执行
      this.reloadFlag = true
      // 以下做法也可以实现白屏效果
      // setTimeout(() => (this.reloadFlag = true), 1000)

      // 滚动到顶部, 模拟刷新
      setTimeout(() => {
        document.documentElement.scrollTo({ left: 0, top: 0 })
        window.$loadingBar.finish()
      }, 100)
    },
    // 加载博客设置信息
    async getBlogInfo() {
      try {
        const res = await api.getBlogConfig()
        if (res.code === 0)
          this.blogConfig = res.data
        else return Promise.reject(res)
      }
      catch (err) {
        return Promise.reject(err)
      }
    },
  },
})
