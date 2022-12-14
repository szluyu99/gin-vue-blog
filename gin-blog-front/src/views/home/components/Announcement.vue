<script setup lang="ts">
import { useAppStore } from '@/store'
const { blogInfo } = storeToRefs(useAppStore())
</script>

<template>
  <div card-view animate-zoom-in hidden lg:block>
    <p mb-10 flex items-center>
      <i-fluent-emoji-flat:bell mr-6 text-20 class="animate-bang" />
      <span text-16> 公告</span>
    </p>
    <div text-6xl leading-25>
      {{ blogInfo.blog_config?.website_notice }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
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
