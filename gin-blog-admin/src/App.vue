<script setup>
import { darkTheme, dateZhCN, NConfigProvider, zhCN } from 'naive-ui'

import themes from '@/assets/themes'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()

// highlight.js 只有操作日志页的 NCode 用得到, 放在这里会进入首屏 chunk (约 38KB),
// 改为在 views/log/operation 里按需引入并直接传给 NCode。

// 上报用户信息, 需要时取消注释, 并补回 onMounted / api / useAuthStore 的导入
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
  >
    <RouterView v-slot="{ Component }">
      <component :is="Component" />
    </RouterView>
  </NConfigProvider>
</template>
