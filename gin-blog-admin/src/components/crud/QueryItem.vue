<script setup>
defineProps({
  label: { type: String, default: '' },
  labelWidth: { type: Number, default: 80 },
  contentWidth: { type: Number, default: 220 },
})
</script>

<template>
  <div class="flex items-center">
    <label v-if="label" class="flex-shrink-0" :style="{ width: `${labelWidth}px` }">
      {{ label }}
    </label>
    <div class="flex-shrink-0" :style="{ width: `${contentWidth}px` }">
      <slot />
    </div>
  </div>
</template>
