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
  <div class="card-view hidden animate-zoom-in animate-duration-600 text-center lg:block space-y-5">
    <div class="mt-4 flex justify-center">
      <img class="w-[105px] duration-600 hover:rotate-360" :src="convertImgUrl(blogConfig.website_avatar)" alt="author avatar">
    </div>
    <div class="space-y-1">
      <p class="text-2xl">
        {{ blogConfig.website_author }}
      </p>
      <p> {{ blogConfig.website_intro }} </p>
    </div>
    <!-- 博客信息 -->
    <div class="flex justify-center">
      <RouterLink to="/archives" class="flex-1">
        <p>文章</p>
        <p class="text-xl">
          {{ articleCount }}
        </p>
      </RouterLink>
      <RouterLink to="/categories" class="flex-1">
        <p>分类</p>
        <p class="text-xl">
          {{ categoryCount }}
        </p>
      </RouterLink>
      <RouterLink to="/tags" class="flex-1">
        <p>标签</p>
        <p class="text-xl">
          {{ tagCount }}
        </p>
      </RouterLink>
    </div>
    <!-- 收藏按钮 -->
    <div class="flex justify-center text-center">
      <button
        class="h-9 w-7/8 f-c-c transform rounded bg-blue-600 text-white transition-200 ease-in-out hover:scale-105 hover:bg-orange-600"
        @click="addToFavorites"
      >
        <span class="i-mdi:bookmark mr-1 text-xl" /> 加入书签
      </button>
    </div>
    <!-- 社交信息 -->
    <div class="text-2xl space-x-4">
      <!-- QQ -->
      <a :href="`http://wpa.qq.com/msgrd?v=3&uin=${blogConfig.qq}&site=qq&menu=yes`" target="_blank">
        <span class="inline-block transition-300 hover:text-red-600">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><path fill="currentColor" d="M21.395 15.035a39.548 39.548 0 0 0-.803-2.264l-1.079-2.695c.001-.032.014-.562.014-.836C19.526 4.632 17.351 0 12 0S4.474 4.632 4.474 9.241c0 .274.013.804.014.836l-1.08 2.695a38.97 38.97 0 0 0-.802 2.264c-1.021 3.283-.69 4.643-.438 4.673c.54.065 2.103-2.472 2.103-2.472c0 1.469.756 3.387 2.394 4.771c-.612.188-1.363.479-1.845.835c-.434.32-.379.646-.301.778c.343.578 5.883.369 7.482.189c1.6.18 7.14.389 7.483-.189c.078-.132.132-.458-.301-.778c-.483-.356-1.233-.646-1.846-.836c1.637-1.384 2.393-3.302 2.393-4.771c0 0 1.563 2.537 2.103 2.472c.251-.03.581-1.39-.438-4.673zM12.662 4.846c.039-1.052.659-1.878 1.385-1.846s1.281.912 1.242 1.964c-.039 1.051-.659 1.878-1.385 1.846s-1.282-.912-1.242-1.964zM9.954 3c.725-.033 1.345.794 1.384 1.846c.04 1.052-.517 1.931-1.242 1.963c-.726.033-1.346-.794-1.385-1.845C8.672 3.912 9.228 3.033 9.954 3zM7.421 8.294c.194-.43 2.147-.908 4.566-.908h.026c2.418 0 4.372.479 4.566.908a.14.14 0 0 1 .014.061c0 .031-.01.059-.026.083c-.163.238-2.333 1.416-4.553 1.416h-.026c-2.221 0-4.39-1.178-4.553-1.416a.136.136 0 0 1-.014-.144zm10.422 8.622c-.22 3.676-2.403 5.987-5.774 6.021h-.137c-3.37-.033-5.554-2.345-5.773-6.021c-.081-1.35.001-2.496.147-3.43c.318.063.638.122.958.176v3.506s1.658.334 3.318.103v-3.225c.488.027.96.04 1.406.034h.025c1.678.021 3.714-.204 5.683-.594c.146.934.227 2.08.147 3.43zM10.48 5.804c.313-.041.542-.409.508-.825c-.033-.415-.314-.72-.629-.679c-.313.04-.541.409-.508.824c.034.417.315.72.629.68zm3.999-.648c.078.037.221.042.289-.146c.035-.095.025-.165-.009-.214c-.023-.033-.133-.118-.371-.176c-.904-.22-1.341.384-1.405.499c-.04.072-.012.176.056.227c.067.051.139.037.179-.006c.58-.628 1.21-.208 1.261-.184z" /></svg>
        </span>
      </a>
      <!-- Github -->
      <a :href="blogConfig.github" target="_blank">
        <span class="inline-block transition-300 hover:text-red-600">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><path fill="currentColor" d="M12 .297c-6.63 0-12 5.373-12 12c0 5.303 3.438 9.8 8.205 11.385c.6.113.82-.258.82-.577c0-.285-.01-1.04-.015-2.04c-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729c1.205.084 1.838 1.236 1.838 1.236c1.07 1.835 2.809 1.305 3.495.998c.108-.776.417-1.305.76-1.605c-2.665-.3-5.466-1.332-5.466-5.93c0-1.31.465-2.38 1.235-3.22c-.135-.303-.54-1.523.105-3.176c0 0 1.005-.322 3.3 1.23c.96-.267 1.98-.399 3-.405c1.02.006 2.04.138 3 .405c2.28-1.552 3.285-1.23 3.285-1.23c.645 1.653.24 2.873.12 3.176c.765.84 1.23 1.91 1.23 3.22c0 4.61-2.805 5.625-5.475 5.92c.42.36.81 1.096.81 2.22c0 1.606-.015 2.896-.015 3.286c0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12" /></svg>
        </span>
      </a>
      <!-- Gitee -->
      <a :href="blogConfig.gitee" target="_blank">
        <span class="inline-block transition-300 hover:text-red-600">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><path fill="currentColor" d="M11.984 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12a12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.016 0zm6.09 5.333c.328 0 .593.266.592.593v1.482a.594.594 0 0 1-.593.592H9.777c-.982 0-1.778.796-1.778 1.778v5.63c0 .327.266.592.593.592h5.63c.982 0 1.778-.796 1.778-1.778v-.296a.593.593 0 0 0-.592-.593h-4.15a.592.592 0 0 1-.592-.592v-1.482a.593.593 0 0 1 .593-.592h6.815c.327 0 .593.265.593.592v3.408a4 4 0 0 1-4 4H5.926a.593.593 0 0 1-.593-.593V9.778a4.444 4.444 0 0 1 4.445-4.444h8.296Z" /></svg>
        </span>
      </a>
    </div>
  </div>
</template>
