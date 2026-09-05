<script setup>
import { NButton, NImage, NInput, NPopconfirm } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import QueryItem from '@/components/crud/QueryItem.vue'
import { convertImgUrl, formatDate, IMG_PLACEHOLDER } from '@/utils'

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
        'src': convertImgUrl(row.info?.avatar),
        'fallback-src': IMG_PLACEHOLDER, // 加载失败时用内联占位图, 不再请求外网
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
      // info 来自 Redis 里反序列化的 UserAuth, 没有 Preload("UserInfo"), 可能为空;
      // 少一个可选链就会让整张表 render 抛错(user/list 那边写的是 row.info?.nickname)
      return h('span', row.info?.nickname || '未知')
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
