<script setup>
defineProps({
  showFooter: {
    type: Boolean,
    default: false,
  },
  showHeader: {
    type: Boolean,
    default: true,
  },
  title: {
    type: String,
    default: undefined,
  },
})

const route = useRoute()
</script>

<template>
  <AppPage :show-footer="showFooter">
    <header
      v-if="showHeader"
      flex justify-between items-center
      px-15 mb-15 min-h-45
    >
      <slot v-if="$slots.header" name="header" />
      <template v-else>
        <h2 color="#333" text-22 font-normal>
          {{ title || route.meta?.title }}
        </h2>
        <div>
          <slot name="action" />
        </div>
      </template>
    </header>
    <n-card rounded-10 flex-1>
      <slot />
    </n-card>
  </AppPage>
</template>
