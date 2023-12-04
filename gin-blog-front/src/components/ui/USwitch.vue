<script setup>
import { computed } from 'vue'
import uniqueId from 'lodash-es/uniqueId'

const props = defineProps({
  modelValue: {
    type: [Boolean, String, Number],
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  checked: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  icon: {
    type: Boolean,
    default: false,
  },
  readonly: {
    type: Boolean,
    default: false,
  },
  checkedValue: {
    type: [Boolean, String, Number],
    default: true,
  },
  uncheckedValue: {
    type: [Boolean, String, Number],
    default: false,
  },
  color: {
    type: String,
    default: 'green',
  },
})

const emits = defineEmits(['update:modelValue'])

const uid = uniqueId('switch-')

const switchValue = computed({
  get: () => {
    if (props.checkedValue && props.modelValue === props.checkedValue)
      return true
    if (props.uncheckedValue && props.modelValue === props.uncheckedValue)
      return false
    return props.modelValue
  },
  set: (val) => {
    if (props.checkedValue && val)
      val = props.checkedValue
    else if (props.uncheckedValue && !val)
      val = props.uncheckedValue
    emits('update:modelValue', val)
  },
})

const styleObj = computed(() => {
  const style = {}
  if (switchValue.value || props.checked) {
    style.backgroundColor = props.color
  } else {
    style.backgroundColor = '#999999'
  }
  return style
})
</script>

<template>
  <div class="flex items-center">
    <label
      :for="uid"
      class="flex select-none items-center"
      :class="[disabled ? 'cursor-not-allowed opacity-60' : readonly ? '' : 'cursor-pointer']"
    >
      <div class="mr-2 font-medium leading-6 text-gray-900 empty:hidden">
        <slot />
      </div>
      <!-- focus:rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 -->
      <div tabindex="0" class="relative">
        <input
          :id="uid"
          v-model="switchValue"
          v-bind="$attrs"
          type="checkbox"
          :disabled="disabled"
          class="sr-only"
        >
        <div class="block h-20 w-40 rounded-full" :style="styleObj" />
        <div
          class="absolute left-0 top-0.5 h-20 w-20 rounded-full bg-white transition"
          :class="[
            (switchValue || checked) ? 'translate-x-25' : 'translate-x-0',
          ]"
        >
          <svg
            v-if="loading"
            class="h-20 w-20 animate-spin"
            xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24"
          ><defs><linearGradient id="mingcuteLoadingLine0" x1="50%" x2="50%" y1="5.271%" y2="91.793%"><stop offset="0%" stop-color="#888888" /><stop offset="100%" stop-color="currentColor" stop-opacity=".55" /></linearGradient><linearGradient id="mingcuteLoadingLine1" x1="50%" x2="50%" y1="8.877%" y2="90.415%"><stop offset="0%" stop-color="currentColor" stop-opacity="0" /><stop offset="100%" stop-color="currentColor" stop-opacity=".55" /></linearGradient></defs><g fill="none"><path d="M24 0v24H0V0h24ZM12.593 23.258l-.011.002l-.071.035l-.02.004l-.014-.004l-.071-.035c-.01-.004-.019-.001-.024.005l-.004.01l-.017.428l.005.02l.01.013l.104.074l.015.004l.012-.004l.104-.074l.012-.016l.004-.017l-.017-.427c-.002-.01-.009-.017-.017-.018Zm.265-.113l-.013.002l-.185.093l-.01.01l-.003.011l.018.43l.005.012l.008.007l.201.093c.012.004.023 0 .029-.008l.004-.014l-.034-.614c-.003-.012-.01-.02-.02-.022Zm-.715.002a.023.023 0 0 0-.027.006l-.006.014l-.034.614c0 .012.007.02.017.024l.015-.002l.201-.093l.01-.008l.004-.011l.017-.43l-.003-.012l-.01-.01l-.184-.092Z" /><path fill="url(#mingcuteLoadingLine0)" d="M8.886.006a1 1 0 0 1 .22 1.988A8.001 8.001 0 0 0 10 17.944v2c-5.523 0-10-4.476-10-10C0 4.838 3.848.566 8.886.007Z" transform="translate(2 2.055)" /><path fill="url(#mingcuteLoadingLine1)" d="M14.322 1.985a1 1 0 0 1 1.392-.248A9.988 9.988 0 0 1 20 9.945c0 5.523-4.477 10-10 10v-2a8 8 0 0 0 4.57-14.567a1 1 0 0 1-.248-1.393Z" transform="translate(2 2.055)" /></g></svg>
          <!-- TODO:  -->
          <svg
            v-else-if="icon"
            class="h-20 w-20"
            :class="[
              (checked || switchValue) ? 'text-red-800' : 'text-gray-200',
            ]"
            fill="currentColor" viewBox="0 0 12 12"
          >
            <path v-if="checked || switchValue" d="M3.707 5.293a1 1 0 00-1.414 1.414l1.414-1.414zM5 8l-.707.707a1 1 0 001.414 0L5 8zm4.707-3.293a1 1 0 00-1.414-1.414l1.414 1.414zm-7.414 2l2 2 1.414-1.414-2-2-1.414 1.414zm3.414 2l4-4-1.414-1.414-4 4 1.414 1.414z" />
            <path v-else d="M4 8l2-2m0 0l2-2M6 6L4 4m2 2l2 2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </div>
      </div>
    </label>
  </div>
</template>
