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
    /*
      只保留失败用例的 console 输出

      很多用例故意让接口 reject, 来验证「失败时不提示成功 / 会回滚」, 源码里的
      console.error 是正确行为, 但会在 CI 日志里堆出一屏 "Error: boom" 堆栈,
      真正失败的那条反而被埋掉。
      用 passed-only 而不是 true: 用例一旦失败, 它的日志照样打出来, 排查信息不丢。
    */
    silent: 'passed-only',
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
