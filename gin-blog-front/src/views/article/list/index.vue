<script setup>
import dayjs from 'dayjs'
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import api from '@/api'
import BannerPage from '@/components/BannerPage.vue'
import UPagination from '@/components/ui/UPagination.vue'
import { convertImgUrl } from '@/utils'

const route = useRoute()

const loading = ref(true)
const articleList = ref([])
const name = ref(route.query.name) // 标题上显示的 标签/分类 名称

// 分页: 以前只取第一页, 分类/标签下超过一页的文章没有任何入口能看到
const PAGE_SIZE = 9
const current = ref(1)
const total = ref(0)
const pageCount = computed(() => Math.ceil(total.value / PAGE_SIZE))

watch(current, () => {
  getArticles()
  window.scrollTo({ behavior: 'smooth', top: 0 })
})

async function getArticles() {
  loading.value = true
  // 失败时 loading 也要复位, 否则页面卡在加载态; data 为空也不能让列表变成 null
  try {
    const resp = await api.getArticles({
      category_id: route.params.categoryId,
      tag_id: route.params.tagId,
      page_num: current.value,
      page_size: PAGE_SIZE,
    })
    articleList.value = resp.data?.page_data ?? []
    total.value = resp.data?.total ?? 0
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  getArticles()
})
</script>

<template>
  <BannerPage :loading="loading" :title="`${route.meta?.title} - ${name}`" label="article_list">
    <div class="grid grid-cols-12 gap-4">
      <div v-for="article of articleList" :key="article.id" class="col-span-12 lg:col-span-4 md:col-span-6">
        <!-- 卡片 -->
        <div class="animate-zoom-in animate-duration-650 rounded-xl bg-surface pb-2 shadow-md transition-300 hover:shadow-2xl">
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
                  <div class="flex items-center text-main transition-300 hover:color-violet">
                    <span class="i-ic:outline-bookmark mr-1" />
                    <span> {{ article.category?.name }} </span>
                  </div>
                </RouterLink>
              </div>
            </div>
            <!-- divider -->
            <div class="my-2 h-0.5 bg-surface-soft" />
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
    <!-- 分页 -->
    <div v-if="pageCount > 1" class="mt-8 flex justify-center">
      <UPagination v-model:page="current" :page-count="pageCount" />
    </div>
    <!-- 空列表提示: 以前没有数据时页面是一片空白 -->
    <div v-if="!loading && !articleList.length" class="py-20 text-center color-muted">
      这里还没有文章
    </div>
  </BannerPage>
</template>
