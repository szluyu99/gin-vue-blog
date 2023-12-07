<script setup>
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'
const { blogInfo } = storeToRefs(useAppStore())
</script>

<template>
  <div class="card-view hidden animate-zoom-in lg:block">
    <p class="mb-2 flex items-center text-lg">
      <span class="i-fluent-emoji-flat:bell mr-2" />
      <span> 公告 </span>
    </p>
    <div class="leading-6">
      {{ blogInfo.blog_config?.website_notice }}
    </div>
  </div>
</template>

<style scoped>
.animate-bang {
  animation: bang 0.8s linear infinite;
}
</style>
