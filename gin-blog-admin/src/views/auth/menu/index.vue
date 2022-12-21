<script setup>
import { NButton, NPopconfirm, NSwitch } from 'naive-ui'
import IconPicker from '@/components/icon/IconPicker.vue'
import TheIcon from '@/components/icon/TheIcon.vue'
import { formatDate, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import api from '@/api'

defineOptions({ name: '菜单管理' })

const $table = $ref(null)
const queryItems = $ref({})

// 表单初始化内容
const initForm = $ref({
  order_num: 1,
  is_hidden: 0,
})

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
  refresh: () => $table?.handleSearch(),
})

onMounted(() => {
  $table?.handleSearch()
})

// 是否展示 "菜单类型"
let showMenuType = $ref(false)

const columns = [
  { title: '菜单名称', key: 'name', width: 80, ellipsis: { tooltip: true } },
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
  { title: '跳转路径', key: 'redirect', width: 80, ellipsis: { tooltip: true } },
  { title: '组件路径', key: 'component', width: 80, ellipsis: { tooltip: true } },
  {
    title: '保活',
    key: 'keep_alive',
    width: 30,
    fixed: 'left',
    render(row) {
      return h('span', row.keep_alive === 1)
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
        loading: !!row.publishing, // 修改 ing 动画
        checkedValue: 1,
        uncheckedValue: 0,
        onUpdateValue: () => handleUpdateHidden(row),
      })
    },
  },
  {
    title: '创建日期',
    key: 'created_at',
    width: 70,
    render(row) {
      return h('span', formatDate(row.created_at))
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
            style: `display: ${row.children ? '' : 'none'};`,
            onClick: () => {
              initForm.is_catelogue = false // 设置非目录(显示组件路径)
              initForm.component = '' // 手动清空组件路径
              initForm.parent_id = row.id // 设置父菜单id
              showMenuType = false
              handleAdd()
            },
          },
          { default: () => '新增', icon: renderIcon('material-symbols:add', { size: 14 }) },
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: () => {
              showMenuType = false
              handleEdit(row)
            },
          },
          { default: () => '编辑', icon: renderIcon('material-symbols:edit-outline', { size: 14 }) },
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete([row.id], false),
          },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'tiny',
                  quaternary: true,
                  type: 'error',
                },
                {
                  default: () => '删除',
                  icon: renderIcon('material-symbols:delete-outline', { size: 14 }),
                },
              ),
            default: () => h('div', {}, '确定删除该菜单吗?'),
          },
        ),
      ]
    },
  },
]

// 修改是否隐藏
async function handleUpdateHidden(row) {
  if (!row.id)
    return
  row.publishing = true
  row.is_hidden = row.is_hidden === 0 ? 1 : 0
  await api.saveOrUpdateMenu(row)
  row.publishing = false
  $message?.success(row.is_hidden ? '已隐藏' : '已取消隐藏')
}

// 新增菜单(可选目录)
function handleClickAdd() {
  showMenuType = true
  initForm.is_catelogue = true // 默认选中"目录"
  initForm.component = 'Layout' // 目录必须是 "Layout", 一级菜单可以是 "Layout"
  initForm.parent_id = 0 // 目录和一级菜单的父id是 0
  handleAdd()
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="菜单管理">
    <template #action>
      <NButton type="primary" @click="handleClickAdd">
        <TheIcon icon="material-symbols:add" :size="18" /> 新建菜单
      </NButton>
    </template>

    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :is-pagination="false"
      :columns="columns"
      :get-data="api.getMenus"
      :single-line="true"
    >
      <template #queryBar>
        <QueryBarItem label="菜单名" :label-width="50">
          <n-input
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="请输入菜单名"
            @keydown.enter="$table?.handleSearch"
          />
        </QueryBarItem>
      </template>
    </CrudTable>

    <!-- 新增/编辑/查看 弹窗 -->
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
      >
        <n-form-item v-if="showMenuType" label="菜单类型" path="type">
          <n-radio-group v-model:value="modalForm.is_catelogue" name="radiogroup">
            <n-space>
              <n-radio :value="true">
                目录
              </n-radio>
              <n-radio :value="false">
                一级菜单
              </n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="菜单名称" path="name">
          <n-input v-model:value="modalForm.name" placeholder="请输入菜单名称" />
        </n-form-item>
        <n-form-item label="菜单图标" path="icon">
          <IconPicker v-model:value="modalForm.icon" w-full />
        </n-form-item>
        <n-form-item v-if="!modalForm.is_catelogue" label="组件路径" path="component">
          <n-input v-model:value="modalForm.component" placeholder="请输入组件路径" />
        </n-form-item>
        <n-form-item label="访问路径" path="path">
          <n-input v-model:value="modalForm.path" placeholder="请输入访问路径" />
        </n-form-item>
        <n-form-item label="跳转路径" path="redirect">
          <n-input
            v-model:value="modalForm.redirect"
            :disabled="modalForm.parent_id !== 0"
            placeholder="只有一级菜单可以设置跳转路径"
          />
        </n-form-item>
        <n-form-item label="显示排序" path="order_num">
          <n-input-number v-model:value="modalForm.order_num" />
        </n-form-item>
        <n-form-item label="是否隐藏" path="is_hidden">
          <NSwitch v-model:value="modalForm.is_hidden" :checked-value="1" :unchecked-value="0" />
        </n-form-item>
        <n-form-item label="KeepAlive" path="keep_alive">
          <NSwitch v-model:value="modalForm.keep_alive" :checked-value="1" :unchecked-value="0" />
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>
