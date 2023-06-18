<script setup>
import { onMounted, ref } from 'vue'
import { NButton, NDropdown, NForm, NFormItem, NImage, NInput } from 'naive-ui'

import CrudModal from '@/components/table/CrudModal.vue'
import UploadOne from '@/components/upload/UploadOne.vue'
import CommonPage from '@/components/page/CommonPage.vue'

import { convertImgUrl, renderIcon } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

const {
  modalVisible,
  modalTitle,
  modalLoading,
  handleAdd,
  handleDelete,
  handleEdit,
  handleSave,
  modalForm,
  modalFormRef,
} = useCRUD({
  name: '页面',
  initForm: {},
  doCreate: api.saveOrUpdatePage,
  doDelete: api.deletePage,
  doUpdate: api.saveOrUpdatePage,
  refresh: fetchData,
})

// ref 语法糖使用
const pageList = ref([])
const reloadFlag = ref(false)
const uploadOneRef = ref(null) // 图片上传 ref 对象

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  const res = await api.getPages()
  pageList.value = res.data
}

// 根据输入的链接刷新预览图片
function refreshImg(img) {
  reloadFlag.value = true
  uploadOneRef.value.previewImg = img
  setTimeout(() => reloadFlag.value = false, 1000)
}

// 注意 handleSelect 第一个参数
function handleSelect(key, page) {
  if (key === 'edit')
    handleEdit(page)
  else if (key === 'delete')
    handleDelete([page.id])
}

const options = [
  {
    label: '编辑',
    key: 'edit',
    icon: renderIcon('mingcute:edit-2-line', { size: '15px' }),
  },
  {
    label: '删除',
    key: 'delete',
    icon: renderIcon('mingcute:delete-back-line', { size: '15px' }),
  },
]
</script>

<template>
  <CommonPage title="页面管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <span class="i-material-symbols:add mr-5 text-18" /> 新建页面
      </NButton>
    </template>
    <div class="flex flex-wrap justify-between">
      <div
        v-for="page of pageList" :key="page.id"
        class="relative my-10 w-300 cursor-pointer text-center"
      >
        <div class="absolute right-10 top-5 text-white">
          <NDropdown :options="options" @select="handleSelect($event, page)">
            <span class="i-ion:ellipsis-horizontal h-20 w-20 text-white hover:text-blue" />
          </NDropdown>
        </div>
        <NImage
          :src="convertImgUrl(page.cover)"
          height="170" width="300"
          :img-props="{ style: { 'border-radius': '5px' } }"
        />
        <p class="text-15">
          {{ page.name }}
        </p>
      </div>
      <div class="h-0 w-300" />
      <div class="h-0 w-300" />
      <div class="h-0 w-300" />
    </div>

    <!-- 新增/编辑/查看 弹窗 -->
    <CrudModal
      v-model:visible="modalVisible"
      width="550px"
      :title="modalTitle"
      :loading="modalLoading"
      @save="handleSave"
    >
      <NForm
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <NFormItem
          label="页面名称"
          path="name"
          :rule="{ required: true, message: '请输入页面名称', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.name" placeholder="页面名称" />
        </NFormItem>
        <NFormItem
          label="页面标签"
          path="label"
          :rule="{ required: true, message: '请输入页面标签', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.label" placeholder="页面标签" />
        </NFormItem>
        <NFormItem
          label="页面封面"
          path="cover"
          :rule="{ required: true, message: '请上传封面图片', trigger: ['input', 'blur'] }"
        >
          <div class="w-full flex items-center justify-between">
            <UploadOne
              ref="uploadOneRef"
              v-model:preview="modalForm.cover"
              :width="300"
              @finish="val => (modalForm.cover = val)"
            />

            <span
              class="i-uiw:reload h-20 w-20 cursor-pointer"
              :class="reloadFlag ? 'animate-spin' : ''"
              @click="refreshImg(modalForm.cover)"
            />
          </div>
        </NFormItem>
        <NFormItem
          label="封面链接"
          path="cover"
          :rule="{ required: true, message: '请输入封面链接', trigger: ['input', 'blur'] }"
        >
          <NInput
            v-model:value="modalForm.cover"
            type="textarea"
            placeholder="图片上传成功自动生成，或者直接复制外链"
          />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
