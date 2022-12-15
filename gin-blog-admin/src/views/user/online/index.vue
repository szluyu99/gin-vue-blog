<script setup>
import { NButton, NImage, NPopconfirm } from 'naive-ui'
import { convertImgUrl, formatDate, renderIcon } from '@/utils'
import api from '@/api'

defineOptions({ name: '在线用户' })

const $table = $ref(null)

onMounted(() => {
  $table?.handleSearch()
})

const columns = [
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
          icon: renderIcon('mdi:update', { size: 18 }),
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
              { size: 'small', type: 'primary' },
              {
                default: () => '下线',
                icon: renderIcon('material-symbols:delete-outline', { size: 14 }),
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
    await api.forceOfflineUser(row)
    $message.success('该用户已被强制下线!')
    $table?.handleSearch()
  }
  catch (err) {
    $message.error('强制下线失败!')
  }
}
</script>

<template>
  <CommonPage show-footer title="在线用户">
    <CrudTable
      ref="$table"
      :columns="columns"
      :get-data="api.getOnlineUsers"
    />
  </CommonPage>
</template>
