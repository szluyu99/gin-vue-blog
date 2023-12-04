<script setup>
import { computed, onUnmounted, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  placement: {
    type: String,
    default: 'left',
    validator: val => ['top', 'right', 'bottom', 'left'].includes(val),
  },
  duration: {
    type: String,
    default: '200',
  },
  width: {
    type: Number,
    default: 300,
  },
  height: {
    type: Number,
    default: 300,
  },
})

const emit = defineEmits(['update:modelValue', 'close'])

watch(
  () => props.modelValue,
  (val) => {
    val && emit('update:modelValue', val)
    document.body.style.overflow = (val ? 'hidden' : 'auto')
  },
)

function close() {
  emit('update:modelValue', false)
  emit('close')
}

onUnmounted(() => {
  document.body.style.overflow = 'auto'
})

const styleObject = computed(() => {
  const style = {}
  if (props.placement === 'top' || props.placement === 'bottom') {
    style.width = '100%'
    style.height = `${props.height}px`
  } else if (props.placement === 'left' || props.placement === 'right') {
    style.width = `${props.width}px`
    style.height = '100%'
  }
  return style
})
</script>

<template>
  <Teleport to="body">
    <!-- FIXME: -->
    <div v-if="modelValue" class="fixed inset-0 z-30" @click="close">
      <div class="absolute inset-0 bg-gray-500 opacity-75" />
    </div>
    <div
      v-bind="$attrs"
      class="fixed z-35 overflow-y-auto bg-white transition-all duration-200"
      :class="[`duration-${duration}`, {
        '-left-full top-0': placement === 'left' && !modelValue,
        'left-0 top-0': placement === 'left' && modelValue,

        '-right-full top-0': placement === 'right' && !modelValue,
        'right-0 top-0': placement === 'right' && modelValue,

        '-top-full left-0': placement === 'top' && !modelValue,
        'top-0 left-0 ': placement === 'top' && modelValue,

        '-bottom-full left-0': placement === 'bottom' && !modelValue,
        'bottom-0 left-0': placement === 'bottom' && modelValue,
      }]"
      :style="styleObject"
    >
      <slot />
    </div>
  </Teleport>
</template>
