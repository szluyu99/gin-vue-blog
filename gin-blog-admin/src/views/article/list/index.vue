<script setup>
import { NButton, NImage, NPopconfirm, NSwitch, NTag } from 'naive-ui'
import { convertImgUrl, downloadFile, formatDate, getToken, renderIcon } from '@/utils'
import { useCRUD } from '@/hooks'
import { articleTypeMap, articleTypeOptions } from '@/constant/data'
import api from '@/api'

// 需要 KeepAlive 必须写 name 属性, 并且和 router 中 name 对应
defineOptions({ name: '文章列表' })

const [route, router] = [useRoute(), useRouter()]
const token = getToken()

let categoryOptions = $ref([])
let tagOptions = $ref([])

const $table = $ref(null)
const queryItems = $ref({}) // 条件搜索
const extraParams = $ref({}) // 控制文章状态: 公开, 私密, 草稿箱, 回收站

const { handleDelete } = useCRUD({
  name: '文章',
  doDelete: updateOrDeleteArticles, // 软删除
  refresh: () => $table?.handleSearch(),
})

onMounted(() => {
  api.getCategoryOption().then(res => (categoryOptions = res.data))
  api.getTagOption().then(res => (tagOptions = res.data))
  handleChangeTab('all') // 默认查看全部
})

// ! 切换页面时, 如果是 [写文章] 页面跳转过来, 会携带 needRefresh 参数
onActivated(() => {
  const { needRefresh } = route.query
  needRefresh && ($table?.handleSearch())
})

const columns = [
  { type: 'selection', width: 20, fixed: 'left' },
  {
    title: '文章封面',
    key: 'img',
    width: 55,
    align: 'center',
    render(row) {
      return h(NImage, {
        imgProps: { style: { 'border-radius': '2px', 'height': '100%', 'width': '100%' } },
        src: convertImgUrl(row.img),
        fallbackSrc: 'http://dummyimage.com/400x400',
        showToolbarTooltip: true,
      })
    },
  },
  {
    title: '文章标题',
    key: 'title',
    width: 120,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: '分类',
    key: 'category.name',
    width: 60,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('div', row.category.name || '无')
    },
  },
  {
    title: '标签',
    key: 'tags',
    width: 100,
    align: 'center',
    render(row) {
      const tags = row.tags ?? []
      const group = []
      for (let i = 0; i < tags.length; i++) {
        group.push(
          h(NTag, { type: 'info', style: { margin: '2px 3px' } }, { default: () => tags[i].name }),
        )
      }
      return h('div', group.length ? group : '无')
    },
  },
  {
    title: '浏览量',
    key: 'view_count',
    width: 40,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: '点赞量',
    key: 'like_count',
    width: 40,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: '类型',
    key: 'type',
    width: 50,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: articleTypeMap[row.type].tag },
        { default: () => articleTypeMap[row.type].name },
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
                await $table?.handleSearch()
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
          { onPositiveClick: () => handleDelete([row.id], false) },
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
  extraParams.is_delete === 0
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
  $table?.handleSearch()
}

// 导出文章
async function exportArticles(ids) {
  // 方式一: 前端根据文章内容和标题进行导出
  const list = $table?.tableData.filter(e => ids.includes(e.id))
  for (const item of list)
    downloadFile(item.content, `${item.title}.md`)

  // 方式二: 后端导出返回链接, 前端根据链接下载
  // const res = await api.exportArticles(ids)
  // for (const url of res.data)
  // downloadFile(url)
}

// 切换标签页: [全部, 公开, 私密, 草稿箱, 回收站]
function handleChangeTab(value) {
  switch (value) {
    case 'all':
      extraParams.is_delete = 0
      extraParams.status = null
      break
    case 'public':
      extraParams.is_delete = 0
      extraParams.status = 1
      break
    case 'secret':
      extraParams.is_delete = 0
      extraParams.status = 2
      break
    case 'draft':
      extraParams.is_delete = 0
      extraParams.status = 3
      break
    case 'delete':
      extraParams.is_delete = 1
      extraParams.status = null
      break
  }
  $table?.handleSearch()
}

// 文件上传前检查类型
function beforeUpload(data) {
  if (!data.file.name.endsWith('.md')) {
    $message.error('只能上传 .md 格式的文件，请重新上传')
    return false
  }
  return true
}

// 文件上传后的操作
function afterUpload({ event }) {
  const respStr = (event?.target).response
  const res = JSON.parse(respStr)
  if (res.code === 0) {
    $table?.handleSearch()
    $message.success('文章导入成功！')
  }
  else {
    $message.error('文章导入失败！')
  }
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
        @click="handleDelete($table?.selections)"
      >
        <TheIcon icon="material-symbols:recycling-rounded" :size="18" mr-5 /> 批量删除
      </NButton>
      <NButton
        ml-20 type="info"
        :disabled="!$table?.selections.length"
        @click="exportArticles($table?.selections)"
      >
        <TheIcon icon="mdi:export" :size="18" mr-5 /> 批量导出
      </NButton>
      <div inline-block>
        <n-upload
          action="/api/article/import"
          :headers="{ Authorization: `Bearer ${token}` }"
          :show-file-list="false"
          multiple
          @before-upload="beforeUpload"
          @finish="afterUpload"
        >
          <NButton ml-20 type="success">
            <TheIcon icon="mdi:import" :size="18" mr-5 /> 批量导入
          </NButton>
        </n-upload>
      </div>
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
            :options="articleTypeOptions"
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
