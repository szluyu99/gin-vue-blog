<script setup>
import { useWindowScroll, useWindowSize } from '@vueuse/core'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import go from 'highlight.js/lib/languages/go'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import { marked } from 'marked'
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'

import Comment from '@/components/comment/Comment.vue'
import AppFooter from '@/components/layout/AppFooter.vue'
import { convertImgUrl } from '@/utils'
import { addCopyButtons } from '@/utils/code-block'
import { typesetMath } from '@/utils/mathjax'
import BannerInfo from './components/BannerInfo.vue'
import Catalogue from './components/Catalogue.vue'
import Copyright from './components/Copyright.vue'
import Forward from './components/Forward.vue'
import LastNext from './components/LastNext.vue'

import LatestList from './components/LatestList.vue'
import Recommend from './components/Recommend.vue'

import Reward from './components/Reward.vue'
import 'highlight.js/styles/a11y-dark.css'

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
    if (!resp.data) {
      return
    }
    // 后端这些字段可能是 null(未分类、没有上一篇、没有推荐), 子组件里直接取 .length / .name,
    // 所以在这里统一退化成默认值, 而不是把 resp.data 整个盖上去
    data.value = {
      ...data.value,
      ...resp.data,
      tags: resp.data.tags ?? [],
      newest_articles: resp.data.newest_articles ?? [],
      recommend_articles: resp.data.recommend_articles ?? [],
      category: resp.data.category ?? {},
      last_article: resp.data.last_article ?? {},
      next_article: resp.data.next_article ?? {},
      // marked 解析 markdown 文本, 正文为空时不能直接丢给 marked
      content: await marked.parse(resp.data.content ?? '', { async: true }),
    }
    await nextTick()
    // highlight.js 代码高亮
    document.querySelectorAll('pre code').forEach(el => hljs.highlightElement(el))
    // 代码块加「复制」按钮
    addCopyButtons(previewRef.value)
    // 正文里有公式才加载 MathJax
    await typesetMath(data.value.content)
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
})

// 阅读进度: 顶部细线, 按滚动位置占可滚动高度的比例
const { y } = useWindowScroll()
const { height: windowHeight } = useWindowSize()
const readProgress = computed(() => {
  const total = document.documentElement.scrollHeight - windowHeight.value
  if (total <= 0) {
    return 0
  }
  return Math.min(100, Math.max(0, (y.value / total) * 100))
})

// 太久没更新的文章给个提示: 技术文章过期得快, 免得读者照着老内容踩坑
const STALE_DAYS = 90
const staleDays = computed(() => {
  if (!data.value.updated_at) {
    return 0
  }
  const days = Math.floor((Date.now() - new Date(data.value.updated_at).getTime()) / 86400000)
  return days >= STALE_DAYS ? days : 0
})

const styleVal = computed(() =>
  data.value.img
    ? `background: url('${convertImgUrl(data.value.img)}') center center / cover no-repeat;`
    : 'background: rgba(0,0,0,0.1) center center / cover no-repeat;',
)
</script>

<template>
  <!-- 阅读进度 -->
  <div
    class="fixed inset-x-0 top-0 z-999 h-0.5 bg-#49b1f5 transition-[width] duration-100"
    :style="{ width: `${readProgress}%` }"
  />
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
        <!-- 老文章提示 -->
        <div
          v-if="!loading && staleDays"
          class="mb-5 border-l-4 border-#f0ad4e rounded bg-#f0ad4e/10 px-4 py-2 text-sm lg:mx-10"
        >
          本文最后更新于 {{ staleDays }} 天前，部分内容可能已经过时。
        </div>
        <article
          ref="previewRef"
          class="max-w-none prose prose-truegray lg:mx-10 dark:prose-invert"
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
        <hr class="my-10 border-2 border-color-divider border-dashed lg:mx-5">
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
