<script setup>
// 无限轮播图
import InfiniteLoading from 'v3-infinite-loading'

import { onMounted, reactive, ref } from 'vue'

import api from '@/api'
import AppFooter from '@/components/layout/AppFooter.vue'
import { stripMarkdown } from '@/utils'
import Announcement from './components/Announcement.vue'
import ArticleCard from './components/ArticleCard.vue'
import AuthorInfo from './components/AuthorInfo.vue'
import HomeBanner from './components/HomeBanner.vue'
import TalkingCarousel from './components/TalkingCarousel.vue'

import WebsiteInfo from './components/WebsiteInfo.vue'

const articleList = ref([])
const loading = ref(false)

// 一页 8 条: 5 条太少, 一屏放得下 2~3 张卡片, 往下滚几下就要再请求一次,
// 每次都要等一小会儿, 滚动看着就是一顿一顿的
const params = reactive({ page_size: 8, page_num: 1 })

// 追加一页, 返回这页拿到几条
async function appendPage() {
  const resp = await api.getArticles(params)
  const list = resp.data?.page_data ?? []
  // 摘要去掉 Markdown 记号
  articleList.value.push(...list.map(e => ({ ...e, content: stripMarkdown(e.content) })))
  params.page_num++
  return list.length
}

// 首屏那次加载的 promise, 供无限加载等待
let firstLoad = null

async function getArticlesInfinite($state) {
  try {
    // 首屏还没回来就先等它: 之前是用 !loading.value 直接跳过, 但跳过时
    // 一个 $state.* 都没调, 组件会一直停在 loading 状态, 要等下一次相交才恢复
    await firstLoad
    const count = await appendPage()
    count ? $state.loaded() : $state.complete()
  }
  catch {
    $state.error()
  }
}

onMounted(() => {
  loading.value = true
  firstLoad = appendPage()
    .catch(err => console.error(err))
    .finally(() => (loading.value = false))
})

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
        <!-- --i 是入场错峰的序号, 见 styles/animate.css 的 card-enter。
             顺序由页面决定, 所以写在这里而不是各个组件里 -->
        <TalkingCarousel :style="{ '--i': 0 }" />
        <!-- 文章列表 -->
        <div class="space-y-5">
          <ArticleCard v-for="(item, idx) in articleList" :key="item.id" :article="item" :idx="idx" />
        </div>
        <!-- 无限加载 -->
        <!-- min-h: 占位高度固定, 否则 loading / 完成提示切换时这一行会跳一下 -->
        <!-- distance: 提前 600px 就开始取下一页, 等滚到底才发请求必然会看到等待 -->
        <!-- firstload=false: 首屏由 onMounted 负责, 不让组件挂载时再打一次 -->
        <div class="min-h-10 f-c-c">
          <InfiniteLoading
            class="mt-2 lg:mt-5"
            :distance="600"
            :firstload="false"
            @infinite="getArticlesInfinite"
          >
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
          <AuthorInfo :style="{ '--i': 0 }" />
          <!-- 公告 -->
          <Announcement :style="{ '--i': 1 }" />
          <!-- 网站资讯 -->
          <WebsiteInfo :style="{ '--i': 2 }" />
        </div>
      </div>
    </div>
  </div>
  <!-- 底部 -->
  <AppFooter />
</template>
