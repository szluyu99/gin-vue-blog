<script setup>
import { NButton, NPopconfirm, NTag } from 'naive-ui'
import { formatDate, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import api from '@/api'

defineOptions({ name: '操作日志' })

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

const $table = $ref(null)
const queryItems = $ref({})

const { modalVisible, modalLoading, handleDelete, modalForm, modalFormRef, handleView } = useCRUD({
  name: '日志',
  doDelete: api.deleteOperationLogs,
  refresh: () => $table?.handleSearch(),
})

onMounted(() => {
  $table?.handleSearch()
})

const columns = [
  { type: 'selection', width: 20, fixed: 'left' },
  { title: '系统模块', key: 'opt_module', width: 70, align: 'center', ellipsis: { tooltip: true } },
  { title: '操作类型', key: 'opt_type', width: 70, align: 'center', ellipsis: { tooltip: true } },
  // { title: '操作描述', key: 'opt_desc', width: 80, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '请求方法',
    key: 'request_method',
    width: 80,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h(
        NTag,
        { type: tagType(row.request_method) }, // 注意这里使用计算属性
        { default: () => row.request_method },
      )
    },
  },
  { title: '操作人员', key: 'nickname', width: 80, align: 'center', ellipsis: { tooltip: true } },
  { title: '登录IP', key: 'ip_address', width: 80, align: 'center', ellipsis: { tooltip: true } },
  { title: '登录地址', key: 'ip_source', width: 80, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '发布时间',
    key: 'created_at',
    align: 'center',
    width: 80,
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.created_at),
          icon: renderIcon('mdi:update', { size: 18 }),
        },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            quaternary: true,
            type: 'info',
            onClick: () => handleView(row),
          },
          {
            default: () => '查看',
            icon: renderIcon('ic:outline-remove-red-eye', { size: 16 }),
          },
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'small',
                  quaternary: true,
                  type: 'error',
                  style: 'margin-left: 15px;',
                },
                {
                  default: () => '删除',
                  icon: renderIcon('material-symbols:delete-outline', { size: 16 }),
                },
              ),
            default: () => h('div', {}, '确定删除该日志吗?'),
          },
        ),
      ]
    },
  },
]

function copyFormatCode(code) {
  const { copy } = useClipboard()
  copy(JSON.stringify(JSON.parse(code), null, 2))
  $message.success('内容已复制到剪切板!')
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="操作日志">
    <template #action>
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
      :get-data="api.getOperationLogs"
    >
      <template #queryBar>
        <QueryBarItem label="模块名" :label-width="50">
          <n-input
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入模块名或描述"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryBarItem>
      </template>
    </CrudTable>

    <!-- 新增/编辑 弹窗 -->
    <CrudModal
      v-model:visible="modalVisible"
      title="日志详情"
      :show-footer="false"
      :loading="modalLoading"
    >
      <!-- 表单 -->
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="90"
        :model="modalForm"
      >
        <n-form-item label="操作模块: " path="opt_module">
          {{ modalForm.opt_module }}
        </n-form-item>
        <n-form-item label="请求地址: " path="opt_url">
          {{ modalForm.opt_url }}
        </n-form-item>
        <n-form-item label="请求方法: " path="request_method">
          <NTag :type="tagType(modalForm.request_method)">
            {{ modalForm.request_method }}
          </NTag>
        </n-form-item>
        <n-form-item label="操作类型: " path="opt_type">
          {{ modalForm.opt_type }}
        </n-form-item>
        <n-form-item label="操作方法: " path="opt_method">
          <n-code
            :code="modalForm.opt_method"
            code-wrap
            language="json"
          />
        </n-form-item>
        <n-form-item label="操作人员: " path="nickname">
          {{ modalForm.nickname }}
        </n-form-item>
        <n-form-item label="请求参数: " path="request_param">
          <n-code
            p-7 cursor-pointer
            word-wrap
            :code="modalForm.request_param"
            language="json"
            @click="copyFormatCode(modalForm.request_param)"
          />
        </n-form-item>
        <n-form-item label="返回数据: " path="response_data">
          <n-code
            p-7 cursor-pointer
            :code="JSON.stringify(JSON.parse(modalForm.response_data), null, 2)"
            language="json"
            @click="copyFormatCode(modalForm.response_data)"
          />
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>
