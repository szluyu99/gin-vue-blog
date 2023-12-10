<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NPopconfirm } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: '分类管理' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

onMounted(() => {
  $table.value?.handleSearch()
})

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
  name: '分类',
  initForm: {},
  doCreate: api.saveOrUpdateCategory,
  doDelete: api.deleteCategory,
  doUpdate: api.saveOrUpdateCategory,
  refresh: () => $table.value?.handleSearch(),
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  { title: '分类名', key: 'name', width: 100, align: 'center', ellipsis: { tooltip: true } },
  { title: '文章量', key: 'article_count', width: 30, align: 'center' },
  {
    title: '创建日期',
    key: 'created_at',
    width: 80,
    align: 'center',
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.created_at),
          icon: () => h('i', { class: 'i-mdi:clock-time-three-outline' }),
        },
      )
    },
  },
  {
    title: '更新日期',
    key: 'updated_at',
    width: 80,
    align: 'center',
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.updated_at),
          icon: () => h('i', { class: 'i-mdi:update' }),
        },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          { size: 'small', type: 'primary', onClick: () => handleEdit(row) },
          { default: () => '编辑', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () => h(
              NButton,
              { size: 'small', type: 'error', style: 'margin-left: 15px;' },
              { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
            ),
            default: () => h('div', {}, '确定删除该分类吗?'),
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <CommonPage title="分类管理">
    <template #action>
      <NButton type="primary" secondary @click="$table?.handleExport()">
        <template #icon>
          <i class="i-mdi:download" />
        </template>
        导出
      </NButton>
      <NButton type="primary" @click="handleAdd">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        新建分类
      </NButton>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <span class="i-material-symbols:playlist-remove" />
        </template>
        批量删除
      </NButton>
    </template>
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getCategorys"
    >
      <template #queryBar>
        <QueryItem label="分类名" :label-width="50">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入分类名"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>

    <CrudModal
      v-model:visible="modalVisible"
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
          label="文章分类"
          path="name"
          :rule="{ required: true, message: '请输入分类名称', trigger: ['input', 'blur'] }"
        >
          <NInput
            v-model:value="modalForm.name"
            placeholder="请输入分类名称"
            clearable
          />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
