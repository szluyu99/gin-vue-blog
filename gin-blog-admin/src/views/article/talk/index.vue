<script setup>
import { NButton, NForm, NFormItem, NInput, NPopconfirm, NRadio, NRadioGroup, NTabPane, NTabs, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'
import { useCRUD } from '@/composables'
import { formatDate } from '@/utils'

defineOptions({ name: '说说管理' })

const $table = ref(null)
const extraParams = ref({ status: 0 }) // 0 表示全部

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
  name: '说说',
  initForm: { content: '', status: 1, is_top: false },
  doCreate: api.saveOrUpdateTalk,
  doDelete: api.deleteTalks,
  doUpdate: api.saveOrUpdateTalk,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  {
    title: '内容',
    key: 'content',
    minWidth: 220,
    ellipsis: { tooltip: true },
  },
  {
    title: '发布者',
    key: 'nickname',
    width: 80,
    align: 'center',
    render(row) {
      return h('span', row.user?.info?.nickname || '-')
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 60,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: row.status === 1 ? 'success' : 'warning' },
        { default: () => (row.status === 1 ? '公开' : '私密') },
      )
    },
  },
  {
    title: '置顶',
    key: 'is_top',
    width: 50,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: row.is_top ? 'info' : 'default', size: 'small' },
        { default: () => (row.is_top ? '是' : '否') },
      )
    },
  },
  {
    title: '发布时间',
    key: 'created_at',
    width: 110,
    align: 'center',
    render(row) {
      return h('span', formatDate(row.created_at))
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 130,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 10px;',
            onClick: () => handleEdit(row),
          },
          { default: () => '编辑', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'small', type: 'error', style: 'margin-left: 10px;' },
                { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
              ),
            default: () => h('div', {}, '确定删除这条说说吗? 它下面的评论会一起删除。'),
          },
        ),
      ]
    },
  },
]

// 切换标签页: [全部, 公开, 私密]
function handleChangeTab(value) {
  extraParams.value.status = Number(value)
  $table.value?.handleSearch()
}
</script>

<template>
  <CommonPage title="说说管理">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <template #icon>
          <span class="i-material-symbols:add" />
        </template>
        新增说说
      </NButton>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <span class="i-material-symbols:recycling-rounded" />
        </template>
        批量删除
      </NButton>
    </template>

    <NTabs type="line" animated @update:value="handleChangeTab">
      <template #prefix>
        状态
      </template>
      <NTabPane name="0" tab="全部" />
      <NTabPane name="1" tab="公开" />
      <NTabPane name="2" tab="私密" />
    </NTabs>

    <CrudTable
      ref="$table"
      :columns="columns"
      :extra-params="extraParams"
      :get-data="api.getTalks"
    />

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
        <NFormItem
          label="内容"
          path="content"
          :rule="{ required: true, message: '请输入说说内容', trigger: ['input', 'blur'] }"
        >
          <NInput
            v-model:value="modalForm.content"
            type="textarea"
            placeholder="说点什么..."
            :autosize="{ minRows: 4, maxRows: 8 }"
            maxlength="1000"
            show-count
          />
        </NFormItem>
        <NFormItem label="状态" path="status">
          <NRadioGroup v-model:value="modalForm.status">
            <NRadio :value="1">
              公开
            </NRadio>
            <NRadio :value="2">
              私密
            </NRadio>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="置顶" path="is_top">
          <NRadioGroup v-model:value="modalForm.is_top">
            <NRadio :value="false">
              不置顶
            </NRadio>
            <NRadio :value="true">
              置顶
            </NRadio>
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
