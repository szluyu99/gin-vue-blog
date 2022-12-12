<script setup lang="ts">
import LinkList from './components/LinkList.vue'
import AddLink from './components/AddLink.vue'
import Comment from '@/components/comment/Comment.vue'
import api from '@/api'

let loading = $ref(true)
let linkList = $ref<any>([])

onMounted(() => {
  api.getLinks().then((res) => {
    linkList = res.data
  }).finally(() => {
    loading = false
  })
})
</script>

<template>
  <BannerPage
    label="link" title="友情链接" card
    :loading="loading"
  >
    <!-- 友链列表 -->
    <LinkList :link-list="linkList" />
    <!-- 添加友链 -->
    <AddLink />
    <!-- 评论 -->
    <Comment mt-30 :type="2" />
  </BannerPage>
</template>

