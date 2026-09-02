<script setup>
import { computed, ref } from 'vue'

import api from '@/api'
import UModal from '@/components/ui/UModal.vue'
import { useAppStore, useUserStore } from '@/store'

const userStore = useUserStore()
const appStore = useAppStore()

const registerFlag = computed({
  get: () => appStore.registerFlag,
  set: val => appStore.setRegisterFlag(val),
})

const loginFlag = computed({
  get: () => appStore.loginFlag,
  set: val => appStore.setLoginFlag(val),
})

const form = ref({
  username: 'test@qq.com',
  password: '11111',
})

// 登录
async function handleLogin() {
  const { username, password } = form.value
  if (!username || !password) {
    window.$message?.warning('请输入用户名和密码')
    return
  }

  const doLogin = async (username, password) => {
    const resp = await api.login({ username, password })
    window.$notify?.success('登录成功!')
    userStore.setToken(resp.data.token)
    // 加载用户信息, 更新 pinia 中信息, 刷新页面
    await userStore.getUserInfo()
    // 清空表单
    form.value = { username: 'test@qq.com', password: '11111' }
    loginFlag.value = false
  }

  doLogin(username, password)
  // 腾讯滑块验证码曾经在这里接入 (index.html 里引入 TCaptcha.js, VITE_USE_CAPTCHA 控制开关),
  // 依赖的 AppID 是原作者的, 已随 assets/config.js 一起删掉, 需要时自己配 env 再接回来。
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
  <UModal v-model="loginFlag" :width="480">
    <div class="mx-2 my-1">
      <div class="mb-4 text-xl font-bold">
        登录
      </div>
      <div class="my-7 space-y-4">
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 用户名 </span>
          <input
            v-model="form.username" required placeholder="用户名"
            class="block w-full border-0 rounded-md p-2 text-main shadow-sm outline-none ring-1 ring-line ring-inset placeholder:text-muted focus:ring-2 focus:ring-emerald"
          >
        </div>
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 密码 </span>
          <input
            v-model="form.password" type="password" placeholder="密码"
            class="block w-full border-0 rounded-md p-2 text-main shadow-sm outline-none ring-1 ring-line ring-inset placeholder:text-muted focus:ring-2 focus:ring-emerald"
          >
        </div>
      </div>
      <div class="my-2 text-center">
        <button class="w-full rounded-lg bg-blue py-2 text-white hover:bg-light-blue" @click="handleLogin">
          登录
        </button>
        <div class="mt-4 flex justify-between">
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
