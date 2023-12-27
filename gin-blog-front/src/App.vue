<script setup>
import { onMounted, ref } from 'vue'

import UToast from '@/components/ui/UToast.vue'
import AppHeader from '@/components/layout/AppHeader.vue'
import GlobalModal from '@/components/modal/index.vue'
import BackToTop from '@/components/BackTop.vue'

import { useAppStore, useUserStore } from '@/store'

const appStore = useAppStore()
const userStore = useUserStore()

const messageRef = ref(null)
const notifyRef = ref(null)

onMounted(() => {
  appStore.getPageList()
  appStore.getBlogInfo()
  userStore.getUserInfo()

  // 挂载全局提示
  window.$message = messageRef.value
  window.$notify = notifyRef.value
})

// 禁止右键菜单
// document.addEventListener('contextmenu', e => e.preventDefault())
</script>

<template>
  <!-- 顶部中间的消息提示 -->
  <UToast ref="messageRef" position="top" align="center" :timeout="3000" closeable />
  <!-- 右上方的消息通知 -->
  <UToast ref="notifyRef" position="top" align="right" :timeout="3000" closeable />

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
