<script setup>
import { computed, ref } from 'vue'
import { convertImgUrl } from '@/utils'
import { useUserStore } from '@/store'

const props = defineProps({
  preview: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:preview'])
const previewImg = ref(props.preview) // 图片预览

// 判断是本地上传的图片或网络资源
// 开发环境可以使用本地文件上传, 生产环境建议使用云存储
const imgUrl = computed(() => convertImgUrl(previewImg.value))

const fileRef = ref(null)

async function handleFileChange() {
  const file = fileRef.value.files[0]
  const formData = new FormData()
  formData.append('file', file)
  try {
    const { token } = useUserStore()
    const response = await fetch('/api/front/upload', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    })

    const responseJSON = await response.json()
    if (responseJSON.code !== 0) {
      window.$message?.error(responseJSON.message)
      return
    }

    previewImg.value = responseJSON.data
    emit('update:preview', previewImg)
  }
  catch (err) {
    console.error(err)
    window.$message?.error('文件上传失败')
  }
}
</script>

<template>
  <!-- TODO: 拖拽文件上传 -->
  <main class="flex items-center justify-center bg-gray-100 font-sans">
    <label for="dropzone-file" class="mx-auto max-w-[300px] w-full cursor-pointer items-center border-1 border-blue-400 rounded-xl border-dashed bg-white p-2 text-center">
      <template v-if="previewImg">
        <div class="group relative">
          <img class="lg:h-[160px] lg:w-[160px]" :src="imgUrl" alt="user avatar">
          <div class="absolute bottom-0 left-0 right-0 top-0 f-c-c cursor-pointer">
            <button class="i-mdi:upload pointer-events-none inline-block text-[50px] text-white opacity-35 duration-200 group-hover:opacity-80" />
          </div>
        </div>
      </template>
      <template v-else>
        <div class="f-c-c lg:h-[160px] lg:w-[160px]">
          <div class="flex flex-col items-center">
            <span class="i-mdi:upload text-[58px] text-blue-500" />
            <span class="text-blue-400"> 点击上传文件</span>
          </div>
        </div>
      </template>
      <input id="dropzone-file" ref="fileRef" type="file" class="hidden" @change="handleFileChange">
    </label>
  </main>
</template>
