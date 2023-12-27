<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'

import BannerPage from '@/components/BannerPage.vue'
import { convertImgUrl } from '@/utils'
import api from '@/api'

const route = useRoute()

const loading = ref(true)
const articleList = ref([])
const name = ref(route.query.name) // 标题上显示的 标签/分类 名称

onMounted(async () => {
  const resp = await api.getArticles({
    category_id: route.params.categoryId,
    tag_id: route.params.tagId,
  })
  articleList.value = resp.data
  loading.value = false
})
</script>

<template>
  <BannerPage :loading="loading" :title="`${route.meta?.title} - ${name}`" label="article_list">
    <div class="grid grid-cols-12 gap-4">
      <div v-for="article of articleList" :key="article.id" class="col-span-12 lg:col-span-4 md:col-span-6">
        <!-- 卡片 -->
        <div class="animate-zoom-in animate-duration-650 rounded-xl bg-white pb-2 shadow-md transition-300 hover:shadow-2xl">
          <!-- 图片 -->
          <div class="overflow-hidden">
            <RouterLink :to="`/article/${article.id}`">
              <img :src="convertImgUrl(article.img)" class="h-[220px] w-full rounded-t-xl transition-600 hover:scale-110">
            </RouterLink>
          </div>
          <!-- 内容 -->
          <div>
            <div class="space-y-1.5">
              <!-- 标题 -->
              <RouterLink :to="`/article/${article.id}`">
                <p class="inline-block px-3 pt-2 hover:color-violet">
                  {{ article.title }}
                </p>
              </RouterLink>
              <div class="flex justify-between px-3">
                <!-- 发布日期 -->
                <span class="flex items-center">
                  <span class="i-mdi:clock-outline mr-1" />
                  <span> {{ dayjs(article.created_at).format('YYYY-MM-DD') }} </span>
                </span>
                <!-- 分类 -->
                <RouterLink :to="`/categories/${article.category_id}?name=${article.category?.name}`">
                  <div class="flex items-center text-#4c4948 transition-300 hover:color-violet">
                    <span class="i-ic:outline-bookmark mr-1" />
                    <span> {{ article.category?.name }} </span>
                  </div>
                </RouterLink>
              </div>
            </div>
            <!-- divider -->
            <div class="my-2 h-0.5 bg-gray-200" />
            <!-- 标签 -->
            <div class="px-3 space-x-1.5">
              <RouterLink v-for="tag of article.tags" :key="tag.id" :to="`/tags/${tag.id}?name=${tag.name}`">
                <span class="inline-block cursor-pointer rounded-xl from-green-400 to-blue-500 bg-gradient-to-r px-2 py-1 text-xs text-white transition-500 hover:scale-110 hover:from-pink-500 hover:to-yellow-500">
                  {{ tag.name }}
                </span>
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </BannerPage>
</template>
