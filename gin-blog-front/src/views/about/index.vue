<script setup lang="ts">
// markdown 相关
import VMdPreview from '@kangc/v-md-editor/lib/preview'
import '@kangc/v-md-editor/lib/style/preview.css'
// markdown 主题
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'
// markdown 高亮
import hljs from 'highlight.js'

import api from '@/api'

import { useAppStore } from '@/store'
const { blogInfo } = storeToRefs(useAppStore())

VMdPreview.use(githubTheme, { Hljs: hljs })

const aboutContent = ref('')
onMounted(() => {
  api.about().then((res) => {
    aboutContent.value = res.data
  })
})
</script>

<template>
  <BannerCard banner-img="https://static.talkxj.com/config/2a56d15dd742ff8ac238a512d9a472a1.jpg">
    <div text-center>
      <n-image
        width="110"
        :src="blogInfo.blog_config?.website_avatar"
        duration-600 hover-rotate-360
      />
      <VMdPreview ref="preview" :text="aboutContent" />
    </div>
  </BannerCard>
</template>
