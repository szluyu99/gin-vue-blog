import { defineStore } from 'pinia'
import { router } from '@/router'
import { sStorage } from '@/utils'

/**
 * @type {string} activeTag 当前激活的标签
 * @description 从 sessionStorage 中获取当前激活的标签
 */
const activeTag = sStorage.get('activeTag')

/**
 * @type {string[]} tags 所有标签
 * @description 从 sessionStorage 中获取所有标签
 */
const tags = sStorage.get('tags')

/**
 * @type {string[]} WITHOUT_TAG_PATHS
 * @description 不添加标签的路由路径
 */
const WITHOUT_TAG_PATHS = ['/404', '/login']

export const useTagsStore = defineStore('tag', {
  state() {
    return {
      tags: tags || [], // 显示的所有标签
      activeTag: activeTag || '', // 当前激活的标签
    }
  },
  getters: {
    activeIndex: state => state.tags.findIndex(item => item.path === state.activeTag),
  },
  actions: {
    /**
     * 设置当前激活的标签
     * @param {string} path 标签对应的路由路径
     */
    setActiveTag(path) {
      this.activeTag = path
      sStorage.set('activeTag', path)
    },
    /**
     * 设置当前显示的所有标签
     * @param {string[]} tags 数组
     */
    setTags(tags) {
      this.tags = tags
      sStorage.set('tags', tags)
    },
    /**
     * 添加标签 (不添加白名单中 和 已存在的)
     * @param {{ name, path, title, icon }} tag 标签对象
     */
    addTag(tag = {}) {
      this.setActiveTag(tag.path)
      // 在 白名单中 或 标签栏中已经有的路径 不添加
      if (WITHOUT_TAG_PATHS.includes(tag.path) || this.tags.some(item => item.path === tag.path))
        return
      this.setTags([...this.tags, tag])
    },
    /**
     * 移除标签 (如果只有一个标签, 无法移除)
     * @param {string} path 标签对应的路由路径
     */
    removeTag(path) {
      // 移除的是当前标签, 计算跳转到哪个标签
      if (path === this.activeTag) {
        // 不是第一个就选左边的标签, 否则选右边的标签
        this.activeIndex > 0
          ? router.push(this.tags[this.activeIndex - 1].path)
          : router.push(this.tags[this.activeIndex + 1].path)
      }
      this.setTags(this.tags.filter(tag => tag.path !== path))
    },
    /**
     * 关闭其他标签, 只留下当前标签
     * @param {string} path
     */
    removeOther(path = this.activeTag) {
      this.setTags(this.tags.filter(tag => tag.path === path))
      // 如果点击的不是当前标签, 则跳转到那个标签
      if (path !== this.activeTag)
        router.push(this.tags[0].path) // 关闭其他后只剩一个标签
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
      if (!filterTags.find(item => item.path === this.activeTag)) // 剩余列表中找不到当前标签
        router.push(filterTags[filterTags.length - 1].path)
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
      if (!filterTags.find(item => item.path === this.activeTag)) // 剩余列表中找不到当前标签
        router.push(filterTags[filterTags.length - 1].path)
    },
    resetTags() {
      this.setTags([])
      this.setActiveTag('')
    },
  },
})
