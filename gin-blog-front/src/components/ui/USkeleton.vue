<script setup>
/**
 * 骨架屏: 列表型页面加载时的占位
 *
 * 原来的加载态是「半透明空卡片 + 一个转圈」, 内容回来时啪一下出现, 落差很大。
 * 骨架屏把版式先占住, 内容替换进来时位置基本不动。
 */
defineProps({
  // 骨架行数
  rows: { type: Number, default: 5 },
  // 是否显示顶部那条更宽的标题行
  title: { type: Boolean, default: true },
})

// 每行宽度错开一点, 免得看起来像表格
const WIDTHS = ['92%', '78%', '85%', '70%', '88%', '75%']
</script>

<template>
  <div class="animate-pulse space-y-4" aria-busy="true" aria-live="polite">
    <span class="sr-only">加载中</span>
    <div v-if="title" class="h-7 w-2/5 rounded bg-surface-soft" />
    <div
      v-for="i of rows" :key="i"
      class="h-5 rounded bg-surface-soft"
      :style="{ width: WIDTHS[(i - 1) % WIDTHS.length] }"
    />
  </div>
</template>
