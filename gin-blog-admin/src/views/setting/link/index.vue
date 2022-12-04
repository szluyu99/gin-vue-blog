<script setup>
import { NButton, NImage } from 'naive-ui'
import { useClipboard } from '@vueuse/core'
import { formatDateTime, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import api from '@/api'

// vueues - 剪切板操作
const { copy } = useClipboard()

defineOptions({ name: '友链管理' })

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
  name: '友链',
  initForm: {},
  doCreate: api.saveOrUpdateLink,
  doDelete: api.deleteLinks,
  doUpdate: api.saveOrUpdateLink,
  refresh: () => handleSearch(),
})

onMounted(() => {
  handleSearch()
})

const columns = [
  { type: 'selection', width: 20, fixed: 'left' },
  {
    title: '头像',
    key: 'avatar',
    width: 40,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 50,
        'imgProps': { style: { 'border-radius': '3px' } },
        'src': row.avatar,
        'fallback-src': 'http://dummyimage.com/400x400', // 加载失败
        'show-toolbar-tooltip': true,
      })
    },
  },
  { title: '链接名', key: 'name', width: 100, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '链接地址',
    key: 'address',
    width: 120,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h(
        'a',
        {
          // href: row['address'],
          // target: '_blank',
          style: 'cursor: pointer',
          onClick: () => {
            copy(row.address)
            $message.info('复制到剪切板!')
          },
        },
        row.address,
      )
    },
  },
  { title: '链接介绍', key: 'intro', width: 120, align: 'center', ellipsis: { tooltip: true } },
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

function handleSearch() {
  selections.value = []
  $table.value?.handleSearch()
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="友链管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <TheIcon icon="material-symbols:add" :size="18" mr-5 /> 新建友链
      </NButton>
      <NButton
        ml-20
        type="error"
        :disabled="!selections.length"
        @click="handleDelete(JSON.stringify(selections))"
      >
        <TheIcon icon="material-symbols:delete" :size="18" mr-5 /> 批量删除
      </NButton>
    </template>

    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getLinks"
      :selections="selections"
      @on-checked="(rowKeys) => (selections = rowKeys)"
    >
      <template #queryBar>
        <QueryBarItem label="友链名" :label-width="50">
          <n-input
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入友链名"
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
          label="链接名"
          path="name"
          :rule="{ required: true, message: '请输入友链名称', trigger: ['input', 'blur'] }"
        >
          <n-input v-model:value="modalForm.name" placeholder="请输入友链名称" />
        </n-form-item>
        <n-form-item
          label="链接头像"
          path="avatar"
          :rule="{ required: true, message: '请输入友链头像', trigger: ['input', 'blur'] }"
        >
          <n-input v-model:value="modalForm.avatar" placeholder="请输入链接头像" />
        </n-form-item>
        <n-form-item
          label="链接地址"
          path="address"
          :rule="{ required: true, message: '请输入友链地址', trigger: ['input', 'blur'] }"
        >
          <n-input v-model:value="modalForm.address" placeholder="请输入链接地址" />
        </n-form-item>
        <n-form-item
          label="链接介绍"
          path="intro"
          :rule="{ required: true, message: '请输入友链介绍', trigger: ['input', 'blur'] }"
        >
          <n-input v-model:value="modalForm.intro" placeholder="请输入友链地址" />
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>
