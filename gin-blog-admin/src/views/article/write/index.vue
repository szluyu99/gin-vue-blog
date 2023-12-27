<script setup>
import { h, nextTick, onActivated, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { NButton, NDynamicTags, NForm, NFormItem, NInput, NRadio, NRadioGroup, NSelect, NSpace, NSwitch, NTag } from 'naive-ui'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

import CommonPage from '@/components/common/CommonPage.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import UploadOne from '@/components//UploadOne.vue'

import { articleTypeOptions } from '@/assets/config'
import { useTagStore } from '@/store'
import api from '@/api'

defineOptions({ name: '发布文章' })

const route = useRoute()
// const router = useRouter()
const tagStore = useTagStore()

const categoryOptions = ref([]) // 分类选项
const tagOptions = ref([]) // 标签选项
let backTagOptions = [] // 备份标签选项

// 解决同时查看多篇文章, 切换标签不刷新的问题
watch(route, async () => tagStore.reloadTag())

onMounted(async () => {
  fetchData()
})

onActivated(async () => {
  fetchData()
})

async function fetchData() {
  getArticleInfo()
  api.getCategoryOption().then((resp) => {
    categoryOptions.value = resp.data.map(e => ({ value: e.label, label: e.label }))
  })
  api.getTagOption().then((resp) => {
    tagOptions.value = resp.data.map(e => ({ value: e.label, label: e.label }))
    backTagOptions = tagOptions.value
  })
  await nextTick()
}

const formRef = ref(null)
const formModel = ref({
  title: '',
  status: 1, // 发布形式: 默认公开
  is_top: false, // 默认不置顶
  type: 1, // 默认原创
  tag_names: [],
  category_name: '',
})
const btnLoading = ref(false)
const modalVisible = ref(false)
const newTag = ref(null) // 新增标签

// 监听已选标签, 实时更新可选择的标签
watch(() => formModel.value.tag_names, (newVal) => {
  tagOptions.value = backTagOptions.filter(e => !newVal.includes(e.label))
}, { deep: true })

// 根据路由中的 id 参数获取文章信息
async function getArticleInfo() {
  const id = route.params.id // 路由中获取参数

  // 没有 id, 表示是新增文章
  if (!id) {
    formModel.value = { status: 1, is_top: false, title: '', type: 1 }
    return
  }

  // 存在 id, 表示是编辑文章
  window.$loadingBar?.start()
  try {
    const resp = await api.getArticleById(id)
    const { category, tags } = resp.data
    formModel.value = resp.data
    formModel.value.tag_names = tags.map(e => e.name)
    formModel.value.category_name = category.name
    window.$loadingBar?.finish()
  }
  catch (err) {
    window.$loadingBar?.error()
    $message?.error('加载失败')
  }
}

// TODO: 保存草稿
function handleDraft() {
  $message.info('保存草稿开发中')
}

// 发布文章
function handlePublish() {
  if (!formModel.value.title || !formModel.value.title?.trim()) {
    formModel.value.title = formModel.value.title?.trim()
    $message.info('请输入标题')
    return
  }
  modalVisible.value = true
}

// 保存
async function handleSave() {
  formRef.value?.validate(async (err) => {
    if (!err) {
      btnLoading.value = true
      // $message.loading('正在保存...')
      try {
        await api.saveOrUpdateArticle(formModel.value)
        modalVisible.value = false
        $message.success('操作成功!')
        // 关闭当前标签, 并跳转回文章列表
        tagStore.removeTag(route.path)
        // await router.replace({ path: '/article/list', query: { needRefresh: true } })
      }
      catch (err) {
        console.error(err)
      }
      finally {
        btnLoading.value = false
      }
    }
  })
}

const rules = {
  category_name: {
    required: true,
    message: '请选择文章分类',
    trigger: ['blur', 'change'],
  },
  tag_names: {
    required: true,
    message: '请选择文章标签',
  },
}

// 渲染标签
function renderTag(tag, index) {
  return h(
    NTag,
    {
      type: 'info',
      disabled: index > 3,
      closable: true,
      onClose: () => formModel.value.tag_names.splice(index, 1),
    },
    { default: () => tag },
  )
}
</script>

<template>
  <CommonPage :show-header="false" title="写文章">
    <div class="mb-4 flex items-center bg-white space-x-2">
      <NInput
        v-model:value="formModel.title"
        type="text"
        class="mr-5 flex-1 py-1 text-lg color-primary font-bold"
        placeholder="输入文章标题..."
      />
      <NButton ghost type="error" :loading="btnLoading" @click="handleDraft">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:uploading-loop" />
        </template>
        保存草稿
      </NButton>
      <NButton type="error" :loading="btnLoading" @click="handlePublish">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:confirm-circle" />
        </template>
        发布文章
      </NButton>
    </div>

    <!-- TODO: 文件上传 -->
    <MdEditor v-model="formModel.content" style="height: calc(100vh - 245px)" />

    <CrudModal
      v-model:visible="modalVisible"
      title="发布文章"
      :loading="btnLoading"
      show-footer
      @save="handleSave"
    >
      <NForm
        ref="formRef"
        label-placement="left"
        label-align="left"
        :label-width="100"
        :model="formModel"
        :rules="rules"
      >
        <NFormItem label="文章分类" path="category_name">
          <NSelect
            v-model:value="formModel.category_name"
            style="width: 50%"
            clearable filterable tag
            placeholder="关键字搜索，enter 添加"
            :options="categoryOptions"
          />
        </NFormItem>
        <NFormItem label="文章标签" path="tag_names">
          <NDynamicTags
            v-model:value="formModel.tag_names"
            :render-tag="renderTag"
            :max="3"
          >
            <template #input="{ submit, deactivate }">
              <NSelect
                v-model:value="newTag"
                size="small" filterable tag clearable
                :options="tagOptions"
                placeholder="标签名称"
                @update:value="{
                  submit($event);
                  newTag = null;
                }"
                @blur="deactivate"
              >
                <template #action>
                  输入标签名搜索，enter 添加自定义标签
                </template>
              </NSelect>
            </template>
          </NDynamicTags>
        </NFormItem>
        <NFormItem label="文章类型" path="type">
          <NSelect
            v-model:value="formModel.type"
            style="width: 50%"
            placeholder="请选择文章分类"
            :options="articleTypeOptions"
          />
        </NFormItem>
        <!-- <n-form-item label="文章描述" path="desc">
          <n-input
            v-model:value="formModel.desc"
            placeholder="请输入文章描述"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </n-form-item> -->
        <NFormItem
          v-if="(formModel.type === 2 || formModel.type === 3)"
          label="原文地址" path="original_url"
        >
          <NInput
            v-model:value="formModel.original_url"
            type="text"
            placeholder="请填写原文连接"
          />
        </NFormItem>
        <NFormItem label="文章缩略图" path="img">
          <UploadOne
            v-model:preview="formModel.img"
            :width="220"
          />
        </NFormItem>
        <NFormItem label="置顶" path="is_top">
          <NSwitch v-model:value="formModel.is_top" />
        </NFormItem>
        <NFormItem label="发布形式" path="status">
          <NRadioGroup v-model:value="formModel.status" name="radiogroup">
            <NSpace>
              <NRadio :value="1">
                公开
              </NRadio>
              <NRadio :value="2">
                私密
              </NRadio>
            </NSpace>
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>

<style lang="scss" scoped>
.md-preview {
  ul,
  ol {
    list-style: revert;
  }
}
</style>
