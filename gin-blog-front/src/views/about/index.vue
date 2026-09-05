<script setup>
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import go from 'highlight.js/lib/languages/go'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import { marked } from 'marked'
import { nextTick, onMounted, ref } from 'vue'
import api from '@/api'

import BannerPage from '@/components/BannerPage.vue'
import { useAppStore } from '@/store'
import { convertImgUrl } from '@/utils'
import { addCopyButtons } from '@/utils/code-block'
import { typesetMath } from '@/utils/mathjax'
import 'highlight.js/styles/a11y-dark.css'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)

// 不解构: blogConfig 是 getter, store 里 blog_config 是整体重新赋值的,
// 解构后拿到的是接口回来之前的旧对象, 头像不会更新
const appStore = useAppStore()
const html = ref('')

onMounted(async () => {
  try {
    const { data } = await api.about()
    // marked 解析 markdown 文本, 内容为空时不能直接丢给 marked
    html.value = await marked.parse(data ?? '', { async: true })
  }
  catch (err) {
    console.error(err)
    return
  }
  await nextTick()
  // higlight.js 代码高亮
  document.querySelectorAll('pre code').forEach(el => hljs.highlightElement(el))
  // 代码块加「复制」按钮, 和文章详情页一致
  addCopyButtons(document.querySelector('article'))
  // 内容里有公式才加载 MathJax
  await typesetMath(html.value)
})
</script>

<template>
  <BannerPage label="about" title="关于我" card>
    <!-- 名片区: 原来这页只有一张头像加一段正文, 空得慌。
         这里用的都是已有数据(站点配置 + 首页那个统计接口), 没有新增请求 -->
    <div class="flex flex-col items-center">
      <img
        :src="convertImgUrl(appStore.blogConfig.website_avatar)"
        class="ring-color-divider h-25 w-25 rounded-full bg-surface-soft object-cover ring-4 transition-600 hover:rotate-360 hover:ring-primary"
        alt="博主头像"
      >
      <p class="mt-4 text-2xl font-bold">
        {{ appStore.blogConfig.website_author }}
      </p>
      <p class="mt-1 text-muted">
        {{ appStore.blogConfig.website_intro }}
      </p>

      <!-- 社交链接: 首页侧栏有, 这一页原来没有 -->
      <div class="mt-4 flex items-center gap-5 text-2xl">
        <a
          :href="`http://wpa.qq.com/msgrd?v=3&uin=${appStore.blogConfig.qq}&site=qq&menu=yes`"
          target="_blank" rel="noopener noreferrer" title="QQ"
        >
          <span class="i-ant-design:qq-circle-filled block transition-300 hover:text-accent" />
        </a>
        <a :href="appStore.blogConfig.github" target="_blank" rel="noopener noreferrer" title="GitHub">
          <span class="i-mdi:github block transition-300 hover:text-accent" />
        </a>
        <a :href="appStore.blogConfig.gitee" target="_blank" rel="noopener noreferrer" title="Gitee">
          <span class="i-simple-icons:gitee block transition-300 hover:text-accent" />
        </a>
      </div>

      <!-- 统计: 三项各自链到对应页面, 访问量没有对应页面所以不做链接 -->
      <div class="grid grid-cols-4 mt-6 max-w-md w-full gap-2 text-center">
        <RouterLink to="/archives" class="transition-300 hover:text-primary">
          <p class="text-xl font-bold">
            {{ appStore.articleCount }}
          </p>
          <p class="text-sm text-muted">
            文章
          </p>
        </RouterLink>
        <RouterLink to="/categories" class="transition-300 hover:text-primary">
          <p class="text-xl font-bold">
            {{ appStore.categoryCount }}
          </p>
          <p class="text-sm text-muted">
            分类
          </p>
        </RouterLink>
        <RouterLink to="/tags" class="transition-300 hover:text-primary">
          <p class="text-xl font-bold">
            {{ appStore.tagCount }}
          </p>
          <p class="text-sm text-muted">
            标签
          </p>
        </RouterLink>
        <div>
          <p class="text-xl font-bold">
            {{ appStore.viewCount }}
          </p>
          <p class="text-sm text-muted">
            访问
          </p>
        </div>
      </div>
    </div>

    <hr class="my-8 border-color-divider border-dashed">

    <div class="flex justify-center">
      <article class="max-w-none prose prose-truegray dark:prose-invert">
        <div v-html="html" />
      </article>
    </div>
    <!-- 后台还没填「关于我」时, 分隔线下面空空的, 看着像坏了 -->
    <p v-if="!html" class="py-4 text-center text-muted">
      博主还没有填写关于信息
    </p>
  </BannerPage>
</template>
