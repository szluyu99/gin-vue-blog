<script setup>
import { computed, defineProps, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { NSpin } from 'naive-ui'
import AppFooter from '../layout/AppFooter.vue'
import { useAppStore } from '@/store'

const {
  label = 'default', // 封面
  showFooter = true, // 默认显示底部
  card = false, // 默认不以卡片视图显示
  loading = false,
  title = useRoute().meta?.title, // 默认从路由加载 title
} = defineProps({
  label: String,
  title: String,
  showFooter: Boolean,
  loading: Boolean,
  card: Boolean,
})

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
    ? `background: url('${page.cover}') center center / cover no-repeat;`
    : 'background: grey center center / cover no-repeat;'
})
</script>

<template>
  <!-- 顶部图片 -->
  <div :style="coverStyle" class="banner-fade-down absolute inset-x-0 top-0 h-280 f-c-c lg:h-400">
    <h1 class="mt-40 animate-fade-in-down animate-duration-800 text-26 font-bold text-light lg:text-40">
      {{ title }}
    </h1>
  </div>
  <!-- 主体内容 -->
  <main class="mx-5 flex-1">
    <!-- 内容在 spin 中 -->
    <NSpin :show="loading" size="large">
      <!-- 卡片视图 -->
      <template v-if="card">
        <div class="card-fade-up mx-auto mb-40 mt-300 mt-440 max-w-970 min-h-180 pb-30 pt-30 pt-50 card-view lg:px-55">
          <slot v-if="!loading" />
        </div>
      </template>
      <!-- 常规视图 -->
      <template v-else>
        <div class="card-fade-up mx-auto mt-260 max-w-1150 min-h-400 px-5 py-40 lg:mt-400">
          <slot />
        </div>
      </template>
    </NSpin>
  </main>
  <!-- 底部 -->
  <footer v-if="showFooter">
    <AppFooter />
  </footer>
</template>
