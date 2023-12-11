<script setup>
import { NIcon } from 'naive-ui'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<template>
  <NIcon size="18" class="cursor-pointer" @click="themeStore.switchCollapsed">
    <i v-if="themeStore.collapsed" class="i-mdi:format-indent-increase" />
    <i v-else class="i-mdi:format-indent-decrease" />
  </NIcon>
</template>
