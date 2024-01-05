<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NImage, NInput, NPopconfirm } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { convertImgUrl, formatDate } from '@/utils'
import api from '@/api'

defineOptions({ name: '在线用户' })

const $table = ref(null)
const queryItems = ref({
  keyword: '', // 用户名 | 昵称
})

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  {
    title: '头像',
    key: 'avatar',
    width: 30,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 30,
        'src': convertImgUrl(row.info.avatar),
        'fallback-src': 'http://dummyimage.com/400x400', // 加载失败
        'show-toolbar-tooltip': true,
      })
    },
  },
  {
    title: '昵称',
    key: 'nickname',
    width: 60,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.info.nickname || '未知')
    },
  },
  {
    title: '登录 IP',
    key: 'ip_address',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.ip_address || '未知')
    },
  },
  {
    title: '登录地址',
    key: 'ip_source',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.ip_source || '未知')
    },
  },
  {
    title: '登录浏览器',
    key: 'browser',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.browser || '未知')
    },
  },
  {
    title: '操作系统',
    key: 'os',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.os || '未知')
    },
  },
  {
    title: '登录时间',
    key: 'last_login_time',
    align: 'center',
    width: 70,
    render(row) {
      return h('span', formatDate(row.last_login_time, 'YYYY-MM-DD HH:mm:ss'))
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
        { onPositiveClick: () => handleForceOffline(row) },
        {
          trigger: () =>
            h(
              NButton,
              { size: 'small', type: 'warning' },
              {
                default: () => '下线',
                icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
              },
            ),
          default: () => h('div', {}, '确定强制该用户下线吗?'),
        },
      )
    },
  },
]

// 强制用户下线
async function handleForceOffline(row) {
  try {
    await api.forceOfflineUser(row.id)
    window.$message.success('该用户已被强制下线!')
    $table.value?.handleSearch()
  }
  catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <CommonPage title="在线用户">
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getOnlineUsers"
      :is-pagination="false"
    >
      <template #queryBar>
        <QueryItem label="用户名 | 昵称" :label-width="100" :content-width="200">
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
  </CommonPage>
</template>
