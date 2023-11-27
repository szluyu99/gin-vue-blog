<script setup>
import { storeToRefs } from 'pinia'
import { NImage } from 'naive-ui'

import { convertImgUrl } from '@/utils'
import { useAppStore } from '@/store'
// 保持响应式, 否则可能 pinia 中还没拿到数据时, 这里显示无数据
const { blogConfig, articleCount, categoryCount, tagCount } = storeToRefs(useAppStore())

function addToFavorites() {
  window.$message?.info('按 CTRL + D 将本页加入书签')
}
</script>

<template>
  <div class="card-view animate-zoom-in text-center hidden lg:block">
    <NImage
      width="110" :src="convertImgUrl(blogConfig.website_avatar)"
      class="duration-600 hover:rotate-360 py-15"
    />
    <p class="text-24">
      {{ blogConfig.website_author }}
    </p>
    <p class="text-14 my-5">
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
        class="flex items-center justify-center leading-32 h-32 w-7/8 text-14 rounded-1rem
        text-white bg-blue-600 hover:bg-orange-600
        transition-500 ease-in-out
        transform hover:-translate-y-1 hover:scale-110"
        @click="addToFavorites"
      >
        <span class="i-material-symbols:bookmark text-18 mr-5" /> 加入书签
      </button>
    </div>
    <!-- 社交信息 -->
    <div class="text-24 my-15 px-40">
      <a :href="`http://wpa.qq.com/msgrd?v=3&uin=${blogConfig.qq}&site=qq&menu=yes`" target="_blank">
        <span class="i-ant-design:qq-circle-filled inline-block hover-text-red-500" />
      </a>
      <a :href="blogConfig.github" target="_blank" class="mx-20">
        <span class="i-mdi:github inline-block hover:text-red-500" />
      </a>
      <a :href="blogConfig.gitee" target="_blank">
        <span class="i-simple-icons:gitee inline-block hover:text-red-500" />
      </a>
    </div>
  </div>
</template>
