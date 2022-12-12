<script setup lang="ts">
import { convertImgUrl } from '@/utils'

interface Props { lastArticle: any; nextArticle: any }
const { lastArticle, nextArticle } = defineProps<Props>()
</script>

<template>
  <div flex flex-wrap bg="#000" text-white lg="flex-nowrap">
    <div
      v-if="lastArticle.id"
      class="art-card" w-full h-150 overflow-hidden relative
    >
      <router-link :to="`/article/${lastArticle.id}`">
        <img
          class="art-img"
          wh-full opacity-40
          :src="convertImgUrl(lastArticle.img)"
        >
        <!-- top:50%; translateY: -50%; 实现绝对定位中的垂直居中 -->
        <div absolute top="1/2" translate-y="-1/2" px-40 w-full leading-25>
          <p> 上一篇 </p>
          <p> {{ lastArticle.title }} </p>
        </div>
      </router-link>
    </div>
    <div
      v-if="nextArticle.id"
      class="art-card" w-full h-150 overflow-hidden relative
    >
      <router-link :to="`/article/${nextArticle.id}`">
        <img
          class="art-img"
          wh-full opacity-40
          :src="convertImgUrl(nextArticle.img)"
        >
        <div absolute top="1/2" translate-y="-1/2" px-40 text-right w-full leading-25>
          <p> 下一篇 </p>
          <p> {{ nextArticle.title }} </p>
        </div>
      </router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.art-img {
  transition: all 0.6s;
  object-fit: cover;
}

.art-card:hover .art-img {
  opacity: 0.8;
  transform: scale(1.1);
}
</style>
