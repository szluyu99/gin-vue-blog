<script setup>
import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'

import { storeToRefs } from 'pinia'
import { onMounted, onUnmounted, ref, watch } from 'vue'

import { useAppStore } from '@/store'

dayjs.extend(duration)
const { blogConfig, viewCount } = storeToRefs(useAppStore())

const runTime = ref('')

// 注意不能用 duration.format('D 天'): D 是"天"这个分量, 建站超过一个月就会归零重算
// (1340 天会显示成 "1 天"), 总天数要用 asDays
function refreshRunTime() {
  const createTime = dayjs(blogConfig.value.website_createtime)
  if (!createTime.isValid()) { // 配置还没拉回来
    runTime.value = '-'
    return
  }
  const diff = dayjs.duration(dayjs().diff(createTime))
  runTime.value = `${Math.floor(diff.asDays())} 天 ${diff.hours()} 时 ${diff.minutes()} 分`
}

// 每 30 秒刷新一次
let timer = null

onMounted(() => {
  refreshRunTime()
  timer = setInterval(refreshRunTime, 30 * 1000)
})

// 建站时间来自接口, 拿到之后立刻重算, 不用等下一个 30 秒
watch(() => blogConfig.value.website_createtime, refreshRunTime)

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <div class="card-view card-enter hidden lg:block space-y-2">
    <p class="flex items-center text-lg">
      <span class="i-icon-park:analysis mr-1.5" />
      <span> 网站资讯 </span>
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
