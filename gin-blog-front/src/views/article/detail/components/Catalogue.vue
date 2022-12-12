<script setup lang="ts">
interface Props { preview: any }
// 文章 ref 对象
const { preview } = defineProps<Props>()

onMounted(() => {
  preview && buildAnchorTitles()
})

// 生成锚点目录
let titles = $ref<any>([])
function buildAnchorTitles() {
  const anchors = preview.$el.querySelectorAll('h1,h2,h3,h4,h5,h6')
  const titleList = Array.from(anchors).filter((t: any) => !!t.innerText.trim())
  if (!titleList.length)
    return
  const hTags = Array.from(new Set(titleList.map((t: any) => t.tagName))).sort()
  titles = titleList.map((el: any) => ({
    title: el.innerText,
    lineIndex: el.getAttribute('data-v-md-line'),
    indent: hTags.indexOf(el.tagName),
  }))
}

// 点击锚点目录
function handleAnchorClick(anchor: any) {
  const { lineIndex } = anchor
  const heading = preview.$el.querySelector(`[data-v-md-line="${lineIndex}"]`)
  if (heading) {
    preview.scrollToTarget({
      target: heading,
      scrollContainer: window,
      top: 40,
    })
  }
}

// TODO: 实现目录高亮当前位置的标题
</script>

<template>
  <n-card hoverable rounded-2rem mb-15>
    <div flex items-center mb-10 text-18>
      <i-fa-solid:list-ul />
      <span ml-10 font-bold>目录</span>
    </div>
    <div
      v-for="anchor in titles" :key="anchor.title"
      py-3 rounded-1rem hover:bg="#00c4b6" hover:text-white
      :style="{ paddingLeft: `${5 + anchor.indent * 15}px` }"
      @click="handleAnchorClick(anchor)"
    >
      <button>{{ anchor.title }}</button>
    </div>
  </n-card>
</template>
