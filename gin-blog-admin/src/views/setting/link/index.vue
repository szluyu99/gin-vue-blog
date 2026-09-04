<script setup>
import { useClipboard } from '@vueuse/core'
import { NButton, NForm, NFormItem, NImage, NInput, NPopconfirm } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { convertImgUrl, formatDate } from '@/utils'

defineOptions({ name: '友链管理' })

// 在 setup 里创建一次: 原来写在列的 render 里, 每次点击都新建一个实例
const { copy } = useClipboard()

const $table = ref(null)
const queryItems = ref({
  keyword: '', // 友链名称 | 地址 | 介绍
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
  name: '友链',
  initForm: {},
  doCreate: api.saveOrUpdateLink,
  doDelete: api.deleteLinks,
  doUpdate: api.saveOrUpdateLink,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  {
    title: '头像',
    key: 'avatar',
    width: 40,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 40,
        'imgProps': { style: { 'border-radius': '3px' } },
        // 本地上传的头像存的是相对路径, 不转换直接给 img 会裂图
        'src': convertImgUrl(row.avatar),
        'fallback-src': 'https://dummyimage.com/400x400', // 加载失败
        'show-toolbar-tooltip': true,
      })
    },
  },
  {
    title: '链接名称',
    key: 'name',
    width: 100,
    align: 'center',
    ellipsis: { tooltip: true },
  },
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
          class: 'hover:underline hover:underline-blue-500 hover:underline-2 hover:underline-solid hover:underline-offset-4 cursor-pointer',
          // href: row.address,
          // target: '_blank',
          onClick: () => {
            copy(row.address)
            $message.info('链接已经复制到剪切板!')
          },
        },
        row.address,
      )
    },
  },
  {
    title: '链接介绍',
    key: 'intro',
    width: 120,
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
          default: () => formatDate(row.created_at),
          icon: () => h('i', { class: 'i-mdi:clock-time-three-outline' }),
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
            default: () => h('div', {}, '确定删除该友链吗?'),
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <CommonPage title="友链管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <template #icon>
          <span class="i-material-symbols:add" />
        </template>
        新建友链
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
      :get-data="api.getLinks"
    >
      <template #queryBar>
        <QueryItem label="友链名称 | 地址 | 介绍" :label-width="150">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="搜索关键字"
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
          label="链接名称"
          path="name"
          :rule="{ required: true, message: '请输入友链名称', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.name" placeholder="请输入友链名称" />
        </NFormItem>
        <NFormItem
          label="链接头像"
          path="avatar"
          :rule="{ required: true, message: '请输入友链头像', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.avatar" placeholder="请输入链接头像" />
        </NFormItem>
        <NFormItem
          label="链接地址"
          path="address"
          :rule="{ required: true, message: '请输入友链地址', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.address" placeholder="请输入链接地址" />
        </NFormItem>
        <NFormItem
          label="链接介绍"
          path="intro"
          :rule="{ required: true, message: '请输入友链介绍', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.intro" placeholder="请输入友链地址" />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
