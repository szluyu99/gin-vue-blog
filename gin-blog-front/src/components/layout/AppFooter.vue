<script setup>
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'
const { blogInfo } = storeToRefs(useAppStore())
</script>

<template>
  <footer class="flex items-center justify-center">
    <div class="footer-wrap w-full p-x-20 p-y-40 text-center text-14 leading-30 text-white">
      <p> &copy;{{ 2022 }} - {{ new Date().getFullYear() }} By 阵、雨</p>
      <a href="https://beian.miit.gov.cn/" target="_blank">
        {{ blogInfo.blog_config?.website_record }}
      </a>
    </div>
  </footer>
</template>

<style scoped>
.footer-wrap {
  background: linear-gradient(-45deg, #ee7752, #ce3e75, #23a6d5, #23d5ab);
  background-size: 400% 400%;
  animation: Gradient 8s ease infinite;
}

@keyframes Gradient {
  0% {
    background-position: 0 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0 50%;
  }
}
</style>
