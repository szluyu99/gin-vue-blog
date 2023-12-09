<script setup>
import { ref } from 'vue'
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { Icon } from '@iconify/vue'

const { y } = useWindowScroll()
const styleVal = ref('')
watchThrottled(y, () => {
  styleVal.value = (y.value > 20) ? 'opacity: 1; transform: translateX(-40px);' : ''
}, { throttle: 100 })

const options = [
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
  <div class="fixed bottom-20 z-4 text-white transition-600 -right-9 space-y-1" :style="styleVal">
    <div
      v-for="item of options" :key="item.icon"
      class="f-c-c cursor-pointer rounded-sm bg-#49b1f5 p-1 duration-300 hover:bg-amber"
    >
      <Icon class="h-5 w-5" :icon="item.icon" @click="item.fn" />
    </div>
  </div>
</template>
