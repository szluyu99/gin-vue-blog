import { defineStore } from 'pinia'
import { useDark } from '@vueuse/core'

const isDark = useDark()
export const useThemeStore = defineStore('theme-store', {
  persist: {
    key: 'gvb_admin_theme',
    paths: ['collapsed', 'watermarked'],
  },
  state: () => ({
    collapsed: false, // 侧边栏折叠
    watermarked: false, // 水印
    darkMode: isDark, // 黑暗模式
  }),
  actions: {
    switchWatermark() {
      this.watermarked = !this.watermarked
    },
    switchCollapsed() {
      this.collapsed = !this.collapsed
    },
    switchDarkMode() {
      this.darkMode = !this.darkMode
    },
  },
})
