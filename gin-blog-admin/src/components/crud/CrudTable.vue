<script setup>
import { nextTick, reactive, ref } from 'vue'
import { NButton, NDataTable, NSpace } from 'naive-ui'
import { utils, writeFile } from 'xlsx'

const props = defineProps({
  /** 是否不设定列的分割线 */
  singleLine: { type: Boolean, default: false },
  /** true: 后端分页 false: 前端分页 */
  remote: { type: Boolean, default: true },
  /** 是否分页 */
  isPagination: { type: Boolean, default: true },
  /** 表格内容的横向宽度 */
  scrollX: { type: Number, default: 1200 },
  /** 主键 name */
  rowKey: { type: String, default: 'id' },
  /** 需要展示的列 */
  columns: { type: Array, required: true },
  /** queryBar 中的参数 */
  queryItems: {
    type: Object,
    default() { return {} },
  },
  /** 补充参数（可选） */
  extraParams: {
    type: Object,
    default() { return {} },
  },
  /**
   * TODO: 如果想要同时有 url 和 body, 怎么处理
   * 获取数据的请求 API
   */
  getData: {
    type: Function,
    required: true,
  },
})

const emit = defineEmits(['update:queryItems', 'checked', 'dataChange', 'sorterChange'])

const loading = ref(false) // 加载
const selections = ref([]) // 多选的 rowKey
const tableData = ref([]) // 表格数据
const initQuery = { ...props.queryItems }

// 分页配置
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [5, 10, 20],
  onChange: (page) => {
    pagination.page = page
    handleQuery()
  },
  onUpdatePageSize: (pageSize) => {
    pagination.page = 1
    pagination.pageSize = pageSize
    handleQuery()
  },
  prefix({ itemCount }) {
    return `共 ${itemCount} 条`
  },
})

async function handleQuery() {
  selections.value = [] // 重置选中

  try {
    loading.value = true
    let paginationParams = {}
    // 如果非分页模式或者使用前端分页, 则无需传分页参数
    if (props.isPagination && props.remote) {
      paginationParams = {
        page_num: pagination.page,
        page_size: pagination.pageSize,
      }
    }
    const { data } = await props.getData({
      ...props.queryItems,
      ...props.extraParams,
      ...paginationParams,
    })
    tableData.value = data?.page_data || data
    pagination.itemCount = data?.total ?? data.length
  }
  catch (error) {
    tableData.value = []
    pagination.itemCount = 0
  }
  finally {
    emit('dataChange', tableData.value)
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1 // 回到第 1 页
  handleQuery()
}

async function handleReset() {
  const queryItems = { ...props.queryItems } // 重置搜索参数
  for (const key in queryItems) {
    queryItems[key] = null // 注意类型
  }
  emit('update:queryItems', { ...queryItems, ...initQuery })
  await nextTick()
  pagination.page = 1 // 回到第 1 页
  handleQuery()
}

function onPageChange(currentPage) {
  pagination.page = currentPage
  props.remote && handleQuery()
}

function onChecked(rowKeys) {
  selections.value = rowKeys
  // 包含 selection
  if (props.columns.some(item => item.type === 'selection')) {
    emit('checked', rowKeys)
  }
}

function onSorterChange(sorter) {
  emit('sorterChange', sorter)
}

function handleExport(columns = props.columns, data = tableData.value) {
  if (!data?.length) {
    return window.$message.warning('没有数据')
  }
  const columnsData = columns.filter(item => !!item.title && !item.hideInExcel)
  const thKeys = columnsData.map(item => item.key)
  const thData = columnsData.map(item => item.title)
  const trData = data.map(item => thKeys.map(key => item[key]))
  const sheet = utils.aoa_to_sheet([thData, ...trData])
  const workBook = utils.book_new()
  utils.book_append_sheet(workBook, sheet, '数据报表')
  writeFile(workBook, '数据报表.xlsx')
}

defineExpose({
  handleQuery,
  handleSearch,
  handleReset,
  handleExport,
  selections,
  tableData,
})
</script>

<template>
  <div
    v-if="$slots.queryBar"
    class="mb-7 min-h-[60px] flex items-start justify-between border border-gray-200 border-gray-400 rounded-2 border-solid bg-gray-50 p-3.5 dark:bg-black dark:bg-opacity-5"
  >
    <NSpace wrap :size="[35, 15]">
      <slot name="queryBar" />
    </NSpace>
    <div class="flex-shrink-0 space-x-4">
      <NButton ghost type="primary" @click="handleReset">
        <template #icon>
          <i class="i-lucide:rotate-ccw" />
        </template>
        重置
      </NButton>
      <NButton type="primary" @click="handleSearch">
        <template #icon>
          <i class="i-fe:search" />
        </template>
        搜索
      </NButton>
      <!-- TODO: 添加额外的插槽，让用户可以自定义按钮 -->
    </div>
  </div>
  <NDataTable
    :remote="remote"
    :loading="loading"
    :scroll-x="scrollX"
    :columns="columns"
    :data="tableData"
    :row-key="(row) => row[rowKey]"
    :single-line="singleLine"
    :pagination="isPagination ? pagination : false"
    :checked-row-keys="selections"
    @update:checked-row-keys="onChecked"
    @update:page="onPageChange"
    @update:sorter="onSorterChange"
  />
</template>
