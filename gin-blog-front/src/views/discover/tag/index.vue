<script setup lang="ts">
import api from '@/api'

let loading = $ref(true)
let tagList = $ref<any>({})

onMounted(() => {
  api.getTags().then((res) => {
    tagList = res.data.list
    loading = false
  })
})

// 随机字体大小
function randomFontSize() {
  return Math.floor(Math.random() * 16) + 15
}
// 随机颜色 16进制
function randomColorHex() {
  return `#${Math.floor(Math.random() * 16777215).toString(16)}`
}
</script>

<template>
  <BannerPage :loading="loading" title="标签" label="tag" card>
    <h2
      text-center leading-20 mb-8rem text-25 mt-10
      lg="text-36"
    >
      标签 - {{ tagList.length }}
    </h2>
    <div text-center>
      <router-link
        v-for="t of tagList" :key="t.id" :to="`tags/${t.id}?name=${t.name}`"
        :style="{
          'font-size': `${randomFontSize()}px`,
          'color': `${randomColorHex()}`,
        }"
        px-8 text-8xl leading-45
        inline-block transition-300
        i_hover:text-lightblue hover:scale-130
      >
        {{ t.name }}
      </router-link>
    </div>
  </BannerPage>
</template>

<style lang="scss" scoped>
a {
  // 实现截断文字效果, 即不会在结束处将一个词语拆开
  display: inline-block;
}
</style>
