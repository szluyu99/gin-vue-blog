<script setup lang="ts">
import { convertImgUrl, formatDate } from '@/utils'

interface Props { idx?: number; article: any }
const { idx = 0, article } = defineProps<Props>()

// 判断图片放置位置 (左 or 右)
const isRightClass = computed(() => idx % 2 === 0
  ? 'rounded-t-3xl md:order-0 md:rounded-l-3xl md:rounded-tr-0'
  : 'rounded-t-3xl md:order-1 md:rounded-r-3xl md:rounded-tl-0')
</script>

<template>
  <div
    items-center mt-20
    w-full bg-white rounded-2rem shadow-md
    transition-600 hover:shadow-2xl
    animate-zoom-in animate-duration-700
    flex-col md:flex-row
    h-430 md:h-280
  >
    <!-- 封面图 -->
    <div
      :class="isRightClass" overflow-hidden
      h-230 w-full md:w="45/100" md:h-full
    >
      <router-link :to="`/article/${article.id}`">
        <img wh-full transition-600 hover:scale-110 :src="convertImgUrl(article.img)">
      </router-link>
    </div>
    <!-- 文章信息 -->
    <div
      mt-20 mb-10 w="9/10"
      md:w="55/100" md:px-45
    >
      <router-link :to="`/article/${article.id}`" text-8xl hover:text-violet>
        《{{ article.title }}》
      </router-link>
      <div flex flex-wrap pt-15 pb-12 color="#858585" text-14>
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
        <router-link :to="`/categories/${article.category_id}?name=${article.category.name}`" flex items-center>
          <i-mdi-inbox-full mr-3 /> {{ article.category.name }}
        </router-link>
        <span mx-7>|</span>
        <!-- 标签 -->
        <router-link
          v-for="tag in article.tags" :key="tag.id"
          :to="`/tags/${tag.id}?name=${tag.name}`"
          flex items-center mx-2
        >
          <i-mdi-tag-multiple mr-2 /> {{ tag.name }}
        </router-link>
      </div>
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
