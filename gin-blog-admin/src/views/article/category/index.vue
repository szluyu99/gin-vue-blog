<script setup>
import { NButton } from 'naive-ui'
import TheIcon from '@/components/icon/TheIcon.vue'
import { formatDateTime, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import api from '@/api'

defineOptions({ name: '分类管理' })

const $table = ref(null)
const queryItems = ref({})
const selections = ref([])

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
  refresh: () => handleSearch(),
})

onMounted(() => {
  handleSearch()
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  { title: '分类名', key: 'name', width: 100, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '文章量',
    key: 'article_count',
    width: 30,
    align: 'center',
    ellipsis: { tooltip: true },
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
          default: () => formatDateTime(row.created_at, 'YYYY-MM-DD'),
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
          default: () => formatDateTime(row.updated_at, 'YYYY-MM-DD'),
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
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 15px;',
            onClick: () => handleEdit(row),
          },
          { default: () => '编辑', icon: renderIcon('material-symbols:edit-outline', { size: 14 }) },
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 15px;',
            onClick: () => handleDelete(JSON.stringify([row.id])),
          },
          {
            default: () => '删除',
            icon: renderIcon('material-symbols:delete-outline', { size: 14 }),
          },
        ),
      ]
    },
  },
]

// 刷新时添加额外逻辑: 清空选中列表
function handleSearch() {
  selections.value = []
  $table.value?.handleSearch()
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="分类管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <TheIcon icon="material-symbols:add" :size="18" mr-5 /> 新建分类
      </NButton>
      <NButton
        ml-20
        type="error"
        :disabled="!selections.length"
        @click="handleDelete(JSON.stringify(selections))"
      >
        <TheIcon icon="material-symbols:playlist-remove" :size="18" mr-5 /> 批量删除
      </NButton>
    </template>

    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getCategorys"
      :selections="selections"
      @on-checked="(rowKeys) => (selections = rowKeys)"
    >
      <template #queryBar>
        <QueryBarItem label="分类名" :label-width="50">
          <n-input
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入分类名"
            @keydown.enter="handleSearch"
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
          label="文章分类"
          path="name"
          :rule="{ required: true, message: '请输入分类名称', trigger: ['input', 'blur'] }"
        >
          <n-input v-model:value="modalForm.name" placeholder="请输入分类名称" />
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>
