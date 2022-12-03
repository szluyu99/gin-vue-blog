import type { PluginOption } from 'vite'
import Vue from '@vitejs/plugin-vue'
import Vuetify from 'vite-plugin-vuetify'
import Pages from 'vite-plugin-pages'
import unocss from 'unocss/vite'
import unplugins from './unplugin'

export function setupVitePlugins(): PluginOption[] {
  return [
    // https://www.npmjs.com/package/@vitejs/plugin-vue
    Vue({ reactivityTransform: true }),
    // https://www.npmjs.com/package/vite-plugin-vuetify
    Vuetify({ autoImport: true }),
    // https://github.com/hannoeru/vite-plugin-pages
    Pages(),
    // https://github.com/antfu/unocss
    unocss(),
    ...unplugins,
  ]
}
