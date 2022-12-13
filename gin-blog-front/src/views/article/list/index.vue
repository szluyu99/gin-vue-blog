<script setup lang="ts">
import { convertImgUrl, formatDate } from '@/utils'
import api from '@/api'

const route = useRoute()

let loading = $ref(true)
let articleList = $ref<any>([])

const name = $ref(route.query.name) // 标题上显示的 标签/分类 名称

onMounted(async () => {
  const res = await api.getArticles({
    category_id: route.params.categoryId,
    tag_id: route.params.tagId,
  })
  articleList = res.data
  loading = false
})
</script>

<template>
  <BannerPage
    :loading="loading"
    :title="`${route.meta?.title} - ${name}`"
    label="article_list"
  >
    <n-grid x-gap="15" y-gap="15" cols="1 m:3" responsive="screen">
      <n-gi v-for="article of articleList" :key="article.id">
        <!-- 卡片 -->
        <div
          shadow-md rounded-2rem bg-white
          animate-zoom-in animate-duration-700
          transition-300 hover:shadow-2xl
        >
          <!-- 图片 -->
          <div overflow-hidden>
            <router-link :to="`/article/${article.id}`">
              <img
                :src="convertImgUrl(article.img)"
                h-220 w-full rounded-t-1rem
                transition-600 hover:scale-110
              >
            </router-link>
          </div>
          <!-- 内容 -->
          <div>
            <!-- 标题 -->
            <router-link :to="`/article/${article.id}`">
              <p px-15 pt-12 pb-6 text-17 hover:color-violet>
                {{ article.title }}
              </p>
            </router-link>
            <p px-15 py-3 flex justify-between>
              <!-- 发布日期 -->
              <span flex items-center>
                <i-mdi:clock-outline text-20 mr-3 />
                <span text-16> {{ formatDate(article.created_at) }} </span>
              </span>
              <!-- 分类 -->
              <router-link :to="`/categories/${article.category_id}?name=${article.category.name}`">
                <div flex items-center text="#4c4948" hover:color-violet>
                  <i-ic:outline-bookmark text-20 mr-3 />
                  <span text-16> {{ article.category.name }} </span>
                </div>
              </router-link>
            </p>
            <div border-1px h-1 my-5 rounded-1rem />
            <!-- 标签 -->
            <p px-15 pt-6 pb-10>
              <router-link v-for="tag of article.tags" :key="tag.id" :to="`/tags/${tag.id}?name=${tag.name}`">
                <span
                  inline-block px-12 py-3 mr-3rem rounded-3rem text-white text-12 cursor-pointer
                  bg-gradient-to-r from-green-400 to-blue-500
                  transition-500 hover:scale-120 hover:from-pink-500 hover:to-yellow-500
                >
                  {{ tag.name }}
                </span>
              </router-link>
            </p>
          </div>
        </div>
      </n-gi>
    </n-grid>
  </BannerPage>
</template>
