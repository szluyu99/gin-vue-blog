<script setup lang="ts">
// markdown 相关
import VMdPreview from '@kangc/v-md-editor/lib/preview'
import '@kangc/v-md-editor/lib/style/preview.css'
// markdown 主题
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'
// markdown 高亮
import hljs from 'highlight.js'

import { convertImgUrl } from '@/utils'
import api from '@/api'

import { useAppStore } from '@/store'
const { blogConfig } = storeToRefs(useAppStore())

VMdPreview.use(githubTheme, { Hljs: hljs })

let content = $ref('')
onMounted(async () => {
  const res = await api.about()
  content = res.data
})
</script>

<template>
  <BannerCard label="about" title="关于我">
    <div text-center>
      <n-image
        width="110"
        :src="convertImgUrl(blogConfig.website_avatar)"
        duration-600 hover-rotate-360
      />
      <VMdPreview :text="content" />
    </div>
  </BannerCard>
</template>
