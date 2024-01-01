<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NPopconfirm, NRadio, NRadioGroup, NSpace, NSwitch, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'
import IconPicker from '@/components/icon/IconPicker.vue'
import TheIcon from '@/components/icon/TheIcon.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: '菜单管理' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

const initForm = {
  order_num: 1,
  is_hidden: false, // 是否隐藏
  is_catalogue: false, // 是否目录
  is_external: false, // 是否外链
  keep_alive: false,
  icon: 'mdi-account',
  order_num: 1,
  name: '',
  path: '',
  redirect: '',
  component: '',
  parent_id: 0,
}

const {
  modalVisible,
  modalTitle,
  modalLoading,
  handleAdd,
  handleDelete,
  handleEdit,
  handleSave,
  modalForm,
  modalFormRef,
} = useCRUD({
  name: '菜单',
  initForm,
  doCreate: api.saveOrUpdateMenu,
  doDelete: api.deleteMenu,
  doUpdate: api.saveOrUpdateMenu,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  {
    title: '菜单名称',
    key: 'name',
    width: 100,
    render: (row) => {
      const groups = []
      groups.push(h('span', row.name))

      if (row.parent_id === 0) {
        groups.push(
          h(
            NTag,
            { type: row.is_catalogue ? 'info' : 'success', class: 'ml-1.5' },
            { default: () => row.is_catalogue ? '目录' : '一级菜单' },
          ),
        )
      }
      else {
        groups.push(
          h(
            NTag,
            { type: 'default', class: 'ml-1.5' },
            { default: () => '子菜单' },
          ),
        )
      }

      if (row.is_external) {
        groups.push(
          h(
            NTag,
            { type: 'warning', class: 'ml-1.5' },
            { default: () => '外链' },
          ),
        )
      }

      return groups
    },
  },
  {
    title: '图标',
    key: 'icon',
    width: 30,
    render(row) {
      return h(TheIcon, { icon: row.icon, size: 20 })
    },
  },
  { title: '排序', key: 'order_num', width: 30, ellipsis: { tooltip: true } },
  { title: '访问路径', key: 'path', width: 60, ellipsis: { tooltip: true } },
  {
    title: '跳转路径',
    key: 'redirect',
    width: 80,
    render(row) {
      if (row.parent_id === 0 && !row.is_catalogue) {
        return h('span', row.redirect)
      }
      return h('span', '-')
    },
  },
  {
    title: '组件路径',
    key: 'component',
    width: 80,
    render(row) {
      if (!row.is_catalogue) {
        return h('span', row.component)
      }
      return h('span', '-')
    },
  },
  {
    title: '保活',
    key: 'keep_alive',
    width: 30,
    fixed: 'left',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.keep_alive,
        loading: !!row.publishing,
        onUpdateValue: () => handleUpdateKeepAlive(row),
      })
    },
  },
  {
    title: '隐藏',
    key: 'is_hidden',
    width: 30,
    fixed: 'left',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.is_hidden,
        loading: !!row.publishing,
        onUpdateValue: () => handleUpdateHidden(row),
      })
    },
  },
  {
    title: '更新日期',
    key: 'updated_at',
    width: 70,
    render(row) {
      return h('span', formatDate(row.updated_at))
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 115,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'primary',
            style: `display: ${!row.is_catalogue && row.parent_id === 0 ? '' : 'none'};`,
            onClick: () => {
              initForm.component = '' // 手动清空组件路径
              initForm.parent_id = row.id // 设置父菜单id
              initForm.is_catalogue = false
              handleAdd()
            },
          },
          { default: () => '新增', icon: () => h('i', { class: 'i-material-symbols:add' }) },
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: () => {
              handleEdit(row)
            },
          },
          { default: () => '编辑', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(row.id, false),
          },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'tiny', quaternary: true, type: 'error' },
                {
                  default: () => '删除',
                  icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
                },
              ),
            default: () => h('div', {}, '确定删除该菜单吗?'),
          },
        ),
      ]
    },
  },
]

async function handleUpdateKeepAlive(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.keep_alive = !row.keep_alive
  try {
    await api.saveOrUpdateMenu(row)
    $message?.success(row.keep_alive ? '已保活' : '已取消保活')
  }
  catch (err) {
    row.keep_alive = !row.keep_alive
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

async function handleUpdateHidden(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.is_hidden = !row.is_hidden
  try {
    await api.saveOrUpdateMenu(row)
    $message?.success(row.is_hidden ? '已隐藏' : '已取消隐藏')
  }
  catch (err) {
    row.is_hidden = !row.is_hidden
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

// 新增菜单(可选目录)
function handleClickAdd() {
  initForm.is_catalogue = true // 默认选中"目录"
  initForm.component = 'Layout' // 目录必须是 "Layout", 一级菜单可以是 "Layout"
  initForm.parent_id = 0 // 目录和一级菜单的父id是 0
  handleAdd()
}
</script>

<template>
  <CommonPage title="菜单管理">
    <template #action>
      <NButton type="primary" @click="handleClickAdd">
        <template #icon>
          <span class="i-material-symbols:add" />
        </template>
        新建菜单
      </NButton>
    </template>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :is-pagination="false"
      :columns="columns"
      :get-data="api.getMenus"
      :single-line="true"
    >
      <template #queryBar>
        <QueryItem label="菜单名" :label-width="50">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入菜单名"
            @keydown.enter="$table?.handleSearch()"
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
      >
        <NFormItem v-if="modalForm.parent_id === 0" label="菜单类型" path="type">
          <NRadioGroup v-model:value="modalForm.is_catalogue" name="radiogroup">
            <NSpace>
              <NRadio :value="true">
                目录
              </NRadio>
              <NRadio :value="false">
                一级菜单
              </NRadio>
            </NSpace>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="菜单名称" path="name">
          <NInput v-model:value="modalForm.name" placeholder="请输入菜单名称" />
        </NFormItem>
        <NFormItem label="菜单图标" path="icon">
          <IconPicker v-model:value="modalForm.icon" />
        </NFormItem>
        <NFormItem v-if="!modalForm.is_catalogue" label="组件路径" path="component">
          <NInput v-model:value="modalForm.component" placeholder="请输入组件路径" />
        </NFormItem>
        <NFormItem label="访问路径" path="path">
          <NInput v-model:value="modalForm.path" placeholder="请输入访问路径" />
        </NFormItem>
        <NFormItem v-if="!modalForm.is_catalogue" label="跳转路径" path="redirect">
          <NInput
            v-model:value="modalForm.redirect"
            :disabled="modalForm.parent_id !== 0"
            placeholder="只有一级菜单可以设置跳转路径"
          />
        </NFormItem>
        <NFormItem label="显示排序" path="order_num">
          <NInputNumber v-model:value="modalForm.order_num" />
        </NFormItem>
        <NFormItem label="是否隐藏" path="is_hidden">
          <NSwitch v-model:value="modalForm.is_hidden" />
        </NFormItem>
        <NFormItem label="是否外链" path="is_external">
          <NSwitch v-model:value="modalForm.is_external" />
        </NFormItem>
        <NFormItem label="KeepAlive" path="keep_alive">
          <NSwitch v-model:value="modalForm.keep_alive" />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
