<script setup>
import { convertImgUrl } from '@/utils'

const { lastArticle, nextArticle } = defineProps({
  lastArticle: {},
  nextArticle: {},
})
</script>

<template>
  <div class="flex flex-wrap bg-#000 text-white lg:flex-nowrap">
    <div v-if="lastArticle.id" class="art-card relative h-150 w-full overflow-hidden">
      <RouterLink :to="`/article/${lastArticle.id}`">
        <img class="art-img h-full w-full opacity-40" :src="convertImgUrl(lastArticle.img)">
        <!-- top:50%; translateY: -50%; 实现绝对定位中的垂直居中 -->
        <div class="absolute top-1/2 w-full px-40 leading-25 -translate-y-1/2">
          <p> 上一篇 </p>
          <p> {{ lastArticle.title }} </p>
        </div>
      </RouterLink>
    </div>
    <div v-if="nextArticle.id" class="art-card relative h-150 w-full overflow-hidden">
      <RouterLink :to="`/article/${nextArticle.id}`">
        <img class="art-img h-full w-full opacity-40" :src="convertImgUrl(nextArticle.img)">
        <div class="absolute top-1/2 w-full px-40 text-right leading-25 -translate-y-1/2">
          <p> 下一篇 </p>
          <p> {{ nextArticle.title }} </p>
        </div>
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.art-img {
  transition: all 0.6s;
  object-fit: cover;
}

.art-card:hover .art-img {
  opacity: 0.8;
  transform: scale(1.1);
}
</style>
