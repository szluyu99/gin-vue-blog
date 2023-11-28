<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { formatDate } from '@/utils'

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
</script>

<template>
  <!-- PC 端显示 -->
  <div class="mx-15 mt-50 hidden text-center text-15 text-light lg:block">
    <h1 class="text-24 lg:text-34">
      {{ article.title }}
    </h1>
    <p class="f-c-c py-8">
      <span class="i-mdi:calendar mr-3 text-18" /> 发布于 {{ formatDate(article.created_at) }}
      <span class="px-6">|</span>
      <span class="i-mdi:update mr-3 text-18" /> 更新于 {{ formatDate(article.updated_at) }}
      <span class="px-6">|</span>
      <RouterLink :to="`/categories/${article.category.id}?name=${article.category.name}`" class="f-c-c">
        <span class="i-material-symbols:menu mr-3 text-18" /> {{ article.category.name }}
      </RouterLink>
    </p>
    <div class="f-c-c">
      <span class="i-ic:twotone-text-snippet mr-3 text-18" /> 字数统计：{{ wordNum }}
      <span class="px-6">|</span>
      <span class="i-mdi:timelapse mr-3 text-18" /> 阅读时长：{{ readTime }}
      <span class="px-6">|</span>
      <span class="i-mdi:eye mr-3 text-18" /> 阅读量：{{ article.view_count }}
      <span class="px-6">|</span>
      <span class="i-ic:outline-insert-comment mr-3 text-18" /> 评论数：{{ article.comment_count }}
    </div>
  </div>
  <!-- 移动端显示 -->
  <div class="mx-15 mt-50 block text-left text-light lg:hidden">
    <h1 class="text-24">
      {{ article.title }}
    </h1>
    <div class="mb-5 mt-10 flex flex-wrap items-center lg:justify-center">
      <span class="i-mdi:calendar mr-3 text-18" /> 发布于 {{ formatDate(article.created_at) }}
      <span class="px-6">|</span>
      <span class="i-mdi:update mr-3 text-18" /> 更新于 {{ formatDate(article.updated_at) }}
    </div>
    <div>
      <RouterLink :to="`/categories/${article.category.id}?name=${article.category.name}`" class="mr-8">
        <span class="i-material-symbols:menu mr-3 text-18" /> {{ article.category.name }}
      </RouterLink>
      <RouterLink v-for="tag of article.tags" :key="tag.id" :to="`/tags/${tag.id}?name=${tag.name}`">
        <span class="mr-8 border-1px border-blue rounded-3rem px-8 py-1 text-white"> {{ tag.name }} </span>
      </RouterLink>
    </div>
    <div class="my-5">
      <span class="i-ic:twotone-text-snippet mr-3 text-18" /> 字数统计：{{ wordNum }}
      <span class="px-6">|</span>
      <span class="i-mdi:timelapse mr-3 text-18" /> 阅读时长：{{ readTime }}
    </div>
    <div>
      <span class="i-mdi:eye mr-3 text-18" /> 阅读量：{{ article.view_count }}
      <span px-6>|</span>
      <span class="i-ic:outline-insert-comment mr-3 text-18" /> 评论数：{{ article.comment_count }}
    </div>
  </div>
</template>
