<script setup lang="ts">
import api from '@/api'
import { convertImgUrl } from '@/utils'
import { useUserStore } from '@/store'
const userStore = useUserStore()

const infoForm = $ref<any>({
  avatar: userStore.avatar,
  nickname: userStore.nickname,
  intro: userStore.intro,
  website: userStore.website,
  email: userStore.email,
})

function uploadAvatar() {
  window.$message?.info('头像上传开发中...')
}

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
  <BannerCard label="user" title="个人中心">
    <p text-24 font-bold>
      基本信息
    </p>
    <n-grid x-gap="15" cols="12">
      <n-gi span="4" f-c-c>
        <n-image
          :src="convertImgUrl(infoForm.avatar)"
          width="140"
          preview-disabled
          rounded-full
          @click="uploadAvatar"
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
  </BannerCard>
</template>
