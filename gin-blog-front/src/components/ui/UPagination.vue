<script setup>
// 前台的分页器
//
// 后台用的是 naive-ui 的 NPagination, 但前台没装 naive-ui —— 之前这里直接写了
// <NPagination>, 组件解析不到, 分页器整个渲染不出来。前台就这一处用得上, 自己写一个。
import { computed } from 'vue'

const props = defineProps({
  page: { type: Number, required: true },
  pageCount: { type: Number, required: true },
  // 页码多时最多同时显示几个数字, 两侧用省略号
  window: { type: Number, default: 5 },
})

const emit = defineEmits(['update:page'])

// 当前页附近的一段页码, 保证长度是 window(总页数不足时就全显示)
const pages = computed(() => {
  const total = props.pageCount
  const size = Math.min(props.window, total)
  let start = props.page - Math.floor(size / 2)
  start = Math.max(1, Math.min(start, total - size + 1))
  return Array.from({ length: size }, (_, i) => start + i)
})

function go(page) {
  if (page < 1 || page > props.pageCount || page === props.page) {
    return
  }
  emit('update:page', page)
}
</script>

<template>
  <nav class="flex select-none items-center gap-1" aria-label="分页">
    <button
      class="btn" :disabled="page <= 1"
      aria-label="上一页" @click="go(page - 1)"
    >
      <span class="i-mdi:chevron-left block text-lg" />
    </button>

    <button v-if="pages[0] > 1" class="btn" @click="go(1)">
      1
    </button>
    <span v-if="pages[0] > 2" class="px-1 color-muted">...</span>

    <button
      v-for="p of pages" :key="p"
      class="btn" :class="p === page ? 'bg-#49b1f5 text-white' : ''"
      :aria-current="p === page ? 'page' : undefined"
      @click="go(p)"
    >
      {{ p }}
    </button>

    <span v-if="pages[pages.length - 1] < pageCount - 1" class="px-1 color-muted">...</span>
    <button v-if="pages[pages.length - 1] < pageCount" class="btn" @click="go(pageCount)">
      {{ pageCount }}
    </button>

    <button
      class="btn" :disabled="page >= pageCount"
      aria-label="下一页" @click="go(page + 1)"
    >
      <span class="i-mdi:chevron-right block text-lg" />
    </button>
  </nav>
</template>

<style scoped>
.btn {
  --uno: min-w-8 h-8 px-2 f-c-c rounded bg-surface text-sm shadow-sm transition-300;
  --uno: hover:bg-#49b1f5 hover:text-white disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-surface disabled:hover:text-inherit;
}
</style>
