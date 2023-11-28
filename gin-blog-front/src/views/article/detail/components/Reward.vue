<script setup>
import { computed, ref, watch } from 'vue'
import { useAppStore, useUserStore } from '@/store'
import api from '@/api'

const { articleId, likeCount } = defineProps({
  articleId: Number,
  likeCount: Number,

})

const [userStore, appStore] = [useUserStore(), useAppStore()]

// 点赞数量
const count = ref(likeCount)
// * 监听父组件传来的 likeCount, 不能直接用 props 中的值初始化 ref 变量
watch(() => likeCount, newVal => count.value = newVal)

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
      count.value--
      window.$message?.info('已取消')
    }
    else {
      count.value++
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
  <div class="f-c-c color-#fff">
    <button class="mr-10 w-110 f-c-c rounded-2 py-6" :class="isLike ? 'bg-red' : 'bg-gray'" @click="likeArticle">
      <span class="i-mdi:thumb-up mr-3" /> 点赞 {{ count }}
    </button>
    <button class="ml-10 w-110 f-c-c rounded-2 bg-blue py-6" @click="rewardArticle">
      <span class="i-mdi:qrcode mr-3" /> 打赏
    </button>
  </div>
</template>
