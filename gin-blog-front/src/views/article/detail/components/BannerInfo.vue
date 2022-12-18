<script setup lang="ts">
import { formatDate } from '@/utils'

const { article } = defineProps<{ article: any }>()

let wordNum = $ref(0) // 字数统计
let readTime = $ref('') // 阅读时间

onMounted(() => {
  // 统计字数
  wordNum = deleteHTMLTag(article.content).length
  // 计算阅读时间
  readTime = `${Math.round(wordNum / 400)}分钟`
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
  <!-- PC 端显示 -->
  <div
    text-center text-15 text-light mt-50 mx-15
    hidden lg:block
  >
    <h1 text-24 lg:text-34>
      {{ article.title }}
    </h1>
    <p py-8 f-c-c>
      <i-mdi:calendar text-18 mr-3 /> 发布于 {{ formatDate(article.created_at) }}
      <span px-6>|</span>
      <i-mdi:update text-18 mr-3 /> 更新于 {{ formatDate(article.updated_at) }}
      <span px-6>|</span>
      <router-link :to="`/categories/${article.category.id}?name=${article.category.name}`" f-c-c>
        <i-material-symbols:menu text-18 mr-3 /> {{ article.category.name }}
      </router-link>
    </p>
    <div f-c-c>
      <i-ic:twotone-text-snippet text-18 mr-3 /> 字数统计：{{ wordNum }}
      <span px-6>|</span>
      <i-mdi:timelapse text-18 mr-3 /> 阅读时长：{{ readTime }}
      <span px-6>|</span>
      <i-mdi:eye text-18 mr-3 /> 阅读量：{{ article.view_count }}
      <span px-6>|</span>
      <i-ic:outline-insert-comment text-18 mr-3 /> 评论数：{{ article.comment_count }}
    </div>
  </div>
  <!-- 移动端显示 -->
  <div
    text-left text-light mt-50 mx-15
    block lg:hidden
  >
    <h1 text-24>
      {{ article.title }}
    </h1>
    <div mt-10 mb-5 flex items-center flex-wrap lg:justify-center>
      <i-mdi:calendar text-18 mr-3 /> 发布于 {{ formatDate(article.created_at) }}
      <span px-6>|</span>
      <i-mdi:update text-18 mr-3 /> 更新于 {{ formatDate(article.updated_at) }}
    </div>
    <div>
      <router-link :to="`/categories/${article.category.id}?name=${article.category.name}`" mr-8>
        <i-material-symbols:menu text-18 mr-3 /> {{ article.category.name }}
      </router-link>

      <router-link
        v-for="tag of article.tags" :key="tag.id"
        :to="`/tags/${tag.id}?name=${tag.name}`"
      >
        <span
          px-8 py-1 mr-8
          border-1px rounded-3rem border-blue text-white
        > {{ tag.name }} </span>
      </router-link>
    </div>
    <div my-5>
      <i-ic:twotone-text-snippet text-18 mr-3 /> 字数统计：{{ wordNum }}
      <span px-6>|</span>
      <i-mdi:timelapse text-18 mr-3 /> 阅读时长：{{ readTime }}
    </div>
    <div>
      <i-mdi:eye text-18 mr-3 /> 阅读量：{{ article.view_count }}
      <span px-6>|</span>
      <i-ic:outline-insert-comment text-18 mr-3 /> 评论数：{{ article.comment_count }}
    </div>
  </div>
</template>
