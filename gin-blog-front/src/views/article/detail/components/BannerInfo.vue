<script setup lang="ts">
import { formatDate } from '@/utils'
interface Props { article: any }
const { article } = defineProps<Props>()

const wordNum = ref(0) // 字数统计
const readTime = ref('') // 阅读时间

onMounted(() => {
  // 统计字数
  wordNum.value = deleteHTMLTag(article.content).length
  // 计算阅读时间
  readTime.value = `${Math.round(wordNum.value / 400)}分钟`
})

// 删除 HTML 标签
function deleteHTMLTag(str: string) {
  return str
    .replace(/<\/?[^>]*>/g, '')
    .replace(/[|]*\n/, '')
    .replace(/&npsp;/gi, '')
}
</script>

<template>
  <div text-center text-light mt-50 text-15>
    <h1 text-34>
      {{ article.title }}
    </h1>
    <p flex items-center justify-center py-8>
      <i-mdi:calendar text-18 mr-3 /> 发布于 {{ formatDate(article.created_at) }}
      <span px-6>|</span>
      <i-mdi:update text-18 mr-3 /> 更新于 {{ formatDate(article.updated_at) }}
      <span px-6>|</span>
      <router-link :to="`/categories/${article.category.id}`" flex items-center>
        <i-material-symbols:menu text-18 mr-3 /> {{ article.category.name }}
      </router-link>
    </p>
    <p flex items-center justify-center>
      <i-ic:twotone-text-snippet text-18 mr-3 /> 字数统计：{{ wordNum }}
      <span px-6>|</span>
      <i-mdi:timelapse text-18 mr-3 /> 阅读时长：{{ readTime }}
      <span px-6>|</span>
      <i-mdi:eye text-18 mr-3 /> 阅读量：{{ article.view_count }}
      <span px-6>|</span>
      <i-ic:outline-insert-comment text-18 mr-3 /> 评论数：{{ 0 }}
    </p>
  </div>
</template>
