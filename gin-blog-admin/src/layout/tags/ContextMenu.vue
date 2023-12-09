<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { NDropdown } from 'naive-ui'

import { useTagStore } from '@/store'
import { renderIcon } from '@/utils'

const props = defineProps({
  show: {
    type: Boolean,
    default: false,
  },
  currentPath: {
    type: String,
    default: '',
  },
  x: {
    type: Number,
    default: 0,
  },
  y: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['update:show'])

const route = useRoute()
const tagStore = useTagStore()

// 右键菜单选项
const options = computed(() => [
  {
    label: '重新加载',
    key: 'reload',
    disabled: props.currentPath !== tagStore.activeTag, // 只能重新加载当前标签
    icon: renderIcon('mdi:refresh', { size: '14px' }),
  },
  {
    label: '关闭',
    key: 'close',
    disabled: tagStore.tags.length <= 1, // 只有一个标签不能关闭
    icon: renderIcon('mdi:close', { size: '14px' }),
  },
  {
    label: '关闭其他',
    key: 'close-other',
    disabled: tagStore.tags.length <= 1, // 只有一个标签不能关闭其他
    icon: renderIcon('mdi:arrow-expand-horizontal', { size: '14px' }),
  },
  {
    label: '关闭左侧',
    key: 'close-left',
    // 只有一个标签 或者 当前选中的是第一个标签, 不能关闭左侧
    disabled: tagStore.tags.length <= 1 || props.currentPath === tagStore.tags[0].path,
    icon: renderIcon('mdi:arrow-expand-left', { size: '14px' }),
  },
  {
    label: '关闭右侧',
    key: 'close-right',
    // 只有一个标签 或者 当前选中的是最后一个标签, 不能关闭右侧
    disabled:
      tagStore.tags.length <= 1
      || props.currentPath === tagStore.tags[tagStore.tags.length - 1].path,
    icon: renderIcon('mdi:arrow-expand-right', { size: '14px' }),
  },
])

// 重置 KeepAlive
function resetKeepAlive() {
  if (route.meta?.keepAlive)
    tagStore.setAliveKey(route.name, +new Date())
}

// TODO: 标签关闭时会重置 KeepAlive
const actionMap = new Map([
  [
    'reload',
    () => {
      resetKeepAlive()
      tagStore.reloadTag()
    },
  ],
  [
    'close',
    () => {
      resetKeepAlive()
      tagStore.removeTag(props.currentPath)
    },
  ],
  [
    'close-other',
    () => tagStore.removeOther(props.currentPath),
  ],
  [
    'close-left',
    () => tagStore.removeLeft(props.currentPath),
  ],
  [
    'close-right',
    () => tagStore.removeRight(props.currentPath),
  ],
])

function handleHideDropdown() {
  emit('update:show', false)
}

function handleSelect(key) {
  const actionFn = actionMap.get(key)
  actionFn && actionFn()
  handleHideDropdown()
}
</script>

<template>
  <NDropdown
    :show="show"
    :options="options"
    :x="x"
    :y="y"
    placement="bottom-start"
    @clickoutside="handleHideDropdown"
    @select="handleSelect"
  />
</template>
