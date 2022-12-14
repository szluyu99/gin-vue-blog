<script setup lang="ts">
import VMdPreview from '@kangc/v-md-editor/lib/preview'
import '@kangc/v-md-editor/lib/style/preview.css'
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'

// 引入全部语言包 (不推荐)
// import hljs from 'highlight.js'
// highlightjs 核心代码
import hljs from 'highlight.js/lib/core'
// 按需引入语言包
import json from 'highlight.js/lib/languages/json'
import javascript from 'highlight.js/lib/languages/javascript'
import go from 'highlight.js/lib/languages/go'
import bash from 'highlight.js/lib/languages/bash'

import BannerInfo from './components/BannerInfo.vue'
import Copyright from './components/Copyright.vue'
import LatestList from './components/LatestList.vue'
import Reward from './components/Reward.vue'
import Forward from './components/Forward.vue'
import LastNext from './components/LastNext.vue'
import Recommend from './components/Recommend.vue'
import Catalogue from './components/Catalogue.vue'
import Comment from '@/components/comment/Comment.vue'

import { convertImgUrl } from '@/utils'
import api from '@/api'

hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
VMdPreview.use(githubTheme, { Hljs: hljs })

let data = $ref<any>({
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
const previewRef = $ref<any>(null)
let loading = $ref(true)

onMounted(async () => {
  // const rendererMD = new marked.Renderer()
  // marked.setOptions({
  //   renderer: rendererMD,
  //   highlight(code) {
  //     return hljs.highlightAuto(code).value
  //   },
  //   pedantic: false,
  //   gfm: true,
  //   // tables: true,
  //   breaks: false,
  //   sanitize: false,
  //   smartLists: true,
  //   smartypants: false,
  //   xhtml: false,
  // })

  window.$loadingBar?.start()
  try {
    const res = await api.getArticleDetail(+useRoute().params.id)
    data = res.data
  }
  finally {
    loading = false
    window.$loadingBar?.finish()
  }
})

const styleVal = $computed(() =>
  data.img
    ? `background: url('${convertImgUrl(data.img)}') center center / cover no-repeat;`
    : 'background: rgba(0,0,0,0.1) center center / cover no-repeat;',
)

onMounted(() => {
  const link = document.createElement('link')
  link.type = 'text/css'
  link.rel = 'stylesheet'
  link.href = 'https://cdn.bootcss.com/github-markdown-css/2.10.0/github-markdown.min.css'
  document.head.appendChild(link)
})
</script>

<template>
  <!-- 头部 -->
  <div
    :style="styleVal" class="banner-fade-down"
    absolute inset-x-0 top-0 h-360 f-c-c
    lg:h-400
  >
    <BannerInfo v-if="!loading" :article="data" />
  </div>
  <!-- 主体内容 -->
  <main flex-1>
    <n-grid
      x-gap="15" cols="12"
      max-w-1200 mt-380 mb-50 mx-auto px-5
      class="card-fade-up"
      responsive="screen"
      item-responsive
      lg:mt-440
    >
      <n-gi span="12 m:9">
        <div card-view pt-30>
          <!-- 文章内容 -->
          <VMdPreview
            ref="previewRef"
            :text="data.content"
            lg:mx-20
          />
          <!-- <div
            ref="previewRef"
            class="markdown-body"
            v-html="marked(data.content)"
          /> -->
          <!-- 版权声明 -->
          <Copyright mb-20 lg:mx-20 />
          <!-- 标签、转发 -->
          <Forward
            :tag-list="data.tags"
            mb-50 lg:mx-20
          />
          <!-- 点赞、打赏 -->
          <Reward
            :article-id="data.id"
            :like-count="data.like_count"
            mb-40
          />
          <!-- 上一篇、下一篇 -->
          <LastNext
            :last-article="data.last_article"
            :next-article="data.next_article"
            lg:mx-20
          />
          <!-- 推荐文章 -->
          <Recommend
            :recommend-list="data.recommend_articles"
            mt-30 lg:mx-20
          />
          <!-- 分隔线 -->
          <hr
            my-40 border-dashed border-2px border-color="#d2ebfd"
            lg:mx-20
          >
          <!-- 文章评论 -->
          <Comment :type="1" lg:mx-20 />
        </div>
      </n-gi>
      <n-gi span="0 m:3">
        <div sticky top-20 hidden lg:block>
          <!-- 目录 -->
          <!-- TODO: v-if 的方法不太好, 想办法解决父组件接口获取数据, 子组件渲染问题 -->
          <Catalogue v-if="!loading" :preview="previewRef" />
          <!-- 最新文章 -->
          <LatestList :article-list="data.newest_articles" />
        </div>
      </n-gi>
    </n-grid>
  </main>
  <!-- 底部 -->
  <footer>
    <AppFooter />
  </footer>
</template>

<style lang="scss">
  .github-markdown-body {
    padding: 0;
  }
</style>
