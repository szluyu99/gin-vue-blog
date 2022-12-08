<script setup lang="ts">
import config from '@/assets/js/config'

import api from '@/api'
import { setToken } from '@/utils'

import { useAppStore, useUserStore } from '@/store'
const [appStore, userStore] = [useAppStore(), useUserStore()]

const loginFlag = computed({
  get: () => appStore.loginFlag,
  set: val => appStore.setLoginFlag(val),
})

const registerFlag = computed({
  get: () => appStore.registerFlag,
  set: val => appStore.setRegisterFlag(val),
})

interface LoginInfo {
  username: string
  password: string
}

const formModel = ref<LoginInfo>({
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
  // 腾讯滑块验证码 (在 index.html 中引入 js 文件)
  const captcha = new (window as any).TencentCaptcha(config.TENCENT_CAPTCHA, async (res: any) => {
    if (res.ret === 0) {
      // 登录
      const res: any = await api.login(formModel.value)
      window.$notification?.success({ title: '登录成功!', duration: 1500 })
      setToken(res.data.token) // 保存在本地
      // 加载用户信息, 更新 pinia 中信息, 刷新页面
      await userStore.getUserInfo()
      // 清空表单
      formModel.value = { username: '', password: '' }

      loginFlag.value = false
    }
  })
  captcha.show()
}

// 立即注册
function openRegister() {
  loginFlag.value = false
  registerFlag.value = true
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
    :block-scroll="false"
    transform-origin="center"
    w-460 px-10
  >
    <!-- <v-form ref="form">
      <v-text-field
        v-model="formModel.username"
        label="邮箱号"
        clearable
        variant="underlined"
        required
      />
      <v-text-field
        v-model="formModel.password"
        type="password"
        label="密码"
        clearable
        variant="underlined"
        required
      />
    </v-form> -->
    <n-form
      :model="formModel"
      :rules="rules"
      label-placement="left"
      label-width="auto"
      require-mark-placement="right-hanging"
    >
      <n-form-item label="用户名" path="username">
        <n-input v-model:value="formModel.username" placeholder="用户名" clearable />
      </n-form-item>
      <n-form-item label="密码" path="password">
        <n-input
          v-model:value="formModel.password"
          type="password"
          show-password-on="click"
          placeholder="密码"
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
