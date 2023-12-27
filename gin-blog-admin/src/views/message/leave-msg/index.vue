<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NImage, NInput, NPopconfirm, NTabPane, NTabs, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { convertImgUrl, formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: '留言管理' })

onMounted(() => {
  handleChangeTab('all') // 默认查看全部
})

const $table = ref(null)
const queryItems = ref({
  nickname: '',
})
const extraParams = ref({
  is_review: null, // 评论状态: 审核中 | 通过
})

const { handleDelete } = useCRUD({
  name: '留言',
  doDelete: api.deleteMessages,
  refresh: () => $table.value?.handleSearch(),
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
        'src': convertImgUrl(row.avatar),
        'fallback-src': 'http://dummyimage.com/400x400', // 加载失败
        'show-toolbar-tooltip': true,
      })
    },
  },
  {
    title: '留言人',
    key: 'nickname',
    width: 60,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: '留言内容',
    key: 'content',
    width: 120,
    align: 'center',
  },
  {
    title: 'IP 地址',
    key: 'ip_address',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'IP 来源',
    key: 'ip_source',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.ip_source || '未知')
    },
  },
  {
    title: '留言时间',
    key: 'created_at',
    align: 'center',
    width: 80,
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.created_at),
          icon: () => h('i', { class: 'i-mdi:update' }),
        },
      )
    },
  },
  {
    title: '状态',
    key: 'is_review',
    width: 50,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: row.is_review ? 'success' : 'error' },
        { default: () => (row.is_review ? '通过' : '审核中') },
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
        row.is_review
          ? h(
            NButton,
            {
              size: 'small',
              type: 'warning',
              onClick: () => handleUpdateReview([row.id], false),
            },
            {
              default: () => '撤下',
              icon: () => h('i', { class: 'i-mi:circle-error' }),
            },
          )
          : h(
            NButton,
            {
              size: 'small',
              type: 'success',
              style: 'margin-left: 15px;',
              onClick: () => handleUpdateReview([row.id], true),
            },
            {
              default: () => '通过',
              icon: () => h('i', { class: 'i-mi:circle-check' }),
            },
          ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'small', type: 'error', style: 'margin-left: 15px;' },
                { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
              ),
            default: () => h('div', {}, '确定删除该条留言吗?'),
          },
        ),
      ]
    },
  },
]

// 修改留言审核
async function handleUpdateReview(ids, is_review) {
  if (!ids.length) {
    $message.info('请选择要审核的数据')
    return
  }

  await api.updateMessageReview(ids, is_review)
  $message?.success(is_review ? '审核成功' : '撤下成功')
  $table.value?.handleSearch()
}

// 切换标签页: [全部, 通过, 审核中]
function handleChangeTab(value) {
  switch (value) {
    case 'all':
      extraParams.value.is_review = null
      break
    case 'has_review': // 通过
      extraParams.value.is_review = 1
      break
    case 'not_review': // 审核中
      extraParams.value.is_review = 0
      break
  }
  $table.value?.handleSearch()
}
</script>

<template>
  <CommonPage title="留言管理">
    <template #action>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <span class="i-material-symbols:recycling-rounded" />
        </template>
        批量删除
      </NButton>
      <NButton
        type="success"
        :disabled="!$table?.selections.length"
        @click="handleUpdateReview($table.selections, true)"
      >
        <template #icon>
          <span class="i-ic:outline-approval" />
        </template>
        批量通过
      </NButton>
    </template>
    <NTabs
      type="line"
      animated
      @update:value="handleChangeTab"
    >
      <template #prefix>
        状态
      </template>
      <NTabPane name="all" tab="全部" />
      <NTabPane name="has_review" tab="通过" />
      <NTabPane name="not_review" tab="审核中" />
    </NTabs>
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :extra-params="extraParams"
      :columns="columns"
      :get-data="api.getMessages"
    >
      <template #queryBar>
        <QueryItem label="用户" :label-width="40" :content-width="180">
          <NInput
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="请输入用户昵称"
            @keydown.enter=" $table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
