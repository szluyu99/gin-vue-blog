import { defineStore } from 'pinia'
import { router } from '@/router'
import { sStorage } from '@/utils'

// 激活标签
const activeTag = sStorage.get('activeTag')
// 获取所有标签
const tags = sStorage.get('tags')
// 不添加标签的路由路径
const WITHOUT_TAG_PATHS = ['/404', '/login']

export const useTagsStore = defineStore('tag', {
  state() {
    return {
      tags: tags || [], // 显示的所有标签
      activeTag: activeTag || '', // 当前激活的标签
    }
  },
  actions: {
    setActiveTag(path) {
      // console.log('useTagsStore, 当前激活标签: ', path)
      this.activeTag = path
      sStorage.set('activeTag', path)
    },
    setTags(tags) {
      // tags 是数组
      this.tags = tags
      sStorage.set('tags', tags)
    },
    /**
     * 添加标签(不添加白名单中 和 已存在的)
     * @param {*} tag { name, path, title }
     */
    addTag(tag = {}) {
      this.setActiveTag(tag.path)
      // 在 白名单中 或 标签栏中已经有的路径 不添加
      if (WITHOUT_TAG_PATHS.includes(tag.path) || this.tags.some((item) => item.path === tag.path))
        return
      this.setTags([...this.tags, tag])
    },
    // 移除标签(如果只有一个标签是无法移除的)
    removeTag(path) {
      // 移除的是当前标签
      if (path === this.activeTag) {
        // 当前激活标签的索引
        const activeIndex = this.tags.findIndex((item) => item.path === path)
        // 取的到左边就取左边的标签
        if (activeIndex > 0) router.push(this.tags[activeIndex - 1].path)
        // 取不到左边取右边的标签
        else router.push(this.tags[activeIndex + 1].path)
      }
      this.setTags(this.tags.filter((tag) => tag.path !== path))
    },
    // 关闭其他标签, 只留下当前标签
    removeOther(curPath = this.activeTag) {
      this.setTags(this.tags.filter((tag) => tag.path === curPath))
      // 如果点击的不是当前标签, 则跳转到那个标签
      if (curPath !== this.activeTag) {
        router.push(this.tags[this.tags.length - 1].path)
      }
    },
    // 关闭左侧
    removeLeft(curPath) {
      const curIndex = this.tags.findIndex((item) => item.path === curPath)
      // 过滤出右边的标签
      const filterTags = this.tags.filter((item, index) => index >= curIndex)
      this.setTags(filterTags)
      // 如果当前浏览的标签被关闭, 打开一个新标签
      if (!filterTags.find((item) => item.path === this.activeTag)) {
        router.push(filterTags[filterTags.length - 1].path)
      }
    },
    // 关闭右侧
    removeRight(curPath) {
      const curIndex = this.tags.findIndex((item) => item.path === curPath)
      // 过滤出左边的标签
      const filterTags = this.tags.filter((item, index) => index <= curIndex)
      this.setTags(filterTags)
      // 如果当前浏览的标签被关闭, 打开一个新标签
      if (!filterTags.find((item) => item.path === this.activeTag)) {
        router.push(filterTags[filterTags.length - 1].path)
      }
    },
    resetTags() {
      this.setTags([])
      this.setActiveTag('')
    },
  },
})
