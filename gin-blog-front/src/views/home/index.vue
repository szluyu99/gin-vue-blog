<script setup lang="ts">
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

import api from '@/api'

let articleList = $ref<any>([])
let loading = $ref(false)

// 无限加载文章
const params = reactive({ page_size: 5, page_num: 1 }) // 列表加载参数
const getArticlesInfinite = async ($state: any) => {
  if (!loading) {
    try {
      const res = await api.getArticles(params)
      // 加载完成
      if (!res.data.length) {
        $state.complete()
        return
      }
      // 非首次加载, 都是往列表中添加数据
      articleList.push(...res.data)
      // 过滤 Markdown 符号
      articleList.forEach((e: any) => e.content = filterMdSymbol(e.content))
      params.page_num++
      $state.loaded()
    }
    catch (error) {
      $state.error()
    }
  }
}

onMounted(async () => {
  loading = true
  // 首次加载
  const res = await api.getArticles(params)
  articleList = res.data
  // 过滤 Markdown 符号
  articleList.forEach((e: any) => e.content = filterMdSymbol(e.content))
  params.page_num++
  loading = false
})

// 过滤 Markdown 符号: 先转 Html 再去除 Html 标签
function filterMdSymbol(mdStr: string) {
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
  <div flex-col justify-center>
    <div max-w-1230 mx-auto mb-40 px-15 style="margin-top: calc(100vh + 30px)">
      <n-grid :x-gap="12" :y-gap="8" cols="9 m:12" responsive="screen">
        <n-gi span="9">
          <!-- 说说轮播 -->
          <TalkingCarousel />
          <!-- 文章列表 -->
          <ArticleCard
            v-for="(item, idx) in articleList" :key="item.id"
            :article="item" :idx="idx"
          />
          <!-- 无限加载 -->
          <div ref="el" f-c-c mt-35>
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
        </n-gi>
        <n-gi span="3">
          <!-- sticky 实现悬浮固定效果 -->
          <div sticky top-20>
            <!-- 博主信息 -->
            <AuthorInfo />
            <!-- 公告 -->
            <Announcement my-25 />
            <!-- 网站资讯 -->
            <WebsiteInfo />
          </div>
        </n-gi>
      </n-grid>
    </div>
    <!-- 底部 -->
    <AppFooter />
  </div>
</template>
