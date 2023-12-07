<script setup>
import { onMounted, reactive, ref } from 'vue'

// 无限轮播图
import InfiniteLoading from 'v3-infinite-loading'
// Markdown => Html
import { marked } from 'marked'

import ArticleCard from './components/ArticleCard.vue'
import AuthorInfo from './components/AuthorInfo.vue'
import WebsiteInfo from './components/WebsiteInfo.vue'
import HomeBanner from './components/HomeBanner.vue'
import Announcement from './components/Announcement.vue'
import TalkingCarousel from './components/TalkingCarousel.vue'
import AppFooter from '@/components/layout/AppFooter.vue'

import api from '@/api'

const articleList = ref([])
const loading = ref(false)

// 无限加载文章
const params = reactive({ page_size: 5, page_num: 1 }) // 列表加载参数
const getArticlesInfinite = async ($state) => {
  if (!loading.value) {
    try {
      const res = await api.getArticles(params)
      // 加载完成
      if (!res.data.length) {
        $state.complete()
        return
      }
      // 非首次加载, 都是往列表中添加数据
      articleList.value.push(...res.data)
      // 过滤 Markdown 符号
      articleList.value.forEach(e => e.content = filterMdSymbol(e.content))
      params.page_num++
      $state.loaded()
    }
    catch (error) {
      $state.error()
    }
  }
}

onMounted(async () => {
  loading.value = true
  // 首次加载
  const res = await api.getArticles(params)
  articleList.value = res.data
  // 过滤 Markdown 符号
  articleList.value.forEach(e => e.content = filterMdSymbol(e.content))
  params.page_num++
  loading.value = false
})

// 过滤 Markdown 符号: 先转 Html 再去除 Html 标签
function filterMdSymbol(mdStr) {
  return marked(mdStr) // 转 HTML
    .replace(/<\/?[^>]*>/g, '') // 正则去除 Html 标签
    .replace(/[|]*\n/, '')
    .replace(/&npsp;/gi, '')
}
</script>

<template>
  <!-- 首页封面图 -->
  <HomeBanner />
  <!-- 内容 -->
  <div class="flex flex-col justify-center">
    <div class="mx-auto mb-0 max-w-[1230px] px-4" style="margin-top: calc(100vh + 30px)">
      <div class="grid grid-cols-12 gap-5">
        <!-- 左半部分 -->
        <div class="col-span-12 lg:col-span-9">
          <!-- 说说轮播 -->
          <TalkingCarousel class="mb-5" />
          <!-- 文章列表 -->
          <div class="space-y-5">
            <ArticleCard
              v-for="(item, idx) in articleList" :key="item.id"
              :article="item" :idx="idx"
            />
          </div>
          <!-- 无限加载 -->
          <div ref="el" class="mt-9 f-c-c">
            <InfiniteLoading @infinite="getArticlesInfinite">
              <!-- TODO: 优化界面 -->
              <template #spinner>
                <!-- <i-eos-icons:bubble-loading text-30 /> -->
              </template>
              <template #complete>
                <span text-gray>
                  <!-- 没有更多文章啦! -->
                </span>
              </template>
            </InfiniteLoading>
          </div>
        </div>
        <!-- 右半部分 -->
        <div class="col-span-0 lg:col-span-3">
          <!-- sticky 实现悬浮固定效果 -->
          <div class="sticky top-5 space-y-5">
            <!-- 博主信息 -->
            <AuthorInfo />
            <!-- 公告 -->
            <Announcement />
            <!-- 网站资讯 -->
            <WebsiteInfo />
          </div>
        </div>
      </div>
    </div>
    <!-- 底部 -->
    <AppFooter />
  </div>
</template>
