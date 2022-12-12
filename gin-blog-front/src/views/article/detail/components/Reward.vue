<script setup lang="ts">
import { useAppStore, useUserStore } from '@/store'
import api from '@/api'

interface Props { articleId: number; likeCount: number }
const { articleId, likeCount } = defineProps<Props>()

// store hook
const [userStore, appStore] = [useUserStore(), useAppStore()]

// 点赞数量
let count = $ref(likeCount)
// * 监听父组件传来的 likeCount, 不能直接用 props 中的值初始化 ref 变量
watch($$(likeCount), newVal => count = newVal)

async function likeArticle() {
  // 判断是否登录
  if (!userStore.userId) {
    appStore.setLoginFlag(true)
    return
  }
  try {
    await api.saveLikeArticle(articleId)
    // 判断是否点赞
    if (userStore.articleLikeSet.includes(articleId)) {
      count--
      window.$message?.info('已取消')
    }
    else {
      count++
      window.$message?.success('已点赞')
    }
    // 维护全局状态中的点赞 Set
    userStore.articleLike(articleId)
  }
  catch (err) {
    console.error(err)
  }
}

// 判断当前用户是否点赞过该文章
const isLike = computed(() => userStore.articleLikeSet.includes(articleId))

function rewardArticle() {
  window.$message?.info('暂时不支持打赏功能')
}
</script>

<template>
  <div f-c-c color="#fff">
    <button
      f-c-c py-6 mr-10 w-110 rounded-2
      :class="isLike ? 'bg-red' : 'bg-gray'"
      @click="likeArticle"
    >
      <i-mdi:thumb-up text-18 mr-3 /> 点赞 {{ count }}
    </button>
    <button f-c-c py-6 ml-10 w-110 rounded-2 bg="blue" @click="rewardArticle">
      <i-mdi:qrcode text-18 mr-3 /> 打赏
    </button>
  </div>
</template>
