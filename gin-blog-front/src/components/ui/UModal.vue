<script setup lang="ts">
import { computed } from 'vue'

// import { useDOMScrollLock } from '../composables/useDOMScrollLock'

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
  dismissible: {
    type: Boolean,
    default: true,
  },
  dismissButton: {
    type: Boolean,
    default: true,
  },
  padded: {
    type: Boolean,
    default: true,
  },
  width: {
    type: Number,
    default: 500,
  },
  zIndex: {
    type: Number,
    default: 40,
  },
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'close'): void
}>()

const isOpen = computed({
  get: () => props.modelValue,
  set: val => emit('update:modelValue', val),
})

function close() {
  if (props.dismissible)
    isOpen.value = false
  emit('close')
}

// Lock DOM scroll when modelValue is `true`
// !️ We need to use type assertion here because of this issue: https://github.com/johnsoncodehk/volar/issues/2219
// useDOMScrollLock(toRef(props, 'modelValue') as Ref<boolean>)
</script>

<template>
  <Teleport to="body">
    <div
      class="fixed inset-0 overflow-y-auto transition-all ease-in"
      :class="[
        isOpen ? 'visible' : 'invisible duration-100 ease-in',
      ]"
      :style="{ 'z-index': zIndex }"
    >
      <!-- overlay -->
      <div
        class="fixed inset-0 bg-black transition-opacity"
        :class="[
          isOpen ? 'opacity-50 duration-75 ease-out' : 'opacity-0 duration-75 ease-in',
        ]"
        @click.self="close"
      />
      <!-- dialog -->
      <div class="min-h-full flex items-center justify-center p-3">
        <div
          v-bind="$attrs"
          class="relative inline-block w-full rounded-lg bg-white shadow-xl transition-all dark:bg-gray-900"
          :class="[
            padded ? 'p-4 lg:py-5 lg:px-7' : 'p-1',
            isOpen
              ? 'translate-y-0 opacity-100 duration-300 sm:scale-100'
              : 'translate-y-4 opacity-0 duration-300 sm:translate-y-0 sm:scale-95',
          ]"
          :style="{
            width: `${width}px`,
          }"
        >
          <button
            v-if="dismissButton"
            class="absolute right-5 top-5 h-6 w-6 rounded-full bg-gray-100 p-1 text-gray-700 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
            @click="close"
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <slot />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style lang="scss">
// html.scroll-lock {
//   position: fixed;
//   overflow-y: scroll;
//   top: var(--window-scroll-top);
//   width: 100%;
// }
</style>
