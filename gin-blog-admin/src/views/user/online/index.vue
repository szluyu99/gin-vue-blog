<template>
  <CommonPage show-footer title="在线用户">
    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="getOnlineUsers"
    >
      <!-- <template #queryBar>
        <QueryBarItem label="昵称" :label-width="40">
          <n-input
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="请输入用户昵称"
            @keydown.enter="$table?.handleSearch"
          />
        </QueryBarItem>
      </template> -->
    </CrudTable>
  </CommonPage>
</template>

<script setup>
import { NButton, NImage, NPopconfirm } from 'naive-ui'
import { formatDateTime, renderIcon } from '@/utils'
// import { useCRUD } from '@/hooks'

import { useUserApi } from '@/api'
const { getOnlineUsers, forceOfflineUser } = useUserApi()

defineOptions({ name: '在线用户' })

const $table = ref(null)
const queryItems = ref({})

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
        height: 50,
        imgProps: { style: { 'border-radius': '3px' } },
        src: row['avatar'],
        'fallback-src': 'http://dummyimage.com/400x400', // 加载失败
        'show-toolbar-tooltip': true,
      })
    },
  },
  { title: '昵称', key: 'nickname', width: 60, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '登录 IP',
    key: 'ip_address',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row['ip_address'] || '未知')
    },
  },
  {
    title: '登录地址',
    key: 'ip_source',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row['ip_source'] || '未知')
    },
  },
  {
    title: '登录浏览器',
    key: 'browser',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row['browser'] || '未知')
    },
  },
  {
    title: '操作系统',
    key: 'os',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row['os'] || '未知')
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
        {
          size: 'small',
          type: 'text',
          ghost: true,
        },
        {
          default: () => formatDateTime(row['last_login_time'], 'YYYY-MM-DD'),
          icon: renderIcon('mdi:update', { size: 18 }),
        }
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
        {
          onPositiveClick: () => handleForceOffline(row),
          onNegativeClick: () => {},
        },
        {
          trigger: () =>
            h(
              NButton,
              {
                size: 'small',
                type: 'primary',
              },
              {
                default: () => '下线',
                icon: renderIcon('material-symbols:delete-outline', { size: 14 }),
              }
            ),
          default: () => h('div', {}, '确定强制该用户下线吗?'),
        }
      )
    },
  },
]

// 强制用户下线
async function handleForceOffline(row) {
  try {
    await forceOfflineUser(row)
    $message.success('该用户已被强制下线!')
    $table.value?.handleSearch()
  } catch (err) {
    $message.error('强制下线失败!')
  }
}

// useCRUD({
//   name: '用户',
//   doUpdate: updateUser,
//   refresh: () => $table.value?.handleSearch(),
// })
</script>
