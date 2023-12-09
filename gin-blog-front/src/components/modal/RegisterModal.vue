<script setup>
import { computed, ref } from 'vue'
import config from '@/assets/config'

import UModal from '@/components/ui/UModal.vue'
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
    window.$notify?.success('注册成功!')
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
  <UModal v-model="registerFlag" :width="480">
    <div class="mx-2 my-1">
      <div class="mb-4 text-xl font-bold">
        注册
      </div>
      <div class="my-7 space-y-4">
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 邮箱号 </span>
          <input
            v-model="form.username" required placeholder="邮箱号，也是用户名"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 验证码 </span>
          <input
            v-model="form.code" required placeholder="请输入 6 位验证码"
            class="block w-5/6 border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
          <button class="w-1/6" @click="sendCode">
            发送
          </button>
        </div>
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 密码 </span>
          <input
            v-model="form.password" required type="password" placeholder="密码"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
      </div>
      <div class="my-2 text-center">
        <button
          class="w-full rounded bg-red py-2 text-white hover:bg-orange"
          @click="handleRegister"
        >
          注册
        </button>
        <div class="mb-2 mt-6 text-left">
          已有账号？
          <button class="duration-300 hover:text-emerald" @click="openLogin">
            登录
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
