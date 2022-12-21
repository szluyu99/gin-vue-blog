<script setup>
import UploadOne from '@/components/upload/UploadOne.vue'
import { useAppStore } from '@/store'
import api from '@/api'

defineOptions({ name: '网站管理' })

const appStore = useAppStore()
appStore.getBlogInfo()

// 表单数据
const form = $ref(appStore.blogConfig)
const formRef = $ref(null)

async function handleSave() {
  formRef?.validate(async (err) => {
    if (!err) {
      try {
        $loadingBar?.start()
        await api.updateBlogConfig(form)
        $loadingBar?.finish()
        $message.success('博客信息更新成功')
        appStore.getBlogInfo() // 重新加载信息
      }
      catch (err) {
        $loadingBar?.error()
      }
    }
  })
}
</script>

<template>
  <CommonPage :show-header="false">
    <n-tabs type="line" animated>
      <n-tab-pane name="website" tab="网站信息">
        <n-form
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          w-500 mt-15
        >
          <n-form-item label="网站头像" path="website_avatar">
            <UploadOne
              v-model:preview="form.website_avatar"
              :width="120"
            />
          </n-form-item>
          <n-form-item label="网站名称" path="website_name">
            <n-input v-model:value="form.website_name" placeholder="请输入网站名称" />
          </n-form-item>
          <n-form-item label="网站作者" path="website_author">
            <n-input v-model:value="form.website_author" placeholder="请输入网站作者" />
          </n-form-item>
          <n-form-item label="网站简介" path="website_intro">
            <n-input v-model:value="form.website_intro" placeholder="请输入网站简介" />
          </n-form-item>
          <n-form-item label="网站创建日期" path="website_createtime">
            <n-date-picker
              v-model:formatted-value="form.website_createtime"
              value-format="yyyy-MM-dd"
              type="date"
            />
          </n-form-item>
          <n-form-item label="网站公告" path="website_notice">
            <n-input
              v-model:value="form.website_notice"
              type="textarea"
              placeholder="请输入网站公告"
              :autosize="{ minRows: 4, maxRows: 6 }"
            />
          </n-form-item>
          <n-form-item label="网站备案号" path="website_record">
            <n-input v-model:value="form.website_record" placeholder="请输入网站备案号" />
          </n-form-item>
          <!-- TODO: 第三方登录 -->
          <!-- <n-form-item label="第三方登录" path="social_login_list">
            <n-checkbox-group v-model:value="cities">
              <n-space item-style="display: flex;">
                <n-checkbox value="QQ" label="QQ" />
                <n-checkbox value="WeiBo" label="微博" />
                <n-checkbox value="WeChat" label="微信" />
              </n-space>
            </n-checkbox-group>
          </n-form-item> -->
          <n-button type="primary" @click="handleSave">
            确认
          </n-button>
        </n-form>
      </n-tab-pane>
      <n-tab-pane name="contact" tab="社交信息">
        <n-form
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          w-500 mt-15
        >
          <n-form-item label="QQ" path="qq">
            <n-input v-model:value="form.qq" placeholder="请输入 QQ" />
          </n-form-item>
          <n-form-item label="Github" path="github">
            <n-input v-model:value="form.github" placeholder="请输入 Github" />
          </n-form-item>
          <n-form-item label="Gitee" path="gitee">
            <n-input v-model:value="form.gitee" placeholder="请输入 Gitee" />
          </n-form-item>
          <n-button type="primary" @click="handleSave">
            确认
          </n-button>
        </n-form>
      </n-tab-pane>
      <n-tab-pane name="other" tab="其他设置">
        <n-form
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          mt-15
        >
          <n-form ref="formRef" label-align="left" :label-width="120" :model="form" inline>
            <n-form-item label="用户头像" path="user_avatar">
              <UploadOne
                v-model:preview="form.user_avatar"
                :width="120"
              />
            </n-form-item>
            <n-form-item label="游客头像" path="tourist_avatar">
              <UploadOne
                v-model:preview="form.tourist_avatar"
                :width="120"
              />
            </n-form-item>
            <!-- <n-form-item label="微信收款码" path="tourist_avatar">
              <n-image border-dashed border-1 text-gray width="120" :src="form.tourist_avatar" />
            </n-form-item>
            <n-form-item label="支付宝收款码" path="tourist_avatar">
              <n-image border-dashed border-1 text-gray width="120" :src="form.tourist_avatar" />
            </n-form-item> -->
          </n-form>
          <n-form-item label-placement="top" label="文章默认封面" path="article_cover">
            <UploadOne
              v-model:preview="form.article_cover"
              :width="300"
            />
          </n-form-item>
          <n-form-item label="评论默认审核" path="is_comment_review">
            <n-radio-group v-model:value="form.is_comment_review" name="is_comment_review">
              <n-radio :value="0">
                关闭
              </n-radio>
              <n-radio :value="1">
                开启
              </n-radio>
            </n-radio-group>
          </n-form-item>
          <n-form-item label="留言默认审核" path="is_message_review">
            <n-radio-group v-model:value="form.is_message_review" name="is_message_review">
              <n-radio :value="0">
                关闭
              </n-radio>
              <n-radio :value="1">
                开启
              </n-radio>
            </n-radio-group>
          </n-form-item>
          <n-form-item label="邮箱通知" path="is_email_notice">
            <n-radio-group v-model:value="form.is_email_notice" name="is_email_notice">
              <n-radio :value="0">
                关闭
              </n-radio>
              <n-radio :value="1">
                开启
              </n-radio>
            </n-radio-group>
          </n-form-item>
          <n-button type="primary" @click="handleSave">
            确认
          </n-button>
        </n-form>
      </n-tab-pane>
    </n-tabs>
  </CommonPage>
</template>
