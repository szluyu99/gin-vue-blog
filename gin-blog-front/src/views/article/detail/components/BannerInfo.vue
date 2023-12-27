<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import dayjs from 'dayjs'

const { article } = defineProps({
  article: {},
})

const wordNum = ref(0) // 字数统计
const readTime = ref('') // 阅读时间

onMounted(() => {
  // 统计字数
  wordNum.value = deleteHTMLTag(article.content).length
  // 计算阅读时间
  readTime.value = `${Math.round(wordNum.value / 400)}分钟`
})

// 删除 HTML 标签
function deleteHTMLTag(str) {
  return str
    .replace(/<\/?[^>]*>/g, '')
    .replace(/[|]*\n/, '')
    .replace(/&npsp;/gi, '')
}

function formatDate(date) {
  return dayjs(date).format('YYYY-MM-DD')
}
</script>

<template>
  <!-- PC 端显示 -->
  <div class="mx-4 mt-12 hidden text-center text-light lg:block">
    <h1 class="text-3xl">
      {{ article.title }}
    </h1>
    <p class="f-c-c py-2">
      <span class="i-mdi:calendar mr-1 text-lg" /> 发布于 {{ formatDate(article.created_at) }}
      <span class="px-2">|</span>
      <span class="i-mdi:update mr-1 text-lg" /> 更新于 {{ formatDate(article.updated_at) }}
      <span class="px-2">|</span>
      <RouterLink :to="`/categories/${article.category?.id}?name=${article.category?.name}`" class="f-c-c">
        <span class="i-material-symbols:menu mr-1 text-lg" /> {{ article.category?.name }}
      </RouterLink>
    </p>
    <div class="f-c-c">
      <span class="i-ic:twotone-text-snippet mr-1 text-lg" /> 字数统计：{{ wordNum }}
      <span class="px-2">|</span>
      <span class="i-mdi:timelapse mr-1 text-lg" /> 阅读时长：{{ readTime }}
      <span class="px-2">|</span>
      <span class="i-mdi:eye mr-1 text-lg" /> 阅读量：{{ article.view_count }}
      <span class="px-2">|</span>
      <span class="i-ic:outline-insert-comment mr-1 text-lg" /> 评论数：{{ article.comment_count }}
    </div>
  </div>
  <!-- 移动端显示 -->
  <div class="mx-4 mt-12 block text-left text-light lg:hidden space-y-1.5">
    <h1 class="text-2xl">
      {{ article.title }}
    </h1>
    <div class="mb-1 mt-2 flex flex-wrap items-center lg:justify-center">
      <span class="i-mdi:calendar mr-1" /> 发布于 {{ formatDate(article.created_at) }}
      <span class="px-2">|</span>
      <span class="i-mdi:update mr-1" /> 更新于 {{ formatDate(article.updated_at) }}
    </div>
    <div class="flex gap-2">
      <div class="f-c-c">
        <RouterLink :to="`/categories/${article.category?.id}?name=${article.category?.name}`">
          <span class="i-material-symbols:menu mr-1" /> {{ article.category?.name }}
        </RouterLink>
      </div>
      <RouterLink v-for="tag of article.tags" :key="tag.id" :to="`/tags/${tag.id}?name=${tag.name}`">
        <span class="border-1px border-blue rounded-xl px-2 py-1 text-sm text-white"> {{ tag.name }} </span>
      </RouterLink>
    </div>
    <div>
      <span class="i-ic:twotone-text-snippet mr-1" /> 字数统计：{{ wordNum }}
      <span class="px-2">|</span>
      <span class="i-mdi:timelapse mr-1" /> 阅读时长：{{ readTime }}
    </div>
    <div>
      <span class="i-mdi:eye mr-1" /> 阅读量：{{ article.view_count }}
      <span class="px-2">|</span>
      <span class="i-ic:outline-insert-comment mr-1" /> 评论数：{{ article.comment_count }}
    </div>
  </div>
</template>
