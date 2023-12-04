<script setup>
import { defineComponent, h } from 'vue'
import { NConfigProvider, NLoadingBarProvider, NMessageProvider, NNotificationProvider, useLoadingBar, useMessage, useNotification } from 'naive-ui'

// 挂载 naive 组件的方法至 window, 方便全局使用
function setupNaiveTools() {
  window.$loadingBar = useLoadingBar()
  window.$message = useMessage()
  window.$notification = useNotification()
}

const NaiveProviderContent = defineComponent({
  setup() {
    setupNaiveTools()
  },
  render() {
    return h('div')
  },
})
</script>

<template>
  <NConfigProvider class="h-full w-full">
    <NLoadingBarProvider>
      <NNotificationProvider>
        <NMessageProvider>
          <slot />
          <NaiveProviderContent />
        </NMessageProvider>
      </NNotificationProvider>
    </NLoadingBarProvider>
  </NConfigProvider>
</template>
