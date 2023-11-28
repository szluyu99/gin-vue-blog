<script setup>
import { computed, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NModal } from 'naive-ui'
import config from '@/assets/js/config'

import api from '@/api'
import { useAppStore } from '@/store'
const appStore = useAppStore()

const registerFlag = computed({
  get: () => appStore.registerFlag,
  set: val => appStore.setRegisterFlag(val),
})

const form = ref({
  username: '',
  password: '',
  code: '',
})

const rules = {}

// 注册
async function handleRegister() {
  const { username, password, code } = form.value
  if (!code) {
    window.$message?.warning('请输入发送到邮箱的验证码')
    return
  }

  if (!username || !password) {
    window.$message?.warning('请输入邮箱号和密码')
    return
  }

  // 腾讯滑块验证码 (在 index.html 中引入 js 文件)
  const doRegister = async () => {
    // 注册
    await api.register(form.value)
    window.$notification?.success({ title: '注册成功!', duration: 1500 })
    form.value = { username: '', password: '', code: '' }
    // 打开登录弹窗
    openLogin()
  }

  if (JSON.parse(import.meta.env.VITE_USE_CAPTCHA)) {
    const captcha = new window.TencentCaptcha(config.TENCENT_CAPTCHA, async (res) => {
      res.ret === 0 && doRegister()
    })
    captcha.show()
  }
  else {
    doRegister()
  }
}

// 发送验证码
async function sendCode() {
  const reg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/
  if (!reg.test(form.value.username)) {
    window.$message?.warning('请输入正确的邮箱格式')
    return
  }
  await api.sendCode({ email: form.value.username })
  window.$message?.success('邮件发送成功')
}

// 登录
function openLogin() {
  appStore.setRegisterFlag(false)
  appStore.setLoginFlag(true)
}
</script>

<template>
  <NModal
    v-model:show="registerFlag"
    display-directive="show"
    preset="card"
    title="注册"
    :block-scroll="appStore.isMobile"
    transform-origin="center"
    class="w-360 px-10 lg:w-460"
  >
    <NForm
      :model="form"
      :rules="rules"
      label-placement="left"
      label-width="70"
      require-mark-placement="right-hanging"
    >
      <NFormItem label="邮箱号" path="username">
        <NInput
          v-model:value="form.username"
          placeholder="邮箱号，也是用户名"
          clearable
        />
      </NFormItem>
      <NFormItem label="验证码" path="code">
        <NInput
          v-model:value="form.code"
          placeholder="请输入 6 位验证码"
          clearable
        />
        <NButton text class="w-50" @click="sendCode">
          发送
        </NButton>
      </NFormItem>
      <NFormItem label="密码" path="password">
        <NInput
          v-model:value="form.password"
          type="password"
          show-password-on="click"
          placeholder="密码"
          clearable
        />
      </NFormItem>
    </NForm>
    <template #footer>
      <div class="px-15 text-center -mt-15">
        <button
          class="w-full rounded-1rem bg-red py-7 text-white hover:bg-orange"
          @click="handleRegister"
        >
          注册
        </button>
        <div class="mb-10 mt-25 text-left">
          已有账号？
          <NButton text @click="openLogin">
            登录
          </NButton>
        </div>
      <!-- TODO: 第三方登录 -->
      <!-- <div text-center text-10 color="#aaa">
          社交帐号登录
        </div> -->
      </div>
    </template>
  </NModal>
</template>
