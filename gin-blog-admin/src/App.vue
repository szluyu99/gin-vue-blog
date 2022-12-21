<script setup>
import api from '@/api'
import { getToken } from '@/utils'

const token = getToken()

onMounted(() => {
  token && api.report() // token 不为空, 上报用户信息
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
