<script setup>
import { computed } from 'vue'
import { renderCustomIcon } from '@/utils'

const props = defineProps({
  icon: {
    type: String,
  },
  color: {
    type: String,
  },
  size: {
    type: Number,
    default: 14,
  },
})

const iconCom = computed(() => renderCustomIcon(props.icon, props))
</script>

<template>
  <component :is="iconCom" />
</template>
