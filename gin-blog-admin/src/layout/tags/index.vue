<script setup>
import { NTag } from 'naive-ui'
import { nextTick, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ScrollX from '@/components/common/ScrollX.vue'
import TheIcon from '@/components/icon/TheIcon.vue'
import { useTagStore } from '@/store'
import ContextMenu from './ContextMenu.vue'

const route = useRoute()
const router = useRouter()
const tagStore = useTagStore()

const scrollXRef = ref(null)
// 按 path 存标签元素: v-for 上的模板 ref 数组 Vue 不保证顺序与 tags 一致,
// 原来是 tabRefs.value[activeIndex], 顺序一旦错位滚动定位就指到别的标签上
const tabEls = new Map()
function setTabRef(path, el) {
  if (el) {
    tabEls.set(path, el.$el ?? el)
  }
  else {
    tabEls.delete(path)
  }
}

const contextMenuOption = reactive({
  show: false,
  x: 0,
  y: 0,
  currentPath: '',
})

// 监听当前路由路径, 发生变化则添加到标签栏
watch(
  () => route.path,
  () => {
    const { name, fullPath: path } = route
    const title = route.meta?.title
    const icon = route.meta?.icon
    tagStore.addTag({ name, path, title, icon })
  },
  { immediate: true },
)

// 监听当前激活的标签, 标签滚动到让其显示的位置
watch(
  () => tagStore.activeTag,
  async (activePath) => {
    await nextTick()
    const activeTabElement = tabEls.get(activePath)
    if (activeTabElement) {
      const { offsetLeft: x, offsetWidth: width } = activeTabElement
      scrollXRef.value?.handleScroll(x + width, width)
    }
  },
  { immediate: true },
)

function handleTagClick(path) {
  tagStore.setActiveTag(path) // 激活当前点击的标签
  router.push(path)
}

// 显示或隐藏右键菜单
function setContextMenuShow(flag) {
  contextMenuOption.show = flag
}

function setContextMenu(x, y, currentPath) {
  // Object.assign(a, b) 将 b 的属性拷贝到 a 身上(相同覆盖), 浅拷贝
  Object.assign(contextMenuOption, { x, y, currentPath })
}

// 右击菜单
async function handleContextMenu(e, tagItem) {
  const { clientX, clientY } = e
  setContextMenuShow(false)
  setContextMenu(clientX, clientY, tagItem.path)
  await nextTick()
  setContextMenuShow(true)
}

function handleRefresh(tag) {
  // 只有当前标签会刷新
  if (route.name === tag.name) {
    tagStore.updateAliveKey(route.name)
    tagStore.reloadTag()
  }
}
</script>

<template>
  <ScrollX ref="scrollXRef" class="bg-white dark:bg-dark!">
    <NTag
      v-for="tag in tagStore.tags" :key="tag.path"
      :ref="el => setTabRef(tag.path, el)"
      class="mx-1 hover:border-blue hover:border-red hover:text-primary"
      :type="tagStore.activeTag === tag.path ? 'primary' : 'default'"
      :closable="tagStore.tags.length > 1"
      @click="handleTagClick(tag.path)"
      @close.stop="tagStore.removeTag(tag.path)"
      @contextmenu.prevent="handleContextMenu($event, tag)"
    >
      <template #icon>
        <!-- 用 useRoute() 拿到的 route, 不再混用模板里的全局 $route -->
        <div :class="{ 'cursor-pointer': route.name === tag.name }" @click="handleRefresh(tag)">
          <TheIcon v-if="tag.icon" :icon="tag.icon" :size="16" />
          <i v-else class="i-mdi:refresh" />
        </div>
      </template>
      <div class="px-0.5">
        {{ tag.title }}
      </div>
    </NTag>
    <ContextMenu
      v-if="contextMenuOption.show"
      v-model:show="contextMenuOption.show"
      :current-path="contextMenuOption.currentPath"
      :x="contextMenuOption.x"
      :y="contextMenuOption.y"
    />
  </ScrollX>
</template>
