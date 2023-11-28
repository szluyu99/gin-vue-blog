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
    <div class="space-x-6">
      <RouterLink v-for="tag of tagList" :key="tag.id" :to="`/tags/${tag.id}?name=${tag.name}`">
        <span class="border-1px border-blue rounded-3rem px-7 py-3 text-12 color-#49b1f5 lg:text-14">
          {{ tag.name }}
        </span>
      </RouterLink>
    </div>

    <div class="flex justify-between text-20 space-x-7" @click="share">
      <span class="i-mdi:wechat cursor-pointer text-green" />
      <span class="i-mdi:sina-weibo cursor-pointer text-orange" />
      <span class="i-mdi:qqchat cursor-pointer text-red" />
      <span class="i-mdi:telegram cursor-pointer text-blue" />
    </div>
  </div>
</template>
