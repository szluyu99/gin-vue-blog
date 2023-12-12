<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import 'highlight.js/styles/a11y-dark.css'
import go from 'highlight.js/lib/languages/go'
import json from 'highlight.js/lib/languages/json'
import javascript from 'highlight.js/lib/languages/javascript'
import bash from 'highlight.js/lib/languages/bash'

import BannerPage from '@/components/BannerPage.vue'
import { useAppStore } from '@/store'
import api from '@/api'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)

const { blogConfig } = useAppStore()
const html = ref('')

onMounted(async () => {
  const { data } = await api.about()
  // marked 解析 markdown 文本
  html.value = await marked.parse(data, { async: true })
  await nextTick()
  // higlight.js 代码高亮
  document.querySelectorAll('pre code').forEach(el => hljs.highlightElement(el))
  // MathJax 渲染公式
  window.MathJax.typeset()
})
</script>

<template>
  <BannerPage label="about" title="关于我" card>
    <div class="flex justify-center">
      <img :src="blogConfig.website_avatar" class="w-25 duration-600 hover:rotate-360" alt="avatar">
    </div>
    <div class="flex justify-center">
      <article class="max-w-none prose prose-truegray">
        <div v-html="html" />
      </article>
    </div>
  </BannerPage>
</template>
