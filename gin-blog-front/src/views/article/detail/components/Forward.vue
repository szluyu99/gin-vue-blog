<script setup lang="ts">
const { tagList = [] } = defineProps<{ tagList: any }>()

function share() {
  window.$message?.info('暂时不支持分享功能')
}
</script>

<template>
  <div flex items-center justify-between>
    <div>
      <router-link
        v-for="tag of tagList" :key="tag.id"
        :to="`/tags/${tag.id}?name=${tag.name}`"
      >
        <span
          px-7 py-3 mr-3 text-12
          border-1px rounded-3rem border-blue color="#49b1f5"
          lg:text-14
        >
          {{ tag.name }}
        </span>
      </router-link>
    </div>

    <div
      flex justify-between w-120 text-25
      lg="text-28 w-140"
      @click="share"
    >
      <i-mdi:wechat text-green cursor-pointer />
      <i-mdi:sina-weibo text-orange cursor-pointer />
      <i-mdi:qqchat text-red cursor-pointer />
      <i-mdi:telegram text-blue cursor-pointer />
    </div>
  </div>
</template>
