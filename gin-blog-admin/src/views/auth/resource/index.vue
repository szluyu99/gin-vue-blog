<script setup>
import { NButton, NPopconfirm, NSwitch, NTag } from 'naive-ui'
import { formatDate, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import api from '@/api'

defineOptions({ name: '接口管理' })

const $table = $ref(null)
const queryItems = $ref({})

const {
  modalVisible,
  modalAction,
  modalTitle,
  modalLoading,
  handleAdd,
  handleDelete,
  handleEdit,
  handleSave,
  modalForm,
  modalFormRef,
} = useCRUD({
  name: '接口',
  doCreate: api.saveOrUpdateResource,
  doDelete: api.deleteResource,
  doUpdate: api.saveOrUpdateResource,
  refresh: () => $table?.handleSearch(),
})

onMounted(() => {
  $table?.handleSearch()
})

// 请求方法
const requestMethods = ['GET', 'POST', 'DELETE', 'PUT']

// 请求方法对应不同类型的标签 (计算属性传参)
const tagType = $computed(() => (type) => {
  switch (type) {
    case 'GET':
      return 'info'
    case 'POST':
      return 'success'
    case 'PUT':
      return 'warning'
    case 'DELETE':
      return 'error'
    default:
      return 'info'
  }
})

const columns = [
  { title: '资源名称', key: 'name', width: 80, ellipsis: { tooltip: true } },
  { title: '资源路径', key: 'url', width: 100, ellipsis: { tooltip: true } },
  {
    title: '请求方式',
    key: 'request_method',
    width: 50,
    align: 'center',
    render(row) {
      return row.children
        ? ''
        : h(
          NTag,
          { type: tagType(row.request_method) }, // 注意这里使用计算属性
          { default: () => row.request_method },
        )
    },
  },
  {
    title: '匿名访问',
    key: 'is_hidden',
    width: 50,
    align: 'center',
    fixed: 'left',
    render(row) {
      return row.children
        ? ''
        : h(NSwitch, {
          size: 'small',
          rubberBand: false,
          value: row.is_anonymous,
          loading: !!row.publishing, // 修改 ing 动画
          checkedValue: 1,
          uncheckedValue: 0,
          onUpdateValue: () => handleUpdateAnonymous(row),
        })
    },
  },
  {
    title: '创建日期',
    key: 'created_at',
    width: 60,
    render(row) {
      return h('span', formatDate(row.created_at))
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 115,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'primary',
            style: `display: ${row.children ? '' : 'none'};`,
            onClick: () => {
              handleAdd() // 新增弹窗
              modalForm.value.parent_id = row.id // 父资源id
            },
          },
          { default: () => '新增', icon: renderIcon('material-symbols:add', { size: 14 }) },
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: () => (row.children ? handleEditModule(row) : handleEdit(row)),
          },
          { default: () => '编辑', icon: renderIcon('material-symbols:edit-outline', { size: 14 }) },
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => {
              handleDelete(row.id, false)
            },
          },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'tiny', quaternary: true, type: 'error' },
                { default: () => '删除', icon: renderIcon('material-symbols:delete-outline', { size: 14 }) },
              ),
            default: () => h('div', {}, '确定删除该接口吗?'),
          },
        ),
      ]
    },
  },
]

// 修改是否允许匿名访问
async function handleUpdateAnonymous(row) {
  if (!row.id)
    return
  row.publishing = true
  row.is_anonymous = row.is_anonymous === 0 ? 1 : 0
  await api.updateResourceAnonymous(row)
  row.publishing = false
  $message?.success(row.is_anonymous ? '已允许匿名访问' : '已禁止匿名访问')
}

// 模块相关
let moduleModalVisible = $ref(false)
function handleAddModule() {
  modalAction.value = 'add'
  modalForm.value = {}
  moduleModalVisible = true
}
function handleEditModule(row) {
  modalAction.value = 'edit'
  modalForm.value = { ...row }
  moduleModalVisible = true
}
async function handleModuleSave() {
  handleSave()
  moduleModalVisible = false
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="接口管理">
    <template #action>
      <NButton type="primary" @click="handleAddModule">
        <TheIcon icon="material-symbols:add" :size="18" /> 新增模块
      </NButton>
    </template>

    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :is-pagination="false"
      :columns="columns"
      :get-data="api.getResources"
      :single-line="true"
    >
      <template #queryBar>
        <QueryBarItem label="资源名" :label-width="50">
          <n-input
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入资源名"
            @keydown.enter="$table?.handleSearch"
          />
        </QueryBarItem>
      </template>
    </CrudTable>

    <!-- 新增/编辑 资源 -->
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
        <n-form-item label="资源名" path="name">
          <n-input v-model:value="modalForm.name" placeholder="请输入资源名" />
        </n-form-item>
        <n-form-item label="资源路径" path="url">
          <n-input v-model:value="modalForm.url" placeholder="请输入资源路径" />
        </n-form-item>
        <n-form-item label="请求方式" path="request_method">
          <n-radio-group v-model:value="modalForm.request_method" name="radiogroup">
            <n-space>
              <n-radio v-for="method of requestMethods" :key="method" :value="method">
                <n-gradient-text :type="tagType(method)">
                  {{ method }}
                </n-gradient-text>
              </n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
      </n-form>
    </CrudModal>

    <!-- 新增/编辑 模块 -->
    <CrudModal
      v-model:visible="moduleModalVisible"
      :title="`${modalAction === 'add' ? '新增' : '编辑'}模块`"
      :loading="modalVisible"
      @on-save="handleModuleSave"
    >
      <!-- 表单 -->
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <n-form-item label="模块名" path="name">
          <n-input v-model:value="modalForm.name" placeholder="请输入模块名" />
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>
