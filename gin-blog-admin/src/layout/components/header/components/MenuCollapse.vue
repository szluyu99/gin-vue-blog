<script setup>
import { useThemeStore } from '@/store'
const themeStore = useThemeStore()
</script>

<template>
  <n-icon size="20" cursor-pointer @click="themeStore.switchCollapsed">
    <icon-mdi:format-indent-increase v-if="themeStore.collapsed" />
    <icon-mdi:format-indent-decrease v-else />
  </n-icon>
</template>
