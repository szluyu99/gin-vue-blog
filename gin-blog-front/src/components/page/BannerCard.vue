<script setup lang="ts">
import { useAppStore } from '@/store'

interface Props {
  label?: string
  title?: string
  showFooter?: boolean
  loading?: boolean
}

const {
  label = 'default',
  showFooter = true,
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
    : 'background: url("https://static.talkxj.com/config/83be0017d7f1a29441e33083e7706936.jpg") center center / cover no-repeat;'
})
</script>

<template>
  <!-- 顶部图片 -->
  <div
    :style="coverStyle"
    absolute inset-x-0 top-0 h-400
    flex items-center justify-center
    class="banner-fade-down"
  >
    <h1 text-40 font-bold text-light>
      {{ title }}
    </h1>
  </div>
  <!-- 主体内容 -->
  <main flex-1>
    <!-- 内容在 card 中 -->
    <n-spin :show="loading" size="large">
      <n-card
        hoverable
        shadow-2xl rounded-2rem
        max-w-970 py-50 px-25
        mt-440 mb-40 mx-auto
        class="card-fade-up"
      >
        <slot v-if="!loading" />
      </n-card>
    </n-spin>
  </main>
  <!-- 底部(可选) -->
  <footer v-if="showFooter">
    <AppFooter />
  </footer>
</template>
