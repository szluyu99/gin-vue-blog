<script setup lang="ts">
const { pageTotal = 0, index, commentId }
= defineProps<{
  pageTotal: number // 总页数
  index: number // 索引
  commentId: number// 评论 id
}>()

const emit = defineEmits(['changeCurrent'])

let show = $ref(false) // 是否显示分页

// 当前页数
let current = $ref(1)
// 计算属性传参: 判断当前页数是否激活
const isActive = computed(() => (i: number) => (i === current))
// 上一页
function prePage() {
  current -= 1
  emit('changeCurrent', current, index, commentId)
}
// 下一页
function nextPage() {
  current += 1
  emit('changeCurrent', current, index, commentId)
}
// 跳转指定页数
function changeCurrent(i: number) {
  current = i
  emit('changeCurrent', current, index, commentId)
}

function setShow(flag: boolean) {
  show = flag
}

defineExpose({ current, setShow })
</script>

<template>
  <!-- TODO: 优化显示 -->
  <div v-show="show" mt-15 mb-10 text-12>
    <span mr-10> 共 {{ pageTotal }} 页 </span>
    <!-- 上一页按钮: 第一页不显示 -->
    <a v-show="current !== 1" @click="prePage"> 上一页 </a>

    <!-- 总页数小于 6 页直接显示 -->
    <template v-if="pageTotal < 6">
      <a
        v-for="i of pageTotal" :key="i" mx-5
        @click="changeCurrent(i)"
      >
        <!-- 当前选中页数 -->
        <span v-if="isActive(i)" color="#00a1d6" font-bold>
          {{ i }}
        </span>
        <span v-else> {{ i }} </span>
      </a>
    </template>

    <!-- TODO: 其他情况 -->

    <!-- 下一页按钮: 最后一页不显示 -->
    <a v-show="current !== pageTotal" @click="nextPage"> 下一页 </a>
  </div>
</template>

<style lang="scss" scoped>
.active {
  color: #00a1d6;
  font-weight: bold;
}
</style>
