<script setup>
import { useThemeStore } from '@/store'
const themeStore = useThemeStore()
</script>

<template>
  <n-icon mr-20 cursor-pointer size="18" @click="themeStore.switchDarkMode">
    <icon-mdi-moon-waning-crescent v-if="themeStore.darkMode" />
    <icon-mdi-white-balance-sunny v-else />
  </n-icon>
</template>
