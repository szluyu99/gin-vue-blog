<script setup lang="ts">
import { onMounted } from 'vue'

import AppProvider from './components/layout/AppProvider.vue'
import AppHeader from '@/components/layout/AppHeader.vue'
import BackToTop from '@/components/other/BackTop.vue'
import GlobalModel from '@/components/modal/index.vue'

import { useAppStore, useUserStore } from '@/store'
const appStore = useAppStore()
const userStore = useUserStore()

onMounted(() => {
  appStore.getBlogInfo()
  userStore.getUserInfo()
})

// 禁止右键菜单
// document.addEventListener('contextmenu', e => e.preventDefault())
</script>

<template>
  <AppProvider>
    <div class="h-full w-full flex-col">
      <!-- 顶部导航栏 -->
      <AppHeader />
      <!-- 中间内容(包含底部信息) -->
      <article class="flex-col flex-1">
        <RouterView v-slot="{ Component, route }">
          <component :is="Component" :key="route.path" />
        </RouterView>
      </article>
    </div>
    <!-- 回到顶部 -->
    <BackToTop />
    <!-- 全局弹窗 -->
    <GlobalModel />
  </AppProvider>
</template>
