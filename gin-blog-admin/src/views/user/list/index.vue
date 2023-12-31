<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NCheckbox, NCheckboxGroup, NForm, NFormItem, NImage, NInput, NSelect, NSpace, NSwitch, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { loginTypeMap, loginTypeOptions } from '@/assets/config'
import { convertImgUrl, formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: '用户列表' })

const $table = ref(null)
const queryItems = ref({
  username: '',
  nickname: '',
  login_type: null,
})

const {
  modalVisible,
  modalLoading,
  handleSave,
  handleEdit,
  modalForm,
  modalFormRef,
} = useCRUD({
  name: '用户',
  doUpdate: api.updateUser,
  refresh: () => $table.value?.handleSearch(),
})

const roleOptions = ref([])

onMounted(() => {
  api.getRoleOption().then(resp => roleOptions.value = resp.data)
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
        'imgProps': { style: { 'border-radius': '3px' } },
        'src': convertImgUrl(row.info?.avatar),
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
      return h('span', row.info?.nickname)
    },
  },
  {
    title: '登录方式',
    key: 'login_type',
    width: 40,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: loginTypeMap[row.login_type]?.tag },
        { default: () => loginTypeMap[row.login_type]?.name || '未知' },
      )
    },
  },
  {
    title: '用户角色',
    key: 'role',
    width: 80,
    align: 'center',
    render(row) {
      if (row.is_super) {
        return h(NTag, { type: 'error' }, { default: () => '超级管理员' })
      }
      const roles = row.roles ?? []
      const groups = []
      for (let i = 0; i < roles.length; i++) {
        groups.push(h(NTag, { type: 'info', style: { margin: '2px 3px' } }, { default: () => roles[i].name }))
      }
      return h('span', groups.length ? groups : '无')
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
          icon: () => h('i', { class: 'i-mdi:update' }),
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
          icon: () => h('i', { class: 'i-mdi:update' }),
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
        loading: !!row.publishing,
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
              row.nickname = row.info?.nickname
              // roles => role_ids
              row.role_ids = row.roles.map(e => e.id)
              handleEdit(row)
            },
          },
          {
            default: () => '编辑',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
          },
        ),
      ]
    },
  },
]

// 修改用户禁用状态
async function handleUpdateDisable(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.is_disable = !row.is_disable
  try {
    await api.updateUserDisable(row.id, row.is_disable)
    $message?.success(row.is_disable ? '已禁用该用户' : '已取消禁用该用户')
    $table.value?.handleSearch()
  }
  catch (err) {
    row.is_disable = !row.is_disable
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}
</script>

<template>
  <CommonPage title="用户列表">
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getUsers"
    >
      <template #queryBar>
        <QueryItem label="昵称" :label-width="40" :content-width="160">
          <NInput
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="请输入昵称"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="用户名" :label-width="60" :content-width="160">
          <NInput
            v-model:value="queryItems.username"
            clearable
            type="text"
            placeholder="请输入用户名"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="登录方式" :label-width="70" :content-width="160">
          <NSelect
            v-model:value="queryItems.login_type"
            clearable
            filterable
            placeholder="请选择登录方式"
            :options="loginTypeOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>

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
                v-for="item in roleOptions"
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
