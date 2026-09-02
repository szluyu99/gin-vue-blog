import path from 'node:path'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vitest/config'

// 单测配置与 vite.config.js 分开: 测试不需要 unocss / gzip / visualizer 这些构建插件
export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(process.cwd(), 'src'),
      '~': path.resolve(process.cwd()),
    },
  },
  plugins: [vue()],
  test: {
    // sessionStorage / document 等浏览器 API 需要 jsdom
    environment: 'jsdom',
    include: ['src/**/*.spec.js'],
    // 源码里通过 import.meta.env 读取, 测试里给一份固定值, 断言才能稳定
    env: {
      VITE_BASE_API: '/api',
      VITE_SERVER_URL: 'http://test-server',
      VITE_USE_MOCK: 'false',
    },
    coverage: {
      include: ['src/utils/**', 'src/store/**'],
    },
  },
})
