<script setup>
import { computed, ref } from 'vue'

const { pageTotal = 0, index, commentId } = defineProps({
  pageTotal: Number, // 总页数
  index: Number, // 索引
  commentId: Number, // 评论 id
})

const emit = defineEmits(['changeCurrent'])

const show = ref(false) // 是否显示分页

// 当前页数
const current = ref(1)
// 计算属性传参: 判断当前页数是否激活
const isActive = computed(() => i => (i === current.value))
// 上一页
function prePage() {
  current.value -= 1
  emit('changeCurrent', current.value, index, commentId)
}
// 下一页
function nextPage() {
  current.value += 1
  emit('changeCurrent', current.value, index, commentId)
}
// 跳转指定页数
function changeCurrent(i) {
  current.value = i
  emit('changeCurrent', current.value, index, commentId)
}

function setShow(flag) {
  show.value = flag
}

defineExpose({
  current,
  setShow,
})
</script>

<template>
  <!-- TODO: 优化显示 -->
  <div v-show="show" class="mb-2.5 mt-4 text-[12px]">
    <span class="mr-10"> 共 {{ pageTotal }} 页 </span>
    <!-- 上一页按钮: 第一页不显示 -->
    <a v-show="current !== 1" @click="prePage"> 上一页 </a>

    <!-- 总页数小于 6 页直接显示 -->
    <template v-if="pageTotal < 6">
      <a v-for="i of pageTotal" :key="i" class="mx-1" @click="changeCurrent(i)">
        <!-- 当前选中页数 -->
        <span v-if="isActive(i)" class="color-#00a1d6 font-bold">
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

<style scoped>
.active {
  color: #00a1d6;
  font-weight: bold;
}
</style>
