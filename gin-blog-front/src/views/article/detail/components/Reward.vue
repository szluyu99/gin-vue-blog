<script setup>
import { computed, ref, watch } from 'vue'
import api from '@/api'
import { useAppStore, useUserStore } from '@/store'

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
  <div class="f-c-c space-x-4">
    <!-- 未点赞用主题色而不是灰: 灰色看着像禁用, 实际是可点的。
         已点赞换成强调色, 和站内其他"已激活"状态一致。

         底色写在静态 class 里、已点赞用 !important 覆盖, 不写成
         `isLike ? 'bg-accent' : 'bg-primary'` —— 那样只要 bg-accent 这条 CSS 没生成
         (比如加 accent token 之前启动的 dev server), 按钮就变成白字透明底, 直接消失。
         这么写最差也只是保持蓝色。!important 也免得和 bg-primary 比生成顺序 -->
    <button
      class="w-[110px] f-c-c rounded-md bg-primary py-1.5 text-sm text-white transition-300 hover:opacity-85"
      :class="isLike && '!bg-accent'"
      :aria-pressed="isLike"
      @click="likeArticle"
    >
      <span class="i-mdi:thumb-up mr-1" /> 点赞 {{ count }}
    </button>
    <!-- 打赏是次要动作, 用描边区分主次 -->
    <button
      class="w-[110px] f-c-c border-1 border-primary rounded-md py-1.5 text-sm text-primary transition-300 hover:bg-primary hover:text-white"
      @click="rewardArticle"
    >
      <span class="i-mdi:qrcode mr-1" /> 打赏
    </button>
  </div>
</template>
