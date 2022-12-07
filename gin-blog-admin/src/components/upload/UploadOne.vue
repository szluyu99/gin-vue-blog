<script setup>
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

// 开发环境可以使用本地文件上传, 生产环境建议使用云存储
const SERVER_URL = 'http://localhost:8765/'

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

// 判断图片是本地上传或网络资源, 本地上传一般用于开发时测试, 生成环境基本都是用网络资源
// function getImgUrl(img) {
//   return img.startsWith('http') ? img : `${SERVER_URL}${img}`
// }
const imgUrl = computed(() => {
  if (previewImg.value.startsWith('http'))
    return previewImg.value
  return `${SERVER_URL}${previewImg.value}`
})
</script>

<template>
  <n-upload
    action="/api/upload"
    @finish="handleImgUpload"
  >
    <template v-if="previewImg">
      <img
        border-dashed border-2 border-color="#d9d9d9" rounded-2rem cursor-pointer p-3
        hover:border-color-lightblue
        :style="{ width: `${props.width}px` }"
        :src="imgUrl" alt="文章封面"
      >
    </template>
    <template v-else>
      <n-upload-dragger>
        <div mb-12>
          <n-icon size="48" :depth="3">
            <icon-mdi:upload />
          </n-icon>
        </div>
        <n-text style="font-size: 14px">
          点击或者拖动文件到该区域来上传
        </n-text>
      </n-upload-dragger>
    </template>
  </n-upload>
</template>
