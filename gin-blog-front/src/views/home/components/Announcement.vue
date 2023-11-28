<script setup>
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'
const { blogInfo } = storeToRefs(useAppStore())
</script>

<template>
  <div class="hidden animate-zoom-in lg:block card-view">
    <p class="mb-10 flex items-center">
      <span class="i-fluent-emoji-flat:bell animate-bang mr-6 text-20" />
      <span class="text-16"> 公告</span>
    </p>
    <div class="text-6xl leading-25">
      {{ blogInfo.blog_config?.website_notice }}
    </div>
  </div>
</template>

<style scoped>
.animate-bang {
  animation: bang 0.8s linear infinite;
}

@keyframes bang {
  0%,
  100% {
    transform: scale(1.1);
  }
  50% {
    transform: scale(0.9);
  }
}
</style>
