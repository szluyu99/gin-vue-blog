<script setup>
import { renderCustomIcon, renderIcon } from '@/utils'

const router = useRouter()
const route = useRoute()

function handleBreadClick(path) {
  if (path !== route.path)
    router.push(path) // 根据 router.js 中设置的 redirect 进行跳转
}

function getIcon(meta) {
  if (meta?.customIcon)
    return renderCustomIcon(meta.customIcon, { size: 18 })
  if (meta?.icon)
    return renderIcon(meta.icon, { size: 18 })
  return null
}
</script>

<template>
  <n-breadcrumb>
    <n-breadcrumb-item
      v-for="item in route.matched.filter((item) => !!item.meta?.title)"
      :key="item.path"
      @click="handleBreadClick(item.path)"
    >
      <component :is="getIcon(item.meta)" />
      {{ item.meta?.title }}
    </n-breadcrumb-item>
  </n-breadcrumb>
</template>
