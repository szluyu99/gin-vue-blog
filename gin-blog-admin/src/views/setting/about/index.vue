<script setup>
import MdEditor from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import api from '@/api'

defineOptions({ name: '关于我' })

let aboutContent = $ref('')
let btnLoading = $ref(false)

onMounted(async () => {
  const res = await api.getAbout()
  aboutContent = res.data
})

async function handleSave() {
  try {
    btnLoading = true
    await api.updateAbout({ content: aboutContent })
    $message.success('更新成功')
  }
  finally {
    btnLoading = false
  }
}
</script>

<template>
  <CommonPage :show-header="false">
    <div mb-15 flex items-center justify-between>
      <div text-23 font-bold>
        关于我
      </div>
      <n-button type="primary" :loading="btnLoading" @click="handleSave">
        <TheIcon v-if="!btnLoading" icon="line-md:confirm-circle" mr-5 :size="18" /> 保存
      </n-button>
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
