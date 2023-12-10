<script setup>
import { useRoute, useRouter } from 'vue-router'
import { NBreadcrumb, NBreadcrumbItem } from 'naive-ui'
import { renderIcon } from '@/utils'

const router = useRouter()
const route = useRoute()

function handleBreadClick(path) {
  if (path !== route.path) {
    router.push(path)
  }
}

function getIcon(meta) {
  if (meta?.icon) {
    return renderIcon(meta.icon, { size: 18 })
  }
  return null
}
</script>

<template>
  <NBreadcrumb>
    <NBreadcrumbItem
      v-for="item in route.matched.filter((item) => !!item.meta?.title)"
      :key="item.path"
      @click="handleBreadClick(item.path)"
    >
      <component :is="getIcon(item.meta)" />
      {{ item.meta?.title }}
    </NBreadcrumbItem>
  </NBreadcrumb>
</template>
