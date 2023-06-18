<script setup>
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/store'

const { reloadFlag, aliveKeys } = storeToRefs(useAppStore())
const router = useRouter()

// 当前路由中所有的路由
const allRoutes = router.getRoutes()
// 缓存的路由名
const keepAliveRouteNames = computed(() => {
  return allRoutes
    .filter(route => route.meta?.keepAlive)
    .map(route => route.name)
})
</script>

<template>
  <router-view v-slot="{ Component, route }">
    <!-- KeepAlive 是一个内置组件, 它的功能是在多个组件动态切换时缓存被移除的组件实例 -->
    <!-- 设置路由的 meta.keepAlive 为 true 会缓存该页面内容 -->
    <KeepAlive :include="keepAliveRouteNames">
      <component
        :is="Component"
        v-if="reloadFlag"
        :key="aliveKeys[route.name] || route.fullPath"
      />
    </KeepAlive>
  </router-view>
</template>
