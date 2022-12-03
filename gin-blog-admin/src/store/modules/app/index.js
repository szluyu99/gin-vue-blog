import { useSettingAPI } from '@/api'
const { getBlogConfig } = useSettingAPI()

// 存放一些全局应用变量
export const useAppStore = defineStore('app', {
  state() {
    return {
      reloadFlag: true,
      collapsed: false,
      // keepAlive 路由的 key, 重新赋值可重置 keepAlive
      aliveKeys: {},
      blogConfig: {}, // 博客设置信息
    }
  },
  actions: {
    // 刷新页面: 效果并非按 F5 刷新整个网页, 而是模拟刷新(nextTick + 滚动到顶部)
    async reloadPage() {
      $loadingBar.start()

      // 配合 v-if="reloadFlag" 有白屏效果
      this.reloadFlag = false
      await nextTick() // 将回调延迟到下次 DOM 更新循环之后执行
      this.reloadFlag = true
      // 以下做法也可以实现白屏效果
      // setTimeout(() => (this.reloadFlag = true), 1000)

      setTimeout(() => {
        // 滚动到顶部, 模拟刷新
        document.documentElement.scrollTo({ left: 0, top: 0 })
        $loadingBar.finish()
      }, 100)
    },
    // 切换页面展开
    switchCollapsed() {
      this.collapsed = !this.collapsed
    },
    setCollapsed(collapsed) {
      this.collapsed = collapsed
    },
    setAliveKeys(key, val) {
      this.aliveKeys[key] = val
    },
    // 加载博客设置信息
    async getBlogInfo() {
      try {
        const res = await getBlogConfig()
        if (res.code === 0) {
          // console.log(res.data)
          this.blogConfig = res.data
        } else return Promise.reject(res)
      } catch (err) {
        return Promise.reject(err)
      }
    },
  },
})
