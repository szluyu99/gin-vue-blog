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
    <!-- 标题 -->
    <p class="flex items-center text-xl">
      <span class="i-mdi:link-variant mr-4 text-blue" />
      <span class="color-main font-bold"> 友情链接 </span>
    </p>
    <!-- 链接列表 -->
    <!-- 友链数量不为 0 -->
    <!-- 只到两列: 外层卡片宽度固定在 970px 左右, 三列时每格只有 279px,
         65px 头像加上名称简介很挤。按视口断点再加一档没用, 卡片宽度不跟视口变 -->
    <div v-if="linkList.length" class="grid grid-cols-1 gap-3 sm:grid-cols-2">
      <div
        v-for="link of linkList" :key="link.id"
        class="group link-wrapper relative rounded-8 transition-300"
      >
        <a :href="link.address" target="_blank" class="flex flex-row items-center gap-4 p-3 hover:text-white">
          <!-- 头像: 必须同时写宽高。只写 w- 时图片没加载出来(图床挂了/网络不通)
               高度就是 0, 整行文字会跟着塌上去 —— 实测过 65x0 -->
          <img
            :src="convertImgUrl(link.avatar)"
            class="z-10 h-[65px] w-[65px] shrink-0 rounded-full bg-surface-soft object-cover duration-600 group-hover:rotate-360"
            loading="lazy" :alt="link.name"
          >
          <!-- 描述: 原来写死 260px 宽, 三列下靠 flex 收缩才没溢出;
               改成占满剩余空间, min-w-0 才能让 truncate 生效 -->
          <div class="z-10 min-w-0 flex-1">
            <p class="truncate text-lg font-bold"> {{ link.name }} </p>
            <p class="truncate text-sm"> {{ link.intro }} </p>
          </div>
        </a>
      </div>
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

<style lang="scss" scoped>
.link-wrapper {
  &:hover {
     box-shadow: 0 2px 20px #49b1f5;
  }
  &:hover::before {
    content: "";
    transform: scale(1);
  }
  &::before {
    content: "";
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    background: #49b1f5;
    border-radius: 8px;
    transition-timing-function: ease-out;
    transition-duration: 0.3s;
    transition-property: transform;
    transform: scale(0);
  }
}
</style>
