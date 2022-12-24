<script setup lang="ts">
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

interface LoginInfo {
  username: string
  password: string
}

let formModel = $ref<LoginInfo>({
  username: '',
  password: '',
})

const rules = {}

// 登录
async function handleLogin() {
  const { username, password } = formModel
  if (!username || !password) {
    window.$message?.warning('请输入用户名和密码')
    return
  }

  const doLogin = async (username: string, password: string) => {
    try {
      const res: any = await api.login({ username, password })
      window.$notification?.success({ title: '登录成功!', duration: 1500 })
      setToken(res.data.token) // 保存在本地
      // 加载用户信息, 更新 pinia 中信息, 刷新页面
      await userStore.getUserInfo()
      // 清空表单
      formModel = { username: '', password: '' }
    }
    finally {
      loginFlag.value = false
    }
  }

  if (JSON.parse(import.meta.env.VITE_USE_CAPTCHA)) {
  // 腾讯滑块验证码 (在 index.html 中引入 js 文件)
    const captcha = new (window as any).TencentCaptcha(config.TENCENT_CAPTCHA,
      async (res: any) => {
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
  <n-modal
    v-model:show="loginFlag"
    display-directive="show"
    preset="card"
    title="登录"
    :block-scroll="appStore.isMobile"
    transform-origin="center"
    px-10 w-360
    lg="w-460"
  >
    <n-form
      :model="formModel"
      :rules="rules"
      label-placement="left"
      label-width="70"
      require-mark-placement="right-hanging"
    >
      <n-form-item label="用户名" path="username">
        <n-input
          v-model:value="formModel.username"
          placeholder="用户名"
          size="large"
          clearable
        />
      </n-form-item>
      <n-form-item label="密码" path="password">
        <n-input
          v-model:value="formModel.password"
          type="password"
          show-password-on="click"
          placeholder="密码"
          size="large"
          clearable
        />
      </n-form-item>
    </n-form>
    <template #footer>
      <div text-center mt="-15">
        <button
          w-full py-7 rounded-1rem bg-blue text-white hover:bg-light-blue
          @click="handleLogin"
        >
          登录
        </button>
        <div mt-25 mb-10 flex justify-between>
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
    </template>
  </n-modal>
</template>
