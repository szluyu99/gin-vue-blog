<script setup>
import { ref } from 'vue'
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import TheIcon from '@/components/icon/TheIcon.vue'

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
  <div
    class="fixed z-4 -right-35 bottom-85 text-20 text-white transition-600"
    :style="styleVal"
  >
    <div
      v-for="item of operations" :key="item.icon"
      class="w-30 h-30 my-3 f-c-c cursor-pointer rounded bg-#49b1f5 hover:bg-amber"
    >
      <TheIcon :icon="item.icon" :size="20" @click="item.fn" />
    </div>
  </div>
</template>
