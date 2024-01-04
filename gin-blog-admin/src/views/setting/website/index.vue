<script setup>
import { onMounted, ref } from 'vue'
import { NButton, NDatePicker, NForm, NFormItem, NInput, NRadio, NRadioGroup, NTabPane, NTabs } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import UploadOne from '@/components//UploadOne.vue'

import api from '@/api'

defineOptions({ name: '网站管理' })

const formRef = ref(null)
const form = ref({
  website_avatar: '',
  website_name: '阵、雨的个人博客',
  website_author: '阵、雨',
  website_intro: '往事随风而去',
  website_notice: '博客后端基于 gin、gorm 开发\n博客前端基于 Vue3、TS、NaiveUI 开发\n努力开发中...冲冲冲！加油！',
  website_createtime: '2023-12-27 22:40:22',
  website_record: '鲁ICP备2022040119号',
  qq: '123456789',
  github: 'https://github.com/szluyu99',
  gitee: 'https://gitee.com/szluyu99',
  tourist_avatar: 'https://cdn.hahacode.cn/16815451239215dc82548dcadcd578a5bbc8d5deaa.jpg',
  user_avatar: 'https://cdn.hahacode.cn/2299fc4d14c94e6183b082973b35855d.png',
  article_cover: 'https://cdn.hahacode.cn/1679461519cc592408198d67faf1290ff8969dc614.png',
  is_comment_review: 1,
  is_message_review: 1,
  // is_email_notice: 0,
  // social_login_list: [],
  // social_url_list: [],
  // is_reward: 0,
  // wechat_qrcode: 'http://dummyimage.com/100x100',
  // alipay_ode: 'http://dummyimage.com/100x100',
})

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  const resp = await api.getConfig()
  form.value = resp.data
}

function handleSave() {
  formRef.value?.validate(async (err) => {
    if (!err) {
      try {
        $loadingBar?.start()
        await api.updateConfig(form.value)
        $loadingBar?.finish()
        $message.success('博客信息更新成功')
        // fetchData()
      }
      catch (err) {
        $loadingBar?.error()
      }
    }
  })
}
</script>

<template>
  <CommonPage :show-header="false" show-footer>
    <NTabs type="line" animated>
      <NTabPane name="website" tab="网站信息">
        <NForm
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          class="mt-4 w-[500px]"
        >
          <NFormItem label="网站头像" path="website_avatar">
            <UploadOne
              v-model:preview="form.website_avatar"
              :width="120"
            />
          </NFormItem>
          <NFormItem label="网站名称" path="website_name">
            <NInput v-model:value="form.website_name" placeholder="请输入网站名称" />
          </NFormItem>
          <NFormItem label="网站作者" path="website_author">
            <NInput v-model:value="form.website_author" placeholder="请输入网站作者" />
          </NFormItem>
          <NFormItem label="网站简介" path="website_intro">
            <NInput v-model:value="form.website_intro" placeholder="请输入网站简介" />
          </NFormItem>
          <NFormItem label="网站创建日期" path="website_createtime">
            <NDatePicker
              v-model:formatted-value="form.website_createtime"
              value-format="yyyy-MM-dd HH:mm:ss"
              type="datetime"
            />
          </NFormItem>
          <NFormItem label="网站公告" path="website_notice">
            <NInput
              v-model:value="form.website_notice"
              type="textarea"
              placeholder="请输入网站公告"
              :autosize="{ minRows: 4, maxRows: 6 }"
            />
          </NFormItem>
          <NFormItem label="网站备案号" path="website_record">
            <NInput v-model:value="form.website_record" placeholder="请输入网站备案号" />
          </NFormItem>
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
          <NButton type="primary" @click="handleSave">
            确认
          </NButton>
        </NForm>
      </NTabPane>
      <NTabPane name="contact" tab="社交信息">
        <NForm
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          class="mt-4 w-[500px]"
        >
          <NFormItem label="QQ" path="qq">
            <NInput v-model:value="form.qq" placeholder="请输入 QQ" />
          </NFormItem>
          <NFormItem label="Github" path="github">
            <NInput v-model:value="form.github" placeholder="请输入 Github" />
          </NFormItem>
          <NFormItem label="Gitee" path="gitee">
            <NInput v-model:value="form.gitee" placeholder="请输入 Gitee" />
          </NFormItem>
          <NButton type="primary" @click="handleSave">
            确认
          </NButton>
        </NForm>
      </NTabPane>
      <NTabPane name="other" tab="其他设置">
        <NForm
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          class="mt-4"
        >
          <NForm ref="formRef" label-align="left" :label-width="120" :model="form" inline>
            <NFormItem label="用户头像" path="user_avatar">
              <UploadOne
                v-model:preview="form.user_avatar"
                :width="120"
              />
            </NFormItem>
            <NFormItem label="游客头像" path="tourist_avatar">
              <UploadOne
                v-model:preview="form.tourist_avatar"
                :width="120"
              />
            </NFormItem>
            <!-- <n-form-item label="微信收款码" path="tourist_avatar">
              <n-image border-dashed border-1 text-gray width="120" :src="form.tourist_avatar" />
            </n-form-item>
            <n-form-item label="支付宝收款码" path="tourist_avatar">
              <n-image border-dashed border-1 text-gray width="120" :src="form.tourist_avatar" />
            </n-form-item> -->
          </NForm>
          <NFormItem label-placement="top" label="文章默认封面" path="article_cover">
            <UploadOne
              v-model:preview="form.article_cover"
              :width="300"
            />
          </NFormItem>
          <NFormItem label="评论默认审核" path="is_comment_review">
            <NRadioGroup v-model:value="form.is_comment_review" name="is_comment_review">
              <NRadio value="true">
                关闭
              </NRadio>
              <NRadio value="false">
                开启
              </NRadio>
            </NRadioGroup>
          </NFormItem>
          <NFormItem label="留言默认审核" path="is_message_review">
            <NRadioGroup v-model:value="form.is_message_review" name="is_message_review">
              <NRadio value="true">
                关闭
              </NRadio>
              <NRadio value="false">
                开启
              </NRadio>
            </NRadioGroup>
          </NFormItem>
          <!-- <NFormItem label="邮箱通知" path="is_email_notice">
            <NRadioGroup v-model:value="form.is_email_notice" name="is_email_notice">
              <NRadio :value="0">
                关闭
              </NRadio>
              <NRadio :value="1">
                开启
              </NRadio>
            </NRadioGroup>
          </NFormItem> -->
          <NButton type="primary" @click="handleSave">
            确认
          </NButton>
        </NForm>
      </NTabPane>
    </NTabs>
  </CommonPage>
</template>
