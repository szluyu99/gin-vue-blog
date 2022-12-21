<script setup>
import UploadOne from '@/components/upload/UploadOne.vue'
import api from '@/api'
import { useUserStore } from '@/store'
const userStore = useUserStore()

// 用户信息的表单
const infoFormRef = $ref(null)
const infoForm = $ref({
  avatar: userStore.avatar,
  nickname: userStore.nickname,
  intro: userStore.intro,
  website: userStore.website,
})
async function updateProfile() {
  infoFormRef?.validate(async (err) => {
    if (!err) {
      await api.updateCurrent(infoForm)
      $message.success('更新成功!')
      userStore.setUserInfo(infoForm)
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
const passwordFormRef = $ref(null)
const passwordForm = $ref({
  old_password: '',
  new_password: '',
  confirm_password: '',
})
async function updatePassword() {
  passwordFormRef?.validate(async (err) => {
    if (!err) {
      await api.updateCurrentPassword(passwordForm)
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
  return !!passwordForm.new_password && passwordForm.new_password.startsWith(value) && passwordForm.new_password.length >= value.length
}
function validatePasswordSame(rule, value) {
  return value === passwordForm.new_password
}
</script>

<template>
  <CommonPage :show-header="false">
    <n-tabs type="line" animated>
      <n-tab-pane name="website" tab="修改信息">
        <div flex items-center m-30>
          <div w-200 mr-30>
            <UploadOne
              v-model:preview="infoForm.avatar"
              :width="130"
            />
          </div>
          <n-form
            ref="infoFormRef"
            label-placement="left"
            label-align="left"
            label-width="100"
            :model="infoForm"
            :rules="infoFormRules"
            w-400
          >
            <n-form-item label="昵称" path="nickname">
              <n-input
                v-model:value="infoForm.nickname"
                type="text"
                placeholder="请填写昵称"
              />
            </n-form-item>
            <n-form-item label="个人简介" path="intro">
              <n-input
                v-model:value="infoForm.intro"
                type="text"
                placeholder="请填写个人简介"
              />
            </n-form-item>
            <n-form-item label="个人网站" path="website">
              <n-input
                v-model:value="infoForm.website"
                type="text"
                placeholder="请填写个人网站"
              />
            </n-form-item>
            <n-button type="primary" @click="updateProfile">
              修改
            </n-button>
          </n-form>
        </div>
      </n-tab-pane>
      <n-tab-pane name="contact" tab="修改密码">
        <n-form
          ref="passwordFormRef"
          label-placement="left"
          label-align="left"
          :model="passwordForm"
          label-width="100"
          :rules="passwordFormRules"
          w-400 m-30
        >
          <n-form-item label="旧密码" path="old_password">
            <n-input
              v-model:value="passwordForm.old_password"
              type="password"
              show-password-on="mousedown"
              placeholder="请输入旧密码"
            />
          </n-form-item>
          <n-form-item label="新密码" path="new_password">
            <n-input
              v-model:value="passwordForm.new_password"
              :disabled="!passwordForm.old_password"
              type="password"
              show-password-on="mousedown"
              placeholder="请输入新密码"
            />
          </n-form-item>
          <n-form-item label="确认密码" path="confirm_password">
            <n-input
              v-model:value="passwordForm.confirm_password"
              :disabled="!passwordForm.new_password"
              type="password"
              show-password-on="mousedown"
              placeholder="请再次输入新密码"
            />
          </n-form-item>
          <n-button type="primary" @click="updatePassword">
            修改
          </n-button>
        </n-form>
      </n-tab-pane>
    </n-tabs>
  </CommonPage>
</template>
