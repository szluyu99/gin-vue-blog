import path from 'path'
import { defineConfig } from 'vite'
import { setupVitePlugins } from './build/plugins'

export default defineConfig(() => {
  return {
    base: '/',
    resolve: {
      alias: {
        '@': path.resolve(path.resolve(process.cwd()), 'src'),
        '~': path.resolve(process.cwd()),
      },
    },
    plugins: setupVitePlugins(),
    server: {
      host: '0.0.0.0',
      open: false,
      hmr: {
        overlay: false,
      },
      // 后端可以解决, 前端就不用设置代理
      // 设置代理后 .env 文件中的 VITE_API 得切换
      // proxy: {
      //   '/api/front': {
      //     target: 'http://47.103.138.81/api/front', // 线上环境
      //     target: 'http://localhost:5678/api/front', // 本地环境
      //     changeOrigin: true,
      //     rewrite: path => path.replace(/^\/api\/front/, ''),
      //   },
      // },
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
    },
  }
})
