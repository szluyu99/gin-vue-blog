<script setup>
import api from '@/api'
import { getToken } from '@/utils'
import { useAppStore } from '@/store'
const appStore = useAppStore()

onMounted(() => {
  if (getToken()) {
    appStore.getBlogInfo() // 获取博客全局设置
    api.report() // 上报用户信息
  }
})
</script>

<template>
  <AppProvider>
    <router-view v-slot="{ Component }">
      <!-- 动态组件 -->
      <component :is="Component" />
    </router-view>
  </AppProvider>
</template>
