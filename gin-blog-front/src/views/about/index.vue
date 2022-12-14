<script setup lang="ts">
// markdown 相关
import VMdPreview from '@kangc/v-md-editor/lib/preview'
import '@kangc/v-md-editor/lib/style/preview.css'
// markdown 主题
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'

// highlightjs 核心代码
import hljs from 'highlight.js/lib/core'
// 按需引入语言包
import json from 'highlight.js/lib/languages/json'
import javascript from 'highlight.js/lib/languages/javascript'
import go from 'highlight.js/lib/languages/go'
import bash from 'highlight.js/lib/languages/bash'

import api from '@/api'
import { useAppStore } from '@/store'
const { blogConfig } = storeToRefs(useAppStore())

hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
VMdPreview.use(githubTheme, { Hljs: hljs })

let content = $ref('')
onMounted(async () => {
  const res = await api.about()
  content = res.data
})
</script>

<template>
  <BannerPage label="about" title="关于我" card>
    <div text-center>
      <n-image
        width="100"
        :src="blogConfig.website_avatar"
        duration-600 hover-rotate-360
      />
      <VMdPreview mt-15 :text="content" />
    </div>
  </BannerPage>
</template>

<style lang="scss">
</style>
