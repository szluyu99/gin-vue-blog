<script setup lang="ts">
import { convertImgUrl, formatDate } from '@/utils'
import api from '@/api'

const route = useRoute()

const articleList = ref<any>([])
const name = ref('') // 标题上显示的 标签/分类 名称
onMounted(async () => {
  await api.getArticles({
    category_id: route.params.categoryId,
    tag_id: route.params.tagId,
  }).then((res) => {
    articleList.value = res.data
    // TODO: 优化, 更好的做法?
    if (route.meta?.title === '标签')
      name.value = res.data[0].tags.find((e: any) => e.id === +route.params.tagId).name
    else name.value = res.data[0].category.name
  })
})
</script>

<template>
  <BannerPage
    :title="`${route.meta?.title} - ${name}`"
    banner-img="https://fastly.jsdelivr.net/gh/jerryc127/CDN@latest/Photo/categories.jpg"
  >
    <n-grid x-gap="15" y-gap="15" :cols="3">
      <n-gi v-for="item of articleList" :key="item.id">
        <!-- 卡片 -->
        <div
          shadow-md rounded-2rem bg-white
          animate-zoom-in animate-duration-700
          transition-300 hover:shadow-2xl
        >
          <!-- 图片 -->
          <div overflow-hidden>
            <router-link :to="`/article/${item.id}`">
              <img
                :src="convertImgUrl(item.img)"
                h-220 w-full rounded-t-1rem
                transition-600 hover:scale-110
              >
            </router-link>
          </div>
          <!-- 内容 -->
          <div>
            <!-- 标题 -->
            <router-link :to="`/article/${item.id}`">
              <p px-15 pt-12 pb-6 text-17 hover:color-violet>
                {{ item.title }}
              </p>
            </router-link>
            <p px-15 py-3 flex justify-between>
              <!-- 发布日期 -->
              <span flex items-center>
                <i-mdi:clock-outline text-20 mr-3 />
                <span text-16> {{ formatDate(item.created_at) }} </span>
              </span>
              <!-- 分类 -->
              <router-link :to="`/categories/${item.category_id}`">
                <div flex items-center text="#4c4948" hover:color-violet>
                  <i-ic:outline-bookmark text-20 mr-3 />
                  <span text-16> {{ item.category.name }} </span>
                </div>
              </router-link>
            </p>
            <div border-1px h-1 my-5 rounded-1rem />
            <!-- 标签 -->
            <p px-15 pt-6 pb-10>
              <router-link v-for="tag of item.tags" :key="tag.id" :to="`/tags/${tag.id}`">
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
