import vue from '@vitejs/plugin-vue'

/**
 * * unocss插件，原子css
 * https://github.com/antfu/unocss
 */
import Unocss from 'unocss/vite'

// TODO 打包分析插件 rollup-plugin-visualizer
// TODO 压缩插件 vite-plugin-compression

import unplugin from './unplugin'
// import { configMockPlugin } from './mock'
import { configHtmlPlugin } from './html'

export function createVitePlugins(viteEnv, isBuild) {
  const plugins = [vue(), ...unplugin, Unocss(), configHtmlPlugin(viteEnv, isBuild)]
  // 读取 .env 中的配置
  // viteEnv?.VITE_USE_MOCK && plugins.push(configMockPlugin(isBuild))
  return plugins
}
