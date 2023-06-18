<script setup>
import { useRoute } from 'vue-router'
import { NCard } from 'naive-ui'
import AppPage from './AppPage.vue'

defineProps({
  showFooter: { type: Boolean, default: true },
  showHeader: { type: Boolean, default: true },
  title: { type: String, default: undefined },
})

const route = useRoute()
</script>

<template>
  <AppPage :show-footer="showFooter">
    <header
      v-if="showHeader"
      class="mb-15 min-h-45 flex items-center justify-between px-15"
    >
      <slot v-if="$slots.header" name="header" />
      <template v-else>
        <h2 class="text-22 font-normal text-[#333] dark:text-[#ccc]">
          {{ title || route.meta?.title }}
        </h2>
        <div class="space-x-20">
          <slot name="action" />
        </div>
      </template>
    </header>
    <NCard class="flex-1 rounded-10">
      <slot />
    </NCard>
  </AppPage>
</template>
