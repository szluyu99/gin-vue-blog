<script setup>
import { NCard } from 'naive-ui'
import AppPage from './AppPage.vue'

defineProps({
  showFooter: { type: Boolean, default: false },
  showHeader: { type: Boolean, default: true },
  title: { type: String, default: undefined },
})
</script>

<template>
  <AppPage :show-footer="showFooter">
    <header
      v-if="showHeader"
      class="mb-3.5 min-h-[45px] flex items-center justify-between px-1"
    >
      <slot v-if="$slots.header" name="header" />
      <template v-else>
        <h2 class="text-2xl text-[#333] font-normal dark:text-[#ccc]">
          {{ title || $route.meta?.title }}
        </h2>
        <div class="space-x-5">
          <slot name="action" />
        </div>
      </template>
    </header>
    <NCard class="flex-1">
      <slot />
    </NCard>
  </AppPage>
</template>
