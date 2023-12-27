<script setup>
import { onMounted, ref } from 'vue'
import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'

const loading = ref(true)
const tagList = ref([])

onMounted(() => {
  api.getTags().then((resp) => {
    tagList.value = resp.data || []
    loading.value = false
  })
})

// 随机字体大小
function randomFontSize() {
  return Math.floor(Math.random() * 16) + 15
}
// 随机颜色 16进制
function randomColorHex() {
  return `#${Math.floor(Math.random() * 16777215).toString(16)}`
}
</script>

<template>
  <BannerPage :loading="loading" title="标签" label="tag" card>
    <h2 class="text-center text-2xl leading-8 lg:text-3xl">
      标签 - {{ tagList.length }}
    </h2>
    <div class="mt-6 text-center">
      <RouterLink
        v-for="t of tagList" :key="t.id" :to="`tags/${t.id}?name=${t.name}`"
        :style="{
          'font-size': `${randomFontSize()}px`,
          'color': `${randomColorHex()}`,
        }"
        class="inline-block px-2 leading-11 transition-300 hover:scale-130 !hover:text-lightblue"
      >
        {{ t.name }}
      </RouterLink>
    </div>
  </BannerPage>
</template>

<style scoped>
/* 实现截断文字效果, 即不会在结束处将一个词语拆开 */
a {
  display: inline-block;
}
</style>
