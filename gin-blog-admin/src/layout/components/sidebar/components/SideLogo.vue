<script setup>
import { useAppStore } from '@/store'
const { collapsed } = storeToRefs(useAppStore())
const title = import.meta.env.VITE_TITLE
</script>

<template>
  <router-link f-c-c h-60 to="/">
    <icon-custom-logo text-36 color-primary />
    <h2
      v-show="!collapsed"
      flex-shrink-0 max-w-140 ml-10 text-16 color-primary font-bold
    >
      {{ title }}
    </h2>
  </router-link>
</template>
