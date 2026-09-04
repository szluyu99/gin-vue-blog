<script setup>
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { onMounted, ref } from 'vue'

const { previewRef } = defineProps({
  previewRef: { type: Object, required: true },
})

onMounted(() => {
  buildAnchors()
})

const selectAnchor = ref('')
const anchors = ref([])
// 正文还没渲染出来时 previewRef 可能是 null, 直接 querySelectorAll 会抛
// 空标题(markdown 里写了 ## 但没写文字)不进目录, 否则目录里会多出一行空白
const headings = previewRef
  ? Array.from(previewRef.querySelectorAll('h1,h2,h3,h4,h5,h6')).filter(t => !!t.textContent.trim())
  : []

function buildAnchors() {
  // 用于确认层级
  const hTags = Array.from(new Set(headings.map(t => t.tagName))).sort()

  let count = 0 // 解决重名问题
  for (let i = 0; i < headings.length; i++) {
    const anchor = headings[i].textContent.trim()
    // 手动构造 id, 在后面加上序号防止重名
    headings[i].id = `${anchor}-${count++}`
    anchors.value.push({
      id: headings[i].id,
      name: anchor,
      indent: hTags.indexOf(headings[i].tagName),
    })
  }
}

function handleClickAnchor(id) {
  const anchorElement = document.getElementById(id)
  // 正文被重新渲染过时可能找不到对应元素
  if (!anchorElement) {
    return
  }
  window.scrollTo({
    behavior: 'smooth',
    top: anchorElement.offsetTop - 40,
  })
  setTimeout(() => selectAnchor.value = id, 600)
}

// 实现目录高亮当前的位置的标题
// 思路: 循环的方式将标题距离顶部距离与滚动条当前位置对比, 来确定高亮的标题
const { y } = useWindowScroll()
watchThrottled(y, () => {
  anchors.value.forEach((e) => {
    const value = headings.find(ee => ee.id === e.id)
    if (value && y.value >= value.offsetTop - 50) { // 控制距离标题多远才算该标题
      selectAnchor.value = value.id
    }
  })
}, { throttle: 200 })
</script>

<template>
  <Transition name="slide-fade" appear>
    <div class="card-view space-y-2">
      <div class="flex items-center">
        <span class="i-fa-solid:list-ul" />
        <span class="ml-2">目录</span>
      </div>
      <ul>
        <li v-for="anchor of anchors" :key="anchor.id">
          <div
            class="cursor-pointer border-l-4 border-transparent rounded py-1 text-sm color-muted hover:bg-#00c4b6 hover:bg-opacity-30"
            :class="anchor.id === selectAnchor && 'bg-#00c4b6 text-white border-l-#009d92'"
            :style="{ paddingLeft: `${5 + anchor.indent * 15}px` }" @click="handleClickAnchor(anchor.id)"
          >
            {{ anchor.name }}
          </div>
        </li>
      </ul>
    </div>
  </Transition>
</template>
