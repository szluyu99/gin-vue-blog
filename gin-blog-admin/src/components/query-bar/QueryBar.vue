<script setup>
import { NButton, NSpace } from 'naive-ui'

const emit = defineEmits(['search', 'reset'])
</script>

<template>
  <div class="min-h-60 flex items-start justify-between border border-gray-200 border-gray-400 rounded-8 border-solid bg-gray-50 p-15 dark:bg-black">
    <NSpace wrap :size="[35, 15]">
      <slot />
    </NSpace>
    <div class="flex-shrink-0 space-x-20">
      <NButton secondary type="primary" @click="emit('reset')">
        重置
      </NButton>
      <NButton type="primary" @click="emit('search')">
        搜索
      </NButton>
    </div>
  </div>
</template>
