<script setup lang="ts">
// 无限加载组件
import InfiniteLoading from 'v3-infinite-loading'

import ArticleCard from './components/ArticleCard.vue'
import AuthorInfo from './components/AuthorInfo.vue'
import WebsiteInfo from './components/WebsiteInfo.vue'
import HomeBanner from './components/HomeBanner.vue'
import Announcement from './components/Announcement.vue'
import TalkingCarousel from './components/TalkingCarousel.vue'

import api from '@/api'

let articleList = $ref<any>([])
let loading = $ref(false)

// 无限加载文章
const params = reactive({ page_size: 5, page_num: 1 }) // 列表加载参数
const getArticlesInfinite = async ($state: any) => {
  if (loading)
    return

  try {
    const res = await api.getArticles(params)

    // 加载完成
    if (!res.data.length) {
      $state.complete()
      return
    }

    // 非首次加载, 都是往列表中添加数据
    articleList.push(...res.data)
    params.page_num++
    $state.loaded()
  }
  catch (error) {
    $state.error()
  }
}

onMounted(async () => {
  loading = true
  // 首次加载
  const res = await api.getArticles(params)
  articleList = res.data
  params.page_num++
  loading = false
})
</script>

<template>
  <!-- 首页轮播图 -->
  <HomeBanner />
  <div flex-col justify-center>
    <div max-w-1200 mx-auto mb-40 style="margin-top: calc(100vh + 30px)">
      <n-grid :x-gap="12" :y-gap="8" cols="12">
        <n-grid-item span="9">
          <!-- 说说轮播 -->
          <TalkingCarousel />
          <!-- 文章列表 -->
          <ArticleCard
            v-for="(item, idx) in articleList" :key="item.id"
            :article="item" :idx="idx"
          />
          <!-- 无限加载 -->
          <div f-c-c mt-35>
            <InfiniteLoading @infinite="getArticlesInfinite">
              <!-- TODO: 优化界面 -->
              <template #spinner>
                <i-eos-icons:bubble-loading text-30 />
              </template>
              <template #complete>
                <span text-gray>
                  没有更多文章啦!
                </span>
              </template>
            </InfiniteLoading>
          </div>
        </n-grid-item>
        <n-grid-item span="3">
          <!-- sticky 实现悬浮固定效果 -->
          <div sticky top-20>
            <!-- 博主信息 -->
            <AuthorInfo />
            <!-- 公告 -->
            <Announcement my-25 />
            <!-- 网站资讯 -->
            <WebsiteInfo />
          </div>
        </n-grid-item>
      </n-grid>
    </div>
    <!-- 底部 -->
    <AppFooter />
  </div>
</template>

<style lang="scss">
</style>
