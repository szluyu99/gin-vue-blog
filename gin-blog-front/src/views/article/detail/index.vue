<script setup>
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import 'highlight.js/styles/a11y-dark.css'
import go from 'highlight.js/lib/languages/go'
import json from 'highlight.js/lib/languages/json'
import javascript from 'highlight.js/lib/languages/javascript'
import bash from 'highlight.js/lib/languages/bash'

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

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)

const route = useRoute()

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

onMounted(async () => {
  try {
    const resp = await api.getArticleDetail(route.params.id)
    data.value = resp.data
    // marked 解析 markdown 文本
    data.value.content = await marked.parse(resp.data.content, { async: true })
    await nextTick()
    // highlight.js 代码高亮
    document.querySelectorAll('pre code').forEach(el => hljs.highlightElement(el))
    // MathJax 渲染公式
    window.MathJax.typeset()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
})

const styleVal = computed(() =>
  data.value.img
    ? `background: url('${convertImgUrl(data.value.img)}') center center / cover no-repeat;`
    : 'background: rgba(0,0,0,0.1) center center / cover no-repeat;',
)
</script>

<template>
  <!-- 头部 -->
  <div :style="styleVal" class="banner-fade-down absolute inset-x-0 top-0 h-[360px] f-c-c lg:h-[400px]">
    <BannerInfo v-if="!loading" :article="data" />
  </div>
  <!-- 主体内容 -->
  <main class="flex-1">
    <div class="card-fade-up grid grid-cols-12 mx-auto mb-3 mt-[380px] gap-4 px-1 lg:mt-[440px] lg:max-w-[1200px]">
      <!-- 文章主体 -->
      <div class="card-view col-span-12 mx-2 pt-7 lg:col-span-9 lg:mx-0">
        <!-- 文章内容 -->
        <article
          ref="previewRef"
          class="max-w-none prose prose-truegray lg:mx-10"
          v-html="data.content"
        />
        <!-- 版权声明 -->
        <Copyright class="my-5 lg:mx-5" />
        <!-- 标签、转发 -->
        <Forward :tag-list="data.tags" class="mb-12 lg:mx-5" />
        <!-- 点赞、打赏 -->
        <Reward
          :article-id="data.id"
          :like-count="data.like_count"
          class="mb-10"
        />
        <!-- 上一篇、下一篇 -->
        <LastNext
          :last-article="data.last_article"
          :next-article="data.next_article"
          class="lg:mx-5"
        />
        <!-- 推荐文章 -->
        <Recommend
          :recommend-list="data.recommend_articles"
          class="mt-7 lg:mx-5"
        />
        <!-- 分隔线 -->
        <hr class="my-10 border-2 border-color-#d2ebfd border-dashed lg:mx-5">
        <!-- 文章评论 -->
        <Comment :type="1" class="lg:mx-5" />
      </div>
      <!-- 文章侧边栏 -->
      <div class="col-span-0 lg:col-span-3">
        <div class="sticky top-5 hidden lg:block space-y-4">
          <!-- 目录 -->
          <!-- TODO: v-if 的方法不太好, 想办法解决父组件接口获取数据, 子组件渲染问题 -->
          <Catalogue v-if="!loading" :preview-ref="previewRef" />
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
