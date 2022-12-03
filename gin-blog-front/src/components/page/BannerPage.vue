<script setup lang="ts">
interface Props {
  title?: string
  bannerImg?: string
  showFooter?: boolean
}

const {
  bannerImg = 'https://static.talkxj.com/config/83be0017d7f1a29441e33083e7706936.jpg',
  showFooter = true,
} = defineProps<Props>()

const styleVal = ref(`background: url('${bannerImg}') center center / cover no-repeat;`)

const route = useRoute()
</script>

<template>
  <!-- 顶部图片 -->
  <div
    :style="styleVal"
    absolute inset-x-0 top-0 h-400
    flex items-center justify-center
    class="banner-fade-down"
  >
    <h1 text-40 font-bold text-light animate-zoom-in>
      {{ title || route.meta?.title }}
    </h1>
  </div>
  <!-- 主体内容 -->
  <main flex-1>
    <div
      max-w-1200 py-50 px-25
      mt-380 mb-40 mx-auto
      class="card-fade-up"
    >
      <slot />
    </div>
    <!-- <n-card
      hoverable
      shadow-2xl rounded-2rem
      max-w-970 py-50 px-40
      mt-440 mb-40 mx-auto
      class="card-fade-up"
    >
      <n-skeleton v-if="loading" :height="100" width="100%" />
      <slot v-else />
    </n-card> -->
  </main>
  <!-- 底部(可选) -->
  <footer v-if="showFooter">
    <AppFooter />
  </footer>
</template>

<style lang="scss" scoped>
</style>
