<script setup>
import { NIcon } from 'naive-ui'
import { useFullscreen } from '@vueuse/core'

const { isFullscreen, toggle } = useFullscreen()
</script>

<template>
  <NIcon size="18" class="mr-20 cursor-pointer" @click="toggle">
    <span v-if="isFullscreen" class="i-ant-design:fullscreen-exit-outlined" />
    <span v-else class="i-ant-design:fullscreen-outlined" />
  </NIcon>
</template>
