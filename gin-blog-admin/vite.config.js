import { resolve } from 'node:path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import unocss from 'unocss/vite'
import visualizer from 'rollup-plugin-visualizer'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) }

  return {
    base: process.env.VITE_PUBLIC_PATH || '/',
    resolve: {
      alias: { '@': resolve(__dirname, 'src') },
    },
    plugins: [
      vue(),
      unocss(),
      visualizer({
        open: true,
        gzipSize: true,
        brotliSize: true,
      }),
    ],
    server: {
      hmr: true,
      host: '0.0.0.0',
      port: 3000,
      proxy: {
        '/api': {
          target: 'http://localhost:8765',
          changeOrigin: true,
        },
      },
    },
    // https://cn.vitejs.dev/guide/api-javascript.html#build
    build: {
      target: 'es2015',
      reportCompressedSize: false,
      chunkSizeWarningLimit: 1024, // chunk 大小警告的限制（单位kb）
    },
    esbuild: {
      drop: ['console', 'debugger'],
    },
  }
})
