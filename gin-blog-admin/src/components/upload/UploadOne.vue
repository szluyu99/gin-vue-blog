<script setup>
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

const emit = defineEmits(['finish'])

const previewImg = ref(props.preview)

// 上传图片
function handleImgUpload({ event }) {
  const respStr = (event?.target).response
  const res = JSON.parse(respStr)
  if (res.code !== 0) {
    $message.error('文件上传失败')
    return
  }
  previewImg.value = res.data
  emit('finish', previewImg.value)
}

// 判断是本地上传的图片或网络资源
// 开发环境可以使用本地文件上传, 生产环境建议使用云存储
const imgUrl = computed(() => {
  return convertImgUrl(previewImg.value)
})

defineExpose({ previewImg })
</script>

<template>
  <div>
    <n-upload
      action="/api/upload"
      @finish="handleImgUpload"
    >
      <template v-if="previewImg">
        <img
          border-dashed border-2 border-color="#d9d9d9" rounded-2rem
          cursor-pointer hover:border-color-lightblue
          :style="{ width: `${props.width}px` }"
          :src="imgUrl" alt="文章封面"
        >
      </template>
      <template v-else>
        <n-upload-dragger>
          <div mb-12>
            <n-icon size="58" :depth="3">
              <icon-mdi:upload />
            </n-icon>
          </div>
          <n-text text-14>
            点击或者拖动文件到该区域来上传
          </n-text>
        </n-upload-dragger>
      </template>
    </n-upload>
  </div>
</template>
