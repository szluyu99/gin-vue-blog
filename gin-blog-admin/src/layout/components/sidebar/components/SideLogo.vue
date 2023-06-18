<script setup>
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
const title = import.meta.env.VITE_TITLE
</script>

<template>
  <router-link class="h-60 flex items-center justify-center" to="/">
    <div class="h-full w-36 flex items-center justify-center">
      <img src="@/assets/svg/logo.svg" alt="logo" class="h-36 w-36">
    </div>
    <h2
      v-show="!themeStore.collapsed"
      class="ml-10 max-w-140 flex-shrink-0 text-16 font-bold text-primary"
    >
      {{ title }}
    </h2>
  </router-link>
</template>
