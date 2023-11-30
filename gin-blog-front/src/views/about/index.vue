<script setup>
import { onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { NImage } from 'naive-ui'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'

import BannerPage from '@/components/page/BannerPage.vue'
import api from '@/api'

import { useAppStore } from '@/store'

const { blogConfig } = storeToRefs(useAppStore())

const content = ref('')

onMounted(async () => {
  const res = await api.about()
  content.value = res.data
  setTimeout(() => {
    document.querySelectorAll('pre code').forEach((el) => {
      hljs.highlightElement(el)
    })
  }, 50)
})
</script>

<template>
  <BannerPage label="about" title="关于我" card>
    <div class="text-center">
      <NImage width="100" :src="blogConfig.website_avatar" class="duration-600 hover-rotate-360" />
    </div>
    <div class="flex justify-center">
      <article class="prose prose-truegray max-w-none">
        <div v-html="marked.parse(content)" />
      </article>
    </div>
  </BannerPage>
</template>
