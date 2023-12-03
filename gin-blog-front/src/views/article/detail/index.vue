<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'

import hljs from 'highlight.js/lib/core'

import BannerInfo from './components/BannerInfo.vue'
import Copyright from './components/Copyright.vue'
import LatestList from './components/LatestList.vue'
import Reward from './components/Reward.vue'
import Forward from './components/Forward.vue'
import LastNext from './components/LastNext.vue'
import Recommend from './components/Recommend.vue'
import Catalogue from './components/Catalogue.vue'

import AppFooter from '@/components/layout/AppFooter.vue'
import Comment from '@/components/comment/Comment.vue'

import { convertImgUrl } from '@/utils'
import api from '@/api'

const data = ref({
  id: 0,
  title: '',
  content: '',
  created_at: '',
  updated_at: '',
  like_count: 0,
  view_count: 0,
  comment_count: 0,
  img: '',
  newest_articles: [],
  tags: [],
  category: {},
  next_article: {},
  last_article: {},
  recommend_articles: [],
})

// 文章内容
const previewRef = ref(null)
const loading = ref(true)

onMounted(() => {
  window.$loadingBar?.start()
  api.getArticleDetail(+useRoute().params.id).then((res) => {
    data.value = res.data
  }).finally(() => {
    loading.value = false
    window.$loadingBar?.finish()
  })

  setTimeout(() => {
    document.querySelectorAll('pre code').forEach((el) => {
      hljs.highlightElement(el)
    })
  }, 100)
})

const styleVal = computed(() =>
  data.value.img
    ? `background: url('${convertImgUrl(data.value.img)}') center center / cover no-repeat;`
    : 'background: rgba(0,0,0,0.1) center center / cover no-repeat;',
)
</script>

<template>
  <!-- 头部 -->
  <div :style="styleVal" class="banner-fade-down absolute inset-x-0 top-0 h-360 f-c-c lg:h-400">
    <BannerInfo v-if="!loading" :article="data" />
  </div>
  <!-- 主体内容 -->
  <main class="flex-1">
    <div class="card-fade-up grid grid-cols-12 mx-auto mb-50 mt-380 gap-15 px-5 lg:mt-440 lg:max-w-1200">
      <!-- 文章主体 -->
      <div class="card-view col-span-12 mx-10 pt-30 lg:col-span-9 lg:mx-0">
        <!-- 文章内容 -->
        <article
          ref="previewRef"
          class="max-w-none prose prose-truegray lg:mx-40"
          v-html="marked(data.content)"
        />
        <!-- 版权声明 -->
        <Copyright class="mb-20 lg:mx-20" />
        <!-- 标签、转发 -->
        <Forward :tag-list="data.tags" class="mb-50 lg:mx-20" />
        <!-- 点赞、打赏 -->
        <Reward
          :article-id="data.id"
          :like-count="data.like_count"
          class="mb-40"
        />
        <!-- 上一篇、下一篇 -->
        <LastNext
          :last-article="data.last_article"
          :next-article="data.next_article"
          class="lg:mx-20"
        />
        <!-- 推荐文章 -->
        <Recommend
          :recommend-list="data.recommend_articles"
          class="mt-30 lg:mx-20"
        />
        <!-- 分隔线 -->
        <hr class="my-40 border-2px border-color-#d2ebfd border-dashed lg:mx-20">
        <!-- 文章评论 -->
        <Comment :type="1" class="lg:mx-20" />
      </div>
      <!-- 文章侧边栏 -->
      <div class="col-span-0 lg:col-span-3">
        <div class="sticky top-20 hidden lg:block">
          <!-- 目录 -->
          <!-- TODO: v-if 的方法不太好, 想办法解决父组件接口获取数据, 子组件渲染问题 -->
          <Catalogue v-if="!loading" :preview="previewRef" />
          <!-- 最新文章 -->
          <LatestList :article-list="data.newest_articles" />
        </div>
      </div>
    </div>
  </main>
  <!-- 底部 -->
  <footer>
    <AppFooter />
  </footer>
</template>
