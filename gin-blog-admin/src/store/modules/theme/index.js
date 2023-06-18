import { defineStore } from 'pinia'
import { useDark } from '@vueuse/core'

const isDark = useDark()
export const useThemeStore = defineStore('theme-store', {
  persist: true, // 刷新保持
  state() {
    return {
      collapsed: false, // 侧边栏折叠
      watermarkFlag: false, // 水印
      darkMode: isDark, // 黑暗模式
    }
  },
  actions: {
    switchWatermark() {
      this.watermarkFlag = !this.watermarkFlag
    },
    switchCollapsed() {
      this.collapsed = !this.collapsed
    },
    switchDarkMode() {
      this.darkMode = !this.darkMode
    },
  },
})
