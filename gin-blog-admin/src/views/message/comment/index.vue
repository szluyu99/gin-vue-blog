<script setup>
import { NButton, NImage, NInput, NPopconfirm, NSelect, NTabPane, NTabs, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import { commentTypeMap, commentTypeOptions } from '@/assets/config'
import CommonPage from '@/components/common/CommonPage.vue'

import CrudTable from '@/components/crud/CrudTable.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { convertImgUrl, formatDate } from '@/utils'

defineOptions({ name: '评论管理' })

onMounted(() => {
  handleChangeTab('all') // 默认查看全部
})

const $table = ref(null)
const queryItems = ref({
  nickname: '',
  type: '',
})
const extraParams = ref({
  is_review: null, // 评论状态: 审核中 | 通过
})

const { handleDelete } = useCRUD({
  name: '评论',
  doDelete: api.deleteComments,
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
        'src': convertImgUrl(row.user?.info?.avatar),
        'fallback-src': 'https://dummyimage.com/400x400', // 加载失败
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
    render(row) {
      return h('span', row.user?.info?.nickname || '无')
    },
  },
  // TODO: 合理的显示评论的文章信息
  {
    title: '回复对象',
    // 这里的 key 只用于导出(CrudTable 按 item[key] 取值), 展示走 render
    key: 'reply_user',
    width: 50,
    align: 'center',
    render(row) {
      return h('span', row.reply_user?.info?.nickname || '-')
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
    title: '来源',
    key: 'type',
    width: 50,
    align: 'center',
    render(row) {
      // 后端新增评论类型时前端的映射表可能还没跟上, 不能直接取 .tag
      const type = commentTypeMap[row.type]
      return h(
        NTag,
        { type: type?.tag ?? 'default' },
        { default: () => type?.name ?? '未知' },
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
            default: () => h('div', {}, '确定删除该条评论吗?'),
          },
        ),
      ]
    },
  },
]

// 修改评论审核
async function handleUpdateReview(ids, is_review) {
  if (!ids.length) {
    window.$message.info('请选择要审核的数据')
    return
  }
  // 失败时拦截器已经弹过提示, 这里只要别让列表刷新和成功提示误报
  try {
    await api.updateCommentReview(ids, is_review)
  }
  catch (err) {
    console.error(err)
    return
  }
  window.$message?.success(is_review ? '审核成功' : '撤下成功')
  $table.value?.handleSearch()
}

// 切换标签页: [全部, 通过, 审核中]
function handleChangeTab(value) {
  switch (value) {
    case 'all':
      extraParams.value.is_review = null
      break
    case 'has_review': // 通过
      extraParams.value.is_review = true
      break
    case 'not_review': // 审核中
      extraParams.value.is_review = false
      break
  }
  $table.value?.handleSearch()
}
</script>

<template>
  <CommonPage title="评论管理">
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
      :get-data="api.getComments"
    >
      <template #queryBar>
        <QueryItem label="用户" :label-width="40" :content-width="180">
          <NInput
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="请输入用户昵称"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="来源" :label-width="40" :content-width="160">
          <NSelect
            v-model:value="queryItems.type"
            clearable
            filterable
            placeholder="请选择评论来源"
            :options="commentTypeOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
