<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'

defineProps({
  showArrow: {
    type: Boolean,
    default: true,
  },
})

const translateX = ref(0)
const content = ref(null)
const wrapper = ref(null)
const isOverflow = ref(false)

const resetTranslateX = debounce((wrapperWidth, contentWidth) => {
  if (!isOverflow.value)
    translateX.value = 0
  else if (-translateX.value > contentWidth - wrapperWidth)
    translateX.value = wrapperWidth - contentWidth
  else if (translateX.value > 0)
    translateX.value = 0
}, 200)

const refreshIsOverflow = debounce(() => {
  const wrapperWidth = wrapper.value?.offsetWidth
  const contentWidth = content.value?.offsetWidth
  isOverflow.value = contentWidth > wrapperWidth
  resetTranslateX(wrapperWidth, contentWidth)
}, 200)

function handleMouseWheel(e) {
  const { wheelDelta } = e
  const wrapperWidth = wrapper.value?.offsetWidth
  const contentWidth = content.value?.offsetWidth
  /**
   * @wheelDelta 平行滚动的值 >0： 右移  <0: 左移
   * @translateX 内容translateX的值
   * @wrapperWidth 容器的宽度
   * @contentWidth 内容的宽度
   */
  if (wheelDelta < 0) {
    if (wrapperWidth > contentWidth && translateX.value < -10)
      return
    if (wrapperWidth <= contentWidth && contentWidth + translateX.value - wrapperWidth < -10)
      return
  }
  if (wheelDelta > 0 && translateX.value > 10)
    return

  translateX.value += wheelDelta
  resetTranslateX(wrapperWidth, contentWidth)
}

const observer = new MutationObserver(refreshIsOverflow)
onMounted(() => {
  refreshIsOverflow()

  window.addEventListener('resize', refreshIsOverflow)
  // 监听内容宽度刷新是否超出
  observer.observe(content.value, { childList: true })
})
onBeforeUnmount(() => {
  window.removeEventListener('resize', refreshIsOverflow)
  observer.disconnect()
})

function handleScroll(x, width) {
  const wrapperWidth = wrapper.value?.offsetWidth
  const contentWidth = content.value?.offsetWidth
  if (contentWidth <= wrapperWidth)
    return

  // 当 x 小于可视范围的最小值时
  if (x < -translateX.value + 150) {
    translateX.value = -(x - 150)
    resetTranslateX(wrapperWidth, contentWidth)
  }

  // 当 x 大于可视范围的最大值时
  if (x + width > -translateX.value + wrapperWidth) {
    translateX.value = wrapperWidth - (x + width)
    resetTranslateX(wrapperWidth, contentWidth)
  }
}

/**
 * @desc  函数防抖
 * @param {Function} method
 * @param {number} wait
 * @param {boolean} immediate
 * @return {*}
 */
function debounce(method, wait, immediate) {
  let timeout
  return function (...args) {
    // eslint-disable-next-line ts/no-this-alias
    const context = this
    if (timeout) {
      clearTimeout(timeout)
    }

    // 立即执行需要两个条件，一是 immediate 为 true，二是 timeout 未被赋值或被置为 null
    if (immediate) {
      /**
       * 如果定时器不存在，则立即执行，并设置一个定时器，wait毫秒后将定时器置为null
       * 这样确保立即执行后wait毫秒内不会被再次触发
       */
      const callNow = !timeout
      timeout = setTimeout(() => {
        timeout = null
      }, wait)
      if (callNow)
        method.apply(context, args)
    }
    else {
      // 如果immediate为false，则函数wait毫秒后执行
      timeout = setTimeout(() => {
        /**
         * args 是一个类数组对象，所以使用 fn.apply
         * 也可写作 method.call(context, ...args)
         */
        method.apply(context, args)
      }, wait)
    }
  }
}

defineExpose({
  handleScroll,
})
</script>

<template>
  <div ref="wrapper" class="wrapper" @mousewheel.prevent="handleMouseWheel">
    <template v-if="showArrow && isOverflow">
      <div class="left" @click="handleMouseWheel({ wheelDelta: 120 })">
        <span class="i-ic:baseline-keyboard-arrow-left" />
      </div>
      <div class="right" @click="handleMouseWheel({ wheelDelta: -120 })">
        <span class="i-ic:baseline-keyboard-arrow-right" />
      </div>
    </template>

    <div
      ref="content"
      class="content"
      :class="{ overflow: isOverflow && showArrow }"
      :style="{
        transform: `translateX(${translateX}px)`,
      }"
    >
      <slot />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.wrapper {
  display: flex;
  background-color: #fff;

  z-index: 9;
  overflow: hidden;
  position: relative;
  .content {
    padding: 0 10px;
    display: flex;
    align-items: center;
    flex-wrap: nowrap;
    transition: transform 0.5s;
    &.overflow {
      padding-left: 30px;
      padding-right: 30px;
    }
  }
  .left,
  .right {
    background-color: #fff;
    position: absolute;
    top: 0;
    bottom: 0;
    margin: auto;

    width: 20px;
    height: 35px;
    display: flex;
    align-items: center;
    justify-content: center;

    font-size: 18px;
    border: 1px solid #e0e0e6;
    border-radius: 2px;

    z-index: 2;
    cursor: pointer;
  }
  .left {
    left: 0;
  }
  .right {
    right: 0;
  }
}
</style>
