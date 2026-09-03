<script setup>
import { NButton, NForm, NFormItem, NInput, NPopconfirm, NSwitch, NTag, NTree } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { formatDate } from '@/utils'

defineOptions({ name: '角色管理' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

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
  initForm: { name: '', label: '', menu_ids: [], resource_ids: [] },
  doCreate: api.saveOrUpdateRole,
  doDelete: api.deleteRole,
  doUpdate: api.saveOrUpdateRole,
  refresh: () => $table.value?.handleSearch(),
})

// 菜单, 资源 跳出菜单的选项不同
const showMenu = ref(true)
const resourceOption = ref([]) // 资源选项
const menuOption = ref([]) // 菜单选项

onMounted(() => {
  $table.value?.handleSearch()
})

// 拉取菜单/资源选项, 已经拉过就不重复请求
async function loadOptions() {
  const tasks = []
  if (!menuOption.value.length) {
    tasks.push(api.getMenuOption().then(resp => (menuOption.value = resp.data ?? [])))
  }
  if (!resourceOption.value.length) {
    tasks.push(api.getResourceOption().then(resp => (resourceOption.value = resp.data ?? [])))
  }
  await Promise.all(tasks)
}

// 新建角色: 后端 SaveRole 已经支持一并写入 role_resource / role_menu,
// 所以这里把两棵树都给出来, 不用再走「先建角色, 再编辑权限」两趟
async function handleAddRole() {
  try {
    await loadOptions()
  }
  catch (err) {
    console.error(err)
    return
  }
  handleAdd()
}

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
            onClick: async () => {
              showMenu.value = true
              await api.getMenuOption().then(resp => (menuOption.value = resp.data))
              handleEdit(row)
            },
          },
          {
            default: () => '菜单权限',
            icon: () => h('i', { class: 'i-material-symbols:edit-outline' }),
          },
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: async () => {
              showMenu.value = false
              await api.getResourceOption().then(resp => (resourceOption.value = resp.data))
              handleEdit(row)
            },
          },
          {
            default: () => '资源权限',
            icon: () => h('i', { class: 'i-ic:baseline-folder-open' }),
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
                  icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
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
  <CommonPage title="角色管理">
    <template #action>
      <NButton type="primary" @click="handleAddRole">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        新建角色
      </NButton>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        批量删除
      </NButton>
    </template>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getRoles"
    >
      <template #queryBar>
        <QueryItem label="角色名" :label-width="50">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入角色名"
            @keydown.enter=" $table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>

    <CrudModal
      v-model:visible="modalVisible"
      :title="modalTitle"
      :loading="modalLoading"
      @save="handleSave"
    >
      <NForm
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
        :disabled="modalAction === 'view'"
      >
        <NFormItem label="角色名" path="name">
          <NInput v-model:value="modalForm.name" placeholder="请输入角色名称" />
        </NFormItem>
        <NFormItem label="角色标签" path="name">
          <NInput v-model:value="modalForm.label" placeholder="请输入角色标签" />
        </NFormItem>
        <!-- 新增时两棵树都给出; 编辑时按入口按钮只展示对应的那棵 -->
        <template v-if="modalAction === 'add'">
          <NFormItem label="菜单权限" path="menu_ids">
            <NTree
              :data="menuOption"
              :checked-keys="modalForm.menu_ids"
              checkable expand-on-click block-line
              @update:checked-keys="(v) => (modalForm.menu_ids = v)"
            />
          </NFormItem>
          <NFormItem label="资源权限" path="resource_ids">
            <NTree
              :data="resourceOption"
              :checked-keys="modalForm.resource_ids"
              block-line checkable expand-on-click cascade accordion
              @update:checked-keys="(v) => (modalForm.resource_ids = v)"
            />
          </NFormItem>
        </template>
        <template v-else-if="modalAction === 'edit'">
          <NFormItem v-if="showMenu" label="菜单权限" path="menu_ids">
            <NTree
              :data="menuOption"
              :checked-keys="modalForm.menu_ids"

              checkable expand-on-click block-line
              @update:checked-keys="(v) => (modalForm.menu_ids = v)"
            />
          </NFormItem>
          <NFormItem v-else label="资源权限" path="resource_ids">
            <NTree
              :data="resourceOption"
              :checked-keys="modalForm.resource_ids"

              block-line checkable expand-on-click cascade accordion
              @update:checked-keys="(v) => (modalForm.resource_ids = v)"
            />
          </NFormItem>
        </template>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
