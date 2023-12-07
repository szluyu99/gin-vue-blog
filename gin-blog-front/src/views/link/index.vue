<script setup>
import { onMounted, ref } from 'vue'

import LinkList from './components/LinkList.vue'
import AddLink from './components/AddLink.vue'
import Comment from '@/components/comment/Comment.vue'
import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'

const loading = ref(true)
const linkList = ref([])

onMounted(() => {
  api.getLinks().then((res) => {
    linkList.value = res.data
  }).finally(() => {
    loading.value = false
  })
})
</script>

<template>
  <BannerPage label="link" title="友情链接" card :loading="loading">
    <div class="space-y-5">
      <!-- 友链列表 -->
      <LinkList :link-list="linkList" />
      <!-- 添加友链 -->
      <AddLink />
      <!-- 评论 -->
      <Comment class="mt-30" :type="2" />
    </div>
  </BannerPage>
</template>

