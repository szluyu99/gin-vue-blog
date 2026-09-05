<script setup>
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useAppStore } from '@/store'

const appStore = useAppStore()

const { y } = useWindowScroll()
const styleVal = ref('')
watchThrottled(y, () => {
  styleVal.value = (y.value > 20) ? 'opacity: 1; transform: translateX(-40px);' : ''
}, { throttle: 100 })

const options = computed(() => [
  {
    // 图标名直接写成 UnoCSS 类名, 构建时就生成 CSS, 不再运行时去 iconify CDN 拉
    // (bi / fluent / ph 三个集合只为几个图标而引入不值得, 换成已加载的 mdi 里的近似图标)
    icon: appStore.isDark ? 'i-mdi:weather-sunny' : 'i-mdi:weather-night',
    fn: () => appStore.toggleTheme(),
  },
  {
    icon: 'i-uiw:setting',
    fn: () => window.$message?.info('设置开发中...'),
  },
  {
    icon: 'i-mdi:arrow-up-bold',
    fn: () => window.scrollTo({ behavior: 'smooth', top: 0 }),
  },
])
</script>

<template>
  <div class="fixed bottom-20 z-4 text-white transition-600 -right-9 space-y-1" :style="styleVal">
    <div
      v-for="item of options" :key="item.icon"
      class="f-c-c cursor-pointer rounded-sm bg-#49b1f5 p-1 duration-300 hover:bg-amber"
    >
      <span class="block h-5 w-5" :class="item.icon" @click="item.fn" />
    </div>
  </div>
</template>
