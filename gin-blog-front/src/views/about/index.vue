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
    <div class="flex justify-center">
      <img :src="appStore.blogConfig.website_avatar" class="w-25 duration-600 hover:rotate-360" alt="avatar">
    </div>
    <div class="flex justify-center">
      <article class="max-w-none prose prose-truegray dark:prose-invert">
        <div v-html="html" />
      </article>
    </div>
    <!-- 后台还没填「关于我」时, 页面上只剩一张头像, 看着像坏了 -->
    <p v-if="!html" class="py-8 text-center text-muted">
      博主还没有填写关于信息
    </p>
  </BannerPage>
</template>
