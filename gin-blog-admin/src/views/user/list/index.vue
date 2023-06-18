<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NCheckbox, NCheckboxGroup, NForm, NFormItem, NImage, NInput, NSelect, NSpace, NSwitch, NTag } from 'naive-ui'

import CommonPage from '@/components/page/CommonPage.vue'
import QueryBarItem from '@/components/query-bar/QueryBarItem.vue'
import CrudModal from '@/components/table/CrudModal.vue'
import CrudTable from '@/components/table/CrudTable.vue'

import { convertImgUrl, formatDate, renderIcon } from '@/utils'
import { useCRUD } from '@/composables'
import { loginTypeMap, loginTypeOptions } from '@/constant/data'
import api from '@/api'

defineOptions({ name: '用户列表' })

const $table = ref(null)
const queryItems = ref({})

const { modalVisible, modalLoading, handleSave, modalForm, modalFormRef, handleEdit } = useCRUD({
  name: '用户',
  doUpdate: api.updateUser,
  refresh: () => $table.value?.handleSearch(),
})

const roleOption = ref([])

onMounted(() => {
  api.getRoleOption().then(res => roleOption.value = res.data)
  $table.value?.handleSearch()
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
    title: '登录方式',
    key: 'login_type',
    width: 40,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: loginTypeMap[row.login_type].tag },
        { default: () => loginTypeMap[row.login_type].name },
      )
    },
  },
  {
    title: '用户角色',
    key: 'role',
    width: 80,
    align: 'center',
    render(row) {
      const roles = row.roles ?? []
      const group = []
      for (let i = 0; i < roles.length; i++)
        group.push(h(NTag, { type: 'info', style: { margin: '2px 3px' } }, { default: () => roles[i].name }))
      return h('span', group)
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
    title: '创建时间',
    key: 'created_at',
    align: 'center',
    width: 70,
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
    title: '上次登录时间',
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
    title: '禁用',
    key: 'is_disable',
    width: 30,
    align: 'center',
    fixed: 'left',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.is_disable,
        loading: !!row.publishing, // 修改 ing 动画
        checkedValue: 1,
        uncheckedValue: 0,
        onUpdateValue: () => handleUpdateDisable(row),
      })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 60,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            onClick: () => {
              // roles => role_ids
              row.role_ids = row.roles.map(e => e.id)
              handleEdit(row)
            },
          },
          {
            default: () => '编辑',
            icon: renderIcon('material-symbols:delete-outline', { size: 16 }),
          },
        ),
      ]
    },
  },
]

// 修改用户禁用状态
async function handleUpdateDisable(row) {
  if (!row.id)
    return
  row.publishing = true
  row.is_disable = row.is_disable === 0 ? 1 : 0
  try {
    await api.updateUserDisable(row)
    $message?.success(row.is_disable ? '已禁用该用户' : '已取消禁用该用户')
    $table.value?.handleSearch()
  }
  catch (err) {
    // 有异常恢复原来的状态
    row.is_disable = row.is_disable === 0 ? 1 : 0
  }
  finally {
    row.publishing = false
  }
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="用户列表">
    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getUsers"
    >
      <template #queryBar>
        <QueryBarItem label="昵称" :label-width="40">
          <NInput
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="请输入用户昵称"
            @keydown.enter="$table?.handleSearch"
          />
        </QueryBarItem>
        <QueryBarItem label="登录方式" :label-width="70" :content-width="180">
          <NSelect
            v-model:value="queryItems.login_type"
            clearable
            filterable
            placeholder="请选择登录方式"
            :options="loginTypeOptions"
            @update:value="$table?.handleSearch"
          />
        </QueryBarItem>
      </template>
    </CrudTable>

    <!-- 新增/编辑 弹窗 -->
    <CrudModal
      v-model:visible="modalVisible"
      title="修改用户"
      :loading="modalLoading"
      @save="handleSave"
    >
      <NForm
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <NFormItem label="用户昵称" path="name">
          <NInput
            v-model:value="modalForm.nickname"
            clearable
            placeholder="请输入用户昵称"
          />
        </NFormItem>
        <NFormItem label="角色" path="role_ids">
          <NCheckboxGroup v-model:value="modalForm.role_ids">
            <NSpace item-style="display: flex;">
              <NCheckbox
                v-for="item in roleOption"
                :key="item.value"
                :value="item.value"
                :label="item.label"
              />
            </NSpace>
          </NCheckboxGroup>
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
