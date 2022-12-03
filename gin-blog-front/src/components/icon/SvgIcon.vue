<script setup lang="ts">
interface Props {
  /** 图标名称(图片的文件名) */
  icon: string
  /** 前缀 */
  prefix?: string
  color?: string
}
const props = withDefaults(defineProps<Props>(), {
  prefix: 'icon-custom',
  color: 'currentColor',
})

// defineOptions({ name: 'SvgIcon' })

const symbolId = computed(() => `#${props.prefix}-${props.icon}`)
</script>

<template>
  <svg aria-hidden="true" width="1em" height="1em">
    <use :xlink:href="symbolId" :fill="color" />
  </svg>
</template>

