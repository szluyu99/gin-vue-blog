<script setup>
import { storeToRefs } from 'pinia'
import { NIcon, NWatermark } from 'naive-ui'
import { useThemeStore, useUserStore } from '@/store'

const { nickname } = storeToRefs(useUserStore())
const themeStore = useThemeStore()
</script>

<template>
  <NIcon size="18" class="mr-20 cursor-pointer" @click="themeStore.switchWatermark">
    <div v-if="themeStore.watermarkFlag" class="i-ic:baseline-branding-watermark" />
    <div v-else class="i-ic:outline-branding-watermark" />
  </NIcon>

  <NWatermark
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
