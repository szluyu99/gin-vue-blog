<script setup>
import { watchDebounced } from '@vueuse/core'
import { NInput, NPopover } from 'naive-ui'
import { ref } from 'vue'

import iconData from '@/assets/icons'
import TheIcon from './TheIcon.vue'

const props = defineProps({ value: String })
const emit = defineEmits(['update:value'])

const choosed = ref(props.value) // 选中值

// 初始就把全部图标列出来。
// 原来是按当前值过滤: icons.js 里 228 项全是 mdi-xxx 短横线格式, 而菜单里存的
// 图标名是冒号格式(如 ic:outline-online-prediction), 过滤结果 0 条, 打开就是空的;
// value 为 undefined 时还会去搜字符串 "undefined", 同样 0 条
const icons = ref([...iconData])

function filterIcons() {
  const keyword = (choosed.value ?? '').trim().toLowerCase()
  icons.value = keyword
    ? iconData.filter(item => item.toLowerCase().includes(keyword))
    : [...iconData]
}

function selectIcon(icon) {
  choosed.value = icon
  emit('update:value', choosed.value)
}

watchDebounced(choosed, () => {
  filterIcons()
  emit('update:value', choosed.value)
}, { debounce: 500 })
</script>

<template>
  <div class="w-full">
    <NPopover trigger="click" placement="bottom-start">
      <template #trigger>
        <NInput
          v-model:value="choosed"
          placeholder="请输入图标名称"
          @update:value="filterIcons"
        >
          <template #prefix>
            <span class="i-mdi:magnify text-base" />
          </template>
          <template #suffix>
            <TheIcon :icon="choosed" :size="18" />
          </template>
        </NInput>
      </template>
      <template #footer>
        更多图标去
        <a class="text-blue" target="_blank" rel="noopener noreferrer" href="https://icones.js.org/collection/all">
          Icones
        </a>
        查看
      </template>
      <ul v-if="icons.length" class="h-[150px] w-[300px] overflow-y-scroll">
        <li
          v-for="(icon, index) in icons" :key="index"
          class="mx-1.5 inline-block cursor-pointer hover:text-cyan"
          @click="selectIcon(icon)"
        >
          <TheIcon :icon="icon" :size="18" />
        </li>
      </ul>
      <div v-else class="w-[300px] py-6 text-center text-gray-400">
        <TheIcon :icon="choosed" :size="18" />
        <p class="mt-1 text-xs">
          没有匹配的图标
        </p>
      </div>
    </NPopover>
  </div>
</template>
