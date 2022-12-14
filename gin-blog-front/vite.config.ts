import path from 'path'
import type { ConfigEnv } from 'vite'
import { defineConfig, loadEnv } from 'vite'
import { convertEnv } from './build/utils'
import { createViteProxy } from './build/config'
import { setupVitePlugins } from './build/plugins'

export default defineConfig((configEnv: ConfigEnv) => {
  // 从环境变量中加载值
  const isBuild = configEnv.command === 'build' // 判断是否是打包
  const viteEnv = convertEnv(loadEnv(configEnv.mode, process.cwd()))
  const { VITE_PORT, VITE_PUBLIC_PATH, VITE_USE_PROXY, VITE_PROXY_TYPE } = viteEnv

  return {
    base: VITE_PUBLIC_PATH,
    resolve: {
      alias: {
        '@': path.resolve(path.resolve(process.cwd()), 'src'),
        '~': path.resolve(process.cwd()),
      },
    },
    plugins: setupVitePlugins(viteEnv, isBuild),
    server: {
      host: '0.0.0.0',
      port: VITE_PORT,
      open: false,
      // 后端可以解决, 前端就不用设置代理
      // 设置代理后 .env 文件中的 VITE_API 得切换
      proxy: createViteProxy(VITE_USE_PROXY, VITE_PROXY_TYPE as ProxyType),
    },
    optimizeDeps: {
      include: ['@kangc/v-md-editor/lib/theme/vuepress.js'],
    },
    build: {
      reportCompressedSize: false,
      sourcemap: false,
      chunkSizeWarningLimit: 1024, // chunk 大小警告的限制 (单位 kb)
      commonjsOptions: {
        ignoreTryCatch: false,
      },
      terserOptions: {
        // 打包后移除 console 和 注释
        compress: {
          drop_console: true,
          drop_debugger: true,
        },
      },
    },
  }
})
