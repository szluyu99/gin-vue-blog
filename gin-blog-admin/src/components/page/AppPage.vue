<script setup>
import { NBackTop } from 'naive-ui'
import AppFooter from '../common/AppFooter.vue'

defineProps({
  showFooter: {
    type: Boolean,
    default: false,
  },
})
</script>

<template>
  <Transition name="fade-slide" mode="out-in" appear>
    <section class="cus-scroll-y h-full w-full flex flex-col bg-[#f5f6fb] p-15 dark:bg-[#121212]">
      <slot />
      <AppFooter v-if="showFooter" class="mt-15" />
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
</style>
