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
  <div class="card-view hidden animate-zoom-in lg:block space-y-2">
    <p class="flex items-center text-lg">
      <span class="animate-bang i-icon-park:analysis mr-2" />
      <span> 网站咨询 </span>
    </p>
    <p class="">
      <span> 运行时间： </span>
      <span class="float-right"> {{ runTime }} </span>
    </p>
    <p class="">
      <span> 总访问量： </span>
      <span class="float-right"> {{ viewCount }} </span>
    </p>
  </div>
</template>

<style scoped>
.animate-bang {
  animation: bang 0.8s linear infinite;
}
</style>
