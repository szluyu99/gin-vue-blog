import type { PluginOption } from 'vite'
import Vue from '@vitejs/plugin-vue'
// import Vuetify from 'vite-plugin-vuetify'
import Pages from 'vite-plugin-pages'
import unocss from 'unocss/vite'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'
import unplugins from './unplugin'

export function setupVitePlugins(viteEnv: ViteEnv, isBuild: boolean): PluginOption[] {
  const plugins = [
    // https://www.npmjs.com/package/@vitejs/plugin-vue
    // 响应式语法糖: https://cn.vuejs.org/guide/extras/reactivity-transform.html
    Vue({ reactivityTransform: true }),
    // https://www.npmjs.com/package/vite-plugin-vuetify
    // Vuetify({ autoImport: true }),
    // https://github.com/hannoeru/vite-plugin-pages
    Pages(),
    // https://github.com/antfu/unocss
    unocss(),
    ...unplugins,
  ]
  // 压缩插件: https://github.com/vbenjs/vite-plugin-compression
  if (viteEnv.VITE_USE_COMPRESS) {
    plugins.push(
      viteCompression({ algorithm: 'gzip' }),
    )
  }
  // 打包可视化插件
  isBuild && plugins.push(visualizer({
    open: true,
    gzipSize: true,
    brotliSize: true,
  }))

  return plugins
}
