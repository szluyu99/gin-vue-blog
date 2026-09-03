<script setup>
import { NLayout, NLayoutSider } from 'naive-ui'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import themes from '@/assets/themes'
import { useTagStore, useThemeStore } from '@/store'
import AppHeader from './header/index.vue'

import Sidebar from './sidebar/index.vue'
import AppTags from './tags/index.vue'

const themeStore = useThemeStore()
const tagStore = useTagStore()
const router = useRouter()

// 缓存的路由名
// router.getRoutes() 不是响应式的, 单靠它这个 computed 会永久缓存首次结果,
// 登录后动态添加的路由进不了 KeepAlive, 必须刷新页面才生效。
// 这里挂上 tagStore.tags 做依赖: 访问过的页面一定先进标签栏, 再需要被缓存
const keepAliveRouteNames = computed(() => {
  void tagStore.tags.length
  return router.getRoutes()
    .filter(route => route.meta?.keepAlive)
    .map(route => route.name)
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
