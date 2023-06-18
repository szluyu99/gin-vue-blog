<script setup>
import { onMounted } from 'vue'

import AppProvider from './components/common/AppProvider.vue'
import api from '@/api'
import { getToken } from '@/utils'

onMounted(() => {
  const token = getToken()
  token && api.report() // token 不为空, 上报用户信息
})
</script>

<template>
  <AppProvider>
    <router-view v-slot="{ Component }">
      <component :is="Component" />
    </router-view>
  </AppProvider>
</template>
