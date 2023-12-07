<script setup>
import { computed } from 'vue'
import dayjs from 'dayjs'

import { convertImgUrl } from '@/utils'

const { idx, article } = defineProps({
  idx: Number,
  article: {},
})

// 判断图片放置位置 (左 or 右)
const isRightClass = computed(() => idx % 2 === 0
  ? 'rounded-t-xl md:order-0 md:rounded-l-xl md:rounded-tr-0'
  : 'rounded-t-xl md:order-1 md:rounded-r-xl md:rounded-tl-0')
</script>

<template>
  <div class="h-[430px] w-full flex flex-col animate-zoom-in animate-duration-700 items-center rounded-xl bg-white shadow-md transition-600 md:h-[280px] md:flex-row hover:shadow-2xl">
    <!-- 封面图 -->
    <div :class="isRightClass" class="h-[230px] w-full overflow-hidden md:h-full md:w-45/100">
      <RouterLink :to="`/article/${article.id}`">
        <img class="h-full w-full transition-600 hover:scale-110" :src="convertImgUrl(article.img)">
      </RouterLink>
    </div>
    <!-- 文章信息 -->
    <div class="mb-3 mt-4 w-9/10 md:w-55/100 space-y-4 md:px-10">
      <RouterLink :to="`/article/${article.id}`" class="text-2xl hover:text-violet">
        《{{ article.title }}》
      </RouterLink>
      <div class="flex flex-wrap text-sm color-#858585">
        <!-- 置顶 -->
        <span v-if="article.is_top === 1" class="flex items-center color-#ff7242">
          <div class="i-carbon:align-vertical-top mr-1" /> 置顶
        </span>
        <span v-if="article.is_top === 1" class="mx-2">|</span>
        <!-- 日期 -->
        <span class="flex items-center">
          <div class="i-mdi-calendar-month-outline mr-1" /> {{ dayjs(article.created_at).format('YYYY-MM-DD') }}
        </span>
        <span class="mx-2">|</span>
        <!-- 分类 -->
        <RouterLink :to="`/categories/${article.category_id}?name=${article.category.name}`" class="flex items-center">
          <div class="i-mdi-inbox-full mr-1" /> {{ article.category.name }}
        </RouterLink>
        <span class="mx-2">|</span>
        <!-- 标签 -->
        <RouterLink
          v-for="tag in article.tags" :key="tag.id"
          :to="`/tags/${tag.id}?name=${tag.name}`"
          class="mr-1 flex items-center"
        >
          <div class="i-mdi-tag-multiple mr-1" /> {{ tag.name }}
        </RouterLink>
      </div>
      <div class="ell-4 text-sm leading-7">
        {{ article.content }}
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 省略文字最多 4 行 */
.ell-4 {
  display: -webkit-box;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
}
</style>
