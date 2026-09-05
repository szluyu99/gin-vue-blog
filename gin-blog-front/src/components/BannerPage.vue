<script setup>
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import USkeleton from '@/components/ui/USkeleton.vue'
import { useAppStore } from '@/store'

import AppFooter from './layout/AppFooter.vue'

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
    <h1 class="mt-[40px] animate-fade-in-down animate-duration-800 text-3xl text-light font-bold lg:text-4xl">
      {{ props.title }}
    </h1>
  </div>
  <!-- 主体内容 -->
  <main class="mx-1 mb-10 flex-1">
    <!-- 加载态用骨架屏而不是「半透明空卡片 + 转圈」: 版式先占住, 内容替换进来时
         位置基本不动; 骨架和内容之间再做一次淡入淡出, 免得啪一下跳出来 -->
    <!-- 卡片视图 -->
    <template v-if="props.card">
      <div class="card-view card-fade-up mx-auto mb-10 mt-[300px] max-w-[970px] min-h-[180px] py-8 lg:mt-[440px] lg:px-[55px]">
        <!-- 必须显式给 duration: 骨架屏上的 animate-pulse 是 2s 无限循环动画,
             Vue 会自动嗅探离场元素上最长的 animation 时长并等它跑完 ——
             实测数据 326ms 就到手了, 内容却要等到 2356ms 才出现 -->
        <Transition name="fade" mode="out-in" :duration="{ enter: 250, leave: 150 }">
          <USkeleton v-if="props.loading" />
          <div v-else>
            <slot />
          </div>
        </Transition>
      </div>
    </template>
    <!-- 常规视图 -->
    <template v-else>
      <div class="card-fade-up mx-auto mt-[260px] max-w-[1150px] min-h-[400px] px-5 py-10 lg:mt-[400px]">
        <!-- 必须显式给 duration: 骨架屏上的 animate-pulse 是 2s 无限循环动画,
             Vue 会自动嗅探离场元素上最长的 animation 时长并等它跑完 ——
             实测数据 326ms 就到手了, 内容却要等到 2356ms 才出现 -->
        <Transition name="fade" mode="out-in" :duration="{ enter: 250, leave: 150 }">
          <USkeleton v-if="props.loading" :rows="8" />
          <div v-else>
            <slot />
          </div>
        </Transition>
      </div>
    </template>
  </main>
  <!-- 底部 -->
  <AppFooter v-if="props.showFooter" />
</template>
