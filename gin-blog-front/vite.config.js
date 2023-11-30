import path from 'path'
import vue from '@vitejs/plugin-vue'
import unocss from 'unocss/vite'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'
import { defineConfig, loadEnv } from 'vite'

export default defineConfig((configEnv) => {
  const env = loadEnv(configEnv.mode, process.cwd())

  return {
    base: env.VITE_PUBLIC_PATH,
    resolve: {
      alias: {
        '@': path.resolve(path.resolve(process.cwd()), 'src'), // 项目 src 路径
        '~': path.resolve(process.cwd()), // 项目根路径
      },
    },
    plugins: [
      vue(),
      unocss(),
      viteCompression({ algorithm: 'gzip' }),
      visualizer({ open: true, gzipSize: true, brotliSize: true }),
    ],
    server: {
      host: '0.0.0.0',
      port: env.VITE_PORT,
      open: false,
      proxy: {
        '/api/front': {
          target: 'http://localhost:5678', // TODO: 后端 URL
          changeOrigin: true,
        },
      },
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
