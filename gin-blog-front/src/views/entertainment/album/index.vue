<script setup>
import BannerPage from '@/components/BannerPage.vue'
</script>

<template>
  <BannerPage label="album" card>
    <div class="text-center">
      <span class="text-8">
        禁止访问
      </span>
      <div class="flex justify-center">
        <img class="w-50" src="/images/404.svg" alt="404">
      </div>
      <button @click="$router.push('/')">
        回到首页
      </button>
    </div>
  </BannerPage>
</template>
