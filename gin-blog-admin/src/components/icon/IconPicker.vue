<script setup>
import { watchDebounced } from '@vueuse/core'
import iconData from '@/assets/js/icons'

const props = defineProps({ value: String })
const emit = defineEmits(['update:value'])

let choosed = $ref(props.value) // 选中值
let icons = $ref(iconData.filter(icon => icon.includes(choosed))) // 可选图标列表

function filterIcons() {
  icons = iconData.filter(item => item.includes(choosed))
}

function selectIcon(icon) {
  choosed = icon
  emit('update:value', choosed)
}

watchDebounced($$(choosed), () => {
  filterIcons()
  emit('update:value', choosed)
}, { debounce: 500 })
</script>

<template>
  <div>
    <n-popover trigger="click" placement="bottom-start">
      <template #trigger>
        <n-input
          v-model:value="choosed"
          placeholder="请输入图标名称"
          @update:value="filterIcons"
        >
          <template #prefix>
            <icon-mdi:magnify />
          </template>
          <template #suffix>
            <TheIcon :icon="choosed" :size="18" />
          </template>
        </n-input>
      </template>
      <template #footer>
        更多图标去 <a color-blue target="_blank" href="https://icones.js.org/collection/all">Icones</a> 查看
      </template>
      <ul v-if="icons.length" w-300 h-150 overflow-y-scroll>
        <li
          v-for="(icon, index) in icons" :key="index"
          mx-5 cursor-pointer hover:text-cyan inline-block
          @click="selectIcon(icon)"
        >
          <TheIcon :icon="icon" :size="18" />
        </li>
      </ul>
      <div v-else>
        <TheIcon :icon="choosed" :size="18" />
      </div>
    </n-popover>
  </div>
</template>
