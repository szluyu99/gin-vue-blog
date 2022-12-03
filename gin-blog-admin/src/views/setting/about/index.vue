<template>
  <CommonPage :show-header="false">
    <div mb-15 flex items-center justify-between>
      <div text-23 font-bold>关于我</div>
      <n-button type="primary" :loading="btnLoading" @click="handleSave">
        <TheIcon v-if="!btnLoading" icon="line-md:confirm-circle" class="mr-5" :size="18" /> 保存
      </n-button>
    </div>
    <MdEditor v-model="aboutContent" style="height: calc(100vh - 305px)" />
  </CommonPage>
</template>

<script setup>
import MdEditor from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

import { useSettingAPI } from '@/api'
const { getAbout, updateAbout } = useSettingAPI()

defineOptions({ name: '关于我' })

onMounted(() => {
  getAbout().then((res) => (aboutContent.value = res.data))
})

const aboutContent = ref('')
const btnLoading = ref(false)

async function handleSave() {
  try {
    btnLoading.value = true
    await updateAbout({ content: aboutContent.value })
    $message.success('更新成功')
  } finally {
    btnLoading.value = false
  }
}
</script>

<style lang="scss">
.md-preview {
  ul,
  ol {
    list-style: revert;
  }
}
</style>
