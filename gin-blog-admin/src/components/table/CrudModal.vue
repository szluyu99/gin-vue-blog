<script setup>
import { computed } from 'vue'
import { NButton, NModal } from 'naive-ui'

const props = defineProps({
  visible: { type: Boolean, required: true },
  width: { type: String, default: '600px' },
  title: { type: String, default: '' },
  showFooter: { type: Boolean, default: true },
  loading: { type: Boolean, default: false },
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
      <footer class="flex justify-end space-x-20">
        <slot name="footer">
          <NButton @click="show = false">
            取消
          </NButton>
          <NButton
            :loading="loading"
            type="primary"
            @click="emit('save')"
          >
            保存
          </NButton>
        </slot>
      </footer>
    </template>
  </NModal>
</template>
