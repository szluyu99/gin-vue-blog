<script setup lang="ts">
// markdown 相关
import VMdPreview from '@kangc/v-md-editor/lib/preview'
import '@kangc/v-md-editor/lib/style/preview.css'
// markdown 主题
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'
// markdown 高亮
import hljs from 'highlight.js'
// 子组件
import BannerInfo from './components/BannerInfo.vue'
import Copyright from './components/Copyright.vue'
import LatestList from './components/LatestList.vue'
import Reward from './components/Reward.vue'
import Forward from './components/Forward.vue'
import LastNext from './components/LastNext.vue'
import Recommend from './components/Recommend.vue'
import Catalogue from './components/Catalogue.vue'
import Comment from '@/components/comment/Comment.vue'

// api
import api from '@/api'

VMdPreview.use(githubTheme, { Hljs: hljs })

const route = useRoute()

const styleVal = ref('background: rgba(0,0,0,0.1) center center / cover no-repeat;')

const data = ref<any>({
  title: '',
  content: '',
  created_at: '',
  updated_at: '',
  like_count: 0,
  view_count: 0,
  img: 'https://static.talkxj.com/articles/3dffb2fcbd541886616ab54c92570de3.jpg',
  newest_articles: [],
  tags: [],
  category: {},
  next_article: {},
  last_article: {},
  recommend_articles: [],
})

// 文章内容
const preview = ref<any>(null)

const loading = ref(true)
onMounted(async () => {
  window.$loadingBar?.start()
  api.getArticleDetail(+route.params.id).then((res) => {
    data.value = res.data
    styleVal.value = `background: url('${data.value.img}') center center / cover no-repeat;`
  }).finally(() => {
    loading.value = false
    window.$loadingBar?.finish()
  })
})
</script>

<template>
  <!-- 头部 -->
  <div
    :style="styleVal" class="banner-fade-down"
    absolute inset-x-0 top-0 h-400 f-c-c
  >
    <!-- 头部信息 -->
    <BannerInfo v-if="!loading" :article="data" />
  </div>
  <!-- 主体内容 -->
  <main flex-1>
    <n-grid
      x-gap="15" cols="12"
      max-w-1200 mt-475 mb-50 mx-auto px-5
      class="card-fade-up"
    >
      <n-grid-item span="9">
        <n-card hoverable rounded-2rem>
          <!-- 文章内容 -->
          <VMdPreview ref="preview" :text="data.content" />
          <!-- 版权声明 -->
          <Copyright mx-20 mb-20 />
          <!-- 标签、转发 -->
          <Forward :tag-list="data.tags" />
          <!-- 点赞、打赏 -->
          <Reward :article-id="data.id" :like-count="data.like_count" />
          <!-- 上一篇、下一篇 -->
          <LastNext :last-article="data.last_article" :next-article="data.next_article" />
          <!-- 推荐文章 -->
          <Recommend :recommend-list="data.recommend_articles" mx-20 mt-30 />
          <!-- 分隔线 -->
          <hr mx-20 my-40 border-dashed border-2px border-color="#d2ebfd">
          <!-- 文章评论 -->
          <Comment :type="1" mx-20 />
        </n-card>
      </n-grid-item>
      <n-grid-item span="3">
        <div sticky top-20>
          <!-- 目录 -->
          <!-- TODO: v-if 的方法不太好, 想办法解决父组件接口获取数据, 子组件渲染问题 -->
          <Catalogue v-if="!loading" :preview="preview" />
          <!-- 最新文章 -->
          <LatestList :article-list="data.newest_articles" />
        </div>
      </n-grid-item>
    </n-grid>
  </main>
  <!-- 底部 -->
  <footer>
    <AppFooter />
  </footer>
</template>
