<script setup>
import { NButton, NImage, NPopconfirm, NSwitch, NTag } from 'naive-ui'
import { convertImgUrl, formatDate, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import { artTypeMap, artTypeOptions } from '@/constant/data'
import api from '@/api'

// 需要 KeepAlive 必须写 name 属性, 并且和 router 中 name 对应
defineOptions({ name: '文章列表' })

const [route, router] = [useRoute(), useRouter()]

const categoryOptions = ref([])
const tagOptions = ref([])

const $table = ref(null)
const queryItems = ref({}) // 条件搜索
const extraParams = ref({}) // 控制文章状态: 公开, 私密, 草稿箱, 回收站

const { handleDelete } = useCRUD({
  name: '文章',
  doDelete: updateOrDeleteArticles, // 软删除
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  api.getCategoryOption().then(res => (categoryOptions.value = res.data))
  api.getTagOption().then(res => (tagOptions.value = res.data))
  handleChangeTab('all') // 默认查看全部
})

// ! 切换页面时, 如果是 [写文章] 页面跳转过来, 会携带 needRefresh 参数
onActivated(() => {
  const { needRefresh } = route.query
  needRefresh && ($table.value?.handleSearch())
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  {
    title: '文章封面',
    key: 'img',
    width: 80,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 100,
        'imgProps': { style: { 'border-radius': '2px' } },
        'src': convertImgUrl(row.img),
        'fallback-src': 'http://dummyimage.com/400x400',
        'show-toolbar-tooltip': true,
      })
    },
  },
  { title: '文章标题', key: 'title', width: 120, align: 'center', ellipsis: { tooltip: true } },
  { title: '分类', key: 'category.name', width: 60, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '标签',
    key: 'tags',
    width: 120,
    align: 'center',
    render(row) {
      const tags = row.tags ?? []
      const group = []
      for (let i = 0; i < tags.length; i++) {
        group.push(
          h(NTag, { type: 'info', style: { margin: '2px 3px' } }, { default: () => tags[i].name }),
        )
      }
      return h('div', group)
    },
  },
  {
    title: '类型',
    key: 'type',
    width: 50,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: artTypeMap[row.type].tag },
        { default: () => artTypeMap[row.type].name },
      )
    },
  },
  {
    title: '发布时间',
    key: 'updateDate',
    align: 'center',
    width: 80,
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.updated_at),
          icon: renderIcon('mdi:update', { size: 18 }),
        },
      )
    },
  },
  {
    title: '置顶',
    key: 'is_top',
    width: 50,
    align: 'center',
    fixed: 'left',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.is_top,
        loading: !!row.publishing,
        checkedValue: 1,
        uncheckedValue: 0,
        onUpdateValue: () => handleUpdateTop(row),
      })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        row.is_delete
          ? h(
            NButton,
            {
              size: 'small',
              type: 'success',
              secondary: true,
              onClick: async () => {
                // 软删除恢复
                await api.softDeleteArticle({ ids: [row.id], is_delete: 0 })
                await $table?.value.handleSearch()
              },
            },
            { default: () => '恢复', icon: renderIcon('majesticons:eye-line', { size: 14 }) },
          )
          : h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              secondary: true,
              onClick: () => router.push(`/article/write/${row.id}`), // 携带参数前往 写文章 页面
            },
            { default: () => '查看', icon: renderIcon('majesticons:eye-line', { size: 14 }) },
          ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete(JSON.stringify([row.id]), false) },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'small', type: 'error', style: 'margin-left: 15px;' },
                { default: () => '删除', icon: renderIcon('material-symbols:delete-outline', { size: 14 }) },
              ),
            default: () => h('div', {}, '确定删除该文章吗?'),
          },
        ),
      ]
    },
  },
]

// extraParams 中的 is_delete 为 0 则软删除, is_delete 为 1 则物理删除
function updateOrDeleteArticles(ids) {
  extraParams.value.is_delete === 0
    ? api.softDeleteArticle({ ids: JSON.parse(ids), is_delete: 1 })
    : api.deleteArticle(ids)
}

// 修改文章置顶
async function handleUpdateTop(row) {
  if (!row.id)
    return
  row.publishing = true
  row.is_top = row.is_top === 0 ? 1 : 0
  await api.updateArticleTop(row)
  row.publishing = false
  $message?.success(row.is_top ? '已成功置顶' : '已取消置顶')
  $table.value?.handleSearch()
}

// 切换标签页: [全部, 公开, 私密, 草稿箱, 回收站]
function handleChangeTab(value) {
  switch (value) {
    case 'all':
      extraParams.value.is_delete = 0
      extraParams.value.status = null
      break
    case 'public':
      extraParams.value.is_delete = 0
      extraParams.value.status = 1
      break
    case 'secret':
      extraParams.value.is_delete = 0
      extraParams.value.status = 2
      break
    case 'draft':
      extraParams.value.is_delete = 0
      extraParams.value.status = 3
      break
    case 'delete':
      extraParams.value.is_delete = 1
      extraParams.value.status = null
      break
  }
  $table.value?.handleSearch()
}
</script>

<template>
  <!-- 业务页面 -->
  <CommonPage show-footer title="文章列表">
    <!-- 操作栏 -->
    <template #action>
      <NButton type="primary" @click="router.replace('/article/write')">
        <TheIcon icon="material-symbols:add" :size="18" mr-5 /> 新建文章
      </NButton>
      <NButton
        ml-20
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete(JSON.stringify($table?.selections))"
      >
        <TheIcon icon="material-symbols:recycling-rounded" :size="18" mr-5 /> 批量删除
      </NButton>
      <NButton
        ml-20 type="info"
        :disabled="!$table?.selections.length"
        @click="() => {}"
      >
        <TheIcon icon="mdi:export" :size="18" mr-5 /> 批量导出
      </NButton>
      <NButton ml-20 type="success" @click="() => {}">
        <TheIcon icon="mdi:import" :size="18" mr-5 /> 批量导入
      </NButton>
    </template>
    <!-- 导航栏 -->
    <n-tabs type="line" animated @update:value="handleChangeTab">
      <template #prefix>
        状态
      </template>
      <n-tab-pane name="all" tab="全部" />
      <n-tab-pane name="public" tab="公开" />
      <n-tab-pane name="secret" tab="私密" />
      <n-tab-pane name="draft" tab="草稿箱" />
      <n-tab-pane name="delete" tab="回收站" />
    </n-tabs>
    <!-- 表格 -->
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :extra-params="extraParams"
      :columns="columns"
      :get-data="api.getArticles"
    >
      <template #queryBar>
        <QueryBarItem label="标题" :label-width="40" :content-width="180">
          <n-input
            v-model:value="queryItems.title"
            clearable
            type="text"
            placeholder="请输入标题"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryBarItem>
        <QueryBarItem label="类型" :label-width="40" :content-width="160">
          <n-select
            v-model:value="queryItems.type"
            clearable
            placeholder="请选择文章类型"
            :options="artTypeOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryBarItem>
        <QueryBarItem label="分类" :label-width="40" :content-width="160">
          <n-select
            v-model:value="queryItems.category_id"
            clearable
            filterable
            placeholder="请选择文章分类"
            :options="categoryOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryBarItem>
        <QueryBarItem label="标签" :label-width="40" :content-width="160">
          <n-select
            v-model:value="queryItems.tag_id"
            clearable
            filterable
            placeholder="请选择文章标签"
            :options="tagOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryBarItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
