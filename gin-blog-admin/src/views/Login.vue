<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStorage } from '@vueuse/core'
import { NButton, NCheckbox, NInput } from 'naive-ui'

import AppPage from '@/components/page/AppPage.vue'
import { lStorage, setToken } from '@/utils'
import { addDynamicRoutes } from '@/router'
import { useAppStore, useUserStore } from '@/store'
import bgImg from '@/assets/images/login_bg.webp'
import config from '@/constant/config'
import api from '@/api'

const title = import.meta.env.VITE_TITLE // 环境变量中读取

const userStore = useUserStore()
const appStore = useAppStore()
const router = useRouter()
const { query } = useRoute()

const loginInfo = ref({
  username: 'test@qq.com',
  password: '11111',
})

initLoginInfo()

// 从 localStorage 中获取记住的用户名和密码
function initLoginInfo() {
  const localLoginInfo = lStorage.get('loginInfo')
  if (localLoginInfo) {
    loginInfo.value = {
      username: localLoginInfo.username || 'test@qq.com',
      password: localLoginInfo.password || '11111',
    }
  }
}

// Reactive LocalStorage/SessionStorage - vueuse
const isRemember = useStorage('isRemember', false)
const loading = ref(false)

async function handleLogin() {
  const { username, password } = loginInfo.value
  if (!username || !password) {
    $message.warning('请输入用户名和密码')
    return
  }

  const doLogin = async (username, password) => {
    loading.value = true
    $message.loading('正在验证...')

    // 登录接口
    try {
      const res = await api.login({ username, password })
      setToken(res.data.token) // 持久化 token
      $message.success('登录成功')

      await userStore.getUserInfo() // 获取用户信息
      await appStore.getBlogInfo() // 获取博客信息

      // "记住我" 功能
      isRemember
        ? lStorage.set('loginInfo', { username, password })
        : lStorage.remove('loginInfo')

      // 动态添加路由
      await addDynamicRoutes()

      // 页面跳转: 根据 URL 中的 redirect 进行跳转
      if (query.redirect) {
        const path = query.redirect
        Reflect.deleteProperty(query, 'redirect') // 从对象身上删除属性
        router.push({ path, query })
      }
      else {
        router.push('/')
      }
    }
    finally {
      loading.value = false
    }
  }

  // 判断是否需要验证码
  if (JSON.parse(import.meta.env.VITE_USE_CAPTCHA)) {
    // 腾讯滑块验证码 (在 index.html 中引入 js 文件)
    const captcha = new TencentCaptcha(config.TENCENT_CAPTCHA,
      async res => res.ret === 0 && doLogin(username, password),
    )
    captcha.show()
  }
  else {
    doLogin(username, password)
  }
}
</script>

<template>
  <AppPage class="bg-cover" :style="{ backgroundImage: `url(${bgImg})` }">
    <div
      style="transform: translateY(25px)"
      class="m-auto max-w-700 min-w-345 flex items-center justify-center rounded-10 bg-white bg-opacity-60 p-15 shadow"
    >
      <div class="hidden w-380 px-20 py-35 md:block">
        <img src="@/assets/images/login_banner.webp" class="w-full" alt="login_banner">
      </div>

      <div class="w-320 flex flex-col px-20 py-35">
        <h5 class="flex items-center justify-center text-24 font-normal text-gray">
          <img src="@/assets/svg/logo.svg" alt="logo" class="mr-10 h-50 w-50">
          {{ title }}
        </h5>
        <div class="mt-30">
          <NInput
            v-model:value="loginInfo.username"
            class="autofocus h-50 items-center pl-10 text-16"
            placeholder="test@qq.com"
            :maxlength="20"
          />
        </div>
        <div class="mt-30">
          <NInput
            v-model:value="loginInfo.password"
            class="h-50 items-center pl-10 text-16"
            type="password"
            show-password-on="mousedown"
            placeholder="11111"
            :maxlength="20"
            @keydown.enter="handleLogin"
          />
        </div>

        <div class="mt-20">
          <NCheckbox
            :checked="isRemember"
            label="记住我"
            :on-update:checked="(val) => (isRemember = val)"
          />
        </div>

        <div class="mt-20">
          <NButton
            class="h-50 w-full rounded-5 text-16"
            type="primary"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </NButton>
        </div>
      </div>
    </div>
  </AppPage>
</template>
