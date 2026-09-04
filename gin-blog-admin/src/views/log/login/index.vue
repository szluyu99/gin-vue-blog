<script setup>
import { NButton, NInput, NPopconfirm, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudTable from '@/components/crud/CrudTable.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { formatDate } from '@/utils'

defineOptions({ name: '登录日志' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

const { handleDelete } = useCRUD({
  name: '登录日志',
  doDelete: api.deleteLoginLogs,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

// 与后端 model.LOGIN_SUCCESS / LOGIN_FAIL 对应
const STATUS_SUCCESS = 1

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  { title: '用户名', key: 'username', width: 100, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '昵称',
    key: 'nickname',
    width: 80,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      // 登录失败时拿不到昵称(用户可能都不存在)
      return h('span', row.nickname || '-')
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 50,
    align: 'center',
    render(row) {
      const success = row.status === STATUS_SUCCESS
      return h(
        NTag,
        { type: success ? 'success' : 'error' },
        { default: () => (success ? '成功' : '失败') },
      )
    },
  },
  {
    title: '失败原因',
    key: 'message',
    width: 100,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.message || '-')
    },
  },
  { title: '登录IP', key: 'ip_address', width: 90, align: 'center', ellipsis: { tooltip: true } },
  { title: '登录地址', key: 'ip_source', width: 90, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '登录时间',
    key: 'created_at',
    width: 80,
    align: 'center',
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.created_at, 'YYYY-MM-DD HH:mm:ss'),
          icon: () => h('i', { class: 'i-mdi:clock-time-three-outline' }),
        },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 60,
    align: 'center',
    fixed: 'right',
    render(row) {
      return h(
        NPopconfirm,
        { onPositiveClick: () => handleDelete([row.id], false) },
        {
          trigger: () => h(
            NButton,
            { size: 'small', quaternary: true, type: 'error' },
            {
              default: () => '删除',
              icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
            },
          ),
          default: () => h('div', {}, '确定删除该条登录日志吗?'),
        },
      )
    },
  },
]
</script>

<template>
  <CommonPage title="登录日志">
    <template #action>
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
      :get-data="api.getLoginLogs"
    >
      <template #queryBar>
        <QueryItem label="关键字" :label-width="60" :content-width="200">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="用户名 / 昵称 / IP"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
