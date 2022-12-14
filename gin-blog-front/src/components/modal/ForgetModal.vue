<script setup lang="ts">
import { useAppStore } from '@/store'
const appStore = useAppStore()

const forgetFlag = computed({
  get: () => appStore.forgetFlag,
  set: val => appStore.setForgetFlag(val),
})

const formRef = ref({
  username: '',
  password: '',
})

const rules = {}

// 登录
function handleLogin() {
  window?.$message?.info('点击了登录!')
}
// 立即注册
function openRegister() {
  appStore.setRegisterFlag(false)
  appStore.setLoginFlag(true)
}
// 忘记密码
function openLogin() {
  appStore.setForgetFlag(false)
  appStore.setLoginFlag(true)
}
</script>

<template>
  <n-modal
    v-model:show="forgetFlag"
    display-directive="show"
    preset="card"
    title="忘记密码"
    :block-scroll="appStore.isMobile"
    transform-origin="center"
    w-360 px-10
    lg:w-460
  >
    <n-form
      ref="formRef"
      :model="formRef"
      :rules="rules"
      label-placement="left"
      label-width="auto"
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
          w-full py-7 btn
          @click="handleLogin"
        >
          登录
        </button>
        <div mt-25 mb-10 flex justify-between>
          <button @click="openRegister">
            立即注册
          </button>
          <button @click="openLogin">
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
