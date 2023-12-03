<script setup>
import { convertImgUrl } from '@/utils'

const { linkList } = defineProps({
  linkList: {
    type: Array,
    default: () => [],
  },
})
</script>

<template>
  <div>
    <!-- 标题 -->
    <p class="mb-30 flex items-center">
      <span class="i-mdi:link-variant mr-5 text-28 text-blue" />
      <span class="text-20 font-bold color-#344c67"> 友情链接 </span>
    </p>
    <!-- 友链数量不为 0 -->
    <!-- 链接列表 -->
    <div v-if="!linkList.length" class="grid grid-cols-3 gap-x-12 gap-y-6">
      <div
        v-for="link of linkList" :key="link.id"
        class="link-wrapper group relative col-span-3 rounded-8 transition-300 sm:col-span-1"
      >
        <a :href="link.address" target="_blank" class="flex flex-row p-5 hover:text-white">
          <!-- 头像 -->
          <div class="z-10 mr-5 w-120 f-c-c">
            <img :src="convertImgUrl(link.avatar)" class="w-65 rounded-full duration-600 group-hover:rotate-360">
          </div>
          <!-- 描述 -->
          <div class="z-10 h-95 w-260 f-c-c">
            <div class="text-center">
              <p class="text-18 font-bold"> {{ link.name }} </p>
              <p class="flex-1"> {{ link.intro }} </p>
            </div>
          </div>
        </a>
      </div>
    </div>
    <!-- 友链数量为 0 -->
    <div v-else class="text-center">
      <img class="inline h-260" src="/images/empty_friend_link.svg" alt="暂无友情链接">
      <div class="mt-3 space-y-4">
        <p class="text-30">
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
