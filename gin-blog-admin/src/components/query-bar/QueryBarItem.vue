<script setup>
defineProps({
  label: {
    type: String,
    default: '',
  },
  labelWidth: {
    type: Number,
    default: 80,
  },
  contentWidth: {
    type: Number,
    default: 220,
  },
})
</script>

<template>
  <div flex items-center>
    <label v-if="label" flex-shrink-0 :style="{ width: `${labelWidth}px` }">
      {{ label }}
    </label>
    <div :style="{ width: `${contentWidth}px` }" flex-shrink-0>
      <slot />
    </div>
  </div>
</template>
