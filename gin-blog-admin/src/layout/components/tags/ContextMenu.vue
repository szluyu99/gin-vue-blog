<script setup>
import { useAppStore, useTagsStore } from '@/store'
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

const tagsStore = useTagsStore()
const appStore = useAppStore()

// 右键菜单选项
const options = computed(() => [
  {
    label: '重新加载',
    key: 'reload',
    disabled: props.currentPath !== tagsStore.activeTag, // 只能重新加载当前标签
    icon: renderIcon('mdi:refresh', { size: '14px' }),
  },
  {
    label: '关闭',
    key: 'close',
    disabled: tagsStore.tags.length <= 1, // 只有一个标签
    icon: renderIcon('mdi:close', { size: '14px' }),
  },
  {
    label: '关闭其他',
    key: 'close-other',
    disabled: tagsStore.tags.length <= 1, // 只有一个标签
    icon: renderIcon('mdi:arrow-expand-horizontal', { size: '14px' }),
  },
  {
    label: '关闭左侧',
    key: 'close-left',
    // 只有一个标签 或者 当前选中的是第一个标签
    disabled: tagsStore.tags.length <= 1 || props.currentPath === tagsStore.tags[0].path,
    icon: renderIcon('mdi:arrow-expand-left', { size: '14px' }),
  },
  {
    label: '关闭右侧',
    key: 'close-right',
    // 只有一个标签 或者 当前选中的是最后一个标签
    disabled:
      tagsStore.tags.length <= 1
      || props.currentPath === tagsStore.tags[tagsStore.tags.length - 1].path,
    icon: renderIcon('mdi:arrow-expand-right', { size: '14px' }),
  },
])

const route = useRoute()
// 重置 KeepAlive
function resetKeepAlive() {
  if (route.meta?.keepAlive)
    appStore.setAliveKeys(route.name, +new Date())
}
// TODO: 标签关闭时会重置 KeepAlive
const actionMap = new Map([
  [
    'reload',
    () => {
      resetKeepAlive()
      appStore.reloadPage()
    },
  ],
  [
    'close',
    () => {
      resetKeepAlive()
      tagsStore.removeTag(props.currentPath)
    },
  ],
  [
    'close-other',
    () => {
      tagsStore.removeOther(props.currentPath)
    },
  ],
  [
    'close-left',
    () => {
      tagsStore.removeLeft(props.currentPath)
    },
  ],
  [
    'close-right',
    () => {
      tagsStore.removeRight(props.currentPath)
    },
  ],
])

function handleHideDropdown() {
  // 配合父组件中使用 v-model:show="show" 来使用
  emit('update:show', false) // 通知父组件更新了 show 的值为 false
}

function handleSelect(key) {
  // console.log('右键菜单选择: ', key)
  const actionFn = actionMap.get(key)
  actionFn && actionFn()
  handleHideDropdown()
}
</script>

<template>
  <n-dropdown
    :show="show"
    :options="options"
    :x="x"
    :y="y"
    placement="bottom-start"
    @clickoutside="handleHideDropdown"
    @select="handleSelect"
  />
</template>
