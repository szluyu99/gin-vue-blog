<script setup>
import { dateZhCN, useDialog, useLoadingBar, useMessage, useNotification, zhCN } from 'naive-ui'
// import { kebabCase } from 'lodash-es'

// 代码高亮插件
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'
import { naiveThemeOverrides } from '~/setting'
import { setupDialog, setupMessage } from '@/utils'
hljs.registerLanguage('json', json)

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
