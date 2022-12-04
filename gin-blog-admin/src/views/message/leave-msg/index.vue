<script setup>
import { NButton, NImage, NPopconfirm, NTag } from 'naive-ui'
import TheIcon from '@/components/icon/TheIcon.vue'
import { formatDateTime, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import api from '@/api'

defineOptions({ name: '留言管理' })

const { handleDelete } = useCRUD({
  name: '留言',
  doDelete: api.deleteMessages,
  refresh: handleSearch,
})

onMounted(async () => {
  handleChangeTab('all') // 默认查看全部
})

const $table = ref(null)
const queryItems = ref({})
const extraParams = ref({})
const selections = ref([])

const columns = [
  { type: 'selection', width: 20, fixed: 'left' },
  {
    title: '头像',
    key: 'avatar',
    width: 50,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 70,
        'imgProps': { style: { 'border-radius': '3px' } },
        'src': row.avatar,
        'fallback-src': 'http://dummyimage.com/400x400', // 加载失败
        'show-toolbar-tooltip': true,
      })
    },
  },
  { title: '留言人', key: 'nickname', width: 60, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '留言内容',
    key: 'content',
    width: 120,
    align: 'center',
  },
  { title: 'IP 地址', key: 'ip_address', width: 70, align: 'center', ellipsis: { tooltip: true } },
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
      // return h('p', [
      //   h(TheIcon, { icon: 'mdi:update', size: 18 }),
      //   h('span', formatDateTime(row['created_at'], 'YYYY-MM-DD')),
      // ])
      return h(
        NButton,
        {
          size: 'small',
          type: 'text',
          ghost: true,
          // onClick: () => {
          //   $message.info('成功复制到剪切板!')
          // },
          // style: 'cursor: default;',
        },
        {
          default: () => formatDateTime(row.created_at, 'YYYY-MM-DD'),
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
        { type: row.is_review === 1 ? 'success' : 'error' },
        { default: () => (row.is_review === 1 ? '通过' : '审核中') },
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
        row.is_review === 0
          ? h(
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
          )
          : h(
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
          ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(JSON.stringify([row.id]), false),
            onNegativeClick: () => {},
          },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  style: 'margin-left: 15px;',
                },
                {
                  default: () => '删除',
                  icon: renderIcon('material-symbols:delete-outline', { size: 14 }),
                },
              ),
            default: () => h('div', {}, '确定删除该条留言吗?'),
          },
        ),
      ]
    },
  },
]

// 修改留言审核: is_review 0-撤下审核, 1-通过审核
async function handleUpdateReview(ids, is_review) {
  await api.updateMessageReview({ ids, is_review })
  $message?.success(is_review ? '审核成功' : '撤下成功')
  handleSearch()
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
  handleSearch()
}

// 刷新时添加额外逻辑: 清空选中列表
function handleSearch() {
  selections.value = []
  $table.value?.handleSearch()
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="留言管理">
    <!-- 操作栏 -->
    <template #action>
      <NButton
        ml-20
        type="error"
        :disabled="!selections.length"
        @click="handleDelete(JSON.stringify(selections))"
      >
        <TheIcon icon="material-symbols:recycling-rounded" :size="18" mr-5 /> 批量删除
      </NButton>
      <NButton
        ml-20
        type="success"
        :disabled="!selections.length"
        @click="handleUpdateReview(selections, 1)"
      >
        <TheIcon icon="ic:outline-approval" :size="18" mr-5 /> 批量通过
      </NButton>
    </template>
    <!-- 标签栏 -->
    <n-tabs type="line" animated @update:value="handleChangeTab">
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
      :get-data="api.getMessages"
      :selections="selections"
      @on-checked="(rowKeys) => (selections = rowKeys)"
    >
      <template #queryBar>
        <QueryBarItem label="用户" :label-width="40" :content-width="180">
          <n-input
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="请输入用户昵称"
            @keydown.enter="handleSearch"
          />
        </QueryBarItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
