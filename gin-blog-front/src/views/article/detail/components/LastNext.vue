<script setup>
import { convertImgUrl } from '@/utils'

const { lastArticle, nextArticle } = defineProps({
  lastArticle: {},
  nextArticle: {},
})
</script>

<template>
  <!-- 封面图压到 40% 不透明度垫在深色底上, 图挂了整块就是底色本身;
       纯黑太生硬, 换成中性深灰 -->
  <div class="flex flex-wrap bg-neutral-800 text-white lg:flex-nowrap">
    <div v-if="lastArticle.id" class="group relative h-[150px] w-full overflow-hidden">
      <RouterLink :to="`/article/${lastArticle.id}`">
        <img :src="convertImgUrl(lastArticle.img)" :alt="lastArticle.title" loading="lazy" class="h-full w-full object-cover opacity-40 transition-600 group-hover:scale-110 group-hover:opacity-80">
        <!-- top:50%; translateY: -50%; 实现绝对定位中的垂直居中 -->
        <div class="absolute top-1/2 w-full px-10 text-left leading-6 -translate-y-1/2">
          <p> 上一篇 </p>
          <p> {{ lastArticle.title }} </p>
        </div>
      </RouterLink>
    </div>
    <div v-if="nextArticle.id" class="group relative h-[150px] w-full overflow-hidden">
      <RouterLink :to="`/article/${nextArticle.id}`">
        <img :src="convertImgUrl(nextArticle.img)" :alt="nextArticle.title" loading="lazy" class="h-full w-full object-cover opacity-40 transition-600 group-hover:scale-110 group-hover:opacity-80">
        <div class="absolute top-1/2 w-full px-10 text-right leading-6 -translate-y-1/2">
          <p> 下一篇 </p>
          <p> {{ nextArticle.title }} </p>
        </div>
      </RouterLink>
    </div>
  </div>
</template>
