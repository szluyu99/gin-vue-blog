<script setup>
import { NButton, NImage, NPopconfirm, NTag } from 'naive-ui'
import { convertImgUrl, formatDate, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import { commentTypeMap, commentTypeOptions } from '@/constant/data'
import api from '@/api'

defineOptions({ name: '评论管理' })

onMounted(() => {
  handleChangeTab('all') // 默认查看全部
})

const $table = $ref(null)
const queryItems = $ref({}) // 条件查询参数
const extraParams = $ref({}) // 额外参数

const { handleDelete } = useCRUD({
  name: '评论',
  doDelete: api.deleteComments,
  refresh: () => $table?.handleSearch(),
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  {
    title: '头像',
    key: 'avatar',
    width: 50,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 50,
        'imgProps': { style: { 'border-radius': '3px' } },
        'src': convertImgUrl(row.avatar),
        'fallback-src': 'http://dummyimage.com/400x400', // 加载失败
        'show-toolbar-tooltip': true,
      })
    },
  },
  {
    title: '评论人',
    key: 'nickname',
    width: 50,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: '回复对象',
    key: 'reply_nick_name',
    width: 50,
    align: 'center',
    render(row) {
      return h('span', row.reply_nickname || '无')
    },
  },
  {
    title: '评论内容',
    key: 'content',
    width: 140,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: '评论时间',
    key: 'created_at',
    align: 'center',
    width: 60,
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
    title: '来源',
    key: 'type',
    width: 50,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: commentTypeMap[row.type].tag },
        { default: () => commentTypeMap[row.type].name },
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
              style: 'margin-left: 15px;',
              onClick: () => handleUpdateReview([row.id], 0),
            },
            {
              default: () => '撤下',
              icon: renderIcon('mi:circle-error', { size: 14 }),
            },
          )
          : h(
            NButton,
            {
              size: 'small',
              type: 'success',
              style: 'margin-left: 15px;',
              onClick: () => handleUpdateReview([row.id], 1),
            },
            {
              default: () => '通过',
              icon: renderIcon('mi:circle-check', { size: 14 }),
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
                { default: () => '删除', icon: renderIcon('material-symbols:delete-outline', { size: 14 }) },
              ),
            default: () => h('div', {}, '确定删除该条评论吗?'),
          },
        ),
      ]
    },
  },
]

// 修改评论审核: is_review 0-撤下审核, 1-通过审核
async function handleUpdateReview(ids, is_review) {
  if (!ids.length) {
    $message.info('请选择要审核的数据')
    return
  }
  await api.updateCommentReview({ ids, is_review })
  $message?.success(is_review ? '审核成功' : '撤下成功')
  $table?.handleSearch()
}

// 切换标签页: [全部, 通过, 审核中]
function handleChangeTab(value) {
  switch (value) {
    case 'all':
      extraParams.is_review = null
      break
    case 'has_review': // 通过
      extraParams.is_review = 1
      break
    case 'not_review': // 审核中
      extraParams.is_review = 0
      break
  }
  $table?.handleSearch()
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="评论管理">
    <!-- 操作栏 -->
    <template #action>
      <NButton
        ml-20
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <TheIcon icon="material-symbols:recycling-rounded" :size="18" mr-5 /> 批量删除
      </NButton>
      <NButton
        ml-20
        type="success"
        :disabled="!$table?.selections.length"
        @click="handleUpdateReview($table.selections, 1)"
      >
        <TheIcon icon="ic:outline-approval" :size="18" mr-5 /> 批量通过
      </NButton>
    </template>
    <!-- 标签栏 -->
    <n-tabs
      type="line"
      animated
      @update:value="handleChangeTab"
    >
      <template #prefix>
        状态
      </template>
      <n-tab-pane name="all" tab="全部" />
      <n-tab-pane name="has_review" tab="通过" />
      <n-tab-pane name="not_review" tab="审核中" />
    </n-tabs>
    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :extra-params="extraParams"
      :columns="columns"
      :get-data="api.getComments"
    >
      <template #queryBar>
        <QueryBarItem label="用户" :label-width="40" :content-width="180">
          <n-input
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="请输入用户昵称"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryBarItem>
        <QueryBarItem label="来源" :label-width="40" :content-width="160">
          <n-select
            v-model:value="queryItems.type"
            clearable
            filterable
            placeholder="请选择评论来源"
            :options="commentTypeOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryBarItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
