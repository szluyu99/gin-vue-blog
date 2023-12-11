<script setup>
import { onMounted, ref } from 'vue'
import { NButton } from 'naive-ui'

import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

import CommonPage from '@/components/common/CommonPage.vue'
import api from '@/api'

defineOptions({ name: '关于我' })

const aboutContent = ref('')
const btnLoading = ref(false)

onMounted(async () => {
  const resp = await api.getAbout()
  aboutContent.value = resp.data
})

async function handleSave() {
  try {
    btnLoading.value = true
    await api.updateAbout({ content: aboutContent.value })
    window.$message.success('更新成功')
  }
  finally {
    btnLoading.value = false
  }
}
</script>

<template>
  <CommonPage :show-header="false">
    <div class="mb-4 flex items-center justify-between">
      <div class="mx-1 text-2xl font-bold">
        关于我
      </div>
      <NButton type="primary" :loading="btnLoading" @click="handleSave">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:confirm-circle" />
        </template>
        保存
      </NButton>
    </div>
    <!-- TODO: 文件上传封装 -->
    <MdEditor v-model="aboutContent" style="height: calc(100vh - 245px)" />
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
