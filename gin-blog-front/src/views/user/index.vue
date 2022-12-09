<script setup lang="ts">
import UploadOne from '@/components/upload/UploadOne.vue'

import { useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()

const infoForm = $ref<any>({
  avatar: userStore.avatar,
  nickname: userStore.nickname,
  intro: userStore.intro,
  website: userStore.website,
  email: userStore.email,
})
async function updateUserInfo() {
  try {
    await api.updateUser(infoForm)
    window.$message?.success('修改成功!')
    userStore.getUserInfo() // 重新获取用户信息
  }
  catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <BannerPage label="user" title="个人中心" card>
    <p text-24 font-bold>
      基本信息
    </p>
    <n-grid x-gap="15" cols="12">
      <n-gi span="4" f-c-c>
        <UploadOne
          v-model:preview="infoForm.avatar"
          :width="140"
          @finish="val => (infoForm.avatar = val)"
        />
      </n-gi>
      <n-gi span="7">
        <v-text-field
          v-model="infoForm.nickname"
          label="昵称"
          placeholder="请输入您的昵称"
          variant="underlined"
          mt-7
        />
        <v-text-field
          v-model="infoForm.website"
          label="个人网站"
          placeholder="请输入个人网站"
          variant="underlined"
          mt-7
        />
        <v-text-field
          v-model="infoForm.intro"
          label="简介"
          placeholder="介绍一下自己吧"
          variant="underlined"
          mt-7
        />
        <v-text-field
          v-model="infoForm.email"
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
  </BannerPage>
</template>
