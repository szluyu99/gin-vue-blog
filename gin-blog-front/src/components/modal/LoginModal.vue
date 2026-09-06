<script setup>
import { computed, ref, watch } from 'vue'

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
  username: userStore.lastUsername || '',
  password: '',
})

/*
每次打开登录框都重新填一次用户名, 优先级: 刚注册的邮箱 > 上次登录成功的用户名 > 空

密码一律不填也不存。用户名是刚注册 / 登录成功过才有, 没成功过的话
prefillUsername 只活在内存里, 刷新就是空的 —— 不会再出现写死的示例账号。
*/
watch(loginFlag, (open) => {
  if (!open) {
    return
  }
  form.value = {
    username: appStore.prefillUsername || userStore.lastUsername || '',
    password: '',
  }
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
    // 登录成功就记住用户名, 放在 getUserInfo 之前:
    // 拉用户信息失败也不该让这次成功的登录白记一次
    userStore.setLastUsername(username)
    appStore.setPrefillUsername('')
    // 加载用户信息, 更新 pinia 中信息, 刷新页面
    await userStore.getUserInfo()
    // 清空表单: 密码不留
    form.value = { username, password: '' }
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
