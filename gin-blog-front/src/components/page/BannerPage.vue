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

const { pageList } = $(storeToRefs(useAppStore()))

onMounted(() => {
  loading && window.$loadingBar?.start()
})

watch($$(loading), (newVal) => {
  newVal
    ? window.$loadingBar?.start()
    : window.$loadingBar?.finish()
})

// 根据后端配置动态获取封面
const coverStyle = computed(() => {
  const page = pageList.find(e => e.label === label)
  return page
    ? `background: url('${page.cover}') center center / cover no-repeat;`
    : 'background: grey center center / cover no-repeat;'
})
</script>

<template>
  <!-- 顶部图片 -->
  <div
    :style="coverStyle" class="banner-fade-down"
    absolute inset-x-0 top-0 f-c-c
    h-280 lg:h-400
  >
    <h1
      mt-40 text-26 font-bold text-light
      animate-fade-in-down animate-duration-800
      lg="text-40"
    >
      {{ title }}
    </h1>
  </div>
  <!-- 主体内容 -->
  <main flex-1 mx-5>
    <!-- 内容在 spin 中 -->
    <n-spin :show="loading" size="large">
      <!-- 卡片视图 -->
      <template v-if="card">
        <div
          card-view pt-30 mt-300 mb-40 mx-auto min-h-180
          class="card-fade-up"
          lg="px-55 pt-50 pb-30 mt-440 max-w-970"
        >
          <slot v-if="!loading" />
        </div>
      </template>
      <!-- 常规视图 -->
      <template v-else>
        <div
          min-h-400 max-w-1150 px-5 py-40
          class="card-fade-up"
          mx-auto mt-260
          lg:mt-400
        >
          <slot />
        </div>
      </template>
    </n-spin>
  </main>
  <!-- 底部 -->
  <footer v-if="showFooter">
    <AppFooter />
  </footer>
</template>

<style lang="scss">
</style>
