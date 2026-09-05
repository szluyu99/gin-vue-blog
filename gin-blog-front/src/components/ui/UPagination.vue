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

// 样式全部写在模板的 class 上, 不要挪进 <style scoped> 里用 --uno 拼:
// scoped 样式注入在 uno.css 之后, 同样是 0,1,0 特异度, 会把 :class 上条件加的
// bg-primary 盖掉 —— 之前选中页就是这样变成白底白字的
const BASE = 'h-9 min-w-9 f-c-c rounded-lg px-3 text-sm shadow-sm transition-300'
const IDLE = 'bg-surface text-main hover:bg-primary hover:text-white hover:shadow-md'

// 当前页附近的一段页码, 保证长度是 window(总页数不足时就全显示)
const pages = computed(() => {
  const total = props.pageCount
  const size = Math.min(props.window, total)
  let start = props.page - Math.floor(size / 2)
  start = Math.max(1, Math.min(start, total - size + 1))
  return Array.from({ length: size }, (_, i) => start + i)
})

function itemClass(active) {
  return active ? `${BASE} bg-primary text-white font-bold shadow-md` : `${BASE} ${IDLE}`
}

// 禁用态不给 hover: 靠 disabled: 变体和 hover 抢特异度, 结果取决于生成顺序
function arrowClass(disabled) {
  return disabled ? `${BASE} bg-surface text-muted cursor-not-allowed opacity-50` : `${BASE} ${IDLE}`
}

function go(page) {
  if (page < 1 || page > props.pageCount || page === props.page) {
    return
  }
  emit('update:page', page)
}
</script>

<template>
  <nav class="flex select-none items-center gap-2" aria-label="分页">
    <button
      :class="arrowClass(page <= 1)" :disabled="page <= 1"
      aria-label="上一页" @click="go(page - 1)"
    >
      <span class="i-mdi:chevron-left block text-lg" />
    </button>

    <template v-if="pages[0] > 1">
      <button :class="itemClass(false)" @click="go(1)">
        1
      </button>
      <span v-if="pages[0] > 2" class="px-1 text-muted">...</span>
    </template>

    <button
      v-for="p of pages" :key="p"
      :class="itemClass(p === page)"
      :aria-current="p === page ? 'page' : undefined"
      @click="go(p)"
    >
      {{ p }}
    </button>

    <template v-if="pages[pages.length - 1] < pageCount">
      <span v-if="pages[pages.length - 1] < pageCount - 1" class="px-1 text-muted">...</span>
      <button :class="itemClass(false)" @click="go(pageCount)">
        {{ pageCount }}
      </button>
    </template>

    <button
      :class="arrowClass(page >= pageCount)" :disabled="page >= pageCount"
      aria-label="下一页" @click="go(page + 1)"
    >
      <span class="i-mdi:chevron-right block text-lg" />
    </button>
  </nav>
</template>
