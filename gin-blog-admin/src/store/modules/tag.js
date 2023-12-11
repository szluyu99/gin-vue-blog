import { nextTick } from 'vue'
import { defineStore } from 'pinia'
import { router } from '@/router'

export const useTagStore = defineStore('tag', {
  persist: {
    key: 'gvb_admin_tag',
    paths: ['tags'],
    storage: window.sessionStorage,
  },
  state: () => ({
    tags: [], // 标签栏的所有标签
    activeTag: '', // 当前激活的标签 path
    reloading: true, // 是否正在刷新
    /**
     * ! keepAlive 路由的 key, 重新赋值可重置 keepAlive
     * key 是 route name
     */
    aliveKeys: {},
  }),
  getters: {
    // 获取当前激活的标签的索引
    activeIndex: state => state.tags.findIndex(tag => tag.path === state.activeTag),
  },
  actions: {
    /**
     * 更新 keepAlive 路由, 让其重新渲染
     * @param {string} name route name
     */
    updateAliveKey(name) {
      this.aliveKeys[name] = (+new Date())
    },
    /**
     * 设置当前激活的标签
     * @param {string} path 标签对应的路由路径
     */
    async setActiveTag(path) {
      await nextTick() // 将回调延迟到下次 DOM 更新循环之后执行
      this.activeTag = path
    },
    /**
     * 设置当前显示的所有标签
     * @param {string[]} tags 数组
     */
    setTags(tags) {
      this.tags = tags
    },
    /**
     * 添加标签 (不添加白名单中 和 已存在的)
     * @param {{ name, path, title, icon }} tag 标签对象
     */
    addTag(tag = {}) {
      const index = this.tags.findIndex(item => item.path === tag.path)
      if (index !== -1) {
        this.tags.splice(index, 1, tag)
      }
      else {
        this.setTags([...this.tags, tag])
      }
      this.setActiveTag(tag.path)
    },
    /**
     * 移除标签 , 如果只有一个标签, 无法移除
     * @param {string} path 标签对应的路由路径
     */
    removeTag(path) {
      // 如果关闭的是当前标签
      if (path === this.activeTag) {
        if (this.activeIndex === 0) { // 如果是第一个标签, 则选中第二个标签
          router.push(this.tags[1].path)
        }
        else { // 否则选中左边的标签
          router.push(this.tags[this.activeIndex - 1].path)
        }
      }
      this.setTags(this.tags.filter(tag => tag.path !== path))
    },
    /**
     * 关闭其他标签
     * @param {string} path
     */
    removeOther(path = this.activeTag) {
      this.setTags(this.tags.filter(tag => tag.path === path))
      // 如果点击的不是当前标签, 会将当前标签关闭, 那么跳转到第一个标签
      if (path !== this.activeTag) {
        router.push(this.tags[0].path) // 关闭其他后只剩一个标签
      }
    },
    /**
     * 关闭左侧标签
     * @param {string} path
     */
    removeLeft(path) {
      const curIndex = this.tags.findIndex(item => item.path === path)
      // 过滤出右边的标签
      const filterTags = this.tags.filter((item, index) => index >= curIndex)
      this.setTags(filterTags)
      // 如果当前浏览的标签被关闭, 打开一个新标签
      if (!filterTags.find(item => item.path === this.activeTag)) {
        router.push(filterTags[filterTags.length - 1].path)
      }
    },
    /**
     * 关闭左侧标签
     * @param {string} path
     */
    removeRight(path) {
      const curIndex = this.tags.findIndex(item => item.path === path)
      // 过滤出左边的标签
      const filterTags = this.tags.filter((item, index) => index <= curIndex)
      this.setTags(filterTags)
      // 如果当前浏览的标签被关闭, 打开一个新标签
      if (!filterTags.find(item => item.path === this.activeTag)) {
        router.push(filterTags[filterTags.length - 1].path)
      }
    },
    /**
     * 重置标签
     */
    resetTags() {
      this.$reset()
    },
    /**
     * 刷新页面
     * @description 效果并非按 F5 刷新整个网页, 而是模拟刷新 (nextTick + 滚动到顶部)
     */
    async reloadTag() {
      window.$loadingBar.start()

      // 配合 v-if="reloadFlag" 实现白屏效果
      this.reloadFlag = false
      await nextTick() // 将回调延迟到下次 DOM 更新循环之后执行
      this.reloadFlag = true

      // 滚动到顶部, 模拟刷新
      setTimeout(() => {
        document.documentElement.scrollTo({ left: 0, top: 0 })
        window.$loadingBar.finish()
      }, 100)
    },
  },
})
