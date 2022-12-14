<script setup lang="ts">
import { useWindowScroll, watchThrottled } from '@vueuse/core'

const { y } = $(useWindowScroll())

let styleVal = $ref('')
watchThrottled($$(y), () => {
  styleVal = (y > 20) ? 'opacity: 1; transform: translateX(-40px);' : ''
}, { throttle: 100 })

const operations = [
  {
    icon: 'bi:moon-stars-fill',
    fn: () => window.$message?.info('黑夜模式开发中...'),
  },
  {
    icon: 'uiw:setting',
    fn: () => window.$message?.info('设置开发中...'),
  },
  {
    icon: 'fluent:arrow-up-12-filled',
    fn: () => window.scrollTo({ behavior: 'smooth', top: 0 }),
  },
]
</script>

<template>
  <div
    fixed z-4 right="-35" bottom-85
    text-20 text-white transition-600
    :style="styleVal"
  >
    <div
      v-for="item of operations" :key="item.icon"
      w-30 h-30 my-3 f-c-c cursor-pointer rounded
      bg="#49b1f5" hover:bg-amber
    >
      <the-icon :icon="item.icon" :size="20" @click="item.fn" />
    </div>
  </div>
</template>
