<script setup>
import { NIcon, NWatermark } from 'naive-ui'
import { useThemeStore, useUserStore } from '@/store'

const userStore = useUserStore()
const themeStore = useThemeStore()
</script>

<template>
  <NIcon size="18" class="cursor-pointer" @click="themeStore.switchWatermark">
    <i v-if="themeStore.watermarked" class="i-ic:baseline-branding-watermark" />
    <i v-else class="i-ic:outline-branding-watermark" />
  </NIcon>

  <NWatermark
    v-if="themeStore.watermarked"
    :content="userStore.nickname"
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
