import { resolve } from 'path'
import DefineOptions from 'unplugin-vue-define-options/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'

/**
 * * unplugin-icons 插件，自动引入 iconify 图标
 * usage: https://github.com/antfu/unplugin-icons
 * 图标库: https://icones.js.org/
 */
import Icons from 'unplugin-icons/vite'
import IconsResolver from 'unplugin-icons/resolver'
import { FileSystemIconLoader } from 'unplugin-icons/loaders'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

import { getSrcPath } from '../utils'
const customIconPath = resolve(getSrcPath(), 'assets/svg')

export default [
  // 在 <script setup> 中使用 Options API
  DefineOptions(),
  // 自动引入 Vue 官方 Api, 只在 .vue 中生效, .js 中不行
  AutoImport({
    imports: ['vue', 'vue-router', 'pinia', '@vueuse/core', 'vue/macros'],
    dts: false,
  }),
  // 自动引入组件: 例如各种 UI 库
  Components({
    resolvers: [
      NaiveUiResolver(),
      // Iconify
      IconsResolver({ customCollections: ['custom'], componentPrefix: 'icon' }),
    ],
    dts: false,
  }),
  // 以组件的形式使用图标
  Icons({
    compiler: 'vue3',
    customCollections: {
      custom: FileSystemIconLoader(customIconPath),
    },
    scale: 1,
    defaultClass: 'inline-block',
  }),
  // 生成 svg 雪碧图
  createSvgIconsPlugin({
    iconDirs: [customIconPath], // 需要缓存的图标文件夹
    symbolId: 'icon-custom-[dir]-[name]', // symboId 格式
    inject: 'body-last', // 自定义插入位置
    customDomId: '__CUSTOM_SVG_ICON__',
  }),
]
