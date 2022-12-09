<script setup lang="ts">
import { useAppStore } from '@/store'

interface Props {
  label?: string
  title?: string
  showFooter?: boolean
  loading?: boolean
  card?: boolean
}

const {
  label = 'default', // 封面
  showFooter = true, // 默认显示底部
  card = false, // 默认不以卡片视图显示
  loading = false,
  title = useRoute().meta?.title, // 默认从路由加载 title
} = defineProps<Props>()

const { pageList } = storeToRefs(useAppStore())

onMounted(() => {
  loading && window.$loadingBar?.start()
})

watch(() => loading, (newVal) => {
  newVal
    ? window.$loadingBar?.start()
    : window.$loadingBar?.finish()
})

// 根据后端配置动态获取封面
const coverStyle = computed(() => {
  const page = pageList.value.find(e => e.label === label)
  return page
    ? `background: url('${page?.cover}') center center / cover no-repeat;`
    : 'background: grey center center / cover no-repeat;'
})
</script>

<template>
  <!-- 顶部图片 -->
  <div
    :style="coverStyle"
    absolute inset-x-0 top-0 h-400 f-c-c
    class="banner-fade-down"
  >
    <h1 text-40 font-bold text-light animate-zoom-in>
      {{ title }}
    </h1>
  </div>
  <!-- 主体内容 -->
  <main flex-1>
    <!-- 内容在 card 中 -->
    <n-spin :show="loading" size="large">
      <!-- 卡片视图 -->
      <template v-if="card">
        <n-card
          hoverable
          shadow-2xl rounded-2rem
          max-w-970 py-50 px-25
          mt-440 mb-40 mx-auto
          class="card-fade-up"
        >
          <slot v-if="!loading" />
        </n-card>
      </template>
      <!-- 常规视图 -->
      <template v-else>
        <div
          min-h-400 max-w-1200 py-50 px-25
          mt-400 mx-auto
          class="card-fade-up"
        >
          <slot />
        </div>
      </template>
    </n-spin>
  </main>
  <!-- 底部 (可选) -->
  <footer v-if="showFooter">
    <AppFooter />
  </footer>
</template>
