<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NGradientText, NInput, NPopconfirm, NRadio, NRadioGroup, NSpace, NSwitch, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: '接口管理' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

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
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

// 请求方法
const requestMethods = ['GET', 'POST', 'DELETE', 'PUT']

// 请求方法对应不同类型的标签 (计算属性传参)
function tagType(type) {
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
}

const columns = [
  {
    title: '资源名称',
    key: 'name',
    width: 80,
    ellipsis: { tooltip: true },
  },
  {
    title: '资源路径',
    key: 'url',
    width: 80,
    ellipsis: { tooltip: true },
    render(row) {
      return row.children ? '-' : h('span', { class: 'color-[#1890ff]' }, row.url)
    },
  },
  {
    title: '请求方式',
    key: 'request_method',
    width: 50,
    align: 'center',
    render(row) {
      return row.children
        ? '-'
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
        ? '-'
        : h(NSwitch, {
          size: 'small',
          rubberBand: false,
          value: row.is_anonymous,
          loading: !!row.publishing, // 修改 ing 动画
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
          { default: () => '新增', icon: () => h('i', { class: 'i-material-symbols:add' }) },
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: () => (row.children ? handleEditModule(row) : handleEdit(row)),
          },
          { default: () => '编辑', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
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
                { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
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
  if (!row.id) {
    return
  }
  row.publishing = true
  row.is_anonymous = !row.is_anonymous
  try {
    await api.updateResourceAnonymous(row)
    $message?.success(row.is_anonymous ? '已允许匿名访问' : '已禁止匿名访问')
  }
  catch (err) {
    row.is_anonymous = !row.is_anonymous
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

// 模块相关
const moduleModalVisible = ref(false)
function handleAddModule() {
  modalAction.value = 'add'
  modalForm.value = {}
  moduleModalVisible.value = true
}
function handleEditModule(row) {
  modalAction.value = 'edit'
  modalForm.value = { ...row }
  moduleModalVisible.value = true
}
async function handleModuleSave() {
  handleSave()
  moduleModalVisible.value = false
}
</script>

<template>
  <CommonPage title="接口管理">
    <template #action>
      <NButton type="primary" @click="handleAddModule">
        <template #icon>
          <span class="i-material-symbols:add" />
        </template>
        新增模块
      </NButton>
    </template>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :is-pagination="false"
      :columns="columns"
      :get-data="api.getResources"
      :single-line="true"
    >
      <template #queryBar>
        <QueryItem label="资源名" :label-width="50">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入资源名"
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
        <NFormItem label="资源名" path="name">
          <NInput v-model:value="modalForm.name" placeholder="请输入资源名" />
        </NFormItem>
        <NFormItem label="资源路径" path="url">
          <NInput v-model:value="modalForm.url" placeholder="请输入资源路径" />
        </NFormItem>
        <NFormItem label="请求方式" path="request_method">
          <NRadioGroup v-model:value="modalForm.request_method" name="radiogroup">
            <NSpace>
              <NRadio v-for="method of requestMethods" :key="method" :value="method">
                <NGradientText :type="tagType(method)">
                  {{ method }}
                </NGradientText>
              </NRadio>
            </NSpace>
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </CrudModal>

    <CrudModal
      v-model:visible="moduleModalVisible"
      :title="`${modalAction === 'add' ? '新增' : '编辑'}模块`"
      :loading="modalVisible"
      @save="handleModuleSave"
    >
      <NForm
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <NFormItem label="模块名" path="name">
          <NInput v-model:value="modalForm.name" placeholder="请输入模块名" />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
