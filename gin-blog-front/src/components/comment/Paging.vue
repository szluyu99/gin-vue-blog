<script setup>
import { computed } from 'vue'

// 受控组件: 当前页数由父组件持有
// 以前 current 和 show 都是组件内部状态, 父组件只能通过 v-for 的模板 ref 数组
// 按下标调 setShow / 读 current, 而 Vue 不保证 ref 数组顺序与源数组一致
const { pageTotal = 0, current = 1 } = defineProps({
  pageTotal: Number, // 总页数
  current: Number, // 当前页数
})

const emit = defineEmits(['changeCurrent'])

// 计算属性传参: 判断当前页数是否激活
const isActive = computed(() => i => (i === current))
</script>

<template>
  <!-- TODO: 优化显示 -->
  <div class="mb-2.5 mt-4 text-[12px]">
    <span class="mr-10"> 共 {{ pageTotal }} 页 </span>
    <!-- 上一页按钮: 第一页不显示 -->
    <a v-show="current !== 1" @click="emit('changeCurrent', current - 1)"> 上一页 </a>

    <!-- 总页数小于 6 页直接显示 -->
    <template v-if="pageTotal < 6">
      <a v-for="i of pageTotal" :key="i" class="mx-1" @click="emit('changeCurrent', i)">
        <!-- 当前选中页数 -->
        <span v-if="isActive(i)" class="color-#00a1d6 font-bold">
          {{ i }}
        </span>
        <span v-else> {{ i }} </span>
      </a>
    </template>

    <!-- TODO: 其他情况 -->

    <!-- 下一页按钮: 最后一页不显示 -->
    <a v-show="current !== pageTotal" @click="emit('changeCurrent', current + 1)"> 下一页 </a>
  </div>
</template>

<style scoped>
.active {
  color: #00a1d6;
  font-weight: bold;
}
</style>
