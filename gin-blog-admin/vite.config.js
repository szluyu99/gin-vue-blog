import { defineConfig, loadEnv } from 'vite'

import { convertEnv, getRootPath, getSrcPath } from './build/utils'
import { createViteProxy, viteDefine } from './build/config'
import { createVitePlugins } from './build/plugin'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  const isBuild = command === 'build' // pnpm build

  // 加载 .env 中的环境变量 (vite.config.js 中无法通过 import.meta.env 获取环境变量, 需要使用 process.env 获取)
  // 设置第三个参数为 '' 可以加载所有环境变量(包括本机的), 默认只会加载 'VITE_' 前缀的变量
  const viteEnv = convertEnv(loadEnv(mode, process.cwd())) // loadEnv(mode, process.cwd(), '')
  const { VITE_PUBLIC_PATH, VITE_PORT, VITE_USE_PROXY, VITE_PROXY_TYPE } = viteEnv

  return {
    base: VITE_PUBLIC_PATH || '/',
    resolve: {
      alias: {
        '~': getRootPath(),
        '@': getSrcPath(),
      },
    },
    define: viteDefine, // 全局常量 (不经过语法分析, 直接替换文本)
    plugins: createVitePlugins(viteEnv, isBuild),
    server: {
      hmr: true,
      host: '0.0.0.0',
      port: VITE_PORT,
      // 设置了代理, 查看 setting/proxy-config.js
      proxy: createViteProxy(VITE_USE_PROXY, VITE_PROXY_TYPE),
    },
    // https://cn.vitejs.dev/guide/api-javascript.html#build
    build: {
      target: 'es2015',
      outDir: 'dist',
      reportCompressedSize: false, // 禁用 gzip 压缩大小报告
      chunkSizeWarningLimit: 1024, // chunk 大小警告的限制（单位kb）
      terserOptions: {
        compress: {
          drop_console: true,
          drop_debugger: true,
        },
      },
      minify: 'terser',
    },
  }
})
