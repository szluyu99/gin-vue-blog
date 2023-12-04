<script setup>
import { ref } from 'vue'
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { Icon } from '@iconify/vue'

const { y } = useWindowScroll()
const styleVal = ref('')
watchThrottled(y, () => {
  styleVal.value = (y.value > 20) ? 'opacity: 1; transform: translateX(-40px);' : ''
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
  <div class="fixed bottom-85 z-4 text-20 text-white transition-600 -right-35" :style="styleVal">
    <div
      v-for="item of operations" :key="item.icon"
      class="my-3 h-30 w-30 f-c-c cursor-pointer rounded bg-#49b1f5 duration-300 hover:bg-amber"
    >
      <Icon class="h-20 w-20" :icon="item.icon" @click="item.fn" />
    </div>
  </div>
</template>
