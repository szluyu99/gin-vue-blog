<script setup>
defineProps({
  show: {
    type: Boolean,
    default: true,
  },
  color: {
    type: String,
    default: 'green',
  },
  thickness: {
    type: Number,
    default: 4,
  },
  size: {
    type: Number,
    default: 32,
  },
})
</script>

<template>
  <!-- TODO:优化旋转的动画 -->
  <div
    :style="{
      height: `${size}px`,
      width: `${size}px`,
    }"
  >
    <slot name="icon">
      <div
        class="inline-block h-full w-full animate-spin border-gray-300 rounded-full bg-transparent"
        :class="[`border-${thickness}`]"
        :style="{
          'border-right-color': color,
          'border-width': `${thickness}px`,
          'border-style': 'solid',
        }"
      >
        <span class="sr-only">loading</span>
      </div>
    </slot>
  </div>
</template>
