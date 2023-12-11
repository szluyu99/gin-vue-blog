<script setup>
import { computed } from 'vue'
import { NButton, NModal } from 'naive-ui'

const props = defineProps({
  visible: { type: Boolean, required: true },
  width: { type: String, default: '600px' },
  title: { type: String, default: '' },
  showFooter: { type: Boolean, default: true },
  loading: { type: Boolean, default: false },
  cancelText: { type: String, default: '取消' },
  okText: { type: String, default: '确定' },
})

const emit = defineEmits(['update:visible', 'save'])

const show = computed({
  get: () => props.visible,
  set: v => emit('update:visible', v),
})
</script>

<template>
  <NModal
    v-model:show="show"
    :style="{ width }"
    preset="card"
    :title="title"
    size="huge"
    :bordered="false"
  >
    <slot />
    <template v-if="showFooter" #footer>
      <footer class="flex justify-end space-x-5">
        <slot name="footer">
          <NButton @click="show = false">
            {{ cancelText }}
          </NButton>
          <NButton :loading="loading" type="primary" @click="emit('save')">
            {{ okText }}
          </NButton>
        </slot>
      </footer>
    </template>
  </NModal>
</template>
