<script setup lang="ts">
import { convertImgUrl, formatDate } from '@/utils'

interface Props {
  idx?: number
  article: any // TODO: 类型
}

const { idx = 0, article } = defineProps<Props>()

// 判断图片放置位置 (左 or 右)
const isRightClass = computed(() => idx % 2 === 0 ? 'order-0 rounded-l-3xl' : 'order-1 rounded-r-3xl')
</script>

<template>
  <div
    flex items-center mt-20
    w-full h-280 bg-white rounded-2rem shadow-md
    transition-600 hover:shadow-2xl
    animate-zoom-in animate-duration-700
  >
    <!-- 封面图 -->
    <div :class="isRightClass" overflow-hidden h-full w="45/100">
      <router-link :to="`/article/${article.id}`">
        <img wh-full transition-600 hover:scale-110 :src="convertImgUrl(article.img)">
      </router-link>
    </div>
    <!-- 文章信息 -->
    <div px-45 w="55/100">
      <router-link :to="`/article/${article.id}`" text-8xl cursor-pointer hover:text-violet>
        《{{ article.title }}》
      </router-link>
      <div flex pt-15 pb-12 color="#858585" text-14>
        <!-- 置顶 -->
        <span v-if="article.is_top === 1" flex items-center color="#ff7242">
          <i-carbon:align-vertical-top mr-3 /> 置顶
        </span>
        <span v-if="article.is_top === 1" mx-7>|</span>
        <!-- 日期 -->
        <span flex items-center>
          <i-mdi-calendar-month-outline mr-3 /> {{ formatDate(article.created_at) }}
        </span>
        <span mx-7>|</span>
        <!-- 分类 -->
        <router-link :to="`/categories/${article.category_id}`" flex items-center>
          <i-mdi-inbox-full mr-3 /> {{ article.category.name }}
        </router-link>
        <span mx-7>|</span>
        <!-- 标签 -->
        <router-link
          v-for="tag in article.tags" :key="tag.id"
          :to="`/tags/${tag.id}`"
          flex items-center mx-2
        >
          <i-mdi-tag-multiple mr-2 /> {{ tag.name }}
        </router-link>
      </div>
      <!-- TODO: 后端? 过滤 Markdown 符号 -->
      <div leading-25 class="ell-4">
        {{ article.content }}
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
// 省略文字最多 4 行
.ell-4 {
  display: -webkit-box;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
}
</style>
