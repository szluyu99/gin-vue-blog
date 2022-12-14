<script setup lang="ts">
import { convertImgUrl } from '@/utils'
import { useAppStore } from '@/store'
// 保持响应式, 否则可能 pinia 中还没拿到数据时, 这里显示无数据
const { blogConfig, articleCount, categoryCount, tagCount } = storeToRefs(useAppStore())

function addToFavourites() {
  window.$message?.info('按 CTRL + D 将本页加入书签')
}
</script>

<template>
  <div card-view animate-zoom-in text-center hidden lg:block>
    <n-image
      width="110" :src="convertImgUrl(blogConfig.website_avatar)"
      duration-600 hover:rotate-360 py-15
    />
    <p text-24>
      {{ blogConfig.website_author }}
    </p>
    <p text-14 my-5>
      {{ blogConfig.website_intro }}
    </p>
    <!-- 博客信息 -->
    <div flex justify-center py-15 text-15>
      <router-link to="/archives" flex-1>
        <p>文章</p>
        <p text-18>
          {{ articleCount }}
        </p>
      </router-link>
      <router-link to="/categories" flex-1>
        <p>分类</p>
        <p text-18>
          {{ categoryCount }}
        </p>
      </router-link>
      <router-link to="/tags" flex-1>
        <p>标签</p>
        <p text-18>
          {{ tagCount }}
        </p>
      </router-link>
    </div>
    <!-- 收藏按钮 -->
    <div flex justify-center text-center>
      <button
        f-c-c leading-32 h-32 w="7/8" text-14 rounded-1rem
        text-white bg-blue-600 hover:bg-orange-600
        transition-500 ease-in-out
        transform hover:-translate-y-1 hover:scale-110
        @click="addToFavourites"
      >
        <i-material-symbols:bookmark text-18 mr-5 /> 加入书签
      </button>
    </div>
    <!-- 社交信息 -->
    <div text-24 my-15 px-40>
      <a :href="`http://wpa.qq.com/msgrd?v=3&uin=${blogConfig.qq}&site=qq&menu=yes`" target="_blank">
        <i-ant-design:qq-circle-filled hover-text-red-500 />
      </a>
      <a :href="blogConfig.github" target="_blank" mx-20>
        <i-mdi:github hover:text-red-500 />
      </a>
      <a :href="blogConfig.gitee" target="_blank">
        <i-simple-icons:gitee hover:text-red-500 />
      </a>
    </div>
  </div>
</template>
