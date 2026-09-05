<script setup>
import { convertImgUrl } from '@/utils'

defineProps({
  linkList: {
    type: Array,
    default: () => [],
  },
})
</script>

<template>
  <div>
    <!-- 标题: 横幅上已经写着"友情链接"了, 这里只说明数量 -->
    <p class="flex items-center text-xl">
      <span class="i-mdi:link-variant mr-4 text-primary" />
      <span class="color-main font-bold"> 共 {{ linkList.length }} 位朋友 </span>
    </p>

    <!-- 只到两列: 外层卡片宽度固定在 970px 左右, 三列时每格只有 279px,
         65px 头像加上名称简介很挤。按视口断点再加一档没用, 卡片宽度不跟视口变 -->
    <div v-if="linkList.length" class="grid grid-cols-1 mt-4 gap-4 sm:grid-cols-2">
      <!-- 常态就该有卡片感: 原来底色透明、没有阴影, 只在 hover 时用 ::before
           铺一层蓝底 + 白字, 平时看不出这是可点的块, 悬停又跳变太大。
           改成浅底卡片 + 悬停上浮, 和分类页的卡片同一套手感 -->
      <a
        v-for="link of linkList" :key="link.id"
        :href="link.address" target="_blank" rel="noopener noreferrer"
        class="group flex items-center gap-4 rounded-xl bg-surface-soft p-3 shadow-sm transition-300 hover:shadow-md hover:-translate-y-1"
      >
        <!-- 头像: 必须同时写宽高。只写 w- 时图片没加载出来(图床挂了/网络不通)
             高度就是 0, 整行文字会跟着塌上去 —— 实测过 65x0 -->
        <img
          :src="convertImgUrl(link.avatar)"
          class="ring-color-divider h-[65px] w-[65px] shrink-0 rounded-full bg-surface object-cover ring-2 transition-600 group-hover:rotate-360 group-hover:ring-primary"
          loading="lazy" :alt="link.name"
        >
        <!-- 描述: min-w-0 才能让 truncate 生效 -->
        <div class="min-w-0 flex-1">
          <p class="truncate text-lg font-bold transition-300 group-hover:text-primary">
            {{ link.name }}
          </p>
          <p class="truncate text-sm text-muted"> {{ link.intro }} </p>
        </div>
        <!-- 明确"跳外站"的语义, 平时淡, 悬停才清晰 -->
        <span class="i-mdi:open-in-new shrink-0 text-muted opacity-40 transition-300 group-hover:text-primary group-hover:opacity-100" />
      </a>
    </div>

    <!-- 友链数量为 0 -->
    <div v-else class="text-center">
      <img class="inline h-[260px]" src="/images/empty_friend_link.svg" alt="暂无友情链接" loading="lazy">
      <div class="mt-1 space-y-1">
        <p class="text-3xl">
          暂无友情链接
        </p>
        <p> 可以在后台添加 </p>
      </div>
    </div>
  </div>
</template>
