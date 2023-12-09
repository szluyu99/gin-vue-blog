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
async function getArticlesInfinite($state) {
  if (!loading.value) {
    try {
      const resp = await api.getArticles(params)
      // 加载完成
      if (!resp.data.length) {
        $state.complete()
        return
      }
      // 非首次加载, 都是往列表中添加数据
      articleList.value.push(...resp.data)
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
  const resp = await api.getArticles(params)
  articleList.value = resp.data
  // 过滤 Markdown 符号
  articleList.value.forEach(e => e.content = filterMdSymbol(e.content))
  params.page_num++
  loading.value = false
})

// 过滤 Markdown 符号: 先转 Html 再去除 Html 标签
function filterMdSymbol(md) {
  return marked(md) // 转 HTML
    .replace(/<\/?[^>]*>/g, '') // 正则去除 Html 标签
    .replace(/[|]*\n/, '')
    .replace(/&npsp;/gi, '')
}

function backTop() {
  window.scrollTo({ behavior: 'smooth', top: 0 })
}
</script>

<template>
  <!-- 首页封面图 -->
  <HomeBanner />
  <!-- 内容 -->
  <div class="mx-auto mb-8 max-w-[1230px] flex flex-col justify-center px-3" style="margin-top: calc(100vh + 30px)">
    <div class="grid grid-cols-12 gap-4">
      <!-- 左半部分 -->
      <div class="col-span-12 lg:col-span-9 space-y-5">
        <!-- 说说轮播 -->
        <TalkingCarousel />
        <!-- 文章列表 -->
        <div class="space-y-5">
          <ArticleCard v-for="(item, idx) in articleList" :key="item.id" :article="item" :idx="idx" />
        </div>
        <!-- 无限加载 -->
        <div class="f-c-c">
          <InfiniteLoading class="mt-2 lg:mt-5" @infinite="getArticlesInfinite">
            <!-- TODO: 优化界面 -->
            <template #spinner>
              <span class="animate-pulse text-xl">
                loading...
              </span>
            </template>
            <template #complete>
              <span class="flex gap-2 text-gray">
                没有更多文章啦!
                <button class="flex items-center" @click="backTop">
                  点击回到顶部 <span class="i-mdi:arrow-up-bold-box ml-1 inline-block text-xl" />
                </button>
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
</template>
