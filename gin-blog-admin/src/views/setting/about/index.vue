<script setup>
import { onMounted, ref } from 'vue'
import { NButton } from 'naive-ui'
import MdEditor from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

import CommonPage from '@/components/page/CommonPage.vue'
import api from '@/api'

defineOptions({ name: '关于我' })

const aboutContent = ref('')
const btnLoading = ref(false)

onMounted(async () => {
  const res = await api.getAbout()
  aboutContent.value = res.data
})

async function handleSave() {
  try {
    btnLoading.value = true
    await api.updateAbout({ content: aboutContent.value })
    $message.success('更新成功')
  }
  finally {
    btnLoading.value = false
  }
}
</script>

<template>
  <CommonPage :show-header="false">
    <div class="mb-15 flex items-center justify-between">
      <div class="text-23 font-bold">
        关于我
      </div>
      <NButton type="primary" :loading="btnLoading" @click="handleSave">
        <span v-if="!btnLoading" class="i-line-md:confirm-circle mr-5 text-18" /> 保存
      </NButton>
    </div>
    <MdEditor v-model="aboutContent" style="height: calc(100vh - 305px)" />
  </CommonPage>
</template>

<style lang="scss">
.md-preview {
  ul,
  ol {
    list-style: revert;
  }
}
</style>
