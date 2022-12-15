<script setup>
import UploadOne from '@/components/upload/UploadOne.vue'
import { convertImgUrl, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
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
let pageList = $ref([])
let reloadFlag = $ref(false)
const uploadOneRef = $ref(null) // 图片上传 ref 对象

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  const res = await api.getPages()
  pageList = res.data
}

// 根据输入的链接刷新预览图片
function refreshImg(img) {
  reloadFlag = true
  uploadOneRef.previewImg = img
  setTimeout(() => reloadFlag = false, 1000)
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
  <CommonPage show-footer title="页面管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <TheIcon icon="material-symbols:add" :size="18" /> 新建页面
      </NButton>
    </template>
    <div flex flex-wrap justify-between>
      <div
        v-for="page of pageList" :key="page.id"
        relative w-300 my-10 text-center cursor-pointer
      >
        <div absolute top-5 right-10 text-white>
          <n-dropdown :options="options" @select="handleSelect($event, page)">
            <icon-ion:ellipsis-horizontal text-white text-20 hover:text-blue />
          </n-dropdown>
        </div>
        <n-image
          :src="convertImgUrl(page.cover)"
          height="170" width="300"
          :img-props="{ style: { 'border-radius': '5px' } }"
        />
        <p text-15>
          {{ page.name }}
        </p>
      </div>
      <div w-300 h-0 />
      <div w-300 h-0 />
      <div w-300 h-0 />
    </div>

    <!-- 新增/编辑/查看 弹窗 -->
    <CrudModal
      v-model:visible="modalVisible"
      width="550px"
      :title="modalTitle"
      :loading="modalLoading"
      @on-save="handleSave"
    >
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <n-form-item
          label="页面名称"
          path="name"
          :rule="{ required: true, message: '请输入页面名称', trigger: ['input', 'blur'] }"
        >
          <n-input v-model:value="modalForm.name" placeholder="页面名称" />
        </n-form-item>
        <n-form-item
          label="页面标签"
          path="label"
          :rule="{ required: true, message: '请输入页面标签', trigger: ['input', 'blur'] }"
        >
          <n-input v-model:value="modalForm.label" placeholder="页面标签" />
        </n-form-item>
        <n-form-item
          label="页面封面"
          path="cover"
          :rule="{ required: true, message: '请上传封面图片', trigger: ['input', 'blur'] }"
        >
          <div flex items-center justify-between w-full>
            <UploadOne
              ref="uploadOneRef"
              v-model:preview="modalForm.cover"
              :width="300"
              @finish="val => (modalForm.cover = val)"
            />
            <icon-uiw:reload
              :class="reloadFlag ? 'animate-spin' : ''"
              text-20 cursor-pointer
              @click="refreshImg(modalForm.cover)"
            />
          </div>
        </n-form-item>
        <n-form-item
          label="封面链接"
          path="cover"
          :rule="{ required: true, message: '请输入封面链接', trigger: ['input', 'blur'] }"
        >
          <n-input
            v-model:value="modalForm.cover"
            type="textarea"
            placeholder="图片上传成功自动生成，或者直接复制外链"
          />
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>
