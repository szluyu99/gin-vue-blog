// import DefineOptions from 'unplugin-vue-define-options/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { NaiveUiResolver, Vuetify3Resolver } from 'unplugin-vue-components/resolvers'

/**
 * * unplugin-icons 插件，自动引入 iconify 图标
 * usage: https://github.com/antfu/unplugin-icons
 * 图标库: https://icones.js.org/
 */
import Icons from 'unplugin-icons/vite'
import IconsResolver from 'unplugin-icons/resolver'

export default [
  // DefineOptions(),
  // https://github.com/antfu/unplugin-auto-import
  AutoImport({
    imports: [
      'vue',
      'vue/macros',
      'vue-router',
      '@vueuse/core',
      'pinia',
    ],
    dts: true,
    dirs: ['./src/composables'],
    vueTemplate: true,
  }),
  Icons({
    compiler: 'vue3',
    scale: 1,
    defaultClass: 'inline-block',
  }),
  // https://github.com/antfu/vite-plugin-components
  Components({
    resolvers: [NaiveUiResolver(), IconsResolver(), Vuetify3Resolver()],
    dts: 'types/components.d.ts',
  }),
]
