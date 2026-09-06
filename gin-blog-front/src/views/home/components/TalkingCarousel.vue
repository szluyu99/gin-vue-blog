<script setup>
import { onMounted, onUnmounted, ref } from 'vue'

import api from '@/api'
import { getOneSentence, getRandomSentence } from '@/utils'

// 首页轮播: 有说说就轮播真实内容, 没有(或接口挂了)退回一言兜底
const sentence = ref(getRandomSentence())
const index = ref(0)
const talks = ref([])
let timer = null

async function rotate() {
  if (talks.value.length) {
    sentence.value = talks.value[index.value % talks.value.length].content
    index.value++
  }
  else {
    sentence.value = await getOneSentence()
  }
}

onMounted(async () => {
  try {
    const resp = await api.getTalks({ page_num: 1, page_size: 5 })
    talks.value = resp.data?.page_data ?? []
  }
  catch (err) {
    console.error(err)
  }

  await rotate()
  timer = setInterval(rotate, 5000)
})

onUnmounted(() => clearInterval(timer))
</script>

<template>
  <div class="card-view card-enter">
    <div class="flex text-center">
      <button class="i-mdi-chat-outline text-xl" />
      <div class="flex-1">
        {{ sentence }}
      </div>
      <RouterLink to="/talks" class="animate-arrow i-mdi-chevron-double-right text-2xl" />
    </div>
  </div>
</template>

<style scoped>
.animate-arrow {
  animation: 1s passing infinite;
}

/* 左 -> 右 闪的特效 */
@keyframes passing {
  0% { transform: translateX(-50%); opacity: 0; }
  50% { transform: translateX(0); opacity: 1; }
  100% { transform: translateX(50%); opacity: 0; }
}
</style>
