<script setup>
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
const title = import.meta.env.VITE_TITLE
</script>

<template>
  <RouterLink class="h-[60px] flex items-center justify-center" to="/">
    <div class="h-full flex items-center justify-center">
      <img src="/image/logo.svg" alt="logo" class="h-[36px] w-[36px]">
    </div>
    <h2
      v-show="!themeStore.collapsed"
      class="ml-2 max-w-[150px] flex-shrink-0 text-base text-primary font-bold dark:text-white"
    >
      {{ title }}
    </h2>
  </RouterLink>
</template>
