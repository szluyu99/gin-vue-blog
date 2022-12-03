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
            filterable
            placeholder="请选择文章分类"
            :options="categoryOptions"
          />
        </n-form-item>
        <n-form-item label="文章标签">
          <n-popselect
            v-model:value="formModel.tag_names"
            multiple
            scrollable
            :options="tagOptions"
          >
            <n-button>
              {{ formModel.tag_names?.length ? formModel.tag_names : '请选择标签' }}
            </n-button>
            <template #empty> 没啥看的，这里是空的 </template>
            <template #action>
              <n-input
                type="text"
                placeholder="enter 可以添加自定义标签"
                @keydown.enter="addNewCategoryName"
              />
            </template>
          </n-popselect>
        </n-form-item>
        <n-form-item label="文章类型" path="type">
          <n-select
            v-model:value="formModel.type"
            style="width: 50%"
            clearable
            filterable
            placeholder="请选择文章分类"
            :options="artTypeOptions"
          />
        </n-form-item>
        <n-form-item label="文章描述" path="desc">
          <n-input
            v-model:value="formModel.desc"
            placeholder="请输入文章描述"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </n-form-item>
        <n-form-item label="文章缩略图" path="img">
          <n-upload
            action="/api/upload"
            list-type="image-card"
            :default-file-list="imgList"
            :max="1"
            @finish="handleImgUpload"
          >
            点击上传
          </n-upload>
        </n-form-item>
        <n-form-item label="置顶" path="is_top">
          <n-switch v-model:value="formModel.is_top" :checked-value="1" :unchecked-value="0" />
        </n-form-item>
        <n-form-item label="发布形式" path="status">
          <n-radio-group v-model:value="formModel.status" name="radiogroup">
            <n-space>
              <n-radio :value="1"> 公开 </n-radio>
              <n-radio :value="2"> 私密 </n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
      </n-form>
    </CrudModal>
  </CommonPage>
</template>

<script setup>
// markdown 编辑器
import MdEditor from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

import { artTypeOptions } from '@/constant/data'
import { useTagsStore } from '@/store'

import { useCategoryApi, useArticleApi, useTagApi } from '@/api'
const { getCategoryOption } = useCategoryApi()
const { saveOrUpdateArticle, getArticleById } = useArticleApi()
const { getTagOption } = useTagApi()

const router = useRouter()
const route = useRoute()

// 需要 KeepAlive 必须写 name 属性
defineOptions({ name: '发布文章' })

const categoryOptions = ref([])
const tagOptions = ref([])

onActivated(async () => {
  getCategoryOption().then(
    (res) => (categoryOptions.value = res.data.map((e) => ({ value: e.label, label: e.label })))
  )
  getTagOption().then(
    (res) => (tagOptions.value = res.data.map((e) => ({ value: e.label, label: e.label })))
  )
  // 根据路由中的 id 参数获取文章信息
  await getArticleInfo()
  // TODO: 优化文章显示白屏加载效果
  await nextTick(() => {})
})

// 根据路由中的 id 参数获取文章信息
async function getArticleInfo() {
  const id = route.params.id // 路由中获取参数
  if (!id) return
  window.$loadingBar?.start()
  $message.loading('加载中...')
  try {
    const res = await getArticleById(id)
    formModel.value = res.data
    // 判断图片是本地上传或网络资源, 本地上传一般用于开发时测试, 基本都是用网络资源
    if (!formModel.value.img.startsWith('http')) {
      formModel.value.img = 'http://localhost:8765/' + formModel.value.img
    }

    // 处理图片显示
    res.data.img &&
      (imgList.value[0] = {
        id: res.data.img,
        name: 'img',
        status: 'finished',
        url: res.data.img,
      })
    window.$loadingBar?.finish()
    $message?.success('加载成功')
  } catch (err) {
    window.$loadingBar?.error()
    $message?.error('加载失败')
  }
}

// FIXME: 完善前端校验
const rules = {
  // category_name: {
  //   required: true,
  //   message: '请选择文章分类',
  //   trigger: ['blur', 'change'],
  // },
  // tag_names: {
  //   required: true,
  //   message: '请选择文章标签',
  //   trigger: ['blur', 'change'],
  // },
  // desc: {
  //   required: true,
  //   message: '请输入文章描述',
  //   trigger: ['input', 'blur'],
  // },
}

const formRef = ref(null)
const formModel = ref({
  status: 1, // 发布形式: 默认公开
})
const btnLoading = ref(false)
const modalVisible = ref(false)

// 添加自定义标签
const newCategoryName = ref('')
function addNewCategoryName() {
  if (!newCategoryName.value.trim()) return
  formModel.value.tag_names.push(newCategoryName.value)
}

// 保存草稿
function handleDraft() {
  // TODO:
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
    if (err) return
    btnLoading.value = true
    $message.loading('正在保存...')
    try {
      console.log(formModel.value)
      await saveOrUpdateArticle(formModel.value)
      modalVisible.value = false
      $message.success('操作成功!')
      // 关闭当前标签, 并跳转回文章列表
      useTagsStore().removeTag(route.path)
      await router.replace({ path: '/article/list', query: { needRefresh: true } })
    } catch (err) {
      console.error(err)
    } finally {
      btnLoading.value = false
    }
  })
}

// 上传图片
const imgList = ref([])
function handleImgUpload({ event }) {
  const respStr = (event?.target).response
  const res = JSON.parse(respStr)
  if (res.code !== 0) {
    $message.error('文件上传失败')
    return
  }
  formModel.value.img = res.data
}
</script>

<style lang="scss">
.md-preview {
  ul,
  ol {
    list-style: revert;
  }
}
</style>
