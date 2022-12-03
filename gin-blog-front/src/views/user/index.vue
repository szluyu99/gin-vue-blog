<script setup lang="ts">
import api from '@/api'
import { useUserStore } from '@/store'

const userStore = useUserStore()

// 修改表单
const form = ref<any>({ avatar: '', nickname: '', website: '', intro: '', email: '' })

// TODO: 不使用延时操作
onMounted(() => {
  setTimeout(() => {
    form.value = { ...userStore.userInfo }
  }, 500)
})

function uploadAvatar() {
  window.$message?.info('头像上传开发中...')
}

function updateUserInfo() {
  try {
    api.updateUser(form.value)
    window.$message?.success('修改成功!')
    userStore.getUserInfo() // 重新获取用户信息
  }
  catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <BannerCard>
    <p text-24 font-bold>
      基本信息
    </p>
    <n-grid x-gap="15" cols="12">
      <n-gi span="4" f-c-c>
        <n-image
          :src="form.avatar"
          width="140"
          preview-disabled
          rounded-full
          @click="uploadAvatar"
        />
      </n-gi>
      <n-gi span="7">
        <v-text-field
          v-model="form.nickname"
          label="昵称"
          placeholder="请输入您的昵称"
          variant="underlined"
          mt-7
        />
        <v-text-field
          v-model="form.website"
          label="个人网站"
          placeholder="请输入个人网站"
          variant="underlined"
          mt-7
        />
        <v-text-field
          v-model="form.intro"
          label="简介"
          placeholder="介绍一下自己吧"
          variant="underlined"
          mt-7
        />
        <v-text-field
          v-model="form.email"
          label="邮箱号"
          disabled
          variant="underlined"
          mt-7
        />
        <v-btn variant="outlined" mt-15 @click="updateUserInfo">
          修改
        </v-btn>
      </n-gi>
    </n-grid>
  </BannerCard>
</template>
