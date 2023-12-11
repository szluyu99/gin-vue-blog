<script setup>
import { NBackTop } from 'naive-ui'
import TheFooter from '@/components/common/TheFooter.vue'

defineProps({
  showFooter: { type: Boolean, default: false },
})
</script>

<template>
  <Transition name="fade-slide" mode="out-in" appear>
    <section class="cus-scroll-y h-full w-full flex flex-col bg-[#f5f6fb] p-4 dark:bg-[#121212]">
      <slot />
      <TheFooter v-if="showFooter" class="mt-5" />
      <NBackTop :bottom="20" />
    </section>
  </Transition>
</template>

<style lang="scss">
/* transition fade-slide */
.fade-slide-leave-active,
.fade-slide-enter-active {
  transition: all 0.3s;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

/* 自定义滚动条样式 */
.cus-scroll {
  overflow: auto;
  &::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }
}
.cus-scroll-x {
  overflow-x: auto;
  &::-webkit-scrollbar {
    width: 0;
    height: 8px;
  }
}
.cus-scroll-y {
  overflow-y: auto;
  &::-webkit-scrollbar {
    width: 8px;
    height: 0;
  }
}
.cus-scroll,
.cus-scroll-x,
.cus-scroll-y {
  &::-webkit-scrollbar-thumb {
    background-color: transparent;
    border-radius: 4px;
  }
  &:hover {
    &::-webkit-scrollbar-thumb {
      background: #bfbfbf;
    }
    &::-webkit-scrollbar-thumb:hover {
      background: var(--primary-color);
    }
  }
}
</style>
