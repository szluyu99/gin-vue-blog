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
            <slot></slot>
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup>
import { defineComponent, h } from 'vue'
import { useLoadingBar, useDialog, useMessage, useNotification, zhCN, dateZhCN } from 'naive-ui'
import { setupDialog, setupMessage } from '@/utils'
import { useCssVar } from '@vueuse/core'
import { kebabCase } from 'lodash-es'
import { naiveThemeOverrides } from '~/setting'

// 代码高亮插件
import hljs from 'highlight.js/lib/core'
import javascript from 'highlight.js/lib/languages/javascript'
hljs.registerLanguage('javascript', javascript)

// 加载主题 -> setting/theme.json
function setupCssVar() {
  const common = naiveThemeOverrides.common
  for (const key in common) {
    useCssVar(`--${kebabCase(key)}`, document.documentElement).value = common[key] || ''
    if (key === 'primarysColor') window.localStorage.setItem('__THEME_COLOR__', common[key] || '')
  }
}

// 挂载 naive 组件的方法至 window, 方便全局使用
function setupNaiveTools() {
  window.$loadingBar = useLoadingBar()
  window.$notification = useNotification()
  window.$message = setupMessage(useMessage())
  window.$dialog = setupDialog(useDialog())
}

const NaiveProviderContent = defineComponent({
  setup() {
    setupCssVar()
    setupNaiveTools()
  },
  render() {
    return h('div')
  },
})
</script>
