<script setup>
import { computed, ref } from 'vue'
import { NForm, NFormItem, NInput } from 'naive-ui'

import UModal from '@/components/ui/UModal.vue'
import config from '@/assets/js/config'

import api from '@/api'
import { setToken } from '@/utils'

import { useAppStore, useUserStore } from '@/store'
const [appStore, userStore] = [useAppStore(), useUserStore()]

const registerFlag = computed({
  get: () => appStore.registerFlag,
  set: val => appStore.setRegisterFlag(val),
})

const loginFlag = computed({
  get: () => appStore.loginFlag,
  set: val => appStore.setLoginFlag(val),
})

const formModel = ref({
  username: '',
  password: '',
})

const rules = {}

// 登录
async function handleLogin() {
  const { username, password } = formModel.value
  if (!username || !password) {
    window.$message?.warning('请输入用户名和密码')
    return
  }

  const doLogin = async (username, password) => {
    try {
      const res = await api.login({ username, password })
      window.$notification?.success({ title: '登录成功!', duration: 1500 })
      setToken(res.data.token) // 保存在本地
      // 加载用户信息, 更新 pinia 中信息, 刷新页面
      await userStore.getUserInfo()
      // 清空表单
      formModel.value = { username: '', password: '' }
    }
    finally {
      loginFlag.value = false
    }
  }

  if (JSON.parse(import.meta.env.VITE_USE_CAPTCHA)) {
  // 腾讯滑块验证码 (在 index.html 中引入 js 文件)
    const captcha = new window.TencentCaptcha(config.TENCENT_CAPTCHA,
      async (res) => {
        res.ret === 0 && doLogin(username, password)
      })
    captcha.show()
  }
  else {
    doLogin(username, password)
  }
}

// 立即注册
function openRegister() {
  registerFlag.value = true
  loginFlag.value = false
}

// TODO:忘记密码
function openForget() {
  window?.$message?.info('暂时不支持找回密码!')
}
</script>

<template>
  <UModal v-model="loginFlag" :width="500">
    <div class="mx-0 my-5">
      <div class="mb-15 text-24 font-bold">
        登录
      </div>
      <NForm
        :model="formModel"
        :rules="rules"
        size="medium"
        label-placement="left"
        label-width="70"
        require-mark-placement="right-hanging"
      >
        <NFormItem label="用户名" path="username">
          <NInput
            v-model:value="formModel.username"
            placeholder="用户名"
            clearable
          />
        </NFormItem>
        <NFormItem label="密码" path="password">
          <NInput
            v-model:value="formModel.password"
            type="password"
            show-password-on="click"
            placeholder="密码"
            clearable
          />
        </NFormItem>
      </NForm>
      <div class="my-10 px-15 text-center">
        <button class="w-full rounded-1rem bg-blue py-7 text-white hover:bg-light-blue" @click="handleLogin">
          登录
        </button>
        <div class="mt-15 flex justify-between">
          <button @click="openRegister">
            立即注册
          </button>
          <button @click="openForget">
            忘记密码？
          </button>
        </div>
      <!-- TODO: 第三方登录 -->
      <!-- <div text-center text-10 color="#aaa">
          社交帐号登录
        </div> -->
      </div>
    </div>
  </UModal>
</template>
