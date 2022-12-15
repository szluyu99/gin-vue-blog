<script setup>
// 父组件传来的方法
const emit = defineEmits(['search', 'reset'])
</script>

<template>
  <div
    flex items-start justify-between
    min-h-60 p-15
    b-1 bc-ccc rounded-8 bg="#fafafc"
  >
    <n-space wrap :size="[35, 15]">
      <slot />
    </n-space>
    <div flex-shrink-0>
      <n-button ml-20 secondary type="primary" @click="emit('reset')">
        重置
      </n-button>
      <n-button ml-20 type="primary" @click="emit('search')">
        搜索
      </n-button>
    </div>
  </div>
</template>
