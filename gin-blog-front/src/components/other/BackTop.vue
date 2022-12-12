<script setup lang="ts">
import { useWindowScroll } from '@vueuse/core'

let isShowStyle = $ref('')

const { y } = $(useWindowScroll())
watch($$(y), () => {
  isShowStyle = (y > 20) ? 'opacity: 1; transform: translateX(-40px);' : ''
})

function switchDark() {
  window.$message?.info('黑夜模式开发中...')
}

function clickSetting() {
  window.$message?.info('设置开发中...')
}

function backTop() {
  window.scrollTo({
    behavior: 'smooth',
    top: 0,
  })
}
</script>

<template>
  <div fixed z-4 right="-35" bottom-85 text-20 :style="isShowStyle" text-white transition-600>
    <div w-30 h-30 f-c-c cursor-pointer bg="#49b1f5" rounded>
      <i-bi:moon-stars-fill @click="switchDark" />
    </div>
    <div w-30 h-30 f-c-c cursor-pointer bg="#49b1f5" rounded my-3>
      <i-uiw:setting @click="clickSetting" />
    </div>
    <div w-30 h-30 f-c-c cursor-pointer bg="#49b1f5" rounded>
      <i-fluent:arrow-up-12-filled @click="backTop" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
@keyframes rightSideOut {
  0% {
    transform: translate(0, 0);
  }
  100% {
    transform: translate(30px, 0);
  }
}
@keyframes rightSideIn {
  0% {
    transform: translate(30px, 0);
  }
  100% {
    transform: translate(0, 0);
  }
}
</style>
