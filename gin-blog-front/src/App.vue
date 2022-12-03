<script setup lang="ts">
import BackToTop from '@/components/other/BackTop.vue'
import GlobalModel from '@/components/modal/index.vue'

import { useAppStore, useUserStore } from '@/store'
const [appStore, userStore] = [useAppStore(), useUserStore()]

appStore.getBlogInfo()
userStore.getUserInfo() // FIXME: 解决 token 后台过期过快问题

// 禁止右键菜单
document.addEventListener('contextmenu', e => e.preventDefault())
</script>

<template>
  <AppProvider>
    <div wh-full flex-col>
      <!-- 顶部导航栏 -->
      <AppHeader />
      <!-- 中间内容(包含底部信息) -->
      <article flex-1 flex-col>
        <router-view v-slot="{ Component, route }">
          <component :is="Component" :key="route.path" />
        </router-view>
      </article>
    </div>
    <!-- 回到顶部 -->
    <BackToTop />
    <!-- 全局弹窗 -->
    <GlobalModel />
  </AppProvider>
</template>
