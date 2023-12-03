<script setup>
import { storeToRefs } from 'pinia'
import { convertImgUrl } from '@/utils'
import { useAppStore } from '@/store'

// 保持响应式, 否则可能 pinia 中还没拿到数据时, 这里显示无数据
const { blogConfig, articleCount, categoryCount, tagCount } = storeToRefs(useAppStore())

function addToFavorites() {
  window.$message?.info('按 CTRL + D 将本页加入书签')
}
</script>

<template>
  <div class="hidden animate-zoom-in text-center lg:block card-view">
    <div class="flex justify-center">
      <img class="w-110 py-15 duration-600 hover:rotate-360" :src="convertImgUrl(blogConfig.website_avatar)" alt="author header">
    </div>
    <p class="text-24">
      {{ blogConfig.website_author }}
    </p>
    <p class="my-5 text-14">
      {{ blogConfig.website_intro }}
    </p>
    <!-- 博客信息 -->
    <div class="flex justify-center py-15 text-15">
      <RouterLink to="/archives" class="flex-1">
        <p>文章</p>
        <p class="text-18">
          {{ articleCount }}
        </p>
      </RouterLink>
      <RouterLink to="/categories" class="flex-1">
        <p>分类</p>
        <p class="text-18">
          {{ categoryCount }}
        </p>
      </RouterLink>
      <RouterLink to="/tags" class="flex-1">
        <p>标签</p>
        <p class="text-18">
          {{ tagCount }}
        </p>
      </RouterLink>
    </div>
    <!-- 收藏按钮 -->
    <div class="flex justify-center text-center">
      <button
        class="h-32 w-7/8 f-c-c transform rounded-1rem bg-blue-600 text-14 leading-32 text-white transition-300 ease-in-out hover:scale-110 hover:bg-orange-600 hover:-translate-y-1"
        @click="addToFavorites"
      >
        <span class="i-material-symbols:bookmark mr-5 text-18" /> 加入书签
      </button>
    </div>
    <!-- 社交信息 -->
    <div class="my-15 px-40 text-24 space-x-20">
      <a :href="`http://wpa.qq.com/msgrd?v=3&uin=${blogConfig.qq}&site=qq&menu=yes`" target="_blank">
        <span class="i-ant-design:qq-circle-filled inline-block hover-text-red-500" />
      </a>
      <a :href="blogConfig.github" target="_blank">
        <span class="i-mdi:github inline-block hover:text-red-500" />
      </a>
      <a :href="blogConfig.gitee" target="_blank">
        <span class="i-simple-icons:gitee inline-block hover:text-red-500" />
      </a>
    </div>
  </div>
</template>
