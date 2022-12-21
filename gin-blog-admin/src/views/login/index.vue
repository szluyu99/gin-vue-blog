<script setup>
import { useStorage } from '@vueuse/core'
import bgImg from '@/assets/images/login_bg.webp'
import { lStorage, setToken } from '@/utils'
import { addDynamicRoutes } from '@/router'
import { useAppStore, useUserStore } from '@/store'
import config from '@/constant/config'
import api from '@/api'

const title = import.meta.env.VITE_TITLE // 环境变量中读取

const [userStore, appStore, router] = [useUserStore(), useAppStore(), useRouter()]
const { query } = useRoute()

let loginInfo = $ref({
  username: 'test@qq.com',
  password: '11111',
})

initLoginInfo()

// 从 localStorage 中获取记住的用户名和密码
function initLoginInfo() {
  const localLoginInfo = lStorage.get('loginInfo')
  if (localLoginInfo) {
    loginInfo = {
      username: localLoginInfo.username || 'test@qq.com',
      password: localLoginInfo.password || '11111',
    }
  }
}

// Reactive LocalStorage/SessionStorage - vueuse
const isRemember = $(useStorage('isRemember', false))
let loading = $ref(false)

async function handleLogin() {
  const { username, password } = loginInfo
  if (!username || !password) {
    $message.warning('请输入用户名和密码')
    return
  }

  const doLogin = async (username, password) => {
    loading = true
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
      loading = false
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
  <AppPage show-footer bg-cover :style="{ backgroundImage: `url(${bgImg})` }">
    <div
      style="transform: translateY(25px)"
      f-c-c min-w-345 max-w-700 m-auto p-15 rounded-10 card-shadow bg-white bg-opacity-60
    >
      <div w-380 hidden md:block px-20 py-35>
        <img src="@/assets/images/login_banner.webp" w-full alt="login_banner">
      </div>

      <div w-320 flex-col px-20 py-35>
        <h5 f-c-c text-24 font-normal color="#6a6a6a">
          <icon-custom-logo mr-10 text-50 color-primary />
          {{ title }}
        </h5>
        <div mt-30>
          <n-input
            v-model:value="loginInfo.username"
            text-16 items-center h-50 pl-10
            autofocus
            placeholder="test@qq.com"
            :maxlength="20"
          />
        </div>
        <div mt-30>
          <n-input
            v-model:value="loginInfo.password"
            text-16 items-center h-50 pl-10
            type="password"
            show-password-on="mousedown"
            placeholder="11111"
            :maxlength="20"
            @keydown.enter="handleLogin"
          />
        </div>

        <div mt-20>
          <n-checkbox
            :checked="isRemember"
            label="记住我"
            :on-update:checked="(val) => (isRemember = val)"
          />
        </div>

        <div mt-20>
          <n-button
            w-full h-50 rounded-5 text-16
            type="primary"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </n-button>
        </div>
      </div>
    </div>
  </AppPage>
</template>

