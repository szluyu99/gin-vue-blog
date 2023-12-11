<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NImage, NPopconfirm } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { convertImgUrl, formatDate } from '@/utils'
import api from '@/api'

defineOptions({ name: '在线用户' })

const $table = ref(null)

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  {
    title: '头像',
    key: 'avatar',
    width: 40,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 40,
        'imgProps': { style: { 'border-radius': '1px' } },
        'src': convertImgUrl(row.avatar),
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
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.last_login_time),
          icon: () => h('i', { class: 'i-mdi:update' }),
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

// FIXME: 无法强制下线
// 强制用户下线
async function handleForceOffline(row) {
  try {
    await api.forceOfflineUser(row)
    window.$message.success('该用户已被强制下线!')
    $table.value?.handleSearch()
  }
  catch (err) {
    window.$message.error('强制下线失败!')
  }
}
</script>

<template>
  <CommonPage title="在线用户">
    <CrudTable
      ref="$table"
      :columns="columns"
      :get-data="api.getOnlineUsers"
    />
  </CommonPage>
</template>
