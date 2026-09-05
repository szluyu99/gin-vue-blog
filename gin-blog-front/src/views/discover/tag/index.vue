<script setup>
import { computed, onMounted, ref } from 'vue'
import api from '@/api'
import BannerPage from '@/components/BannerPage.vue'

const loading = ref(true)
const tagList = ref([])

onMounted(() => {
  // 补 catch 并把 loading 放到 finally: 以前失败会留下未捕获的 rejection,
  // 而且页面永远停在加载态
  api.getTags().then((resp) => {
    tagList.value = resp.data || []
  }).catch((err) => {
    console.error(err)
  }).finally(() => {
    loading.value = false
  })
})

// 字号按文章数映射
//
// 之前是 15~30px 纯随机, 标签云"字越大越热门"的语义完全没了 ——
// 实测出现过 6 篇文章的标签比 2 篇的还小。接口的 article_count 直接能用。
const MIN_SIZE = 15
const MAX_SIZE = 30

const countRange = computed(() => {
  const counts = tagList.value.map(t => t.article_count ?? 0)
  return { min: Math.min(...counts), max: Math.max(...counts) }
})

function fontSize(tag) {
  const { min, max } = countRange.value
  // 所有标签文章数一样时没有可比性, 统一取中间值, 不要都顶到最大
  if (!Number.isFinite(min) || max === min) {
    return Math.round((MIN_SIZE + MAX_SIZE) / 2)
  }
  const ratio = ((tag.article_count ?? 0) - min) / (max - min)
  return Math.round(MIN_SIZE + (MAX_SIZE - MIN_SIZE) * ratio)
}

// 配色改成固定色板, 按 id 取
//
// 原来是 `#${Math.floor(Math.random() * 16777215).toString(16)}`, 两个问题:
// 1. 约 6% 的概率生成不足 6 位的十六进制(如 #a3f2c), 是无效色值, 浏览器直接忽略
// 2. 随机色不保证对比度, 可能生成接近背景的浅色, 暗色模式下更没保障
// 按 id 取还有个好处: 同一个标签每次进来颜色是一样的, 不会跳
const COLORS = [
  '#49b1f5', // 主题蓝
  '#ff7242', // 橙
  '#7c5cff', // 紫
  '#00a97f', // 绿
  '#e0508a', // 粉
  '#0e7fd4', // 深蓝
  '#c9971c', // 金
  '#00a3a3', // 青
]

function color(tag) {
  return COLORS[tag.id % COLORS.length]
}
</script>

<template>
  <BannerPage :loading="loading" title="标签" label="tag" card>
    <!-- 横幅上已经写着"标签"了, 这里只说明数量 -->
    <p class="text-center text-muted">
      共 {{ tagList.length }} 个标签
    </p>
    <div class="mt-6 text-center">
      <!-- hover 从 scale-130 收到 110: 站内其他地方都是 110, 放大三成在标签云里
           会把相邻标签挤开; 悬停色也从 lightblue(在白底上几乎看不清)换成强调色 -->
      <RouterLink
        v-for="t of tagList" :key="t.id" :to="`tags/${t.id}?name=${t.name}`"
        :style="{ 'font-size': `${fontSize(t)}px`, 'color': color(t) }"
        :title="`${t.name} - ${t.article_count ?? 0} 篇`"
        class="inline-block px-2 leading-11 transition-300 hover:scale-110 !hover:text-accent"
      >
        {{ t.name }}
      </RouterLink>
    </div>

    <div v-if="!loading && !tagList.length" class="py-10 text-center text-muted">
      还没有标签
    </div>
  </BannerPage>
</template>

<style scoped>
/* 实现截断文字效果, 即不会在结束处将一个词语拆开 */
a {
  display: inline-block;
}
</style>
