<script setup>
import { NButton, NPopconfirm, NSwitch, NTag } from 'naive-ui'
import { formatDate, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import api from '@/api'

defineOptions({ name: '角色管理' })

const $table = $ref(null)
const queryItems = $ref({})

const {
  modalVisible,
  modalAction,
  modalTitle,
  modalLoading,
  handleAdd,
  handleDelete,
  handleEdit,
  handleSave,
  modalForm,
  modalFormRef,
} = useCRUD({
  name: '角色',
  initForm: {},
  doCreate: api.saveOrUpdateRole,
  doDelete: api.deleteRole,
  doUpdate: api.saveOrUpdateRole,
  refresh: () => $table?.handleSearch(),
})

// 菜单, 资源 跳出菜单的选项不同
let showMenu = $ref(true)
let resourceOption = $ref([]) // 资源选项
let menuOption = $ref([]) // 菜单选项

onMounted(() => {
  $table?.handleSearch()
  api.getResourceOption().then(res => (resourceOption = res.data))
  api.getMenuOption().then(res => (menuOption = res.data))
})

const columns = [
  {
    type: 'selection',
    width: 15,
    fixed: 'left',
  },
  {
    title: '角色名',
    key: 'name',
    width: 80,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: '角色标签',
    key: 'label',
    width: 80,
    align: 'center',
    render(row) {
      return h(NTag, { type: 'info' }, { default: () => row.label })
    },
  },
  {
    title: '创建日期',
    key: 'created_at',
    width: 60,
    align: 'center',
    render(row) {
      return h('span', formatDate(row.created_at))
    },
  },
  {
    title: '是否禁用',
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
        onUpdateValue: () => $message.info('这个功能暂时还不支持~'),
      })
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
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: () => {
              showMenu = true
              api.getMenuOption().then(res => (menuOption = res.data))
              handleEdit(row)
            },
          },
          {
            default: () => '菜单权限',
            icon: renderIcon('material-symbols:edit-outline', { size: 14 }),
          },
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: () => {
              showMenu = false
              api.getResourceOption().then(res => (resourceOption = res.data))
              handleEdit(row)
            },
          },
          {
            default: () => '资源权限',
            icon: renderIcon('ic:baseline-folder-open', { size: 14 }),
          },
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete([row.id], false),
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
            default: () => h('div', {}, '确定删除该角色吗?'),
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="角色管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <TheIcon icon="material-symbols:add" :size="18" /> 新建角色
      </NButton>
      <NButton
        ml-20
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <TheIcon icon="material-symbols:playlist-remove" :size="18" /> 批量删除
      </NButton>
    </template>

    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getRoles"
    >
      <template #queryBar>
        <QueryBarItem label="角色名" :label-width="50">
          <n-input
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入角色名"
            @keydown.enter=" $table?.handleSearch()"
          />
        </QueryBarItem>
      </template>
    </CrudTable>

    <!-- 新增/编辑 弹窗 -->
    <CrudModal
      v-model:visible="modalVisible"
      :title="modalTitle"
      :loading="modalLoading"
      @on-save="handleSave"
    >
      <!-- 表单 -->
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
        :disabled="modalAction === 'view'"
      >
        <n-form-item label="角色名" path="name">
          <n-input v-model:value="modalForm.name" placeholder="请输入角色名称" />
        </n-form-item>
        <n-form-item label="角色标签" path="name">
          <n-input v-model:value="modalForm.label" placeholder="请输入角色标签" />
        </n-form-item>
        <n-form-item v-if="showMenu" label="菜单权限" path="menu_ids">
          <n-tree
            :data="menuOption"
            :checked-keys="modalForm.menu_ids"
            block-line
            checkable
            expand-on-click
            @update:checked-keys="(v) => (modalForm.menu_ids = v)"
          />
        </n-form-item>
        <n-form-item v-else label="资源权限" path="resource_ids">
          <n-tree
            :data="resourceOption"
            :checked-keys="modalForm.resource_ids"
            block-line
            checkable
            cascade
            accordion
            expand-on-click
            @update:checked-keys="(v) => (modalForm.resource_ids = v)"
          />
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>
