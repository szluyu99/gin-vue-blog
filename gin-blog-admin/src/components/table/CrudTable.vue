<script setup>
const props = defineProps({
  /** 是否不设定列的分割线 **/
  singleLine: {
    type: Boolean,
    default: false,
  },
  /** true: 后端分页  false: 前端分页 */
  remote: {
    type: Boolean,
    default: true,
  },
  /** 是否分页 */
  isPagination: {
    type: Boolean,
    default: true,
  },
  /** 表格内容的横向宽度 */
  scrollX: {
    type: Number,
    default: 1200,
  },
  rowKey: {
    type: String,
    default: 'id',
  },
  /** 需要展示的列 */
  columns: {
    type: Array,
    required: true,
  },
  /** 选中的列 **/
  // selections: {
  //   type: Array,
  //   default() {
  //     return []
  //   },
  // },
  /** queryBar 中的参数 */
  queryItems: {
    type: Object,
    default() {
      return {}
    },
  },
  /** 补充参数（可选） */
  extraParams: {
    type: Object,
    default() {
      return {}
    },
  },
  /**
   * ! 约定接口入参出参
   * * 分页模式需约定分页接口入参
   *    @param {number} pageSize 分页参数：一页展示多少条，默认 10
   *    @param {number} pageNum   分页参数：页码，默认 1
   * * 需约定接口出参
   *    @param {number} pageData 分页模式必须, 非分页模式如果没有 pageData 则取上一层 data
   *    @param {number} total    分页模式必须，非分页模式如果没有 total 则取上一层 data.length
   */
  getData: {
    type: Function,
    required: true,
  },
})

const emit = defineEmits(['update:queryItems', 'onChecked'])

const selections = ref([]) // 多选

const loading = ref(false)
const initQuery = { ...props.queryItems }
const tableData = ref([])
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
    // 如果非分页模式或者使用前端分页,则无需传分页参数
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
    tableData.value = data?.pageData || data
    pagination.itemCount = data?.total ?? data.length
  }
  catch (error) {
    tableData.value = []
    pagination.itemCount = 0
  }
  finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1 // 搜索回到第 1 页
  handleQuery()
}

async function handleReset() {
  const queryItems = { ...props.queryItems }
  for (const key in queryItems)
    queryItems[key] = null // 注意类型

  emit('update:queryItems', { ...queryItems, ...initQuery })
  await nextTick()
  pagination.page = 1
  handleQuery()
}

function onPageChange(currentPage) {
  pagination.page = currentPage
  // 后端分页
  if (props.remote)
    handleQuery()
}

function onChecked(rowKeys) {
  selections.value = rowKeys
  // 包含 selection
  if (props.columns.some(item => item.type === 'selection'))
    emit('onChecked', rowKeys)
}

defineExpose({
  handleSearch,
  handleReset,
  selections,
})
</script>

<template>
  <QueryBar
    v-if="$slots.queryBar"
    mb-30
    @search="handleSearch"
    @reset="handleReset"
  >
    <slot name="queryBar" />
  </QueryBar>

  <n-data-table
    :remote="remote"
    :loading="loading"
    :scroll-x="scrollX"
    :columns="columns"
    :data="tableData"
    :row-key="(row) => row[rowKey]"
    :bordered="true"
    :single-line="singleLine"
    :pagination="isPagination ? pagination : false"
    :checked-row-keys="selections"
    @update:checked-row-keys="onChecked"
    @update:page="onPageChange"
  />
</template>
