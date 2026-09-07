<script setup>
import { NButton, NInput, NPopconfirm, NSelect, NTag, NTooltip } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudTable from '@/components/crud/CrudTable.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { formatDate } from '@/utils'

defineOptions({ name: '前端错误' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
  source: '',
})

const { handleDelete } = useCRUD({
  name: '错误日志',
  doDelete: api.deleteErrorLogs,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

// 与后端 model.ERR_SOURCE_FRONT / ERR_SOURCE_ADMIN 对应
const sourceOptions = [
  { label: '博客前台', value: 'front' },
  { label: '博客后台', value: 'admin' },
]

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  {
    title: '来源',
    key: 'source',
    width: 60,
    align: 'center',
    render(row) {
      const front = row.source === 'front'
      return h(
        NTag,
        { type: front ? 'info' : 'warning', size: 'small' },
        { default: () => (front ? '前台' : '后台') },
      )
    },
  },
  {
    title: '错误信息',
    key: 'message',
    minWidth: 200,
    ellipsis: { tooltip: true },
  },
  {
    title: '调用栈',
    key: 'stack',
    width: 90,
    align: 'center',
    // 栈很长, 列表里只给一个悬浮查看, 不占版面
    render(row) {
      if (!row.stack) {
        return h('span', '-')
      }
      return h(NTooltip, { style: { maxWidth: '600px' } }, {
        trigger: () => h(NButton, { size: 'small', text: true, type: 'primary' }, { default: () => '查看' }),
        default: () => h('pre', { style: { margin: 0, whiteSpace: 'pre-wrap', fontSize: '12px' } }, row.stack),
      })
    },
  },
  {
    title: '出错页面',
    key: 'url',
    width: 160,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.url || '-')
    },
  },
  {
    // 聚合之后这一列才是判断严重程度的主要依据
    title: '次数',
    key: 'count',
    width: 50,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: row.count > 10 ? 'error' : 'default', size: 'small' },
        { default: () => String(row.count ?? 1) },
      )
    },
  },
  {
    title: '用户',
    key: 'user_id',
    width: 50,
    align: 'center',
    render(row) {
      // 匿名访客上报时为 0
      return h('span', row.user_id ? String(row.user_id) : '游客')
    },
  },
  { title: 'IP', key: 'ip_address', width: 90, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '首次出现',
    key: 'created_at',
    width: 100,
    align: 'center',
    render(row) {
      return h('span', formatDate(row.created_at, 'MM-DD HH:mm:ss'))
    },
  },
  {
    // 排序按的是这一列: 关心的是"最近还在发生"
    title: '最近出现',
    key: 'updated_at',
    width: 100,
    align: 'center',
    render(row) {
      return h('span', formatDate(row.updated_at, 'MM-DD HH:mm:ss'))
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
          default: () => h('div', {}, '确定删除该条错误日志吗?'),
        },
      )
    },
  },
]
</script>

<template>
  <CommonPage title="前端错误">
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
      :get-data="api.getErrorLogs"
    >
      <template #queryBar>
        <QueryItem label="关键字" :label-width="60" :content-width="200">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="错误信息 / 页面地址"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="来源" :label-width="40" :content-width="140">
          <NSelect
            v-model:value="queryItems.source"
            clearable
            placeholder="请选择来源"
            :options="sourceOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
