<script setup>
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'

const { blogInfo } = storeToRefs(useAppStore())
</script>

<template>
  <div class="card-view hidden animate-zoom-in animate-duration-600 lg:block space-y-2">
    <p class="flex items-center text-lg">
      <span class="i-fluent-emoji-flat:bell mr-1.5" />
      <span> 公告 </span>
    </p>
    <div class="leading-6">
      {{ blogInfo.blog_config?.website_notice }}
    </div>
  </div>
</template>
