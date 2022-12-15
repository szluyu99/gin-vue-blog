<script setup>
import ContextMenu from './ContextMenu.vue'
import { useAppStore, useTagsStore } from '@/store'
import ScrollX from '@/components/common/ScrollX.vue'

const route = useRoute()
const router = useRouter()
const tagsStore = useTagsStore()
const appStore = useAppStore()

const scrollXRef = $ref(null)
const tabRefs = $ref([])

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
    const { name, path } = route
    const title = route.meta?.title
    tagsStore.addTag({ name, path, title })
  },
  { immediate: true },
)

// 监听当前激活的标签, 标签滚动到让其显示的位置
watch(
  () => tagsStore.activeIndex,
  async (activeIndex) => {
    await nextTick()
    const activeTabElement = tabRefs[activeIndex]?.$el
    if (!activeTabElement)
      return
    const { offsetLeft: x, offsetWidth: width } = activeTabElement
    scrollXRef?.handleScroll(x + width, width)
  },
  { immediate: true },
)

const handleTagClick = (path) => {
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
  setContextMenuShow(false)
  setContextMenu(e.clientX, e.clientY, tagItem.path)
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
  <ScrollX ref="scrollXRef">
    <n-tag
      v-for="tag in tagsStore.tags"
      ref="tabRefs"
      :key="tag.path"
      px-13 mx-5 rounded-4 cursor-pointer hover:color-primary
      :type="tagsStore.activeTag === tag.path ? 'primary' : 'default'"
      :closable="tagsStore.tags.length > 1"
      @click="handleTagClick(tag.path)"
      @close.stop="tagsStore.removeTag(tag.path)"
      @contextmenu.prevent="handleContextMenu($event, tag)"
    >
      <!-- 刷新图标 -->
      <div f-c-c>
        <icon-mdi:refresh text-17 mr-3 @click="handleRefresh" />
        {{ tag.title }}
      </div>
    </n-tag>
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

<style>
.n-tag__close {
  box-sizing: content-box;
  border-radius: 50%;
  font-size: 12px;
  padding: 2px;
  transform: scale(0.9);
  transform: translateX(5px);
  transition: all 0.3s;
}
</style>
