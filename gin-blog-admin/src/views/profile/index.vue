<script setup>
import { onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NTabPane, NTabs } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import UploadOne from '@/components//UploadOne.vue'
import { useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()

const infoFormRef = ref(null)
const infoForm = ref({
  avatar: userStore.avatar,
  nickname: userStore.nickname,
  intro: userStore.intro,
  website: userStore.website,
})

onMounted(async () => {
  await userStore.getUserInfo()
  infoForm.value = {
    avatar: userStore.avatar,
    nickname: userStore.nickname,
    intro: userStore.intro,
    website: userStore.website,
  }
})

async function updateProfile() {
  infoFormRef.value?.validate(async (err) => {
    if (!err) {
      await api.updateCurrent(infoForm.value)
      $message.success('更新成功!')
      userStore.getUserInfo()
    }
  })
}
const infoFormRules = {
  nickname: [
    {
      required: true,
      message: '请输入昵称',
      trigger: ['input', 'blur', 'change'],
    },
  ],
}

// 修改密码的表单
const passwordFormRef = ref(null)
const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

function updatePassword() {
  passwordFormRef.value?.validate(async (err) => {
    if (!err) {
      await api.updateCurrentPassword(passwordForm.value)
      $message.success('修改成功!')
    }
  })
}
const passwordFormRules = {
  old_password: [
    {
      required: true,
      message: '请输入旧密码',
      trigger: ['input', 'blur', 'change'],
    },
  ],
  new_password: [
    {
      required: true,
      message: '请输入新密码',
      trigger: ['input', 'blur', 'change'],
    },
  ],
  confirm_password: [
    {
      required: true,
      message: '请再次输入密码',
      trigger: ['input', 'blur'],
    },
    {
      validator: validatePasswordStartWith,
      message: '两次密码输入不一致',
      trigger: 'input',
    },
    {
      validator: validatePasswordSame,
      message: '两次密码输入不一致',
      trigger: ['blur', 'password-input'],
    },
  ],
}
function validatePasswordStartWith(rule, value) {
  return !!passwordForm.value.new_password && passwordForm.value.new_password.startsWith(value) && passwordForm.value.new_password.length >= value.length
}
function validatePasswordSame(rule, value) {
  return value === passwordForm.value.new_password
}
</script>

<template>
  <CommonPage :show-header="false">
    <NTabs type="line" animated>
      <NTabPane name="website" tab="修改信息">
        <div class="m-7 flex items-center">
          <div class="mr-7 w-50">
            <UploadOne
              v-model:preview="infoForm.avatar"
              :width="130"
            />
          </div>
          <NForm
            ref="infoFormRef"
            label-placement="left"
            label-align="left"
            label-width="100"
            :model="infoForm"
            :rules="infoFormRules"
            class="w-80"
          >
            <NFormItem label="昵称" path="nickname">
              <NInput
                v-model:value="infoForm.nickname"
                type="text"
                placeholder="请填写昵称"
              />
            </NFormItem>
            <NFormItem label="个人简介" path="intro">
              <NInput
                v-model:value="infoForm.intro"
                type="text"
                placeholder="请填写个人简介"
              />
            </NFormItem>
            <NFormItem label="个人网站" path="website">
              <NInput
                v-model:value="infoForm.website"
                type="text"
                placeholder="请填写个人网站"
              />
            </NFormItem>
            <NButton type="primary" @click="updateProfile">
              修改
            </NButton>
          </NForm>
        </div>
      </NTabPane>
      <NTabPane name="contact" tab="修改密码">
        <NForm
          ref="passwordFormRef"
          label-placement="left"
          label-align="left"
          :model="passwordForm"
          label-width="100"
          :rules="passwordFormRules"
          class="m-[30px] w-[400px]"
        >
          <NFormItem label="旧密码" path="old_password">
            <NInput
              v-model:value="passwordForm.old_password"
              type="password"
              show-password-on="mousedown"
              placeholder="请输入旧密码"
            />
          </NFormItem>
          <NFormItem label="新密码" path="new_password">
            <NInput
              v-model:value="passwordForm.new_password"
              :disabled="!passwordForm.old_password"
              type="password"
              show-password-on="mousedown"
              placeholder="请输入新密码"
            />
          </NFormItem>
          <NFormItem label="确认密码" path="confirm_password">
            <NInput
              v-model:value="passwordForm.confirm_password"
              :disabled="!passwordForm.new_password"
              type="password"
              show-password-on="mousedown"
              placeholder="请再次输入新密码"
            />
          </NFormItem>
          <NButton type="primary" @click="updatePassword">
            修改
          </NButton>
        </NForm>
      </NTabPane>
    </NTabs>
  </CommonPage>
</template>
