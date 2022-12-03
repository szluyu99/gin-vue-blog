<script setup lang="ts">
import LinkList from './components/LinkList.vue'
import AddLink from './components/AddLink.vue'
import Comment from '@/components/comment/Comment.vue'
import api from '@/api'

const loading = ref(true)
const linkList = ref<any>([])
onMounted(() => {
  api.getLinks().then((res) => {
    linkList.value = res.data
  }).finally(() => {
    loading.value = false
  })
})
</script>

<template>
  <!-- https://static.talkxj.com/config/9034edddec5b8e8542c2e61b0da1c1da.jpg -->
  <BannerCard
    banner-img="https://static.talkxj.com/config/a741b0656a9a3db2e2ba5c2f4140eb6c.jpg"
    :loading="loading"
  >
    <!-- 友链列表 -->
    <LinkList :link-list="linkList" />
    <!-- 添加友链 -->
    <AddLink />
    <!-- 评论 -->
    <Comment mt-25 :type="2" />
  </BannerCard>
</template>

