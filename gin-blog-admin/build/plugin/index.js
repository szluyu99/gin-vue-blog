import vue from '@vitejs/plugin-vue'

// https://github.com/antfu/unocss
import Unocss from 'unocss/vite'

// 压缩
import viteCompression from 'vite-plugin-compression'

// rollup 打包分析插件
import visualizer from 'rollup-plugin-visualizer'

import unplugin from './unplugin'
import { configHtmlPlugin } from './html'

export function createVitePlugins(viteEnv, isBuild) {
  const plugins = [
    vue({
      reactivityTransform: true, // 启用响应式语法糖
    }),
    ...unplugin,
    Unocss(),
    configHtmlPlugin(viteEnv, isBuild),
  ]

  viteEnv?.VITE_USE_COMPRESS && plugins.push(viteCompression({ algorithm: 'gzip' }))

  isBuild && plugins.push(
    visualizer({
      open: true,
      gzipSize: true,
      brotliSize: true,
    }))

  return plugins
}
