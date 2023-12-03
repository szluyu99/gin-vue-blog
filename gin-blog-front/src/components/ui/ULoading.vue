<script setup>
import Spin from './USpin.vue'

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
  <template v-if="!$slots.default">
    <Spin v-if="show" :color="color" :size="size" :thickness="thickness">
      <template v-if="$slots.icon" #icon>
        <slot name="icon" />
      </template>
      <slot />
    </Spin>
  </template>
  <template v-else>
    <div class="relative">
      <div :class="{ 'opacity-50': show }">
        <slot />
      </div>
      <div class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
        <Spin v-if="show" :color="color" :size="size" :thickness="thickness">
          <template #icon>
            <slot name="icon" />
          </template>
        </Spin>
      </div>
    </div>
  </template>
</template>
