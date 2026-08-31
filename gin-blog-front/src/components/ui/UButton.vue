<script setup>
defineProps({
  // 'button' | 'a'
  as: {
    type: String,
    default: 'button',
  },
  // 'default' | 'success' | 'info' | 'warning' | 'error' | 'primary' | 'secondary' | 'accent'
  type: {
    type: String,
    default: 'default',
  },
  // 'sm' | 'md' | 'lg'
  size: {
    type: String,
    default: 'md',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})
</script>

<script>
export default {
  name: 'UButton',
  inheritAttrs: false,
}
</script>

<template>
  <div class="flex items-center">
    <component
      :is="as"
      v-bind="$attrs"
      type="button"
      :disabled="disabled"
      :aria-disabled="disabled"
      class="inline-flex cursor-pointer items-center justify-center whitespace-nowrap rounded-md text-xs font-semibold leading-4 font-sans shadow-sm disabled:cursor-not-allowed disabled:opacity-30 focus:outline-none focus:ring-2 focus:ring-offset-2"
      :class="[
        type === 'default'
          ? 'text-gray-700 bg-gray-100 hover:bg-gray-200 focus:ring-gray-500'
          : `text-white bg-${type}-500 hover:bg-${type}-600 focus:ring-${type}-500`,
        {
          'px-2 py-1': size === 'sm',
          'px-3 py-2': size === 'md',
          'px-4 py-3': size === 'lg',
        }]"
    >
      <slot />
    </component>
  </div>
</template>
