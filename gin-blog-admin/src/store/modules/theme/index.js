import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme-store', {
  persist: true,
  state() {
    return {
      collapsed: false, // 侧边栏折叠
      watermarkFlag: true, // 水印
    }
  },
  actions: {
    // 切换水印显示
    switchWatermark() {
      this.watermarkFlag = !this.watermarkFlag
    },
    // 切换页面展开
    switchCollapsed() {
      this.collapsed = !this.collapsed
    },
  },
})
