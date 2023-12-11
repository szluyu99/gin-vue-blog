<script setup>
import { computed, h } from 'vue'
import { useRoute } from 'vue-router'
import { NDropdown } from 'naive-ui'

import { useTagStore } from '@/store'

const props = defineProps({
  show: { type: Boolean, default: false },
  currentPath: { type: String, default: '' },
  x: { type: Number, default: 0 },
  y: { type: Number, default: 0 },
})

const emit = defineEmits(['update:show'])

const route = useRoute()
const tagStore = useTagStore()

const options = computed(() => [
  {
    label: '重新加载',
    key: 'reload',
    disabled: props.currentPath !== tagStore.activeTag, // 只能重新加载当前标签
    icon: () => h('i', { class: 'i-mdi:refresh' }),
  },
  {
    label: '关闭',
    key: 'close',
    disabled: tagStore.tags.length <= 1, // 只有一个标签时, 不能关闭
    icon: () => h('i', { class: 'i-mdi:close' }),
  },
  {
    label: '关闭其他',
    key: 'close-other',
    disabled: tagStore.tags.length <= 1, // 只有一个标签时, 不能关闭其他
    icon: () => h('i', { class: 'i-mdi:arrow-expand-horizontal' }),
  },
  {
    label: '关闭左侧',
    key: 'close-left',
    // 只有一个标签 或者 当前选中的是第一个标签, 不能关闭左侧
    disabled: tagStore.tags.length <= 1 || props.currentPath === tagStore.tags[0].path,
    icon: () => h('i', { class: 'i-mdi:arrow-expand-left' }),
  },
  {
    label: '关闭右侧',
    key: 'close-right',
    // 只有一个标签 或者 当前选中的是最后一个标签, 不能关闭右侧
    disabled: tagStore.tags.length <= 1 || props.currentPath === tagStore.tags[tagStore.tags.length - 1].path,
    icon: () => h('i', { class: 'i-mdi:arrow-expand-right' }),
  },
])

const actionMap = new Map([
  [
    'reload',
    () => {
      // 重新加载, 不管是不是 keepAlive, 都要重新获取数据
      tagStore.updateAliveKey(route.name)
      tagStore.reloadTag()
    },
  ],
  [
    'close',
    () => {
      // resetKeepAlive()
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
