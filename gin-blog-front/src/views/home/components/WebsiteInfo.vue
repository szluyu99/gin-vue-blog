<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { storeToRefs } from 'pinia'

import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'

import { useAppStore } from '@/store'

dayjs.extend(duration)
const { blogConfig, viewCount } = storeToRefs(useAppStore())

// 每秒刷新时间
const runTime = ref('')

// 每秒刷新当前时间
const timer = setInterval(() => {
  const createTime = dayjs(blogConfig.value.website_createtime)
  runTime.value = dayjs.duration(dayjs().diff(createTime)).format('D 天 H 时 m 分')
}, 30 * 1000)

onMounted(() => {
  const createTime = dayjs(blogConfig.value.website_createtime)
  runTime.value = dayjs.duration(dayjs().diff(createTime)).format('D 天 H 时 m 分')
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <div class="card-view hidden animate-zoom-in animate-duration-600 lg:block space-y-2">
    <p class="flex items-center text-lg">
      <span class="i-icon-park:analysis mr-1.5" />
      <span> 网站咨询 </span>
    </p>
    <div class="space-y-1">
      <p>
        <span> 运行时间： </span>
        <span class="float-right"> {{ runTime }} </span>
      </p>
      <p>
        <span> 总访问量： </span>
        <span class="float-right"> {{ viewCount }} </span>
      </p>
    </div>
  </div>
</template>
