<script setup>
import { NIcon } from 'naive-ui'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<template>
  <NIcon class="cursor-pointer" size="18" @click="themeStore.switchDarkMode">
    <span v-if="themeStore.darkMode" class="i-mdi-moon-waning-crescent" />
    <span v-else class="i-mdi-white-balance-sunny" />
  </NIcon>
</template>
