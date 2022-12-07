<script setup>
// markdown 编辑器
import MdEditor from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

// 单文件上传组件
import UploadOne from '@/components/upload/UploadOne.vue'

import { artTypeOptions } from '@/constant/data'
import { useTagsStore } from '@/store'
import api from '@/api'

defineOptions({ name: '发布文章' })

const [router, route] = [useRouter(), useRoute()]

const categoryOptions = ref([])
const tagOptions = ref([])

onActivated(async () => {
  api.getCategoryOption().then(
    res => (categoryOptions.value = res.data.map(e => ({ value: e.label, label: e.label }))),
  )
  api.getTagOption().then(
    res => (tagOptions.value = res.data.map(e => ({ value: e.label, label: e.label }))),
  )
  // 根据路由中的 id 参数获取文章信息
  await getArticleInfo()
  // TODO: 优化文章显示白屏加载效果
  await nextTick(() => {})
})

const formRef = ref(null)
const formModel = ref({
  status: 1, // 发布形式: 默认公开
  is_top: 0, // 默认不置顶
})
const btnLoading = ref(false)
const modalVisible = ref(false)

// 根据路由中的 id 参数获取文章信息
async function getArticleInfo() {
  const id = route.params.id // 路由中获取参数

  // 没有 id, 表示是新增文章
  if (!id) {
    formModel.value = { status: 1, is_top: 0, title: '' }
    return
  }

  // 存在 id, 表示是编辑文章
  window.$loadingBar?.start()
  $message.loading('加载中...')
  try {
    const res = await api.getArticleById(id)
    formModel.value = res.data

    window.$loadingBar?.finish()
    $message?.success('加载成功')
  }
  catch (err) {
    window.$loadingBar?.error()
    $message?.error('加载失败')
  }
}

// 删除标签
function removeTag(name) {
  formModel.value.tag_names = formModel.value.tag_names.filter(e => e !== name)
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
    if (err)
      return
    btnLoading.value = true
    $message.loading('正在保存...')
    try {
      await api.saveOrUpdateArticle(formModel.value)
      modalVisible.value = false
      $message.success('操作成功!')
      // 关闭当前标签, 并跳转回文章列表
      useTagsStore().removeTag(route.path)
      await router.replace({ path: '/article/list', query: { needRefresh: true } })
    }
    catch (err) {
      console.error(err)
    }
    finally {
      btnLoading.value = false
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
</script>

<template>
  <CommonPage :show-header="false" show-footer title="写文章">
    <div mb-15 flex items-center bg-white>
      <n-input
        v-model:value="formModel.title"
        type="text"
        class="flex-1 py-5 mr-20 text-18 font-bold color-primary"
        placeholder="输入文章标题..."
      />
      <n-button ghost type="error" mr-10 :loading="btnLoading" @click="handleDraft">
        <TheIcon v-if="!btnLoading" icon="line-md:uploading-loop" mr-5 :size="18" />
        保存草稿
      </n-button>
      <n-button type="error" :loading="btnLoading" @click="handlePublish">
        <TheIcon v-if="!btnLoading" icon="line-md:confirm-circle" mr-5 :size="18" />
        发布文章
      </n-button>
    </div>
    <!-- markdown 编辑器 -->
    <MdEditor v-model="formModel.content" style="height: calc(100vh - 305px)" />

    <CrudModal
      v-model:visible="modalVisible"
      title="发布文章"
      :loading="btnLoading"
      show-footer
      @on-save="handleSave"
    >
      <!-- 表单 -->
      <n-form
        ref="formRef"
        label-placement="left"
        label-align="left"
        :label-width="100"
        :model="formModel"
        :rules="rules"
      >
        <n-form-item label="文章分类" path="category_name">
          <n-select
            v-model:value="formModel.category_name"
            style="width: 50%"
            clearable
            tag
            filterable
            placeholder="关键字搜索，enter 添加"
            :options="categoryOptions"
          />
        </n-form-item>
        <n-form-item label="文章标签" path="tag_names">
          <n-tag
            v-for="tag of formModel.tag_names" :key="tag"
            mr-8 closable type="info"
            @close="removeTag(tag)"
          >
            {{ tag }}
          </n-tag>
          <n-popover
            v-if="(formModel.tag_names.length < 3)"
            trigger="click"
            placement="bottom-start"
          >
            <template #trigger>
              <n-button secondary type="success" size="small">
                添加标签
              </n-button>
            </template>
            <div w-340>
              <n-select
                v-model:value="formModel.tag_names"
                filterable multiple tag
                :options="tagOptions"
                placeholder="输入标签名搜索，enter 添加自定义标签"
              >
                <template #action>
                  请输入标签名搜索，enter 可添加自定义标签
                </template>
              </n-select>
            </div>
          </n-popover>
        </n-form-item>
        <n-form-item label="文章类型" path="type">
          <n-select
            v-model:value="formModel.type"
            style="width: 50%"
            placeholder="请选择文章分类"
            :options="artTypeOptions"
          />
        </n-form-item>
        <!-- <n-form-item label="文章描述" path="desc">
          <n-input
            v-model:value="formModel.desc"
            placeholder="请输入文章描述"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </n-form-item> -->
        <n-form-item
          v-if="(formModel.type === 2 || formModel.type === 3)"
          label="原文地址" path="original_url"
        >
          <n-input
            v-model:value="formModel.original_url"
            type="text"
            placeholder="请填写原文连接"
          />
        </n-form-item>
        <n-form-item label="文章缩略图" path="img">
          <UploadOne
            v-model:preview="formModel.img"
            :width="180"
            @finish="val => (formModel.img = val)"
          />
        </n-form-item>
        <n-form-item label="置顶" path="is_top">
          <n-switch v-model:value="formModel.is_top" :checked-value="1" :unchecked-value="0" />
        </n-form-item>
        <n-form-item label="发布形式" path="status">
          <n-radio-group v-model:value="formModel.status" name="radiogroup">
            <n-space>
              <n-radio :value="1">
                公开
              </n-radio>
              <n-radio :value="2">
                私密
              </n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
      </n-form>
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
