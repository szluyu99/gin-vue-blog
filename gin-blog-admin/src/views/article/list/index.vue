<script setup>
import { defineOptions, h, onActivated, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NImage, NInput, NPopconfirm, NSelect, NSwitch, NTabPane, NTabs, NTag, NUpload } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { convertImgUrl, formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import { articleTypeMap, articleTypeOptions } from '@/assets/config'
import api from '@/api'

// 需要 KeepAlive 必须写 name 属性, 并且和 router 中 name 对应
defineOptions({ name: '文章列表' })

const route = useRoute()
const router = useRouter()

const categoryOptions = ref([])
const tagOptions = ref([])

const $table = ref(null)

const queryItems = ref({
  title: '', // 标题
  type: null, // 类型
  category_id: null, // 分类
  tag_id: null, // 标签
})

const extraParams = ref({
  is_delete: null, // 未删除 | 回收站
  status: null, // null-all, 1-公开, 2-私密, 3-草稿
})

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
        { type: articleTypeMap[row.type]?.tag },
        { default: () => articleTypeMap[row.type]?.name },
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
          icon: () => h('i', { class: 'i-mdi:update' }),
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
                await api.softDeleteArticle([row.id], false)
                await $table.value?.handleSearch()
              },
            },
            { default: () => '恢复', icon: () => h('i', { class: 'i-majesticons:eye-line' }) },
          )
          : h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              secondary: true,
              onClick: () => router.push(`/article/write/${row.id}`), // 携带参数前往 写文章 页面
            },
            { default: () => '查看', icon: () => h('i', { class: 'i-majesticons:eye-line' }) },
          ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'small', type: 'error', style: 'margin-left: 15px;' },
                { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
              ),
            default: () => h('div', {}, '确定删除该文章吗?'),
          },
        ),
      ]
    },
  },
]

function updateOrDeleteArticles(ids) {
  extraParams.value.is_delete
    ? api.deleteArticle(ids)
    : api.softDeleteArticle(JSON.parse(ids), true)
}

// 修改文章置顶
async function handleUpdateTop(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.is_top = !row.is_top
  try {
    await api.updateArticleTop(row.id, row.is_top)
    $message?.success(row.is_top ? '已成功置顶' : '已取消置顶')
    $table.value?.handleSearch()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

// 导出文章
async function exportArticles(ids) {
  // 方式一: 前端根据文章内容和标题进行导出
  const list = $table.value?.tableData.filter(e => ids.includes(e.id))
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
    $table.value?.handleSearch()
    $message.success('文章导入成功！')
  }
  else {
    $message.error('文章导入失败！')
  }
}

function downloadFile(content, fileName) {
  const aEle = document.createElement('a') // 创建下载链接
  aEle.download = fileName // 设置下载的名称
  aEle.style.display = 'none'// 隐藏的可下载链接
  // 字符内容转变成 blob 地址
  const blob = new Blob([content])
  aEle.href = URL.createObjectURL(blob)
  // 绑定点击时间
  document.body.appendChild(aEle)
  aEle.click()
  // 然后移除
  document.body.removeChild(aEle)
}
</script>

<template>
  <CommonPage title="文章列表">
    <template #action>
      <NButton type="primary" @click="$router.replace('/article/write')">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        新建文章
      </NButton>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <i class="i-material-symbols:recycling-rounded" />
        </template>
        批量删除
      </NButton>
      <NButton
        type="info"
        :disabled="!$table?.selections.length"
        @click="exportArticles($table?.selections)"
      >
        <template #icon>
          <i class="i-mdi:export" />
        </template>
        批量导出
      </NButton>
      <div class="inline-block">
        <NUpload
          action="/api/article/import"
          :show-file-list="false"
          multiple
          @before-upload="beforeUpload"
          @finish="afterUpload"
        >
          <NButton type="success">
            <template #icon>
              <i class="i-mdi:import" />
            </template>
            批量导入
          </NButton>
        </NUpload>
      </div>
    </template>

    <NTabs type="line" animated @update:value="handleChangeTab">
      <template #prefix>
        状态
      </template>
      <NTabPane name="all" tab="全部" />
      <NTabPane name="public" tab="公开" />
      <NTabPane name="secret" tab="私密" />
      <NTabPane name="draft" tab="草稿箱" />
      <NTabPane name="delete" tab="回收站" />
    </NTabs>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :extra-params="extraParams"
      :columns="columns"
      :get-data="api.getArticles"
    >
      <template #queryBar>
        <QueryItem label="标题" :label-width="40" :content-width="180">
          <NInput
            v-model:value="queryItems.title"
            clearable
            type="text"
            placeholder="请输入标题"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="类型" :label-width="40" :content-width="160">
          <NSelect
            v-model:value="queryItems.type"
            clearable
            placeholder="请选择文章类型"
            :options="articleTypeOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="分类" :label-width="40" :content-width="160">
          <NSelect
            v-model:value="queryItems.category_id"
            clearable
            filterable
            placeholder="请选择文章分类"
            :options="categoryOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="标签" :label-width="40" :content-width="160">
          <NSelect
            v-model:value="queryItems.tag_id"
            clearable
            filterable
            placeholder="请选择文章标签"
            :options="tagOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
