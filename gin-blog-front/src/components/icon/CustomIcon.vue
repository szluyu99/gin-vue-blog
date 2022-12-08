<script setup lang="ts">
import { renderCustomIcon } from '@/utils'

interface Props {
  /** 图标名称(图片的文件名) */
  icon: string
  color?: string
  size?: number
}
const props = withDefaults(defineProps<Props>(), {
  size: 14,
})

const iconCom = computed(() => renderCustomIcon(props.icon, props))
</script>

<template>
  <component :is="iconCom" />
</template>
