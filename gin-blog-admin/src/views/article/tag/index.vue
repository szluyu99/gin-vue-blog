<script setup>
import { NButton, NPopconfirm, NTag } from 'naive-ui'
import { formatDate, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import api from '@/api'

defineOptions({ name: '标签管理' })

const $table = $ref(null)
const queryItems = $ref({})

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
  name: '标签',
  initForm: {},
  doCreate: api.saveOrUpdateTag,
  doDelete: api.deleteTag,
  doUpdate: api.saveOrUpdateTag,
  refresh: () => $table?.handleSearch(),
})

onMounted(() => {
  $table?.handleSearch()
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  {
    title: '标签名',
    key: 'name',
    width: 100,
    align: 'center',
    render(row) {
      return h(NTag, { type: 'info' }, { default: () => row.name })
    },
  },
  {
    title: '文章量',
    key: 'article_count',
    width: 30,
    align: 'center',
  },
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
          icon: renderIcon('mdi:clock-time-three-outline', { size: 17 }),
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
          icon: renderIcon('mdi:update', { size: 18 }),
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
          { default: () => '编辑', icon: renderIcon('material-symbols:edit-outline', { size: 14 }) },
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () => h(
              NButton,
              { size: 'small', type: 'error', style: 'margin-left: 15px;' },
              { default: () => '删除', icon: renderIcon('material-symbols:delete-outline', { size: 14 }) },
            ),
            default: () => h('div', {}, '确定删除该标签吗?'),
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="标签管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <TheIcon icon="material-symbols:add" :size="18" /> 新建标签
      </NButton>
      <NButton
        ml-20
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <TheIcon icon="material-symbols:playlist-remove" :size="18" /> 批量删除
      </NButton>
    </template>

    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getTags"
    >
      <template #queryBar>
        <QueryBarItem label="标签名" :label-width="50">
          <n-input
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入标签名"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryBarItem>
      </template>
    </CrudTable>

    <!-- 新增/编辑/查看 弹窗 -->
    <CrudModal
      v-model:visible="modalVisible"
      :title="modalTitle"
      :loading="modalLoading"
      @on-save="handleSave"
    >
      <!-- 表单 -->
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <n-form-item
          label="文章标签"
          path="name"
          :rule="{ required: true, message: '请输入标签名称', trigger: ['input', 'blur'] }"
        >
          <n-input
            v-model:value="modalForm.name"
            placeholder="请输入标签名称"
            clearable
          />
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>
