<script setup>
import { onMounted } from 'vue'
import { NConfigProvider, darkTheme, dateZhCN, zhCN } from 'naive-ui'
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'

import { useAuthStore, useThemeStore } from '@/store'
import themes from '@/assets/themes'
import api from '@/api'

hljs.registerLanguage('json', json)
const themeStore = useThemeStore()

// onMounted(() => {
//   const { accessToken } = useAuthStore()
//   accessToken && api.report() // 上报用户信息
// })

// FIXME: 每次 Docker 打包完运行会继承之前的 localStorage
// TODO: 每次如果发现当前没有路由信息，就跳转到登录页
</script>

<template>
  <NConfigProvider
    class="h-full w-full"
    :theme="themeStore.darkMode ? darkTheme : undefined"
    :theme-overrides="themes.naiveThemeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :hljs="hljs"
  >
    <RouterView v-slot="{ Component }">
      <component :is="Component" />
    </RouterView>
  </NConfigProvider>
</template>
