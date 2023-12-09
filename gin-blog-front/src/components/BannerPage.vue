<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'

import AppFooter from './layout/AppFooter.vue'
import ULoading from '@/components/ui/ULoading.vue'

import { useAppStore } from '@/store'

// 注意, 如果使用了解构赋值的形式, watch 会失效
// const {
//   label = 'default', // 封面
//   showFooter = true, // 默认显示底部
//   card = false, // 默认不以卡片视图显示
//   loading = false,
//   title = useRoute().meta?.title, // 默认从路由加载 title
// }

const props = defineProps({
// 封面
  label: {
    type: String,
    default: 'default',
  },
  // 默认显示底部
  showFooter: {
    type: Boolean,
    default: true,
  },
  // 默认不以卡片视图显示
  card: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: () => useRoute().meta?.title, // 默认从路由加载 title
  },
})

const { pageList } = storeToRefs(useAppStore())

// 根据后端配置动态获取封面
const coverStyle = computed(() => {
  const page = pageList.value.find(e => e.label === props.label)
  return page
    ? `background: url('${page.cover}') center center / cover no-repeat;`
    : 'background: grey center center / cover no-repeat;'
})
</script>

<template>
  <!-- 顶部图片 -->
  <div :style="coverStyle" class="banner-fade-down absolute inset-x-0 top-0 h-[280px] f-c-c lg:h-[400px]">
    <h1 class="mt-[40px] animate-fade-in-down animate-duration-800 text-3xl font-bold text-light lg:text-4xl">
      {{ props.title }}
    </h1>
  </div>
  <!-- 主体内容 -->
  <main class="mx-1 mb-10 flex-1">
    <!-- 内容在 spin 中 -->
    <ULoading :show="props.loading">
      <!-- 卡片视图 -->
      <template v-if="props.card">
        <div class="card-view card-fade-up mx-auto mb-10 mt-[300px] max-w-[970px] min-h-[180px] py-8 lg:mt-[440px] lg:px-[55px]">
          <slot v-if="!props.loading" />
        </div>
      </template>
      <!-- 常规视图 -->
      <template v-else>
        <div class="card-fade-up mx-auto mt-[260px] max-w-[1150px] min-h-[400px] px-5 py-10 lg:mt-[400px]">
          <slot />
        </div>
      </template>
    </ULoading>
  </main>
  <!-- 底部 -->
  <AppFooter v-if="props.showFooter" />
</template>
