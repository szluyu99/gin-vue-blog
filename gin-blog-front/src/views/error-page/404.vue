<script setup>
import BannerPage from '@/components/page/BannerPage.vue'
</script>

<template>
  <BannerPage label="404" card>
    <div class="text-center">
      <span class="text-32">
        404 资源不存在
      </span>
      <div class="flex justify-center">
        <img class="w-200" src="../../../public/images/404.svg" alt="404">
      </div>
      <button @click="$router.push('/')">
        回到首页
      </button>
    </div>
  </BannerPage>
</template>
