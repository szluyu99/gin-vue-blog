<script setup>
import { ref } from 'vue'
import { storeToRefs } from 'pinia'

import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'

import { useAppStore } from '@/store'

dayjs.extend(duration)
const { blogConfig, viewCount } = storeToRefs(useAppStore())

// 每秒刷新时间
const runTime = ref('')

// 每秒刷新当前时间
setInterval(() => {
  const createTime = dayjs(blogConfig.value.website_createtime)
  runTime.value = dayjs.duration(dayjs().diff(createTime)).format('D 天 H 时 m 分 s 秒')
}, 1000)
</script>

<template>
  <div class="hidden animate-zoom-in lg:block space-y-10 card-view">
    <p class="flex items-center">
      <span class="animate-bang i-icon-park:analysis mr-5 text-18" />
      <span class="text-16"> 网站咨询 </span>
    </p>
    <p class="text-14">
      <span> 运行时间： </span>
      <span class="float-right"> {{ runTime }} </span>
    </p>
    <p class="text-14">
      <span> 总访问量： </span>
      <span class="float-right"> {{ viewCount }} </span>
    </p>
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
