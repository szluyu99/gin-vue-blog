<script setup lang="ts">
interface Props { tagList: any }
const { tagList = [] } = defineProps<Props>()

function share() {
  window.$message?.info('暂时不支持分享功能')
}
</script>

<template>
  <div mx-20 mb-50 flex items-center justify-between>
    <div>
      <router-link v-for="tag of tagList" :key="tag.id" :to="`/tags/${tag.id}`">
        <span
          px-9 py-2 mr-8
          border-1px rounded-3rem border-blue color="#49b1f5"
        > {{ tag.name }} </span>
      </router-link>
    </div>

    <div w-150 flex justify-between text-33 cursor-pointer @click="share">
      <i-mdi:wechat text-green />
      <i-mdi:sina-weibo text-orange />
      <i-mdi:qqchat text-red />
      <i-mdi:telegram text-blue />
    </div>
  </div>
</template>
