<script setup>
import { NButton, NDropdown, NEmpty, NForm, NFormItem, NImage, NInput } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import UploadOne from '@/components//UploadOne.vue'
import CommonPage from '@/components/common/CommonPage.vue'

import CrudModal from '@/components/crud/CrudModal.vue'
import { useCRUD } from '@/composables'
import { convertImgUrl } from '@/utils'

// FIXME: 只有这个页面的 KeepAlive 为什么没有生效？

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

const pageList = ref([])
// 首屏拉取期间不显示空状态, 否则会闪一下「还没有页面」
const loading = ref(true)
const reloadFlag = ref(false)
const uploadOneRef = ref(null) // 图片上传 ref 对象

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  // 裸 await 会在接口失败时产生未捕获的 rejection; data 为空时也不能让 pageList 变成 undefined
  loading.value = true
  try {
    const resp = await api.getPages()
    pageList.value = resp.data ?? []
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
}

// 根据输入的链接刷新预览图片
function refreshImg(img) {
  // 弹窗未打开时 UploadOne 还没挂载, 直接取 .previewImg 会抛异常
  if (!uploadOneRef.value) {
    return
  }
  reloadFlag.value = true
  uploadOneRef.value.previewImg = img
  setTimeout(() => reloadFlag.value = false, 600)
}

function handleSelect(key, page) {
  if (key === 'edit') {
    handleEdit(page)
  }
  else if (key === 'delete') {
    handleDelete([page.id])
  }
}

const options = [
  {
    label: '编辑',
    key: 'edit',
    icon: () => h('i', { class: 'i-mingcute:edit-2-line' }),
  },
  {
    label: '删除',
    key: 'delete',
    icon: () => h('i', { class: 'i-mingcute:delete-back-line' }),
  },
]
</script>

<template>
  <CommonPage title="页面管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        新建页面
      </NButton>
    </template>
    <!-- 原来是 flex + justify-between, 末尾还要垫三个空 div 才对得齐;
         换成 grid, 列数按容器宽度自适应, 卡片数变化也不会跳 -->
    <div class="grid gap-4" style="grid-template-columns: repeat(auto-fill, minmax(300px, 1fr))">
      <div
        v-for="page of pageList" :key="page.id"
        class="relative my-2 cursor-pointer text-center"
      >
        <div class="absolute right-2 top-1 text-white">
          <NDropdown :options="options" @select="handleSelect($event, page)">
            <span class="i-ion:ellipsis-horizontal h-5 w-5 text-white hover:text-blue" />
          </NDropdown>
        </div>
        <NImage
          :src="convertImgUrl(page.cover)"
          height="170" width="300"
          :img-props="{ style: { 'border-radius': '5px' } }"
        />
        <p class="text-base">
          {{ page.name }}
        </p>
      </div>
    </div>

    <NEmpty v-if="!loading && !pageList.length" class="py-16" description="还没有页面, 点右上角新建" />

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
              class="i-uiw:reload h-5 w-5 cursor-pointer"
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
