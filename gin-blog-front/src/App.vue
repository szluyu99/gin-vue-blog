<script setup>
import { onMounted } from 'vue'

import AppHeader from '@/components/layout/AppHeader.vue'
import BackToTop from '@/components/other/BackTop.vue'
import GlobalModal from '@/components/modal/index.vue'

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
  <div class="h-full w-full flex flex-col">
    <!-- 顶部导航栏 -->
    <AppHeader />
    <!-- 中间内容(包含底部信息) -->
    <article class="flex flex-1 flex-col">
      <RouterView v-slot="{ Component, route }">
        <component :is="Component" :key="route.path" />
      </RouterView>
    </article>
  </div>
  <!-- 回到顶部 -->
  <BackToTop />
  <!-- 全局弹窗 -->
  <GlobalModal />
</template>
