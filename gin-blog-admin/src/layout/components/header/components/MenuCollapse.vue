<script setup>
import { NIcon } from 'naive-ui'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<template>
  <NIcon size="20" class="cursor-pointer" @click="themeStore.switchCollapsed">
    <span v-if="themeStore.collapsed" class="i-mdi:format-indent-increase" />
    <span v-else class="i-mdi:format-indent-decrease" />
  </NIcon>
</template>
