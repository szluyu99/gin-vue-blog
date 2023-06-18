<script setup>
import { NConfigProvider, darkTheme, dateZhCN, zhCN } from 'naive-ui'

import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'
import { naiveThemeOverrides } from '@/../setting/theme.json'
import { useThemeStore } from '@/store'

hljs.registerLanguage('json', json)
const themeStore = useThemeStore()
</script>

<template>
  <NConfigProvider
    class="h-full w-full"
    :theme="themeStore.darkMode ? darkTheme : undefined"
    :theme-overrides="naiveThemeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :hljs="hljs"
  >
    <slot />
  </NConfigProvider>
</template>
