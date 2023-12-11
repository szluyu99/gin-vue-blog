<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { NLayout, NLayoutSider } from 'naive-ui'

import AppHeader from './header/index.vue'
import Sidebar from './sidebar/index.vue'
import AppTags from './tags/index.vue'

import { useTagStore, useThemeStore } from '@/store'
import themes from '@/assets/themes'

const themeStore = useThemeStore()
const tagStore = useTagStore()
const router = useRouter()

// 缓存的路由名
const keepAliveRouteNames = computed(() => {
  const allRoutes = router.getRoutes()
  const names = allRoutes.filter(route => route.meta?.keepAlive).map(route => route.name)
  return names
})
</script>

<template>
  <NLayout has-sider class="h-full w-full">
    <!-- 左侧边栏 -->
    <NLayoutSider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      :native-scrollbar="false"
      :collapsed="themeStore.collapsed"
    >
      <Sidebar />
    </NLayoutSider>
    <!-- 右半部分 -->
    <article class="flex flex-1 flex-col overflow-hidden">
      <!-- 头部 -->
      <header
        class="flex items-center border-b-1 border-gray-200 border-b-solid px-4"
        :style="{ height: `${themes.header.height}px` }"
      >
        <AppHeader />
      </header>
      <!-- 标签栏 -->
      <section v-if="themes.tags.visible" class="border-b border-gray-200 border-b-solid">
        <AppTags :style="{ height: `${themes.tags.height}px` }" />
      </section>
      <!-- 主体内容 -->
      <section class="flex-1 overflow-hidden">
        <RouterView v-slot="{ Component, route }">
          <KeepAlive :include="keepAliveRouteNames">
            <component
              :is="Component"
              v-if="tagStore.reloading"
              :key="tagStore.aliveKeys[route.name] || route.fullPath"
            />
          </keepalive>
        </RouterView>
      </section>
    </article>
  </NLayout>
</template>
