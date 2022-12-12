<script setup lang="ts">
import config from '@/assets/js/config'

import api from '@/api'
import { useAppStore } from '@/store'
const appStore = useAppStore()

const registerFlag = computed({
  get: () => appStore.registerFlag,
  set: val => appStore.setRegisterFlag(val),
})

let formRef = $ref({
  username: '',
  password: '',
})

const rules = {}

// 注册
async function handleRegister() {
  const { username, password } = formRef
  if (!username || !password) {
    window.$message?.warning('请输入用户名和密码')
    return
  }

  // 腾讯滑块验证码 (在 index.html 中引入 js 文件)
  const captcha = new (window as any).TencentCaptcha(config.TENCENT_CAPTCHA, async (res: any) => {
    if (res.ret === 0) {
      // 注册
      await api.register(formRef)
      window.$notification?.success({ title: '注册成功!', duration: 1500 })
      formRef = { username: '', password: '' }
      // 打开登录弹窗
      openLogin()
    }
  })
  captcha.show()
}
// 登录
function openLogin() {
  appStore.setRegisterFlag(false)
  appStore.setLoginFlag(true)
}
</script>

<template>
  <n-modal
    v-model:show="registerFlag"
    display-directive="show"
    preset="card"
    title="注册"
    :block-scroll="appStore.isMobile"
    transform-origin="center"
    px-10 w-370
    lg="w-460"
  >
    <n-form
      :model="formRef"
      :rules="rules"
      label-placement="left"
      label-width="70"
      require-mark-placement="right-hanging"
    >
      <n-form-item label="用户名" path="username">
        <n-input
          v-model:value="formRef.username"
          placeholder="用户名"
          clearable
        />
      </n-form-item>
      <n-form-item label="密码" path="password">
        <n-input
          v-model:value="formRef.password"
          type="password"
          show-password-on="click"
          placeholder="密码"
          clearable
        />
      </n-form-item>
    </n-form>
    <template #footer>
      <div text-center px-15 mt="-15">
        <button
          w-full py-7 rounded-1rem text-white bg-red hover:bg-orange
          @click="handleRegister"
        >
          注册
        </button>
        <div mt-25 mb-10 text-left>
          已有账号？
          <button @click="openLogin">
            登录
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
