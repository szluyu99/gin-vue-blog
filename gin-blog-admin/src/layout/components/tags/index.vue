<script setup>
import { nextTick, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NTag } from 'naive-ui'

import ContextMenu from './ContextMenu.vue'
import TheIcon from '@/components/icon/TheIcon.vue'
import ScrollX from '@/components/common/ScrollX.vue'
import { useAppStore, useTagsStore } from '@/store'

const route = useRoute()
const router = useRouter()
const tagsStore = useTagsStore()
const appStore = useAppStore()

const scrollXRef = ref(null)
const tabRefs = ref([])

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
    tagsStore.addTag({ name, path, title, icon })
  },
  { immediate: true },
)

// 监听当前激活的标签, 标签滚动到让其显示的位置
watch(
  () => tagsStore.activeIndex,
  async (activeIndex) => {
    await nextTick()
    const activeTabElement = tabRefs.value[activeIndex]?.$el
    if (activeTabElement) {
      const { offsetLeft: x, offsetWidth: width } = activeTabElement
      scrollXRef.value?.handleScroll(x + width, width)
    }
  },
  { immediate: true },
)

function handleTagClick(path) {
  tagsStore.setActiveTag(path) // 激活当前点击的标签
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

// 刷新
function handleRefresh() {
  // 重置 keepAlive
  if (route.meta?.keepAlive)
    appStore.setAliveKeys(route.name, +new Date())
  appStore.reloadPage()
}
</script>

<template>
  <ScrollX ref="scrollXRef" class="bg-white dark:bg-dark!">
    <NTag
      v-for="tag in tagsStore.tags" :key="tag.path"
      ref="tabRefs"
      class="mx-5 cursor-pointer rounded-4 px-13 hover:text-primary"
      :type="tagsStore.activeTag === tag.path ? 'primary' : 'default'"
      :closable="tagsStore.tags.length > 1"
      @click="handleTagClick(tag.path)"
      @close.stop="tagsStore.removeTag(tag.path)"
      @contextmenu.prevent="handleContextMenu($event, tag)"
    >
      <template #icon>
        <TheIcon
          v-if="tag.icon"
          :icon="tag.icon"
          :size="16"
          class="cursor-pointer"
          @click="handleRefresh"
        />
        <div
          v-else
          class="i-mdi:refresh cursor-pointer text-17"
          @click="handleRefresh"
        />
      </template>
      <div class="px-3">
        {{ tag.title }}
      </div>
    </NTag>
    <!-- 自定义右键菜单 -->
    <ContextMenu
      v-if="contextMenuOption.show"
      v-model:show="contextMenuOption.show"
      :current-path="contextMenuOption.currentPath"
      :x="contextMenuOption.x"
      :y="contextMenuOption.y"
    />
  </ScrollX>
</template>
