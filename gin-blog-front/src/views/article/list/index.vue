<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { NGi, NGrid } from 'naive-ui'

import { convertImgUrl, formatDate } from '@/utils'
import api from '@/api'

const route = useRoute()

const loading = ref(true)
const articleList = ref([])
const name = ref(route.query.name) // 标题上显示的 标签/分类 名称

onMounted(async () => {
  const res = await api.getArticles({
    category_id: route.params.categoryId,
    tag_id: route.params.tagId,
  })
  articleList.value = res.data
  loading.value = false
})
</script>

<template>
  <BannerPage :loading="loading" :title="`${route.meta?.title} - ${name}`" label="article_list">
    <NGrid x-gap="15" y-gap="15" cols="1 m:3" responsive="screen">
      <NGi v-for="article of articleList" :key="article.id">
        <!-- 卡片 -->
        <div class="animate-zoom-in animate-duration-700 rounded-2rem bg-white shadow-md transition-300 hover:shadow-2xl">
          <!-- 图片 -->
          <div class="overflow-hidden">
            <RouterLink :to="`/article/${article.id}`">
              <img :src="convertImgUrl(article.img)" class="h-220 w-full rounded-t-1rem transition-600 hover:scale-110">
            </RouterLink>
          </div>
          <!-- 内容 -->
          <div>
            <!-- 标题 -->
            <RouterLink :to="`/article/${article.id}`">
              <p class="px-15 pb-6 pt-12 text-17 hover:color-violet">
                {{ article.title }}
              </p>
            </RouterLink>
            <p class="flex justify-between px-15 py-3">
              <!-- 发布日期 -->
              <span class="flex items-center">
                <span class="i-mdi:clock-outline mr-3 text-20" />
                <span class="text-16"> {{ formatDate(article.created_at) }} </span>
              </span>
              <!-- 分类 -->
              <RouterLink :to="`/categories/${article.category_id}?name=${article.category.name}`">
                <div class="flex items-center text-#4c4948 hover:color-violet">
                  <span class="i-ic:outline-bookmark mr-3 text-20" />
                  <span class="text-16"> {{ article.category.name }} </span>
                </div>
              </RouterLink>
            </p>
            <div class="my-5 h-1 border-1px rounded-1rem" />
            <!-- 标签 -->
            <p class="px-15 pb-10 pt-6">
              <RouterLink v-for="tag of article.tags" :key="tag.id" :to="`/tags/${tag.id}?name=${tag.name}`">
                <span class="mr-3rem inline-block cursor-pointer rounded-3rem from-green-400 to-blue-500 bg-gradient-to-r px-12 py-3 text-12 text-white transition-500 hover:scale-120 hover:from-pink-500 hover:to-yellow-500">
                  {{ tag.name }}
                </span>
              </RouterLink>
            </p>
          </div>
        </div>
      </NGi>
    </NGrid>
  </BannerPage>
</template>
