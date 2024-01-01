<script setup>
import { computed, ref, watch } from 'vue'
import { NIcon, NText, NUpload, NUploadDragger } from 'naive-ui'
import { useAuthStore } from '@/store'
import { convertImgUrl } from '@/utils'

const props = defineProps({
  preview: {
    type: String,
    default: '',
  },
  width: {
    type: Number,
    default: 120,
  },
})

const emit = defineEmits(['update:preview'])

const { token } = useAuthStore()
const previewImg = ref(props.preview)

watch(() => props.preview, val => previewImg.value = val)

// 上传图片
function handleImgUpload({ event }) {
  const respStr = (event?.target).response
  const res = JSON.parse(respStr)
  if (res.code !== 0) {
    $message?.error(res.message)
    return
  }
  previewImg.value = res.data
  emit('update:preview', previewImg.value)
}

// 判断是本地上传的图片或网络资源
// 开发环境可以使用本地文件上传, 生产环境建议使用云存储
const imgUrl = computed(() => convertImgUrl(previewImg.value))

defineExpose({ previewImg })
</script>

<template>
  <div>
    <NUpload
      action="/api/upload"
      :headers="{ Authorization: `Bearer ${token}` }"
      :show-file-list="false"
      @finish="handleImgUpload"
    >
      <template v-if="previewImg">
        <img
          border-color="#d9d9d9"
          class="cursor-pointer border-2 rounded-lg border-dashed hover:border-color-lightblue"
          :style="{ width: `${props.width}px` }"
          :src="imgUrl"
          alt="文章封面"
        >
      </template>
      <template v-else>
        <NUploadDragger>
          <div class="mb-3">
            <NIcon size="50" :depth="3">
              <span class="i-mdi:upload" />
            </NIcon>
          </div>
          <NText>
            点击或者拖动文件到该区域来上传
          </NText>
        </NUploadDragger>
      </template>
    </NUpload>
  </div>
</template>
