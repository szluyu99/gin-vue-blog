<script setup>
const { tagList } = defineProps({
  tagList: {
    type: Array,
    default: () => [],
  },
})

function share() {
  window.$message?.info('暂时不支持分享功能')
}
</script>

<template>
  <div class="flex items-center justify-between">
    <div class="space-x-2">
      <RouterLink v-for="tag of tagList" :key="tag.id" :to="`/tags/${tag.id}?name=${tag.name}`">
        <span class="border-1px border-blue rounded-xl px-2 py-1 text-xs color-#49b1f5 lg:text-sm">
          {{ tag.name }}
        </span>
      </RouterLink>
    </div>

    <div class="flex justify-between text-xl space-x-2" @click="share">
      <span class="i-mdi:wechat cursor-pointer text-green" />
      <span class="i-mdi:sina-weibo cursor-pointer text-orange" />
      <span class="i-mdi:qqchat cursor-pointer text-red" />
      <span class="i-mdi:telegram cursor-pointer text-blue" />
    </div>
  </div>
</template>
