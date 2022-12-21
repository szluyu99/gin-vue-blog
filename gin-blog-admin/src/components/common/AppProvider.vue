<script setup>
import { dateZhCN, useDialog, useLoadingBar, useMessage, useNotification, zhCN } from 'naive-ui'
// import { kebabCase } from 'lodash-es'

// 代码高亮插件
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'
import { useThemeStore } from '@/store'
import { naiveThemeOverrides } from '~/setting'
import { setupDialog, setupMessage } from '@/utils'
hljs.registerLanguage('json', json)

const themeStore = useThemeStore()

// function kebabCase(str) {
//   return str.replace(/([A-Z])/g, '-$1').toLowerCase()
// }

// function setupCssVar() {
//   const common = naiveThemeOverrides.common
//   for (const key in common) {
//     useCssVar(`--${kebabCase(key)}`, document.documentElement).value = common[key] || ''
//     if (key === 'primarysColor')
//       window.localStorage.setItem('__THEME_COLOR__', common[key] || '')
//   }
// }

// 挂载 naive 组件的方法至 window, 方便全局使用
function setupNaiveTools() {
  window.$loadingBar = useLoadingBar()
  window.$notification = useNotification()
  window.$message = setupMessage(useMessage())
  window.$dialog = setupDialog(useDialog())
}

const NaiveProviderContent = defineComponent({
  setup() {
    // setupCssVar()
    setupNaiveTools()
  },
  render() {
    return h('div')
  },
})

// 监听黑暗模式
watch(() => themeStore.darkMode, (v) => {
  v
    ? document.documentElement.classList.add('dark')
    : document.documentElement.classList.remove('dark')
})

// 默认是否加载黑夜模式
onMounted(() => {
  themeStore.darkMode
    ? document.documentElement.classList.add('dark')
    : document.documentElement.classList.remove('dark')
})
</script>

<template>
  <n-config-provider
    wh-full
    :theme-overrides="naiveThemeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :hljs="hljs"
  >
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <NaiveProviderContent />
            <slot />
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>
