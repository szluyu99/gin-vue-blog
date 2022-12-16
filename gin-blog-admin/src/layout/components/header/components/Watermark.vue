<script setup>
import { storeToRefs } from 'pinia'
import { useThemeStore, useUserStore } from '@/store'

const { nickname } = storeToRefs(useUserStore())
const themeStore = useThemeStore()
</script>

<template>
  <n-icon size="18" mr-20 cursor-pointer @click="themeStore.switchWatermark">
    <icon-ic:baseline-branding-watermark v-if="themeStore.watermarkFlag" />
    <icon-ic:outline-branding-watermark v-else />
  </n-icon>

  <n-watermark
    v-if="themeStore.watermarkFlag"
    :content="nickname"
    cross
    fullscreen
    :font-size="16"
    :line-height="16"
    :width="184"
    :height="184"
    :x-offset="12"
    :y-offset="60"
    :rotate="-15"
  />
</template>
